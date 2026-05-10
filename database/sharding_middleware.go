package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/XingMenTech/common/logger"
	"github.com/beego/beego/v2/client/orm"
)

// ==================== Beego ORM 通用分库分表中间件 ====================

// ShardMiddlewareConfig 分片中间件配置
type ShardMiddlewareConfig struct {
	ShardConfig          ShardConfig   // 分片配置
	DBPrefix             string        // 数据库前缀
	TablePrefix          string        // 表前缀
	DefaultDBAlias       string        // 默认数据库别名
	ConnectionPoolConfig PoolConfig    // 连接池配置
	EnableQueryCache     bool          // 是否启用查询缓存
	CacheExpire          time.Duration // 缓存过期时间
	MaxRetryCount        int           // 最大重试次数
	RetryInterval        time.Duration // 重试间隔
}

// PoolConfig 连接池配置
type PoolConfig struct {
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int // 秒
	ConnMaxIdleTime int // 秒
}

// ShardMiddleware 分片中间件
type ShardMiddleware struct {
	config     ShardMiddlewareConfig
	router     *ShardRouter
	dbs        map[string]orm.Ormer // 多个数据库连接
	mu         sync.RWMutex
	queryCache sync.Map // 简单的查询缓存
}

// NewShardMiddleware 创建分片中间件
func NewShardMiddleware(cfg ShardMiddlewareConfig) (*ShardMiddleware, error) {
	if err := cfg.ShardConfig.Validate(); err != nil {
		return nil, fmt.Errorf("分片配置验证失败: %w", err)
	}

	// 设置默认值
	if cfg.DBPrefix == "" {
		cfg.DBPrefix = "db"
	}
	if cfg.TablePrefix == "" {
		cfg.TablePrefix = "users"
	}
	if cfg.DefaultDBAlias == "" {
		cfg.DefaultDBAlias = "default"
	}
	if cfg.MaxRetryCount <= 0 {
		cfg.MaxRetryCount = 3
	}
	if cfg.RetryInterval <= 0 {
		cfg.RetryInterval = 100 * time.Millisecond
	}

	router, err := NewShardRouter(cfg.ShardConfig, cfg.DBPrefix, cfg.TablePrefix)
	if err != nil {
		return nil, fmt.Errorf("创建分片路由器失败: %w", err)
	}

	return &ShardMiddleware{
		config: cfg,
		router: router,
		dbs:    make(map[string]*orm.Database),
	}, nil
}

// RegisterDatabase 注册数据库连接
func (sm *ShardMiddleware) RegisterDatabase(alias string, config *MysqlConfig) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// 初始化数据库连接
	if err := InitMysql(config); err != nil {
		return fmt.Errorf("初始化数据库 %s 失败: %w", alias, err)
	}

	// 创建 ORM 实例
	o := orm.NewOrm()
	o.Using(alias)

	sm.dbs[alias] = o
	logger.LOG.Infof("数据库 %s 已注册到分片中间件", alias)

	return nil
}

// GetOrmer 根据用户ID获取对应的 Ormer
func (sm *ShardMiddleware) GetOrmer(userID int64) (orm.Ormer, error) {
	dbName := sm.router.GetDBName(userID)

	sm.mu.RLock()
	o, ok := sm.dbs[dbName]
	sm.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("数据库 %s 未注册", dbName)
	}

	return o, nil
}

// GetTableName 获取表名
func (sm *ShardMiddleware) GetTableName(userID int64) string {
	return sm.router.GetTableName(userID)
}

// GetFullTableName 获取完整表名
func (sm *ShardMiddleware) GetFullTableName(userID int64) string {
	return sm.router.GetFullTableName(userID)
}

// Insert 插入数据（自动路由）
func (sm *ShardMiddleware) Insert(userID int64, data interface{}) (int64, error) {
	var result int64
	var lastErr error

	for retry := 0; retry < sm.config.MaxRetryCount; retry++ {
		o, err := sm.GetOrmer(userID)
		if err != nil {
			lastErr = err
			time.Sleep(sm.config.RetryInterval)
			continue
		}

		tableName := sm.GetTableName(userID)
		id, err := o.QueryTable(tableName).Insert(data)
		if err != nil {
			lastErr = err
			time.Sleep(sm.config.RetryInterval)
			continue
		}

		result = id
		lastErr = nil
		break
	}

	if lastErr != nil {
		return 0, fmt.Errorf("插入数据失败（重试%d次）: %w", sm.config.MaxRetryCount, lastErr)
	}

	return result, nil
}

// InsertMulti 批量插入数据
func (sm *ShardMiddleware) InsertMulti(userID int64, data interface{}, batchSize int) (int64, error) {
	var result int64
	var lastErr error

	for retry := 0; retry < sm.config.MaxRetryCount; retry++ {
		o, err := sm.GetOrmer(userID)
		if err != nil {
			lastErr = err
			time.Sleep(sm.config.RetryInterval)
			continue
		}

		tableName := sm.GetTableName(userID)
		num, err := o.InsertMulti(batchSize, data)
		if err != nil {
			lastErr = err
			time.Sleep(sm.config.RetryInterval)
			continue
		}

		result = num
		lastErr = nil
		break
	}

	if lastErr != nil {
		return 0, fmt.Errorf("批量插入数据失败（重试%d次）: %w", sm.config.MaxRetryCount, lastErr)
	}

	return result, nil
}

// QueryOne 查询单条记录
func (sm *ShardMiddleware) QueryOne(userID int64, to interface{}, conditions ...string) error {
	cacheKey := fmt.Sprintf("query_one:%d:%v", userID, conditions)

	// 检查缓存
	if sm.config.EnableQueryCache {
		if cached, ok := sm.queryCache.Load(cacheKey); ok {
			if cacheData, ok := cached.(*cacheItem); ok {
				if time.Now().Before(cacheData.expireAt) {
					// 复制缓存数据
					if err := copyStruct(cacheData.data, to); err == nil {
						return nil
					}
				} else {
					sm.queryCache.Delete(cacheKey)
				}
			}
		}
	}

	var lastErr error
	for retry := 0; retry < sm.config.MaxRetryCount; retry++ {
		o, err := sm.GetOrmer(userID)
		if err != nil {
			lastErr = err
			time.Sleep(sm.config.RetryInterval)
			continue
		}

		tableName := sm.GetTableName(userID)
		qs := o.QueryTable(tableName)

		// 应用查询条件
		for i := 0; i < len(conditions); i += 2 {
			if i+1 < len(conditions) {
				qs = qs.Filter(conditions[i], conditions[i+1])
			}
		}

		err = qs.One(to)
		if err != nil {
			lastErr = err
			time.Sleep(sm.config.RetryInterval)
			continue
		}

		// 缓存结果
		if sm.config.EnableQueryCache {
			sm.queryCache.Store(cacheKey, &cacheItem{
				data:     to,
				expireAt: time.Now().Add(sm.config.CacheExpire),
			})
		}

		lastErr = nil
		break
	}

	if lastErr != nil {
		return fmt.Errorf("查询数据失败（重试%d次）: %w", sm.config.MaxRetryCount, lastErr)
	}

	return nil
}

// QueryAll 查询多条记录
func (sm *ShardMiddleware) QueryAll(userID int64, to interface{}, conditions ...string) (int64, error) {
	var lastErr error
	var num int64

	for retry := 0; retry < sm.config.MaxRetryCount; retry++ {
		o, err := sm.GetOrmer(userID)
		if err != nil {
			lastErr = err
			time.Sleep(sm.config.RetryInterval)
			continue
		}

		tableName := sm.GetTableName(userID)
		qs := o.QueryTable(tableName)

		// 应用查询条件
		for i := 0; i < len(conditions); i += 2 {
			if i+1 < len(conditions) {
				qs = qs.Filter(conditions[i], conditions[i+1])
			}
		}

		num, err = qs.All(to)
		if err != nil {
			lastErr = err
			time.Sleep(sm.config.RetryInterval)
			continue
		}

		lastErr = nil
		break
	}

	if lastErr != nil {
		return 0, fmt.Errorf("查询数据失败（重试%d次）: %w", sm.config.MaxRetryCount, lastErr)
	}

	return num, nil
}

// Update 更新数据
func (sm *ShardMiddleware) Update(userID int64, data interface{}, cols ...string) (int64, error) {
	var result int64
	var lastErr error

	for retry := 0; retry < sm.config.MaxRetryCount; retry++ {
		o, err := sm.GetOrmer(userID)
		if err != nil {
			lastErr = err
			time.Sleep(sm.config.RetryInterval)
			continue
		}

		tableName := sm.GetTableName(userID)
		num, err := o.QueryTable(tableName).Update(data, cols...)
		if err != nil {
			lastErr = err
			time.Sleep(sm.config.RetryInterval)
			continue
		}

		result = num
		lastErr = nil
		break
	}

	if lastErr != nil {
		return 0, fmt.Errorf("更新数据失败（重试%d次）: %w", sm.config.MaxRetryCount, lastErr)
	}

	return result, nil
}

// Delete 删除数据
func (sm *ShardMiddleware) Delete(userID int64, conditions ...string) (int64, error) {
	var result int64
	var lastErr error

	for retry := 0; retry < sm.config.MaxRetryCount; retry++ {
		o, err := sm.GetOrmer(userID)
		if err != nil {
			lastErr = err
			time.Sleep(sm.config.RetryInterval)
			continue
		}

		tableName := sm.GetTableName(userID)
		qs := o.QueryTable(tableName)

		// 应用删除条件
		for i := 0; i < len(conditions); i += 2 {
			if i+1 < len(conditions) {
				qs = qs.Filter(conditions[i], conditions[i+1])
			}
		}

		num, err := qs.Delete()
		if err != nil {
			lastErr = err
			time.Sleep(sm.config.RetryInterval)
			continue
		}

		result = num
		lastErr = nil
		break
	}

	if lastErr != nil {
		return 0, fmt.Errorf("删除数据失败（重试%d次）: %w", sm.config.MaxRetryCount, lastErr)
	}

	return result, nil
}

// BroadcastQuery 广播查询（跨所有分片）
func (sm *ShardMiddleware) BroadcastQuery(queryFunc func(o orm.Ormer, tableName string) (interface{}, error)) ([]interface{}, error) {
	var results []interface{}
	var mu sync.Mutex
	var wg sync.WaitGroup
	var firstErr error

	sm.mu.RLock()
	dbs := make(map[string]orm.Ormer)
	for k, v := range sm.dbs {
		dbs[k] = v
	}
	sm.mu.RUnlock()

	for dbName := range dbs {
		for tblIdx := 0; tblIdx < sm.config.ShardConfig.TableCount; tblIdx++ {
			wg.Add(1)
			go func(dbName string, tblIdx int) {
				defer wg.Done()

				o := dbs[dbName]

				// 计算数据库索引
				dbIndex := 0
				fmt.Sscanf(dbName, "%*[^_]_%d", &dbIndex)

				tableName := fmt.Sprintf("%s_%d_%d", sm.config.TablePrefix, dbIndex, tblIdx)

				result, err := queryFunc(o, tableName)
				if err != nil {
					mu.Lock()
					if firstErr == nil {
						firstErr = err
					}
					mu.Unlock()
					return
				}

				if result != nil {
					mu.Lock()
					results = append(results, result)
					mu.Unlock()
				}
			}(dbName, tblIdx)
		}
	}

	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}

	return results, nil
}

// ClearCache 清除缓存
func (sm *ShardMiddleware) ClearCache() {
	sm.queryCache.Range(func(key, value interface{}) bool {
		sm.queryCache.Delete(key)
		return true
	})
}

// GetRouter 获取分片路由器
func (sm *ShardMiddleware) GetRouter() *ShardRouter {
	return sm.router
}

// GetStats 获取中间件统计信息
func (sm *ShardMiddleware) GetStats() map[string]interface{} {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	stats := map[string]interface{}{
		"registered_databases": len(sm.dbs),
		"shard_config": map[string]interface{}{
			"db_count":     sm.config.ShardConfig.DBCount,
			"table_count":  sm.config.ShardConfig.TableCount,
			"total_shards": sm.config.ShardConfig.TotalShards(),
		},
		"cache_enabled": sm.config.EnableQueryCache,
		"max_retry":     sm.config.MaxRetryCount,
	}

	return stats
}

// cacheItem 缓存项
type cacheItem struct {
	data     interface{}
	expireAt time.Time
}

// copyStruct 复制结构体数据（简单实现）
func copyStruct(src, dst interface{}) error {
	// 这里可以使用反射或序列化/反序列化来实现深拷贝
	// 简化版本，实际项目中建议使用更完善的实现
	return nil
}

// WithContext 支持 Context 的查询方法
func (sm *ShardMiddleware) QueryOneWithContext(ctx context.Context, userID int64, to interface{}, conditions ...string) error {
	done := make(chan error, 1)

	go func() {
		done <- sm.QueryOne(userID, to, conditions...)
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

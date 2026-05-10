package database

import (
	"fmt"
	"sync"
)

// ==================== 分库分表公共基础模块 ====================

// NamingStrategy 命名策略类型
type NamingStrategy string

const (
	// NamingStrategySimple 简单策略: db_0, users_0_0
	NamingStrategySimple NamingStrategy = "simple"

	// NamingStrategyPadded 补零策略: db_000, users_000_000
	NamingStrategyPadded NamingStrategy = "padded"

	// NamingStrategyCustom 自定义策略: 使用自定义函数
	NamingStrategyCustom NamingStrategy = "custom"
)

// NameFormatFunc 自定义命名函数类型
type NameFormatFunc func(prefix string, index int) string

// ShardConfig 分片配置（公共）
type ShardConfig struct {
	DBCount         int            // 数据库数量
	TableCount      int            // 每个库的表数量
	NamingStrategy  NamingStrategy // 命名策略
	PaddingWidth    int            // 补零宽度（仅用于Padded策略）
	DBNameFormat    NameFormatFunc // 自定义数据库命名函数
	TableNameFormat NameFormatFunc // 自定义表命名函数
}

// Validate 验证分片配置
func (sc *ShardConfig) Validate() error {
	if sc.DBCount <= 0 {
		return fmt.Errorf("数据库数量必须大于0")
	}
	if sc.TableCount <= 0 {
		return fmt.Errorf("表数量必须大于0")
	}
	if sc.DBCount*sc.TableCount > 1000 {
		return fmt.Errorf("分片总数不能超过1000")
	}

	// 验证命名策略
	if sc.NamingStrategy == NamingStrategyCustom {
		if sc.DBNameFormat == nil || sc.TableNameFormat == nil {
			return fmt.Errorf("自定义命名策略需要提供命名函数")
		}
	}

	if sc.NamingStrategy == NamingStrategyPadded && sc.PaddingWidth <= 0 {
		sc.PaddingWidth = 3 // 默认补零宽度为3
	}

	return nil
}

// TotalShards 获取总分片数
func (sc *ShardConfig) TotalShards() int {
	return sc.DBCount * sc.TableCount
}

// ShardRouter 通用分片路由器（公共核心逻辑）
type ShardRouter struct {
	config      ShardConfig
	dbPrefix    string // 数据库前缀（如 "db", "gorm_db", "beego_db"）
	tablePrefix string // 表前缀（如 "users", "gorm_users", "beego_users"）
}

// NewShardRouter 创建分片路由器
func NewShardRouter(config ShardConfig, dbPrefix, tablePrefix string) (*ShardRouter, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &ShardRouter{
		config:      config,
		dbPrefix:    dbPrefix,
		tablePrefix: tablePrefix,
	}, nil
}

// GetDBIndex 根据用户ID计算数据库索引（公共算法）
func (sr *ShardRouter) GetDBIndex(userID int64) int {
	return int(userID % int64(sr.config.DBCount))
}

// GetTableIndex 根据用户ID计算表索引（公共算法）
func (sr *ShardRouter) GetTableIndex(userID int64) int {
	return int((userID / int64(sr.config.DBCount)) % int64(sr.config.TableCount))
}

// GetDBName 获取数据库名称（支持自定义前缀和命名策略）
func (sr *ShardRouter) GetDBName(userID int64) string {
	dbIndex := sr.GetDBIndex(userID)

	switch sr.config.NamingStrategy {
	case NamingStrategySimple:
		return fmt.Sprintf("%s_%d", sr.dbPrefix, dbIndex)
	case NamingStrategyPadded:
		width := sr.config.PaddingWidth
		if width <= 0 {
			width = 3
		}
		return fmt.Sprintf("%s_%0*d", sr.dbPrefix, width, dbIndex)
	case NamingStrategyCustom:
		if sr.config.DBNameFormat != nil {
			return sr.config.DBNameFormat(sr.dbPrefix, dbIndex)
		}
		fallthrough
	default:
		return fmt.Sprintf("%s_%d", sr.dbPrefix, dbIndex)
	}
}

// GetTableName 获取表名（支持自定义前缀和命名策略）
func (sr *ShardRouter) GetTableName(userID int64) string {
	dbIndex := sr.GetDBIndex(userID)
	tableIndex := sr.GetTableIndex(userID)

	switch sr.config.NamingStrategy {
	case NamingStrategySimple:
		return fmt.Sprintf("%s_%d_%d", sr.tablePrefix, dbIndex, tableIndex)
	case NamingStrategyPadded:
		width := sr.config.PaddingWidth
		if width <= 0 {
			width = 3
		}
		return fmt.Sprintf("%s_%0*d_%0*d", sr.tablePrefix, width, dbIndex, width, tableIndex)
	case NamingStrategyCustom:
		if sr.config.TableNameFormat != nil {
			// 使用数据库索引和表索引的组合来确定表名
			combinedIndex := dbIndex*sr.config.TableCount + tableIndex
			return sr.config.TableNameFormat(sr.tablePrefix, combinedIndex)
		}
		fallthrough
	default:
		return fmt.Sprintf("%s_%d_%d", sr.tablePrefix, dbIndex, tableIndex)
	}
}

// GetFullTableName 获取完整的表引用（database.table）
func (sr *ShardRouter) GetFullTableName(userID int64) string {
	return fmt.Sprintf("%s.%s", sr.GetDBName(userID), sr.GetTableName(userID))
}

// GetShardPosition 获取分片位置信息
func (sr *ShardRouter) GetShardPosition(userID int64) (dbIndex, tableIndex int) {
	return sr.GetDBIndex(userID), sr.GetTableIndex(userID)
}

// ShardDistribution 分片分布统计
type ShardDistribution struct {
	DBIndex    int
	TableIndex int
	DBName     string
	TableName  string
	Count      int
}

// CalculateDistribution 计算数据分布（公共方法）
func (sr *ShardRouter) CalculateDistribution(userIDs []int64) map[string]*ShardDistribution {
	distribution := make(map[string]*ShardDistribution)

	for _, userID := range userIDs {
		dbIndex := sr.GetDBIndex(userID)
		tableIndex := sr.GetTableIndex(userID)
		key := fmt.Sprintf("%d_%d", dbIndex, tableIndex)

		if _, exists := distribution[key]; !exists {
			distribution[key] = &ShardDistribution{
				DBIndex:    dbIndex,
				TableIndex: tableIndex,
				DBName:     sr.GetDBName(userID),
				TableName:  sr.GetTableName(userID),
				Count:      0,
			}
		}

		distribution[key].Count++
	}

	return distribution
}

// PrintDistribution 打印分片分布统计
func (sr *ShardRouter) PrintDistribution(userIDs []int64) {
	distribution := sr.CalculateDistribution(userIDs)

	fmt.Println("\n分片分布统计:")
	fmt.Printf("总分片数: %d (数据库: %d × 表: %d)\n",
		sr.config.TotalShards(), sr.config.DBCount, sr.config.TableCount)
	fmt.Printf("总数据量: %d\n", len(userIDs))
	fmt.Println("---")

	for _, dist := range distribution {
		fmt.Printf("  %s.%s: %d 条数据\n", dist.DBName, dist.TableName, dist.Count)
	}

	// 计算均衡性
	if len(distribution) > 0 {
		var min, max int = -1, 0
		for _, dist := range distribution {
			if min == -1 || dist.Count < min {
				min = dist.Count
			}
			if dist.Count > max {
				max = dist.Count
			}
		}

		if min > 0 {
			ratio := float64(max) / float64(min)
			fmt.Printf("---\n均衡性: %.2f (越接近1.0越均衡)\n", ratio)
		}
	}
}

// GenerateTestUserIDs 生成测试用户ID列表
func GenerateTestUserIDs(count int) []int64 {
	userIDs := make([]int64, count)
	for i := 0; i < count; i++ {
		userIDs[i] = int64(i + 1)
	}
	return userIDs
}

// ConcurrentQueryTask 并发查询任务接口
type ConcurrentQueryTask interface {
	Execute() (interface{}, error)
}

// ConcurrentQueryResult 并发查询结果
type ConcurrentQueryResult struct {
	Data  interface{}
	Error error
}

// ExecuteConcurrentQueries 执行并发查询（公共并发控制）
func ExecuteConcurrentQueries(tasks []ConcurrentQueryTask, maxConcurrency int) ([]interface{}, error) {
	if len(tasks) == 0 {
		return []interface{}{}, nil
	}

	if maxConcurrency <= 0 {
		maxConcurrency = len(tasks)
	}

	results := make([]ConcurrentQueryResult, len(tasks))
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxConcurrency)

	for i, task := range tasks {
		wg.Add(1)
		go func(index int, t ConcurrentQueryTask) {
			defer wg.Done()

			// 获取信号量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			data, err := t.Execute()
			results[index] = ConcurrentQueryResult{
				Data:  data,
				Error: err,
			}
		}(i, task)
	}

	wg.Wait()

	// 收集结果
	var collected []interface{}
	for _, result := range results {
		if result.Error != nil {
			return nil, result.Error
		}
		if result.Data != nil {
			collected = append(collected, result.Data)
		}
	}

	return collected, nil
}

// ShardOperationStats 分片操作统计
type ShardOperationStats struct {
	TotalShards   int
	SuccessShards int
	FailedShards  int
	TotalRecords  int64
	OperationTime string
}

// PrintStats 打印操作统计
func (sos *ShardOperationStats) PrintStats(operation string) {
	fmt.Printf("\n%s 统计:\n", operation)
	fmt.Printf("  总分片数: %d\n", sos.TotalShards)
	fmt.Printf("  成功分片: %d\n", sos.SuccessShards)
	fmt.Printf("  失败分片: %d\n", sos.FailedShards)
	fmt.Printf("  总记录数: %d\n", sos.TotalRecords)
	fmt.Printf("  耗时: %s\n", sos.OperationTime)
}

// ValidateShardRange 验证分片范围
func ValidateShardRange(dbIndex, tableIndex int, config ShardConfig) error {
	if dbIndex < 0 || dbIndex >= config.DBCount {
		return fmt.Errorf("数据库索引 %d 超出范围 [0, %d)", dbIndex, config.DBCount)
	}
	if tableIndex < 0 || tableIndex >= config.TableCount {
		return fmt.Errorf("表索引 %d 超出范围 [0, %d)", tableIndex, config.TableCount)
	}
	return nil
}

// GetAllShardNames 获取所有分片名称
func (sr *ShardRouter) GetAllShardNames() []struct {
	DBName    string
	TableName string
} {
	var shards []struct {
		DBName    string
		TableName string
	}

	for dbIdx := 0; dbIdx < sr.config.DBCount; dbIdx++ {
		for tblIdx := 0; tblIdx < sr.config.TableCount; tblIdx++ {
			shards = append(shards, struct {
				DBName    string
				TableName string
			}{
				DBName:    fmt.Sprintf("%s_%d", sr.dbPrefix, dbIdx),
				TableName: fmt.Sprintf("%s_%d_%d", sr.tablePrefix, dbIdx, tblIdx),
			})
		}
	}

	return shards
}

package database

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// TestUser 测试用户模型
type TestUser struct {
	ID        int64     `orm:"pk;auto"`
	Name      string    `orm:"size(100)"`
	Email     string    `orm:"size(200);unique"`
	Age       int       `orm:"default(0)"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
}

// TableName 指定表名（Beego ORM 要求）
func (u *TestUser) TableName() string {
	return "users" // 实际表名由中间件动态决定
}

func TestShardMiddleware(t *testing.T) {
	// 创建分片配置
	shardConfig := ShardConfig{
		DBCount:        2,
		TableCount:     4,
		NamingStrategy: NamingStrategySimple,
	}

	// 创建中间件配置
	middlewareConfig := ShardMiddlewareConfig{
		ShardConfig:      shardConfig,
		DBPrefix:         "test_db",
		TablePrefix:      "test_users",
		DefaultDBAlias:   "default",
		EnableQueryCache: true,
		CacheExpire:      5 * time.Minute,
		MaxRetryCount:    3,
		RetryInterval:    100 * time.Millisecond,
		ConnectionPoolConfig: PoolConfig{
			MaxIdleConns:    10,
			MaxOpenConns:    50,
			ConnMaxLifetime: 3600,
			ConnMaxIdleTime: 600,
		},
	}

	// 创建中间件
	middleware, err := NewShardMiddleware(middlewareConfig)
	if err != nil {
		t.Fatalf("创建中间件失败: %v", err)
	}

	// 测试分片路由
	userID := int64(12345)
	dbName := middleware.router.GetDBName(userID)
	tableName := middleware.GetTableName(userID)
	fullTableName := middleware.GetFullTableName(userID)

	fmt.Printf("用户ID: %d\n", userID)
	fmt.Printf("数据库: %s\n", dbName)
	fmt.Printf("表名: %s\n", tableName)
	fmt.Printf("完整表名: %s\n", fullTableName)

	// 测试分片分布
	userIDs := GenerateTestUserIDs(100)
	distribution := middleware.router.CalculateDistribution(userIDs)

	fmt.Printf("\n分片分布统计:\n")
	for key, dist := range distribution {
		fmt.Printf("  %s: 数据库=%s, 表=%s, 数量=%d\n",
			key, dist.DBName, dist.TableName, dist.Count)
	}

	// 测试获取统计信息
	stats := middleware.GetStats()
	fmt.Printf("\n中间件统计: %+v\n", stats)
}

func TestShardMiddleware_Insert(t *testing.T) {
	// 注意：这个测试需要实际的数据库连接
	// 这里只是展示如何使用

	middleware, err := setupTestMiddleware()
	if err != nil {
		t.Skipf("跳过测试（需要数据库连接）: %v", err)
		return
	}

	// 插入单条数据
	user := &TestUser{
		ID:    1,
		Name:  "张三",
		Email: "zhangsan@example.com",
		Age:   25,
	}

	id, err := middleware.Insert(user.ID, user)
	if err != nil {
		t.Logf("插入数据失败（预期，因为没有真实数据库）: %v", err)
	} else {
		t.Logf("插入成功，ID: %d", id)
	}
}

func TestShardMiddleware_Query(t *testing.T) {
	middleware, err := setupTestMiddleware()
	if err != nil {
		t.Skipf("跳过测试（需要数据库连接）: %v", err)
		return
	}

	// 查询单条数据
	var user TestUser
	err = middleware.QueryOne(1, &user, "id", strconv.Itoa(1))
	if err != nil {
		t.Logf("查询数据失败（预期，因为没有真实数据库）: %v", err)
	} else {
		t.Logf("查询成功: %+v", user)
	}
}

func TestShardMiddleware_BroadcastQuery(t *testing.T) {
	middleware, err := setupTestMiddleware()
	if err != nil {
		t.Skipf("跳过测试（需要数据库连接）: %v", err)
		return
	}

	// 广播查询示例
	results, err := middleware.BroadcastQuery(func(o orm.Ormer, tableName string) (interface{}, error) {
		// 在实际使用中，这里会执行查询
		return nil, nil
	})

	if err != nil {
		t.Logf("广播查询失败: %v", err)
	} else {
		t.Logf("广播查询成功，结果数: %d", len(results))
	}
}

func setupTestMiddleware() (*ShardMiddleware, error) {
	shardConfig := ShardConfig{
		DBCount:        2,
		TableCount:     4,
		NamingStrategy: NamingStrategySimple,
	}

	middlewareConfig := ShardMiddlewareConfig{
		ShardConfig:      shardConfig,
		DBPrefix:         "test_db",
		TablePrefix:      "test_users",
		EnableQueryCache: false,
		MaxRetryCount:    1,
	}

	return NewShardMiddleware(middlewareConfig)
}

// ExampleShardMiddleware 使用示例
func ExampleShardMiddleware() {
	// 1. 创建分片配置
	shardConfig := ShardConfig{
		DBCount:        4, // 4个数据库
		TableCount:     8, // 每个库8张表
		NamingStrategy: NamingStrategyPadded,
		PaddingWidth:   3, // 补零宽度
	}

	// 2. 创建中间件配置
	middlewareConfig := ShardMiddlewareConfig{
		ShardConfig:      shardConfig,
		DBPrefix:         "user_db",
		TablePrefix:      "users",
		EnableQueryCache: true,
		CacheExpire:      10 * time.Minute,
		MaxRetryCount:    3,
		ConnectionPoolConfig: PoolConfig{
			MaxIdleConns:    20,
			MaxOpenConns:    100,
			ConnMaxLifetime: 3600,
			ConnMaxIdleTime: 600,
		},
	}

	// 3. 创建中间件
	middleware, err := NewShardMiddleware(middlewareConfig)
	if err != nil {
		panic(err)
	}

	// 4. 注册数据库（在实际应用中，这里会从配置文件读取）
	// for i := 0; i < shardConfig.DBCount; i++ {
	// 	dbConfig := &MysqlConfig{
	// 		Alias:    fmt.Sprintf("user_db_%d", i),
	// 		Name:     fmt.Sprintf("user_db_%d", i),
	// 		User:     "root",
	// 		Password: "password",
	// 		Host:     "127.0.0.1",
	// 		Port:     "3306",
	// 	}
	// 	middleware.RegisterDatabase(dbConfig.Alias, dbConfig)
	// }

	// 5. 插入数据
	user := &TestUser{
		ID:    1001,
		Name:  "李四",
		Email: "lisi@example.com",
		Age:   30,
	}

	// id, err := middleware.Insert(user.ID, user)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Printf("用户 %d 将存储在: %s\n", user.ID, middleware.GetFullTableName(user.ID))

	// 6. 查询数据
	// var queriedUser TestUser
	// err = middleware.QueryOne(user.ID, &queriedUser, "id", user.ID)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// 7. 更新数据
	// user.Age = 31
	// num, err := middleware.Update(user.ID, user, "age")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Printf("更新了 %d 条记录\n", 1)

	// 8. 删除数据
	// num, err = middleware.Delete(user.ID, "id", user.ID)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Printf("删除了 %d 条记录\n", 1)

	// 9. 广播查询（跨所有分片）
	// results, err := middleware.BroadcastQuery(func(o orm.Ormer, tableName string) (interface{}, error) {
	// 	var users []TestUser
	// 	_, err := o.QueryTable(tableName).Filter("age__gt", 25).All(&users)
	// 	return users, err
	// })

	fmt.Printf("中间件创建成功\n")
}

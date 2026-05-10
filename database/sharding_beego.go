package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/XingMenTech/common/logger"
	"github.com/beego/beego/v2/client/orm"
)

// ==================== Beego ORM 分库分表解决方案 ====================

// BeegoShardedUser Beego分片用户模型
type BeegoShardedUser struct {
	ID        int64     `orm:"pk;auto"`
	Name      string    `orm:"size(100)"`
	Email     string    `orm:"size(200);unique"`
	Age       int       `orm:"default(0)"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
}

func (u *BeegoShardedUser) TableName() string {
	return "beego_users"
}

// BeegoShardRouter Beego分片路由器（扩展公共路由器）
type BeegoShardRouter struct {
	*ShardRouter                      // 嵌入公共路由器
	dbs          map[string]orm.Ormer // 多个数据库连接
}

// NewBeegoShardRouter 创建Beego分片路由器
func NewBeegoShardRouter(config ShardConfig) (*BeegoShardRouter, error) {
	commonRouter, err := NewShardRouter(config, "beego_db", "beego_users")
	if err != nil {
		return nil, err
	}

	return &BeegoShardRouter{
		ShardRouter: commonRouter,
		dbs:         make(map[string]orm.Ormer),
	}, nil
}

// RegisterDB 注册数据库连接
func (bsr *BeegoShardRouter) RegisterDB(name string, ormer orm.Ormer) {
	bsr.dbs[name] = ormer
}

// GetDB 获取对应的数据库连接
func (bsr *BeegoShardRouter) GetDB(userID int64) (orm.Ormer, error) {
	dbName := bsr.GetDBName(userID)
	db, ok := bsr.dbs[dbName]
	if !ok {
		return nil, fmt.Errorf("数据库 %s 未注册", dbName)
	}
	return db, nil
}

// 方案11: Beego分库分表插入数据
func beegoShardedInsert(router *BeegoShardRouter, users []BeegoShardedUser) error {
	fmt.Println("开始Beego分库分表插入...")

	// 按数据库分组
	dbGroups := make(map[string][]BeegoShardedUser)
	for _, user := range users {
		dbName := router.GetDBName(user.ID)
		dbGroups[dbName] = append(dbGroups[dbName], user)
	}

	totalInserted := 0
	for dbName, dbUsers := range dbGroups {
		db, err := router.GetDB(dbUsers[0].ID)
		if err != nil {
			return fmt.Errorf("获取数据库连接失败: %w", err)
		}

		// 在该数据库中按表分组
		tableGroups := make(map[string][]BeegoShardedUser)
		for _, user := range dbUsers {
			tableName := router.GetTableName(user.ID)
			tableGroups[tableName] = append(tableGroups[tableName], user)
		}

		// 批量插入到各个表
		for tableName, tableUsers := range tableGroups {
			err := beegoInsertToTable(db, tableName, tableUsers)
			if err != nil {
				return fmt.Errorf("插入表 %s 失败: %w", tableName, err)
			}
			totalInserted += len(tableUsers)
			fmt.Printf("  数据库 %s, 表 %s: 插入 %d 条数据\n", dbName, tableName, len(tableUsers))
		}
	}

	fmt.Printf("Beego分库分表插入完成，共插入 %d 条数据\n", totalInserted)
	return nil
}

// Beego向指定表批量插入数据
func beegoInsertToTable(o orm.Ormer, tableName string, users []BeegoShardedUser) error {
	// 创建表（如果不存在）
	createSQL := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(200) NOT NULL UNIQUE,
		age INT NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		INDEX idx_created_at (created_at)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`, tableName)

	_, err := o.Raw(createSQL).Exec()
	if err != nil {
		return fmt.Errorf("创建表失败: %w", err)
	}

	// 批量插入
	batchSize := 1000
	for i := 0; i < len(users); i += batchSize {
		end := i + batchSize
		if end > len(users) {
			end = len(users)
		}

		batch := users[i:end]
		if _, err := o.InsertMulti(len(batch), batch); err != nil {
			return fmt.Errorf("批量插入失败: %w", err)
		}
	}

	return nil
}

// 方案12: Beego分库分表查询单条记录
func beegoShardedGetUser(router *BeegoShardRouter, userID int64) (*BeegoShardedUser, error) {
	db, err := router.GetDB(userID)
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %w", err)
	}

	tableName := router.GetTableName(userID)
	var user BeegoShardedUser

	err = db.QueryTable(tableName).Filter("id", userID).One(&user)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("查询失败: %w", err)
	}

	return &user, nil
}

// 方案13: Beego分库分表范围查询（跨分片）
func beegoShardedQueryRange(router *BeegoShardRouter, startTime, endTime time.Time) ([]BeegoShardedUser, error) {
	var allUsers []BeegoShardedUser

	// 遍历所有数据库
	for dbName, db := range router.dbs {
		// 遍历该数据库中的所有表
		for i := 0; i < router.config.TableCount; i++ {
			tableName := fmt.Sprintf("beego_users_%d_%d",
				int(dbName[len(dbName)-1])-48, i)

			var users []BeegoShardedUser
			_, err := db.QueryTable(tableName).
				Filter("created_at__gte", startTime).
				Filter("created_at__lte", endTime).
				All(&users)

			if err != nil {
				// 表可能不存在，跳过
				continue
			}

			allUsers = append(allUsers, users...)
		}
	}

	return allUsers, nil
}

// 方案14: Beego分库分表广播查询（并发）
func beegoShardedBroadcastQuery(router *BeegoShardRouter, age int) ([]BeegoShardedUser, error) {
	var allUsers []BeegoShardedUser

	type result struct {
		users []BeegoShardedUser
		err   error
	}

	results := make(chan result, len(router.dbs)*router.config.TableCount)

	// 并发查询所有分片
	for dbName, db := range router.dbs {
		for i := 0; i < router.config.TableCount; i++ {
			go func(dbName string, db orm.Ormer, tableIndex int) {
				tableName := fmt.Sprintf("beego_users_%d_%d",
					int(dbName[len(dbName)-1])-48, tableIndex)

				var users []BeegoShardedUser
				_, err := db.QueryTable(tableName).
					Filter("age__gt", age).
					All(&users)

				results <- result{users, err}
			}(dbName, db, i)
		}
	}

	// 收集结果
	for i := 0; i < len(router.dbs)*router.config.TableCount; i++ {
		res := <-results
		if res.err != nil {
			return nil, fmt.Errorf("查询分片失败: %w", res.err)
		}
		allUsers = append(allUsers, res.users...)
	}

	return allUsers, nil
}

// 方案15: Beego分库分表批量更新
func beegoShardedBatchUpdate(router *BeegoShardRouter, ageThreshold int, increment int) error {
	fmt.Println("开始Beego分库分表批量更新...")

	totalUpdated := int64(0)

	// 遍历所有数据库和表
	for dbName, db := range router.dbs {
		for i := 0; i < router.config.TableCount; i++ {
			tableName := fmt.Sprintf("beego_users_%d_%d", int(dbName[len(dbName)-1])-48, i)

			num, err := db.QueryTable(tableName).
				Filter("age__lt", ageThreshold).
				Update(orm.Params{
					"age": orm.ColValue(orm.ColAdd, increment),
				})

			if err != nil {
				return fmt.Errorf("更新表 %s 失败: %w", tableName, err)
			}

			totalUpdated += num
			fmt.Printf("  表 %s: 更新 %d 条数据\n", tableName, num)
		}
	}

	fmt.Printf("Beego分库分表批量更新完成，共更新 %d 条数据\n", totalUpdated)
	return nil
}

// 初始化Beego分库分表环境
func initBeegoShardedDBs(router *BeegoShardRouter) error {
	fmt.Println("初始化Beego分库分表环境...")

	// 创建多个数据库连接
	for i := 0; i < router.config.DBCount; i++ {
		dbName := fmt.Sprintf("beego_db_%d", i)

		// DSN配置
		dsn := fmt.Sprintf("root:password@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&loc=Local", dbName)

		// 注册数据库
		err := orm.RegisterDataBase(dbName, "mysql", dsn)
		if err != nil {
			return fmt.Errorf("注册数据库 %s 失败: %w", dbName, err)
		}

		// 设置连接池
		orm.SetMaxIdleConns(dbName, 30)
		orm.SetMaxOpenConns(dbName, 100)

		// 创建ORM实例（Beego ORM v2会自动使用注册的默认数据库）
		o := orm.NewOrm()

		// 尝试创建数据库（使用database/sql）
		rootDSN := "root:password@tcp(127.0.0.1:3306)/"
		rootDB, err := sql.Open("mysql", rootDSN)
		if err == nil {
			defer rootDB.Close()
			_, _ = rootDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4", dbName))
		}

		router.RegisterDB(dbName, o)
		fmt.Printf("  Beego数据库 %s 已注册\n", dbName)
	}

	return nil
}

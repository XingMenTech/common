package database

import (
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"gitlab.novgate.com/common/common/logger"
)

// ------------------[mysql]-------------------
type MysqlConfig struct {
	//数据库类别
	DatabaseType string `yaml:"db_type"`
	//连接名称
	Alias string `yaml:"db_alias"`
	//数据库名称
	Name string `yaml:"db_name"`
	//数据库连接用户名
	User string `yaml:"db_user"`
	//数据库连接用户名
	Password string `yaml:"db_pwd"`
	//数据库IP（域名）
	Host string `yaml:"db_host"`
	//数据库端口
	Port string `yaml:"db_port"`
	//字符集类型
	Charset string `yaml:"db_charset"`
	//搜索最大条数限制,-1不限制
	DefaultRowsLimit int `yaml:"default_rows_limit"`
	//是否调试模式
	Debug bool `yaml:"db_debug"`
	//表前缀
	TablePrefix string `yaml:"db_table_prefix"`
}

func InitMysql(config *MysqlConfig) error {

	//config := config.GetDatabaseConfig()
	if config == nil {
		return errors.New("init database fail. can not find database config")
	}
	//config.DatabaseConfigGlobal = databaseConfig

	//数据库类别
	dbType := "mysql"
	//连接名称
	dbAlias := config.Alias
	//数据库名称
	dbName := config.Name
	//数据库连接用户名
	dbUser := config.User
	//数据库连接用户名
	dbPwd := config.Password
	//数据库IP（域名）
	dbHost := config.Host
	//数据库端口
	dbPort := config.Port
	//字符集
	dbCharset := config.Charset
	err := orm.RegisterDriver("mysql", orm.DRMySQL)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&loc=Local", dbUser, dbPwd, dbHost, dbPort, dbName, dbCharset)
	logger.LOG.Debugf("数据库链接：%s \n", path)
	if err := orm.RegisterDataBase(dbAlias, dbType, path); err != nil {
		return err
	}
	orm.DefaultRowsLimit = -1

	//如果是开发模式，则显示命令信息
	if config.Debug {
		orm.Debug = true
	}
	return nil
}

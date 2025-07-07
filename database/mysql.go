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
	DatabaseType     string `yaml:"db_type" json:"type" comment:"数据库类别"`
	Alias            string `yaml:"db_alias" json:"alias" comment:"连接名称"`
	Name             string `yaml:"db_name" json:"name" comment:"数据库名称"`
	User             string `yaml:"db_user" json:"user" comment:"数据库连接用户名"`
	Password         string `yaml:"db_pwd" json:"password" comment:"数据库连接用户名"`
	Host             string `yaml:"db_host" json:"host" comment:"数据库IP（域名）"`
	Port             string `yaml:"db_port" json:"port" comment:"数据库端口"`
	Charset          string `yaml:"db_charset" json:"charset" comment:"字符集类型"`
	DefaultRowsLimit int    `yaml:"default_rows_limit" json:"defaultRowsLimit" comment:"搜索最大条数限制,-1不限制"`
	Debug            string `yaml:"db_debug" json:"debug" comment:"是否调试模式"`
	TablePrefix      string `yaml:"db_table_prefix" json:"tablePrefix" comment:"表前缀"`
}

func (c *MysqlConfig) Url() string {
	path := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&loc=Local", c.User, c.Password, c.Host, c.Port, c.Name, c.Charset)
	logger.LOG.Debugf("数据库链接：%s \n", path)
	return path
}
func InitMysql(config *MysqlConfig) error {

	if config == nil {
		return errors.New("init database fail. can not find database config")
	}

	err := orm.RegisterDriver(config.DatabaseType, orm.DRMySQL)
	if err != nil {
		return err
	}

	if err := orm.RegisterDataBase(config.Alias, config.DatabaseType, config.Url()); err != nil {
		return err
	}
	orm.DefaultRowsLimit = -1

	//如果是开发模式，则显示命令信息
	if config.Debug == "true" {
		orm.Debug = true
	}
	prefix = config.TablePrefix
	return nil
}

var prefix string

func TableName(tableName string) string {
	return prefix + tableName
}

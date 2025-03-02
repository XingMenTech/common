package database

import (
	"github.com/zhangyuanCloud/common/logger"
	"testing"
)

func TestMysql(t *testing.T) {
	logger.InitializeLogger(&logger.LogConfig{})
	config := &MysqlConfig{
		DatabaseType:     "mysql",
		Alias:            "default",
		Name:             "xxx",
		User:             "user",
		Password:         "password",
		Host:             "dbhost",
		Port:             "port",
		Charset:          "utf8",
		DefaultRowsLimit: 1,
		Debug:            true,
		TablePrefix:      "",
	}
	err := InitMysql(config)
	if err != nil {
		t.Error(err)
	}
}

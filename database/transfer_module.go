package database

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/sirupsen/logrus"
)

type TransferModule struct {
	Or  orm.TxOrmer
	Log *logrus.Entry
	Err error
}

func (o *TransferModule) InitTran() {
	//o.Or = orm.NewOrm()
}

func (o *TransferModule) Begin() bool {
	if o.Or, o.Err = orm.NewOrm().Begin(); o.Err != nil {
		o.Log.WithField("error", o.Err).Errorf("启动事务失败")
		return false
	}
	return true
}

func (o *TransferModule) Commit() bool {
	if o.Err = o.Or.Commit(); o.Err != nil {
		o.Log.WithField("error", o.Err).Error("提交事务失败")
		o.Rollback()
		return false
	}
	return true
}

func (o *TransferModule) Rollback() {
	if err := o.Or.Rollback(); err != nil && err != orm.ErrTxDone {
		o.Log.WithField("error", err).Error("回滚失败")
	}
}

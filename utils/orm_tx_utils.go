package utils

import (
	"github.com/beego/beego/v2/client/orm"
)

type TxFunc func(o orm.TxOrmer) error

type OrmTx struct {
	o orm.TxOrmer
}

func NewOrmTx() *OrmTx {
	ormer, _ := orm.NewOrm().Begin()
	return &OrmTx{
		o: ormer,
	}
}

func (tx *OrmTx) Execute(f TxFunc) error {
	err := f(tx.o)
	if err != nil {
		_ = tx.o.Rollback()
		return err
	}

	return tx.o.Commit()
}

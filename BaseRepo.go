package common

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/sirupsen/logrus"
)

type BaseRepo struct {
	TableName string
	Log       *logrus.Entry
}

func (r *BaseRepo) InsertOne(o orm.TxOrmer, m interface{}) (i int64, err error) {
	if o == nil {
		i, err = orm.NewOrm().Insert(m)
		if err != nil {
			err = NewMsgError(CommonDbInsertError, err.Error())
		}
		return
	}

	i, err = o.Insert(m)
	if err != nil {
		err = NewMsgError(CommonDbInsertError, err.Error())
	}
	return
}

func (r *BaseRepo) InsertBatch(o orm.TxOrmer, bulk int, m interface{}) (i int64, err error) {
	if o == nil {
		i, err = orm.NewOrm().InsertMulti(bulk, m)
		if err != nil {
			err = NewMsgError(CommonDbInsertError, err.Error())
		}
		return
	}
	i, err = o.InsertMulti(bulk, m)
	if err != nil {
		err = NewMsgError(CommonDbInsertError, err.Error())
	}
	return
}

func (r *BaseRepo) ReadOne(m interface{}, cols ...string) error {
	return orm.NewOrm().Read(m, cols...)
}

func (r *BaseRepo) Update(o orm.TxOrmer, m interface{}, cols ...string) (err error) {
	if o == nil {
		_, err = orm.NewOrm().Update(m, cols...)
		if err != nil {
			err = NewMsgError(CommonDbUpdateError, err.Error())
		}
		return
	}

	_, err = o.Update(m, cols...)
	if err != nil {
		err = NewMsgError(CommonDbUpdateError, err.Error())
	}
	return
}

func (r *BaseRepo) UpdateByCondition(o orm.TxOrmer, cond *orm.Condition, param orm.Params) (i int64, err error) {
	if len(param) <= 0 {
		return 0, orm.ErrArgs
	}
	var query orm.QuerySeter
	if o == nil {
		query = orm.NewOrm().QueryTable(r.TableName)
	} else {
		query = o.QueryTable(r.TableName)
	}

	if cond != nil && !cond.IsEmpty() {
		query = query.SetCond(cond)
	}
	i, err = query.Update(param)
	if err != nil {
		err = NewMsgError(CommonDbUpdateError, err.Error())
	}
	return
}

func (r *BaseRepo) Delete(o orm.TxOrmer, m interface{}, cols ...string) (err error) {
	if o == nil {
		_, err = orm.NewOrm().Delete(m, cols...)
		return
	}
	_, err = o.Delete(m, cols...)
	return
}

func (r *BaseRepo) DeleteByCondition(o orm.TxOrmer, cond *orm.Condition) (err error) {
	if cond.IsEmpty() {
		return orm.ErrArgs
	}
	var query orm.QuerySeter
	if o == nil {
		query = orm.NewOrm().QueryTable(r.TableName)
	} else {
		query = o.QueryTable(r.TableName)
	}

	_, err = query.SetCond(cond).Delete()
	return
}

func (r *BaseRepo) Count(cond *orm.Condition) int64 {
	query := orm.NewOrm().QueryTable(r.TableName).SetCond(cond)
	total, err := query.Count()
	if err != nil {
		r.Log.Errorf("数据表 %s 数据行数查询失败 error: %s", r.TableName, err.Error())
		return 0
	}
	return total
}

func (r *BaseRepo) List(cond *orm.Condition, sort string, container interface{}) (total int64, err error) {
	query := orm.NewOrm().QueryTable(r.TableName)
	if cond != nil {
		query = query.SetCond(cond)
	}
	total, err = query.Count()
	if err != nil {
		r.Log.Errorf("数据表 %s 数据行数查询失败 error: %s", r.TableName, err.Error())
		return
	}
	if total == 0 {
		return
	}
	if len(sort) > 0 {
		query = query.OrderBy(sort)
	}
	_, err = query.All(container)
	if err != nil {
		r.Log.Errorf("数据表 %s 数据查询失败 error: %s", r.TableName, err.Error())
		return 0, err
	}

	return total, nil
}

func (r *BaseRepo) PageList(cond *orm.Condition, pageParam *BaseQueryParam, sort string, container interface{}) (total int64, err error) {

	query := orm.NewOrm().QueryTable(r.TableName).SetCond(cond)
	total, err = query.Count()
	if err != nil {
		r.Log.Errorf("数据表 %s 数据行数查询失败 error: %s", r.TableName, err.Error())
		return
	}
	if total == 0 {
		return
	}

	if len(sort) > 0 {
		query = query.OrderBy(sort)
	}

	if pageParam != nil && pageParam.IsValid() {
		limit, offset := pageParam.GetLimit()
		query = query.Limit(limit).Offset(offset)
	}

	_, err = query.All(container)
	if err != nil {
		r.Log.Errorf("数据表 %s 数据查询失败 error: %s", r.TableName, err.Error())
	}

	return
}

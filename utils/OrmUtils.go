package utils

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/client/orm/clauses/order_clause"
	"strings"
	"time"
)

const (
	DefaultPage     int = 1
	DefaultPageSize int = 20
	MaxPageSize     int = 500
)

type ListParam struct {
	Param *orm.Condition
	Page  *PageParam
	Time  *TimeParam
	Order []*order_clause.Order
}

// BaseQueryParam 用于查询的类
type PageParam struct {
	Page     int `json:"page" form:"page" binding:"required"`
	PageSize int `json:"pageSize" form:"pageSize" binding:"required"`
}

func (bqp *PageParam) IsValid() bool {
	return bqp.Page > 0 && bqp.PageSize > 0
}

func (bqp *PageParam) Offset() int {
	offset := 0
	if bqp.Page > 1 {
		offset = (bqp.Page - 1) * bqp.PageSize
	}
	return offset
}

func (bqp *PageParam) GetLimit() (limit, offset int) {
	if bqp.Page < DefaultPage {
		bqp.Page = DefaultPage
	}
	if bqp.PageSize <= 0 {
		bqp.PageSize = DefaultPageSize
	}
	if bqp.PageSize > MaxPageSize {
		bqp.PageSize = MaxPageSize
	}

	limit = bqp.PageSize
	offset = (bqp.Page - 1) * bqp.PageSize
	return
}

type TimeParam struct {
	Column    string `json:"column" form:"column"`
	StartTime string `json:"startTime" form:"startTime"` //开始时间
	EndTime   string `json:"endTime" form:"endTime"`     //结束时间
}

func (req *TimeParam) IsValid() bool {
	if req == nil {
		return false
	}
	return req.StartTime != "" && req.EndTime != ""
}

func (req *TimeParam) GetTime() (start, end time.Time) {
	t1 := ParseLocalTime(strings.TrimSpace(req.StartTime))
	t2 := ParseLocalTime(strings.TrimSpace(req.EndTime))
	return t1, t2
}

func FindOne[T any](id int64) *T {
	var model T
	err := orm.NewOrm().QueryTable(&model).Filter("id", id).One(&model)
	if err != nil {
		return nil
	}
	return &model
}

func FindAll[T any](form ListParam) (list []*T, total int64, err error) {
	query := orm.NewOrm().QueryTable(new(T))
	if form.Param != nil && !form.Param.IsEmpty() {
		query = query.SetCond(form.Param)
	}
	timeParam := form.Time
	if timeParam != nil && timeParam.IsValid() {
		column := timeParam.Column
		start, end := timeParam.GetTime()
		query = query.Filter(fmt.Sprintf("%s__gte", column), start).Filter(fmt.Sprintf("%s__lt", column), end)
	}

	total, err = query.Count()
	if err != nil {
		return
	}

	if total == 0 {
		return
	}

	if len(form.Order) > 0 {
		query = query.OrderClauses(form.Order...)
	}

	if form.Page != nil && form.Page.IsValid() {
		limit, offset := form.Page.GetLimit()
		query = query.Limit(limit, offset)
	}

	list = make([]*T, 0)
	_, err = query.All(&list)

	return
}

func Count[T any](cond *orm.Condition) (int64, error) {
	query := orm.NewOrm().QueryTable(new(T)).SetCond(cond)
	return query.Count()
}

func Update[T any](o orm.TxOrmer, form T, columns ...string) (err error) {
	if o == nil {
		_, err = orm.NewOrm().Update(form, columns...)
	} else {
		_, err = o.Update(form, columns...)
	}
	//if err != nil {
	//	err = common.NewMsgError(common.CommonDbUpdateError, err.Error())
	//}
	return err
}

func UpdateByCondition[T any](o orm.TxOrmer, cond *orm.Condition, param orm.Params) (err error) {
	if len(param) <= 0 {
		return orm.ErrArgs
	}
	var query orm.QuerySeter
	if o == nil {
		query = orm.NewOrm().QueryTable(new(T))
	} else {
		query = o.QueryTable(new(T))
	}

	if cond != nil && !cond.IsEmpty() {
		query = query.SetCond(cond)
	}
	_, err = query.Update(param)
	//if err != nil {
	//	err = common.NewMsgError(common.CommonDbUpdateError, err.Error())
	//}
	return
}

func Delete[T any](o orm.TxOrmer, form T) (err error) {
	if o == nil {
		_, err = orm.NewOrm().Delete(form)
	} else {
		_, err = o.Delete(form)
	}
	//if err != nil {
	//	err = common.NewMsgError(common.CommonDbUpdateError, err.Error())
	//}
	return
}

func DeleteByCondition[T any](o orm.TxOrmer, cond *orm.Condition) (err error) {
	if cond.IsEmpty() {
		return orm.ErrArgs
	}

	var query orm.QuerySeter
	if o == nil {
		query = orm.NewOrm().QueryTable(new(T))
	} else {
		query = o.QueryTable(new(T))
	}

	_, err = query.SetCond(cond).Delete()
	//if err != nil {
	//	err = common.NewMsgError(common.CommonDbUpdateError, err.Error())
	//}
	return
}

func Insert[T any](o orm.TxOrmer, form T) (err error) {
	if o == nil {
		_, err = orm.NewOrm().Insert(form)
	} else {
		_, err = o.Insert(form)
	}
	//if err != nil {
	//	err = common.NewMsgError(common.CommonDbInsertError, err.Error())
	//}
	return
}

func InsertBatch(o orm.TxOrmer, bulk int, m interface{}) (i int64, err error) {
	if o == nil {
		i, err = orm.NewOrm().InsertMulti(bulk, m)
	} else {
		i, err = o.InsertMulti(bulk, m)
	}

	//if err != nil {
	//	err = common.NewMsgError(common.CommonDbInsertError, err.Error())
	//}

	return
}

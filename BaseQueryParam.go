package common

import (
	"time"

	"github.com/XingMenTech/common/utils"
)

const (
	DefaultPage     int = 1
	DefaultPageSize int = 20
	MaxPageSize     int = 500
)

// PageParam 用于查询的类
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

// TimeParam 时间区间参数
type TimeParam struct {
	Column    string `json:"column" form:"column"`
	StartTime string `json:"startTime" form:"startTime" binding:"required"` //开始时间
	EndTime   string `json:"endTime" form:"endTime" binding:"required"`     //结束时间
}

func (req *TimeParam) IsValid() bool {
	if req == nil {
		return false
	}
	return req.StartTime != "" && req.EndTime != ""
}

func (req *TimeParam) GetTime() (start, end time.Time) {
	t1 := utils.ParseLocalTime(req.StartTime)
	t2 := utils.ParseLocalTime(req.EndTime)
	return t1, t2
}

func (req *TimeParam) DiffDays() int {
	t1, t2 := req.GetTime()
	return int(t2.Sub(t1).Hours() / 24)
}

type IdParam struct {
	Id int64 `json:"id" form:"id" binding:"required"`
}

func IdParamError() map[string]string {
	return map[string]string{
		"Id.required": "主键参数不能为空",
	}
}

type IdArrParam struct {
	Id []int64 `json:"id" form:"id" binding:"required"`
}

package common

import (
	"time"
)

const (
	DefaultPage     int = 1
	DefaultPageSize int = 20
	MaxPageSize     int = 500
)

// BaseQueryParam 用于查询的类
type BaseQueryParam struct {
	Page     int `json:"page" form:"page" binding:"required"`
	PageSize int `json:"pageSize" form:"pageSize" binding:"required"`
}

func (bqp *BaseQueryParam) IsValid() bool {
	return bqp.Page > 0 && bqp.PageSize > 0
}

func (bqp *BaseQueryParam) Offset() int {
	offset := 0
	if bqp.Page > 1 {
		offset = (bqp.Page - 1) * bqp.PageSize
	}
	return offset
}

func (bqp *BaseQueryParam) GetLimit() (limit, offset int) {
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

const DateTimeFormat = "2006-01-02 15:04:05"
const DateFormat = "2006-01-02"

// BaseTimeRequest 时间区间参数
type BaseTimeRequest struct {
	StartTime string `json:"startTime" form:"startTime"` //开始时间
	EndTime   string `json:"endTime" form:"endTime"`     //结束时间
}

func (req *BaseTimeRequest) IsValid() bool {
	if req == nil {
		return false
	}
	return req.StartTime != "" && req.EndTime != ""
}

func (req *BaseTimeRequest) GetTime() (start, end time.Time) {
	t1, _ := time.ParseInLocation(DateTimeFormat, req.StartTime, time.UTC)
	t2, _ := time.ParseInLocation(DateTimeFormat, req.EndTime, time.UTC)
	return t1, t2
}

func (req *BaseTimeRequest) GetDate() (start, end time.Time) {
	t1, _ := time.ParseInLocation(DateFormat, req.StartTime, time.UTC)
	t2, _ := time.ParseInLocation(DateFormat, req.EndTime, time.UTC)
	return t1, t2
}

func (req *BaseTimeRequest) GetFormatTime() (start, end string) {
	t1, t2 := req.GetTime()
	return t1.Format(DateTimeFormat), t2.Format(DateTimeFormat)
}

func (req *BaseTimeRequest) GetFormatDate() (start, end string) {
	t1, t2 := req.GetDate()
	return t1.Format(DateFormat), t2.Format(DateFormat)
}

func (req *BaseTimeRequest) DiffDays() int {
	t1, t2 := req.GetTime()
	return int(t2.Sub(t1).Hours() / 24)
}

type BaseIdParam struct {
	Id int64 `json:"id" form:"id" binding:"required"`
}

func BaseIdParamError() map[string]string {
	return map[string]string{
		"Id.required": "主键参数不能为空",
	}
}

type BaseIdArrParam struct {
	Id []int64 `json:"id" form:"id" binding:"required"`
}

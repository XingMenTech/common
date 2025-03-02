package common

const (
	StatusDelete  = 3
	StatusDisable = 2
	StatusEnable  = 1
)

var Status = map[int]string{
	StatusDelete:  "删除",
	StatusDisable: "禁用",
	StatusEnable:  "启用",
}

// 在线状态
const (
	UserOnlineStatus  = 1
	UserOfflineStatus = 2
)

const (
	Monday    = 1
	Tuesday   = 2
	Wednesday = 3
	Thursday  = 4
	Friday    = 5
	Saturday  = 6
	Sunday    = 7
)

var NoticeWeekDayPush = map[int]string{
	Monday:    "MON",
	Tuesday:   "TUE",
	Wednesday: "WED",
	Thursday:  "THU",
	Friday:    "FRI",
	Saturday:  "SAT",
	Sunday:    "SUN",
}

const (
	CompleteStatusYes = 1
	CompleteStatusNo  = 2
)

var BetCompleteStatus = map[int]string{
	CompleteStatusYes: "完成",
	CompleteStatusNo:  "未完成",
}

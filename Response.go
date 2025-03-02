package common

type TypeOption struct {
	Type     int    `json:"type"`
	TypeName string `json:"typeName"`
}

type SelectOption struct {
	Type     string `json:"type"`
	TypeName string `json:"typeName"`
}

type SelectTree struct {
	Id       int           `json:"id"`
	Label    string        `json:"label"`
	Children []*SelectTree `json:"children"`
}

type Response struct {
	Code    int    `json:"code" doc:"响应码,400"`
	Message string `json:"message" doc:"响应消息,status bad request"`
}

type DataResponse struct {
	Response
	Data interface{} `json:"data"`
}

type PageResponse struct {
	TotalCount int64       `json:"totalCount"`
	TotalPage  int64       `json:"totalPage"`
	PageSize   int         `json:"pageSize"`
	List       interface{} `json:"list"`
}

type PageStatisticResponse struct {
	TotalCount int64       `json:"totalCount"`
	TotalPage  int64       `json:"totalPage"`
	PageSize   int         `json:"pageSize"`
	Statistic  interface{} `json:"statistic"`
	List       interface{} `json:"list"`
}

// 提现
type WithdrawPageResponse struct {
	TotalCount           int64       `json:"totalCount"`
	TotalPage            int64       `json:"totalPage"`
	Size                 int         `json:"size"`
	CurrentPageAmountSum float64     `json:"currentPageAmountSum"`
	AllAmountSum         float64     `json:"allAmountSum"`
	List                 interface{} `json:"list"`
}

// 人工入款
type DepositPageResponse struct {
	TotalCount             int
	TotalPage              int
	Size                   int
	CurrentPageAmountSum   float64     `json:"currentPageAmountSum"`
	CurrentPageDiscountSum float64     `json:"currentPageDiscountSum"`
	AllAmountSum           float64     `json:"allAmountSum"`
	AllDiscountSum         float64     `json:"allDiscountSum"`
	List                   interface{} `json:"list"`
}

// 银行充值
type CardRechargePageResponse struct {
	TotalCount             int64       `json:"totalCount"`
	TotalPage              int64       `json:"totalPage"`
	Size                   int         `json:"size"`
	CurrentPageAmountSum   float64     `json:"currentPageAmountSum"`
	CurrentPageDiscountSum float64     `json:"currentPageDiscountSum"`
	AllAmountSum           float64     `json:"allAmountSum"`
	AllDiscountSum         float64     `json:"allDiscountSum"`
	List                   interface{} `json:"list"`
}

// 注单列表
type GameRecordPageResponse struct {
	TotalCount int64       `json:"totalCount"`
	TotalPage  int64       `json:"totalPage"`
	PageSize   int         `json:"pageSize"`
	List       interface{} `json:"list"`
}

// 打码量总和
type RequireBetAmountSumResponse struct {
	IncomeSum    float64 `json:"incomeSum"`    //入款总和
	DiscountSum  float64 `json:"discountSum"`  //优惠总和
	RequiredSum  float64 `json:"requiredSum"`  //打码量总和
	CompletedSum float64 `json:"completedSum"` //已完成总和
}

// 注单总和
type GameRecordSumResponse struct {
	BetAmountSum         float64 `json:"betAmountSum"`
	WinsAmountSum        float64 `json:"winsAmountSum"`
	ProfitSum            float64 `json:"profitSum"`
	JackpotBonusSum      float64 `json:"jackpotBonusSum"`
	JackpotContributeSum float64 `json:"jackpotContributeSum"`
}

type BaifuLotterysPageResponse struct {
	TotalCount int64       `json:"totalCount"`
	TotalPage  int64       `json:"totalPage"`
	PageSize   int         `json:"pageSize"`
	List       interface{} `json:"list"`
}

type Captcha struct {
	CaptchaId string `json:"captchaId"`
	ImageUrl  string `json:"imageUrl"`
}

// 體育注單總和
type TyLotterySumResponse struct {
	ValidBetAmountSum float64 `json:"validBetAmountSum"`
	BetAmountSum      float64 `json:"betAmountSum"`
	WinAmountSum      float64 `json:"winAmountSum"`
	WinLoseSum        float64 `json:"winLoseSum"`
}

type BaifuLotterySumResponse struct {
	BetAmountSum     float64 `json:"betAmountSum"`
	WinningBounusSum float64 `json:"winningBounusSum"`
	RealWinSum       float64 `json:"realWinSum"`
}

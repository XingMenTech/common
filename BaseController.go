package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/XingMenTech/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

type BaseController struct {
	validate *validator.Validate
}

func (baseController *BaseController) ReturnData(ctx *gin.Context, code int, data interface{}, errorMessage ...interface{}) {
	returnData := new(DataResponse)
	returnData.Code = code
	if nil != errorMessage && len(errorMessage) != 0 {
		returnData.Message = CodeMapMessage[code] + fmt.Sprint(errorMessage)
	} else {
		returnData.Message = CodeMapMessage[code]
	}

	returnData.Data = data
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.JSON(http.StatusOK, returnData)
}
func (baseController *BaseController) ReturnErrorCode(ctx *gin.Context, code int) {
	returnData := new(DataResponse)
	returnData.Code = code
	returnData.Message = CodeMapMessage[code]
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.JSON(http.StatusOK, returnData)
}
func (baseController *BaseController) ReturnErrorData(ctx *gin.Context, err error) {
	returnData := new(DataResponse)
	switch errType := err.(type) {
	case Error:
		returnData.Code = errType.ErrorCode()
		returnData.Message = errType.Error()
	default:
		returnData.Code = CommonSystemError
		returnData.Message = CodeMapMessage[CommonSystemError] + err.Error()
	}
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.JSON(http.StatusOK, returnData)
}
func (baseController *BaseController) CheckForm(form interface{}, formError map[string]string) error {

	//if config.ConfigGlobal.App.RunMode == "dev" {
	fmt.Println("==========post data===============================")
	formJson, _ := json.Marshal(form)
	fmt.Println(string(formJson))
	fmt.Println("==========post data===============================")
	//}

	//vconfig := &validator.Config{
	//	TagName: "binding",
	//}
	//参数名 英文 ==》 中文
	var paramsMap = make(map[string]interface{})
	//formJson, _ := json.Marshal(form)
	err := json.Unmarshal(formJson, &paramsMap)
	if err != nil {
		return NewMsgError(CommonParamError, err.Error())
	}
	//for formKey, formValue := range paramsMap {
	//	//判断是否已经配置语言
	//	if chineseTag, ok := utils.EnglishToChinese[formKey]; ok {
	//		paramsMap[chineseTag] = formValue
	//		delete(paramsMap, formKey)
	//	}
	//}

	validate := validator.New(func(v *validator.Validate) {
		v.SetTagName("binding")
	})
	errs := validate.Struct(form)

	if errs != nil {
		var errInfo validator.ValidationErrors
		var errMessage Error
		errors.As(errs, &errInfo)
		for _, info := range errInfo {
			errTag := info.Field() + "." + info.Tag()
			errMessage = NewMsgError(CommonParamError, formError[errTag])
			break
			//errData[err.Name] = formError[errTag]
		}
		return errMessage
	}
	return nil
}

func NewSysMonitorController() *ServerController {
	return &ServerController{
		startTime: time.Now(),
	}
}

type ServerController struct {
	BaseController
	startTime time.Time
}

func (ctrl *ServerController) Heath(ctx *gin.Context) {
	fmt.Println("服务健康检查。。。。。。。。")
	ctx.Writer.WriteString("OK")
	return
}

func (ctrl *ServerController) Monitor(ctx *gin.Context) {

	info := utils.ServerInfo()
	goStartTime := ctrl.startTime.Format(utils.TimeFormat) //启动时间
	goRunTime := time.Now().Unix() - ctrl.startTime.Unix() //运行时长（秒）
	info["goStartTime"] = goStartTime
	info["goRunTime"] = goRunTime

	ctrl.ReturnData(ctx, Success, info)
}

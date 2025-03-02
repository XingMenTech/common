package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/zhangyuanCloud/common"
	"github.com/zhangyuanCloud/common/logger"
    "{{.Project}}/pkg/{{.ModuleName}}/service"
    "{{.Project}}/pkg/{{.ModuleName}}/validate"
)
var {{.VarFieldName}} *{{.ModelName}}

type {{.ModelName}} struct {
	common.BaseController
	service *service.{{.ModelName}}Service
	log     *logrus.Entry
}

func New{{.ModelName}}() *{{.ModelName}} {
	if {{.VarFieldName}} == nil {
		{{.VarFieldName}} = &{{.ModelName}}{
			service: service.New{{.ModelName}}Service(),
			log:  logger.LOG.WithField("model", "{{.ModelName}}Controller"),
		}
	}
	return {{.VarFieldName}}
}

func (ctrl *{{.ModelName}}) GetList(c *gin.Context) {
	form := validate.{{.ModelName}}ListForm{}
	formErr := validate.{{.ModelName}}ListFormError()
	_ = c.ShouldBind(&form)

	if errData := ctrl.CheckForm(&form, formErr); errData != nil {
		ctrl.ReturnErrorData(c, errData)
		return
	}
	list, total,err := ctrl.service.PageList(&form)
	if err != nil {
	    ctrl.ReturnErrorData(c, err)
		return
	}
	ctrl.ReturnData(c, common.Success, &common.PageResponse{
		TotalCount: total,
		PageSize:   form.PageSize,
		List:       list,
	})
}

func (ctrl *{{.ModelName}}) Add(c *gin.Context) {
    form := validate.{{.ModelName}}AddForm{}
	formErr := validate.{{.ModelName}}AddFormError()
	_ = c.ShouldBind(&form)

	if errData := ctrl.CheckForm(&form, formErr); errData != nil {
		ctrl.ReturnErrorData(c, errData)
		return
	}
	err := ctrl.service.Add(&form)
    if err != nil {
	    ctrl.ReturnErrorData(c, err)
		return
	}
	ctrl.ReturnData(c, common.Success, "")
}

func (ctrl *{{.ModelName}}) Edit(c *gin.Context) {
    form := validate.{{.ModelName}}EditForm{}
	formErr := validate.{{.ModelName}}EditFormError()
	_ = c.ShouldBind(&form)

	if errData := ctrl.CheckForm(&form, formErr); errData != nil {
		ctrl.ReturnErrorData(c, errData)
		return
	}
	err := ctrl.service.Edit(&form)
    if err != nil {
	    ctrl.ReturnErrorData(c, err)
		return
	}
	ctrl.ReturnData(c, common.Success, "")
}


func (ctrl *{{.ModelName}}) Del(c *gin.Context) {
	var pkParam *common.BaseIdParam
     _ = c.ShouldBindQuery(&pkParam)

	var formErr= common.BaseIdParamError()
    if errData := ctrl.CheckForm(&pkParam, formErr); errData != nil {
		ctrl.ReturnErrorData(c, errData)
		return
	}
	err := ctrl.service.Delete(pkParam.Id)
	if err != nil {
	    ctrl.ReturnErrorData(c, err)
		return
	}

	ctrl.ReturnData(c, common.Success, "")
}

func (ctrl *{{.ModelName}}) Info(c *gin.Context) {

	var pkParam *common.BaseIdParam
     _ = c.ShouldBindQuery(&pkParam)
	var formErr= common.BaseIdParamError()
    if errData := ctrl.CheckForm(&pkParam, formErr); errData != nil {
		ctrl.ReturnErrorData(c, errData)
		return
	}

	menu := ctrl.service.FindOne(pkParam.Id)
	ctrl.ReturnData(c, common.Success, menu)
}

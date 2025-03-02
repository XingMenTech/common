package validate

import (
	"github.com/zhangyuanCloud/common"
	{{if .IsTime}} "time" {{end}}
)

type {{.ModelName}}ListForm struct {
	*common.BaseQueryParam
	*common.BaseTimeRequest
}
func {{.ModelName}}ListFormError() map[string]string {
	formError := make(map[string]string)
	formError["Page.required"] = "分页参数必须填写"
	formError["PageSize.required"] = "分页参数必须填写"
	return formError
}

type {{.ModelName}}AddForm struct {
	{{range .Fields}}
	{{if .IsPk}}{{.Name}} {{.Type}} `{{range .Tags }}{{.Name}}:"{{.Value}}" {{end}}` {{ end }}{{ end }}

}

type {{.ModelName}}EditForm struct {
	common.BaseIdParam
	{{.ModelName}}AddForm
}

func {{.ModelName}}AddFormError() map[string]string {
	formError := make(map[string]string)
	return formError
}

func {{.ModelName}}EditFormError() map[string]string {
	formError := {{.ModelName}}AddFormError()
	formError["Id.required"] = "Id缺失"
	return formError
}
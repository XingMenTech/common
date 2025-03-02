package {{.ModuleName}}

import (
    "github.com/gin-gonic/gin"
	"github.com/beego/beego/v2/client/orm"
	"{{.Project}}/pkg/{{.ModuleName}}/models"
	"{{.Project}}/pkg/{{.ModuleName}}/controller"
)
func Init(r *gin.RouterGroup) {
	initModels()
	initRouter(r)
}

func initModels() {
	orm.RegisterModel(
	{{range $model := .Models}}new(models.{{$model}}),{{end}}
	)
}

func initRouter(r *gin.RouterGroup) {
    g := r.Group("/{{.ModuleName}}")
	{{range $model := .Models}}
	{{$model}}Ctl := controller.New{{$model}}()
	g.POST("/{{$model}}/add", {{$model}}Ctl.Add)
	g.POST("/{{$model}}/edit", {{$model}}Ctl.Edit)
	g.GET("/{{$model}}/pageList", {{$model}}Ctl.GetList)
	g.GET("/{{$model}}/find", {{$model}}Ctl.Info)
	g.DELETE("/{{$model}}/delete", {{$model}}Ctl.Del)
	{{end}}
}
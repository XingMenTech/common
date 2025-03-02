// 自动生成模板{{.ModelName}}
package models
import (
    {{if .IsTime}} "time" {{end}}
)

type {{.ModelName}} struct {
      {{range .Fields}}
      {{.Name}} {{.Type}} `{{range .Tags }}{{.Name}}:"{{.Value}}" {{end}}` {{ end }}

}
//设置表名
func (v {{.ModelName}}) TableName() string {
	return "{{.TableName}}"
}
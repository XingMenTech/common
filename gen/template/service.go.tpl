// 自动生成模板{{.ModelName}}
package service

import (
    "{{.Project}}/pkg/{{.ModuleName}}/models"
	"{{.Project}}/pkg/{{.ModuleName}}/repository"
	"{{.Project}}/pkg/{{.ModuleName}}/validate"
	"github.com/sirupsen/logrus"
	"gitlab.novgate.com/common/common/logger"
	"github.com/beego/beego/v2/client/orm"
	"errors"
)

var {{.VarFieldName}}Service *{{.ModelName}}Service

type {{.ModelName}}Service struct {
	repo *repository.{{.ModelName}}Repo
	log  *logrus.Entry
}

func New{{.ModelName}}Service() *{{.ModelName}}Service {
	if {{.VarFieldName}}Service == nil {
		{{.VarFieldName}}Service = &{{.ModelName}}Service{
			repo: repository.New{{.ModelName}}Repo(),
			log:  logger.LOG.WithField("model", "{{.ModelName}}Service"),
		}
	}
	return {{.VarFieldName}}Service
}
func (s *{{.ModelName}}Service) FindOne(id int64) *models.{{.ModelName}} {
	model := &models.{{.ModelName}}{
	    {{.PkField}}: id,
	}
	if err := s.repo.ReadOne(model); err != nil {
		return nil
	}
	return model
}

func (s *{{.ModelName}}Service) PageList(param *validate.{{.ModelName}}ListForm) ([]*models.{{.ModelName}}, int64, error) {
    cond := orm.NewCondition()
    if param.BaseTimeRequest != nil && param.BaseTimeRequest.IsValid() {
    		start, end := param.BaseTimeRequest.GetTime()
    		cond = cond.And("create_time__gte", start).And("create_time__lt", end)
    }
    list := make([]*models.{{.ModelName}}, 0)
    total, err := s.repo.PageList(cond, param.BaseQueryParam, "", &list)
    return list, total, err
}

func (s *{{.ModelName}}Service) Add(form *validate.{{.ModelName}}AddForm) error {
	m := &models.{{.ModelName}}{
		{{range $field := .Fields}}
		{{if $field.IsPk}}{{$field.Name}}: form.{{$field.Name}},{{end}}
		{{end}}
	}
	if _, err := s.repo.InsertOne(nil, m); err != nil {
        return err
    }
	return nil
}

func (s *{{.ModelName}}Service) Edit(form *validate.{{.ModelName}}EditForm) error {

	old := s.FindOne(form.{{.PkField}})
	if old == nil {
		return errors.New("账号信息不存在")
	}
	fields := make([]string, 0)
	{{range $field := .Fields}}
	    {{if $field.IsPk}}
		if old.{{$field.Name}} != form.{{$field.Name}}{
			old.{{$field.Name}} = form.{{$field.Name}}
			fields = append(fields, "{{$field.ColumnName}}")
		}
	    {{end}}
    {{end}}
    return s.repo.Update(nil, old,fields...)

}

func (s *{{.ModelName}}Service) Delete(id int64) error {
    one := s.FindOne(id)
	if one == nil {
		return nil
	}

	return s.repo.Delete(nil, one)
}

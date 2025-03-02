package gen

import "os"

type TemplateModel struct {
	Project      string // 项目名称
	ModelName    string // 模型名称
	VarFieldName string //
	ModuleName   string // 模块名称
	TableName    string // 表名称
	TableSchema  string // 数据库名称
	PkColumn     string // 主键表字段
	PkField      string // 主键结构体字段
	IsTime       bool
	Fields       []*ModelField // 字段集合
}

func (m TemplateModel) getFile(tablePackage string) string {
	path := "pkg/" + m.ModuleName + "/" + tablePackage
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return ""
	}
	return path + "/" + m.ModelName + tablePackage + ".go"
}

type Tag struct {
	Name  string
	Value string
}

type ModelField struct {
	Name       string `json:"Name"`
	Type       string `json:"Type"`
	ColumnName string `json:"columnName"`
	IsPk       bool   `json:"isPk"`
	Tags       []*Tag `json:"Tags"`
	FormTags   []*Tag `json:"FormTags"`
}

type InitTemplateModel struct {
	Project    string
	ModuleName string   // 模块名称
	Models     []string // 结构体集合
}

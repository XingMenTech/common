package gen

import (
	"fmt"
	"strings"
	"text/template"
)

type Table struct {
	TableSchema  string `json:"TABLE_SCHEMA" orm:"column(TABLE_SCHEMA)"`
	TableName    string `json:"TABLE_NAME" orm:"column(TABLE_NAME)"`
	TableComment string `json:"TABLE_COMMENT" orm:"column(TABLE_COMMENT)"`
	ColumnName   string `json:"COLUMN_NAME" orm:"column(COLUMN_NAME)"`
	ModelName    string
	ModuleName   string
	Columns      []*TableColumn
}

type TableColumn struct {
	TableSchema   string `json:"TABLE_SCHEMA" orm:"column(TABLE_SCHEMA)"`
	TableName     string `json:"TABLE_NAME" orm:"column(TABLE_NAME)"`
	ColumnName    string `json:"COLUMN_NAME" orm:"column(COLUMN_NAME)"`
	DataType      string `json:"COLUMN_TYPE" orm:"column(DATA_TYPE)"`
	ColumnComment string `json:"COLUMN_COMMENT" orm:"column(COLUMN_COMMENT)"`
	IsNullable    string `json:"IS_NULLABLE" orm:"column(IS_NULLABLE)"`
}

type GoTpl struct {
	Package string            `json:"package"`
	Tpl     template.Template `json:"tpl"`
}

func (t *Table) getModelName() string {
	if t.ModelName != "" {
		return t.ModelName
	}
	index := strings.Index(t.TableName, "_")
	ms := t.TableName[index+1:]
	t.ModelName = camelString(ms)
	return t.ModelName
}

func (t *Table) getModuleName() string {
	if t.ModuleName != "" {
		return t.ModuleName
	}
	index := strings.Index(t.TableName, "_")
	t.ModuleName = t.TableName[:index]
	return t.ModuleName
}

func (t *Table) BuildModelFields(projectName string) *TemplateModel {
	tmpl := &TemplateModel{
		Project:      projectName,
		ModelName:    t.getModelName(),
		VarFieldName: firstLowerCase(t.ModelName),
		ModuleName:   t.getModuleName(),
		TableName:    t.TableName,
		TableSchema:  t.TableSchema,
		PkColumn:     t.ColumnName,
		PkField:      camelString(t.ColumnName),
	}
	tmpl.Fields, tmpl.IsTime = t.buildField()

	return tmpl
}

func (t *Table) buildField() ([]*ModelField, bool) {

	fields := make([]*ModelField, 0)
	isTime := false
	for _, column := range t.Columns {
		isPk := column.ColumnName != t.ColumnName
		fieldType := camelType(column.DataType)
		isTime = isTime || fieldType == "time.Time"
		fields = append(fields, &ModelField{
			Name:       camelString(column.ColumnName),
			Type:       fieldType,
			ColumnName: column.ColumnName,
			IsPk:       isPk,
			Tags:       column.buildFiledTags(isPk),
			FormTags:   column.buildFormTags(),
		})
	}
	return fields, isTime
}

func (c *TableColumn) buildFiledTags(isPk bool) []*Tag {

	ormTag := &Tag{Name: "orm"}
	if !isPk {
		ormTag.Value = fmt.Sprintf("pk,column(%s)", c.ColumnName)
	} else {
		ormTag.Value = fmt.Sprintf("column(%s)", c.ColumnName)
	}

	return []*Tag{ormTag,
		{
			Name:  "json",
			Value: camelJSONTag(c.ColumnName),
		},
		{
			Name:  "comment",
			Value: c.ColumnComment,
		},
	}
}

func (c *TableColumn) buildFormTags() []*Tag {
	tags := []*Tag{
		{
			Name:  "form",
			Value: camelJSONTag(c.ColumnName),
		},
		{
			Name:  "json",
			Value: camelJSONTag(c.ColumnName),
		},
		{
			Name:  "comment",
			Value: c.ColumnComment,
		},
	}
	if c.IsNullable == "NO" {
		tags = append(tags, &Tag{
			Name:  "binding",
			Value: "required",
		})
	}
	return tags
}

package tool_mongo

import "strings"

type Info struct {
	Name      string              // struct 名称
	TypeName  string              // struct 名称(带包名)
	LowerName string              // 小写
	Fields    []*Field            // 字段
	Package   map[string]struct{} // 需要导入的包
	PkgName   string              // 包路径
	PkField   *Field              // 主键
}

type Field struct {
	DbField   string // 数据库字段名
	CamelName string // 数据库字段名转驼峰命名
	FieldName string // struct字段名
	FieldType string // 字段类型
	IsSlice   bool   // 是否为切片类型
	IsMap     bool   // 是否为map类型
}

var tpl = `package {{{.PkgName}}}

import (
	{{{- range $pkg, $v := .Package}}}
	"{{{$pkg}}}"
	{{{- end}}}
)

const (
	{{{- range .Fields}}}
	{{{$.Name}}}Field{{{.FieldName}}} = "{{{.DbField}}}"
	{{{- end}}}
)

// {{{.Name}}} 数据表内存模型
type {{{.Name}}}Model struct {
	{{{- if .PkField}}}
	{{{- else}}}
	MongoId    primitive.ObjectID ${quote}bson:"_id,omitempty"${quote}
	{{{- end}}}
	{{{.TypeName}}} ${quote}bson:",inline"${quote}
	setter     map[string]interface{} ${quote}bson:"-"${quote}
}

// 获取完整数据
func (m *{{{.Name}}}Model) Object() *{{{.TypeName}}} {
	return &(m.{{{.Name}}})
}

{{{- range .Fields}}}
// 获取 {{{.DbField}}} 的值
func (m *{{{$.Name}}}Model) Get{{{.FieldName}}}() {{{.FieldType}}} {
	return m.{{{$.Name}}}.{{{.FieldName}}}
}
{{{- end}}}

{{{- range .Fields}}}
// 设置 {{{.DbField}}} 的值
func (m *{{{$.Name}}}Model) Set{{{.FieldName}}}({{{.CamelName}}} {{{.FieldType}}}) *{{{$.Name}}}Model {
	m.setter[{{{$.Name}}}Field{{{.FieldName}}}] = {{{.CamelName}}}
	m.{{{.FieldName}}} = {{{.CamelName}}}
	return m
}
{{{- end}}}

// 获取表名
func (m *{{{.Name}}}Model) TableName() string {
	if v, ok := interface{}(m.{{{.Name}}}).(Table); ok {
		return v.TableName()
	}
	return "{{{.LowerName}}}"
}

func New{{{.Name}}}(v *{{{.TypeName}}}) *{{{.Name}}}Model {
	return &{{{.Name}}}Model{
		//MongoId: primitive.ObjectID{},
		{{{.Name}}}:   *v,
		setter: map[string]interface{}{},
	}
}

// 插入新数据
func (m *{{{.Name}}}Model) Insert(ctx context.Context) error {
	result, err := GetCollection(m.TableName()).InsertOne(ctx, m.Object())
	if err != nil {
		return err
	}
	//fmt.Printf("{{{.Name}}}Model Insert result:%+v\n", result)
	{{{- if .PkField}}}
	m.{{{.PkField.FieldName}}} = result.InsertedID.({{{.PkField.FieldType}}})
	{{{- else}}}
	m.MongoId = result.InsertedID.(primitive.ObjectID)
	{{{- end}}}
	return nil
}

// 删除数据
func (m *{{{.Name}}}Model) Delete(ctx context.Context) error {
	{{{- if .PkField}}}
	_, err := GetCollection(m.TableName()).DeleteOne(ctx, bson.D{{"_id", m.{{{.PkField.FieldName}}}}})
	{{{- else}}}
	_, err := GetCollection(m.TableName()).DeleteOne(ctx, bson.D{{"_id", m.MongoId}})
	{{{- end}}}
	if err != nil {
		return err
	}
	//fmt.Printf("{{{.Name}}}Model Delete result:%+v\n", result)
	return nil
}

// 更新数据
func (m *{{{.Name}}}Model) Update(ctx context.Context) error {
	var update bson.D
	for key, value := range m.setter {
		update = append(update, bson.E{
			Key:   key,
			Value: value,
		})
	}
	{{{- if .PkField}}}
	_, err := GetCollection(m.TableName()).UpdateOne(ctx, bson.D{{"_id", m.{{{.PkField.FieldName}}}}}, bson.D{{"$set", update}})
	{{{- else}}}
	_, err := GetCollection(m.TableName()).UpdateOne(ctx, bson.D{{"_id", m.MongoId}}, bson.D{{"$set", update}})
	{{{- end}}}
	if err != nil {
		return err
	}
	//fmt.Printf("{{{.Name}}}Model Update result:%+v\n", result)
	return nil
}

{{{$f := index .Fields 0}}}
// 更新或新增
func (m *{{{.Name}}}Model) UpsertBy{{{if .PkField}}}{{{.PkField.FieldName}}}{{{else}}}{{{$f.FieldName}}}{{{end}}}(ctx context.Context) error {
	result, err := GetCollection(m.TableName()).
		{{{- if .PkField}}}
		UpdateOne(ctx, bson.D{{"_id", m.{{{.PkField.FieldName}}}}}, bson.D{{"$set", m.Object()}}, options.Update().SetUpsert(true))
		{{{- else}}}
		UpdateOne(ctx, bson.D{{"{{{$f.DbField}}}", m.{{{$f.FieldName}}}}}, bson.D{{"$set", m.Object()}}, options.Update().SetUpsert(true))
		{{{- end}}}
	if err != nil {
		return err
	}
	//fmt.Printf("{{{.Name}}}Model Update result:%+v\n", result)
	if result.MatchedCount != 0 {
		//fmt.Println("matched and replaced an existing document")
		return nil
	}
	if result.UpsertedCount != 0 {
		//fmt.Printf("inserted a new document with ID %v\n", result.UpsertedID)
		{{{- if .PkField}}}
		m.{{{.PkField.FieldName}}} = result.UpsertedID.({{{.PkField.FieldType}}})
		{{{- else}}}
		m.MongoId = result.UpsertedID.(primitive.ObjectID)
		{{{- end}}}
		return nil
	}
	return nil
}

{{{- if .PkField}}}
// 根据 {{{.PkField.FieldName}}} 查询
func (m *{{{.Name}}}Model) Get{{{.Name}}}By{{{.PkField.FieldName}}}(ctx context.Context, {{{.PkField.DbField}}} {{{.PkField.FieldType}}}) error {
	return GetCollection(m.TableName()).FindOne(ctx, bson.D{{"_id", {{{.PkField.DbField}}}}}).Decode(m)
}
{{{- else}}}
// 根据 {{{$f.FieldName}}} 查询
func (m *{{{.Name}}}Model) Get{{{.Name}}}By{{{$f.FieldName}}}(ctx context.Context, {{{$f.CamelName}}} {{{$f.FieldType}}}) error {
	return GetCollection(m.TableName()).FindOne(ctx, bson.D{{"{{{$f.DbField}}}", {{{$f.CamelName}}}}}).Decode(m)
}
{{{- end}}}

`

func init() {
	tpl = strings.ReplaceAll(tpl, "${quote}", "`")
}

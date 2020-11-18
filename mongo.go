package tool_mongo

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

/**
 * @param models 需要生成的对象
 * @param pkgName 生成的目录名
 */
func Create(models []interface{}, pkgName string) error {
	for _, m := range models {
		t := reflect.TypeOf(m)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		if t.Kind() != reflect.Struct {
			return errors.New("check type error not struct")
		}

		d := &Info{
			Name:      t.Name(),
			LowerName: strings.ToLower(t.Name()),
			TypeName:  t.String(),
			Package: map[string]struct{}{
				"context":                                    {},
				"go.mongodb.org/mongo-driver/bson":           {},
				"go.mongodb.org/mongo-driver/mongo/options":  {},
				"go.mongodb.org/mongo-driver/bson/primitive": {},
			},
			PkgName: pkgName,
			PkField: nil,
		}

		if !strings.HasSuffix(t.PkgPath(), "/"+pkgName) {
			d.Package[t.PkgPath()] = struct{}{}
		}
		//fieldNum := t.NumField()
		for i := 0; i < t.NumField(); i++ {
			//fmt.Printf("FieldName: %s, TypeString: %s, bson:%s\r\n", t.Field(i).Name, t.Field(i).Type.String(), t.Field(i).Tag.Get("bson"))
			//fmt.Printf("TypeName: %s, KindString: %s\r\n", t.Field(i).Type.Name(), t.Field(i).Type.Kind().String())
			//fmt.Printf("TypeName: %s\r\n", t.Field(i).Type.Name())

			f := &Field{
				DbField:   t.Field(i).Tag.Get("bson"),
				FieldName: t.Field(i).Name,
				FieldType: t.Field(i).Type.String(),
			}
			f.CamelName = ToCamel(f.DbField)
			switch t.Field(i).Type.Kind() {
			case reflect.Map:
				f.IsMap = true
			case reflect.Slice:
				f.IsSlice = true
			}
			if f.DbField == "_id" {
				d.PkField = f
				f.DbField = strings.ToLower(f.FieldName)
				delete(d.Package, "go.mongodb.org/mongo-driver/bson/primitive")
			} else {
				d.Fields = append(d.Fields, f)
			}
		}

		if err := os.MkdirAll("./"+pkgName, 0777); err != nil {
			return err
		}
		file := filepath.Join("./"+pkgName, strings.ToLower(t.Name())+".go")

		if err := WriteFile("model", tpl, file, d); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("√ 成功导出数据表 [%s] 的内存模型 ...\n", t.Name())
		}
	}
	return nil
}

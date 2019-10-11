package minecache

import (
	"fmt"
	"reflect"
	"strings"
)

type StructInfo struct {
	Name       string
	TableName  string
	Fields     []*Field
	PrimaryKey string
}

// 根据where和表达式生成redis key尾部
func getRedisKeyTail(where ...interface{}) (bool, string) {
	key := ""
	for i := 0; i < len(where); i++ {
		tempWhere := where[i]
		t := reflect.TypeOf(tempWhere)
		if t.Kind() == reflect.Ptr {
			v := reflect.ValueOf(tempWhere)

			t = v.Type()
			tempWhere = v.Elem()
		}

		switch t.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Bool, reflect.String,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			key += fmt.Sprintf(":%v", tempWhere)
		default:
			return false, ""
		}
	}
	return true, key
}

func getStructInfo(value interface{}) *StructInfo {
	structInfo := &StructInfo{
		Fields: make([]*Field, 0),
	}

	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)
	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}

	structInfo.Name = t.Name()
	structInfo.TableName = snakeString(structInfo.Name)

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		fv := v.Field(i)

		fieldInfo := &Field{}
		fieldInfo.Name = sf.Name
		fieldInfo.ColumnName = snakeString(fieldInfo.Name)
		fieldInfo.IsIgnore, fieldInfo.IsPrimaryKey = dealTag(sf.Tag)
		fieldInfo.Kind, fieldInfo.Value = getVal(fv)
		if fieldInfo.IsPrimaryKey {
			structInfo.PrimaryKey = fieldInfo.Name
		}

		structInfo.Fields = append(structInfo.Fields, fieldInfo)
	}
	return structInfo
}

func dealTag(tag reflect.StructTag) (bool, bool) {
	orm := tag.Get("orm")
	if orm == "-" {
		return true, false
	}

	ormTags := strings.Split(orm, ",")
	for i := 0; i < len(ormTags); i++ {
		if ormTags[i] == "pk" {
			return false, true
		}
	}
	return false, false
}

func getVal(value reflect.Value) (reflect.Kind, interface{}) {
	tempVal := value
	if value.Kind() == reflect.Ptr {
		v := value.Elem()

		tempVal = v
	}

	if value.Kind() == reflect.Struct {
		fmt.Println(tempVal, ">>>", value.Kind())
		return reflect.Struct, getStructInfo(tempVal)
	}

	if value.Kind() == reflect.String && tempVal.String() == "" {
		return reflect.String, nil
	}

	if value.Kind() == reflect.Int64 && tempVal.Int() == 0 {
		return reflect.Int64, nil
	}

	return value.Kind(), tempVal.Interface()
}

func snakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

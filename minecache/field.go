package minecache

import "reflect"

// Field model field definition
type Field struct {
	Name         string
	ColumnName   string
	Kind         reflect.Kind
	Value        interface{}
	IsIgnore     bool
	IsPrimaryKey bool
}

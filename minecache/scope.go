package minecache

import (
	"fmt"
	"reflect"
	"strings"
)

type scope struct {
	sql        string
	sqlVars    []interface{}
	redisKey   string
	primaryKey string
	//redisValue string // json
	//primaryKeyFields []*Field
	//fields          []*Field
	//selectAttrs     *[]string
}

func getUpdateScope(value interface{}, expression string, where ...interface{}) *scope {
	valid, tailKey := getRedisKeyTail(where)
	if !valid {
		return nil
	}

	structInfo := getStructInfo(value)

	setParams := make([]interface{}, 0)
	setSql := make([]string, 0)
	for i := 0; i < len(structInfo.Fields); i++ {
		if structInfo.Fields[i].Kind == reflect.Struct {
			fmt.Println("struct is not support")
			continue
		}
		if structInfo.Fields[i].IsIgnore {
			continue
		}

		if structInfo.Fields[i].Value == nil {
			continue
		}

		setSql = append(setSql, fmt.Sprintf("%s=?", structInfo.Fields[i].ColumnName))
		setParams = append(setParams, structInfo.Fields[i].Value)
	}

	if len(setSql) == 0 {
		return nil
	}

	setStr := strings.Replace(strings.Trim(fmt.Sprint(setSql), "[]"), " ", ",", -1)
	sql := fmt.Sprintf("update %s set %s where %s", structInfo.TableName, setStr, expression)
	key := fmt.Sprintf("project:%s:%s%s:%%d", structInfo.TableName, expression, tailKey)

	sp := &scope{}
	sp.redisKey = key
	sp.sql = sql
	sp.primaryKey = structInfo.PrimaryKey
	sp.sqlVars = setParams
	sp.sqlVars = append(sp.sqlVars, where...)
	return sp
}

func getAddScope(value interface{}) *scope {
	structInfo := getStructInfo(value)

	sqlAttr := make([]string, 0)
	sqlValue := make([]interface{}, 0)
	keys := make([]string, 0)
	for i := 0; i < len(structInfo.Fields); i++ {
		if structInfo.Fields[i].Value == nil {
			continue
		}
		if structInfo.Fields[i].Kind == reflect.Struct {
			fmt.Println("struct is not support")
			continue
		}
		if structInfo.Fields[i].IsIgnore {
			continue
		}

		sqlAttr = append(sqlAttr, structInfo.Fields[i].ColumnName)
		sqlValue = append(sqlValue, structInfo.Fields[i].Value)

		if structInfo.Fields[i].IsPrimaryKey {
			keys = append(keys, structInfo.Fields[i].ColumnName)
		}
	}
	if len(sqlAttr) == 0 {
		return nil
	}

	placeholders := ""
	for i := 0; i < len(sqlAttr); i++ {
		placeholders += "?"
		if len(sqlAttr)-1 != i {
			placeholders += ","
		}
	}
	attrStr := strings.Replace(strings.Trim(fmt.Sprint(sqlAttr), "[]"), " ", ",", -1)
	sql := fmt.Sprintf("insert into %s(%s) values(%s)", structInfo.TableName, attrStr, placeholders)
	keyStr := strings.Replace(strings.Trim(fmt.Sprint(keys), "[]"), " ", ",", -1)

	sp := &scope{}
	sp.redisKey = fmt.Sprintf("project:%s%s:%%d", structInfo.TableName, keyStr)
	sp.sql = sql
	sp.sqlVars = sqlValue
	sp.primaryKey = structInfo.PrimaryKey
	return sp
}

func getDeleteScope(value interface{}, expression string, where ...interface{}) *scope {
	valid, tailKey := getRedisKeyTail(where)
	if !valid {
		return nil
	}

	structInfo := getStructInfo(value)

	sp := &scope{}
	sp.redisKey = fmt.Sprintf("project:%s%s:%%d", structInfo.TableName, tailKey)
	sp.sql = fmt.Sprintf("delete from %s where %s", structInfo.TableName, expression)
	sp.sqlVars = where
	sp.primaryKey = structInfo.PrimaryKey
	return sp
}

func getFindScope(value interface{}, expression string, where ...interface{}) *scope {
	valid, tailKey := getRedisKeyTail(where)
	if !valid {
		return nil
	}

	structInfo := getStructInfo(value)

	sqlAttr := make([]string, 0)
	for i := 0; i < len(structInfo.Fields); i++ {
		if structInfo.Fields[i].Value == nil {
			continue
		}
		if structInfo.Fields[i].Kind == reflect.Struct {
			fmt.Println("struct is not support")
			continue
		}
		if structInfo.Fields[i].IsIgnore {
			continue
		}

		sqlAttr = append(sqlAttr, structInfo.Fields[i].ColumnName)
	}

	attrStr := strings.Replace(strings.Trim(fmt.Sprint(sqlAttr), "[]"), " ", ",", -1)

	sp := &scope{}
	sp.redisKey = fmt.Sprintf("project:%s%s:%%d", structInfo.TableName, tailKey)
	sp.sql = fmt.Sprintf("select %s from %s where %s", attrStr, structInfo.TableName, expression)
	sp.sqlVars = where
	sp.primaryKey = structInfo.PrimaryKey
	return sp
}

func getFindFirstScope(value interface{}, expression string, where ...interface{}) *scope {
	sp := getFindScope(value, expression, where...)
	if sp == nil {
		return nil
	}

	sp.sql += " limit 1"
	return sp
}

package insert

import (
	"errors"
	"reflect"
	"strings"
)

var errInvalidEntity = errors.New("invalid entity")

func InsertStmt(entity interface{}) (string, interface{}, error) {
	if entity == nil {
		return "", nil, errInvalidEntity
	}
	val := reflect.ValueOf(entity)
	typ := val.Type()
	// 如果是非结构体返回""
	if typ.Kind() != reflect.Struct {
		return "", nil, errInvalidEntity
	}
	// 如果是一级指针
	if typ.Kind() == reflect.Pointer {
		val = val.Elem()
		typ = val.Type()
	}
	bd := strings.Builder{}
	_, _ = bd.WriteString("INSERT INTO `")
	bd.WriteString(typ.Name())
	bd.WriteString("`(")
	fields,fieldVals := fieldNameAndValues(val)
	for i, name := range fields {
		if i > 0 {
			bd.WriteRune(',')
		}
		bd.WriteRune('`')
		bd.WriteString(name)
		bd.WriteRune('`')
	}
	bd.WriteString(") VALUES (")

	args := make([]interface{}, 0, len(fieldVals))
	for i,fd := range fields {
		if i > 0 {
			bd.WriteRune(',')
		}
		bd.WriteRune('?')
		args = append(args,fieldVals[fd])
	}
	if len(args) == 0 {
		return "", nil, errInvalidEntity
	}
	bd.WriteString(");")
	return bd.String(), args, nil
}

func fieldNameAndValues(val reflect.Value) ([]string, map[string]interface{}) {

	typ := val.Type()
	fieldNum := val.NumField()
	fields := make([]string, 0, fieldNum)
	values := make(map[string]interface{},fieldNum)
	for i := 0; i < fieldNum; i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		// Anonymous 判断是否是匿名结构体
		if field.Type.Kind() == reflect.Struct && field.Anonymous {
			subField, subValus := fieldNameAndValues(fieldVal)
			for _,k := range subField {
				if _,ok := values[k]; ok {
					// 重复字段，只会出现在组合情况下，直接忽略重复字段
					continue
				}
				fields = append(fields, k)
				values[k] = subValus[k]
			}
		}
		fields = append(fields,field.Name)
		values[field.Name] = fieldVal.Interface()
	}
	return fields, values
}

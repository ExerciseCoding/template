package reflect

import (
	"errors"
	"reflect"
)

func IterateFields(entity any) (map[string]any, error) {
	if entity == nil {
		return nil, errors.New("不支持 nil")
	}
	typ := reflect.TypeOf(entity)
	val := reflect.ValueOf(entity)
	if val.IsZero() {
		// 零值：对象占据的bit空间全是零
		return nil, errors.New("不支持零值")
	}
	for typ.Kind() == reflect.Pointer {
		// Elem 方法如果typ是指针，Elem会取出指针指向的值
		typ = typ.Elem()
		val = val.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("不支持类型")
	}

	numField := typ.NumField()
	res := make(map[string]any, numField)
	for i := 0; i < numField; i++ {
		// 字段的类型
		fieldType := typ.Field(i)
		// 字段的值
		fieldValue := val.Field(i)
		// 判断字段是否是可以导出的(私有字段)
		if fieldType.IsExported() {
			res[fieldType.Name] = fieldValue.Interface()
		} else {
			res[fieldType.Name] = reflect.Zero(fieldValue.Type()).Interface()
		}
	}
	return res, nil
}


func SetField(entity any, field string, newValue any) error {
	val := reflect.ValueOf(entity)
	for val.Type().Kind() == reflect.Pointer {
		// Elem方法是reflect包中的一个方法，它返回接口持有的值的reflect.Value，或者指针、数组、切片、字典的元素的reflect.Value。
		val = val.Elem()
	}

	fieldValue := val.FieldByName(field)
	// CanSet字符安可不可以被修改
	if !fieldValue.CanSet() {
		return errors.New("不可修改字段")
	}
	fieldValue.Set(reflect.ValueOf(newValue))
	return nil
}
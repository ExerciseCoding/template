package reflect_demo

import (
	"errors"
	"fmt"
	"reflect"
)

// 这样写没法测试断言，因为没有返回值
//func IterateFields(val any) {
//
//}

func IterateFields(val any) {
	res, err := iterateFields(val)
	if err != nil {
		fmt.Println(err)
		return
	}
	for k, v := range res {
		fmt.Println(k, v)
	}
}

func iterateFields(val any) (map[string]any, error) {
	if val == nil {
		return nil, errors.New("不能为nil")
	}
	// 获取字段的type信息
	typ := reflect.TypeOf(val)
	// 获取字段的val信息
	refVal := reflect.ValueOf(val)
	// 一级指针
	//if typ.Kind() == reflect.Ptr {
	//	typ = typ.Elem()
	//	refVal = refVal.Elem()
	//}

	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		refVal = refVal.Elem()
	}
	numField := typ.NumField()
	res := make(map[string]any, numField)
	for i := 0; i < numField; i++ {
		fdType := typ.Field(i)
		res[fdType.Name] = refVal.Field(i).Interface()
	}
	return res, nil
}

func SetField(entity any, field string, newVal any) error {
	val := reflect.ValueOf(entity)
	typ := val.Type()
	// 只能是一级指针，类似*User
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Ptr {
		return errors.New("非法类型")
	}

	typ = typ.Elem()
	val = val.Elem()
	// 这个地方判断不出来 field在不在
	fd := val.FieldByName(field)
	// 利用type来判断field在不在
	if _, found := typ.FieldByName(field); !found {
		return errors.New("字段不存在")
	}
	if !fd.CanSet() {
		return errors.New("不可修改字段")
	}
	fd.Set(reflect.ValueOf(newVal))
	return nil
}

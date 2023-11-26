package entity_copy_myself

import (
	"errors"
	"fmt"
	"reflect"
)

func copy(input any) (any, error) {
	if input == nil {
		return nil, errors.New("输入不能为nil")
	}
	typ := reflect.TypeOf(input)
	val := reflect.ValueOf(input)
	if typ.Kind() != reflect.Struct && (typ.Kind() == reflect.Ptr && typ.Elem().Kind() != reflect.Struct) {
		return nil, errors.New("非法类型")
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	numsField := typ.NumField()
	//copyEntity := reflect.New(typ.Elem()).Interface()

	for i := 0; i < numsField; i++ {
		valField := val.Field(i)
		fmt.Println(valField)
	}

	return nil, nil
}

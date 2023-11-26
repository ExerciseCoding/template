package reflect_demo

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"template/internal/reflect_demo/types"
	"testing"
)

func TestReflectUser(t *testing.T) {
	r := reflect.TypeOf(&User{})
	if r.Kind() == reflect.Struct {
		fmt.Println(r.NumField())
	} else if r.Kind() == reflect.Ptr {
		fmt.Println("指针")
	}
}

type User struct {
	Name string
}

func TestIterateFields(t *testing.T) {
	u1 := &User{
		Name: "haohao",
	}
	u2 := &u1
	tests := []struct {
		// 名字
		name string
		//输入部署
		val any

		//输出部分
		wantRes map[string]any
		wantErr error
	}{
		{
			name:    "nil",
			val:     nil,
			wantErr: errors.New("不能为nil"),
		},
		{
			name: "user",
			val: User{
				Name: "Tom",
			},
			wantErr: nil,
			wantRes: map[string]any{
				"Name": "Tom",
			},
		},
		{
			// 指针
			name:    "pointer",
			val:     &User{Name: "Jerry"},
			wantErr: nil,
			wantRes: map[string]any{
				"Name": "Jerry",
			},
		},
		{
			// 指针
			name:    "multiple pointer",
			val:     u2,
			wantErr: nil,
			wantRes: map[string]any{
				"Name": "haohao",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := iterateFields(tt.val)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}

			assert.Equal(t, tt.wantRes, res)
		})
	}
}

func TestSetField(t *testing.T) {
	testCases := []struct {
		name string

		field  string
		entity any

		newVal  any
		wantErr error
	}{
		{
			name:    "struct",
			entity:  types.User{},
			field:   "Name",
			wantErr: errors.New("非法类型"),
		},
		{
			name:    "private field",
			entity:  &types.User{},
			field:   "age",
			wantErr: errors.New("不可修改字段"),
		},
		{
			name:    "invalid field",
			entity:  &types.User{},
			field:   "invalid_field",
			wantErr: errors.New("字段不存在"),
		},
		{
			name:   "pass",
			entity: &types.User{},
			field:  "Name",
			newVal: "Tom",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := SetField(tc.entity, tc.field, tc.newVal)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

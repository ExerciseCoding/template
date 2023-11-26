package entity_copy_myself

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	//"template/internal/homework/entity_copy_myself/types"
	"testing"
)

func TestCopy(t *testing.T) {
	user := types.User{
		Name: "ming",
	}
	typ := reflect.TypeOf(user)
	val := reflect.ValueOf(user)
	fmt.Println(typ, val)

	user1 := &types.User{
		Name: "ning",
	}
	typ = reflect.TypeOf(user1)
	val = reflect.ValueOf(user1)
	fmt.Println(typ, typ.Elem(), val, val.Elem())
}

func TestCopyCase(t *testing.T) {
	user := &types.User{}
	testCases := []struct {
		name string

		//输入
		entity any

		//输出
		copyEntity any
		wantErr    error
	}{
		{
			name:       "二级指针",
			entity:     &user,
			copyEntity: nil,
			wantErr:    errors.New("非法类型"),
		},
		{
			name:       "二级指针",
			entity:     &types.User{Name: "lihao"},
			copyEntity: nil,
			wantErr:    errors.New("非法类型"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := copy(tc.entity)
			assert.Equal(t, tc.wantErr, err)

		})
	}
}

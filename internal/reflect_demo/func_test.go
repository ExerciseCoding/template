package reflect_demo

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

// 考虑： nil, 基本类型，内置类型(切片，map，channel)之类
// 结构体，结构体指针，多级指针

// 决策支持：结构体、或者结构体指针
func TestIterateFunc(t *testing.T) {
	type args struct {
		val any
	}

	tests := []struct {
		name string

		args args

		want    map[string]*FuncInfo
		wantErr error
	}{
		{
			name:    "nil",
			wantErr: errors.New("输入 nil"),
		},
		{
			name:    "basic types",
			wantErr: errors.New("不支持类型"),
		},
		{
			name: "basic types",
			args: args{
				val: 123,
			},
			wantErr: errors.New("不支持类型"),
		},
		{
			name: "struct typs",
			args: args{
				val: Order{
					buyer: 18,
				},
			},
			want: map[string]*FuncInfo{
				"GetBuyer": {
					Name:   "GetBuyer",
					In:     []reflect.Type{reflect.TypeOf(Order{})},
					Out:    []reflect.Type{reflect.TypeOf(int64(0))},
					Result: []any{int64(18)},
				},
				// 私有方法在不同包里反射也是拿不到的
				"getSeller": {
					Name:   "GetBuyer",
					In:     []reflect.Type{reflect.TypeOf(Order{})},
					Out:    []reflect.Type{reflect.TypeOf(int64(0))},
					Result: []any{int64(18)},
				},
			},
		},

		{
			name: "ptr typs",
			args: args{
				val: &OrderV1{
					buyer: 18,
				},
			},
			want: map[string]*FuncInfo{
				"GetBuyer": {
					Name:   "GetBuyer",
					In:     []reflect.Type{reflect.TypeOf(&OrderV1{})},
					Out:    []reflect.Type{reflect.TypeOf(int64(0))},
					Result: []any{int64(18)},
				},
			},
		},
		// 指针可以访问到结构体方法
		// 结构体访问不到指针方法
		{
			name: "pointer type but input struct",
			args: args{
				val: OrderV1{
					buyer: 18,
				},
			},
			want: map[string]*FuncInfo{
				"GetBuyer": {
					Name:   "GetBuyer",
					In:     []reflect.Type{reflect.TypeOf(&OrderV1{})},
					Out:    []reflect.Type{reflect.TypeOf(int64(0))},
					Result: []any{int64(18)},
				},
			},
		},

		{
			name: "struct type but input ptr",
			args: args{
				val: &OrderV1{
					buyer: 18,
				},
			},
			want: map[string]*FuncInfo{
				"GetBuyer": {
					Name:   "GetBuyer",
					In:     []reflect.Type{reflect.TypeOf(OrderV1{})},
					Out:    []reflect.Type{reflect.TypeOf(int64(0))},
					Result: []any{int64(18)},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IterateFunc(tt.args.val)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

type Order struct {
	buyer  int64
	seller int64
}

// 在反射层面上： func GetBuyer(o Order) int64 {
func (o Order) GetBuyer() int64 {
	return o.buyer
}

func (o Order) getSeller() int64 {
	return o.seller
}

// 指针
//func (o *Order) GetBuyer() int64 {
//	return o.buyer
//}

type OrderV1 struct {
	buyer int64
}

// 在反射层面上： func GetBuyer(o Order) int64 {
func (o *OrderV1) GetBuyer() int64 {
	return o.buyer
}

type MyInterface interface {
	Abc()
}

//var _ MyInterface = abcImpl{}

var _ MyInterface = &abcImpl{}

type abcImpl struct {
}

func (a *abcImpl) Abc() {
	panic("")
}

type MyService struct {
	GetById func()
}

func Proxy() {
	myService := &MyService{}
	myService.GetById = func() {
		// 发起RPC调用
		// 解析请求
	}
}

// 数据库直接对应
type UserEntity struct {
	Name string
}

func Copy(src any, dst any) error {
	// 反射操作，一个个字段复制
	return nil
}

// ignoreFields 忽略一些字段不复制
func CopyV1(src any, dst any, ignoreFields ...string) error {
	// 反射操作，一个个字段复制
	return nil
}

type Copier struct {
	src          any
	dst          any
	ignoreFields []string
}

func (c Copier) Copy() error {
	return nil
}

type UserV1 struct {
	Name  string
	Email string
}

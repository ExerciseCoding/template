package reflect_demo

import (
	"errors"
	"io"
	"reflect"
)

// 反射中，一个实例可以看成两部分
// - 值
// - 实际类型

func IterateFunc(val any) (map[string]*FuncInfo, error) {
	//w := &bytes.Buffer{}
	//Print(w, "aaa")
	//w.String() = "aaa"
	if val == nil {
		return nil, errors.New("输入 nil")
	}

	typ := reflect.TypeOf(val)
	//if typ.Kind() == reflect.Ptr {
	//	typ = typ.Elem()
	//}

	// 不是结构体也不是一级指针
	if typ.Kind() != reflect.Struct &&
		!(typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct) {
		return nil, errors.New("不支持类型")
	}

	//获取方法数量
	numMethod := typ.NumMethod()
	// 遍历方法
	res := make(map[string]*FuncInfo, numMethod)
	for i := 0; i < numMethod; i++ {
		// method :=
		method := typ.Method(i)
		mt := method.Type
		numIn := mt.NumIn()

		in := make([]reflect.Type, 0, numIn)
		for j := 0; j < numIn; j++ {
			in = append(in, mt.In(j))
		}

		numOut := mt.NumOut()
		out := make([]reflect.Type, 0, numOut)
		for j := 0; j < numOut; j++ {
			out = append(out, mt.Out(j))
		}

		callRes := method.Func.Call([]reflect.Value{reflect.ValueOf(val)})

		retVals := make([]any, 0, len(callRes))
		for _, cr := range callRes {
			retVals = append(retVals, cr.Interface())
		}
		res[method.Name] = &FuncInfo{
			Name:   method.Name,
			In:     in,
			Out:    out,
			Result: retVals,
		}
	}

	return res, nil

}

func Print(writer io.Writer, msg string) {

}

type FuncInfo struct {
	Name string
	In   []reflect.Type
	Out  []reflect.Type

	//反射调用得到的结果
	Result []any
}

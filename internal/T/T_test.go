package T

// 泛型方法：
// 1.形式在方法名后面加类型参数，FuncName[T constraint]
// 2.类型参数可以有多个
// 2.constraint是类型约束，必须是接口。可以是普通接口，也可以是Numeric这种复合接口

//type Numeric interface {
//	type int, int64, int32
//}
//
//func Sum[T Numeric](values []T) T {
//	var res T
//	for _, val := range values {
//		res = res + val
//	}
//	return res
//}


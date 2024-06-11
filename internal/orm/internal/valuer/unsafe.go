package valuer

import (
	"database/sql"
	"github.com/ExerciseCoding/template/internal/orm/internal/errs"
	"github.com/ExerciseCoding/template/internal/orm/model"
	"reflect"
	"unsafe"
)

type unsafeValue struct {
	model   *model.Model
	val     any
	address unsafe.Pointer
}

// 此处的用作是防止有地方修改了签名
// 例如：func NewUnsafeValue(model *orm.Model, val any, a int) Value {}
var _ Creator = NewUnsafeValue

func NewUnsafeValue(model *model.Model, val any) Value {
	address := reflect.ValueOf(val).UnsafePointer()
	return unsafeValue{
		model:   model,
		val:     val,
		address: address,
	}
}

func (u unsafeValue) SetColums(rows *sql.Rows) error {
	// 怎么知道SELECT 出来了哪些列
	// 拿到SELECT出来的列
	cs, err := rows.Columns()
	if err != nil {
		return err
	}

	// 怎么利用cs来解决顺序问题和类型问题

	vals := make([]any, 0, len(cs))
	//valsElem := make([]reflect.Value, 0, len(cs))

	// reflect.ValueOf(tp).Pointer() Pointer()方法返回的是一个uintptr类型的值，这个值表示一个地址，但它不是一个真正的指针。你不能对这个值进行解引用操作，也不能使用它来访问内存。Pointer()方法只能用于Chan, Func, Map, Ptr, UnsafePointer, Slice类型的值。
	// reflect.ValueOf(tp).UnsafePointer() UnsafeAddr()方法返回的也是一个uintptr类型的值，这个值表示一个地址，但它不是一个真正的指针。你不能对这个值进行解引用操作，也不能使用它来访问内存。UnsafeAddr()方法可以用于任何类型的值，它返回的是值的地址，而不是值本身的地址。

	// reflect.ValueOf(tp)取到的是tp的值也就是tp指向的Mystruct, 然后UnsafePointer取到的就是这个值得地址也就是MyStruct的地址

	//reflect.ValueOf(t).Pointer()：
	//
	//这行代码将 t 转换为 reflect.Value 类型，并使用 Pointer() 方法获取 t 变量所指向的 MyStruct 实例的地址。
	//结果是一个 uintptr 类型的整数，表示 MyStruct 实例的地址。
	//reflect.ValueOf(t).UnsafePointer()：
	//
	//这行代码将 t 转换为 reflect.Value 类型，并使用 UnsafePointer() 方法获取 t 变量所指向的 MyStruct 实例的地址。
	//结果是一个 uintptr 类型的整数，表示 MyStruct 实例的地址。

	for _, c := range cs {
		fd, ok := u.model.ColumnMap[c]
		if !ok {
			return errs.NewErrUnkownColumn(c)
		}
		// reflect.ValueOf(&myStruct)
		fdAddress := unsafe.Pointer(uintptr(u.address) + fd.Offset)
		// 反射创建一个实例,反射在特点的地址上，创建一个特定类型的实例
		// 这里创建的实例是原本类型的指针类型
		// 例如: fd.Type = int, 那么val 是 *int
		val := reflect.NewAt(fd.Type, fdAddress)
		vals = append(vals, val.Interface())
	}
	err = rows.Scan(vals...)
	return err
}

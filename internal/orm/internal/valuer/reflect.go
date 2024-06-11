package valuer

import (
	"database/sql"
	"github.com/ExerciseCoding/template/internal/orm/internal/errs"
	"github.com/ExerciseCoding/template/internal/orm/model"
	"reflect"
)

type reflectValue struct {
	model *model.Model
	val   any
}

var _Creator = NewReflectValue

func NewReflectValue(model *model.Model, val any) Value {
	return reflectValue{
		model: model,
		val:   val,
	}
}
func (r reflectValue) SetColums(rows *sql.Rows) error {
	// 怎么知道SELECT 出来了哪些列
	// 拿到SELECT出来的列
	cs, err := rows.Columns()
	if err != nil {
		return err
	}

	// 怎么利用cs来解决顺序问题和类型问题

	// 通过cs 来构造vals
	vals := make([]any, 0, len(cs))
	valsElem := make([]reflect.Value, 0, len(cs))
	for _, c := range cs {
		fd, ok := r.model.ColumnMap[c]
		if !ok {
			return errs.NewErrUnkownColumn(c)
		}
		val := reflect.New(fd.Type)
		vals = append(vals, val.Interface())
		// 记得调用Elem, 因为fd.Type = int, 那么val是*int
		valsElem = append(valsElem, val.Elem())
		//for _, fd :=range s.model.fieldMap {
		//	if fd.colName == c {
		//		// 反射创建一个实例
		//		// 这里创建的实力是原本类型的指针类型
		//		// 例如：fd.Type = int, 那么val 是 *int指针
		//		val := reflect.New(fd.typ)
		//		vals = append(vals, val.Interface())
		//	}
		//}
	}

	err = rows.Scan(vals...)
	if err != nil {
		return err
	}

	// 想办法将vals 塞进去 结果tp里面
	tpValue := reflect.ValueOf(r.val).Elem()

	for i, c := range cs {
		fd, ok := r.model.ColumnMap[c]
		if !ok {
			return errs.NewErrUnkownColumn(c)
		}

		tpValue.FieldByName(fd.GoName).Set(valsElem[i])
		//for _, fd := range s.model.fieldMap {
		//	if fd.colName == c {
		//		tpValue.Elem().FieldByName(fd.goName).Set(reflect.ValueOf(vals[i]).Elem())
		//	}
		//}
	}
	return nil
}

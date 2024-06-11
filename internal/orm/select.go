package orm

import (
	"context"
	_ "database/sql"
	"github.com/ExerciseCoding/template/internal/orm/internal/errs"
	model2 "github.com/ExerciseCoding/template/internal/orm/model"
	"reflect"
	"strings"
	"unsafe"
)

// Selectable 是一个标记接口
// 它代表的是查找的列，或者聚合函数等
// SELECT XXX 部分
type Selectable interface {
	selectable()
}


type Selector[T any] struct {
	table string
	model *model2.Model
	where []Predicate
	sb    *strings.Builder
	args  []any
	columns []Selectable
	db *DB
	//r *registry
}

func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		sb: &strings.Builder{},
		db: db,
	}
}

func (s *Selector[T]) BuildOld() (*Query, error) {

	var sb strings.Builder
	sb.WriteString("SELECT * FROM ")
	if s.table == "" {
		var t T
		typ := reflect.TypeOf(t)
		sb.WriteByte('`')
		sb.WriteString(typ.Name())
		sb.WriteByte('`')
	} else {
		//sb.WriteByte('`')
		sb.WriteString(s.table)
		//sb.WriteByte('`')
	}
	args := make([]any, 0, 4)
	if len(s.where) > 0 {
		sb.WriteString(" WHERE ")
		p := s.where[0]
		for i := 1; i < len(s.where); i++ {
			p.And(s.where[i])
		}
		// 在这里处理 p
		// p.left 构建好
		// p.op 构建好
		// p.right 构建好
		switch left := p.left.(type) {
		case Column:
			sb.WriteByte('`')
			sb.WriteString(left.name)
			sb.WriteByte('`')
		}
		sb.WriteString(p.op.String())
		switch right := p.right.(type) {
		case Value:
			sb.WriteByte('?')
			args = append(args, right.val)
		}
	}
	sb.WriteByte(';')
	return &Query{
		SQL:  sb.String(),
		Args: args,
	}, nil
}

// Build BuildOld的改造版
func (s *Selector[T]) Build() (*Query, error) {
	s.sb = &strings.Builder{}
	var err error
	s.model, err = s.db.r.Get(new(T))
	if err != nil {
		return nil, err
	}

	sb := s.sb
	// sb.WriteString("SELECT * FROM ")
	sb.WriteString("SELECT ")

	if err := s.buildColumns(); err != nil {
		return nil, err
	}


	sb.WriteString(" FROM ")
	if s.table == "" {
		sb.WriteByte('`')
		sb.WriteString(s.model.TableName)
		sb.WriteByte('`')
	} else {
		//sb.WriteByte('`')
		sb.WriteString(s.table)
		//sb.WriteByte('`')
	}
	if len(s.where) > 0 {
		sb.WriteString(" WHERE ")
		p := s.where[0]
		for i := 1; i < len(s.where); i++ {
			p = p.And(s.where[i])
		}
		if err := s.buildExpression(p); err != nil {
			return nil, err
		}
	}
	sb.WriteByte(';')
	return &Query{
		SQL:  sb.String(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) buildColumns() error {
	if len(s.columns) == 0 {
		// 没有指定列
		s.sb.WriteByte('*')
	}
	for i, col := range s.columns {
		if i > 0 {
			s.sb.WriteByte(',')
		}
		switch c := col.(type) {
		case Column:
			err := s.buildColumn(c.name)
			if err != nil {
				return err
			}
		case Aggregate:
			s.sb.WriteString(c.fn)
			s.sb.WriteByte('(')
			err := s.buildColumn(c.arg)
			if err != nil {
				return err
			}
			s.sb.WriteByte(')')
		}
	}

	return nil
}

func (s *Selector[T]) buildExpression(expr Expression) error {
	switch exp := expr.(type) {
	case nil:
	case Predicate:
		// 在这里处理 p
		// p.left 构建好
		// p.op 构建好
		// p.right 构建好
		_, ok := exp.left.(Predicate)
		if ok {
			s.sb.WriteByte('(')
		}

		if err := s.buildExpression(exp.left); err != nil {
			return err
		}

		if ok {
			s.sb.WriteByte(')')
		}
		s.sb.WriteByte(' ')
		s.sb.WriteString(exp.op.String())
		s.sb.WriteByte(' ')

		_, ok = exp.right.(Predicate)
		if ok {
			s.sb.WriteByte('(')
		}

		if err := s.buildExpression(exp.right); err != nil {
			return err
		}

		//switch left := expr.left.(type) {
		//case Column:
		//	sb.WriteByte('`')
		//	sb.WriteString(left.name)
		//	sb.WriteByte('`')
		//}
		//
		//switch right := expr.right.(type) {
		//case Value:
		//	sb.WriteByte('?')
		//	args = append(args, right.val)
		//}
		if ok {
			s.sb.WriteByte(')')
		}

	case Column:
		s.buildColumn(exp.name)

	case Value:
		s.sb.WriteByte('?')
		s.addArgs(exp.val)

	default:
		return errs.NewErrUnsupportedExpression(expr)
	}
	return nil
}

func (s *Selector[T]) buildColumn(c string) error {
	fd, ok := s.model.FieldMap[c]
	// 字段不对，或者列不对
	if !ok {
		return errs.NewErrUnkownField(c)
	}
	s.sb.WriteByte('`')
	s.sb.WriteString(fd.ColName)
	s.sb.WriteByte('`')
	return nil
}

func (s *Selector[T]) addArgs(val any) {
	if s.args == nil {
		s.args = make([]any, 0, 8)
	}
	s.args = append(s.args, val)
}


//func (s *Selector[T]) SelectV1(cols string) *Selector[T] {
//
//}
//
//func (s *Selector[T]) Select(cols ...string) *Selector[T] {
//	s.column = cols
//	return s
//}

func (s *Selector[T]) Select(cols ...Selectable) *Selector[T] {
	s.columns = cols
	return s
}

func (s *Selector[T]) From(table string) *Selector[T] {
	s.table = table
	return s
}

func (s *Selector[T]) GetV2(ctx context.Context) (*T, error) {
	q, err := s.Build()
	if err != nil {
		return nil, err
	}

	db := s.db.db
	// 发起查询，并且处理结果集
	rows, err := db.QueryContext(ctx, q.SQL, q.Args...)

	// 这个是查询错误
	if err != nil {
		return nil, err
	}
	// 继续处理结果集

	// 确认是否有数据
	if !rows.Next() {
		// 是否返回error
		// 返回error,和sql包语义保持一致
		return nil, errs.ErrNoRows

	}
	// 在这里处理结果集

	// 怎么知道SELECT 出来了哪些列
	// 拿到SELECT出来的列
	//cs, err := rows.Columns()
	//if err != nil {
	//	return nil, err
	//}

	// 怎么利用cs来解决顺序问题和类型问题

	tp := new(T)
	val := s.db.creator(s.model, tp)
	err = val.SetColums(rows)
	// 接口定义好之后，就两件事，一个是用新接口的方法改造上层
	// 一个就是提供不同的实现
	return tp, err

}

func (s *Selector[T]) GetV1(ctx context.Context) (*T, error) {
	q, err := s.Build()
	if err != nil {
		return nil, err
	}

	db := s.db.db
	// 发起查询，并且处理结果集
	rows, err := db.QueryContext(ctx, q.SQL, q.Args...)

	// 这个是查询错误
	if err != nil {
		return nil, err
	}
	// 继续处理结果集

	// 确认是否有数据
	if !rows.Next() {
		// 是否返回error
		// 返回error,和sql包语义保持一致
		return nil, errs.ErrNoRows

	}
	// 在这里处理结果集

	// 怎么知道SELECT 出来了哪些列
	// 拿到SELECT出来的列
	cs, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// 怎么利用cs来解决顺序问题和类型问题

	tp := new(T)
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

	address := reflect.ValueOf(tp).UnsafePointer()
	for _, c := range cs {
		fd, ok := s.model.ColumnMap[c]
		if !ok {
			return nil, errs.NewErrUnkownColumn(c)
		}
		// reflect.ValueOf(&myStruct)
		fdAddress := unsafe.Pointer(uintptr(address) + fd.Offset)
		// 反射创建一个实例,反射在特点的地址上，创建一个特定类型的实例
		// 这里创建的实例是原本类型的指针类型
		// 例如: fd.Type = int, 那么val 是 *int
		val := reflect.NewAt(fd.Type, fdAddress)
		vals = append(vals, val.Interface())
	}
	err = rows.Scan(vals...)
	return tp, err
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	q, err := s.Build()
	if err != nil {
		return nil, err
	}

	db := s.db.db
	// 发起查询，并且处理结果集
	rows, err := db.QueryContext(ctx, q.SQL, q.Args...)

	// 这个是查询错误
	if err != nil {
		return nil, err
	}
	// 继续处理结果集

	// 确认是否有数据
	if !rows.Next() {
		// 是否返回error
		// 返回error,和sql包语义保持一致
		return nil, errs.ErrNoRows

	}
	// 在这里处理结果集

	// 怎么知道SELECT 出来了哪些列
	// 拿到SELECT出来的列
	cs, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// 怎么利用cs来解决顺序问题和类型问题

	tp := new(T)

	// 通过cs 来构造vals
	vals := make([]any, 0, len(cs))
	valsElem := make([]reflect.Value, 0, len(cs))
	for _, c := range cs {
		fd, ok := s.model.ColumnMap[c]
		if !ok {
			return nil, errs.NewErrUnkownColumn(c)
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
		return nil, err
	}

	// 想办法将vals 塞进去 结果tp里面
	tpValue := reflect.ValueOf(tp).Elem()

	for i, c := range cs {
		fd, ok := s.model.ColumnMap[c]
		if !ok {
			return nil, errs.NewErrUnkownColumn(c)
		}

		tpValue.FieldByName(fd.GoName).Set(valsElem[i])
		//for _, fd := range s.model.fieldMap {
		//	if fd.colName == c {
		//		tpValue.Elem().FieldByName(fd.goName).Set(reflect.ValueOf(vals[i]).Elem())
		//	}
		//}
	}

	return tp, nil
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	q, err := s.Build()
	if err != nil {
		return nil, err
	}

	db := s.db.db
	// 发起查询，并且处理结果集
	rows, err := db.QueryContext(ctx, q.SQL, q.Args)
	// 继续处理结果集
	for rows.Next() {
		// 构造[]*T
	}
	return nil, err
}

func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}

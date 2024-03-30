package orm

import (
	"context"
	_"database/sql"
	"github.com/ExerciseCoding/template/internal/orm/internal/errs"
	"reflect"
	"strings"
)

type Selector[T any] struct {
	table string
	model *Model
	where []Predicate
	sb    *strings.Builder
	args  []any

	db *DB
	//r *registry
}

func NewSelector[T any] (db *DB) *Selector[T] {
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
	sb.WriteString("SELECT * FROM ")
	if s.table == "" {
		sb.WriteByte('`')
		sb.WriteString(s.model.tableName)
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
		fd, ok := s.model.fieldMap[exp.name]
		// 字段不对，或者列不对
		if !ok {
			return errs.NewErrUnkownField(exp.name)
		}
		s.sb.WriteByte('`')
		s.sb.WriteString(fd.colName)
		s.sb.WriteByte('`')

	case Value:
		s.sb.WriteByte('?')
		s.addArgs(exp.val)

	default:
		return errs.NewErrUnsupportedExpression(expr)
	}
	return nil
}

func (s *Selector[T]) addArgs(val any) {
	if s.args == nil {
		s.args = make([]any, 0, 8)
	}
	s.args = append(s.args, val)
}

func (s *Selector[T]) From(table string) *Selector[T] {
	s.table = table
	return s
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
		fd, ok := s.model.columnMap[c]
		if !ok {
			return nil, errs.NewErrUnkownColumn(c)
		}
		val := reflect.New(fd.typ)
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
		fd, ok := s.model.columnMap[c]
		if !ok{
			return nil, errs.NewErrUnkownColumn(c)
		}

		tpValue.FieldByName(fd.goName).Set(valsElem[i])
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

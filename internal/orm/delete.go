package orm

import (
	"fmt"
	"strings"
)

type Delete[T any] struct{
	builder
	where []Predicate
	table string

	r *registry
}

func (d *Delete[T]) Build() (*Query, error) {
	d.sb = &strings.Builder{}
	var err error
	fmt.Println("11")
	d.model, err = d.r.Register(new(T))
	fmt.Println("222")
	if err != nil {
		return nil, err
	}
	sb := d.sb
	sb.WriteString("DELETE FROM ")
	if d.table == "" {
		sb.WriteByte('`')
		sb.WriteString(d.model.tableName)
		sb.WriteByte('`')
	} else {
		sb.WriteString(d.table)
	}

	if len(d.where) > 0 {
		sb.WriteString(" WHERE ")
		err = d.buildPredicates(d.where)
		if err != nil {
			return nil, err
		}
		//p := d.where[0]
		//for i := 1; i < len(d.where); i++ {
		//	p = p.And(d.where[i])
		//}
		//if err := d.buildExpression(p); err != nil {
		//	return nil, err
		//}
	}
	sb.WriteByte(';')
	return &Query{
		SQL: sb.String(),
		Args: d.args,
	}, nil
}




//func (d *Delete[T]) buildExpression(expr Expression) error {
//	switch exp := expr.(type) {
//	case nil:
//	case Predicate:
//		// 在这里处理 p
//		// p.left 构建好
//		// p.op 构建好
//		// p.right 构建好
//		_, ok := exp.left.(Predicate)
//		if ok {
//			d.sb.WriteByte('(')
//		}
//
//		if err := d.buildExpression(exp.left); err != nil {
//			return err
//		}
//
//		if ok {
//			d.sb.WriteByte(')')
//		}
//		d.sb.WriteByte(' ')
//		d.sb.WriteString(exp.op.String())
//		d.sb.WriteByte(' ')
//
//		_, ok = exp.right.(Predicate)
//		if ok {
//			d.sb.WriteByte('(')
//		}
//
//		if err := d.buildExpression(exp.right); err != nil {
//			return err
//		}
//
//		//switch left := expr.left.(type) {
//		//case Column:
//		//	sb.WriteByte('`')
//		//	sb.WriteString(left.name)
//		//	sb.WriteByte('`')
//		//}
//		//
//		//switch right := expr.right.(type) {
//		//case Value:
//		//	sb.WriteByte('?')
//		//	args = append(args, right.val)
//		//}
//		if ok {
//			d.sb.WriteByte(')')
//		}
//
//	case Column:
//		fd, ok := d.model.fields[exp.name]
//		// 字段不对，或者列不对
//		if !ok {
//			return errs.NewErrUnkownField(exp.name)
//		}
//		d.sb.WriteByte('`')
//		d.sb.WriteString(fd.colName)
//		d.sb.WriteByte('`')
//
//	case Value:
//		d.sb.WriteByte('?')
//		d.addArgs(exp.val)
//
//	default:
//		return errs.NewErrUnsupportedExpression(expr)
//	}
//	return nil
//}


func (d *Delete[T]) From(table string) (*Delete[T]) {
	d.table = table
	return d
}

func (d *Delete[T]) Where(ps ...Predicate) *Delete[T] {
	d.where = ps
	return d
}


//func (d *Delete[T]) addArgs(val any) {
//	if d.args == nil {
//		d.args = make([]any, 0, 8)
//	}
//	d.args = append(d.args, val)
//}
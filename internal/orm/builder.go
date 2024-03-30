package orm

import (
	"github.com/ExerciseCoding/template/internal/orm/internal/errs"
	"strings"
)

type builder struct {
	sb *strings.Builder
	model *Model
	args []any
}

// 方法一：buildPredicates抽离方式
//type Predicates []Predicate
//func (p Predicates) buildPredicates(s *strings.Builder) error {
//
//}


// 方法二：buildPredicates抽离方式
//type Predicates struct {
//	// WHERE 或者 HAVING
//	prefix string
//	ps []Predicate
//}

//func (p Predicates) buildPredicates(s *strings.Builder) error {
//	// 拼接WHERE或者HAVING的部分
//}


func (d *builder) buildPredicates(ps []Predicate) error {
	p := ps[0]
	for i := 1; i < len(ps); i++ {
		p = p.And(ps[i])
	}
	return d.buildExpression(p)
}


func (d *builder) buildExpression(expr Expression) error {
	switch exp := expr.(type) {
	case nil:
	case Predicate:
		// 在这里处理 p
		// p.left 构建好
		// p.op 构建好
		// p.right 构建好
		_, ok := exp.left.(Predicate)
		if ok {
			d.sb.WriteByte('(')
		}

		if err := d.buildExpression(exp.left); err != nil {
			return err
		}

		if ok {
			d.sb.WriteByte(')')
		}
		d.sb.WriteByte(' ')
		d.sb.WriteString(exp.op.String())
		d.sb.WriteByte(' ')

		_, ok = exp.right.(Predicate)
		if ok {
			d.sb.WriteByte('(')
		}

		if err := d.buildExpression(exp.right); err != nil {
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
			d.sb.WriteByte(')')
		}

	case Column:
		fd, ok := d.model.fieldMap[exp.name]
		// 字段不对，或者列不对
		if !ok {
			return errs.NewErrUnkownField(exp.name)
		}
		d.sb.WriteByte('`')
		d.sb.WriteString(fd.colName)
		d.sb.WriteByte('`')

	case Value:
		d.sb.WriteByte('?')
		d.addArgs(exp.val)

	default:
		return errs.NewErrUnsupportedExpression(expr)
	}
	return nil
}

func (d *builder) addArgs(val any) {
	if d.args == nil {
		d.args = make([]any, 0, 8)
	}
	d.args = append(d.args, val)
}
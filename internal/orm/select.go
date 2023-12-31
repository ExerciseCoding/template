package orm

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

type Selector[T any] struct {
	table string
	where []Predicate
	sb    *strings.Builder
	args  []any
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
	sb := s.sb
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
		s.sb.WriteByte('`')
		s.sb.WriteString(exp.name)
		s.sb.WriteByte('`')

	case Value:
		s.sb.WriteByte('?')
		s.addArgs(exp.val)

	default:
		return fmt.Errorf("orm: 不支持的表达式类型 %v", expr)
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
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}

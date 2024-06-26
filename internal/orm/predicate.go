package orm

// 衍生类型
type op string

const (
	opEq  op = "="
	opLt op = "<"
	opNot op = "NOT"
	opAnd op = "AND"
	opOr  op = "OR"
)

func (o op) String() string {
	return string(o)
}

// 别名
// type op=string
//type Predicate struct {
//	c   Column
//	op  op
//	arg any
//}

type Predicate struct {
	left  Expression
	op    op
	right Expression
}

// Eq("id",12)
// Eq(sub, "id", 12)
// Eq(sub.id, "id", 12)
//func Eq(column string, arg any) Predicate {
//	return Predicate{
//		Column: column,
//		Op:     "=",
//		Arg:    arg,
//	}
//}



// C("id").Eq(12)
// sub.C("id").Eq(12)
//func (c Column) Eq(arg any) Predicate {
//	return Predicate{
//		c:   c,
//		op:  opEq,
//		arg: arg,
//	}
//}


// Not(C("name").Eq("Tom"))
func Not(p Predicate) Predicate {
	return Predicate{
		op:    opNot,
		right: p,
	}
}

// And C("id").Eq(12).And(C("name").Eq("Tom"))
func (left Predicate) And(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opAnd,
		right: right,
	}
}

// Or C("id").Eq(12).Or(C("name").Eq("Tom"))
func (left Predicate) Or(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opOr,
		right: right,
	}
}

func (Predicate) expr() {}

// Expression 是一个标记接口，代表表达式
type Expression interface {
	expr()
}

type Value struct {
	val any
}

func (Value) expr() {}

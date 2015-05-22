package data

type Clause interface {
	Clause() string
	Arg() interface{}
}

type PrefixClause struct {
	Col, Val string
}

func (q PrefixClause) Clause() string {
	return q.Col + " LIKE ?"
}

func (q PrefixClause) Arg() interface{} {
	return q.Val + "%"
}

type LessThanClause struct {
	Col string
	Val interface{}
}

func (q LessThanClause) Clause() string {
	return q.Col + " < ?"
}

func (q LessThanClause) Arg() interface{} {
	return q.Val
}

type GreaterThanClause struct {
	Col string
	Val interface{}
}

func (q GreaterThanClause) Clause() string {
	return q.Col + " > ?"
}

func (q GreaterThanClause) Arg() interface{} {
	return q.Val
}

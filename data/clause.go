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

package main

type Expr interface{}

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

func NewBinary(left Expr, operator Token, right Expr) Binary {
	return Binary{
		left,
		operator,
		right,
	}
}

type Grouping struct {
	expression Expr
}

func NewGrouping(expression Expr) Grouping {
	return Grouping{
		expression,
	}
}

type Literal struct {
	value interface{}
}

func NewLiteral(value interface{}) Literal {
	return Literal{
		value,
	}
}

type Unary struct {
	operator Token
	right    Expr
}

func NewUnary(operator Token, right Expr) Unary {
	return Unary{
		operator,
		right,
	}
}

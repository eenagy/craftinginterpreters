package main

type Expr interface {
	Accept(v Visitor) (interface{}, error)
}
type Visitor interface {
	VisitBinaryExpr(expr Binary) (interface{}, error)
	VisitGroupingExpr(expr Grouping) (interface{}, error)
	VisitLiteralExpr(expr Literal) (interface{}, error)
	VisitUnaryExpr(expr Unary) (interface{}, error)
}

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
func (a Binary) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitBinaryExpr(a)
}

type Grouping struct {
	expression Expr
}

func NewGrouping(expression Expr) Grouping {
	return Grouping{
		expression,
	}
}
func (a Grouping) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitGroupingExpr(a)
}

type Literal struct {
	value interface{}
}

func NewLiteral(value interface{}) Literal {
	return Literal{
		value,
	}
}
func (a Literal) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitLiteralExpr(a)
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
func (a Unary) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitUnaryExpr(a)
}

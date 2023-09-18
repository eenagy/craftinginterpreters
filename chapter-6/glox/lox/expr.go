package main

type Expr interface {
	Accept(v Visitor) interface{}
}
type Visitor interface {
	VisitBinaryExpr(expr Binary) interface{}
	VisitGroupingExpr(expr Grouping) interface{}
	VisitLiteralExpr(expr Literal) interface{}
	VisitUnaryExpr(expr Unary) interface{}
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
func (a Binary) Accept(visitor Visitor) interface{} {
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
func (a Grouping) Accept(visitor Visitor) interface{} {
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
func (a Literal) Accept(visitor Visitor) interface{} {
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
func (a Unary) Accept(visitor Visitor) interface{} {
	return visitor.VisitUnaryExpr(a)
}

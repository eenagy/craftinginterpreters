package main

type Expr interface {
	Accept(v ExprVisitor) (interface{}, error)
}
type ExprVisitor interface {
	VisitAssignExpr(expr Assign) (interface{}, error)
	VisitBinaryExpr(expr Binary) (interface{}, error)
	VisitGroupingExpr(expr Grouping) (interface{}, error)
	VisitLiteralExpr(expr Literal) (interface{}, error)
	VisitLogicalExpr(expr Logical) (interface{}, error)
	VisitUnaryExpr(expr Unary) (interface{}, error)
	VisitVariableExpr(expr Variable) (interface{}, error)
}

type Assign struct {
	name  Token
	value Expr
}

func NewAssign(name Token, value Expr) Assign {
	return Assign{
		name,
		value,
	}
}
func (a Assign) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitAssignExpr(a)
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
func (a Binary) Accept(visitor ExprVisitor) (interface{}, error) {
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
func (a Grouping) Accept(visitor ExprVisitor) (interface{}, error) {
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
func (a Literal) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitLiteralExpr(a)
}

type Logical struct {
	left     Expr
	operator Token
	right    Expr
}

func NewLogical(left Expr, operator Token, right Expr) Logical {
	return Logical{
		left,
		operator,
		right,
	}
}
func (a Logical) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitLogicalExpr(a)
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
func (a Unary) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitUnaryExpr(a)
}

type Variable struct {
	name Token
}

func NewVariable(name Token) Variable {
	return Variable{
		name,
	}
}
func (a Variable) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitVariableExpr(a)
}

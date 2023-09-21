package main

type Stmt interface {
	Accept(v StmtVisitor) (interface{}, error)
}
type StmtVisitor interface {
	VisitBlockStmt(stmt Block) (interface{}, error)
	VisitExpressionStmt(stmt Expression) (interface{}, error)
	VisitPrintStmt(stmt Print) (interface{}, error)
	VisitVarStmt(stmt Var) (interface{}, error)
}

type Block struct {
	statements []Stmt
}

func NewBlock(statements []Stmt) Block {
	return Block{
		statements,
	}
}
func (a Block) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitBlockStmt(a)
}

type Expression struct {
	expression Expr
}

func NewExpression(expression Expr) Expression {
	return Expression{
		expression,
	}
}
func (a Expression) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitExpressionStmt(a)
}

type Print struct {
	expression Expr
}

func NewPrint(expression Expr) Print {
	return Print{
		expression,
	}
}
func (a Print) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitPrintStmt(a)
}

type Var struct {
	name        Token
	initializer Expr
}

func NewVar(name Token, initializer Expr) Var {
	return Var{
		name,
		initializer,
	}
}
func (a Var) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitVarStmt(a)
}

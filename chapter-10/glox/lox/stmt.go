package main

type Stmt interface {
	Accept(v StmtVisitor) (interface{}, error)
}
type StmtVisitor interface {
	VisitBlockStmt(stmt Block) (interface{}, error)
	VisitExpressionStmt(stmt Expression) (interface{}, error)
	VisitIfStmt(stmt If) (interface{}, error)
	VisitPrintStmt(stmt Print) (interface{}, error)
	VisitVarStmt(stmt Var) (interface{}, error)
	VisitWhileStmt(stmt While) (interface{}, error)
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

type If struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func NewIf(condition Expr, thenBranch Stmt, elseBranch Stmt) If {
	return If{
		condition,
		thenBranch,
		elseBranch,
	}
}
func (a If) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitIfStmt(a)
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

type While struct {
	condition Expr
	body      Stmt
}

func NewWhile(condition Expr, body Stmt) While {
	return While{
		condition,
		body,
	}
}
func (a While) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitWhileStmt(a)
}

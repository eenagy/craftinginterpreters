package main

import (
	"fmt"
	"strconv"
)

type Interpreter struct {
	environment *Environment
}

type RuntimeError struct {
	Operator Token
	Message  string
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("%s: %s", e.Operator, e.Message)
}

func NewInterpreter() Interpreter {
	return Interpreter{
		environment: NewEnvironment(nil),
	}
}

func (i *Interpreter) Interpret(stmts []Stmt) {
	for _, stmt := range stmts {
		_, err := i.execute(stmt)
		if err != nil {
			ReportRuntimeError(err.(RuntimeError))
		}
	}

}
func (i *Interpreter) VisitBinaryExpr(expr Binary) (interface{}, error) {
	var err error
	var right interface{}
	var left interface{}
	right, err = i.evaluate(expr.right)
	if err != nil {
		return nil, err
	}
	left, err = i.evaluate(expr.left)
	if err != nil {
		return nil, err
	}
	switch expr.operator.Type {
	case GREATER:
		err = checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) >
			right.(float64), nil
	case GREATER_EQUAL:
		err = checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case LESS:
		err = checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case LESS_EQUAL:
		err = checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case MINUS:
		err = checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		value := left.(float64) - right.(float64)
		return strconv.FormatFloat(value, 'f', -1, 64), nil
	case BANG_EQUAL:
		return !isEqual(left, right), nil
	case EQUAL_EQUAL:
		return isEqual(left, right), nil
	case PLUS:
		r := isNumber(right)
		l := isNumber(left)
		if l && r {
			value := left.(float64) + right.(float64)
			return value, nil
		} else {
			var leftvalue string
			var rightvalue string
			if l {
				leftvalue = strconv.FormatFloat(left.(float64), 'f', -1, 64)
			} else {
				leftvalue = left.(string)
			}
			if r {
				rightvalue = strconv.FormatFloat(right.(float64), 'f', -1, 64)
			} else {
				rightvalue = right.(string)
			}
			value := leftvalue + rightvalue
			return value, nil
		}
	case SLASH:
		err = checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		value := left.(float64) / right.(float64)
		return strconv.FormatFloat(value, 'f', -1, 64), nil
	case STAR:
		err = checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		value := left.(float64) * right.(float64)
		return strconv.FormatFloat(value, 'f', -1, 64), nil
	}
	// unreachable
	return nil, nil
}
func (i *Interpreter) VisitGroupingExpr(expr Grouping) (interface{}, error) {
	return i.evaluate(expr.expression)
}
func (i *Interpreter) VisitLiteralExpr(expr Literal) (interface{}, error) {
	return expr.value, nil

}
func (i *Interpreter) VisitUnaryExpr(expr Unary) (interface{}, error) {
	var err error
	var right interface{}
	right, err = i.evaluate(expr.right)
	if err != nil {
		return nil, err
	}
	switch expr.operator.Type {
	case BANG:
		return !i.isTruthy(right), nil
	case MINUS:
		err = checkNumberOperand(expr.operator, right)
		if err != nil {
			return nil, err
		}
		return strconv.FormatFloat(right.(float64), 'f', -1, 64), nil
	}
	// unreachable
	return nil, nil
}
func (i *Interpreter) VisitVariableExpr(stmt Variable) (interface{}, error) {
	return i.environment.Get(stmt.name)
}

func (i *Interpreter) VisitAssignExpr(stmt Assign) (interface{}, error) {
	value, err := i.evaluate(stmt.value)
	if err != nil {
		return nil, err
	}
	i.environment.Assign(stmt.name, value)
	return value, nil
}
func (i *Interpreter) VisitVarStmt(stmt Var) (interface{}, error) {
	var value interface{}
	var err error
	if stmt.initializer != nil {
		value, err = i.evaluate(stmt.initializer)
		if err != nil {
			return nil, err
		}
	}
	i.environment.Define(stmt.name.Lexeme, value)
	return nil, nil
}
func (i *Interpreter) VisitExpressionStmt(stmt Expression) (interface{}, error) {
	_, err := i.evaluate(stmt.expression)
	return nil, err
}
func (i *Interpreter) VisitBlockStmt(stmt Block) (interface{}, error) {
	_, err := i.executeBlock(stmt.statements, *NewEnvironment(i.environment))
	return nil, err
}
func (i *Interpreter) VisitPrintStmt(stmt Print) (interface{}, error) {
	value, err := i.evaluate(stmt.expression)
	if err == nil {
		fmt.Println(stringify(value))
		return nil, nil
	} else {
		return nil, err
	}
}

func (i *Interpreter) execute(stmt Stmt) (interface{}, error) {
	return stmt.Accept(i)
}

func (i *Interpreter) executeBlock(statements []Stmt, environment Environment) (interface{}, error) {

	previous := i.environment.Copy()
	i.environment = &environment
	for _, statement := range statements {
		_, err := i.execute(statement)
		if err != nil {
			i.environment = previous
			return nil, nil
		}
	}
	i.environment = previous
	return nil, nil
}

func (i *Interpreter) evaluate(expr Expr) (interface{}, error) {
	return expr.Accept(i)
}

func (i *Interpreter) isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}
	if object == false {
		return false
	}
	return true
}

func isNumber(object interface{}) bool {
	switch object.(type) {
	case int:
		return true
	case float32, float64:
		return true
	default:
		return false
	}
}

func isEqual(left interface{}, right interface{}) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil {
		return false
	}
	return left == right
}

func checkNumberOperands(operator Token, left interface{}, right interface{}) error {
	if isNumber(left) && isNumber(right) {
		return nil
	}
	return &RuntimeError{Operator: operator, Message: "Operands must be a number."}

}
func checkNumberOperand(operator Token, operand interface{}) error {
	if isNumber(operand) {
		return nil
	}
	return &RuntimeError{Operator: operator, Message: "Operands must be a number."}

}

func stringify(object interface{}) string {
	if object == nil {
		return "nil"
	}

	switch value := object.(type) {
	case int:
		return strconv.Itoa(value)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case string:
		return fmt.Sprintf("\"%v\"", object)
	default:
		return fmt.Sprintf("Unknown type")
	}
}

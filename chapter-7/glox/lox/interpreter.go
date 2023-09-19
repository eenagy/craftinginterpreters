package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Interpreter struct {
}

type RuntimeError struct {
	Operator Token
	Message  string
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("%s: %s", e.Operator, e.Message)
}

func NewInterpreter() Interpreter {
	return Interpreter{}
}

func (i Interpreter) Interpret(expr Expr) {
	value, err := i.evaluate(expr)
	if err != nil {
		ReportRuntimeError(err.(RuntimeError))
	} else {
		fmt.Println(stringify(value))
	}
}
func (i Interpreter) VisitBinaryExpr(expr Binary) (interface{}, error) {
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
			return strconv.FormatFloat(value, 'f', -1, 64), nil
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
func (i Interpreter) VisitGroupingExpr(expr Grouping) (interface{}, error) {
	return i.evaluate(expr.expression)
}
func (i Interpreter) VisitLiteralExpr(expr Literal) (interface{}, error) {
	return expr.value, nil

}
func (i Interpreter) VisitUnaryExpr(expr Unary) (interface{}, error) {
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

func (i Interpreter) evaluate(expr Expr) (interface{}, error) {
	return expr.Accept(i)
}

func (i Interpreter) isTruthy(object interface{}) bool {
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

	// Check if the object is of type float64 (equivalent to Java Double)
	if floatValue, ok := object.(float64); ok {
		text := fmt.Sprintf("%v", floatValue)
		// Remove trailing ".0" if it exists
		if strings.HasSuffix(text, ".0") {
			return text[:len(text)-2]
		}
		return text
	}

	// If it's not a float64, use the default Go string representation
	return fmt.Sprintf("\"%v\"", object)
}

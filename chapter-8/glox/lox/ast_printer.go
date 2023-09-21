package main

import (
	"strconv"
	"strings"
)

type AstPrinter struct {
}

func NewAstPrinter() AstPrinter {
	return AstPrinter{}
}
func (a AstPrinter) Print(expr Expr) (string, error) {
	r, err := expr.Accept(a)
	if err != nil {
		return "", err
	}
	return r.(string), nil
}
func (a AstPrinter) VisitBinaryExpr(expr Binary) (interface{}, error) {
	return parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}
func (a AstPrinter) VisitGroupingExpr(expr Grouping) (interface{}, error) {
	return parenthesize("group", expr.expression)
}

func (a AstPrinter) VisitAssignExpr(expr Assign) (interface{}, error) {
	return nil, nil
}
func (a AstPrinter) VisitVariableExpr(expr Variable) (interface{}, error) {
	return nil, nil
}
func (a AstPrinter) VisitLiteralExpr(expr Literal) (interface{}, error) {
	if expr.value == nil {
		return "nil", nil
	}
	switch v := expr.value.(type) {
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	default:
		return expr.value.(string), nil
	}
}
func (a AstPrinter) VisitUnaryExpr(expr Unary) (interface{}, error) {
	return parenthesize(expr.operator.Lexeme, expr.right)
}

func parenthesize(name string, exprs ...Expr) (string, error) {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)

	for _, expr := range exprs {
		builder.WriteString(" ")
		printer := AstPrinter{}
		r, err := expr.Accept(printer)
		if err != nil {
			return "", err
		}
		builder.WriteString((r).(string))
	}

	builder.WriteString(")")

	return builder.String(), nil
}

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
func (a AstPrinter) Print(expr Expr) string {
	return expr.Accept(a).(string)
}
func (a AstPrinter) VisitBinaryExpr(expr Binary) interface{} {
	return parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}
func (a AstPrinter) VisitGroupingExpr(expr Grouping) interface{} {
	return parenthesize("group", expr.expression)
}
func (a AstPrinter) VisitLiteralExpr(expr Literal) interface{} {
	if expr.value == nil {
		return "nil"
	}
	switch v := expr.value.(type) {
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return expr.value.(string)
	}
}
func (a AstPrinter) VisitUnaryExpr(expr Unary) interface{} {
	return parenthesize(expr.operator.Lexeme, expr.right)
}

func parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)

	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(expr.Accept(AstPrinter{}).(string))
	}

	builder.WriteString(")")

	return builder.String()
}

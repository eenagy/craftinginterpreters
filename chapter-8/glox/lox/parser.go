package main

import (
	"errors"
)

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) Parser {
	return Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() []Stmt {
	var statements []Stmt
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err == nil {
			statements = append(statements, stmt)
		}
	}
	return statements
}
func (p *Parser) declaration() (Stmt, error) {
	var err error
	var result Stmt
	if p.match(VAR) {
		result, err = p.varDeclaration()
		if err != nil {
			p.synchronize()
			return nil, nil
		}
		return result, nil
	}
	result, err = p.statement()
	if err != nil {
		p.synchronize()
		return nil, nil
	}
	return result, nil
}

func (p *Parser) varDeclaration() (Stmt, error) {
	err := p.consume(IDENTIFIER, "Expect variable name")
	if err != nil {
		return nil, err
	}
	name := p.previous()
	var initializer Expr
	if p.match(EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	err = p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}
	return NewVar(name, initializer), nil
}
func (p *Parser) statement() (Stmt, error) {
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(LEFT_BRACE) {
		return p.blockStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) blockStatement() (Stmt, error) {
	var statements []Stmt
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	err := p.consume(RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}
	return NewBlock(statements), nil
}
func (p *Parser) printStatement() (Stmt, error) {
	var value Expr
	var err error
	value, err = p.expression()
	if err != nil {
		return nil, err
	}
	err = p.consume(SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}
	return NewPrint(value), nil
}
func (p *Parser) expressionStatement() (Stmt, error) {
	var value Expr
	var err error
	value, err = p.expression()
	if err != nil {
		return nil, err
	}
	err = p.consume(SEMICOLON, "Expect ';' after expression.")
	if err != nil {
		return nil, err
	}
	return NewExpression(value), nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}
	if p.match(EQUAL) {
		equals := p.previous()
		var value Expr
		value, err = p.assignment()
		if err != nil {
			return nil, err
		}
		variable, isVariable := expr.(Variable)
		if isVariable {
			name := variable.name
			return NewAssign(name, value), nil
		}
		TokenError(equals, "Invalid Assignment target.")
	}
	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		var right Expr
		right, err = p.comparison()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, right)
	}
	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		var left Expr
		left, err = p.term()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, left)
	}
	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		var left Expr
		left, err = p.factor()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, left)
	}
	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.match(SLASH, STAR) {
		operator := p.previous()
		var left Expr
		left, err = p.unary()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, left)
	}
	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	for p.match(BANG, MINUS) {
		operator := p.previous()
		left, err := p.primary()
		return NewUnary(operator, left), err
	}
	return p.primary()
}
func (p *Parser) primary() (Expr, error) {

	if p.match(FALSE) {
		return NewLiteral(false), nil
	}
	if p.match(TRUE) {
		return NewLiteral(true), nil
	}
	if p.match(NIL) {
		return NewLiteral(nil), nil
	}
	if p.match(NUMBER, STRING) {
		return NewLiteral(p.previous().Literal), nil
	}
	if p.match(IDENTIFIER) {
		return NewVariable(p.previous()), nil
	}
	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return NewGrouping(expr), err
		}
		err = p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return NewGrouping(expr), err
	}
	return nil, p.error(p.peek(), "Expect expression")
}

func (p *Parser) match(types ...TokenType) bool {
	for _, token_type := range types {
		if p.check(token_type) {
			p.advance()
			return true
		}
	}
	return false
}
func (p *Parser) check(token_type TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == token_type
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()

}
func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}
func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}
func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) consume(token_type TokenType, message string) error {
	if p.check(token_type) {
		p.advance()
		return nil
	}
	return p.error(p.peek(), message)
}

func (p *Parser) error(token Token, message string) error {
	TokenError(token, message)
	return errors.New("ParseError")
}

func (p *Parser) synchronize() {
}

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

func (p *Parser) Parse() Expr {
	result, err := p.expression()
	if err != nil {
		return nil
	}
	return result
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		var right Expr
		right, err = p.comparison()
		expr = NewBinary(expr, operator, right)
	}
	return expr, err
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		var left Expr
		left, err = p.term()
		expr = NewBinary(expr, operator, left)
	}
	return expr, err
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	for p.match(MINUS, PLUS) {
		operator := p.previous()
		var left Expr
		left, err = p.factor()
		expr = NewBinary(expr, operator, left)
	}
	return expr, err
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	for p.match(SLASH, STAR) {
		operator := p.previous()
		var left Expr
		left, err = p.unary()
		expr = NewBinary(expr, operator, left)
	}
	return expr, err
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
	if p.match(STRING, NUMBER) {
		return NewLiteral(p.previous().Literal), nil
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
	}
	return p.error(p.peek(), message)
}

func (p *Parser) error(token Token, message string) error {
	TokenError(token, message)
	return errors.New("ParseError")
}

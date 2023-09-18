package main

import (
	"strconv"
	"unicode"
)

type Scanner struct {
	source   string
	tokens   []Token
	start    int
	current  int
	line     int
	keywords map[string]TokenType
}

func NewScanner(source string) Scanner {
	keywords := map[string]TokenType{
		"and":    AND,
		"class":  CLASS,
		"else":   ELSE,
		"false":  FALSE,
		"for":    FOR,
		"fun":    FUN,
		"if":     IF,
		"nil":    NIL,
		"or":     OR,
		"print":  PRINT,
		"return": RETURN,
		"super":  SUPER,
		"this":   THIS,
		"true":   TRUE,
		"var":    VAR,
		"while":  WHILE,
	}
	return Scanner{
		source:   source,
		start:    0,
		current:  0,
		line:     1,
		keywords: keywords,
	}
}

// Define a method on the Person type
func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		// We are at the beginning of the next lexeme.
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, NewToken(EOF, "", nil, s.line))
	return s.tokens
}

func (s Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
		break
	case ')':
		s.addToken(RIGHT_PAREN)
		break
	case '{':
		s.addToken(LEFT_BRACE)
		break
	case '}':
		s.addToken(RIGHT_BRACE)
		break
	case ',':
		s.addToken(COMMA)
		break
	case '.':
		s.addToken(DOT)
		break
	case '-':
		s.addToken(MINUS)
		break
	case '+':
		s.addToken(PLUS)
		break
	case ';':
		s.addToken(SEMICOLON)
		break
	case '*':
		s.addToken(STAR)
		break
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
		break
	case ' ':
	case '\r':
	case '\t':
		// Ignore whitespace.
		break

	case '\n':
		s.line++
		break
	case '"':
		s.string()
		break
	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			Error(s.line, "Unexpected character.")
		}
		break
	}
}

func (s *Scanner) advance() rune {
	s.current++
	return rune(s.source[s.current-1])
}

func (s *Scanner) addToken(tokenType TokenType, literals ...interface{}) {
	text := s.source[s.start:s.current]
	var literal interface{}
	if len(literals) > 0 {
		literal = literals[0]
	}
	s.tokens = append(s.tokens, NewToken(tokenType, text, literal, s.line))

}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if rune(s.source[s.current]) != expected {
		return false
	}

	s.current++
	return true
}

func (s Scanner) peek() rune {
	if s.isAtEnd() {
		return '\x00'
	}
	return rune(s.source[s.current])
}

func (s Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '\x00'
	}
	return rune(s.source[s.current+1])
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		Error(s.line, "Unterminated string.")
		return
	}
	// The closing "
	s.advance()

	// Trim the surrounding quotes
	start := s.start
	end := s.current
	value := s.source[start:end]
	s.addToken(STRING, value)
}

func (s Scanner) isDigit(c rune) bool {
	return unicode.IsDigit(c)
}
func (s Scanner) isAlpha(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c)
}
func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}
	number, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err == nil {
		s.addToken(NUMBER, number)
	}
}
func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	token_type, exists := s.keywords[text]
	if !exists {
		token_type = IDENTIFIER
	}
	s.addToken(token_type)
}

func (s Scanner) isAlphaNumeric(c rune) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

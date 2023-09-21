package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var hadError = false
var hadRuntimeError = false
var interpreter = NewInterpreter()

func LoxMain(args []string) {
	length := len(args)

	if length > 1 {
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	} else if length == 1 {
		runFile(args[0])
	} else {
		runPrompt()
	}
}

func runFile(filePath string) {
	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	run(string(content))
	if hadError {
		os.Exit(65)
	}
	if hadRuntimeError {
		os.Exit(70)
	}
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			scanner.Err()
			break
		}
		line := scanner.Text()
		run(line)
		hadError = false
	}
}

func run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()
	// fmt.Println("Tokens: ")
	// for _, token := range tokens {
	// 	fmt.Printf("Type: %-12s Lexeme: %-10s Literal: %v \n", TokenName[token.Type], token.Lexeme, token.Literal)
	// }
	// fmt.Println("")
	parser := NewParser(tokens)
	statements := parser.Parse()
	if hadError {
		return
	}
	interpreter.Interpret(statements)
}

func Error(line int, message string) {
	report(line, "", message)
}

func TokenError(token Token, message string) {
	if token.Type == EOF {
		report(token.Line, " at the end ", message)
	} else {
		report(token.Line, " at '"+token.Lexeme+"'", message)
	}
}

func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s\n", line, where, message)
	hadError = true
}

func ReportRuntimeError(err RuntimeError) {
	fmt.Fprintf(os.Stderr, "%s\n[line %d]\n", err.Message, err.Operator.Line)
	hadRuntimeError = true
}
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var hadError = false

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
	for _, token := range tokens {
		fmt.Println(token)
	}
}

func Error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s\n", line, where, message)
	hadError = true
}

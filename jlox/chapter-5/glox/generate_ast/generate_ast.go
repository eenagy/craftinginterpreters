package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GeneretaAst() {
	args := os.Args
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage generate_ast <output directory>\n")
		os.Exit(64)
	}
	outdir := args[1]
	defineAst(outdir, "Expr", []string{
		"Binary   : left Expr, operator Token, right Expr",
		"Grouping : expression Expr",
		"Literal  : value interface {}",
		"Unary    : operator Token, right Expr",
	})
}

func defineAst(outdir string, baseName string, types []string) {
	path := filepath.Join(outdir, strings.ToLower(baseName)+".go")
	// Create or open the file for writing
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a bufio.Writer for efficient writing
	writer := bufio.NewWriter(file)
	writeContent(writer, `package main
	type Expr interface {}
	`)
	for _, stype := range types {
		substrings := strings.Split(stype, ":")

		defineType(writer, strings.Title(substrings[0]), substrings[1])
	}
	// Flush the buffer to ensure data is written to the file
	err = writer.Flush()
	if err != nil {
		panic(err)
	}

}
func defineType(writer *bufio.Writer, className string, fields string) {

	formatString := `
    type %s struct {
		%s
	}
	func New%s(%s) %s {
		return %s{
			%s
		}
	}
	`
	substrings := strings.Split(fields, ",")
	fieldNames := ""
	for _, field := range substrings {
		f := strings.Split(field, " ")
		fieldName := f[1]
		fmt.Println(fieldName)
		if fieldName != "" {

			fieldNames = fieldNames + fieldName + ",\n"
		}
	}

	content := fmt.Sprintf(formatString, className, strings.ReplaceAll(fields, ",", "\n"), className,
		fields, className, className, fieldNames)
	writeContent(writer, content)
}

func writeContent(writer *bufio.Writer, content string) {
	// Write content to the file
	_, err := writer.WriteString(content)
	if err != nil {
		// panic as we don't handle errors
		panic(err)
	}
}

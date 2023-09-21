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

	defineAst(outdir, "Stmt", []string{
		"Block      : statements []Stmt",
		"Expression : expression Expr",
		"Print      : expression Expr",
		"Var        : name Token, initializer Expr",
	})

	defineAst(outdir, "Expr", []string{
		"Assign   : name Token, value Expr",
		"Binary   : left Expr, operator Token, right Expr",
		"Grouping : expression Expr",
		"Literal  : value interface {}",
		"Unary    : operator Token, right Expr",
		"Variable : name Token",
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
	writeContent(writer, fmt.Sprintf(`package main
	type %s interface {
		Accept (v %sVisitor) (interface{},error)
	}
	`, baseName, baseName))

	writeLine(writer, fmt.Sprintf("type %sVisitor interface {", baseName))
	for _, stype := range types {
		substrings := strings.Split(stype, ":")
		visitorFormatString := `Visit%s%s(%s %s) (interface{}, error)`
		className := strings.ReplaceAll(strings.Title(substrings[0]), " ", "")
		//line := "visit" + className + baseName + "(" + strings.ToLower(baseName) + " " + className + " )"
		line := fmt.Sprintf(visitorFormatString, className, baseName, strings.ToLower(baseName), className)
		writeLine(writer, line)
	}
	writeLine(writer, "}")
	for _, stype := range types {
		substrings := strings.Split(stype, ":")
		className := strings.ReplaceAll(strings.Title(substrings[0]), " ", "")

		defineType(writer, baseName, className, substrings[1])
	}
	// Flush the buffer to ensure data is written to the file
	err = writer.Flush()
	if err != nil {
		panic(err)
	}

}
func defineType(writer *bufio.Writer, baseName string, className string, fields string) {

	formatString := `
    type %s struct {
		%s
	}
	func New%s(%s) %s {
		return %s{
			%s
		}
	}
	func (a %s) Accept(visitor %sVisitor) (interface{}, error) {
		return visitor.Visit%s%s(a)
	}
	`
	substrings := strings.Split(fields, ",")
	fieldNames := ""
	for _, field := range substrings {
		f := strings.Split(field, " ")
		fieldName := f[1]
		if fieldName != "" {

			fieldNames = fieldNames + fieldName + ",\n"
		}
	}

	content := fmt.Sprintf(formatString, className, strings.ReplaceAll(fields, ",", "\n"), className,
		fields, className, className, fieldNames, className, baseName, className, baseName)
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

func writeLine(writer *bufio.Writer, line string) {
	// Write content to the file
	_, err := fmt.Fprintln(writer, line)
	if err != nil {
		// panic as we don't handle errors
		panic(err)
	}
}

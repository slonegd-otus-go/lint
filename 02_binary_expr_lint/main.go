package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("must be one argument: file name for lint")
	}
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, os.Args[1], nil, 0)
	if err != nil {
		log.Fatal(err.Error())
	}

	ast.Inspect(file, func(node ast.Node) bool {
		binaryExpr, ok := node.(*ast.BinaryExpr)
		if !ok {
			return true
		}

		if binaryExpr.Op != token.ADD {
			return true
		}

		if _, ok := binaryExpr.X.(*ast.BasicLit); !ok {
			return true
		}

		if _, ok := binaryExpr.Y.(*ast.BasicLit); !ok {
			return true
		}

		position := fileset.Position(binaryExpr.Pos())
		// %q заворачивает в кавычки
		fmt.Fprintf(os.Stdout, "%s: integer addition found: %q\n", position, render(fileset, binaryExpr))

		return true
	})
}

func render(fileset *token.FileSet, node ast.Node) string {
	var buffer bytes.Buffer
	err := printer.Fprint(&buffer, fileset, node)
	if err != nil {
		log.Fatal(err.Error())
	}
	return buffer.String()
}

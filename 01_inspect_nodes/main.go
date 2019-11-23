package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

func main() {
	src, err := os.Open("struct.go")
	if err != nil {
		log.Fatal(err.Error())
	}

	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, "demo", src, parser.ParseComments)
	if err != nil {
		log.Fatal(err.Error())
	}

	ast.Inspect(file, func(x ast.Node) bool {
		structType, ok := x.(*ast.StructType)
		if !ok {
			return true
		}

		for _, field := range structType.Fields.List {
			fmt.Printf("Field: %s\n", field.Names[0].Name)
			if field.Tag != nil {
				fmt.Printf("Tag:   %s\n", field.Tag.Value)
			}
		}
		return true
	})
}

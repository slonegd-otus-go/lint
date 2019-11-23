package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func main() {
	fileset := token.NewFileSet()
	// для работы с пакетом ast
	file, err := parser.ParseFile(fileset, "struct.go", nil, 0)
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

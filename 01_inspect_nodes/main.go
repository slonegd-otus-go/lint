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

	ast.Inspect(file, func(node ast.Node) bool {
		spec, ok := node.(*ast.TypeSpec)
		if !ok {
			return true
		}

		structType, ok := spec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		log.Printf("got struct name: %s", spec.Name.Name)

		for _, field := range structType.Fields.List {
			fmt.Printf("Field: %s\n", field.Names[0].Name)
			if field.Tag != nil {
				fmt.Printf("Tag:   %s\n", field.Tag.Value)
			}
		}
		return true
	})
}

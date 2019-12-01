package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
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

		beginLine := fileset.Position(structType.Pos()).Line
		endLine := fileset.Position(structType.End()).Line
		log.Printf("got begin line: %d, end line: %d", beginLine, endLine)

		beginOffset := fileset.Position(structType.Pos()).Offset
		endOffset := fileset.Position(structType.End()).Offset
		log.Printf("got begin offset: %d, end offset: %d", beginOffset, endOffset)

		logFieldTag(structType)
		addTags(structType)
		return true
	})

	// вывод изменённого ast
	err = format.Node(os.Stdout, fileset, file)
	if err != nil {
		log.Printf("format file failed: %s", err)
	}
}

func logFieldTag(structType *ast.StructType) {
	for _, field := range structType.Fields.List {
		message := fmt.Sprintf("Field: %s", field.Names[0].Name)
		if field.Tag != nil {
			message = fmt.Sprintf("%s, Tag: %s", message, field.Tag.Value)
		}
		log.Print(message)
	}
}

func addTags(structType *ast.StructType) {
	for _, field := range structType.Fields.List {
		if field.Tag != nil {
			continue
		}

		field.Tag = &ast.BasicLit{}
		field.Tag.Value = fmt.Sprintf("`newtag:\"%s\"", field.Names[0].Name)
	}
}

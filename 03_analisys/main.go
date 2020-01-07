package main

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"log"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	// пакет analysis предоставляет всю необходимую инфу для линтера
	// не надо самим собирать их, как это было в 02_binary_expr_lint
	analyzer := &analysis.Analyzer{
		Name: "addlint",
		Doc:  "reports integer additions",
		Run:  run,
	}
	// предоставляет стандартизированный  cli для линтера
	singlechecker.Main(analyzer)
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			binaryExpr, ok := node.(*ast.BinaryExpr)
			if !ok {
				return true
			}

			if binaryExpr.Op != token.ADD {
				return true
			}

			x, ok := binaryExpr.X.(*ast.BasicLit)
			if !ok {
				return true
			}

			y, ok := binaryExpr.Y.(*ast.BasicLit)
			if !ok {
				return true
			}

			if x.Kind != token.INT || y.Kind != token.INT {
				return true
			}

			// %q заворачивает в кавычки
			pass.Reportf(binaryExpr.Pos(), "integer addition found: %q", render(pass.Fset, binaryExpr))

			return true
		})
	}
	return nil, nil
}

func render(fileset *token.FileSet, node ast.Node) string {
	var buffer bytes.Buffer
	err := printer.Fprint(&buffer, fileset, node)
	if err != nil {
		log.Fatal(err.Error())
	}
	return buffer.String()
}

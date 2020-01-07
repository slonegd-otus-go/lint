package main

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"log"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/singlechecker"
	"golang.org/x/tools/go/ast/inspector"
)

func main() {
	analyzer := &analysis.Analyzer{
		Name: "addlint",
		Doc:  "reports integer additions",
		Run:  run,
		// для запуска одного линтера это совсем не обязательно
		// выгода есть, если несколько линтеров будут пользоваться плодами этого
		// но у него есть несколько методов, которые позволяют писать чуть проще
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
	singlechecker.Main(analyzer)
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	binaryExprFilter := []ast.Node{
		(*ast.BinaryExpr)(nil),
	}

	inspect.Preorder(binaryExprFilter, func(node ast.Node) {
		binaryExpr := node.(*ast.BinaryExpr)

		if binaryExpr.Op != token.ADD {
			return
		}

		x, ok := binaryExpr.X.(*ast.BasicLit)
		if !ok {
			return
		}

		y, ok := binaryExpr.Y.(*ast.BasicLit)
		if !ok {
			return
		}

		if x.Kind != token.INT || y.Kind != token.INT {
			return
		}

		// %q заворачивает в кавычки
		pass.Reportf(binaryExpr.Pos(), "integer addition found: %q", render(pass.Fset, binaryExpr))
	})

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

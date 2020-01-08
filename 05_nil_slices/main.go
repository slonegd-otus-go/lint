package main

import (
	"bytes"
	"go/ast"
	"log"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/singlechecker"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	analyzer := &analysis.Analyzer{
		Name:     "nilSliceReturn",
		Doc:      "reports return nil slice",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
	singlechecker.Main(analyzer)
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	filter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	spew.Config.Indent = "    "
	inspect.Preorder(filter, func(node ast.Node) {
		funcDecl := node.(*ast.FuncDecl)

		if funcDecl.Type.Results == nil {
			return // функция ничего не возвращает
		}

		var returnSliceIndexes []int
		var returnErrorIndexes []int
		for i, field := range funcDecl.Type.Results.List {
			if _, ok := field.Type.(*ast.ArrayType); ok {
				returnSliceIndexes = append(returnSliceIndexes, i)
			}
			if ident, ok := field.Type.(*ast.Ident); ok {
				if ident.Name == "error" {
					returnErrorIndexes = append(returnErrorIndexes, i)
				}
			}
		}
		if len(returnSliceIndexes) == 0 {
			return // не возвращает слайс
		}

		checker := Checker{pass, returnSliceIndexes, returnErrorIndexes}

		checker.checkBody(funcDecl.Body, nil)

		buffer := &bytes.Buffer{}
		spew.Fdump(buffer, funcDecl)
		log.Printf("%s", buffer)
	})

	return nil, nil
}

func copySlice(in []string) []string {
	out := make([]string, len(in))
	copy(out, in)
	return out
}

type Checker struct {
	pass               *analysis.Pass
	returnSliceIndexes []int
	returnErrorIndexes []int
}

func (checker Checker) checkBody(body *ast.BlockStmt, notNilSlices []string) {
	notNilSlices = copySlice(notNilSlices)
	for _, tmp := range body.List {
		switch tmp.(type) {
		case *ast.ReturnStmt:
			returnStmt := tmp.(*ast.ReturnStmt)
			checker.checkReturn(returnStmt, notNilSlices)

		case *ast.IfStmt:
			ifStmt := tmp.(*ast.IfStmt)
			checker.checkBody(ifStmt.Body, notNilSlices)
			if ifStmt.Else != nil {
				checker.checkBody(ifStmt.Else.(*ast.BlockStmt), notNilSlices)
			}

		case *ast.RangeStmt:
			rangeStmt := tmp.(*ast.RangeStmt)
			checker.checkBody(rangeStmt.Body, notNilSlices)
		}

	}
}

func (checker Checker) checkReturn(returnStmt *ast.ReturnStmt, notNilSlices []string) {
	for _, index := range checker.returnSliceIndexes {
		if len(returnStmt.Results) < index {
			continue
		}
		ident, ok := returnStmt.Results[index].(*ast.Ident)
		if !ok {
			continue // TODO проверить другие варианты (функция CallExpr, переменная)
		}
		if ident.Name == "nil" {
			// %q заворачивает в кавычки
			checker.pass.Reportf(ident.Pos(), "nil slice in return statement")
		}
	}
}

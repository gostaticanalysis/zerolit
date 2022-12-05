package zerolit

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "zerolit finds return zero values but they are not literal"

// Analyzer finds return zero values but they are not literal.
var Analyzer = &analysis.Analyzer{
	Name: "zerolit",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.ReturnStmt)(nil),
		(*ast.AssignStmt)(nil),
		(*ast.FuncType)(nil),
	}

	isNG := make(map[types.Object]bool)
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.ReturnStmt:
			checkReturnStmt(pass, isNG, n)
		case *ast.AssignStmt:
			checkAssignStmt(pass, isNG, n)
		case *ast.FuncType:
			checkFuncType(pass, isNG, n)
		}
	})

	return nil, nil
}

func checkReturnStmt(pass *analysis.Pass, isNG map[types.Object]bool, ret *ast.ReturnStmt) {
	for _, v := range ret.Results {
		id, _ := v.(*ast.Ident)
		obj := pass.TypesInfo.ObjectOf(id)
		if obj == nil {
			continue
		}

		for defID, def := range pass.TypesInfo.Defs {
			if id != defID && def == obj {
				if isNG[obj] {
					pass.Reportf(v.Pos(), "zero value should return as a literal")
				}
				isNG[obj] = true
				break
			}
		}
	}
}

func checkAssignStmt(pass *analysis.Pass, isNG map[types.Object]bool, assign *ast.AssignStmt) {
	for _, v := range assign.Lhs {
		id, _ := v.(*ast.Ident)
		obj := pass.TypesInfo.ObjectOf(id)
		if obj == nil {
			continue
		}

		for useID, use := range pass.TypesInfo.Uses {
			if id != useID && use == obj {
				if isNG[obj] {
					pass.Reportf(useID.Pos(), "zero value should return as a literal")
				}
				isNG[obj] = true
				break
			}
		}
	}
}

func checkFuncType(pass *analysis.Pass, isNG map[types.Object]bool, funcType *ast.FuncType) {
	if funcType.Results == nil {
		return
	}

	for _, ret := range funcType.Results.List {
		for _, id := range ret.Names {
			obj := pass.TypesInfo.ObjectOf(id)
			if obj == nil {
				continue
			}

			for useID, use := range pass.TypesInfo.Uses {
				if id != useID && use == obj {
					if isNG[obj] {
						pass.Reportf(useID.Pos(), "zero value should return as a literal")
					}
					isNG[obj] = true
					break
				}
			}
		}
	}
}

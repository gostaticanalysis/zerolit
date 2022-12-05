package zerolit

import (
	"go/ast"
	"go/constant"
	"go/types"
	"reflect"

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
	}

	isNonZero := make(map[types.Object]bool)
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.AssignStmt:
			checkAssignStmt(pass, isNonZero, n)
		}
	})

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.ReturnStmt:
			checkReturnStmt(pass, isNonZero, n)
		}
	})

	return nil, nil
}

func checkAssignStmt(pass *analysis.Pass, isNonZero map[types.Object]bool, assign *ast.AssignStmt) {
	if len(assign.Lhs) != len(assign.Rhs) {
		return
	}

	for i := range assign.Lhs {
		id, _ := assign.Lhs[i].(*ast.Ident)
		obj := pass.TypesInfo.ObjectOf(id)
		if obj != nil && !isZero(pass, obj.Type(), assign.Rhs[i]) {
			isNonZero[obj] = true
		}
	}
}

func isZero(pass *analysis.Pass, typ types.Type, v ast.Expr) bool {
	switch typ.Underlying().(type) {
	case *types.Basic:
		tv := pass.TypesInfo.Types[v]
		return tv.Value != nil && reflect.ValueOf(constant.Val(tv.Value)).IsZero()
	case *types.Pointer, *types.Slice, *types.Map, *types.Chan, *types.Signature, *types.Interface:
		return types.Identical(pass.TypesInfo.TypeOf(v).Underlying(), types.Typ[types.UntypedNil])
	case *types.Struct, *types.Array:
		clit, _ := v.(*ast.CompositeLit)
		return clit != nil && len(clit.Elts) == 0
	}
	return false
}

func checkReturnStmt(pass *analysis.Pass, isNonZero map[types.Object]bool, ret *ast.ReturnStmt) {
	for _, v := range ret.Results {
		id, _ := v.(*ast.Ident)
		obj := pass.TypesInfo.ObjectOf(id)
		if obj == nil {
			continue
		}

		for defID, def := range pass.TypesInfo.Defs {
			if id != defID && def == obj && !isNonZero[obj] {
				pass.Reportf(v.Pos(), "zero value should return as a literal")
				break
			}
		}
	}
}

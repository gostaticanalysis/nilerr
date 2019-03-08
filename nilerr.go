package nilerr

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/gostaticanalysis/comment"
	"github.com/gostaticanalysis/comment/passes/commentmap"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/ssa"
)

// flags
var annotation = "return nil"

func init() {
	Analyzer.Flags.StringVar(&annotation, "annotation", annotation, "annotation for explicit return nil")
}

var Analyzer = &analysis.Analyzer{
	Name: "nilerr",
	Doc:  Doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		buildssa.Analyzer,
		commentmap.Analyzer,
	},
}

const Doc = "nilerr checks returning nil when err is not nil"

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	funcs := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).SrcFuncs
	cmaps := pass.ResultOf[commentmap.Analyzer].(comment.Maps)

	returns := map[token.Pos]ast.Node{}
	nodeFilter := []ast.Node{
		(*ast.ReturnStmt)(nil),
	}

	inspector.Preorder(nodeFilter, func(n ast.Node) {
		returns[n.Pos()] = n
	})

	for i := range funcs {
		for _, b := range funcs[i].Blocks {
			if errIsNotNil(b) {
				if ret := isReturnNil(b.Succs[0]); ret != nil {
					n, ok := returns[ret.Pos()]
					if ok && !cmaps.Annotated(n, annotation) {
						pass.Reportf(ret.Pos(), "error is not nil but it returns nil")
					}
				}
			}
		}
	}

	return nil, nil
}

var errType = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

func errIsNotNil(b *ssa.BasicBlock) bool {
	if len(b.Instrs) == 0 {
		return false
	}

	ifinst, ok := b.Instrs[len(b.Instrs)-1].(*ssa.If)
	if !ok {
		return false
	}

	binop, ok := ifinst.Cond.(*ssa.BinOp)
	if !ok {
		return false
	}

	if binop.Op != token.NEQ {
		return false
	}

	if !types.Implements(binop.X.Type(), errType) {
		return false
	}

	if !types.Implements(binop.Y.Type(), errType) {
		return false
	}

	xIsConst, yIsConst := isConst(binop.X), isConst(binop.Y)
	if (!xIsConst && !yIsConst) || (xIsConst && yIsConst) {
		return false
	}

	return true
}

func isConst(v ssa.Value) bool {
	_, ok := v.(*ssa.Const)
	return ok
}

func isReturnNil(b *ssa.BasicBlock) *ssa.Return {
	if len(b.Instrs) == 0 {
		return nil
	}

	ret, ok := b.Instrs[len(b.Instrs)-1].(*ssa.Return)
	if !ok {
		return nil
	}

	if len(ret.Results) == 0 {
		return nil
	}

	v, ok := ret.Results[len(ret.Results)-1].(*ssa.Const)
	if !ok {
		return nil
	}

	if !types.Implements(v.Type(), errType) {
		return nil
	}

	if !v.IsNil() {
		return nil
	}

	return ret
}

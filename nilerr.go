package nilerr

import (
	"go/token"
	"go/types"

	"github.com/gostaticanalysis/comment"
	"github.com/gostaticanalysis/comment/passes/commentmap"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
)

var Analyzer = &analysis.Analyzer{
	Name: "nilerr",
	Doc:  Doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
		commentmap.Analyzer,
	},
}

const Doc = "nilerr checks returning nil when err is not nil"

func run(pass *analysis.Pass) (interface{}, error) {
	funcs := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).SrcFuncs
	cmaps := pass.ResultOf[commentmap.Analyzer].(comment.Maps)

	for i := range funcs {
		for _, b := range funcs[i].Blocks {
			if v := binOpErrNil(b, token.NEQ); v != nil {
				if ret := isReturnNil(b.Succs[0]); ret != nil {
					pos := ret.Pos()
					line := pass.Fset.File(pos).Line(pos)
					if !cmaps.IgnoreLine(pass.Fset, line, "nilerr") {
						pass.Reportf(pos, "error is not nil but it returns nil")
					}
				}
			} else if v := binOpErrNil(b, token.EQL); v != nil {
				if ret := isReturnError(b.Succs[0], v); ret != nil {
					pos := ret.Pos()
					line := pass.Fset.File(pos).Line(pos)
					if !cmaps.IgnoreLine(pass.Fset, line, "nilerr") {
						pass.Reportf(pos, "error is nil but it returns error")
					}
				}
			}

		}
	}

	return nil, nil
}

var errType = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

func binOpErrNil(b *ssa.BasicBlock, op token.Token) ssa.Value {
	if len(b.Instrs) == 0 {
		return nil
	}

	ifinst, ok := b.Instrs[len(b.Instrs)-1].(*ssa.If)
	if !ok {
		return nil
	}

	binop, ok := ifinst.Cond.(*ssa.BinOp)
	if !ok {
		return nil
	}

	if binop.Op != op {
		return nil
	}

	if !types.Implements(binop.X.Type(), errType) {
		return nil
	}

	if !types.Implements(binop.Y.Type(), errType) {
		return nil
	}

	xIsConst, yIsConst := isConst(binop.X), isConst(binop.Y)
	switch {
	case !xIsConst && yIsConst: // err != nil or err == nil
		return binop.X
	case xIsConst && !yIsConst: // nil != err or nil == err
		return binop.Y
	}

	return nil
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

func isReturnError(b *ssa.BasicBlock, errVal ssa.Value) *ssa.Return {
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

	v := ret.Results[len(ret.Results)-1]
	if v != errVal {
		return nil
	}

	return ret
}

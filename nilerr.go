package nilerr

import (
	"fmt"
	"go/token"
	"go/types"
	"slices"

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

func run(pass *analysis.Pass) (any, error) {
	funcs := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).SrcFuncs
	cmaps := pass.ResultOf[commentmap.Analyzer].(comment.Maps)

	reportFail := func(v ssa.Value, ret *ssa.Return, format string) {
		pos := ret.Pos()
		if !cmaps.IgnorePos(pos, "nilerr") {
			seen := map[string]struct{}{}
			errLines := getValueLineNumbers(pass, v, seen)

			var errLineText string
			switch len(errLines) {
			case 0: // unknown position, e.g. a captured variable
			case 1:
				errLineText = fmt.Sprintf(" (line %d)", errLines[0])
			default:
				errLineText = fmt.Sprintf(" (lines %v)", errLines)
			}
			pass.Reportf(pos, format, errLineText)
		}
	}

	for i := range funcs {
		for _, b := range funcs[i].Blocks {
			if v := binOpErrNil(b, token.NEQ); v != nil {
				if ret := isReturnNil(b.Succs[0]); ret != nil {
					if !usesErrorValue(b.Succs[0], v) {
						reportFail(v, ret, "error is not nil%s but it returns nil")
					}
				}
			} else if v := binOpErrNil(b, token.EQL); v != nil {
				if len(b.Succs[0].Preds) == 1 { // if there are multiple conditions, this may be false positive
					if ret := isReturnError(b.Succs[0], v); ret != nil {
						reportFail(v, ret, "error is nil%s but it returns error")
					}
				}
			}

		}
	}

	return nil, nil
}

// getValueLineNumbers returns the line numbers.
// `seen` is used to avoid infinite loop.
func getValueLineNumbers(pass *analysis.Pass, v ssa.Value, seen map[string]struct{}) []int {
	if phi, ok := v.(*ssa.Phi); ok {
		result := make([]int, 0, len(phi.Edges))

		for _, edge := range phi.Edges {
			if _, ok := seen[edge.Name()]; ok {
				if edge.Pos() == token.NoPos {
					// Skip elements without a position.
					continue
				}

				result = append(result, pass.Fset.File(edge.Pos()).Line(edge.Pos()))
				continue
			}

			seen[edge.Name()] = struct{}{}

			result = append(result, getValueLineNumbers(pass, edge, seen)...)
		}

		slices.Sort(result)

		return result
	}

	value := v
	if extract, ok := value.(*ssa.Extract); ok {
		value = extract.Tuple
	}

	pos := value.Pos()

	if pos == token.NoPos {
		return nil
	}

	return []int{pass.Fset.File(pos).Line(pos)}
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

	errorReturnValues := 0
	for _, res := range ret.Results {
		if !types.Implements(res.Type(), errType) {
			continue
		}

		errorReturnValues++
		v, ok := deferredResult(b, res).(*ssa.Const)
		if !ok {
			return nil
		}

		if !v.IsNil() {
			return nil
		}
	}

	if errorReturnValues == 0 {
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

	if slices.ContainsFunc(ret.Results, func(v ssa.Value) bool { return deferredResult(b, v) == errVal }) {
		return ret
	}

	return nil
}

// deferredResult returns the value which is stored to the result variable
// just before the deferred functions are run.
//
// A function which has defer statements stores its results into result variables,
// runs the deferred functions, and returns the loaded values of the variables.
// If the value is not such a load, or the variable may be modified
// by a deferred function (it is captured by a closure), v is returned as it is.
func deferredResult(b *ssa.BasicBlock, v ssa.Value) ssa.Value {
	load, ok := v.(*ssa.UnOp)
	if !ok || load.Op != token.MUL {
		return v
	}

	alloc, ok := load.X.(*ssa.Alloc)
	if !ok {
		return v
	}

	if referrers := alloc.Referrers(); referrers != nil {
		for _, r := range *referrers {
			switch r.(type) {
			case *ssa.Store, *ssa.UnOp, *ssa.DebugRef:
			default:
				return v // e.g. captured by a closure
			}
		}
	}

	runDefersIdx := slices.IndexFunc(b.Instrs, func(instr ssa.Instruction) bool {
		_, ok := instr.(*ssa.RunDefers)
		return ok
	})
	if runDefersIdx < 0 {
		return v
	}

	for i := runDefersIdx - 1; i >= 0; i-- {
		if store, ok := b.Instrs[i].(*ssa.Store); ok && store.Addr == alloc {
			return store.Val
		}
	}

	return v
}

func usesErrorValue(b *ssa.BasicBlock, errVal ssa.Value) bool {
	for _, instr := range b.Instrs {
		if callInstr, ok := instr.(*ssa.Call); ok {
			for _, arg := range callInstr.Call.Args {
				if isUsedInValue(arg, errVal) {
					return true
				}

				sliceArg, ok := arg.(*ssa.Slice)
				if ok {
					if isUsedInSlice(sliceArg, errVal) {
						return true
					}
				}
			}
			continue
		}

		// The error may also be handed over by a store (field, pointer, element),
		// a channel send, a go/defer statement or a closure binding.
		switch instr.(type) {
		case *ssa.Store, *ssa.Send, *ssa.Go, *ssa.Defer, *ssa.MakeClosure:
			var buf [10]*ssa.Value
			for _, operand := range instr.Operands(buf[:0]) {
				if *operand != nil && isUsedInValue(*operand, errVal) {
					return true
				}
			}

			// A closure captures the variable which holds the error.
			if closure, ok := instr.(*ssa.MakeClosure); ok {
				for _, binding := range closure.Bindings {
					alloc, ok := binding.(*ssa.Alloc)
					if !ok || alloc.Referrers() == nil {
						continue
					}
					if load, ok := errVal.(*ssa.UnOp); ok && load.Op == token.MUL && load.X == alloc {
						return true // the variable itself is captured
					}
					for _, r := range *alloc.Referrers() {
						if store, ok := r.(*ssa.Store); ok && store.Addr == alloc && isUsedInValue(store.Val, errVal) {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

type ReferrersHolder interface {
	Referrers() *[]ssa.Instruction
}

var _ ReferrersHolder = (ssa.Node)(nil)
var _ ReferrersHolder = (ssa.Value)(nil)

func isUsedInSlice(sliceArg *ssa.Slice, errVal ssa.Value) bool {
	var valueBuf [10]*ssa.Value
	operands := sliceArg.Operands(valueBuf[:0])

	var valuesToInspect []ssa.Value
	addValueForInspection := func(value ssa.Value) {
		if value != nil {
			valuesToInspect = append(valuesToInspect, value)
		}
	}

	var nodesToInspect []ssa.Node
	visitedNodes := map[ssa.Node]bool{}
	addNodeForInspection := func(node ssa.Node) {
		if !visitedNodes[node] {
			visitedNodes[node] = true
			nodesToInspect = append(nodesToInspect, node)
		}
	}
	addReferrersForInspection := func(h ReferrersHolder) {
		if h == nil {
			return
		}

		referrers := h.Referrers()
		if referrers == nil {
			return
		}

		for _, r := range *referrers {
			if node, ok := r.(ssa.Node); ok {
				addNodeForInspection(node)
			}
		}
	}

	for _, operand := range operands {
		addReferrersForInspection(*operand)
		addValueForInspection(*operand)
	}

	for i := 0; i < len(nodesToInspect); i++ {
		switch node := nodesToInspect[i].(type) {
		case *ssa.IndexAddr:
			addReferrersForInspection(node)
		case *ssa.Store:
			addValueForInspection(node.Val)
		}
	}

	for _, value := range valuesToInspect {
		if isUsedInValue(value, errVal) {
			return true
		}
	}
	return false
}

func isUsedInValue(value, lookedFor ssa.Value) bool {
	if value == lookedFor {
		return true
	}

	switch value := value.(type) {
	case *ssa.ChangeInterface:
		return isUsedInValue(value.X, lookedFor)
	case *ssa.MakeInterface:
		return isUsedInValue(value.X, lookedFor)
	case *ssa.Call:
		if value.Call.IsInvoke() {
			return isUsedInValue(value.Call.Value, lookedFor)
		}
	}

	return false
}

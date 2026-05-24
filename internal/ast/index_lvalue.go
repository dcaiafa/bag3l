package ast

import "github.com/dcaiafa/bag3l/internal/vm"

// IndexLValue is the assignable form of IndexExpr (e.g. `a[i] = x`). Unlike
// IndexExpr it is not an expression and always emits a reference to the indexed
// element. Its Target and Index remain regular (value) expressions.
type IndexLValue struct {
	PosImpl

	Target Expr
	Index  Expr
}

func (e *IndexLValue) IsLValue() {}

func (e *IndexLValue) RunPass(ctx *Context, pass Pass) {
	switch pass {
	case Check:
		if !CheckNoOptional(ctx, e.Target) {
			return
		}
	}

	ctx.RunPassChild(e, e.Target, pass)
	ctx.RunPassChild(e, e.Index, pass)

	switch pass {
	case Emit:
		emitter := ctx.Emitter()
		emitter.Emit(e.Pos(), vm.OpObjectGetRef, 0, 0)
	}
}

package ast

import "github.com/dcaiafa/bag3l/internal/vm"

type IndexExpr struct {
	PosImpl

	Target   Expr
	Index    Expr
	Optional bool
}

func (e *IndexExpr) isExpr() {}

func (e *IndexExpr) RunPass(ctx *Context, pass Pass) {
	switch pass {
	case Check:
		if !CheckNoOptional(ctx, e.Target) {
			return
		}
	case Emit:
		if e.Optional {
			PropagateOptional(e.Target)
		}
	}

	ctx.RunPassChild(e, e.Target, pass)
	ctx.RunPassChild(e, e.Index, pass)

	switch pass {
	case Emit:
		emitter := ctx.Emitter()
		if _, isLValue := ctx.Parent().(*LValue); isLValue {
			emitter.Emit(e.Pos(), vm.OpObjectGetRef, 0, 0)
		} else {
			var flags uint16
			if e.Optional {
				flags |= vm.OptionalIndexFlag
			}
			emitter.Emit(e.Pos(), vm.OpObjectGet, 0, flags)
		}
	}
}

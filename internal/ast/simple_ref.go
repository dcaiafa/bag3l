package ast

import (
	"github.com/dcaiafa/bag3l/internal/symbol"
	"github.com/dcaiafa/bag3l/internal/token"
)

type SimpleRef struct {
	PosImpl
	ID token.Token

	sym symbol.Symbol

	Import *symbol.Import
}

func (r *SimpleRef) isExpr() {}

func (r *SimpleRef) RunPass(ctx *Context, pass Pass) {
	switch pass {
	case Check:
		symName := r.ID.Str

		// The blank identifier is only valid as an assignment target, never as a
		// value to read.
		if symName == "_" {
			ctx.Failf(r.Pos(), "_ cannot be used as a value")
			return
		}

		r.sym = ctx.FindSymbol(symName)
		if r.sym == nil {
			ctx.Failf(r.Pos(), "Symbol %q not found.", symName)
			return
		}

		var ok bool
		r.Import, ok = r.sym.(*symbol.Import)
		if ok {
			// An import name is only valid as the target of a member access,
			// in value (MemberAccess) or assignable (MemberAccessLValue) form.
			switch ctx.Parent().(type) {
			case *MemberAccess, *MemberAccessLValue:
			default:
				ctx.Failf(
					r.Pos(),
					"%v is an import, and cannot be used as a value",
					r.ID.Str)
				return
			}
		}

	case Emit:
		if r.Import == nil {
			emitSymbolPush(r.Pos(), ctx.Emitter(), r.sym)
		}
	}
}

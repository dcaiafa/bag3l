package ast

import (
	"github.com/dcaiafa/bag3l/internal/symbol"
	"github.com/dcaiafa/bag3l/internal/token"
)

// SimpleRefLValue is the assignable form of SimpleRef: a bare identifier on the
// left-hand side of an assignment. Unlike SimpleRef it is not an expression and
// always emits a reference to the symbol.
type SimpleRefLValue struct {
	PosImpl
	ID token.Token

	sym symbol.Symbol
}

func (r *SimpleRefLValue) IsLValue() {}

func (r *SimpleRefLValue) RunPass(ctx *Context, pass Pass) {
	switch pass {
	case Check:
		symName := r.ID.Str

		r.sym = ctx.FindSymbol(symName)
		if r.sym == nil {
			ctx.Failf(r.Pos(), "Symbol %q not found.", symName)
			return
		}

		if _, ok := r.sym.(*symbol.Import); ok {
			ctx.Failf(
				r.Pos(),
				"%v is an import, and cannot be used as a value",
				r.ID.Str)
			return
		}

		if r.sym.ReadOnly() {
			ctx.Failf(
				r.Pos(),
				"%v is read-only and cannot be assigned to",
				r.ID.Str)
			return
		}

	case Emit:
		emitSymbolStore(r.Pos(), ctx.Emitter(), r.sym)
	}
}

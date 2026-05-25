package ast

import (
	"github.com/dcaiafa/bag3l/internal/symbol"
	"github.com/dcaiafa/bag3l/internal/token"
)

type PreResolvedReference struct {
	PosImpl

	// TODO: for consistency, rename every `sym` everywhere to `symbol`.
	Symbol symbol.Symbol
}

// TODO: maybe replace all these with a ExprImpl.
func (r *PreResolvedReference) isExpr() {}

func (r *PreResolvedReference) RunPass(ctx *Context, pass Pass) {
	switch pass {
	case Emit:
		emitSymbolPush(token.Pos{}, ctx.Emitter(), r.Symbol)
	}
}

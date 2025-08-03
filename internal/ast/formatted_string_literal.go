package ast

import (
	"slices"

	"github.com/dcaiafa/bag3l/internal/token"
)

type FormattedStringLiteral struct {
	PosImpl

	Parts []Expr
}

func (e *FormattedStringLiteral) isExpr() {}

func (e *FormattedStringLiteral) RunPass(ctx *Context, pass Pass) {
	if pass == Rewrite {
		e.rewrite()
	}

	RunPassChildren(ctx, e, e.Parts, pass)
}

func (e *FormattedStringLiteral) rewrite() {
	for i := 0; i < len(e.Parts)-1; i++ {
		if isStringLiteral(e.Parts[i]) && isStringLiteral(e.Parts[i+1]) {
			a := e.Parts[i].(*LiteralExpr)
			b := e.Parts[i+1].(*LiteralExpr)
			a.Val.Str += b.Val.Str
			e.Parts = slices.Delete(e.Parts, i+1, i+2)
			i--
		}
	}

	concat := &FuncCallExpr{
		Target: &SimpleRef{
			ID: token.Token{
				Type: token.String,
				Str:  "$concat",
			},
		},
		Args: e.Parts,
		RetN: 1,
	}

	e.Parts = []Expr{concat}
}

func isStringLiteral(e Expr) bool {
	l, ok := e.(*LiteralExpr)
	if !ok {
		return false
	}
	return l.Val.Type == token.String
}

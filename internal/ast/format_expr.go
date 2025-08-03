package ast

import "github.com/dcaiafa/bag3l/internal/token"

type FormatExpr struct {
	PosImpl
	Target Expr
	Format string
}

func (e *FormatExpr) isExpr() {}

func (e *FormatExpr) RunPass(ctx *Context, pass Pass) {
	if pass == Rewrite {
		formatExpr := &FuncCallExpr{
			Target: &SimpleRef{
				ID: token.Token{
					Type: token.String,
					Str:  "$format",
				},
			},
			Args: []Expr{
				e.Target,
				&LiteralExpr{
					Val: token.Token{
						Type: token.String,
						Str:  e.Format,
					},
				},
			},
			RetN: 1,
		}
		e.Target = formatExpr
	}

	ctx.RunPassChild(e, e.Target, pass)
}

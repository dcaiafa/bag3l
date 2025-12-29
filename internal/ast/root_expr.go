package ast

import "github.com/dcaiafa/bag3l/internal/token"

type RootExprType int

const (
	RootExprMemberAccess RootExprType = iota
	RootExprIndex
)

type RootExpr struct {
	PosImpl

	Type RootExprType

	Target *RootExpr

	// Member access

	Member token.Token

	// Index

	Index Expr

	// Slice

	Begin Expr
	End   Expr

	Optional bool

	rewritten Expr
}

func (e *RootExpr) isExpr() {}

func (e *RootExpr) RunPass(ctx *Context, pass Pass) {
	if pass == Rewrite {
		fn := WithPos(e, &AnonFuncExpr{})
		fn.Params = []*FuncParam{{Name: "$"}}
		fn.Block = WithPos(e, &StmtBlock{
			Stmts: ASTs{
				WithPos(e, &ReturnStmt{
					Values: Exprs{
						e.rewrite(ctx),
					},
				}),
			},
		})
		e.rewritten = fn
	}

	ctx.RunPassChild(e, e.rewritten, pass)
}

func (e *RootExpr) rewrite(ctx *Context) Expr {
	var target Expr
	if e.Target != nil {
		target = e.Target.rewrite(ctx)
	} else {
		target = WithPos(e, &SimpleRef{
			ID: token.StringToken("$"),
		})
	}

	switch e.Type {
	case RootExprMemberAccess:
		return WithPos(e, &MemberAccess{
			Target:   target,
			Member:   e.Member,
			Optional: e.Optional,
		})
	case RootExprIndex:
		return WithPos(e, &IndexExpr{
			Target:   target,
			Index:    e.Index,
			Optional: e.Optional,
		})
	default:
		panic("not reached")
	}
}

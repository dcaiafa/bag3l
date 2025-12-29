package ast

import "github.com/dcaiafa/bag3l/internal/errlogger"

type PipeExpr struct {
	PosImpl

	Left  Expr
	Right Expr

	rewritten *FuncCallExpr
}

func (e *PipeExpr) isExpr() {}

func (e *PipeExpr) SetRetN(logger errlogger.ErrLogger, n int) {
	e.rewritten.SetRetN(logger, n)
}

func (e *PipeExpr) RunPass(ctx *Context, pass Pass) {
	if pass == Rewrite {
		e.rewrite(ctx)
	}

	ctx.RunPassChild(e, e.rewritten, pass)
}

func (e *PipeExpr) rewrite(ctx *Context) {
	if funcCall, ok := Unwrap(e.Right).(*FuncCallExpr); ok {
		funcCall.SetPos(e.Pos())
		funcCall.Args = append(Exprs{e.Left}, funcCall.Args...)
		funcCall.Pipeline = true
		e.rewritten = funcCall
		e.Left = nil
		e.Right = nil
		return
	}

	e.rewritten = &FuncCallExpr{
		Target:   e.Right,
		Args:     Exprs{e.Left},
		RetN:     1,
		Pipeline: true,
	}
	e.rewritten.SetPos(e.Pos())

	e.Left = nil
	e.Right = nil
}

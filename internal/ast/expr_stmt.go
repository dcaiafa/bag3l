package ast

type ExprStmt struct {
	PosImpl
	Expr Expr
}

func (c *ExprStmt) RunPass(ctx *Context, pass Pass) {
	switch pass {
	case Check:
		switch expr := c.Expr.(type) {
		case MultiValueExpr:
			expr.SetRetN(ctx, 0)
		default:
			ctx.Failf(c.Pos(), "Expression not allowed as a statement")
		}
	}

	ctx.Push(c)
	c.Expr.RunPass(ctx, pass)
	ctx.Pop()
}

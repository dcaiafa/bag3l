package ast

type Repeatable struct {
	PosImpl

	Stmt AST
}

func (r *Repeatable) RunPass(ctx *Context, pass Pass) {
	ctx.RunPassChild(r, r.Stmt, pass)
}

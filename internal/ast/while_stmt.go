package ast

type WhileStmt struct {
	PosImpl

	Predicate Expr
	Block     *StmtBlock

	loop *Loop
}

func (s *WhileStmt) RunPass(ctx *Context, pass Pass) {
	if pass == Rewrite {
		s.loop = &Loop{
			PosImpl:   s.PosImpl,
			Predicate: s.Predicate,
			Block:     s.Block,
		}

		s.Predicate = nil
		s.Block = nil
	}

	ctx.RunPassChild(s, s.loop, pass)
}

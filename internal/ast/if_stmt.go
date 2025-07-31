package ast

import "github.com/dcaiafa/bag3l/internal/vm"

type IfStmt struct {
	PosImpl
	Sections []*IfSection

	end *vm.Label
}

func (s *IfStmt) RunPass(ctx *Context, pass Pass) {
	switch pass {
	case Emit:
		s.end = ctx.Emitter().NewLabel()
	}

	RunPassChildren(ctx, s, s.Sections, pass)

	switch pass {
	case Emit:
		ctx.Emitter().ResolveLabel(s.end)
	}
}

type IfSection struct {
	PosImpl
	Pred  Expr
	Block AST

	end *vm.Label
}

func (b *IfSection) RunPass(ctx *Context, pass Pass) {
	switch pass {
	case Emit:
		b.end = ctx.Emitter().NewLabel()
	}

	// The "else" block does not have a predicate expression.
	if b.Pred != nil {
		ctx.RunPassChild(b, b.Pred, pass)
	}

	switch pass {
	case Emit:
		if b.Pred != nil {
			ctx.Emitter().EmitJump(b.Pos(), vm.OpJumpIfFalse, b.end, 0)
		}
	}

	ctx.RunPassChild(b, b.Block, pass)

	switch pass {
	case Emit:
		ifStmtEnd := ctx.Parent().(*IfStmt).end
		ctx.Emitter().EmitJump(b.Pos(), vm.OpJump, ifStmtEnd, 0)
		ctx.Emitter().ResolveLabel(b.end)
	}
}

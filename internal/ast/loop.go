package ast

import (
	"github.com/dcaiafa/bag3l/internal/token"
	"github.com/dcaiafa/bag3l/internal/vm"
)

type Loop struct {
	PosImpl

	Predicate     Expr
	Block         AST
	ContinueBlock AST

	begin *vm.Label
	cont  *vm.Label
	end   *vm.Label
}

func (l *Loop) IsRepeatableScope() {}

func (l *Loop) RunPass(ctx *Context, pass Pass) {
	if pass == Rewrite {
		if l.Block != nil {
			l.Block = &Repeatable{
				Stmt: l.Block,
			}
		}
		if l.ContinueBlock != nil {
			l.ContinueBlock = &Repeatable{
				Stmt: l.ContinueBlock,
			}
		}
	}

	var emitter *vm.Emitter

	if pass == Emit {
		emitter = ctx.Emitter()
		l.begin = emitter.NewLabel()
		l.cont = emitter.NewLabel()
		l.end = emitter.NewLabel()
		emitter.ResolveLabel(l.begin)
	}

	ctx.RunPassChild(l, l.Predicate, pass)

	if pass == Emit {
		emitter.EmitJump(l.Predicate.Pos(), vm.OpJumpIfFalse, l.end, 0)
	}

	ctx.RunPassChild(l, l.Block, pass)

	if pass == Emit {
		emitter.ResolveLabel(l.cont)
	}

	if l.ContinueBlock != nil {
		ctx.RunPassChild(l, l.ContinueBlock, pass)
	}

	if pass == Emit {
		emitter.EmitJump(l.Block.Pos(), vm.OpJump, l.begin, 0)
		emitter.ResolveLabel(l.end)
	}
}

func (l *Loop) EmitBreak(pos token.Pos, e *vm.Emitter) {
	e.EmitJump(pos, vm.OpJump, l.end, 0)
}

func (l *Loop) EmitContinue(pos token.Pos, e *vm.Emitter) {
	e.EmitJump(pos, vm.OpJump, l.cont, 0)
}

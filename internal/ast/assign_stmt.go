package ast

type AssignStmt struct {
	PosImpl
	Lvalues []LValue
	Rvalues Exprs
}

func (s *AssignStmt) RunPass(ctx *Context, pass Pass) {
	switch pass {
	case Check:
		if len(s.Lvalues) != len(s.Rvalues) {
			if funcCall, ok := Unwrap(s.Rvalues[0]).(*FuncCallExpr); ok && len(s.Rvalues) == 1 {
				funcCall.RetN = len(s.Lvalues)
			} else {
				ctx.Failf(
					s.Pos(),
					"Left side of assignment expects %v values, "+
						"but right side produces %v values",
					len(s.Lvalues), len(s.Rvalues))
			}
		}
	}

	if pass == Emit {
		// Evaluate the right-hand side first, leaving v1..vN on the stack, then
		// store into the lvalues right-to-left. Processing in reverse exposes
		// each value on top of the stack just as its lvalue's prefix
		// (container/key) is pushed above it, so the store opcode finds the
		// value directly beneath the prefix.
		ctx.RunPassChild(s, s.Rvalues, pass)
		for i := len(s.Lvalues) - 1; i >= 0; i-- {
			ctx.RunPassChild(s, s.Lvalues[i], pass)
		}
		return
	}

	RunPassChildren(ctx, s, s.Lvalues, pass)
	ctx.RunPassChild(s, s.Rvalues, pass)
}

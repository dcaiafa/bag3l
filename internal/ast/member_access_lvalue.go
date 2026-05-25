package ast

import (
	"github.com/dcaiafa/bag3l/internal/symbol"
	"github.com/dcaiafa/bag3l/internal/token"
	"github.com/dcaiafa/bag3l/internal/vm"
)

// MemberAccessLValue is the assignable form of MemberAccess (e.g. `a.b = x`).
// Unlike MemberAccess it is not an expression and always emits a reference to
// the member. Its Target remains a regular (value) expression.
type MemberAccessLValue struct {
	PosImpl

	Target Expr
	Member token.Token

	ModuleMember symbol.Symbol
}

func (a *MemberAccessLValue) IsLValue() {}

func (a *MemberAccessLValue) RunPass(ctx *Context, pass Pass) {
	switch pass {
	case Check:
		if !CheckNoOptional(ctx, a.Target) {
			return
		}
	}

	ctx.RunPassChild(a, a.Target, pass)

	switch pass {
	case Check:
		simpleRef, ok := a.Target.(*SimpleRef)
		if ok {
			if simpleRef.Import != nil {
				a.ModuleMember = simpleRef.Import.GetSymbol(a.Member.Str)
				if a.ModuleMember == nil {
					ctx.Failf(a.Pos(), "%v not declared by package %v",
						a.Member.Str,
						simpleRef.Import.Name())
					return
				}
			}
		}

	case Emit:
		emitter := ctx.Emitter()
		if a.ModuleMember == nil {
			// The value to store was pushed by AssignStmt before the Target;
			// push the member name as the key so the value sits beneath the
			// container/key pair expected by OpStoreIndex.
			emitter.Emit(
				a.Pos(), vm.OpLoadLiteral,
				uint32(emitter.AddLiteral(vm.NewString(a.Member.Str))), 0)
			emitter.Emit(a.Pos(), vm.OpStoreIndex, 0, 0)
		} else {
			ctx.Failf(a.Pos(), "cannot assign to module")
			return
		}
	}
}

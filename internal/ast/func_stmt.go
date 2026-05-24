package ast

import (
	"github.com/dcaiafa/bag3l/internal/scope"
	"github.com/dcaiafa/bag3l/internal/symbol"
)

type FuncStmt struct {
	Func

	Name string

	sym symbol.Symbol
}

func (s *FuncStmt) RunPass(ctx *Context, pass Pass) {
	parentFn := ctx.CurrentFunc()

	switch pass {
	case CreateGlobals:
		if parentFn == nil {
			s.sym = &symbol.LiteralSymbol{}
			s.sym.SetReadOnly(true)
			s.sym.SetName(s.Name)
			s.sym.SetPos(s.Pos())

			if !ctx.GetScope(scope.Package).PutSymbol(ctx, s.sym) {
				return
			}
		}

	case Check:
		s.DebugName = s.Name

		if parentFn != nil {
			s.sym = parentFn.NewLocal()
			s.sym.SetName(s.Name)
			s.sym.SetPos(s.Pos())
			s.IsClosure = true

			if !ctx.GetScope(scope.Block).PutSymbol(ctx, s.sym) {
				return
			}
		}

	case Emit:
		if localSym, ok := s.sym.(*symbol.LocalVarSymbol); ok {
			emitVariableInit(ctx, s.Pos(), localSym)
		}
	}

	ctx.RunPassChild(s, &s.Func, pass)

	switch pass {
	case Check:
		if fnSym, ok := s.sym.(*symbol.LiteralSymbol); ok {
			fnSym.PackageIdx = 0 // Local package
			fnSym.LiteralIdx = s.IdxFunc()
		}

	case Emit:
		if localSym, ok := s.sym.(*symbol.LocalVarSymbol); ok {
			// `Func` emitted the closure, leaving it on top of the stack. Store
			// it into the local var.
			emitSymbolStore(s.Pos(), ctx.Emitter(), localSym)
		}
	}
}

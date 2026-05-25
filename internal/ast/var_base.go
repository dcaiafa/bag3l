package ast

import (
	"github.com/dcaiafa/bag3l/internal/symbol"
	"github.com/dcaiafa/bag3l/internal/token"
	"github.com/dcaiafa/bag3l/internal/vm"
)

func emitVariableInit(ctx *Context, pos token.Pos, sym symbol.Symbol) {
	if sym.Lifted() {
		switch sym := sym.(type) {
		case *symbol.LocalVarSymbol:
			ctx.Emitter().Emit(pos, vm.OpInitLiftedLocal, uint32(sym.LocalNdx), 0)
		default:
			panic("unreachable")
		}
	} else if ctx.IsInRepeatableScope() {
		switch sym := sym.(type) {
		case *symbol.LocalVarSymbol:
			ctx.Emitter().Emit(pos, vm.OpInitLocal, uint32(sym.LocalNdx), 0)
		case *symbol.GlobalVarSymbol:
			ctx.Emitter().Emit(pos, vm.OpInitGlobal, uint32(sym.GlobalNdx), 0)
		default:
			panic("unreachable")
		}
	}
}

func emitSymbolPush(pos token.Pos, emitter *vm.Emitter, sym symbol.Symbol) {
	switch sym := sym.(type) {
	case *symbol.LiteralSymbol:
		emitter.Emit(pos, vm.OpLoadLiteral, uint32(sym.LiteralIdx), uint16(sym.PackageIdx))

	case *symbol.GlobalVarSymbol:
		emitter.Emit(pos, vm.OpLoadGlobal, uint32(sym.GlobalNdx), 0)

	case *symbol.LocalVarSymbol:
		if sym.Lifted() {
			emitter.Emit(pos, vm.OpLoadLocalDeref, uint32(sym.LocalNdx), 0)
		} else {
			emitter.Emit(pos, vm.OpLoadLocal, uint32(sym.LocalNdx), 0)
		}

	case *symbol.CaptureSymbol:
		emitter.Emit(pos, vm.OpLoadCapture, uint32(sym.CaptureNdx), 0)

	case *symbol.ParamSymbol:
		if sym.Lifted() {
			emitter.Emit(pos, vm.OpLoadArgDeref, uint32(sym.ParamNdx), 0)
		} else {
			emitter.Emit(pos, vm.OpLoadArg, uint32(sym.ParamNdx), 0)
		}

	default:
		panic("not implemented")
	}
}

// emitSymbolStore stores the value currently on top of the stack into the
// variable denoted by sym, popping the value. Lifted locals/params are written
// through their boxes (ValueRef); everything else is written directly.
func emitSymbolStore(pos token.Pos, emitter *vm.Emitter, sym symbol.Symbol) {
	switch sym := sym.(type) {
	case *symbol.GlobalVarSymbol:
		emitter.Emit(pos, vm.OpStoreGlobal, uint32(sym.GlobalNdx), 0)

	case *symbol.LocalVarSymbol:
		if sym.Lifted() {
			emitter.Emit(pos, vm.OpStoreLocalDeref, uint32(sym.LocalNdx), 0)
		} else {
			emitter.Emit(pos, vm.OpStoreLocal, uint32(sym.LocalNdx), 0)
		}

	case *symbol.CaptureSymbol:
		emitter.Emit(pos, vm.OpStoreCapture, uint32(sym.CaptureNdx), 0)

	case *symbol.ParamSymbol:
		if sym.Lifted() {
			emitter.Emit(pos, vm.OpStoreArgDeref, uint32(sym.ParamNdx), 0)
		} else {
			emitter.Emit(pos, vm.OpStoreArg, uint32(sym.ParamNdx), 0)
		}

	default:
		panic("unreachable")
	}
}

// emitCaptureBoxPush pushes the box (ValueRef) of a captured variable onto the
// stack so it can be collected into a closure's or iterator's capture list.
// Captured symbols are always lifted, so the box already lives in the source
// slot.
func emitCaptureBoxPush(pos token.Pos, emitter *vm.Emitter, sym symbol.Symbol) {
	switch sym := sym.(type) {
	case *symbol.LocalVarSymbol:
		if !sym.Lifted() {
			panic("captured local is not lifted")
		}
		emitter.Emit(pos, vm.OpLoadLocal, uint32(sym.LocalNdx), 0)

	case *symbol.ParamSymbol:
		if !sym.Lifted() {
			panic("captured param is not lifted")
		}
		emitter.Emit(pos, vm.OpLoadArg, uint32(sym.ParamNdx), 0)

	case *symbol.CaptureSymbol:
		emitter.Emit(pos, vm.OpLoadCaptureBox, uint32(sym.CaptureNdx), 0)

	default:
		panic("unreachable")
	}
}

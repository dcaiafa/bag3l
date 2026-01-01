package parser

import (
	_i0 "github.com/dcaiafa/bag3l/internal/stub/ast"
)

var _rules = []int32{
	0, 1, 1, 2, 3, 4, 4, 4, 4, 4, 5, 6, 7, 8,
	9, 10, 11, 12, 12, 13, 14, 15, 16, 16, 17, 18, 18, 19,
	19, 20, 20, 21, 21, 22, 22, 23, 23, 24, 24, 25, 25, 26,
	26, 27, 27, 28, 28, 29, 29,
}

var _termCounts = []int32{
	1, 1, 1, 2, 2, 1, 1, 1, 1, 1, 5, 1, 3, 5,
	2, 3, 2, 3, 1, 6, 4, 2, 3, 1, 2, 1, 0, 2,
	1, 1, 0, 2, 1, 1, 0, 1, 0, 3, 1, 1, 0, 1,
	0, 1, 0, 3, 1, 1, 0,
}

var _actions = []int32{
	74, 79, 82, 85, 98, 101, 104, 107, 110, 113, 116, 119, 132, 145,
	148, 161, 174, 187, 200, 213, 226, 229, 232, 235, 242, 245, 258, 283,
	286, 299, 304, 309, 314, 327, 332, 226, 335, 340, 343, 348, 361, 364,
	377, 382, 385, 390, 395, 418, 441, 458, 471, 484, 489, 492, 509, 226,
	512, 517, 520, 226, 533, 546, 559, 572, 579, 332, 584, 589, 594, 599,
	604, 226, 609, 622, 4, 1, 1, 4, 2, 2, 0, -2, 2, 18,
	19, 12, 2, 6, 0, -26, 3, 7, 8, 8, 6, 9, 5, 10,
	2, 0, 2147483647, 2, 0, -1, 2, 19, 20, 2, 19, 24, 2, 19,
	21, 2, 19, 22, 2, 19, 23, 12, 2, -5, 0, -5, 3, -5,
	8, -5, 6, -5, 5, -5, 12, 2, -28, 0, -28, 3, -28, 8,
	-28, 6, -28, 5, -28, 2, 0, -3, 12, 2, 6, 0, -25, 3,
	7, 8, 8, 6, 9, 5, 10, 12, 2, -9, 0, -9, 3, -9,
	8, -9, 6, -9, 5, -9, 12, 2, -6, 0, -6, 3, -6, 8,
	-6, 6, -6, 5, -6, 12, 2, -8, 0, -8, 3, -8, 8, -8,
	6, -8, 5, -8, 12, 2, -7, 0, -7, 3, -7, 8, -7, 6,
	-7, 5, -7, 12, 2, -4, 0, -4, 3, -4, 8, -4, 6, -4,
	5, -4, 2, 19, 26, 2, 18, 28, 2, 10, 29, 6, 19, -34,
	12, 30, 18, -34, 2, 14, 33, 12, 2, -27, 0, -27, 3, -27,
	8, -27, 6, -27, 5, -27, 24, 11, -48, 16, -48, 2, -48, 15,
	-48, 0, -48, 9, -48, 3, -48, 19, -48, 8, -48, 17, 46, 6,
	-48, 5, -48, 2, 9, 34, 12, 2, -12, 0, -12, 3, -12, 8,
	-12, 6, -12, 5, -12, 4, 11, -30, 19, 35, 4, 19, -33, 18,
	-33, 4, 19, 39, 18, 40, 12, 2, -15, 0, -15, 3, -15, 8,
	-15, 6, -15, 5, -15, 4, 15, -36, 19, 44, 2, 18, 48, 4,
	11, -32, 19, -32, 2, 11, 50, 4, 11, -29, 19, 35, 12, 2,
	-18, 0, -18, 3, -18, 8, -18, 6, -18, 5, -18, 2, 13, 52,
	12, 2, -16, 0, -16, 3, -16, 8, -16, 6, -16, 5, -16, 4,
	16, 57, 15, -35, 2, 15, 53, 4, 20, 54, 19, -42, 4, 16,
	-38, 15, -38, 22, 11, -47, 16, -47, 2, -47, 15, -47, 0, -47,
	9, -47, 3, -47, 19, -47, 8, -47, 6, -47, 5, -47, 22, 11,
	-24, 16, -24, 2, -24, 15, -24, 0, -24, 9, -24, 3, -24, 19,
	-24, 8, -24, 6, -24, 5, -24, 16, 16, -11, 2, -11, 15, -11,
	0, -11, 3, -11, 8, -11, 6, -11, 5, -11, 12, 2, -10, 0,
	-10, 3, -10, 8, -10, 6, -10, 5, -10, 12, 2, -13, 0, -13,
	3, -13, 8, -13, 6, -13, 5, -13, 4, 11, -14, 19, -14, 2,
	19, 58, 16, 2, -40, 0, -40, 3, -40, 19, 26, 8, -40, 14,
	59, 6, -40, 5, -40, 2, 19, -41, 4, 11, -31, 19, -31, 2,
	19, 44, 12, 2, -17, 0, -17, 3, -17, 8, -17, 6, -17, 5,
	-17, 12, 2, -39, 0, -39, 3, -39, 8, -39, 6, -39, 5, -39,
	12, 2, -19, 0, -19, 3, -19, 8, -19, 6, -19, 5, -19, 12,
	2, -23, 0, -23, 3, -23, 8, -23, 6, -23, 5, -23, 6, 16,
	-44, 15, -44, 9, 65, 4, 16, -37, 15, -37, 4, 16, -43, 15,
	-43, 4, 16, -20, 15, -20, 4, 16, 71, 15, 72, 4, 16, -46,
	15, -46, 4, 16, -21, 15, -21, 12, 2, -22, 0, -22, 3, -22,
	8, -22, 6, -22, 5, -22, 4, 16, -45, 15, -45,
}

var _goto = []int32{
	74, 81, 81, 82, 81, 81, 81, 81, 81, 81, 81, 81, 81, 81,
	99, 81, 81, 81, 81, 81, 112, 81, 81, 115, 81, 81, 120, 81,
	81, 123, 81, 130, 81, 133, 140, 143, 81, 81, 146, 81, 81, 81,
	81, 81, 149, 81, 81, 81, 81, 81, 81, 81, 81, 152, 81, 159,
	81, 162, 81, 165, 81, 81, 81, 170, 81, 175, 81, 81, 81, 81,
	81, 178, 81, 81, 6, 3, 3, 1, 4, 2, 5, 0, 16, 5,
	11, 4, 12, 18, 13, 19, 14, 13, 15, 7, 16, 8, 17, 10,
	18, 12, 5, 11, 4, 25, 13, 15, 7, 16, 8, 17, 10, 18,
	2, 17, 27, 4, 22, 31, 11, 32, 2, 29, 47, 6, 9, 36,
	20, 37, 21, 38, 2, 12, 41, 6, 24, 42, 23, 43, 14, 45,
	2, 6, 49, 2, 17, 51, 2, 9, 56, 2, 26, 55, 6, 16,
	60, 25, 61, 17, 62, 2, 17, 63, 2, 14, 64, 4, 28, 68,
	17, 69, 4, 15, 66, 27, 67, 2, 6, 70, 2, 17, 73,
}

type _Bounds struct {
	Begin Token
	End   Token
	Empty bool
}

func _cast[T any](v any) T {
	cv, _ := v.(T)
	return cv
}

type Error struct {
	Token    Token
	Expected []int
}

func _Find(table []int32, y, x int32) (int32, bool) {
	i := int(table[int(y)])
	count := int(table[i])
	i++
	end := i + count
	for ; i < end; i += 2 {
		if table[i] == x {
			return table[i+1], true
		}
	}
	return 0, false
}

type _Lexer interface {
	ReadToken() (Token, int)
}

type _item struct {
	State int32
	Sym   any
}

type lox struct {
	_lex   _Lexer
	_stack _Stack[_item]

	_la    int
	_lasym any

	_qla    int
	_qlasym any
}

func (p *parser) parse(lex _Lexer) bool {
	const accept = 2147483647

	p._lex = lex
	p._qla = -1
	p._stack.Push(_item{})

	p._readToken()

	for {
		topState := p._stack.Peek(0).State
		action, ok := _Find(_actions, topState, int32(p._la))
		if !ok {
			if !p._recover() {
				return false
			}
			continue
		}
		if action == accept {
			break
		} else if action >= 0 { // shift
			p._stack.Push(_item{
				State: action,
				Sym:   p._lasym,
			})
			p._readToken()
		} else { // reduce
			prod := -action
			termCount := _termCounts[int(prod)]
			rule := _rules[int(prod)]
			res := p._act(prod)
			p._stack.Pop(int(termCount))
			topState = p._stack.Peek(0).State
			nextState, _ := _Find(_goto, topState, rule)
			p._stack.Push(_item{
				State: nextState,
				Sym:   res,
			})
		}
	}

	return true
}

// recoverLookahead can be called during an error production action (an action
// for a production that has a @error term) to recover the lookahead that was
// possibly lost in the process of reducing the error production.
func (p *parser) recoverLookahead(typ int, tok Token) {
	if p._qla != -1 {
		panic("recovered lookahead already pending")
	}

	p._qla = p._la
	p._qlasym = p._lasym
	p._la = typ
	p._lasym = tok
}

func (p *parser) _readToken() {
	if p._qla != -1 {
		p._la = p._qla
		p._lasym = p._qlasym
		p._qla = -1
		p._qlasym = nil
		return
	}

	p._lasym, p._la = p._lex.ReadToken()
	if p._la == ERROR {
		p._lasym = p._makeError()
	}
}

func (p *parser) _recover() bool {
	errSym, ok := p._lasym.(Error)
	if !ok {
		errSym = p._makeError()
	}

	for p._la == ERROR {
		p._readToken()
	}

	for {
		save := p._stack

		for len(p._stack) >= 1 {
			state := p._stack.Peek(0).State

			for {
				action, ok := _Find(_actions, state, int32(ERROR))
				if !ok {
					break
				}

				if action < 0 {
					prod := -action
					rule := _rules[int(prod)]
					state, _ = _Find(_goto, state, rule)
					continue
				}

				state = action

				_, ok = _Find(_actions, state, int32(p._la))
				if !ok {
					break
				}

				p._qla = p._la
				p._qlasym = p._lasym
				p._la = ERROR
				p._lasym = errSym
				return true
			}

			p._stack.Pop(1)
		}

		if p._la == EOF {
			return false
		}

		p._stack = save
		p._readToken()
	}
}

func (p *parser) _makeError() Error {
	e := Error{
		Token: p._lasym.(Token),
	}

	// Compile list of allowed tokens at this state.
	// See _Find for the format of the _actions table.
	s := p._stack.Peek(0).State
	i := int(_actions[int(s)])
	count := int(_actions[i])
	i++
	end := i + count
	for ; i < end; i += 2 {
		e.Expected = append(e.Expected, int(_actions[i]))
	}

	return e
}

func (p *parser) _act(prod int32) any {
	switch prod {
	case 1:
		return p.on_s(
			_cast[*_i0.Unit](p._stack.Peek(0).Sym),
		)
	case 2:
		return p.on_s__error(
			_cast[Error](p._stack.Peek(0).Sym),
		)
	case 3:
		return p.on_unit(
			_cast[*_i0.Unit](p._stack.Peek(1).Sym),
			_cast[[]_i0.AST](p._stack.Peek(0).Sym),
		)
	case 4:
		return p.on_package(
			_cast[Token](p._stack.Peek(1).Sym),
			_cast[Token](p._stack.Peek(0).Sym),
		)
	case 5:
		return p.on_decl(
			_cast[_i0.AST](p._stack.Peek(0).Sym),
		)
	case 6:
		return p.on_decl(
			_cast[_i0.AST](p._stack.Peek(0).Sym),
		)
	case 7:
		return p.on_decl(
			_cast[_i0.AST](p._stack.Peek(0).Sym),
		)
	case 8:
		return p.on_decl(
			_cast[_i0.AST](p._stack.Peek(0).Sym),
		)
	case 9:
		return p.on_decl(
			_cast[_i0.AST](p._stack.Peek(0).Sym),
		)
	case 10:
		return p.on_const_decl(
			_cast[Token](p._stack.Peek(4).Sym),
			_cast[Token](p._stack.Peek(3).Sym),
			_cast[*_i0.TypeRef](p._stack.Peek(2).Sym),
			_cast[Token](p._stack.Peek(1).Sym),
			_cast[*_i0.ConstValue](p._stack.Peek(0).Sym),
		)
	case 11:
		return p.on_const_value(
			_cast[Token](p._stack.Peek(0).Sym),
		)
	case 12:
		return p.on_import_decl(
			_cast[Token](p._stack.Peek(2).Sym),
			_cast[Token](p._stack.Peek(1).Sym),
			_cast[Token](p._stack.Peek(0).Sym),
		)
	case 13:
		return p.on_struct_decl(
			_cast[Token](p._stack.Peek(4).Sym),
			_cast[Token](p._stack.Peek(3).Sym),
			_cast[Token](p._stack.Peek(2).Sym),
			_cast[[]*_i0.StructField](p._stack.Peek(1).Sym),
			_cast[Token](p._stack.Peek(0).Sym),
		)
	case 14:
		return p.on_struct_field(
			_cast[Token](p._stack.Peek(1).Sym),
			_cast[*_i0.TypeRef](p._stack.Peek(0).Sym),
		)
	case 15:
		return p.on_type_decl(
			_cast[Token](p._stack.Peek(2).Sym),
			_cast[Token](p._stack.Peek(1).Sym),
			_cast[*_i0.GoType](p._stack.Peek(0).Sym),
		)
	case 16:
		return p.on_go_type(
			_cast[Token](p._stack.Peek(1).Sym),
			_cast[*_i0.GoType](p._stack.Peek(0).Sym),
		)
	case 17:
		return p.on_simple_go_type__full(
			_cast[Token](p._stack.Peek(2).Sym),
			_cast[Token](p._stack.Peek(1).Sym),
			_cast[Token](p._stack.Peek(0).Sym),
		)
	case 18:
		return p.on_simple_go_type(
			_cast[Token](p._stack.Peek(0).Sym),
		)
	case 19:
		return p.on_func_decl(
			_cast[Token](p._stack.Peek(5).Sym),
			_cast[Token](p._stack.Peek(4).Sym),
			_cast[Token](p._stack.Peek(3).Sym),
			_cast[[]*_i0.FuncParam](p._stack.Peek(2).Sym),
			_cast[Token](p._stack.Peek(1).Sym),
			_cast[[]*_i0.TypeRef](p._stack.Peek(0).Sym),
		)
	case 20:
		return p.on_func_param(
			_cast[Token](p._stack.Peek(3).Sym),
			_cast[Token](p._stack.Peek(2).Sym),
			_cast[*_i0.TypeRef](p._stack.Peek(1).Sym),
			_cast[*_i0.FuncParam](p._stack.Peek(0).Sym),
		)
	case 21:
		return p.on_func_param_default(
			_cast[Token](p._stack.Peek(1).Sym),
			_cast[*_i0.ConstValue](p._stack.Peek(0).Sym),
		)
	case 22:
		return p.on_func_rets__multi(
			_cast[Token](p._stack.Peek(2).Sym),
			_cast[[]*_i0.TypeRef](p._stack.Peek(1).Sym),
			_cast[Token](p._stack.Peek(0).Sym),
		)
	case 23:
		return p.on_func_rets__single(
			_cast[*_i0.TypeRef](p._stack.Peek(0).Sym),
		)
	case 24:
		return p.on_type_ref(
			_cast[Token](p._stack.Peek(1).Sym),
			_cast[Token](p._stack.Peek(0).Sym),
		)
	case 25: // ZeroOrMore
		return _cast[[]_i0.AST](p._stack.Peek(0).Sym)
	case 26: // ZeroOrMore
		{
			var zero []_i0.AST
			return zero
		}
	case 27: // OneOrMore
		return append(
			_cast[[]_i0.AST](p._stack.Peek(1).Sym),
			_cast[_i0.AST](p._stack.Peek(0).Sym),
		)
	case 28: // OneOrMore
		return []_i0.AST{
			_cast[_i0.AST](p._stack.Peek(0).Sym),
		}
	case 29: // ZeroOrMore
		return _cast[[]*_i0.StructField](p._stack.Peek(0).Sym)
	case 30: // ZeroOrMore
		{
			var zero []*_i0.StructField
			return zero
		}
	case 31: // OneOrMore
		return append(
			_cast[[]*_i0.StructField](p._stack.Peek(1).Sym),
			_cast[*_i0.StructField](p._stack.Peek(0).Sym),
		)
	case 32: // OneOrMore
		return []*_i0.StructField{
			_cast[*_i0.StructField](p._stack.Peek(0).Sym),
		}
	case 33: // ZeroOrOne
		return _cast[Token](p._stack.Peek(0).Sym)
	case 34: // ZeroOrOne
		{
			var zero Token
			return zero
		}
	case 35: // ZeroOrOne
		return _cast[[]*_i0.FuncParam](p._stack.Peek(0).Sym)
	case 36: // ZeroOrOne
		{
			var zero []*_i0.FuncParam
			return zero
		}
	case 37: // List
		return append(
			_cast[[]*_i0.FuncParam](p._stack.Peek(2).Sym),
			_cast[*_i0.FuncParam](p._stack.Peek(0).Sym),
		)
	case 38: // List
		return []*_i0.FuncParam{
			_cast[*_i0.FuncParam](p._stack.Peek(0).Sym),
		}
	case 39: // ZeroOrOne
		return _cast[[]*_i0.TypeRef](p._stack.Peek(0).Sym)
	case 40: // ZeroOrOne
		{
			var zero []*_i0.TypeRef
			return zero
		}
	case 41: // ZeroOrOne
		return _cast[Token](p._stack.Peek(0).Sym)
	case 42: // ZeroOrOne
		{
			var zero Token
			return zero
		}
	case 43: // ZeroOrOne
		return _cast[*_i0.FuncParam](p._stack.Peek(0).Sym)
	case 44: // ZeroOrOne
		{
			var zero *_i0.FuncParam
			return zero
		}
	case 45: // List
		return append(
			_cast[[]*_i0.TypeRef](p._stack.Peek(2).Sym),
			_cast[*_i0.TypeRef](p._stack.Peek(0).Sym),
		)
	case 46: // List
		return []*_i0.TypeRef{
			_cast[*_i0.TypeRef](p._stack.Peek(0).Sym),
		}
	case 47: // ZeroOrOne
		return _cast[Token](p._stack.Peek(0).Sym)
	case 48: // ZeroOrOne
		{
			var zero Token
			return zero
		}
	default:
		panic("unreachable")
	}
}

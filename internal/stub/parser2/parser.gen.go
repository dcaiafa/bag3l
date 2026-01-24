package parser

import (
	_i0 "github.com/dcaiafa/bag3l/internal/stub/ast"
	_i1 "github.com/dcaiafa/bag3l/internal/token"
)

var _rules = []int32{
	0, 1, 1, 2, 3, 4, 4, 4, 4, 4, 5, 6, 7, 8,
	9, 10, 11, 12, 12, 13, 14, 15, 16, 16, 17, 18, 18, 19,
	19, 20, 20, 21, 21, 22, 22, 23, 23, 24, 24, 25, 25, 26,
	26, 27, 27, 28, 28, 29, 29, 30, 30,
}

var _termCounts = []int32{
	1, 1, 1, 2, 2, 1, 1, 1, 1, 1, 5, 1, 3, 5,
	2, 3, 2, 3, 1, 6, 4, 2, 3, 1, 2, 1, 1, 1,
	0, 2, 1, 1, 0, 2, 1, 1, 0, 1, 0, 3, 1, 1,
	0, 1, 0, 1, 0, 3, 1, 1, 0,
}

var _actions = []int32{
	76, 81, 84, 89, 102, 105, 108, 111, 114, 117, 120, 123, 136, 149,
	152, 165, 178, 191, 204, 217, 236, 255, 268, 84, 271, 274, 283, 286,
	299, 324, 327, 340, 345, 352, 359, 372, 84, 268, 377, 382, 385, 390,
	403, 416, 419, 424, 427, 432, 437, 460, 483, 496, 513, 526, 531, 534,
	551, 268, 554, 559, 562, 268, 575, 588, 601, 614, 621, 84, 626, 631,
	636, 641, 646, 268, 651, 664, 4, 1, 1, 4, 2, 2, 0, -2,
	4, 19, 19, 18, 20, 12, 2, 6, 0, -28, 3, 7, 8, 8,
	6, 9, 5, 10, 2, 0, 2147483647, 2, 0, -1, 2, 20, 22, 2,
	20, 26, 2, 20, 23, 2, 20, 24, 2, 20, 25, 12, 2, -5,
	0, -5, 3, -5, 8, -5, 6, -5, 5, -5, 12, 2, -30, 0,
	-30, 3, -30, 8, -30, 6, -30, 5, -30, 2, 0, -3, 12, 2,
	6, 0, -27, 3, 7, 8, 8, 6, 9, 5, 10, 12, 2, -9,
	0, -9, 3, -9, 8, -9, 6, -9, 5, -9, 12, 2, -6, 0,
	-6, 3, -6, 8, -6, 6, -6, 5, -6, 12, 2, -8, 0, -8,
	3, -8, 8, -8, 6, -8, 5, -8, 12, 2, -7, 0, -7, 3,
	-7, 8, -7, 6, -7, 5, -7, 18, 16, -26, 2, -26, 15, -26,
	13, -26, 0, -26, 3, -26, 8, -26, 6, -26, 5, -26, 18, 16,
	-25, 2, -25, 15, -25, 13, -25, 0, -25, 3, -25, 8, -25, 6,
	-25, 5, -25, 12, 2, -4, 0, -4, 3, -4, 8, -4, 6, -4,
	5, -4, 2, 20, 28, 2, 10, 31, 8, 20, -36, 19, -36, 12,
	32, 18, -36, 2, 14, 35, 12, 2, -29, 0, -29, 3, -29, 8,
	-29, 6, -29, 5, -29, 24, 11, -50, 16, -50, 2, -50, 15, -50,
	0, -50, 9, -50, 3, -50, 20, -50, 8, -50, 17, 48, 6, -50,
	5, -50, 2, 9, 36, 12, 2, -12, 0, -12, 3, -12, 8, -12,
	6, -12, 5, -12, 4, 11, -32, 20, 37, 6, 20, -35, 19, -35,
	18, -35, 6, 20, 41, 19, 19, 18, 20, 12, 2, -15, 0, -15,
	3, -15, 8, -15, 6, -15, 5, -15, 4, 15, -38, 20, 46, 4,
	11, -34, 20, -34, 2, 11, 52, 4, 11, -31, 20, 37, 12, 2,
	-18, 0, -18, 3, -18, 8, -18, 6, -18, 5, -18, 12, 2, -16,
	0, -16, 3, -16, 8, -16, 6, -16, 5, -16, 2, 13, 54, 4,
	16, 59, 15, -37, 2, 15, 55, 4, 21, 56, 20, -44, 4, 16,
	-40, 15, -40, 22, 11, -49, 16, -49, 2, -49, 15, -49, 0, -49,
	9, -49, 3, -49, 20, -49, 8, -49, 6, -49, 5, -49, 22, 11,
	-24, 16, -24, 2, -24, 15, -24, 0, -24, 9, -24, 3, -24, 20,
	-24, 8, -24, 6, -24, 5, -24, 12, 2, -10, 0, -10, 3, -10,
	8, -10, 6, -10, 5, -10, 16, 16, -11, 2, -11, 15, -11, 0,
	-11, 3, -11, 8, -11, 6, -11, 5, -11, 12, 2, -13, 0, -13,
	3, -13, 8, -13, 6, -13, 5, -13, 4, 11, -14, 20, -14, 2,
	20, 60, 16, 2, -42, 0, -42, 3, -42, 20, 28, 8, -42, 14,
	61, 6, -42, 5, -42, 2, 20, -43, 4, 11, -33, 20, -33, 2,
	20, 46, 12, 2, -17, 0, -17, 3, -17, 8, -17, 6, -17, 5,
	-17, 12, 2, -41, 0, -41, 3, -41, 8, -41, 6, -41, 5, -41,
	12, 2, -19, 0, -19, 3, -19, 8, -19, 6, -19, 5, -19, 12,
	2, -23, 0, -23, 3, -23, 8, -23, 6, -23, 5, -23, 6, 16,
	-46, 15, -46, 9, 67, 4, 16, -39, 15, -39, 4, 16, -45, 15,
	-45, 4, 16, -20, 15, -20, 4, 16, 73, 15, 74, 4, 16, -48,
	15, -48, 4, 16, -21, 15, -21, 12, 2, -22, 0, -22, 3, -22,
	8, -22, 6, -22, 5, -22, 4, 16, -47, 15, -47,
}

var _goto = []int32{
	76, 83, 84, 87, 83, 83, 83, 83, 83, 83, 83, 83, 83, 83,
	104, 83, 83, 83, 83, 83, 83, 83, 117, 120, 83, 123, 83, 83,
	128, 83, 83, 131, 83, 138, 83, 143, 150, 155, 83, 83, 158, 83,
	83, 83, 83, 83, 161, 83, 83, 83, 83, 83, 83, 83, 83, 164,
	83, 171, 83, 174, 83, 177, 83, 83, 83, 182, 83, 187, 83, 83,
	83, 83, 83, 192, 83, 83, 6, 3, 3, 1, 4, 2, 5, 0,
	2, 18, 21, 16, 5, 11, 4, 12, 19, 13, 20, 14, 13, 15,
	7, 16, 8, 17, 10, 18, 12, 5, 11, 4, 27, 13, 15, 7,
	16, 8, 17, 10, 18, 2, 17, 29, 2, 18, 30, 4, 23, 33,
	11, 34, 2, 30, 49, 6, 9, 38, 21, 39, 22, 40, 4, 12,
	42, 18, 43, 6, 25, 44, 24, 45, 14, 47, 4, 6, 50, 18,
	51, 2, 17, 53, 2, 9, 58, 2, 27, 57, 6, 16, 62, 26,
	63, 17, 64, 2, 17, 65, 2, 14, 66, 4, 29, 70, 17, 71,
	4, 15, 68, 28, 69, 4, 6, 72, 18, 51, 2, 17, 75,
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
			_cast[_i1.Token](p._stack.Peek(0).Sym),
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
			_cast[_i1.Token](p._stack.Peek(0).Sym),
		)
	case 12:
		return p.on_import_decl(
			_cast[Token](p._stack.Peek(2).Sym),
			_cast[Token](p._stack.Peek(1).Sym),
			_cast[_i1.Token](p._stack.Peek(0).Sym),
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
			_cast[_i1.Token](p._stack.Peek(2).Sym),
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
	case 25:
		return p.on_string(
			_cast[Token](p._stack.Peek(0).Sym),
		)
	case 26:
		return p.on_string(
			_cast[Token](p._stack.Peek(0).Sym),
		)
	case 27: // ZeroOrMore
		return _cast[[]_i0.AST](p._stack.Peek(0).Sym)
	case 28: // ZeroOrMore
		{
			var zero []_i0.AST
			return zero
		}
	case 29: // OneOrMore
		return append(
			_cast[[]_i0.AST](p._stack.Peek(1).Sym),
			_cast[_i0.AST](p._stack.Peek(0).Sym),
		)
	case 30: // OneOrMore
		return []_i0.AST{
			_cast[_i0.AST](p._stack.Peek(0).Sym),
		}
	case 31: // ZeroOrMore
		return _cast[[]*_i0.StructField](p._stack.Peek(0).Sym)
	case 32: // ZeroOrMore
		{
			var zero []*_i0.StructField
			return zero
		}
	case 33: // OneOrMore
		return append(
			_cast[[]*_i0.StructField](p._stack.Peek(1).Sym),
			_cast[*_i0.StructField](p._stack.Peek(0).Sym),
		)
	case 34: // OneOrMore
		return []*_i0.StructField{
			_cast[*_i0.StructField](p._stack.Peek(0).Sym),
		}
	case 35: // ZeroOrOne
		return _cast[Token](p._stack.Peek(0).Sym)
	case 36: // ZeroOrOne
		{
			var zero Token
			return zero
		}
	case 37: // ZeroOrOne
		return _cast[[]*_i0.FuncParam](p._stack.Peek(0).Sym)
	case 38: // ZeroOrOne
		{
			var zero []*_i0.FuncParam
			return zero
		}
	case 39: // List
		return append(
			_cast[[]*_i0.FuncParam](p._stack.Peek(2).Sym),
			_cast[*_i0.FuncParam](p._stack.Peek(0).Sym),
		)
	case 40: // List
		return []*_i0.FuncParam{
			_cast[*_i0.FuncParam](p._stack.Peek(0).Sym),
		}
	case 41: // ZeroOrOne
		return _cast[[]*_i0.TypeRef](p._stack.Peek(0).Sym)
	case 42: // ZeroOrOne
		{
			var zero []*_i0.TypeRef
			return zero
		}
	case 43: // ZeroOrOne
		return _cast[Token](p._stack.Peek(0).Sym)
	case 44: // ZeroOrOne
		{
			var zero Token
			return zero
		}
	case 45: // ZeroOrOne
		return _cast[*_i0.FuncParam](p._stack.Peek(0).Sym)
	case 46: // ZeroOrOne
		{
			var zero *_i0.FuncParam
			return zero
		}
	case 47: // List
		return append(
			_cast[[]*_i0.TypeRef](p._stack.Peek(2).Sym),
			_cast[*_i0.TypeRef](p._stack.Peek(0).Sym),
		)
	case 48: // List
		return []*_i0.TypeRef{
			_cast[*_i0.TypeRef](p._stack.Peek(0).Sym),
		}
	case 49: // ZeroOrOne
		return _cast[Token](p._stack.Peek(0).Sym)
	case 50: // ZeroOrOne
		{
			var zero Token
			return zero
		}
	default:
		panic("unreachable")
	}
}

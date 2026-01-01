package parser

const (
	EOF      int = 0
	ERROR    int = 1
	CONST    int = 2
	FUNC     int = 3
	PACKAGE  int = 4
	TYPE     int = 5
	STRUCT   int = 6
	NIL      int = 7
	IMPORT   int = 8
	EQ       int = 9
	OCURLY   int = 10
	CCURLY   int = 11
	STAR     int = 12
	DOT      int = 13
	OPAREN   int = 14
	CPAREN   int = 15
	COMMA    int = 16
	QMARK    int = 17
	STRING   int = 18
	ID       int = 19
	ELLIPSIS int = 20
)

func _TokenToString(t int) string {
	switch t {
	case EOF:
		return "EOF"
	case ERROR:
		return "ERROR"
	case CONST:
		return "CONST"
	case FUNC:
		return "FUNC"
	case PACKAGE:
		return "PACKAGE"
	case TYPE:
		return "TYPE"
	case STRUCT:
		return "STRUCT"
	case NIL:
		return "NIL"
	case IMPORT:
		return "IMPORT"
	case EQ:
		return "EQ"
	case OCURLY:
		return "OCURLY"
	case CCURLY:
		return "CCURLY"
	case STAR:
		return "STAR"
	case DOT:
		return "DOT"
	case OPAREN:
		return "OPAREN"
	case CPAREN:
		return "CPAREN"
	case COMMA:
		return "COMMA"
	case QMARK:
		return "QMARK"
	case STRING:
		return "STRING"
	case ID:
		return "ID"
	case ELLIPSIS:
		return "ELLIPSIS"
	default:
		return "???"
	}
}

type _Stack[T any] []T

func (s *_Stack[T]) Push(x T) {
	*s = append(*s, x)
}

func (s *_Stack[T]) Pop(n int) {
	*s = (*s)[:len(*s)-n]
}

func (s _Stack[T]) Peek(n int) T {
	return s[len(s)-n-1]
}

func (s _Stack[T]) PeekSlice(n int) []T {
	return s[len(s)-n:]
}

package parser

const (
	EOF                 int = 0
	ERROR               int = 1
	INFO                int = 2
	PARAM               int = 3
	FLAG                int = 4
	AND                 int = 5
	BREAK               int = 6
	CATCH               int = 7
	CONTINUE            int = 8
	DEFER               int = 9
	ELSE                int = 10
	FALSE               int = 11
	FOR                 int = 12
	FUNC                int = 13
	IF                  int = 14
	IMPORT              int = 15
	NIL                 int = 16
	NOT                 int = 17
	OR                  int = 18
	RETURN              int = 19
	THROW               int = 20
	TRUE                int = 21
	TRY                 int = 22
	VAR                 int = 23
	WHILE               int = 24
	YIELD               int = 25
	ASSIGN              int = 26
	ASSIGN_ADD          int = 27
	ASSIGN_SUB          int = 28
	ASSIGN_MUL          int = 29
	ASSIGN_DIV          int = 30
	ASSIGN_MOD          int = 31
	EQ                  int = 32
	NE                  int = 33
	LT                  int = 34
	LE                  int = 35
	GT                  int = 36
	GE                  int = 37
	ADD                 int = 38
	SUB                 int = 39
	MUL                 int = 40
	DIV                 int = 41
	INC                 int = 42
	DEC                 int = 43
	QUESTION_MARK       int = 44
	SEMICOLON           int = 45
	COMMA               int = 46
	COLON               int = 47
	PERIOD              int = 48
	OPAREN              int = 49
	CPAREN              int = 50
	OBRACKET            int = 51
	CBRACKET            int = 52
	OCURLY              int = 53
	CCURLY              int = 54
	ARROW               int = 55
	LAMBDA              int = 56
	PIPE                int = 57
	EXPAND              int = 58
	NUMBER              int = 59
	ID                  int = 60
	REGEX               int = 61
	NEWLINE             int = 62
	CHAR                int = 63
	SHEBANG             int = 64
	FORMAT              int = 65
	STRING              int = 66
	RAW_STRING          int = 67
	FSTR_OPEN           int = 68
	FBSTR_LITERAL       int = 69
	FBSTR_ESC_LITERAL   int = 70
	FSTR_CLOSE          int = 71
	EXEC_OPEN           int = 72
	EXEC_WS             int = 73
	EXEC_HOME           int = 74
	EXEC_LITERAL        int = 75
	EXEC_CLOSE          int = 76
	EXEC_DQUOTE_LITERAL int = 77
	EXEC_SQUOTE_LITERAL int = 78
)

func _TokenToString(t int) string {
	switch t {
	case EOF:
		return "EOF"
	case ERROR:
		return "ERROR"
	case INFO:
		return "INFO"
	case PARAM:
		return "PARAM"
	case FLAG:
		return "FLAG"
	case AND:
		return "AND"
	case BREAK:
		return "BREAK"
	case CATCH:
		return "CATCH"
	case CONTINUE:
		return "CONTINUE"
	case DEFER:
		return "DEFER"
	case ELSE:
		return "ELSE"
	case FALSE:
		return "FALSE"
	case FOR:
		return "FOR"
	case FUNC:
		return "FUNC"
	case IF:
		return "IF"
	case IMPORT:
		return "IMPORT"
	case NIL:
		return "NIL"
	case NOT:
		return "NOT"
	case OR:
		return "OR"
	case RETURN:
		return "RETURN"
	case THROW:
		return "THROW"
	case TRUE:
		return "TRUE"
	case TRY:
		return "TRY"
	case VAR:
		return "VAR"
	case WHILE:
		return "WHILE"
	case YIELD:
		return "YIELD"
	case ASSIGN:
		return "ASSIGN"
	case ASSIGN_ADD:
		return "ASSIGN_ADD"
	case ASSIGN_SUB:
		return "ASSIGN_SUB"
	case ASSIGN_MUL:
		return "ASSIGN_MUL"
	case ASSIGN_DIV:
		return "ASSIGN_DIV"
	case ASSIGN_MOD:
		return "ASSIGN_MOD"
	case EQ:
		return "EQ"
	case NE:
		return "NE"
	case LT:
		return "LT"
	case LE:
		return "LE"
	case GT:
		return "GT"
	case GE:
		return "GE"
	case ADD:
		return "ADD"
	case SUB:
		return "SUB"
	case MUL:
		return "MUL"
	case DIV:
		return "DIV"
	case INC:
		return "INC"
	case DEC:
		return "DEC"
	case QUESTION_MARK:
		return "QUESTION_MARK"
	case SEMICOLON:
		return "SEMICOLON"
	case COMMA:
		return "COMMA"
	case COLON:
		return "COLON"
	case PERIOD:
		return "PERIOD"
	case OPAREN:
		return "OPAREN"
	case CPAREN:
		return "CPAREN"
	case OBRACKET:
		return "OBRACKET"
	case CBRACKET:
		return "CBRACKET"
	case OCURLY:
		return "OCURLY"
	case CCURLY:
		return "CCURLY"
	case ARROW:
		return "ARROW"
	case LAMBDA:
		return "LAMBDA"
	case PIPE:
		return "PIPE"
	case EXPAND:
		return "EXPAND"
	case NUMBER:
		return "NUMBER"
	case ID:
		return "ID"
	case REGEX:
		return "REGEX"
	case NEWLINE:
		return "NEWLINE"
	case CHAR:
		return "CHAR"
	case SHEBANG:
		return "SHEBANG"
	case FORMAT:
		return "FORMAT"
	case STRING:
		return "STRING"
	case RAW_STRING:
		return "RAW_STRING"
	case FSTR_OPEN:
		return "FSTR_OPEN"
	case FBSTR_LITERAL:
		return "FBSTR_LITERAL"
	case FBSTR_ESC_LITERAL:
		return "FBSTR_ESC_LITERAL"
	case FSTR_CLOSE:
		return "FSTR_CLOSE"
	case EXEC_OPEN:
		return "EXEC_OPEN"
	case EXEC_WS:
		return "EXEC_WS"
	case EXEC_HOME:
		return "EXEC_HOME"
	case EXEC_LITERAL:
		return "EXEC_LITERAL"
	case EXEC_CLOSE:
		return "EXEC_CLOSE"
	case EXEC_DQUOTE_LITERAL:
		return "EXEC_DQUOTE_LITERAL"
	case EXEC_SQUOTE_LITERAL:
		return "EXEC_SQUOTE_LITERAL"
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

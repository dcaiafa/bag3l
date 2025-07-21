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
	MOD                 int = 42
	INC                 int = 43
	DEC                 int = 44
	QUESTION_MARK       int = 45
	SEMICOLON           int = 46
	COMMA               int = 47
	COLON               int = 48
	PERIOD              int = 49
	OPAREN              int = 50
	CPAREN              int = 51
	OBRACKET            int = 52
	CBRACKET            int = 53
	OCURLY              int = 54
	CCURLY              int = 55
	ARROW               int = 56
	LAMBDA              int = 57
	PIPE                int = 58
	EXPAND              int = 59
	BACKTICK            int = 60
	NUMBER              int = 61
	ID                  int = 62
	REGEX               int = 63
	NEWLINE             int = 64
	CHAR                int = 65
	STRING              int = 66
	EXEC_PREFIX         int = 67
	EXEC_WS             int = 68
	EXEC_HOME           int = 69
	EXEC_LITERAL        int = 70
	EXEC_DQUOTE_LITERAL int = 71
	EXEC_SQUOTE_LITERAL int = 72
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
	case MOD:
		return "MOD"
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
	case BACKTICK:
		return "BACKTICK"
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
	case STRING:
		return "STRING"
	case EXEC_PREFIX:
		return "EXEC_PREFIX"
	case EXEC_WS:
		return "EXEC_WS"
	case EXEC_HOME:
		return "EXEC_HOME"
	case EXEC_LITERAL:
		return "EXEC_LITERAL"
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

package parser

import (
	gotoken "go/token"

	"github.com/dcaiafa/loxlex/simplelexer"
)

type lexer struct {
	lexer *simplelexer.Lexer

	stashed  Token
	eolToken bool
}

func newLexer(file *gotoken.File, filename string, input []byte) *lexer {
	return &lexer{
		lexer: simplelexer.New(simplelexer.Config{
			StateMachine: &_LexerStateMachine{},
			File:         file,
			Input:        input,
		}),
	}
}

func (l *lexer) ReadToken() (Token, int) {
	var t Token

	for {
		if l.stashed.Type != 0 {
			t = l.stashed
			l.stashed.Type = 0
		} else {
			t, _ = l.lexer.ReadToken()
		}

		switch t.Type {
		case NEWLINE, EOF:
			if l.eolToken {
				l.eolToken = false
				l.stashed = t
				t = Token{
					Type: SEMICOLON,
					Pos:  l.stashed.Pos,
				}
			}

		case CPAREN,
			CBRACKET,
			CCURLY,
			EXPAND,
			NUMBER,
			STRING,
			ID,
			TRUE,
			FALSE,
			BREAK,
			CONTINUE,
			RETURN,
			REGEX,
			CHAR,
			INC,
			DEC,
			NIL:
			l.eolToken = true

		default:
			l.eolToken = false
		}

		if t.Type != NEWLINE {
			break
		}
	}

	return t, t.Type
}

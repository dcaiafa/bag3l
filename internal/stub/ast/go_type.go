package ast

import "github.com/dcaiafa/bag3l/internal/token"

type GoType struct {
	Ref     bool
	Package token.Token
	ID      token.Token
}

func (t *GoType) RunPass(ctx *Context, pass Pass) {}

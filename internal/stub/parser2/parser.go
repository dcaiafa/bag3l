package parser

import (
	"fmt"
	gotoken "go/token"
	"regexp"

	"github.com/dcaiafa/bag3l/internal/errlogger"
	"github.com/dcaiafa/bag3l/internal/stub/ast"
	"github.com/dcaiafa/bag3l/internal/token"
	"github.com/dcaiafa/loxlex/simplelexer"
)

func Parse(
	filename string,
	input []byte,
	errLogger errlogger.ErrLogger,
) (*ast.Unit, error) {
	file := gotoken.NewFileSet().AddFile(filename, -1, len(input))
	parser := &parser{
		file:      file,
		errLogger: errlogger.NewErrLoggerBase(errLogger),
	}

	lexer := simplelexer.New(simplelexer.Config{
		StateMachine: new(_LexerStateMachine),
		File:         file,
		Input:        input,
	})

	ok := parser.parse(lexer)
	if parser.errLogger.Error() != nil {
		return nil, parser.errLogger.Error()
	}
	if !ok {
		return nil, fmt.Errorf("parsing failed")
	}
	return parser.unit, nil
}

type Token = simplelexer.Token

type parser struct {
	lox

	file      *gotoken.File
	errLogger *errlogger.ErrLoggerWrapper
	unit      *ast.Unit
}

func (p *parser) on_s(u *ast.Unit) *ast.Unit {
	p.unit = u
	return u
}

func (p *parser) on_s__error(err Error) *ast.Unit {
	if err.Token.Type == ERROR {
		p.errLogger.Failf(p.tokenPos(err.Token), "%v", err.Token.Err)
	} else {
		p.errLogger.Failf(
			p.tokenPos(err.Token), "Unexpected %v", _TokenToString(err.Token.Type))
	}
	return &ast.Unit{}
}

func (p *parser) on_unit(u *ast.Unit, decls []ast.AST) *ast.Unit {
	u.Decls = decls
	return u
}

func (p *parser) on_package(_ Token, n token.Token) *ast.Unit {
	return &ast.Unit{
		Package: n,
	}
}

func (p *parser) on_decl(decl ast.AST) ast.AST {
	return decl
}

func (p *parser) on_const_decl(_ Token, id Token, t *ast.TypeRef, _ Token, v *ast.ConstValue) *ast.ConstDecl {
	return &ast.ConstDecl{
		ID:      p.tokenToNitro(id),
		TypeRef: t,
		Value:   v,
	}
}

func (p *parser) on_const_value(v token.Token) *ast.ConstValue {
	return &ast.ConstValue{Expr: v.Str}
}

func (p *parser) on_import_decl(_ Token, a Token, i token.Token) *ast.ImportDecl {
	return &ast.ImportDecl{
		Alias:  p.tokenToNitro(a).Str,
		Import: i.Str,
	}
}

func (p *parser) on_struct_decl(_, n, _ Token, fields []*ast.StructField, _ Token) *ast.StructDecl {
	return &ast.StructDecl{
		Name:   p.tokenToNitro(n).Str,
		Fields: fields,
	}
}

func (p *parser) on_struct_field(n Token, t *ast.TypeRef) *ast.StructField {
	return &ast.StructField{
		Name: p.tokenToNitro(n).Str,
		Type: t,
	}
}

func (p *parser) on_type_decl(_ Token, n Token, gt *ast.GoType) *ast.TypeDecl {
	return &ast.TypeDecl{
		ID:     p.tokenToNitro(n),
		GoType: gt,
	}
}

func (p *parser) on_go_type(ref Token, gt *ast.GoType) *ast.GoType {
	gt.Ref = ref.Type == STAR
	return gt
}

func (p *parser) on_simple_go_type__full(pkg token.Token, _ Token, id Token) *ast.GoType {
	return &ast.GoType{
		Package: pkg,
		ID:      p.tokenToNitro(id),
	}
}

func (p *parser) on_simple_go_type(id Token) *ast.GoType {
	return &ast.GoType{
		ID: p.tokenToNitro(id),
	}
}

func (p *parser) on_func_decl(_, n, _ Token, ps []*ast.FuncParam, _ Token, rets []*ast.TypeRef) *ast.FuncDecl {
	return &ast.FuncDecl{
		ID:     p.tokenToNitro(n),
		Params: ps,
		Rets:   rets,
	}
}

func (p *parser) on_func_param(n Token, vararg Token, t *ast.TypeRef, def *ast.FuncParam) *ast.FuncParam {
	fp := def
	if fp == nil {
		fp = &ast.FuncParam{}
	}

	fp.ID = p.tokenToNitro(n)
	fp.Type = t
	fp.VarArg = vararg.Type == ELLIPSIS

	return fp
}

func (p *parser) on_func_param_default(_ Token, v *ast.ConstValue) *ast.FuncParam {
	return &ast.FuncParam{
		DefaultValue: v,
	}
}

func (p *parser) on_func_rets__multi(_ Token, rets []*ast.TypeRef, _ Token) []*ast.TypeRef {
	return rets
}

func (p *parser) on_func_rets__single(ret *ast.TypeRef) []*ast.TypeRef {
	return []*ast.TypeRef{ret}
}

func (p *parser) on_type_ref(n Token, nilable Token) *ast.TypeRef {
	return &ast.TypeRef{
		ID:      p.tokenToNitro(n),
		Nilable: nilable.Type == QMARK,
	}
}

func (p *parser) tokenToNitro(at Token) token.Token {
	var err error

	t := token.Token{}
	switch at.Type {
	case STRING:
		s := string(at.Str)
		s = s[1 : len(s)-1] // remove quotes
		t.Type = token.String
		t.Str = s
		t.Str, err = expandEscapeSequences(s)
		if err != nil {
			p.errLogger.Failf(p.tokenPos(at), "Invalid string literal: %v", err)
		}

	case RSTRING:
		s := string(at.Str)
		s = s[1 : len(s)-1] // remove quotes
		t.Type = token.String
		t.Str = s

	default:
		t.Type = token.String
		t.Str = string(at.Str)
	}

	t.Pos = p.tokenPos(at)
	return t
}

func (p *parser) on_string(n Token) token.Token {
	return p.tokenToNitro(n)
}

var escapeSeqRegex = regexp.MustCompile(
	`\\x[A-Fa-f0-9]{2}|\\u[A-Fa-f0-9]{2}|\\U[A-Fa-f0-9]{4}|\\.`)

func expandEscapeSequences(s string) (string, error) {
	var err error
	s = escapeSeqRegex.ReplaceAllStringFunc(s, func(s string) string {
		switch s[1] {
		case 'n':
			return "\n"
		case 'r':
			return "\r"
		case 't':
			return "\t"
		case '"':
			return "\""
		case '\'':
			return "'"
		case '\\':
			return "\\"
		default:
			if err == nil {
				err = fmt.Errorf("invalid escape sequence %s", s)
			}
			return ""
		}
	})

	if err != nil {
		return "", err
	}

	return s, nil
}

func (p *parser) tokenPos(t Token) token.Pos {
	position := p.file.Position(t.Pos)
	return token.Pos{
		Filename: position.Filename,
		Line:     position.Line,
		Col:      position.Column,
	}
}

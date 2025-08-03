package parser

import (
	"fmt"
	gotoken "go/token"
	"strconv"
	"strings"

	"github.com/dcaiafa/bag3l/internal/ast"
	"github.com/dcaiafa/bag3l/internal/errlogger"
	"github.com/dcaiafa/bag3l/internal/token"
	"github.com/dcaiafa/bag3l/internal/vm"
	"github.com/dcaiafa/loxlex/simplelexer"
)

func ParseUnit(filename string, input []byte, errLogger errlogger.ErrLogger) (*ast.Unit, error) {
	file := gotoken.NewFileSet().AddFile(filename, -1, len(input))
	parser := &parser{
		file:      file,
		errLogger: errlogger.NewErrLoggerBase(errLogger),
	}

	lexer := newLexer(file, filename, input)
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

func (p *parser) on_start(unit *ast.Unit) any {
	p.unit = unit
	return nil
}

func (p *parser) on_start__error(e Error) any {
	if e.Token.Type == ERROR {
		p.errLogger.Failf(p.tokenPos(e.Token), "%v", e.Token.Err)
		return nil
	} else {
		p.errLogger.Failf(
			p.tokenPos(e.Token), "unexpected %v", _TokenToString(e.Token.Type))
	}
	return nil
}

func (p *parser) on_unit(_ Token, meta []ast.AST, imports []*ast.Import, stmts *ast.StmtBlock) *ast.Unit {
	return &ast.Unit{
		Meta:    meta,
		Imports: imports,
		Block:   stmts.Stmts,
	}
}

func (p *parser) on_meta_directive(directive ast.AST) ast.AST {
	return directive
}

func (p *parser) on_meta_info(info Token, attribs []ast.AST, _ Token) ast.AST {
	p.errLogger.Failf(p.tokenPos(info), "Info metadata not implemented")
	return nil
}

func (p *parser) on_meta_param(_ Token, id Token, value ast.Expr, attribs []ast.AST, _ Token) ast.AST {
	return &ast.MetaParam{
		Name:    string(id.Str),
		IsFlag:  false,
		Default: value,
		Attribs: attribs,
	}
}

func (p *parser) on_meta_flag(_ Token, id Token, value ast.Expr, attribs []ast.AST, _ Token) ast.AST {
	return &ast.MetaParam{
		Name:    string(id.Str),
		IsFlag:  true,
		Default: value,
		Attribs: attribs,
	}
}

func (p *parser) on_meta_value(_ Token, value ast.Expr) ast.Expr {
	return value
}

func (p *parser) on_meta_attribs(_ Token, attribs []ast.AST, _ Token) []ast.AST {
	return attribs
}

func (p *parser) on_meta_attrib_list(attribs []ast.AST, _ Token) []ast.AST {
	return attribs
}

func (p *parser) on_meta_attrib(id string, value vm.Value) ast.AST {
	return &ast.MetaAttrib{
		Name:  id,
		Value: value,
	}
}

func (p *parser) on_meta_attrib_value(_ Token, value vm.Value) vm.Value {
	return value
}

func (p *parser) on_meta_literal(lit Token) vm.Value {
	t := p.tokenToNitro(lit)
	switch t.Type {
	case token.Nil:
		return nil
	case token.Int:
		return vm.NewInt(t.Int)
	case token.Float:
		return vm.NewFloat(t.Float)
	case token.String:
		return vm.NewString(t.Str)
	case token.Bool:
		return vm.NewBool(t.Bool)
	default:
		panic("unreachable")
	}
}

func (p *parser) on_import_stmt(_ Token, alias Token, pkg Token, _ Token) *ast.Import {
	return &ast.Import{
		Alias:   string(alias.Str),
		Package: p.tokenToNitro(pkg).Str,
	}
}

func (p *parser) on_stmts(stmts ast.ASTs) *ast.StmtBlock {
	return &ast.StmtBlock{
		Stmts: stmts,
	}
}

func (p *parser) on_stmt_list(stmts []ast.AST, _ Token) ast.ASTs {
	return ast.ASTs(stmts)
}

func (p *parser) on_stmt(stmt ast.AST) ast.AST {
	return stmt
}

func (p *parser) on_expr_stmt(expr ast.Expr) ast.AST {
	return &ast.ExprStmt{Expr: expr}
}

func (p *parser) on_assignment_stmt(lvalueExprs []ast.Expr, _ Token, rvalues []ast.Expr) ast.AST {
	lvalues := make(ast.ASTs, len(lvalueExprs))
	for i, lvalueExpr := range lvalueExprs {
		lvalues[i] = p.lvalue(lvalueExpr)
	}
	return &ast.AssignStmt{
		Lvalues: lvalues,
		Rvalues: ast.Exprs(rvalues),
	}
}

func (p *parser) on_assignment_op_stmt(lvalueExpr ast.Expr, op Token, rvalue ast.Expr) ast.AST {
	opAssign := &ast.AssignOpStmt{
		LValue: p.lvalue(lvalueExpr),
		RValue: rvalue,
	}

	switch op.Type {
	case ASSIGN_ADD:
		opAssign.Op = ast.BinOpPlus
	case ASSIGN_SUB:
		opAssign.Op = ast.BinOpMinus
	case ASSIGN_MUL:
		opAssign.Op = ast.BinOpMult
	case ASSIGN_DIV:
		opAssign.Op = ast.BinOpDiv
	default:
		panic("unreachable")
	}

	return opAssign
}

func (p *parser) on_assign_op(op Token) Token {
	return op
}

func (p *parser) on_var_decl_stmt(_ Token, vars []Token, initValues []ast.Expr) ast.AST {
	return &ast.VarDeclStmt{
		Vars:       p.tokens(vars),
		InitValues: initValues,
	}
}

func (p *parser) on_var_decl_init(_ Token, initValues []ast.Expr) []ast.Expr {
	return initValues
}

func (p *parser) on_for_stmt(_ Token, vars []Token, _ Token, iterExpr ast.Expr, _ Token, block *ast.StmtBlock, _ Token) ast.AST {
	forVars := make([]*ast.ForVar, len(vars))
	for i, v := range vars {
		name := p.tokenToNitro(v)
		forVars[i] = &ast.ForVar{
			VarName: name,
		}
		forVars[i].SetPos(name.Pos)
	}
	return &ast.ForStmt{
		ForVars:  forVars,
		IterExpr: iterExpr,
		Block:    block,
	}
}

func (p *parser) on_while_stmt(_ Token, pred ast.Expr, _ Token, block *ast.StmtBlock, _ Token) ast.AST {
	return &ast.WhileStmt{
		Predicate: pred,
		Block:     block,
	}
}

func (p *parser) on_if_stmt(ifkw Token, pred ast.Expr, _ Token, block *ast.StmtBlock, _ Token, elseSections []*ast.IfSection) ast.AST {
	ifStmt := &ast.IfStmt{}

	ifSection := &ast.IfSection{
		Pred:  pred,
		Block: block,
	}
	ifSection.SetPos(p.tokenPos(ifkw))

	for i, elseSection := range elseSections {
		if elseSection.Pred == nil && i != len(elseSections)-1 {
			p.errLogger.Failf(
				elseSections[i+1].Pos(),
				"Unexpected else clause")
		}
	}

	ifStmt.Sections = append([]*ast.IfSection{ifSection}, elseSections...)

	return ifStmt
}

func (p *parser) on_if_else__if(_, _ Token, pred ast.Expr, _ Token, block *ast.StmtBlock, _ Token) *ast.IfSection {
	return &ast.IfSection{
		Pred:  pred,
		Block: block,
	}
}

func (p *parser) on_if_else(_, _ Token, block *ast.StmtBlock, _ Token) *ast.IfSection {
	return &ast.IfSection{
		Block: block,
	}
}

func (p *parser) on_func_stmt(_ Token, name Token, _ Token, params []*ast.FuncParam, _, _ Token, block *ast.StmtBlock, _ Token) ast.AST {
	f := &ast.FuncStmt{}
	f.Name = string(name.Str)
	f.Params = params
	f.Block = block
	return f
}

func (p *parser) on_param_list(params []Token) []*ast.FuncParam {
	fparams := make([]*ast.FuncParam, len(params))
	for i, param := range params {
		fparams[i] = &ast.FuncParam{
			Name: string(param.Str),
		}
	}
	return fparams
}

func (p *parser) on_return_stmt(_ Token, values []ast.Expr) ast.AST {
	return &ast.ReturnStmt{
		Values: values,
	}
}

func (p *parser) on_try_catch_stmt(_, _ Token, tryBlock *ast.StmtBlock, _, _ Token, catchVarToken Token, _ Token, catchBlock *ast.StmtBlock, _ Token) ast.AST {
	var catchVar *token.Token
	if catchVarToken.Type != 0 {
		catchVar = new(token.Token)
		*catchVar = p.tokenToNitro(catchVarToken)
	}
	return ast.NewTryCatchStmt(tryBlock, catchVar, catchBlock)
}

func (p *parser) on_throw_stmt(_ Token, expr ast.Expr) ast.AST {
	return &ast.ThrowStmt{
		Expr: expr,
	}
}

func (p *parser) on_defer_stmt(_ Token, expr ast.Expr) ast.AST {
	callExpr, ok := expr.(*ast.FuncCallExpr)
	if !ok {
		p.errLogger.Failf(expr.Pos(), "Deferred expression must be a function call")
		return nil
	}
	return ast.NewDeferStmt(callExpr)
}

func (p *parser) on_yield_stmt(_ Token, values []ast.Expr) ast.AST {
	return &ast.YieldStmt{
		Values: values,
	}
}

func (p *parser) on_break_stmt(_ Token) ast.AST {
	return &ast.BreakStmt{}
}

func (p *parser) on_continue_stmt(_ Token) ast.AST {
	return &ast.ContinueStmt{}
}

func (p *parser) on_inc_dec_stmt(lvalue ast.Expr, op Token) ast.AST {
	var opType ast.IncDecOp
	switch op.Type {
	case INC:
		opType = ast.IncDecOpInc
	case DEC:
		opType = ast.IncDecOpDec
	default:
		panic("unreachable")
	}
	return &ast.IncDecStmt{
		LValue: p.lvalue(lvalue),
		Op:     opType,
	}
}

func (p *parser) on_expr(e ast.Expr) ast.Expr {
	return e
}

func (p *parser) on_expr3__ternary(pred ast.Expr, _ Token, thenValue ast.Expr, _ Token, elseValue ast.Expr) ast.Expr {
	return &ast.TernaryExpr{
		Condition: pred,
		Then:      thenValue,
		Else:      elseValue,
	}
}

func (p *parser) on_expr3__binary(e ast.Expr) ast.Expr {
	return e
}

func (p *parser) on_binary_expr__unary(e ast.Expr) ast.Expr {
	return e
}

func (p *parser) on_binary_expr__binary(l ast.Expr, op Token, r ast.Expr) ast.Expr {
	if op.Type == PIPE {
		if funcCall, ok := r.(*ast.FuncCallExpr); ok {
			funcCall.Args = append(ast.Exprs{l}, funcCall.Args...)
			funcCall.Pipeline = true
			return funcCall
		}

		funcCall := &ast.FuncCallExpr{
			Target:   r,
			Args:     ast.Exprs{l},
			RetN:     1,
			Pipeline: true,
		}
		return funcCall
	}

	binExpr := &ast.BinaryExpr{
		Left:  l,
		Right: r,
	}

	switch op.Type {
	case MUL:
		binExpr.Op = ast.BinOpMult
	case DIV:
		binExpr.Op = ast.BinOpDiv
	case ADD:
		binExpr.Op = ast.BinOpPlus
	case SUB:
		binExpr.Op = ast.BinOpMinus
	case LT:
		binExpr.Op = ast.BinOpLT
	case LE:
		binExpr.Op = ast.BinOpLE
	case GT:
		binExpr.Op = ast.BinOpGT
	case GE:
		binExpr.Op = ast.BinOpGE
	case EQ:
		binExpr.Op = ast.BinOpEq
	case NE:
		binExpr.Op = ast.BinOpNE
	case AND:
		return &ast.AndExpr{Left: l, Right: r}
	case OR:
		return &ast.OrExpr{Left: l, Right: r}
	default:
		panic("unreachable")
	}

	return binExpr
}

func (p *parser) on_unary_expr__op(op Token, term ast.Expr) ast.Expr {
	var opType ast.UnaryOp
	switch op.Type {
	case NOT:
		opType = ast.UnaryOpNot
	case ADD:
		opType = ast.UnaryOpPlus
	case SUB:
		opType = ast.UnaryOpMinus
	default:
		panic("unreachable")
	}

	return &ast.UnaryExpr{
		Term: term,
		Op:   opType,
	}
}

func (p *parser) on_unary_expr__format(e ast.Expr, f Token) ast.Expr {
	if f.Type == FORMAT {
		e = &ast.FormatExpr{
			Target: e,
			Format: string(f.Str),
		}
	}

	return e
}

func (p *parser) on_primary_expr__simple_ref(id Token) ast.Expr {
	return &ast.SimpleRef{
		ID: p.tokenToNitro(id),
	}
}

func (p *parser) on_primary_expr__member_access(target ast.Expr, _ Token, member Token) ast.Expr {
	return &ast.MemberAccess{
		Target: target,
		Member: p.tokenToNitro(member),
	}
}

func (p *parser) on_primary_expr__index_or_call(target ast.Expr, delim Token, index ast.Expr, _ Token) ast.Expr {
	switch delim.Type {
	case OBRACKET:
		return &ast.IndexExpr{
			Target: target,
			Index:  index,
		}
	case OPAREN:
		funcCall, _ := index.(*ast.FuncCallExpr)
		if funcCall == nil {
			funcCall = &ast.FuncCallExpr{}
		}
		funcCall.Target = target
		funcCall.RetN = 1
		return funcCall
	default:
		panic("unreachable")
	}
}

func (p *parser) on_primary_expr__slice(target ast.Expr, _ Token, begin ast.Expr, _ Token, end ast.Expr, _ Token) ast.Expr {
	return &ast.SliceExpr{
		Target: target,
		Begin:  begin,
		End:    end,
	}
}

func (p *parser) on_primary_expr__forward(e ast.Expr) ast.Expr {
	return e
}

func (p *parser) on_primary_expr__parenthesis(_ Token, e ast.Expr, _ Token) ast.Expr {
	return e
}

func (p *parser) on_regex(r Token) ast.Expr {
	return &ast.RegexLiteral{
		RegexDef: string(r.Str),
	}
}

func (p *parser) on_simple_literal(lit Token) ast.Expr {
	return &ast.LiteralExpr{
		Val: p.tokenToNitro(lit),
	}
}

func (p *parser) on_simple_literal__expr(e ast.Expr) ast.Expr {
	return e
}

func (p *parser) on_string_literal__token(str Token) ast.Expr {
	v := string(str.Str)
	v = v[1 : len(v)-1] // remove quotes

	var err error
	switch str.Type {
	case STRING:
		v, err = expandEscapeSequences(v)
		if err != nil {
			p.errLogger.Failf(p.tokenPos(str), "invalid string literal: %w", err)
		}
	case RAW_STRING:
		v = strings.ReplaceAll(v, "\r", "")
	default:
		panic("unreachable")
	}

	t := token.Token{
		Type: token.String,
		Str:  v,
	}
	return &ast.LiteralExpr{
		Val: t,
	}
}

func (p *parser) on_string_literal__format(_ Token, parts []ast.Expr, _ Token) ast.Expr {
	return &ast.FormattedStringLiteral{
		Parts: parts,
	}
}

func (p *parser) on_fstr_part__token(t Token) ast.Expr {
	token := token.Token{
		Type: token.String,
	}

	switch t.Type {
	case FBSTR_LITERAL:
		token.Str = string(t.Str)
	case FBSTR_ESC_LITERAL:
		var err error
		token.Str, err = expandEscapeSequences(string(t.Str))
		if err != nil {
			p.errLogger.Failf(p.tokenPos(t), "%v", err)
		}
	default:
		panic("unreached")
	}

	return &ast.LiteralExpr{
		Val: token,
	}
}

func (p *parser) on_fstr_part__expr(_ Token, e ast.Expr, _ Token) ast.Expr {
	return e
}

func (p *parser) on_func_call_arg_list(args []ast.Expr, _ Token, expand Token) *ast.FuncCallExpr {
	return &ast.FuncCallExpr{
		Args:   args,
		Expand: expand.Type != 0,
	}
}

func (p *parser) on_lambda_expr(_, _ Token, params []*ast.FuncParam, _, _ Token, block *ast.StmtBlock, _ Token) ast.Expr {
	lambda := &ast.AnonFuncExpr{}
	lambda.Params = params
	lambda.Block = block
	return lambda
}

func (p *parser) on_short_lambda_expr(_ Token, params []*ast.FuncParam, _ Token, e ast.Expr) ast.Expr {
	lambda := &ast.AnonFuncExpr{}
	lambda.Params = params
	lambda.Block = &ast.StmtBlock{
		Stmts: ast.ASTs{
			&ast.ReturnStmt{
				Values: ast.Exprs{e},
			},
		},
	}
	return lambda
}

func (p *parser) on_exec_expr(_ Token, args []ast.Expr, _ Token) ast.Expr {
	// TODO: this should be an AST and the transformation should happen in the
	// Rewrite phase.

	var argm execArgMaker
	argm.Reserve(len(args))

	for _, arg := range args {
		err := argm.AddArg(arg)
		if err != nil {
			p.errLogger.Failf(arg.Pos(), "%s", err.Error())
			return nil
		}
	}
	argm.AddArg(nil)

	return &ast.FuncCallExpr{
		Target: &ast.SimpleRef{
			ID: token.Token{
				Str:  "$exec",
				Type: token.String,
			},
		},
		Args: ast.Exprs{argm.ArrayLiteral()},
		RetN: 1,
	}
}

func (p *parser) on_exec_expr_arg__token(t Token) ast.Expr {
	switch t.Type {
	case EXEC_LITERAL,
		EXEC_DQUOTE_LITERAL,
		EXEC_SQUOTE_LITERAL:
		return &ast.LiteralExpr{
			Val: p.tokenToNitro(t),
		}
	case EXEC_WS:
		return &ast.ExecWS{}
	case EXEC_HOME:
		return &ast.ExecHome{}
	default:
		panic("unreachable")
	}
}

func (p *parser) on_exec_expr_arg__expr(_ Token, e ast.Expr, expand Token, _ Token) ast.Expr {
	if expand.Type != 0 {
		e = &ast.ExecExpand{Expr: e}
	}
	return e
}

func (p *parser) on_object_literal(_ Token, fields []*ast.ObjectField, _ Token) ast.Expr {
	return &ast.ObjectLiteral{
		FieldBlock: &ast.ObjectFieldBlock{
			Fields: fields,
		},
	}
}

func (p *parser) on_object_fields(fields []*ast.ObjectField, _ Token) []*ast.ObjectField {
	return fields
}

func (p *parser) on_object_field__id_name(name string, _ Token, value ast.Expr) *ast.ObjectField {
	return &ast.ObjectField{
		NameID: name,
		Val:    value,
	}
}

func (p *parser) on_object_field__expr_name(_ Token, name ast.Expr, _, _ Token, value ast.Expr) *ast.ObjectField {
	return &ast.ObjectField{
		NameExpr: name,
		Val:      value,
	}
}

func (p *parser) on_array_literal(_ Token, elems []*ast.ArrayElement, _ Token) ast.Expr {
	return &ast.ArrayLiteral{
		Block: &ast.ArrayElementBlock{
			Elements: elems,
		},
	}
}

func (p *parser) on_array_elems(elems []*ast.ArrayElement, _ Token) []*ast.ArrayElement {
	return elems
}

func (p *parser) on_array_elem(elem ast.Expr, expand Token) *ast.ArrayElement {
	return &ast.ArrayElement{
		Val:    elem,
		Expand: expand.Type != 0,
	}
}

func (p *parser) on_comma_semicolon(t Token) Token {
	return t
}

func (p *parser) on_id_or_keyword(t Token) string {
	return string(t.Str)
}

func (p *parser) _onBounds(v any, begin, end Token) {
	ast, ok := v.(ast.AST)
	if !ok || ast == nil {
		return
	}
	ast.SetPos(p.tokenPos(begin))
}

func (p *parser) tokenPos(t Token) token.Pos {
	position := p.file.Position(t.Pos)
	return token.Pos{
		Filename: position.Filename,
		Line:     position.Line,
		Col:      position.Column,
	}
}

func (p *parser) tokenToNitro(at Token) token.Token {
	var err error

	t := token.Token{}
	switch at.Type {
	case NIL:
		t.Type = token.Nil

	case STRING,
		EXEC_DQUOTE_LITERAL,
		EXEC_SQUOTE_LITERAL:
		s := string(at.Str)
		s = s[1 : len(s)-1] // remove quotes
		t.Type = token.String
		t.Str = s
		t.Str, err = expandEscapeSequences(s)
		if err != nil {
			p.errLogger.Failf(p.tokenPos(at), "Invalid string literal: %v", err)
		}

	case CHAR:
		s := string(at.Str)
		s = s[1 : len(s)-1] // remove quotes
		r, err := expandCharEscapeSequences(s)
		if err != nil {
			p.errLogger.Failf(p.tokenPos(at), "Invalid character literal: %v", err)
		}
		t.Type = token.Int
		t.Int = int64(r)

	case NUMBER:
		s := string(at.Str)
		if strings.IndexByte(s, '.') == -1 {
			t.Type = token.Int
			t.Int, _ = strconv.ParseInt(s, 10, 64)
		} else {
			t.Type = token.Float
			t.Float, _ = strconv.ParseFloat(s, 64)
		}

	case TRUE:
		t.Type = token.Bool
		t.Bool = true

	case FALSE:
		t.Type = token.Bool
		t.Bool = false

	default:
		t.Type = token.String
		t.Str = string(at.Str)
	}

	t.Pos = p.tokenPos(at)
	return t
}

func (p *parser) tokens(ts []Token) []token.Token {
	tokens := make([]token.Token, len(ts))
	for i, t := range ts {
		tokens[i] = p.tokenToNitro(t)
	}
	return tokens
}

func (p *parser) lvalue(expr ast.Expr) *ast.LValue {
	switch expr.(type) {
	case *ast.SimpleRef:
	case *ast.MemberAccess:
	case *ast.IndexExpr:
	default:
		p.errLogger.Failf(expr.Pos(), "Expression is not lvalue")
		return nil
	}
	lv := &ast.LValue{
		Expr: expr,
	}
	lv.SetPos(expr.Pos())
	return lv
}

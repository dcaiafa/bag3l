package ast

import (
	"github.com/dcaiafa/bag3l/internal/symbol"
	"github.com/dcaiafa/bag3l/internal/token"
)

type ForStmt struct {
	PosImpl
	ForVars  []*ForVar
	IterExpr Expr
	Block    *StmtBlock
}

func (s *ForStmt) RunPass(ctx *Context, pass Pass) {
	if pass == Rewrite {
		s.rewrite(ctx)
	}

	// ForVars and IterExpr have been rewritten at this point.
	// Only Block remains.

	ctx.RunPassChild(s, s.Block, pass)
}

func (s *ForStmt) rewrite(ctx *Context) {
	/*

		   for a, b in e {
		   	...
		   }

		   =>

		   {
		   	var $iter = iterate(e)
		   	try {
					var a, b, $ok = next($iter)
		   		loop $ok {
		   			{
		   				...
		   			}

					cont:
						a, b, $ok = next($iter)
		   		}

		   		$close_iter($iter, true)
		   	} catch e {
		     	$close_iter($iter, false)
		   		throw e
		   	}
		   }

	*/

	forVars := []token.Token{}
	for _, v := range s.ForVars {
		forVars = append(forVars, v.VarName)
	}
	forVars = append(forVars, token.Token{
		Type: token.String,
		Str:  "$ok",
	})
	forRefs := make([]AST, len(forVars))
	for i, v := range forVars {
		forRefs[i] = &LValue{
			Expr: &SimpleRef{
				ID: v,
			},
		}
	}

	e := token.StringToken("e")

	block := &StmtBlock{
		PosImpl: s.PosImpl,
		Stmts: []AST{
			// var $iter = iterate(e)
			&VarDeclStmt{
				PosImpl: s.PosImpl,
				Vars: []token.Token{
					token.StringToken("$iter"),
				},
				InitValues: []Expr{
					&FuncCallExpr{
						Target: &SimpleRef{
							ID: token.Token{
								Type: token.String,
								Str:  "iterate",
							},
						},
						Args: []Expr{
							s.IterExpr,
						},
						RetN: 1,
					},
				},
			},

			// try {
			NewTryCatchStmt(
				&StmtBlock{
					Stmts: []AST{
						// var ..., $ok = next($iter)
						&VarDeclStmt{
							Vars: forVars,
							InitValues: []Expr{
								&FuncCallExpr{
									Target: &SimpleRef{ID: token.StringToken("next")},
									Args: []Expr{
										&SimpleRef{ID: token.StringToken("$iter")},
									},
									RetN: len(forVars) + 1,
								},
							},
						},

						// loop $ok {
						&Loop{
							PosImpl:   s.PosImpl,
							Predicate: &SimpleRef{ID: token.StringToken("$ok")},
							Block: &StmtBlock{
								PosImpl: s.PosImpl,
								Stmts: []AST{
									// if not $ok {
									&IfStmt{
										Sections: []*IfSection{
											&IfSection{
												Pred: &UnaryExpr{
													Op:   UnaryOpNot,
													Term: &SimpleRef{ID: token.StringToken("$ok")},
												},
												// break
												Block: &StmtBlock{
													Stmts: []AST{
														&BreakStmt{},
													},
												},
											},
										},
									},

									// ...
									s.Block,
								},
							},
							ContinueBlock: &StmtBlock{
								PosImpl: s.PosImpl,
								Stmts: []AST{
									// ..., $ok = next($iter)
									&AssignStmt{
										Lvalues: forRefs,
										Rvalues: []Expr{
											&FuncCallExpr{
												Target: &SimpleRef{ID: token.StringToken("next")},
												Args: []Expr{
													&SimpleRef{ID: token.StringToken("$iter")},
												},
												RetN: len(forVars) + 1,
											},
										},
									},
								},
							},
						},

						// $close_iter($iter, true)
						&ExprStmt{
							Expr: &FuncCallExpr{
								Target: &SimpleRef{ID: token.StringToken("$close_iter")},
								Args: []Expr{
									&SimpleRef{ID: token.StringToken("$iter")},
									&LiteralExpr{Val: token.BoolToken(true)},
								},
							},
						},
					},
				},
				// catch e {
				&e,
				&StmtBlock{
					Stmts: []AST{
						// $close_iter($iter, true)
						&ExprStmt{
							Expr: &FuncCallExpr{
								Target: &SimpleRef{ID: token.StringToken("$close_iter")},
								Args: []Expr{
									&SimpleRef{ID: token.StringToken("$iter")},
									&LiteralExpr{Val: token.BoolToken(false)},
								},
							},
						},

						// throw e
						&ThrowStmt{
							PosImpl: s.PosImpl,
							Expr:    &SimpleRef{ID: e},
						},
					},
				}),
		},
	}

	s.ForVars = nil
	s.IterExpr = nil
	s.Block = block
}

type ForVar struct {
	PosImpl
	VarName token.Token

	sym symbol.Symbol
}

func (s *ForVar) RunPass(ctx *Context, pass Pass) {}

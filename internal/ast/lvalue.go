package ast

// LValue is implemented by AST nodes that may appear on the left-hand side of
// an assignment. IsLValue is a sentinel marker used to distinguish lvalues from
// expressions; it is never called.
type LValue interface {
	AST
	IsLValue()
}

// lvalueReadExpr returns an expression that reads the current value held by the
// lvalue. It is used to desugar compound assignment (`a += b` => `a = a + b`)
// and increment/decrement (`a++` => `a = a + 1`), where the target is both read
// and written. The Target/Index sub-expressions are shared with the lvalue,
// mirroring how the previous wrapper-based design reused a single node.
func lvalueReadExpr(lv LValue) Expr {
	var e Expr
	switch lv := lv.(type) {
	case *SimpleRefLValue:
		e = &SimpleRef{ID: lv.ID}
	case *MemberAccessLValue:
		e = &MemberAccess{Target: lv.Target, Member: lv.Member}
	case *IndexLValue:
		e = &IndexExpr{Target: lv.Target, Index: lv.Index}
	default:
		panic("unreachable")
	}
	e.SetPos(lv.Pos())
	return e
}

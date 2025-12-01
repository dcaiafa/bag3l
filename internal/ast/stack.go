package ast

import "iter"

type Stack struct {
	stack []AST
}

func (s *Stack) Push(ast AST) {
	s.stack = append(s.stack, ast)
}

func (s *Stack) Pop() {
	s.stack = s.stack[:len(s.stack)-1]
}

func (s *Stack) All() iter.Seq[AST] {
	return func(yield func(AST) bool) {
		for i := len(s.stack) - 1; i >= 0; i-- {
			if !yield(s.stack[i]) {
				return
			}
		}
	}
}

package lib

import (
	"fmt"
	"math/rand"
	gosort "sort"

	"github.com/dcaiafa/bag3l/internal/vm"
)

type sortExpr struct {
	pathExpr *PathExpr
	desc     bool
}

type sorter struct {
	m     *vm.VM
	arr   *vm.List
	err   error
	less  vm.Callable
	exprs []*sortExpr
}

func (s *sorter) Len() int {
	return s.arr.Len()
}

func (s *sorter) Less(i, j int) bool {
	a := s.arr.Get(i)
	b := s.arr.Get(j)

	if s.err != nil {
		return false
	}

	if len(s.exprs) > 0 {
		for i, expr := range s.exprs {
			a, err := expr.pathExpr.Eval(s.m, a)
			if err != nil {
				s.err = err
				return false
			}
			if a == nil {
				return !expr.desc
			}
			b, err := expr.pathExpr.Eval(s.m, b)
			if err != nil {
				s.err = err
				return false
			}
			if b == nil {
				return expr.desc
			}

			if i < len(s.exprs)-1 {
				areEq, err := evalCmpOp(vm.OpEq, a, b)
				if err != nil {
					s.err = err
					return false
				}
				if areEq {
					continue
				}
			}

			op := vm.OpLT
			if expr.desc {
				op = vm.OpGT
			}
			res, err := evalCmpOp(op, a, b)
			if err != nil {
				s.err = err
				return false
			}
			return res
		}
	}

	if s.less != nil {
		res, err := s.m.Call(s.less, []vm.Value{a, b}, 1)
		if err != nil {
			s.err = err
			return false
		}
		return res[0].(vm.Bool).Bool()
	}

	if a == nil {
		return true
	}
	if b == nil {
		return false
	}

	res, err := vm.EvalOp(vm.OpLT, a, b)
	if err != nil {
		s.err = err
		return false
	}

	return res.(vm.Bool).Bool()
}

func (s *sorter) Swap(i, j int) {
	t := s.arr.Get(i)
	s.arr.Put(i, s.arr.Get(j))
	s.arr.Put(j, t)
}

func sort(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if len(args) < 1 {
		return nil, errNotEnoughArgs
	}

	arr, err := ListFromIter(m, args[0])
	if err != nil {
		return nil, err
	}

	s := &sorter{
		m:   m,
		arr: arr,
	}

	if len(args) >= 2 {
		switch arg2 := args[1].(type) {
		case vm.Callable:
			s.less = arg2
			if len(args) != 2 {
				return nil, errTooManyArgs
			}

		case vm.String:
			s.exprs = make([]*sortExpr, 0, len(args)-1)
			for _, arg := range args[1:] {
				sortExpr := new(sortExpr)
				var extra string
				sortExpr.pathExpr, extra, err = ParsePathExpr(arg)
				if err != nil {
					return nil, err
				}
				if len(extra) != 0 {
					if extra == "d" {
						sortExpr.desc = true
					} else if extra != "a" {
						return nil, fmt.Errorf("invalid sort qualifier %q", extra)
					}
				}
				s.exprs = append(s.exprs, sortExpr)
			}
		}
	}

	gosort.Sort(s)

	if s.err != nil {
		return nil, s.err
	}

	return []vm.Value{arr}, nil
}

func evalCmpOp(op vm.Op, operand1, operand2 vm.Value) (bool, error) {
	res, err := vm.EvalOp(op, operand1, operand2)
	if err != nil {
		return false, err
	}
	boolRes, ok := res.(vm.Bool)
	if !ok {
		return false, fmt.Errorf(
			"expected operation to return bool; returned %v instead",
			vm.TypeName(res))
	}
	return boolRes.Bool(), nil
}

func shuffle(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if len(args) < 1 {
		return nil, errNotEnoughArgs
	}

	arr, err := ListFromIter(m, args[0])
	if err != nil {
		return nil, err
	}

	rand.Shuffle(arr.Len(), func(i, j int) {
		t := arr.Get(i)
		arr.Put(i, arr.Get(j))
		arr.Put(j, t)
	})

	return []vm.Value{arr}, nil
}

package lib

import (
	"fmt"

	nitro "github.com/dcaiafa/bag3l"
)

type mapReduceSpec struct {
	Reduce nitro.Callable
	Pick   *PathExpr
}

func (s *mapReduceSpec) Convert(v nitro.Value) error {
	mapv, ok := v.(*nitro.Object)
	if !ok {
		return fmt.Errorf("reduce spec must be a map")
	}

	if mapv.Len() < 1 || mapv.Len() > 2 {
		return fmt.Errorf("invalid spec")
	}

	reduce, ok := mapv.Get(nitro.NewString("reduce"))
	if !ok {
		return fmt.Errorf("spec must specify a reduce function")
	}
	s.Reduce, ok = reduce.(nitro.Callable)
	if !ok {
		return fmt.Errorf("spec field reduce must be a function")
	}

	if mapv.Len() == 2 {
		pick, ok := mapv.Get(nitro.NewString("pick"))
		if !ok {
			return fmt.Errorf("invalid spec")
		}

		var err error
		s.Pick, _, err = ParsePathExpr(pick)
		if err != nil {
			return fmt.Errorf("spec field pick is invalid: %w", err)
		}
	}

	return nil
}

// map_reduce(iter, map_expr, []{ reduce: func, pick: func })

func mapReduce(m *nitro.VM, args []nitro.Value, nRet int) ([]nitro.Value, error) {
	if len(args) != 3 {
		return nil, errInvalidNumberOfArgs
	}

	iter, err := nitro.MakeIterator(m, args[0])
	if err != nil {
		return nil, fmt.Errorf("invalid argument 1: %w", err)
	}
	defer m.IterClose(iter)

	mapExpr, _, err := ParsePathExpr(args[1])
	if err != nil {
		return nil, fmt.Errorf("invalid argument 2: %w", err)
	}

	specsValue, err := getListArg(args, 2)
	if err != nil {
		return nil, err
	}

	specs := make([]*mapReduceSpec, specsValue.Len())

	for i := 0; i < specsValue.Len(); i++ {
		spec := new(mapReduceSpec)
		err = spec.Convert(specsValue.Get(i))
		if err != nil {
			return nil, fmt.Errorf("invalid reduce spec: %w", err)
		}
		specs[i] = spec
	}

	res := nitro.NewObject()
	for {
		cur, err := m.IterNext(iter, 1)
		if err != nil {
			return nil, err
		}
		if cur == nil {
			break
		}

		mapKey, err := mapExpr.Eval(m, cur[0])
		if err != nil {
			return nil, err
		}

		var accumList *nitro.Array
		accumRef, _ := res.IndexRef(mapKey)
		if *accumRef.Ref == nil {
			accumList = nitro.NewArrayFromSlice(make([]nitro.Value, len(specs)))
			*accumRef.Ref = accumList
		} else {
			accumList = (*accumRef.Ref).(*nitro.Array)
		}

		for i, spec := range specs {
			val := cur[0]
			if spec.Pick != nil {
				val, err = spec.Pick.Eval(m, cur[0])
				if err != nil {
					return nil, fmt.Errorf("error evaluting pick expression: %w", err)
				}
			}

			partialRes, err := m.Call(
				spec.Reduce, []nitro.Value{accumList.Get(i), val}, 1)
			if err != nil {
				return nil, err
			}
			accumList.Put(i, partialRes[0])
		}
	}

	err = nil
	res.ForEach(func(k, accum nitro.Value) bool {
		accumList := accum.(*nitro.Array)
		for i, spec := range specs {
			finalRes, callErr := m.Call(spec.Reduce, []nitro.Value{accumList.Get(i), nil}, 1)
			if callErr != nil {
				err = callErr
				return false
			}
			accumList.Put(i, finalRes[0])
		}
		return true
	})

	if err != nil {
		return nil, fmt.Errorf("final reduce failed: %w", err)
	}

	return []nitro.Value{res}, nil
}

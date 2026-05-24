package lib

import "github.com/dcaiafa/bag3l/internal/vm"

func unique(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if err := expectArgCount(args, 1, 2); err != nil {
		return nil, err
	}

	e, err := getIterArg(m, args, 0)
	if err != nil {
		return nil, err
	}

	var keyFn *PathExpr
	if len(args) == 2 {
		var err error
		keyFn, _, err = ParsePathExpr(args[1])
		if err != nil {
			return nil, err
		}
	}

	set := make(map[vm.Value]vm.Value)

	for {
		v, err := m.IterNext(e, 1)
		if err != nil {
			return nil, err
		}
		if v == nil {
			break
		}
		key := v[0]
		if keyFn != nil {
			key, err = keyFn.Eval(m, v[0])
			if err != nil {
				return nil, err
			}
		}

		set[key] = v[0]
	}

	arr := vm.NewListWithSlice(make([]vm.Value, 0, len(set)))
	for _, v := range set {
		arr.Add(v)
	}

	return []vm.Value{arr}, nil
}

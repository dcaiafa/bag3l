package lib

import (
	"github.com/dcaiafa/bag3l/internal/vm"
)

type histogramAccum struct {
	hist map[vm.Value]int
}

func newHistAccum() *histogramAccum {
	return &histogramAccum{
		hist: make(map[vm.Value]int),
	}
}

func (a *histogramAccum) String() string    { return "<histogram>" }
func (a *histogramAccum) Type() string      { return "histogram" }
func (a *histogramAccum) Traits() vm.Traits { return vm.TraitNone }

func (a *histogramAccum) Process(v vm.Value) {
	a.hist[v]++
}

func (a *histogramAccum) ToResult() *vm.Map {
	r := vm.NewMap()
	for k, v := range a.hist {
		r.Put(k, vm.NewInt(int64(v)))
	}
	return r
}

func hist(m *vm.VM, args []vm.Value, nret int) ([]vm.Value, error) {
	if len(args) > 2 {
		return nil, errTooManyArgs
	} else if len(args) < 1 {
		return nil, errNotEnoughArgs
	}

	if len(args) == 2 {
		var accum *histogramAccum
		if args[0] == nil {
			accum = newHistAccum()
		} else {
			var ok bool
			accum, ok = args[0].(*histogramAccum)
			if !ok {
				return nil, errExpectedArg(0, args[0], "hist")
			}
		}

		if args[1] == nil {
			return []vm.Value{accum.ToResult()}, nil
		}

		accum.Process(args[1])
		return []vm.Value{accum}, nil
	}

	iter, err := getIterArg(m, args, 0)
	if err != nil {
		return nil, err
	}

	accum := newHistAccum()

	for {
		v, err := m.IterNext(iter, 1)
		if err != nil {
			return nil, err
		}
		if v == nil {
			break
		}
		accum.Process(v[0])
	}

	return []vm.Value{accum.ToResult()}, nil
}

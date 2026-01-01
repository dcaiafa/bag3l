package global

import _p0 "github.com/dcaiafa/bag3l/internal/export"
import _p1 "github.com/dcaiafa/bag3l/internal/stub"
import _p2 "github.com/dcaiafa/bag3l/internal/vm"

func _args(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
	var err error
	_ = err
	if len(args) > 0 {
		return nil, _p1.ErrTooManyArgs
	}
	{
		_r0, err := args0(vm)
		if err != nil {
			return nil, err
		}
		return []_p2.Value{_r0}, nil
	}
}
func _batch(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
	var err error
	_ = err
	if len(args) < 2 {
		return nil, _p1.ErrInsufficientArgs
	}
	switch _a0 := args[0].(type) {
	case _p2.Iterable, _p2.Iterator:
		switch _a1 := args[1].(type) {
		case _p2.Int:
			if len(args) > 2 {
				return nil, _p1.ErrTooManyArgs
			}
			{
				_ta0 := _p1.MustMakeIter(vm, _a0)
				_ta1 := (_a1).Int64()
				_r0, err := batch0(vm, _ta0, _ta1)
				if err != nil {
					return nil, err
				}
				return []_p2.Value{_r0}, nil
			}
		default:
			return nil, _p1.InvalidArg(args, 1)
		}
	default:
		return nil, _p1.InvalidArg(args, 0)
	}
}

var Exports = _p0.Exports{
	{N: "args", T: _p0.Func, F: _args},
	{N: "batch", T: _p0.Func, F: _batch},
}

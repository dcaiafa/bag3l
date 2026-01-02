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
func _close(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
	var err error
	_ = err
	if len(args) < 1 {
		return nil, _p1.ErrInsufficientArgs
	}
	switch _a0 := args[0].(type) {
	case _p2.Value:
		if len(args) > 1 {
			return nil, _p1.ErrTooManyArgs
		}
		{
			_ta0 := _a0
			err := close0(vm, _ta0)
			if err != nil {
				return nil, err
			}
			return []_p2.Value{}, nil
		}
	default:
		return nil, _p1.InvalidArg(args, 0)
	}
}
func _color(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
	var err error
	_ = err
	if len(args) < 1 {
		return nil, _p1.ErrInsufficientArgs
	}
	var _a0 []_p2.Value = args[0:]
	{
		_ta0 := make([]string, len(_a0))
		for i := range _a0 {
			switch _a := _a0[i].(type) {
			case _p2.String:
				_ta0[i] = (_a).String()
			default:
				return nil, _p1.InvalidArg(args, 0)
			}
		}
		_r0, err := color0(vm, _ta0)
		if err != nil {
			return nil, err
		}
		return []_p2.Value{_r0}, nil
	}
}
func _deep_equal(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
	var err error
	_ = err
	if len(args) < 2 {
		return nil, _p1.ErrInsufficientArgs
	}
	switch _a0 := args[0].(type) {
	case _p2.Value:
		switch _a1 := args[1].(type) {
		case _p2.Value:
			if len(args) > 2 {
				return nil, _p1.ErrTooManyArgs
			}
			{
				_ta0 := _a0
				_ta1 := _a1
				_r0, err := deep_equal0(vm, _ta0, _ta1)
				if err != nil {
					return nil, err
				}
				return []_p2.Value{_p2.NewBool(_r0)}, nil
			}
		default:
			return nil, _p1.InvalidArg(args, 1)
		}
	default:
		return nil, _p1.InvalidArg(args, 0)
	}
}
func _discard(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
	var err error
	_ = err
	if len(args) < 1 {
		return nil, _p1.ErrInsufficientArgs
	}
	switch _a0 := args[0].(type) {
	case _p2.Value:
		if len(args) > 1 {
			return nil, _p1.ErrTooManyArgs
		}
		{
			_ta0 := _a0
			err := discard0(vm, _ta0)
			if err != nil {
				return nil, err
			}
			return []_p2.Value{}, nil
		}
	default:
		return nil, _p1.InvalidArg(args, 0)
	}
}
func _do(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
	var err error
	_ = err
	if len(args) < 2 {
		return nil, _p1.ErrInsufficientArgs
	}
	switch _a0 := args[0].(type) {
	case _p2.Iterable, _p2.Iterator:
		switch _a1 := args[1].(type) {
		case _p2.Callable:
			if len(args) > 2 {
				return nil, _p1.ErrTooManyArgs
			}
			{
				_ta0 := _p1.MustMakeIter(vm, _a0)
				_ta1 := _a1
				err := do0(vm, _ta0, _ta1)
				if err != nil {
					return nil, err
				}
				return []_p2.Value{}, nil
			}
		default:
			return nil, _p1.InvalidArg(args, 1)
		}
	default:
		return nil, _p1.InvalidArg(args, 0)
	}
}
func _dont_close(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
	var err error
	_ = err
	if len(args) < 1 {
		return nil, _p1.ErrInsufficientArgs
	}
	switch _a0 := args[0].(type) {
	case _p2.Iterable, _p2.Iterator:
		if len(args) > 1 {
			return nil, _p1.ErrTooManyArgs
		}
		{
			_ta0 := _p1.MustMakeIter(vm, _a0)
			_r0, err := dont_close0(vm, _ta0)
			if err != nil {
				return nil, err
			}
			return []_p2.Value{_r0}, nil
		}
	default:
		return nil, _p1.InvalidArg(args, 0)
	}
}
func _enumerate(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
	var err error
	_ = err
	if len(args) < 1 {
		return nil, _p1.ErrInsufficientArgs
	}
	switch _a0 := args[0].(type) {
	case _p2.Iterable, _p2.Iterator:
		if len(args) > 1 {
			return nil, _p1.ErrTooManyArgs
		}
		{
			_ta0 := _p1.MustMakeIter(vm, _a0)
			_r0, err := enumerate0(vm, _ta0)
			if err != nil {
				return nil, err
			}
			return []_p2.Value{_r0}, nil
		}
	default:
		return nil, _p1.InvalidArg(args, 0)
	}
}

var Exports = _p0.Exports{
	{N: "args", T: _p0.Func, F: _args},
	{N: "batch", T: _p0.Func, F: _batch},
	{N: "close", T: _p0.Func, F: _close},
	{N: "color", T: _p0.Func, F: _color},
	{N: "deep_equal", T: _p0.Func, F: _deep_equal},
	{N: "discard", T: _p0.Func, F: _discard},
	{N: "do", T: _p0.Func, F: _do},
	{N: "dont_close", T: _p0.Func, F: _dont_close},
	{N: "enumerate", T: _p0.Func, F: _enumerate},
}

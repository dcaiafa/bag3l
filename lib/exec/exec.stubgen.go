package exec

import _p0 "github.com/dcaiafa/bag3l/internal/export"
import _p1 "github.com/dcaiafa/bag3l/internal/stub"
import _p2 "github.com/dcaiafa/bag3l/internal/vm"

func _exec(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
	var err error
	_ = err
	if len(args) < 1 {
		return nil, _p1.ErrInsufficientArgs
	}
	switch _a0 := args[0].(type) {
	case *_p2.List:
		if len(args) > 1 {
			return nil, _p1.ErrTooManyArgs
		}
		{
			_ta0 := _a0
			_r0, err := exec0(vm, _ta0)
			if err != nil {
				return nil, err
			}
			return []_p2.Value{_r0}, nil
		}
	case _p2.Reader, _p2.Readable:
		if len(args) <= 1 {
			return nil, _p1.ErrInsufficientArgs
		}
		switch _a1 := args[1].(type) {
		case *_p2.List:
			if len(args) > 2 {
				return nil, _p1.ErrTooManyArgs
			}
			{
				_ta0 := _p1.MustMakeReader(vm, _a0)
				_ta1 := _a1
				_r0, err := exec1(vm, _ta0, _ta1)
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
func _with_combined_output(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
	var err error
	_ = err
	if len(args) < 1 {
		return nil, _p1.ErrInsufficientArgs
	}
	switch _a0 := args[0].(type) {
	case *Process:
		if len(args) > 1 {
			return nil, _p1.ErrTooManyArgs
		}
		{
			_ta0 := _a0
			_r0, err := with_combined_output0(vm, _ta0)
			if err != nil {
				return nil, err
			}
			return []_p2.Value{_r0}, nil
		}
	default:
		return nil, _p1.InvalidArg(args, 0)
	}
}
func _with_dir(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
	var err error
	_ = err
	if len(args) < 2 {
		return nil, _p1.ErrInsufficientArgs
	}
	switch _a0 := args[0].(type) {
	case *Process:
		switch _a1 := args[1].(type) {
		case _p2.String:
			if len(args) > 2 {
				return nil, _p1.ErrTooManyArgs
			}
			{
				_ta0 := _a0
				_ta1 := (_a1).String()
				_r0, err := with_dir0(vm, _ta0, _ta1)
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
func _with_env(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
	var err error
	_ = err
	if len(args) < 1 {
		return nil, _p1.ErrInsufficientArgs
	}
	switch _a0 := args[0].(type) {
	case *Process:
		var _a1 []_p2.Value = args[1:]
		{
			_ta0 := _a0
			_ta1 := make([]string, len(_a1))
			for i := range _a1 {
				switch _a := _a1[i].(type) {
				case _p2.String:
					_ta1[i] = (_a).String()
				default:
					return nil, _p1.InvalidArg(args, 1)
				}
			}
			_r0, err := with_env0(vm, _ta0, _ta1)
			if err != nil {
				return nil, err
			}
			return []_p2.Value{_r0}, nil
		}
	default:
		return nil, _p1.InvalidArg(args, 0)
	}
}
func _with_partial_env(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
	var err error
	_ = err
	if len(args) < 1 {
		return nil, _p1.ErrInsufficientArgs
	}
	switch _a0 := args[0].(type) {
	case *Process:
		var _a1 []_p2.Value = args[1:]
		{
			_ta0 := _a0
			_ta1 := make([]string, len(_a1))
			for i := range _a1 {
				switch _a := _a1[i].(type) {
				case _p2.String:
					_ta1[i] = (_a).String()
				default:
					return nil, _p1.InvalidArg(args, 1)
				}
			}
			_r0, err := with_partial_env0(vm, _ta0, _ta1)
			if err != nil {
				return nil, err
			}
			return []_p2.Value{_r0}, nil
		}
	default:
		return nil, _p1.InvalidArg(args, 0)
	}
}
func _with_stderr(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
	var err error
	_ = err
	if len(args) < 2 {
		return nil, _p1.ErrInsufficientArgs
	}
	switch _a0 := args[0].(type) {
	case *Process:
		switch _a1 := args[1].(type) {
		case _p2.Writer:
			if len(args) > 2 {
				return nil, _p1.ErrTooManyArgs
			}
			{
				_ta0 := _a0
				_ta1 := _a1
				_r0, err := with_stderr0(vm, _ta0, _ta1)
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
	{N: "exec", T: _p0.Func, F: _exec},
	{N: "with_combined_output", T: _p0.Func, F: _with_combined_output},
	{N: "with_dir", T: _p0.Func, F: _with_dir},
	{N: "with_env", T: _p0.Func, F: _with_env},
	{N: "with_partial_env", T: _p0.Func, F: _with_partial_env},
	{N: "with_stderr", T: _p0.Func, F: _with_stderr},
}

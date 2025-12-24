package os

import (
	_p0 "github.com/dcaiafa/bag3l/internal/export"
	_p1 "github.com/dcaiafa/bag3l/internal/stub"
	_p2 "github.com/dcaiafa/bag3l/internal/vm"
)

func _get_workdir(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
	var err error
	_ = err
	if len(args) > 0 {
		return nil, _p1.ErrTooManyArgs
	}
	{
		_r0, err := get_workdir0(vm)
		if err != nil {
			return nil, err
		}
		return []_p2.Value{_p2.NewString(_r0)}, nil
	}
}
func _home_dir(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
	var err error
	_ = err
	if len(args) > 0 {
		return nil, _p1.ErrTooManyArgs
	}
	{
		_r0, err := home_dir0(vm)
		if err != nil {
			return nil, err
		}
		return []_p2.Value{_p2.NewString(_r0)}, nil
	}
}
func _set_workdir(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
	var err error
	_ = err
	if len(args) < 1 {
		return nil, _p1.ErrInsufficientArgs
	}
	switch _a0 := args[0].(type) {
	case _p2.String:
		if len(args) > 1 {
			return nil, _p1.ErrTooManyArgs
		}
		{
			_ta0 := (_a0).String()
			err := set_workdir0(vm, _ta0)
			if err != nil {
				return nil, err
			}
			return []_p2.Value{}, nil
		}
	default:
		return nil, _p1.InvalidArg(args, 0)
	}
}

var Exports = _p0.Exports{
	{N: "get_workdir", T: _p0.Func, F: _get_workdir},
	{N: "home_dir", T: _p0.Func, F: _home_dir},
	{N: "set_workdir", T: _p0.Func, F: _set_workdir},
}

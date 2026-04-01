package io

import _p0 "github.com/dcaiafa/bag3l/internal/export"
import _p1 "github.com/dcaiafa/bag3l/internal/stub"
import _p2 "github.com/dcaiafa/bag3l/internal/vm"
import _p3 "runtime"

func _err(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
	var err error
	_ = err
	if len(args) == 0 {
		var _a0 _p2.Reader = nil
		{
			_ta0 := _p1.MustMakeReader(vm, _a0)
			_r0, err := err0(vm, _ta0)
			if err != nil {
				return nil, err
			}
			return []_p2.Value{_r0}, nil
		}
	}
	switch _a0 := args[0].(type) {
	case _p2.Reader, _p2.Readable:
		if len(args) > 1 {
			return nil, _p1.ErrTooManyArgs
		}
		{
			_ta0 := _p1.MustMakeReader(vm, _a0)
			_r0, err := err0(vm, _ta0)
			if err != nil {
				return nil, err
			}
			return []_p2.Value{_r0}, nil
		}
	default:
		return nil, _p1.InvalidArg(args, 0)
	}
}
func _from_crlf(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
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
			_r0, err := from_crlf0(vm, _ta0)
			if err != nil {
				return nil, err
			}
			return []_p2.Value{_p2.NewString(_r0)}, nil
		}
	case _p2.Writer:
		if len(args) > 1 {
			return nil, _p1.ErrTooManyArgs
		}
		{
			_ta0 := _a0
			_r0, err := from_crlf2(vm, _ta0)
			if err != nil {
				return nil, err
			}
			return []_p2.Value{_r0}, nil
		}
	case _p2.Reader, _p2.Readable:
		if len(args) > 1 {
			return nil, _p1.ErrTooManyArgs
		}
		{
			_ta0 := _p1.MustMakeReader(vm, _a0)
			_r0, err := from_crlf1(vm, _ta0)
			if err != nil {
				return nil, err
			}
			return []_p2.Value{_r0}, nil
		}
	default:
		return nil, _p1.InvalidArg(args, 0)
	}
}
func _in(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
	var err error
	_ = err
	if len(args) > 0 {
		return nil, _p1.ErrTooManyArgs
	}
	{
		_r0, err := in0(vm)
		if err != nil {
			return nil, err
		}
		return []_p2.Value{_r0}, nil
	}
}
func _out(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
	var err error
	_ = err
	if len(args) == 0 {
		var _a0 _p2.Reader = nil
		{
			_ta0 := _p1.MustMakeReader(vm, _a0)
			_r0, err := out0(vm, _ta0)
			if err != nil {
				return nil, err
			}
			return []_p2.Value{_r0}, nil
		}
	}
	switch _a0 := args[0].(type) {
	case _p2.Reader, _p2.Readable:
		if len(args) > 1 {
			return nil, _p1.ErrTooManyArgs
		}
		{
			_ta0 := _p1.MustMakeReader(vm, _a0)
			_r0, err := out0(vm, _ta0)
			if err != nil {
				return nil, err
			}
			return []_p2.Value{_r0}, nil
		}
	default:
		return nil, _p1.InvalidArg(args, 0)
	}
}
func _to_crlf(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
	var err error
	_ = err
	if len(args) < 1 {
		return nil, _p1.ErrInsufficientArgs
	}
	switch _a0 := args[0].(type) {
	case _p2.String:
		if len(args) == 1 {
			var _a1 _p2.Bool = _p2.NewBool(_p3.GOOS == "windows")
			{
				_ta0 := (_a0).String()
				_ta1 := (_a1).Bool()
				_r0, err := to_crlf0(vm, _ta0, _ta1)
				if err != nil {
					return nil, err
				}
				return []_p2.Value{_p2.NewString(_r0)}, nil
			}
		}
		switch _a1 := args[1].(type) {
		case _p2.Bool:
			if len(args) > 2 {
				return nil, _p1.ErrTooManyArgs
			}
			{
				_ta0 := (_a0).String()
				_ta1 := (_a1).Bool()
				_r0, err := to_crlf0(vm, _ta0, _ta1)
				if err != nil {
					return nil, err
				}
				return []_p2.Value{_p2.NewString(_r0)}, nil
			}
		default:
			return nil, _p1.InvalidArg(args, 1)
		}
	case _p2.Reader, _p2.Readable:
		if len(args) == 1 {
			var _a1 _p2.Bool = _p2.NewBool(_p3.GOOS == "windows")
			{
				_ta0 := _p1.MustMakeReader(vm, _a0)
				_ta1 := (_a1).Bool()
				_r0, err := to_crlf1(vm, _ta0, _ta1)
				if err != nil {
					return nil, err
				}
				return []_p2.Value{_r0}, nil
			}
		}
		switch _a1 := args[1].(type) {
		case _p2.Bool:
			if len(args) > 2 {
				return nil, _p1.ErrTooManyArgs
			}
			{
				_ta0 := _p1.MustMakeReader(vm, _a0)
				_ta1 := (_a1).Bool()
				_r0, err := to_crlf1(vm, _ta0, _ta1)
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
	{N: "err", T: _p0.Func, F: _err},
	{N: "from_crlf", T: _p0.Func, F: _from_crlf},
	{N: "in", T: _p0.Func, F: _in},
	{N: "out", T: _p0.Func, F: _out},
	{N: "to_crlf", T: _p0.Func, F: _to_crlf},
}

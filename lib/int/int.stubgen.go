package int

import _p0 "github.com/dcaiafa/bag3l/internal/export"
import _p1 "github.com/dcaiafa/bag3l/internal/stub"
import _p2 "github.com/dcaiafa/bag3l/internal/vm"

func _into(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
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
			_r0, err := into0(vm, _ta0)
			if err != nil {
				return nil, err
			}
			return []_p2.Value{_p2.NewInt(_r0)}, nil
		}
	default:
		return nil, _p1.InvalidArg(args, 0)
	}
}
func _parse(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
	var err error
	_ = err
	if len(args) < 1 {
		return nil, _p1.ErrInsufficientArgs
	}
	switch _a0 := args[0].(type) {
	case _p2.String:
		if len(args) == 1 {
			var _a1 _p2.Int = _p2.NewInt(0)
			{
				_ta0 := (_a0).String()
				_ta1 := (_a1).Int64()
				_r0, err := parse0(vm, _ta0, _ta1)
				if err != nil {
					return nil, err
				}
				return []_p2.Value{_p2.NewInt(_r0)}, nil
			}
		}
		switch _a1 := args[1].(type) {
		case _p2.Int:
			if len(args) > 2 {
				return nil, _p1.ErrTooManyArgs
			}
			{
				_ta0 := (_a0).String()
				_ta1 := (_a1).Int64()
				_r0, err := parse0(vm, _ta0, _ta1)
				if err != nil {
					return nil, err
				}
				return []_p2.Value{_p2.NewInt(_r0)}, nil
			}
		default:
			return nil, _p1.InvalidArg(args, 1)
		}
	default:
		return nil, _p1.InvalidArg(args, 0)
	}
}

var Exports = _p0.Exports{
	{N: "into", T: _p0.Func, F: _into},
	{N: "parse", T: _p0.Func, F: _parse},
}

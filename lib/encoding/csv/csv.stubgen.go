package csv

import _p0 "github.com/dcaiafa/bag3l/internal/export"
import _p1 "github.com/dcaiafa/bag3l/internal/stub"
import _p2 "github.com/dcaiafa/bag3l/internal/vm"

type DecodeOptions struct {
	Columns *_p2.List
}

func (m *DecodeOptions) FromMap(v *_p2.Map) error {
	var err error
	_ = err
	v.ForEach(func(k, v _p2.Value) bool {
		n, ok := k.(_p2.String)
		if !ok {
			err = _p1.ErrMapKeyMustBeStr
			return false
		}
		switch n.String() {
		case "columns":
			cv, ok := v.(*_p2.List)
			if !ok {
				err = _p1.ErrInvalidFieldType
				return false
			}
			tv := cv
			m.Columns = tv
		default:
			err = _p1.StructDoesNotHaveField("DecodeOptions", n.String())
			return false
		}
		return true
	})
	if err != nil {
		return err
	}
	return nil
}
func _decode(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
	var err error
	_ = err
	if len(args) < 1 {
		return nil, _p1.ErrInsufficientArgs
	}
	switch _a0 := args[0].(type) {
	case _p2.Reader, _p2.Readable:
		if len(args) == 1 {
			var _a1 *_p2.Map = nil
			{
				_ta0 := _p1.MustMakeReader(vm, _a0)
				var _ta1 *DecodeOptions
				if _a1 != nil {
					_ta1 = new(DecodeOptions)
					err = _ta1.FromMap(_a1)
					if err != nil {
						return nil, _p1.InvalidArgErr(args, 1, err)
					}
				}
				_r0, err := decode0(vm, _ta0, _ta1)
				if err != nil {
					return nil, err
				}
				return []_p2.Value{_r0}, nil
			}
		}
		switch _a1 := args[1].(type) {
		case *_p2.Map:
			if len(args) > 2 {
				return nil, _p1.ErrTooManyArgs
			}
			{
				_ta0 := _p1.MustMakeReader(vm, _a0)
				var _ta1 *DecodeOptions
				if _a1 != nil {
					_ta1 = new(DecodeOptions)
					err = _ta1.FromMap(_a1)
					if err != nil {
						return nil, _p1.InvalidArgErr(args, 2, err)
					}
				}
				_r0, err := decode0(vm, _ta0, _ta1)
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
func _encode(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
	var err error
	_ = err
	if len(args) < 1 {
		return nil, _p1.ErrInsufficientArgs
	}
	switch _a0 := args[0].(type) {
	case _p2.Iterable, _p2.Iterator:
		if len(args) == 1 {
			var _a1 _p2.Writer = nil
			{
				_ta0 := _p1.MustMakeIter(vm, _a0)
				_ta1 := _a1
				err := encode0(vm, _ta0, _ta1)
				if err != nil {
					return nil, err
				}
				return []_p2.Value{}, nil
			}
		}
		switch _a1 := args[1].(type) {
		case _p2.Writer:
			if len(args) > 2 {
				return nil, _p1.ErrTooManyArgs
			}
			{
				_ta0 := _p1.MustMakeIter(vm, _a0)
				_ta1 := _a1
				err := encode0(vm, _ta0, _ta1)
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

var Exports = _p0.Exports{
	{N: "decode", T: _p0.Func, F: _decode},
	{N: "encode", T: _p0.Func, F: _encode},
}

package runtime

import _p0 "github.com/dcaiafa/bag3l/internal/export"
import _p1 "github.com/dcaiafa/bag3l/internal/stub"
import _p2 "github.com/dcaiafa/bag3l/internal/vm"

func _script_dir(vm *_p2.VM, args []_p2.Value, nret int) ([]_p2.Value, error) {
	var err error
	_ = err
	if len(args) > 0 {
		return nil, _p1.ErrTooManyArgs
	}
	{
		_r0, err := script_dir0(vm)
		if err != nil {
			return nil, err
		}
		return []_p2.Value{_p2.NewString(_r0)}, nil
	}
}

var Exports = _p0.Exports{
	{N: "script_dir", T: _p0.Func, F: _script_dir},
}

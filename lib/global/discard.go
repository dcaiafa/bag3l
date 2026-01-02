package global

import (
	"fmt"
	"io"

	"github.com/dcaiafa/bag3l/internal/vm"
	"github.com/dcaiafa/bag3l/lib/core"
)

func discard0(m *vm.VM, v vm.Value) error {
	if reader, ok := v.(io.Reader); ok {
		_, err := io.Copy(io.Discard, reader)
		if err != nil {
			return err
		}
		core.CloseReader(reader)
	} else if iter, err := vm.MakeIterator(m, v); err == nil {
		for {
			r, err := m.IterNext(iter, 1)
			if err != nil {
				return err
			}
			if r == nil {
				break
			}
		}
	} else {
		return fmt.Errorf(
			"arg #1 %v is not reader or iter",
			vm.TypeName(v))
	}

	return nil
}

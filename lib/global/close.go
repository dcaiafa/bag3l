package global

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/dcaiafa/bag3l/internal/vm"
)

func close0(m *vm.VM, v vm.Value) error {
	switch arg := v.(type) {
	case io.Closer:
		err := arg.Close()
		if err != nil && !errors.Is(err, os.ErrClosed) {
			return err
		}

	case vm.Iterator:
		err := m.IterClose(arg)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf(
			"arg #1 %v is not closeable",
			vm.TypeName(v))
	}

	return nil
}

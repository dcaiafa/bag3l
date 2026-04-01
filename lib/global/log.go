package global

import (
	"github.com/dcaiafa/bag3l/internal/vm"
	libio "github.com/dcaiafa/bag3l/lib/io"
)

func log0(m *vm.VM, args []vm.Value) error {
	return basePrint(libio.Stderr(m), m, args)
}

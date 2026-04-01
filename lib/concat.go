package lib

import (
	"strings"

	"github.com/dcaiafa/bag3l/internal/vm"
)

// $concat - used by the compiler.
func concat(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	str := strings.Builder{}

	for _, arg := range args {
		if arg == nil {
			str.WriteString("<nil>")
		} else {
			str.WriteString(arg.String())
		}
	}

	return []vm.Value{vm.NewString(str.String())}, nil
}

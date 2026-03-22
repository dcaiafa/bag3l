package hex

import (
	"encoding/hex"

	"github.com/dcaiafa/bag3l/internal/vm"
)

//go:generate go run ../../../internal/stub/stubgen hex.stubgen

func encode0(m *vm.VM, v string) (string, error) {
	return hex.EncodeToString([]byte(v)), nil
}

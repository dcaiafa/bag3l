package float

import (
	"strconv"

	"github.com/dcaiafa/bag3l/internal/vm"
)

//go:generate go run ../../internal/stub/stubgen float.stubgen

func parse0(m *vm.VM, v string) (float64, error) {
	return strconv.ParseFloat(v, 64)
}

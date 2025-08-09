package lib

import (
	"testing"

	"github.com/dcaiafa/bag3l/internal/vm"
)

func TestDontClose(t *testing.T) {
	RunSubO(t, "with", `
		var it = range(4)

		print(it | dont_close | take(2) | to_list)
		print(it | to_list)
`, `
[0 1]
[2 3]
`)

	RunSubErr(t, "without", `
		var it = range(4)

		print(it | take(2) | to_list)
		print(it | to_list)
`, vm.ErrIteratorClosed)
}

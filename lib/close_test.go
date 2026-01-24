package lib

import (
	"testing"

	"github.com/dcaiafa/bag3l/internal/vm"
)

func TestClose(t *testing.T) {
	RunSubO(t, "close_iterator", `
		var it = range(4)
		var v, ok = next(it)
		print(v)
		close(it)
		print(f"iter_closed={harness.is_iter_closed(it)}")
`, `
0
iter_closed=true
`)

	RunSubErr(t, "use_closed_iterator", `
		var it = range(4)
		close(it)
		next(it)
`, vm.ErrIteratorClosed)

	RunSubO(t, "close_file", `
		var f = file.create_temp()
		defer file.remove(f)
		close(f)
`, `
`)

	RunSubO(t, "close_already_closed", `
		var f = file.create_temp()
		defer file.remove(f)
		close(f)
		close(f)
`, `
`)
}

package lib

import "testing"

func TestTakeUntil(t *testing.T) {
	RunSubO(t, "basic", `
		var it = [1, 2, 3, 4] |
			take_until(&v -> v == 2)
		it | print
		print(f"iter_closed={harness.is_iter_closed(it)}")
`, `
1
2
iter_closed=true
`)
	RunSubO(t, "close", `
		var it = [1, 2, 3, 4] | iterate
		it |
			take_until(&v -> v == 3) |
			take(1) |
			print
		print(f"iter_closed={harness.is_iter_closed(it)}")
`, `
1
iter_closed=true
`)

}

package lib

import "testing"

func TestTakeWhile(t *testing.T) {
	RunSubO(t, "exclusive", `
		[1, 2, 3, 4] |
			take_while(&v -> v != 2) |
			print
`, `
1
`)
	RunSubO(t, "inclusive", `
		[1, 2, 3, 4] |
			take_while(&v -> v != 2, true) |
			print
`, `
1
2
`)
	RunSubO(t, "close", `
		var it = [1, 2, 3, 4] | iterate
		it |
			take_while(&v -> v != 2) |
			take(1) |
			print
		print(f"iter_closed={harness.is_iter_closed(it)}")
`, `
1
iter_closed=true
`)

}

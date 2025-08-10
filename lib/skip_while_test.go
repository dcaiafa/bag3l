package lib

import "testing"

func TestSkipWhile(t *testing.T) {
	RunSubO(t, "inclusive", `
		[1, 2, 3, 4] |
			skip_while(&v -> v != 2) |
			print
`, `
2
3
4
`)
	RunSubO(t, "exclusive", `
		[1, 2, 3, 4] |
			skip_while(&v -> v != 2, true) |
			print
`, `
3
4
`)
	RunSubO(t, "close", `
  	var it = [1, 2, 3, 4] | iterate
		it |
			skip_while(&v -> v != 2) |
			take(1) |
			print
		print(f"iter_closed={harness.is_iter_closed(it)}")
`, `
2
iter_closed=true
`)

}

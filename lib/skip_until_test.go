package lib

import "testing"

func TestSkipUntil(t *testing.T) {
	RunSubO(t, "basic", `
		[1, 2, 3, 4, 5] |
			skip_until(&v -> v == 3) |
			print
`, `
4
5
`)
	RunSubO(t, "close", `
  	var it = [1, 2, 3, 4, 5] | iterate
		it |
			skip_until(&v -> v == 3) |
			take(1) |
			print
		print(f"iter_closed={harness.is_iter_closed(it)}")
`, `
4
iter_closed=true
`)

}

package global_test

import (
	"testing"

	"github.com/dcaiafa/bag3l/internal/btesting"
)

func TestPrint(t *testing.T) {
	btesting.RunSubO(t, "single_string", `
		print("hello")
	`, `hello`)

	btesting.RunSubO(t, "single_int", `
		print(42)
	`, `42`)

	btesting.RunSubO(t, "single_float", `
		print(3.14)
	`, `3.14`)

	btesting.RunSubO(t, "single_bool", `
		print(true)
	`, `true`)

	btesting.RunSubO(t, "multiple_args", `
		print("hello", "world")
	`, `hello world`)

	btesting.RunSubO(t, "mixed_types", `
		print("count", 3, "pi", 3.14)
	`, `count 3 pi 3.14`)

	btesting.RunSubO(t, "iterator", `
		print(range(3))
	`, `0
1
2`)

	btesting.RunSubO(t, "iterator_multi_ret", `
		print(enumerate(["a", "b", "c"]))
	`, `0 a
1 b
2 c`)

	btesting.RunSubO(t, "reader", `
		var r = ["hello", "world"] | map(str.into) | stream
		print(r)
	`, `hello
world`)

	btesting.RunSubO(t, "list", `
		print([1, 2, 3])
	`, `[1 2 3]`)

	btesting.RunSubO(t, "map", `
		print({a: 1})
	`, `{a: 1}`)
}

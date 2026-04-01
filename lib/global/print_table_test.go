package global_test

import (
	"testing"

	"github.com/dcaiafa/bag3l/internal/btesting"
)

func TestPrintTable(t *testing.T) {
	runSubO(t, "basic", `
		[{name: "alice", age: 30}, {name: "bob", age: 25}] | print_table
	`, `
name  age
alice 30
bob   25`)

	runSubO(t, "single_row", `
		[{x: 1, y: 2}] | print_table
	`, `
x y
1 2`)

	runSubO(t, "missing_key", `
		[{a: 1, b: 2}, {a: 3}] | print_table
	`, `
a b
1 2
3`)

	runSubO(t, "nil_value", `
		[{a: 1, b: nil}] | print_table
	`, `
a b
1 <nil>`)

	runSubO(t, "with_padding", `
		[{a: 1, b: 2}] | print_table({padding: 3})
	`, `
a   b
1   2`)

	btesting.RunSubErr(t, "invalid_arg", `
		print_table(123)
	`, nil)
}

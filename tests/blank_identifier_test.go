package tests

import "testing"

// TestBlankIdentifierDiscard covers `_` as a discard target: the assigned value
// is dropped instead of stored. It is valid in assignments, var declarations,
// and for loops, and never needs (or creates) a declaration.
func TestBlankIdentifierDiscard(t *testing.T) {
	RunSubO(t, "assign_discard_first", `
		func two() { return 1, 2 }
		var b
		_, b = two()
		print(b)
	`, `2`)

	RunSubO(t, "assign_discard_second", `
		func two() { return 1, 2 }
		var a
		a, _ = two()
		print(a)
	`, `1`)

	RunSubO(t, "assign_discard_single", `
		func two() { return 1, 2 }
		_ = two()
		print("ok")
	`, `ok`)

	RunSubO(t, "assign_discard_pre_assignment_value", `
		var a = 1
		var b = 2
		a, _ = b, a
		print(a)
	`, `2`)

	RunSubO(t, "var_decl_discard", `
		func two() { return 1, 2 }
		var _, b = two()
		print(b)
	`, `2`)

	RunSubO(t, "var_decl_no_init_noop", `
		var _
		print("ok")
	`, `ok`)

	RunSubO(t, "for_discard_index", `
		for _, e in enumerate(["a", "b"]) {
			print(e)
		}
	`, `
a
b
`)

	RunSubO(t, "for_discard_value", `
		for i, _ in enumerate(["a", "b"]) {
			print(i)
		}
	`, `
0
1
`)

	RunSubO(t, "for_discard_only_var", `
		var n = 0
		for _ in range(1, 4) {
			n += 1
		}
		print(n)
	`, `3`)
}

// TestBlankIdentifierForbidden verifies that `_` cannot be declared as a
// parameter or read as a value.
func TestBlankIdentifierForbidden(t *testing.T) {
	RunSubErr(t, "param", `
		func f(_) { return 1 }
		print(f(0))
	`, nil)

	RunSubErr(t, "read", `print(_)`, nil)

	RunSubErr(t, "read_in_expr", `var x = _ + 1`, nil)

	RunSubErr(t, "compound_assign", `_ += 1`, nil)

	RunSubErr(t, "increment", `_++`, nil)
}

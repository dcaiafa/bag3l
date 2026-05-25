package tests

import "testing"

// TestStoreMultiAssign exercises the direct-store model for assignment
// statements. Multi-assignment evaluates the whole right-hand side first and
// then stores into the lvalues right-to-left, so swaps and self-referential
// assignments must still observe the pre-assignment values.
func TestStoreMultiAssign(t *testing.T) {
	RunSubO(t, "swap_simple", `
		var a, b = 1, 2
		a, b = b, a
		print(a, b)
	`, `2 1`)

	RunSubO(t, "swap_array_elems", `
		var a = [1, 2, 3]
		a[0], a[2] = a[2], a[0]
		print(a[0], a[1], a[2])
	`, `3 2 1`)

	RunSubO(t, "mixed_simple_and_index", `
		var a = [10, 20]
		var b
		a[0], b = 99, a[1]
		print(a[0], a[1], b)
	`, `99 20 20`)

	RunSubO(t, "two_index_lvalues", `
		var a = [1, 2]
		var b = [3, 4]
		a[0], b[1] = b[0], a[1]
		print(a[0], a[1], b[0], b[1])
	`, `3 2 3 2`)

	RunSubO(t, "member_swap", `
		var o = { x: 1, y: 2 }
		o.x, o.y = o.y, o.x
		print(o.x, o.y)
	`, `2 1`)

	RunSubO(t, "multi_return_into_index", `
		func two() { return 100, 200 }
		var a = [0, 0]
		a[0], a[1] = two()
		print(a[0], a[1])
	`, `100 200`)
}

// TestStoreIndexAndMember covers single index/member stores backed by
// Indexable.SetIndex, including inserting new map keys.
func TestStoreIndexAndMember(t *testing.T) {
	RunSubO(t, "array_set", `
		var a = [1, 2, 3]
		a[1] = 10
		print(a[0], a[1], a[2])
	`, `1 10 3`)

	RunSubO(t, "map_insert_new_key", `
		var m = {}
		m["a"] = 1
		m.b = 2
		print(m.a, m.b, len(m))
	`, `1 2 2`)

	RunSubO(t, "compound_and_incdec_on_index", `
		var a = [1, 2, 3]
		a[1] += 10
		a[1]++
		print(a[0], a[1], a[2])
	`, `1 13 3`)

	RunSubErr(t, "array_set_out_of_range", `
		var a = [1, 2]
		a[5] = 9
	`, nil)

	RunSubErr(t, "string_is_immutable", `
		var s = "abc"
		s[0] = 1
	`, nil)
}

// TestStoreVariableKinds exercises each store target: global, local, parameter,
// lifted local, lifted parameter, and capture.
func TestStoreVariableKinds(t *testing.T) {
	RunSubO(t, "global", `
		var g = 1
		func bump() { g = g + 1 }
		bump()
		bump()
		print(g)
	`, `3`)

	RunSubO(t, "param", `
		func f(x) {
			x = x + 5
			return x
		}
		print(f(10))
	`, `15`)

	RunSubO(t, "lifted_local_via_capture", `
		func f() {
			var x = 1
			var g = func() { x = x + 10 }
			g()
			return x
		}
		print(f())
	`, `11`)

	RunSubO(t, "lifted_param", `
		func f(x) {
			var g = func() { return x }
			x = x + 1
			return g()
		}
		print(f(10))
	`, `11`)

	RunSubO(t, "capture_store_counter", `
		func make_counter() {
			var n = 0
			return func() {
				n = n + 1
				return n
			}
		}
		var c = make_counter()
		print(c(), c(), c())
	`, `1 2 3`)
}

// TestStoreCatchVar verifies the thrown error is stored into the catch variable,
// including when the catch variable is lifted by a closure that captures it.
func TestStoreCatchVar(t *testing.T) {
	RunSubO(t, "simple", `
		try {
			throw "boom"
		} catch e {
			print(e.error)
		}
	`, `boom`)

	RunSubO(t, "lifted_catch_var", `
		func f() {
			var g
			try {
				throw "oops"
			} catch e {
				g = func() { return e.error }
			}
			return g()
		}
		print(f())
	`, `oops`)
}

package tests

import (
	"testing"

	"github.com/dcaiafa/bag3l/internal/vm"
)

func TestStackOverflow(t *testing.T) {
	RunSubO(t, "deep_recursion_ok", `
		func f(n) {
			if n == 0 {
				return 0
			}
			return 1 + f(n - 1)
		}
		print(f(100))
`, `100`)

	RunSubErr(t, "infinite_recursion", `
		func f(n) {
			return 1 + f(n + 1)
		}
		f(0)
`, vm.ErrStackOverflow)

	RunSubErr(t, "infinite_recursion_many_locals", `
		func f(n) {
			var a = n
			var b = a
			var c = b
			return [a, b, c, f(n + 1)]
		}
		f(0)
`, vm.ErrStackOverflow)

	// Each iterator runs on its own operand stack, so only the call-depth
	// limit stops this one.
	RunSubErr(t, "infinite_recursive_iterator", `
		func f(n) {
			for v in f(n + 1) {
				yield v
			}
			yield n
		}
		for v in f(0) {
			print(v)
		}
`, vm.ErrStackOverflow)

	RunSubErr(t, "native_call_recursion", `
		import "harness"
		func f(g) {
			return harness.call(g, g)
		}
		f(f)
`, vm.ErrStackOverflow)

	RunSubErr(t, "arg_expansion", `
		import "list"
		func f() {
			return narg()
		}
		var l = []
		for i in range(5000) {
			list.append(l, i)
		}
		print(f(l...))
`, vm.ErrStackOverflow)

	RunSubErr(t, "not_catchable", `
		func f(n) {
			return 1 + f(n + 1)
		}
		try {
			f(0)
		} catch e {
			print("caught")
		}
`, vm.ErrStackOverflow)

	RunSubO(t, "arg_expansion_ok", `
		import "list"
		func f() {
			return narg()
		}
		var l = []
		for i in range(100) {
			list.append(l, i)
		}
		print(f(l...))
`, `100`)
}

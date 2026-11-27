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

	// Deep enough to grow the operand stack several times.
	RunSubO(t, "deeper_recursion_ok", `
		func f(n) {
			var a = n
			var b = a + 1
			if n == 0 {
				return 0
			}
			return b - a + f(n - 1)
		}
		print(f(5000))
`, `5000`)

	// Calls made from an iterator's body run on the iterator's own stack,
	// which must keep its locals intact as it grows between yields.
	RunSubO(t, "iterator_stack_growth", `
		func depth(n) {
			if n == 0 {
				return 0
			}
			return 1 + depth(n - 1)
		}
		func gen() {
			var a = "a"
			yield a, depth(10)
			var b = "b"
			yield a + b, depth(200)
			yield a + b, depth(2000)
		}
		for s, d in gen() {
			print(s, d)
		}
`, `
a 10
ab 200
ab 2000
`)

	RunSubO(t, "coroutine_stack_growth", `
		import "co"
		func f(n) {
			if n == 0 {
				return 0
			}
			return 1 + f(n - 1)
		}
		co.start(func() {
			print(f(3000))
		})
		print(f(4000))
`, `
3000
4000
`)

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
		func f() {
			return narg()
		}
		var l = [range(1100000)...]
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
		for i in range(5000) {
			list.append(l, i)
		}
		print(f(l...))
`, `5000`)
}

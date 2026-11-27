package tests

import (
	"testing"

	"github.com/dcaiafa/bag3l/internal/vm"
)

func TestIter(t *testing.T) {
	RunSubO(t, "simple", `
		func iterable() {
			yield 1
			yield 2
		}
		for i in iterable() {
    	print(i)
		}
`, `
1
2
`)

	RunSubO(t, "simple2", `
		func iterable() {
			yield 1,10
			yield 2,20
		}
		for i,j in iterable() {
    	print(i,j)
		}
`, `
1 10
2 20
`)

	RunSubO(t, "a_little_more_complex", `
		func counter(n, m) {
			var last
			for i in range(n) {
				var res = mod(i, m)
				yield res, last
				last = res
			}
		}
		for i, l in counter(5, 3) {
			print(i, l)
		}
	`, `
0 <nil>
1 0
2 1
0 2
1 0
`)

	RunSubO(t, "return", `
		func iter(x) {
			for i in range(1000) {
				if i > x {
					return
				}
				yield i
			}
		}
		for i in iter(3) {
    	print(i)
		}
`, `
0
1
2
3
`)

	RunSubErr(t, "call_again_after_ended", `
		func iter() {
			yield 1
			yield 2
		}
		var i = iter()
		for x in i {
		}
		for x in i {
		}
`, vm.ErrIteratorClosed)

	RunSubO(t, "nil", `
		print("begin")
		for x, y in nil {
    	print("not reached")
		}
		print("end")
`, `
begin
end
`)

	RunSubO(t, "nested", `
func sum(iter) {
	var s = 0
	for v in iter {
		s = s + v
	}
	return s
}

func pass(iter) {
	for v in iter {
		yield v
	}
}

range(4) | pass | sum | print
`, `
6
`)
}

// Calling an iterator value directly (instead of through next or a for loop)
// yields one value per call, and false once the body has finished.
func TestIterDirectCall(t *testing.T) {
	RunSubO(t, "no_locals", `
		func it() {
			yield 1
			yield 2
		}
		var i = it()
		print(i())
		print(i())
		print(i())
		print(i())
`, `
1
2
false
false
`)

	RunSubO(t, "with_locals", `
		func it() {
			var a = 10
			var b = 20
			yield a
			yield b
		}
		var i = it()
		print(i())
		print(i())
		print(i())
		print(i())
`, `
10
20
false
false
`)

	RunSubO(t, "multiple_values", `
		func it() {
			var a = 1
			yield a, "x"
		}
		var i = it()
		var p, q = i()
		print(p, q)
		p, q = i()
		print(p, q)
`, `
1 x
false <nil>
`)

	RunSubO(t, "finishes_in_catch", `
		func it() {
			var a = 5
			try {
				yield 1
				throw "boom"
			} catch e {
				yield a
				yield e.error
			}
		}
		var i = it()
		print(i())
		print(i())
		print(i())
		print(i())
`, `
1
5
boom
false
`)
}

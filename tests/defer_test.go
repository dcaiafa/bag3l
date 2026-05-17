package tests

import "testing"

func TestDefer(t *testing.T) {
	RunSubO(t, "from_main", `
		var x = 1
		defer print(x)
	`, `1`)

	RunSubO(t, "from_func", `
		func f() {
			var x = 1
      defer print(x)
			print(x)
			x = x + 1
		}
		f()
	`, `
1
2
`)

	RunSubO(t, "multi", `
		func f() {
			var x = 1
			defer print("bye")
			defer print(x)
			print(x)
			x = x + 1
		}
		f()
`, `
1
2
bye
`)

	RunSubO(t, "exception", `
		func f() {
			var x = 1
			defer func() {
				print("bye from", x)
			}()

			print("bomb")
			throw "boom"
		}

		try {
			f()
		} catch e {
			print("it did go", e.error)
		}
`, `
bomb
bye from 1
it did go boom
`)
}

func TestDeferInIterator(t *testing.T) {
	RunSubO(t, "runs_once_at_end", `
		func it() {
			defer print("cleanup")
			yield 1
			yield 2
			yield 3
		}
		for x in it() {
			print(x)
		}
	`, `
1
2
3
cleanup
`)

	RunSubO(t, "runs_on_break", `
		func it() {
			defer print("cleanup")
			yield 1
			yield 2
			yield 3
		}
		for x in it() {
			print(x)
			if x == 2 {
				break
			}
		}
		print("after")
	`, `
1
2
cleanup
after
`)

	RunSubO(t, "runs_on_explicit_close", `
		func it() {
			defer print("cleanup")
			yield 1
			yield 2
			yield 3
		}
		var i = it()
		var v1, ok = next(i)
		print(v1)
		var v2
		v2, ok = next(i)
		print(v2)
		close(i)
		print("after")
	`, `
1
2
cleanup
after
`)

	RunSubO(t, "lifo_order", `
		func it() {
			defer print("a")
			defer print("b")
			defer print("c")
			yield 1
			yield 2
		}
		for x in it() {
			print(x)
		}
	`, `
1
2
c
b
a
`)

	RunSubO(t, "throw_escapes", `
		func it() {
			defer print("cleanup")
			yield 1
			throw "boom"
		}

		try {
			for x in it() {
				print(x)
			}
		} catch e {
			print("caught", e.error)
		}
	`, `
1
cleanup
caught boom
`)

	RunSubO(t, "close_is_idempotent", `
		func it() {
			defer print("cleanup")
			yield 1
			yield 2
		}
		var i = it()
		var v, ok = next(i)
		print(v)
		close(i)
		close(i)
		print("after")
	`, `
1
cleanup
after
`)
}

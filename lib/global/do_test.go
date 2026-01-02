package global_test

import (
	"testing"

	"github.com/dcaiafa/bag3l/internal/btesting"
)

func TestDo(t *testing.T) {
	btesting.RunSubO(t, "do_basic", `
    var sum = 0
    func add(v) { sum = sum + v }
    range(5) | do(add)
    print(sum)
  `, `10`)

	btesting.RunSubO(t, "do_with_print", `
    func f(v) { print(v) }
    range(3) | do(f)
  `, `0
1
2`)

	btesting.RunSubO(t, "do_empty_iter", `
    func f(v) { print(v) }
    [] | do(f)
    print("done")
  `, `done`)

	btesting.RunSubO(t, "do_with_map", `
    func f(m) { print(m.a) }
    [{a: 1}, {a: 2}] | do(f)
  `, `1
2`)
}

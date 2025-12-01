package tests

import "testing"

func TestIndexExpr(t *testing.T) {
	RunSubO(t, "array", `
		var a = [1, 2, 3]
		a[1] = 10
		print(a[0], a[1], a[2])
	`, `1 10 3`)

	RunSubErr(t, "array_bad_index", `
		var a = [1, 2, 3]
		print(a[3])
`, nil)

	RunSubO(t, "object", `
		var o = { 
    	a: {
				b: "foo"
				c: "bar"
			}
		}
		o.a.b = "fur"
		print(o["a"]["b"], o["a"]["c"], o["a"]["z"]?, o["z"]["a"]?)
	`, `fur bar <nil> <nil>`)

	RunSubO(t, "nil", `
		var a
		print(a[1]?)
	`, `<nil>`)
}

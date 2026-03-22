package tests

import "testing"

func TestOptional(t *testing.T) {
	RunSubO(t, "valid", `
		var o = { 
    	a: {
				b: "foo"
				c: [ { x: 1 } ],
			}
		}
		print(o.a.c[0].x?)
		print(o.a.c[1].x?)
		print(o.w.c[0].x?)

		var x = o.a.c[0].x?
		print(x)
`, `
1
<nil>
<nil>
1
`)
	RunSubErr(t, "invalid1", `
		var o = { 
    	a: {
				b: "foo"
				c: [ { x: 1 } ],
			}
		}
		print(o.a?.c[0].x)
`, nil)
	RunSubErr(t, "invalid2", `
		var o = { 
    	a: {
				b: "foo"
				c: [ { x: 1 } ],
			}
		}
		print(o.a.c?[0].x)
`, nil)
	RunSubErr(t, "invalid3", `
		var o = { 
    	a: {
				b: "foo"
				c: [ { x: 1 } ],
			}
		}
		print(o.a.c[0]?.x)
`, nil)
}

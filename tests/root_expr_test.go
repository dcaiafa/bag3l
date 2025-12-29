package tests

import "testing"

func TestRootExpr(t *testing.T) {
	RunSubO(t, "member_access", `
{
	hello: {
  	world: 1
	}
} | .hello.world | print
`, `1`)
	RunSubErr(t, "member_access_neg", `
{
	hello: {
  	world: 1
	}
} | .hello.wurld | print
`, nil)
	RunSubO(t, "member_access_opt", `
{
	hello: {
  	world: 1
	}
} | .hello.world? | print
`, `1`)
	RunSubO(t, "member_access_opt_neg", `
{
	hello: {
  	world: 1
	}
} | .hello.wurld? | print
`, `<nil>`)
	RunSubO(t, "index", `
[1, 2, 3] | .[1] | print
`, `2`)
	RunSubO(t, "index_string", `
{
	["the-thing"]: 123
} | .["the-thing"] | print
`, `123`)
	RunSubErr(t, "index_neg", `
[1, 2, 3] | .[4] | print
`, nil)
	RunSubO(t, "index_opt", `
[1, 2, 3] | .[1]? | print
`, `2`)
	RunSubO(t, "index_opt_neg", `
[1, 2, 3] | .[4]? | print
`, `<nil>`)
	RunSubO(t, "complex", `
{
	x: [
	  nil
		{ y: "hi" }
	]
} | .x[1].y | print
`, `hi`)
	RunSubO(t, "complex_opt", `
{
	x: [
	  nil
		{ y: "hi" }
	]
} | .x[0].y? | print
`, `<nil>`)
}

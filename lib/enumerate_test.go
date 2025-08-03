package lib

import "testing"

func TestEnumerate(t *testing.T) {
	RunSubO(t, "basic", `
		var a = ["a", "b", "c"]
		for i, e in enumerate(a) {
    	print(i, e)
		}
`, `
0 a
1 b
2 c
`)

	RunSubO(t, "just_index", `
		var a = ["a", "b", "c"]
		for i in enumerate(a) {
    	print(i, a[i])
		}
`, `
0 a
1 b
2 c
`)
}

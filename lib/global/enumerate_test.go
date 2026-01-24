package global_test

import (
	"testing"

	"github.com/dcaiafa/bag3l/internal/btesting"
)

func TestEnumerate(t *testing.T) {
	btesting.RunSubO(t, "basic", `
		var a = ["a", "b", "c"]
		for i, e in enumerate(a) {
    	print(i, e)
		}
`, `0 a
1 b
2 c`)

	btesting.RunSubO(t, "just_index", `
		var a = ["a", "b", "c"]
		for i in enumerate(a) {
    	print(i, a[i])
		}
`, `0 a
1 b
2 c`)
}

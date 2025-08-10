package lib

import (
	"testing"

	"github.com/dcaiafa/bag3l/internal/stub"
)

func TestIndex(t *testing.T) {
	RunSubO(t, "positive", `
		var v, ok = {a:1} | index("a")
		print(v, ok)
`, `1 true`)
	RunSubO(t, "negative", `
		var v, ok = {a:1} | index("b")
		print(v, ok)
`, `<nil> false`)
	RunSubErr(t, "invalid_target", `
		var v, ok = 1 | index("b")
`, stub.ErrInvalidArg)
	RunSubErr(t, "invalid_key", `
		var v, ok = [1,2,3] | index("b")
`, nil)
}

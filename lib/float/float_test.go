package float_test

import (
	"testing"

	"github.com/dcaiafa/bag3l/internal/btesting"
)

func TestFloatParse(t *testing.T) {
	btesting.RunSubO(t, "integer", `float.parse("42") | print`, `42`)
	btesting.RunSubO(t, "decimal", `float.parse("3.14") | print`, `3.14`)
	btesting.RunSubO(t, "negative", `float.parse("-2.5") | print`, `-2.5`)
	btesting.RunSubO(t, "scientific", `float.parse("1.5e2") | print`, `150`)
	btesting.RunSubErr(t, "invalid", `float.parse("abc")`, nil)
}

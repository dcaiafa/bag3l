package bool_test

import (
	"testing"

	"github.com/dcaiafa/bag3l/internal/btesting"
)

func TestBoolInto(t *testing.T) {
	btesting.RunSubO(t, "true", `bool.into(true) | print`, `true`)
	btesting.RunSubO(t, "false", `bool.into(false) | print`, `false`)
	btesting.RunSubO(t, "int_nonzero", `bool.into(1) | print`, `true`)
	btesting.RunSubO(t, "int_zero", `bool.into(0) | print`, `true`)
	btesting.RunSubO(t, "string_nonempty", `bool.into("hello") | print`, `true`)
	btesting.RunSubO(t, "string_empty", `bool.into("") | print`, `true`)
	btesting.RunSubO(t, "list", `bool.into([1, 2]) | print`, `true`)
}

package int_test

import (
	"testing"

	"github.com/dcaiafa/bag3l/internal/btesting"
)

func TestIntParse(t *testing.T) {
	btesting.RunSubO(t, "decimal", `int.parse("42") | print`, `42`)
	btesting.RunSubO(t, "negative", `int.parse("-7") | print`, `-7`)
	btesting.RunSubO(t, "hex", `int.parse("ff", 16) | print`, `255`)
	btesting.RunSubO(t, "binary", `int.parse("1010", 2) | print`, `10`)
	btesting.RunSubO(t, "octal", `int.parse("77", 8) | print`, `63`)
	btesting.RunSubO(t, "zero_prefix_auto", `int.parse("0xff") | print`, `255`)
	btesting.RunSubErr(t, "invalid", `int.parse("abc")`, nil)
}

func TestIntInto(t *testing.T) {
	btesting.RunSubO(t, "int", `int.into(42) | print`, `42`)
	btesting.RunSubO(t, "float", `int.into(3.7) | print`, `3`)
	btesting.RunSubErr(t, "string", `int.into("hello")`, nil)
	btesting.RunSubErr(t, "bool", `int.into(true)`, nil)
}

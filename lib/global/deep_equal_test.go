package global_test

import (
	"testing"

	"github.com/dcaiafa/bag3l/internal/btesting"
)

func TestDeepEqual(t *testing.T) {
	btesting.RunSubO(t, "basic_positive", `
var x = { foo: "bar" }
var a = {
	a: 1
	b: [1, { c: "hi" }]
	d: x
}
var b = {
	a: 1
	b: [1, { c: "hi" }]
	d: x
}
print(deep_equal(a, b))
`, "true")
	btesting.RunSubO(t, "basic_negative", `
var x = { foo: "bar" }
var a = {
	a: 1
	b: [1, { c: "hi" }]
	d: x
}
var b = {
	a: 1
	b: [1, { c: "bye" }]
	d: x
}
print(deep_equal(a, b))
`, "false")

}

package lib

import (
	"testing"
)

func TestColor(t *testing.T) {
	RunSubO(t, "color_type", `
		var c = color("red")
		print(type(c))
`, `
color
`)

	RunSubO(t, "color_multiple_attrs", `
		var c = color("bold", "red")
		print(type(c))
`, `
color
`)

	RunSubErr(t, "color_invalid_attr", `
		color("invalid_color_name")
`, nil)

	RunSubErr(t, "color_no_args", `
		color()
`, nil)
}

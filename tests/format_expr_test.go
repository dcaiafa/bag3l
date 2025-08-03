package tests

import "testing"

func TestFormatExpr(t *testing.T) {
	RunSubO(t, "", `print(3.1415%.2f)`, "3.14")
}

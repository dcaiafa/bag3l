package lib

import "testing"

func TestFlatten(t *testing.T) {
	RunSubO(t, "simple", "[[1,2,3],4,[5,6]]|flatten|list.from_iter|print", "[1 2 3 4 5 6]")
	RunSubO(t, "empty", "[]|flatten|list.from_iter|print", "[]")
	RunSubErr(t, "usage", "123|flatten|list.from_iter|print", nil)
}

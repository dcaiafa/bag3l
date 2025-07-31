package lib

import "testing"

func TestFlatten(t *testing.T) {
	RunSubO(t, "simple", "[[1,2,3],4,[5,6]]|flatten|to_list|print", "[1 2 3 4 5 6]")
	RunSubO(t, "empty", "[]|flatten|to_list|print", "[]")
	RunSubErr(t, "usage", "123|flatten|to_list|print", nil)
}

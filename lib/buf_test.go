package lib

import "testing"

func TestBuf(t *testing.T) {
	RunSubO(t, ``, `
		var b = buf.new()
		["hello", "world"] | stream | b()
		print(b)
	`, `
hello
world
	`)
}

package lib

import "testing"

func TestBuf(t *testing.T) {
	RunSubO(t, ``, `
		var b = buf.new()
		["hello", "world"] | stream_lines | b()
		print(b)
	`, `
hello
world
	`)
}

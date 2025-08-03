package tests

import "testing"

func TestStringLiteral(t *testing.T) {
	// TODO: unicode sequences
	RunSubO(t, "simple_escape_seq", `print("Hello\n\r\t\\World")`, "Hello\n\r\t\\World")
	RunSubO(t, "hex", `print("Hello\x2cWorld")`, "Hello,World")
	RunSubO(t, "char", `print('A', '\n', '\'')`, `65 10 39`)
	RunSubO(t, "raw", "print(`hello\\nworld`)", `hello\nworld`)
}

func TestFormattedString(t *testing.T) {
	RunSubO(t, "", `print(f"One plus two is {1+2}!")`, "One plus two is 3!")
	RunSubO(t, "", `print(f"{"hello"} world")`, "hello world")
	RunSubO(t, "", `print(f"")`, "")
	RunSubO(t, "", `print(f"Bye\n{"bye"}")`, "Bye\nbye")
	RunSubO(t, "", `print(f"a{"b"}c")`, "abc")
	RunSubO(t, "", `print(f"a\{1}c")`, `a{1}c`)
}

func TestStringSlice(t *testing.T) {
	RunSubO(t, "full", `
		print("hello"[0:4])
	`, `hell`)

	RunSubO(t, "implicit_start", `
		print("hello"[:4])
	`, `hell`)

	RunSubO(t, "implicit_end", `
		print("hello"[1:])
	`, `ello`)

	RunSubO(t, "relative_end", `
		print("hello"[:-1])
	`, `hell`)

	RunSubO(t, "relative_end2", `
		print("hello"[1:-1])
	`, `ell`)

	RunSubO(t, "relative_end3", `
		print("hello"[1:-2])
	`, `el`)

	RunSubO(t, "large_start", `
		print("hello"[100:])
	`, ``)

	RunSubO(t, "large_end", `
		print("hello"[:100])
	`, `hello`)

	RunSubO(t, "intersect", `
	  print("hello"[3:2])
	`, ``)

	RunSubErr(t, "err_negative_start", `
		print("hello"[-1:])
	`, nil)
}

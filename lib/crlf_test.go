package lib

import (
	"strings"
	"testing"
)

func TestFromCRLF(t *testing.T) {
	run := func(name, input, output string) {
		RunSubO(t, name+"-s", `"`+input+`" | from_crlf | (&r->"["+r+"]") | print`, "["+output+"]")
		RunSubO(t, name+"-r", strings.ReplaceAll(`
	var b = buf.new()
	$input$ | b
	print(f"[{b | from_crlf | read}]")
`, `$input$`, `"`+input+`"`), "["+output+"]")
		RunSubO(t, name+"-w", strings.ReplaceAll(`
	var b = buf.new()
	var b2 = b | from_crlf

	$input$ | b2
	print(f"[{b | read}]")
`, `$input$`, `"`+input+`"`), "["+output+"]")
	}

	run("crlf", `abc\r\ndef\r\n`, "abc\ndef\n")
	run("crlflf", `abc\r\n\ndef\r\n`, "abc\n\ndef\n")
	run("lf", `abc\ndef\r\n`, "abc\ndef\n")
	run("lf_end", `abc\ndef\n`, "abc\ndef\n")
	run("cr", `abc\rdef\r\n`, "abcdef\n")
}

func TestToCRLF(t *testing.T) {
	run := func(name, input, output string) {
		RunSubO(t, name, `"`+input+`" | to_crlf(true) | (&r->"["+r+"]") | print`, "["+output+"]")
	}

	run("crlf", `abc\r\ndef\r\n`, "abc\r\ndef\r\n")
	run("crcrlf", `abc\r\r\ndef\r\n`, "abc\r\r\ndef\r\n")
	run("crlflf", `abc\r\n\ndef\r\n`, "abc\r\n\r\ndef\r\n")
	run("lf", `abc\ndef\r\n`, "abc\r\ndef\r\n")
	run("lf_end", `abc\ndef\n`, "abc\r\ndef\r\n")
	run("cr", `abc\rdef\n`, "abc\rdef\r\n")
	run("cr_end", `abc\ndef\r`, "abc\r\ndef\r")

	RunSubO(t, "reader", `
	var b = buf.new()
	"abc\ndef\n" | b
	b = b | to_crlf(true)

	print(f"[{b | read(4)}{b | read()}]")
`, "[abc\r\ndef\r\n]")
}

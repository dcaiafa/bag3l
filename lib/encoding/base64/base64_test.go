package base64_test

import (
	"testing"

	"github.com/dcaiafa/bag3l/internal/btesting"
)

func TestBase64_Decode(t *testing.T) {
	btesting.RunSubO(t, "std", `"\x14\xfb\x9c\x03\xd9\x7e" | base64.encode | print`, `FPucA9l+`)
	btesting.RunSubO(t, "url", `"\x14\xfb\x9c\x03\xd9\x7e" | base64.encode({mode:"url"}) | print`, `FPucA9l-`)

	btesting.RunSubO(t, "reader", `
		var data = buf.new()
		"hello world" | data
		data | base64.encode | print
`, `aGVsbG8gd29ybGQ=`)
}

func TestBase64_Encode(t *testing.T) {
	btesting.RunSubO(t, "std", `"FPucA9l+" | base64.decode | print`, "\x14\xfb\x9c\x03\xd9\x7e")
	btesting.RunSubO(t, "url", `"FPucA9l-" | base64.decode({mode:"url"}) | print`, "\x14\xfb\x9c\x03\xd9\x7e")

	btesting.RunSubO(t, "reader", `
		var data = buf.new()
		"aGVsbG8gd29ybGQ=" | data
		data | base64.decode | print
`, `hello world`)
}

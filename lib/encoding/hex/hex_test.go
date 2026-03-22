package hex_test

import (
	"testing"

	"github.com/dcaiafa/bag3l/internal/btesting"
)

func TestHexEncode(t *testing.T) {
	btesting.RunSubO(t, "basic", `"hello" | hex.encode | print`, `68656c6c6f`)
	btesting.RunSubO(t, "empty", `"" | hex.encode | print`, ``)
	btesting.RunSubO(t, "binary", `"\x00\xff" | hex.encode | print`, `00ff`)
}

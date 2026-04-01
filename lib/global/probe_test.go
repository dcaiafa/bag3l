package global_test

import (
	"testing"

	"github.com/dcaiafa/bag3l/internal/btesting"
)

func TestProbe(t *testing.T) {
	btesting.RunSubO(t, "passthrough_value", `
		var x = probe(42)
		print(x)
	`, `42`)

	btesting.RunSubO(t, "passthrough_string", `
		var x = probe("hello")
		print(x)
	`, `hello`)

	btesting.RunSubO(t, "passthrough_iter", `
		var r = range(3) | probe | list.into
		print(r)
	`, `[0 1 2]`)

	btesting.RunSubO(t, "in_pipeline", `
		var r = range(3) | probe | map(&v -> v * 2) | list.into
		print(r)
	`, `[0 2 4]`)

	btesting.RunSubO(t, "with_label", `
		var x = probe(42, "debug:")
		print(x)
	`, `42`)
}

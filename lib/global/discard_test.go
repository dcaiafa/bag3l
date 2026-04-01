package global_test

import (
	"testing"

	"github.com/dcaiafa/bag3l/internal/btesting"
)

func TestDiscard(t *testing.T) {
	btesting.RunSubO(t, "discard_iter", `
    range(10) | discard
    print("done")
  `, `done`)

	btesting.RunSubO(t, "discard_reader", `
    range(10) | map(str.into) | stream | discard
    print("done")
  `, `done`)

	btesting.RunSubErr(t, "discard_invalid_arg", `
    discard(123)
  `, nil)
}

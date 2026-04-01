package global_test

import (
	"regexp"
	"testing"

	"github.com/dcaiafa/bag3l/internal/btesting"
)

// runSubO is like btesting.RunSubO but normalizes trailing whitespace per line
// before comparing. This is needed for functions like print_table whose
// tabwriter output includes trailing spaces.
func runSubO(t *testing.T, name string, prog string, expectedOutput string) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		t.Helper()
		output := btesting.Run(t, prog)
		expected := normalizeOutput(expectedOutput)
		actual := normalizeOutput(output)
		if actual != expected {
			t.Fatalf("Expected output:\n%v\nActual:\n%v", expected, actual)
		}
	})
}

var trailingSpaceRe = regexp.MustCompile(`(?m)(^\s+)|(\s+$)`)

func normalizeOutput(v string) string {
	return trailingSpaceRe.ReplaceAllString(v, "")
}

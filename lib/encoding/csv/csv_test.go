package csv_test

import (
	"testing"

	"github.com/dcaiafa/bag3l/internal/btesting"
)

func TestCsvDecode(t *testing.T) {
	btesting.RunSubO(t, "basic", `
		var data = buf.new()
		"a,b,c\n1,2,3\n" | data
		for row in data | csv.decode {
			print(row)
		}
	`, `
[a b c]
[1 2 3]
`)

	btesting.RunSubO(t, "columns", `
		var data = buf.new()
		"a,b,c\n1,2,3\n" | data
		for row in data | csv.decode({columns: [0, 2]}) {
			print(row)
		}
	`, `
[a c]
[1 3]
`)
}

func TestCsvEncode(t *testing.T) {
	btesting.RunSubO(t, "basic", `
		var data = buf.new()
		[["a", "b", "c"], ["1", "2", "3"]] | csv.encode(data)
		print(data)
	`, `
a,b,c
1,2,3
`)
}

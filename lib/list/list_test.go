package list_test

import (
	"testing"

	"github.com/dcaiafa/bag3l/internal/btesting"
)

func TestListAppend(t *testing.T) {
	btesting.RunSubO(t, "single", `[1, 2, 3] | list.append(4) | print`, `[1 2 3 4]`)
	btesting.RunSubO(t, "multiple", `[1] | list.append(2, 3, 4) | print`, `[1 2 3 4]`)
	btesting.RunSubO(t, "empty", `[] | list.append(1) | print`, `[1]`)
	btesting.RunSubO(t, "mixed_types", `[] | list.append(1, "two", true) | print`, `[1 two true]`)
}

func TestListAppendIter(t *testing.T) {
	btesting.RunSubO(t, "range", `[10] | list.append_iter(range(3)) | print`, `[10 0 1 2]`)
	btesting.RunSubO(t, "empty_iter", `[1, 2] | list.append_iter(range(0)) | print`, `[1 2]`)
}

func TestListFind(t *testing.T) {
	btesting.RunSubO(t, "found", `["a", 3, "c", 4] | list.find("c") | print`, `2`)
	btesting.RunSubO(t, "not_found", `["a", 3, "c", 4] | list.find(5) | print`, `<nil>`)
	btesting.RunSubO(t, "first_element", `[10, 20, 30] | list.find(10) | print`, `0`)
	btesting.RunSubO(t, "last_element", `[10, 20, 30] | list.find(30) | print`, `2`)
}

func TestListFromIter(t *testing.T) {
	btesting.RunSubO(t, "range", `list.from_iter(range(4)) | print`, `[0 1 2 3]`)
	btesting.RunSubO(t, "empty", `list.from_iter(range(0)) | print`, `[]`)
}

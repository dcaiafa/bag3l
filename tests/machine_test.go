package tests

import (
	"errors"
	"testing"

	"github.com/dcaiafa/bag3l/internal/vm"
	"github.com/stretchr/testify/require"
)

func TestErrorStack(t *testing.T) {
	const prog = `
	var x
	func h() {
		x()
	}
	func g() {
		h()
	}
	func f() {
		g()
	}
	f()
`

	_, err := run(prog, nil)
	require.Error(t, err)

	var rerr *vm.RuntimeError

	require.True(t, errors.As(err, &rerr))

	expectedStack := []vm.FrameInfo{
		{Filename: "main.b3", Line: 4, Func: "h"},
		{Filename: "main.b3", Line: 7, Func: "g"},
		{Filename: "main.b3", Line: 10, Func: "f"},
		{Filename: "main.b3", Line: 12, Func: "main"},
		{Filename: "", Line: 0, Func: "$main"},
	}

	require.Equal(t, expectedStack, rerr.Stack)
}

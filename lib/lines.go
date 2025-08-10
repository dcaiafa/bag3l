package lib

import (
	"bufio"
	"fmt"
	"io"

	nitro "github.com/dcaiafa/bag3l"
	"github.com/dcaiafa/bag3l/lib/core"
)

type linesOptions struct {
	MaxLineSize int64 `nitro:"maxlinesize"`
}

var linesOptionsConv core.Value2Structer

func lines(m *nitro.VM, args []nitro.Value, nRet int) ([]nitro.Value, error) {
	if len(args) != 1 && len(args) != 2 {
		return nil, errNotEnoughArgs
	}
	input, err := nitro.MakeReader(m, args[0])
	if err != nil {
		return nil, fmt.Errorf("invalid argument #1: %w", err)
	}

	var opt linesOptions
	if len(args) == 2 {
		optv, ok := args[1].(*nitro.Object)
		if !ok {
			return nil, fmt.Errorf("invalid options arg")
		}
		err := linesOptionsConv.Convert(optv, &opt)
		if err != nil {
			return nil, err
		}
	}

	l := &linesIter{
		input:   input,
		scanner: bufio.NewScanner(input),
	}

	if opt.MaxLineSize != 0 {
		l.scanner.Buffer(nil, int(opt.MaxLineSize))
	}

	outIter := nitro.NewIterator(l.Next, l.Close, 1)
	return []nitro.Value{outIter}, nil
}

type linesIter struct {
	input   io.Reader
	scanner *bufio.Scanner
}

func (l *linesIter) Next(m *nitro.VM, args []nitro.Value, nRet int) ([]nitro.Value, error) {
	if !l.scanner.Scan() {
		l.Close(m)
		if l.scanner.Err() != nil {
			return nil, l.scanner.Err()
		}
		return nil, nil
	}
	return []nitro.Value{
		nitro.NewString(l.scanner.Text()),
	}, nil
}

func (l *linesIter) Close(m *nitro.VM) error {
	core.CloseReader(l.input)
	return nil
}

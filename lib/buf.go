package lib

import (
	"bytes"
	"io"

	"github.com/dcaiafa/bag3l/internal/vm"
	"github.com/dcaiafa/bag3l/lib/core"
)

type Buffer struct {
	core.WriterBase
	buf *bytes.Buffer
}

func (b *Buffer) Read(buf []byte) (int, error) {
	return b.buf.Read(buf)
}

func (b *Buffer) Len() int {
	return b.buf.Len()
}

func (b *Buffer) String() string {
	return b.buf.String()
}

func newBuffer(data string) *Buffer {
	b := &Buffer{}
	if data != "" {
		b.buf = bytes.NewBufferString(data)
	} else {
		b.buf = new(bytes.Buffer)
	}
	b.WriterBase = core.NewWriterBase(b.buf)
	return b
}

func getBufferArg(args []vm.Value, ndx int) (*Buffer, error) {
	if ndx >= len(args) {
		return nil, errNotEnoughArgs
	}
	v, ok := args[ndx].(*Buffer)
	if !ok {
		return nil, errExpectedArg(ndx, args[ndx], "duration")
	}
	return v, nil
}

func bufNew(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	var err error
	if err = expectArgCount(args, 0, 1); err != nil {
		return nil, err
	}

	init := ""
	if len(args) == 1 {
		init, err = getStringArg(args, 0)
		if err != nil {
			return nil, err
		}
	}

	b := newBuffer(init)

	return []vm.Value{b}, nil
}

func bufRead(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	return read(m, args, nRet)
}

func bufReadByte(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if err := expectArgCount(args, 1, 1); err != nil {
		return nil, err
	}

	buf, err := getBufferArg(args, 0)
	if err != nil {
		return nil, err
	}

	b, err := buf.buf.ReadByte()
	if err != nil {
		if err == io.EOF {
			return []vm.Value{nil}, nil
		}
	}

	return []vm.Value{vm.NewInt(int64(b))}, nil
}

func bufReadRune(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if err := expectArgCount(args, 1, 1); err != nil {
		return nil, err
	}

	buf, err := getBufferArg(args, 0)
	if err != nil {
		return nil, err
	}

	r, l, err := buf.buf.ReadRune()
	if err != nil {
		if err == io.EOF {
			return []vm.Value{nil}, nil
		}
	}

	return []vm.Value{vm.NewInt(int64(r)), vm.NewInt(int64(l))}, nil
}

func bufReadFrom(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if err := expectArgCount(args, 1, 1); err != nil {
		return nil, err
	}

	buf, err := getBufferArg(args, 0)
	if err != nil {
		return nil, err
	}

	r, err := getReaderArg(m, args, 1)
	if err != nil {
		return nil, err
	}

	_, err = buf.buf.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func bufLen(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if err := expectArgCount(args, 1, 1); err != nil {
		return nil, err
	}

	buf, err := getBufferArg(args, 0)
	if err != nil {
		return nil, err
	}

	l := buf.Len()

	return []vm.Value{vm.NewInt(int64(l))}, nil
}

func bufCap(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if err := expectArgCount(args, 1, 1); err != nil {
		return nil, err
	}

	buf, err := getBufferArg(args, 0)
	if err != nil {
		return nil, err
	}

	c := buf.buf.Cap()

	return []vm.Value{vm.NewInt(int64(c))}, nil
}

func bufUnreadByte(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if err := expectArgCount(args, 1, 1); err != nil {
		return nil, err
	}

	buf, err := getBufferArg(args, 0)
	if err != nil {
		return nil, err
	}

	err = buf.buf.UnreadByte()
	if err != nil {
		return nil, err
	}

	return nil, nil
}

package lib

import (
	"io"

	"github.com/dcaiafa/bag3l/internal/bytequeue"
	"github.com/dcaiafa/bag3l/internal/vm"
)

func stream(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if len(args) != 1 {
		return nil, errInvalidNumberOfArgs
	}

	iter, err := getIterArg(m, args, 0)
	if err != nil {
		return nil, err
	}

	reader := &iterReader{
		m: m,
		e: iter,
	}

	return []vm.Value{reader}, nil
}

type iterReader struct {
	m   *vm.VM
	e   vm.Iterator
	buf bytequeue.ByteQueue
}

func newIterReader(m *vm.VM, iter vm.Iterator) *iterReader {
	return &iterReader{
		m: m,
		e: iter,
	}
}

func (r *iterReader) String() string    { return "<reader>" }
func (r *iterReader) Type() string      { return "reader" }
func (r *iterReader) Traits() vm.Traits { return vm.TraitNone }

func (r *iterReader) Read(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}

	if len(r.buf.Peek()) < len(b) && !r.e.IsClosed() {
		v, err := r.m.IterNext(r.e, 1)
		if err != nil {
			return 0, err
		}
		if v != nil {
			r.buf.Write([]byte(v[0].String()))
			r.buf.Write([]byte{'\n'})
		} else {
			err := r.m.IterClose(r.e)
			if err != nil {
				return 0, err
			}
		}
	}

	if len(r.buf.Peek()) == 0 {
		return 0, io.EOF
	}

	n := len(r.buf.Peek())
	if n > len(b) {
		n = len(b)
	}

	copy(b, r.buf.Peek()[:n])
	r.buf.Pop(n)

	return n, nil
}

func (r *iterReader) Close() error {
	return r.m.IterClose(r.e)
}

package lib

import (
	"io"

	nitro "github.com/dcaiafa/bag3l"
	"github.com/dcaiafa/bag3l/internal/bytequeue"
	"github.com/dcaiafa/bag3l/internal/vm"
)

func streamLines(m *nitro.VM, args []nitro.Value, nRet int) ([]nitro.Value, error) {
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

	return []nitro.Value{reader}, nil
}

type iterReader struct {
	m   *vm.VM
	e   vm.Iterator
	buf bytequeue.ByteQueue
}

func newIterReader(vm *vm.VM, iter vm.Iterator) *iterReader {
	return &iterReader{
		m: vm,
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

	if !r.e.IsClosed() {
		for len(r.buf.Peek()) < len(b) {
			v, err := r.m.IterNext(r.e, 1)
			if err != nil {
				return 0, err
			}
			if v == nil {
				err := r.m.IterClose(r.e)
				if err != nil {
					return 0, err
				}
				break
			}
			r.buf.Write([]byte(v[0].String()))
			r.buf.Write([]byte{'\n'})
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

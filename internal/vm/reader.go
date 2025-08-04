package vm

import (
	"fmt"
	"io"
)

type Readable interface {
	Value
	MakeReader() Reader
}

func IsReadable(v Value) bool {
	if v == nil {
		return true
	}
	switch v.(type) {
	case Readable, Reader, Iterator:
		return true
	default:
		return false
	}
}

func MakeReader(vm *VM, v Value) (Reader, error) {
	if v == nil {
		return emptyReader, nil
	}
	switch v := v.(type) {
	case Reader:
		return v, nil
	case Readable:
		return v.MakeReader(), nil
	default:
		return nil, fmt.Errorf("type %q %w", TypeName(v), ErrIsNotReadable)
	}
}

type Reader interface {
	Value
	io.Reader
}

type emptyReaderImpl struct{}

var emptyReader = &emptyReaderImpl{}

func (r *emptyReaderImpl) String() string { return "<reader>" }
func (r *emptyReaderImpl) Type() string   { return "reader" }
func (r *emptyReaderImpl) Traits() Traits { return TraitNone }

func (r *emptyReaderImpl) MakeReader() (Reader, error) {
	return r, nil
}

func (r *emptyReaderImpl) Read(b []byte) (int, error) {
	return 0, io.EOF
}

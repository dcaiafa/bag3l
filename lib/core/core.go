package core

import (
	"errors"
	"fmt"
	"io"

	"github.com/dcaiafa/bag3l/internal/vm"
)

type NativeReader interface {
	GetNativeReader() io.Reader
}

type NativeWriter interface {
	GetNativeWriter() io.Writer
}

type WriterBase struct {
	io.Writer
	typ string
}

func NewWriterBase(w io.Writer) WriterBase {
	return WriterBase{Writer: w}
}

func (b *WriterBase) String() string    { return "<Writer>" }
func (b *WriterBase) Type() string      { return "Writer" }
func (b *WriterBase) Traits() vm.Traits { return vm.TraitNone }

var ErrWriterCallUsage = errors.New(
	`invalid usage. Expected <writer>(reader)`)

func (b *WriterBase) Call(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if len(args) != 1 {
		return nil, ErrWriterCallUsage
	}

	reader, err := vm.MakeReader(m, args[0])
	if err != nil {
		return nil, ErrWriterCallUsage
	}

	n, err := io.Copy(b.Writer, reader)
	if err != nil {
		return nil, err
	}

	CloseReader(reader)

	return []vm.Value{vm.NewInt(n)}, nil
}

func CloseReader(r io.Reader) {
	if c, ok := r.(io.Closer); ok {
		c.Close()
	}
}

var ErrNotEnoughArgs = errors.New("not enough arguments")
var ErrTooManyArgs = errors.New("too many arguments")

func CheckArgCount(args []vm.Value, min, max int) error {
	if len(args) < min {
		return ErrNotEnoughArgs
	} else if len(args) > max {
		return ErrTooManyArgs
	}
	return nil
}

func GetString(args []vm.Value, n int) (vm.String, error) {
	arg, ok := args[n].(vm.String)
	if !ok {
		return vm.String{}, fmt.Errorf(
			"arg %v: expected string, got %v",
			n, vm.TypeName(args[n]))
	}
	return arg, nil
}

func GetInt(args []vm.Value, n int) (vm.Int, error) {
	arg, ok := args[n].(vm.Int)
	if !ok {
		return vm.Int{}, fmt.Errorf(
			"arg %v: expected int, got %v",
			n, vm.TypeName(args[n]))
	}
	return arg, nil
}

func GetFloat(args []vm.Value, n int) (vm.Float, error) {
	arg, ok := args[n].(vm.Float)
	if !ok {
		return vm.Float{}, fmt.Errorf(
			"arg %v: expected float, got %v",
			n, vm.TypeName(args[n]))
	}
	return arg, nil
}

func GetReader(m *vm.VM, args []vm.Value, n int) (vm.Reader, error) {
	arg, err := vm.MakeReader(m, args[n])
	if err != nil {
		return nil, fmt.Errorf("arg %v: %w", n, err)
	}
	return arg, nil
}

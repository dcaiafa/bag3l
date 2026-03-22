package io

import (
	"io"
	"os"

	nitro "github.com/dcaiafa/bag3l"
	"github.com/dcaiafa/bag3l/internal/vm"
	"github.com/dcaiafa/bag3l/lib/core"
)

//go:generate go run ../../internal/stub/stubgen io.stubgen

type stdinWrapper struct {
	*os.File
}

func (w *stdinWrapper) String() string             { return "<Reader>" }
func (w *stdinWrapper) Type() string               { return "Reader" }
func (w *stdinWrapper) Traits() vm.Traits          { return vm.TraitNone }
func (w *stdinWrapper) GetNativeReader() io.Reader { return w.File }

var DefaultStdin = &stdinWrapper{
	File: os.Stdin,
}

type stdoutWrapper struct {
	*os.File
}

func (w *stdoutWrapper) String() string             { return "<Writer>" }
func (w *stdoutWrapper) Type() string               { return "Writer" }
func (w *stdoutWrapper) Traits() vm.Traits          { return vm.TraitNone }
func (w *stdoutWrapper) GetNativeWriter() io.Writer { return w.File }

var DefaultOut = &stdoutWrapper{
	File: os.Stdout,
}

var DefaultErr = &stdoutWrapper{
	File: os.Stderr,
}

func wrapWriter(w io.Writer) vm.Writer {
	v, ok := w.(vm.Writer)
	if ok {
		return v
	}
	wb := core.NewWriterBase(w)
	return &wb
}

type stdoutUserDataKey struct{}

type stdoutStack struct {
	def   vm.Writer
	stack []vm.Writer
}

func Stdout(m *nitro.VM) vm.Writer {
	stdout, ok := m.GetUserData(stdoutUserDataKey{}).(vm.Writer)
	if !ok {
		panic("stdout is not set")
	}
	return stdout
}

func SetStdout(m *vm.VM, w io.Writer) {
	m.SetUserData(stdoutUserDataKey{}, wrapWriter(w))
}

func Stderr(m *nitro.VM) vm.Writer {
	return DefaultErr
}

func in0(vm *vm.VM) (vm.Reader, error) {
	return DefaultStdin, nil
}

func out0(vm *nitro.VM, r vm.Reader) (vm.Writer, error) {
	out := Stdout(vm)
	if r == nil {
		return out, nil
	}
	_, err := io.Copy(out, r)
	if err != nil {
		return nil, err
	}
	core.CloseReader(r)
	return out, nil
}

func err0(vm *nitro.VM, r vm.Reader) (vm.Writer, error) {
	out := Stderr(vm)
	if r == nil {
		return out, nil
	}
	_, err := io.Copy(out, r)
	if err != nil {
		return nil, err
	}
	core.CloseReader(r)
	return out, nil
}

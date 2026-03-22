package io

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/dcaiafa/bag3l/internal/vm"
	"github.com/dcaiafa/bag3l/lib/core"
)

// from_crlf(s Str) Str
func from_crlf0(m *vm.VM, s string) (string, error) {
	converter := newFromCRLFReaderWriter(strings.NewReader(s))
	converted, err := io.ReadAll(converter)
	if err != nil {
		return "", err
	}
	return string(converted), nil
}

// from_crlf(r Reader) Reader
func from_crlf1(m *vm.VM, r vm.Reader) (vm.Reader, error) {
	return newFromCRLFReaderWriter(r), nil
}

// from_crlf(w Writer) Writer
func from_crlf2(m *vm.VM, w vm.Writer) (vm.Writer, error) {
	return newFromCRLFReaderWriter(w), nil
}

// to_crlf(s Str, convert Bool) Str
func to_crlf0(m *vm.VM, s string, convert bool) (string, error) {
	if !convert {
		return s, nil
	}
	converter := newToCRLFReader(strings.NewReader(s))
	converted, err := io.ReadAll(converter)
	if err != nil {
		return "", err
	}
	return string(converted), nil
}

// to_crlf(r Reader, convert Bool) Reader
func to_crlf1(m *vm.VM, r vm.Reader, convert bool) (vm.Reader, error) {
	if !convert {
		return r, nil
	}
	return newToCRLFReader(r), nil
}

// fromCRLFReaderWriter strips carriage returns from reads and writes.
type fromCRLFReaderWriter struct {
	c   io.Closer
	r   io.Reader
	w   io.Writer
	buf *bufio.Reader
}

func (c *fromCRLFReaderWriter) String() string    { return "<from_crlf>" }
func (c *fromCRLFReaderWriter) Type() string      { return "from_crlf" }
func (c *fromCRLFReaderWriter) Traits() vm.Traits { return vm.TraitNone }

func newFromCRLFReaderWriter(v any) *fromCRLFReaderWriter {
	c := &fromCRLFReaderWriter{}
	c.c, _ = v.(io.Closer)
	c.r, _ = v.(io.Reader)
	c.w, _ = v.(io.Writer)
	if c.r != nil {
		var ok bool
		c.buf, ok = c.r.(*bufio.Reader)
		if !ok {
			c.buf = bufio.NewReader(c.r)
		}
	}
	return c
}

func (c *fromCRLFReaderWriter) Read(buf []byte) (int, error) {
	if c.r == nil {
		return 0, fmt.Errorf("value not readable")
	}
	n := 0
	for n < len(buf) {
		b, err := c.buf.ReadByte()
		if err != nil {
			return n, err
		}
		if b != '\r' {
			buf[n] = b
			n++
		}
	}
	return n, nil
}

func (c *fromCRLFReaderWriter) Write(buf []byte) (int, error) {
	if c.w == nil {
		return 0, fmt.Errorf("value not writable")
	}
	n := 0
	for len(buf) > 0 {
		if buf[0] == '\r' {
			buf = buf[1:]
			n++
		}
		i := bytes.IndexByte(buf, '\r')
		if i == -1 {
			i = len(buf)
		}
		if i > 0 {
			pn, err := c.w.Write(buf[:i])
			n += pn
			if err != nil {
				return n, err
			}
			buf = buf[i:]
		}
	}
	return n, nil
}

func (c *fromCRLFReaderWriter) Call(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if len(args) != 1 {
		return nil, core.ErrWriterCallUsage
	}
	reader, err := vm.MakeReader(m, args[0])
	if err != nil {
		return nil, core.ErrWriterCallUsage
	}
	n, err := io.Copy(c, reader)
	if err != nil {
		return nil, err
	}
	return []vm.Value{vm.NewInt(n)}, nil
}

func (c *fromCRLFReaderWriter) Close() error {
	if c.c != nil {
		return c.c.Close()
	}
	return nil
}

// toCRLFReader converts LF to CRLF during reads.
type toCRLFReader struct {
	r   io.Reader
	buf *bufio.Reader
	cr  bool
}

func newToCRLFReader(r io.Reader) *toCRLFReader {
	c := &toCRLFReader{
		r: r,
	}
	var ok bool
	c.buf, ok = r.(*bufio.Reader)
	if !ok {
		c.buf = bufio.NewReader(r)
	}
	return c
}

func (c *toCRLFReader) String() string    { return "<to_crlf_reader>" }
func (c *toCRLFReader) Type() string      { return "to_crlf_reader" }
func (c *toCRLFReader) Traits() vm.Traits { return vm.TraitNone }

func (c *toCRLFReader) Read(buf []byte) (int, error) {
	n := 0
	for n < len(buf) {
		b, err := c.readByte()
		if err != nil {
			return n, err
		}
		buf[n] = b
		n++
	}
	return n, nil
}

func (c *toCRLFReader) Close() error {
	if closer, ok := c.r.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func (c *toCRLFReader) readByte() (byte, error) {
	b, err := c.buf.ReadByte()
	if err != nil {
		return 0, err
	}
	if b == '\n' && !c.cr {
		c.cr = true
		c.buf.UnreadByte()
		return '\r', nil
	}
	c.cr = (b == '\r')
	return b, nil
}

package base64

import (
	"encoding/base64"
	"errors"
	"io"

	"github.com/dcaiafa/bag3l/internal/vm"
)

//go:generate go run ../../../internal/stub/stubgen base64.stubgen

func decode0(m *vm.VM, v string, opts *Options) (string, error) {
	enc, err := newEncoding(opts)
	if err != nil {
		return "", err
	}
	dec, err := enc.DecodeString(v)
	if err != nil {
		return "", err
	}
	return string(dec), nil
}

func decode1(m *vm.VM, r vm.Reader, opts *Options) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	return decode0(m, string(data), opts)
}

func encode0(m *vm.VM, v string, opts *Options) (string, error) {
	enc, err := newEncoding(opts)
	if err != nil {
		return "", err
	}

	encData := enc.EncodeToString([]byte(v))
	if err != nil {
		return "", err
	}

	return encData, nil
}

func encode1(m *vm.VM, r vm.Reader, opts *Options) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	return encode0(m, string(data), opts)
}

var defaultOptions = &Options{
	Mode:   "std",
	Strict: false,
}

func newEncoding(opts *Options) (*base64.Encoding, error) {
	if opts == nil {
		opts = defaultOptions
	}

	var enc *base64.Encoding
	switch opts.Mode {
	case "std":
		enc = base64.StdEncoding
	case "url":
		enc = base64.URLEncoding
	case "rawstd":
		enc = base64.RawStdEncoding
	case "rawurl":
		enc = base64.RawURLEncoding
	default:
		return nil, errors.New("invalid mode")
	}

	if opts.Strict {
		enc = enc.Strict()
	}

	return enc, nil
}

package lib

import (
	cryptosha1 "crypto/sha1"
	"errors"
	"io"

	"github.com/dcaiafa/bag3l/internal/vm"
)

var errSha1Usage = errors.New(
	`invalid usage. Expected sha1(reader)`)

func sha1(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if len(args) != 1 {
		return nil, errSha1Usage
	}

	input, err := vm.MakeReader(m, args[0])
	if err != nil {
		return nil, errSha1Usage
	}

	h := cryptosha1.New()
	_, err = io.Copy(h, input)
	if err != nil {
		return nil, err
	}

	return []vm.Value{vm.NewString(string(h.Sum(nil)))}, nil
}

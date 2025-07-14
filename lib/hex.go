package lib

import (
	"encoding/hex"
	"errors"

	"github.com/dcaiafa/bag3l"
)

var errToHexUsage = errors.New(
	`invalid usage. Expected to_hex(string)`)

func toHex(vm *bag3l.VM, args []bag3l.Value, nRet int) ([]bag3l.Value, error) {
	if len(args) != 1 {
		return nil, errToHexUsage
	}

	input, ok := args[0].(bag3l.String)
	if !ok {
		return nil, errSha1Usage
	}

	res := hex.EncodeToString([]byte(input.String()))
	return []bag3l.Value{bag3l.NewString(res)}, nil
}

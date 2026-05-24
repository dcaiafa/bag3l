package lib

import (
	"errors"
	"regexp"

	"github.com/dcaiafa/bag3l/internal/vm"
)

var errRegexUsage = errors.New(
	`invalid usage. Expected regex(string)`)

func regex(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if len(args) != 1 {
		return nil, errRegexUsage
	}

	reStr, ok := args[0].(vm.String)
	if !ok {
		return nil, errRegexUsage
	}

	re, err := regexp.Compile(reStr.String())
	if err != nil {
		return nil, errRegexUsage
	}

	return []vm.Value{vm.NewRegex(re)}, nil
}

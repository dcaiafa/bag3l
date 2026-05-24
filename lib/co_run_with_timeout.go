package lib

import (
	"context"

	"github.com/dcaiafa/bag3l/internal/vm"
)

func runWithTimeout(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if len(args) > 2 {
		return nil, errTooManyArgs
	}

	dur, err := getDurationArg(args, 1)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(m.Context(), dur.Duration())
	defer cancel()

	m.PushContext(ctx)
	defer m.PopContext()

	_, err = m.Call(args[0], nil, 0)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

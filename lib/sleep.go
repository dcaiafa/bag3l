package lib

import (
	"context"
	"errors"
	"time"

	"github.com/dcaiafa/bag3l/internal/vm"
	btime "github.com/dcaiafa/bag3l/lib/time"
)

var errSleepUsage = errors.New(
	`invalid usage. Expected sleep(duration|int)`)

func sleep(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if len(args) != 1 {
		return nil, errSleepUsage
	}

	var dur time.Duration

	if durArg, ok := args[0].(btime.Duration); ok {
		dur = durArg.Duration()
	} else if intArg, ok := args[0].(vm.Int); ok {
		dur = time.Duration(intArg.Int64()) * time.Second
	} else {
		return nil, errSleepUsage
	}

	timer := time.NewTimer(dur)
	defer timer.Stop()

	var err error
	m.Block(func(ctx context.Context) {
		select {
		case <-timer.C:
		case <-ctx.Done():
			err = ctx.Err()
		}
	})

	return nil, err
}

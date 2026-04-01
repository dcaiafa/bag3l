package global

import (
	"fmt"
	"io"

	"github.com/dcaiafa/bag3l/internal/vm"
	libio "github.com/dcaiafa/bag3l/lib/io"
)

func print0(m *vm.VM, args []vm.Value) error {
	return basePrint(libio.Stdout(m), m, args)
}

func basePrint(out io.Writer, m *vm.VM, args []vm.Value) error {
	fprint := fmt.Fprintln

	if len(args) == 1 {
		switch arg := args[0].(type) {
		case vm.Iterator:
			for {
				nret := arg.IterNRet()
				v, err := m.IterNext(arg, nret)
				if err != nil {
					return err
				}
				if v == nil {
					break
				}
				iargs := valuesToInterface(v)
				_, err = fprint(out, iargs...)
				if err != nil {
					return err
				}
			}
			return nil

		case io.Reader:
			_, err := io.Copy(out, arg)
			if err != nil {
				return err
			}
			fmt.Fprintln(out, "")
			return nil
		}
	}

	iargs := valuesToInterface(args)
	fprint(out, iargs...)

	return nil
}

func valuesToInterface(values []vm.Value) []any {
	ivalues := make([]any, len(values))
	for i, v := range values {
		switch v := v.(type) {
		case vm.Int:
			ivalues[i] = v.Int64()
		case vm.Float:
			ivalues[i] = v.Float64()
		case vm.String:
			ivalues[i] = v.String()
		default:
			ivalues[i] = v
		}
	}
	return ivalues
}

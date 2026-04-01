package global

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/dcaiafa/bag3l/internal/vm"
	libio "github.com/dcaiafa/bag3l/lib/io"
)

func print_table0(m *vm.VM, iter vm.Iterator, opts *PrintTableOptions) error {
	defer m.IterClose(iter)

	var alignright bool
	var minwidth int64
	var padding int64 = 1
	var padchar int64 = ' '

	if opts != nil {
		alignright = opts.Alignright
		minwidth = opts.Minwidth
		if opts.Padding != 0 {
			padding = opts.Padding
		}
		if opts.Padchar != 0 {
			padchar = opts.Padchar
		}
	}

	var flags uint
	if alignright {
		flags = flags | tabwriter.AlignRight
	}

	tabw := tabwriter.NewWriter(
		libio.Stdout(m),
		int(minwidth),
		0, /*tabwidth*/
		int(padding),
		byte(padchar),
		flags)

	defer tabw.Flush()

	buf := bytes.Buffer{}
	writeRecord := func(vs []vm.Value) error {
		buf.Reset()
		for _, v := range vs {
			buf.WriteString(v.String())
			buf.WriteByte('\t')
		}
		buf.WriteByte('\n')
		_, err := tabw.Write(buf.Bytes())
		return err
	}

	first := true
	var headers []vm.Value
	var values []vm.Value
	for {
		vs, err := m.IterNext(iter, 1)
		if err != nil {
			return err
		}
		if vs == nil {
			break
		}

		rec := vs[0]

		if mrec, ok := rec.(*vm.Map); ok {
			if first {
				headers = make([]vm.Value, 0, mrec.Len())
				values = make([]vm.Value, mrec.Len())
				err = nil
				mrec.ForEach(func(k, v vm.Value) bool {
					kstr, ok := k.(vm.String)
					if !ok {
						err = fmt.Errorf(
							"only string keys are supported, but map has key type %v",
							vm.TypeName(k))
						return false
					}
					headers = append(headers, kstr)
					return true
				})
				if err != nil {
					return err
				}
				writeRecord(headers)
			}
			for i, k := range headers {
				v, ok := mrec.Get(k)
				if ok {
					if vstr, ok := v.(vm.String); ok {
						values[i] = vstr
					} else if str, ok := v.(fmt.Stringer); ok {
						values[i] = vm.NewString(str.String())
					} else if v == nil {
						values[i] = vm.NewString("<nil>")
					} else {
						values[i] = vm.NewString("")
					}
				} else {
					values[i] = vm.NewString("")
				}
			}
			writeRecord(values)
		} else if _, ok := rec.(*vm.List); ok {
			panic("not impl")
		} else {
			return fmt.Errorf(
				"expected iterator values of array or map; received %v",
				vm.TypeName(rec))
		}

		first = false
	}

	err := tabw.Flush()
	if err != nil {
		return err
	}

	return nil
}

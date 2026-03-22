package lib

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"text/tabwriter"

	nitro "github.com/dcaiafa/bag3l"
	"github.com/dcaiafa/bag3l/lib/core"
	libio "github.com/dcaiafa/bag3l/lib/io"
)

func print(vm *nitro.VM, args []nitro.Value, nRet int) ([]nitro.Value, error) {
	return basePrint(libio.Stdout(vm), vm, args, nRet)
}

func basePrint(out io.Writer, m *nitro.VM, args []nitro.Value, nRet int) ([]nitro.Value, error) {
	fprint := fmt.Fprintln

	if len(args) == 1 {
		switch arg := args[0].(type) {
		case nitro.Iterator:
			for {
				nret := arg.IterNRet()
				v, err := m.IterNext(arg, nret)
				if err != nil {
					return nil, err
				}
				if v == nil {
					break
				}
				iargs := valuesToInterface(v)
				_, err = fprint(out, iargs...)
				if err != nil {
					return nil, err
				}
			}
			return nil, nil

		case io.Reader:
			_, err := io.Copy(out, arg)
			if err != nil {
				return nil, err
			}
			fmt.Fprintln(out, "")

			return nil, nil
		}
	}

	iargs := valuesToInterface(args)
	fprint(out, iargs...)

	return nil, nil
}

type printTableOptions struct {
	AlignRight bool  `nitro:"alignright"`
	MinWidth   int64 `nitro:"minwidth"`
	Padding    int64 `nitro:"padding"`
	PadChar    int64 `nitro:"padchar"`
}

var printTableConv core.Value2Structer

var errPrintTableUsage = errors.New(
	`invalid usage. Expected print_table(iter, map?)`)

func printTable(m *nitro.VM, args []nitro.Value, nRet int) ([]nitro.Value, error) {
	if len(args) != 1 && len(args) != 2 {
		return nil, errPrintTableUsage
	}

	inIter, err := nitro.MakeIterator(m, args[0])
	if err != nil {
		return nil, errPrintTableUsage
	}
	defer m.IterClose(inIter)

	opts := printTableOptions{
		Padding: -1,
	}

	if len(args) == 2 {
		err = printTableConv.Convert(args[1], &opts)
		if err != nil {
			return nil, err
		}
	}

	var flags uint
	if opts.AlignRight {
		flags = flags | tabwriter.AlignRight
	}

	if opts.Padding == -1 {
		opts.Padding = 1
	}
	if opts.PadChar == 0 {
		opts.PadChar = ' '
	}

	tabw := tabwriter.NewWriter(
		libio.Stdout(m),
		int(opts.MinWidth),
		0, /*tabwidth*/
		int(opts.Padding),
		byte(opts.PadChar),
		flags)

	defer tabw.Flush()

	buf := bytes.Buffer{}
	writeRecord := func(vs []nitro.Value) error {
		buf.Reset()
		for _, v := range vs {
			// TODO: replace \t in value.
			buf.WriteString(v.String())
			buf.WriteByte('\t')
		}
		buf.WriteByte('\n')
		_, err := tabw.Write(buf.Bytes())
		return err
	}

	first := true
	var headers []nitro.Value
	var values []nitro.Value
	for {
		vs, err := m.IterNext(inIter, 1)
		if err != nil {
			return nil, err
		}
		if vs == nil {
			break
		}

		rec := vs[0]

		if mrec, ok := rec.(*nitro.Object); ok {
			if first {
				headers = make([]nitro.Value, 0, mrec.Len())
				values = make([]nitro.Value, mrec.Len())
				err = nil
				mrec.ForEach(func(k, v nitro.Value) bool {
					kstr, ok := k.(nitro.String)
					if !ok {
						err = fmt.Errorf(
							"only string keys are supported, but map has key type %v",
							nitro.TypeName(k))
						return false
					}
					headers = append(headers, kstr)
					return true
				})
				writeRecord(headers)
			}
			for i, k := range headers {
				v, ok := mrec.Get(k)
				if ok {
					if vstr, ok := v.(nitro.String); ok {
						values[i] = vstr
					} else if str, ok := v.(fmt.Stringer); ok {
						values[i] = nitro.NewString(str.String())
					} else if v == nil {
						values[i] = nitro.NewString("<nil>")
					} else {
						values[i] = nitro.NewString("")
					}
				} else {
					values[i] = nitro.NewString("")
				}
			}
			writeRecord(values)
		} else if _, ok := rec.(*nitro.Array); ok {
			panic("not impl")
		} else {
			return nil, fmt.Errorf(
				"expected iterator values of array or map; received %v",
				nitro.TypeName(rec))
		}

		first = false
	}

	err = tabw.Flush()
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func valuesToInterface(values []nitro.Value) []interface{} {
	ivalues := make([]interface{}, len(values))
	for i, v := range values {
		switch v := v.(type) {
		case nitro.Int:
			ivalues[i] = v.Int64()
		case nitro.Float:
			ivalues[i] = v.Float64()
		case nitro.String:
			ivalues[i] = v.String()
		default:
			ivalues[i] = v
		}
	}
	return ivalues
}

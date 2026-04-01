package csv

import (
	"encoding/csv"
	"fmt"
	"io"

	nitro "github.com/dcaiafa/bag3l"
	"github.com/dcaiafa/bag3l/internal/vm"
	"github.com/dcaiafa/bag3l/lib/core"
	libio "github.com/dcaiafa/bag3l/lib/io"
)

//go:generate go run ../../../internal/stub/stubgen csv.stubgen

func decode0(m *vm.VM, input vm.Reader, opts *DecodeOptions) (vm.Iterator, error) {
	var columns []int
	if opts != nil && opts.Columns != nil {
		columns = make([]int, opts.Columns.Len())
		for i := 0; i < opts.Columns.Len(); i++ {
			intVal, ok := opts.Columns.Get(i).(vm.Int)
			if !ok {
				return nil, fmt.Errorf("columns must contain integers")
			}
			columns[i] = int(intVal.Int64())
		}
	}

	i := &csvIter{
		origReader: input,
		csvReader:  csv.NewReader(input),
		columns:    columns,
	}
	i.csvReader.ReuseRecord = true

	return vm.NewIterator(i.next, i.close, 1), nil
}

type csvIter struct {
	origReader io.Reader
	csvReader  *csv.Reader
	columns    []int
}

func (i *csvIter) next(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	record, err := i.csvReader.Read()
	if err != nil {
		core.CloseReader(i.origReader)
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}

	var res []vm.Value
	if i.columns != nil {
		res = make([]vm.Value, len(i.columns))
		for j, col := range i.columns {
			if col < 0 || col >= len(record) {
				return nil, fmt.Errorf(
					"column %v requested, but csv record has %v columns",
					col, len(record))
			}
			res[j] = vm.NewString(record[col])
		}
	} else {
		res = make([]vm.Value, len(record))
		for j, entry := range record {
			res[j] = vm.NewString(entry)
		}
	}

	return []vm.Value{vm.NewListWithSlice(res)}, nil
}

func (i *csvIter) close(vm *vm.VM) error {
	core.CloseReader(i.origReader)
	return nil
}

func encode0(m *nitro.VM, iter vm.Iterator, out vm.Writer) error {
	var w io.Writer
	if out == nil {
		w = libio.Stdout(m)
	} else {
		w = out
	}

	var records []string
	csvOut := csv.NewWriter(w)

	for {
		vals, err := m.IterNext(iter, 1)
		if err != nil {
			return err
		}
		if vals == nil {
			break
		}

		recordValues, ok := vals[0].(*vm.List)
		if !ok {
			return fmt.Errorf(
				"iterator must return lists, but returned %v",
				vm.TypeName(vals[0]))
		}

		if records == nil {
			records = make([]string, recordValues.Len())
		} else if len(records) != recordValues.Len() {
			return fmt.Errorf("all csv rows must have the same number of records")
		}

		for j := 0; j < recordValues.Len(); j++ {
			records[j] = recordValues.Get(j).String()
		}

		err = csvOut.Write(records)
		if err != nil {
			return err
		}
	}

	csvOut.Flush()

	return csvOut.Error()
}

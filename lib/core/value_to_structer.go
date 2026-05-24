package core

import (
	"fmt"
	golog "log"
	"reflect"
	"sync"

	"github.com/dcaiafa/bag3l/internal/vm"
)

var typeValue = reflect.Indirect(reflect.ValueOf(new(vm.Value))).Type()

type fieldMapping struct {
	fieldIndex int
	typ        reflect.Type
}

type Value2Structer struct {
	once   sync.Once
	fields map[string]*fieldMapping
}

func (s *Value2Structer) Convert(from vm.Value, to interface{}) error {
	s.once.Do(func() { s.init(to) })

	rv := reflect.Indirect(reflect.ValueOf(to))
	fromObj, ok := from.(*vm.Map)
	if !ok {
		return fmt.Errorf("expected object; received %v", vm.TypeName(from))
	}

	var err error
	fromObj.ForEach(func(k, v vm.Value) bool {
		kstr, ok := k.(vm.String)
		if !ok {
			err = fmt.Errorf("invalid option %v", k.String())
			return false
		}
		field := s.fields[kstr.String()]
		if field == nil {
			err = fmt.Errorf("invalid option %v", kstr.String())
			return false
		}

		rfield := rv.Field(field.fieldIndex)

		if rfield.Type().Kind() == reflect.Ptr {
			zeroValue := reflect.New(rfield.Type().Elem())
			rfield.Set(zeroValue)
			rfield = reflect.Indirect(rfield)
		}

		switch rfield.Type().Kind() {
		case reflect.Bool:
			vb, ok := v.(vm.Bool)
			if !ok {
				err = fmt.Errorf("option %v expected bool; received %v",
					kstr.String(), vm.TypeName(v))
				return false
			}
			rfield.SetBool(vb.Bool())

		case reflect.String:
			vstr, ok := v.(vm.String)
			if !ok {
				err = fmt.Errorf("option %v expected string; received %v",
					kstr.String(), vm.TypeName(v))
				return false
			}
			rfield.SetString(vstr.String())

		case reflect.Int64, reflect.Int:
			vint, ok := v.(vm.Int)
			if !ok {
				err = fmt.Errorf("option %v expected int; received %v",
					kstr.String(), vm.TypeName(v))
				return false
			}
			rfield.SetInt(vint.Int64())

		case reflect.Slice:
			va, ok := v.(*vm.List)
			if !ok {
				err = fmt.Errorf("option %v expected array; received %v",
					kstr.String(), vm.TypeName(v))
				return false
			}

			var newSlice reflect.Value
			for i := 0; i < va.Len(); i++ {
				entry := va.Get(i)

				switch rfield.Type().Elem().Kind() {
				case reflect.String:
					entryStr, ok := entry.(vm.String)
					if !ok {
						err = fmt.Errorf(
							"option %v expected array of string; but entry %d was %v",
							kstr.String(), i, vm.TypeName(entry))
						return false
					}
					newSlice = reflect.Append(rfield, reflect.ValueOf(entryStr.String()))

				case reflect.Int64:
					entryInt, ok := entry.(vm.Int)
					if !ok {
						err = fmt.Errorf(
							"option %v expected array of int; but entry %d was %v",
							kstr.String(), i, vm.TypeName(entry))
						return false
					}
					newSlice = reflect.Append(rfield, reflect.ValueOf(entryInt.Int64()))

				default:
					golog.Panicf("Unsupported option array elem type %v", rfield.Elem().Type())
				}
				rfield.Set(newSlice)
			}

		case reflect.Interface:
			rfield.Set(reflect.ValueOf(v))

		default:
			panic("unreachable")
		}

		return true
	})

	return err
}

func (s *Value2Structer) init(o interface{}) {
	typ := reflect.Indirect(reflect.ValueOf(o)).Type()

	if typ.Kind() != reflect.Struct {
		panic("value is not a struct")
	}

	s.fields = make(map[string]*fieldMapping)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		nitroName := field.Tag.Get("nitro")
		if nitroName == "" {
			continue
		}

		s.fields[nitroName] = &fieldMapping{
			fieldIndex: i,
			typ:        field.Type,
		}
	}
}

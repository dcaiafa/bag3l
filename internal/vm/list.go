package vm

import (
	"fmt"
	"math"
)

type List struct {
	list []Value
}

func NewList() *List {
	return &List{}
}

func NewListWithSlice(s []Value) *List {
	return &List{list: s}
}

func (a *List) Type() string   { return "list" }
func (a *List) Traits() Traits { return TraitNone }

func (a *List) Add(v Value) {
	a.list = append(a.list, v)
}

func (a *List) Get(index int) Value {
	if index >= len(a.list) {
		return nil
	}
	return a.list[index]
}

func (a *List) Put(index int, v Value) {
	a.list[index] = v
}

func (a *List) Find(v Value) int {
	for i, entry := range a.list {
		if entry == v {
			return i
		}
	}
	return -1
}

func (a *List) Index(key Value) (Value, error) {
	switch key := key.(type) {
	case Int:
		idx := int(key.Int64())
		if idx < 0 {
			idx = len(a.list) + idx
		}
		if idx < 0 || idx >= len(a.list) {
			return nil, nil
		}
		return a.list[idx], nil

	default:
		return nil, fmt.Errorf(
			"cannot index string using key type %v",
			TypeName(key))
	}
}

func (a *List) IndexRef(key Value) (ValueRef, error) {
	index, ok := key.(Int)
	if !ok {
		return ValueRef{}, fmt.Errorf(
			"cannot index list: index must be Int, but it is %v",
			TypeName(key))
	}
	if index.Int64() < 0 || index.Int64() > math.MaxInt32 {
		return ValueRef{}, fmt.Errorf(
			"cannot index list: invalid index %v",
			index.Int64())
	}

	i := int(index.Int64())
	if i >= len(a.list) {
		return ValueRef{}, fmt.Errorf(
			"cannot index list: index %v is greater than list size %v",
			i, len(a.list))
	}
	return NewValueRef(&a.list[i]), nil
}

func (a *List) Slice(b, e Value) (Value, error) {
	bi, ok := b.(Int)
	ei, ok2 := e.(Int)
	if !ok || !ok2 {
		return nil, fmt.Errorf(
			"slice indices must be Int; instead they are %q and %q",
			TypeName(b), TypeName(e))
	}

	begin := int(bi.Int64())
	end := int(ei.Int64())

	if begin < 0 {
		return nil, fmt.Errorf(
			"invalid slice begin index %v; begin index must be >= 0",
			begin)
	}

	if end < 0 {
		end = len(a.list) + end
	}
	if end > len(a.list) {
		end = len(a.list)
	}
	if end < begin {
		begin = 0
		end = 0
	}

	return NewListWithSlice(a.list[begin:end]), nil
}

func (a *List) AddIter(vm *VM, iter Iterator) error {
	for {
		v, err := vm.IterNext(iter, 1)
		if err != nil {
			return err
		} else if v == nil {
			break
		}
		a.Add(v[0])
	}
	return nil
}

func (a *List) Len() int {
	return len(a.list)
}

func (a *List) String() string {
	return formatObject(a)
}

func (a *List) MakeIterator() Iterator {
	i := &listIter{
		list: a,
		next: 0,
	}
	return NewIterator(i.Next, nil, 2)
}

type listIter struct {
	list *List
	next int
}

func (i *listIter) Next(m *VM, args []Value, nret int) ([]Value, error) {
	if i.next >= i.list.Len() {
		return nil, nil
	}

	idx := i.next
	i.next++

	v := i.list.Get(idx)

	return []Value{v, NewInt(int64(idx))}, nil
}

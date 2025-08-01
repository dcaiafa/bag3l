package vm

import (
	"strconv"
	"strings"
)

type mapNode struct {
	next, prev *mapNode
	key        Value
	value      Value
}

func newMapNodeList() *mapNode {
	l := &mapNode{}
	l.prev = l
	l.next = l
	return l
}

func (n *mapNode) InsertAfter(o *mapNode) {
	n.prev = o
	n.next = o.next
	o.next.prev = n
	o.next = n
}

func (n *mapNode) Remove() {
	n.prev.next = n.next
	n.next.prev = n.prev
	n.next = nil
	n.prev = nil
}

type Map struct {
	data map[Value]*mapNode
	list *mapNode
}

var _ Iterable = (*Map)(nil)

func NewMap() *Map {
	return NewMapWithCapacity(0)
}

func NewMapWithCapacity(c int) *Map {
	return &Map{
		data: make(map[Value]*mapNode, c),
		list: newMapNodeList(),
	}
}

func (o *Map) Type() string   { return "map" }
func (o *Map) Traits() Traits { return TraitNone }

func (o *Map) Len() int {
	return len(o.data)
}

func (o *Map) Has(k Value) bool {
	_, ok := o.data[k]
	return ok
}

func (o *Map) Get(k Value) (Value, bool) {
	n, ok := o.data[k]
	if !ok {
		return nil, false
	}
	return n.value, true
}

func (o *Map) Put(k, v Value) {
	n := o.data[k]
	if n == nil {
		n = &mapNode{key: k}
		n.InsertAfter(o.list.prev)
		o.data[k] = n
	}
	n.value = v
}

func (o *Map) Clone() *Map {
	r := NewMapWithCapacity(o.Len())
	o.ForEach(func(k, v Value) bool {
		r.Put(k, v)
		return true
	})
	return r
}

func (o *Map) Index(k Value) (Value, error) {
	n := o.data[k]
	if n == nil {
		return nil, nil
	}
	return n.value, nil
}

func (o *Map) IndexRef(k Value) (ValueRef, error) {
	n := o.data[k]
	if n == nil {
		n = &mapNode{key: k}
		n.InsertAfter(o.list.prev)
		o.data[k] = n
	}
	return NewValueRef(&n.value), nil
}

func (o *Map) GetFirst() (key Value, val Value) {
	if len(o.data) == 0 {
		return nil, nil
	}
	return o.list.next.key, o.list.next.value
}

func (o *Map) GetNext(key Value) (nextKey Value, nextVal Value) {
	node := o.data[key]
	if node == nil || node.next == o.list {
		return nil, nil
	}
	return node.next.key, node.next.value
}

func (o *Map) Delete(key Value) {
	node := o.data[key]
	if node == nil {
		return
	}
	node.Remove()
	delete(o.data, key)
}

func (o *Map) ForEach(f func(k, v Value) bool) {
	for n := o.list.next; n != o.list; n = n.next {
		if !f(n.key, n.value) {
			return
		}
	}
}

func (o *Map) String() string {
	return formatMap(o)
}

func (o *Map) MakeIterator() Iterator {
	key, _ := o.GetFirst()

	i := &mapIter{
		obj: o,
		key: key,
	}

	return NewIterator(i.Next, nil, 2)
}

type mapFormatter struct {
	visited map[Value]bool
	w       *strings.Builder
}

func formatMap(v Value) string {
	of := &mapFormatter{
		visited: make(map[Value]bool),
		w:       &strings.Builder{},
	}
	of.format(v)
	return of.w.String()
}

func (of *mapFormatter) format(v Value) {
	switch v := v.(type) {
	case String:
		// TODO: escape special characters.
		of.str(`"` + v.String() + `"`)
	case Int:
		of.str(strconv.FormatInt(v.Int64(), 10))
	case Float:
		of.str(strconv.FormatFloat(v.Float64(), 'g', -1, 64))
	case Bool:
		if v.Bool() {
			of.str("true")
		} else {
			of.str("false")
		}
	case *Map:
		if of.visited[v] {
			of.str("<cycle>")
			return
		}
		of.visited[v] = true
		of.str("{")
		first := true
		v.ForEach(func(k, v Value) bool {
			if !first {
				of.str(", ")
			} else {
				first = false
			}
			if ks, ok := k.(String); ok {
				of.str(ks.String())
			} else {
				of.str("[")
				of.format(k)
				of.str("]")
			}
			of.str(": ")
			of.format(v)
			return true
		})
		of.str("}")
	case *List:
		of.str("[")
		first := true
		for i := 0; i < v.Len(); i++ {
			if !first {
				of.str(" ")
			} else {
				first = false
			}
			of.str(ToString(v.Get(i)))
		}
		of.str("]")

	default:
		of.str(ToString(v))
	}
}

func (of *mapFormatter) str(s string) {
	of.w.WriteString(s)
}

type mapIter struct {
	obj *Map
	key Value
}

func (i *mapIter) Next(m *VM, args []Value, nRet int) ([]Value, error) {
	if i.key == nil {
		return nil, nil
	}

	// TODO: This wouldn't work with nil keys/values. Maybe change Index to return
	// a third 'ok' result.

	key := i.key

	val, err := i.obj.Index(key)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, nil
	}

	i.key, _ = i.obj.GetNext(key)

	return []Value{key, val}, nil
}

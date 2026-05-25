package time

import (
	"fmt"
	"time"

	"github.com/dcaiafa/bag3l/internal/vm"
)

//go:generate go run ../../internal/stub/stubgen time.stubgen

type Time struct {
	time time.Time
}

var _ /* implements */ vm.Value = Time{}

func NewTime(t time.Time) Time {
	return Time{time: t}
}

func (t Time) Time() time.Time { return t.time }

func (t Time) String() string    { return t.time.String() }
func (t Time) Type() string      { return "time" }
func (t Time) Traits() vm.Traits { return vm.TraitEq }

func (t Time) EvalOp(op vm.Op, operand vm.Value) (vm.Value, error) {
	switch op {
	case vm.OpEq, vm.OpSub, vm.OpLT, vm.OpLE, vm.OpGT, vm.OpGE:
		operandTime, ok := operand.(Time)
		if !ok {
			if op == vm.OpEq {
				return vm.NewBool(false), nil
			}
			return nil, fmt.Errorf(
				"invalid operation between time and %v",
				vm.TypeName(operand))
		}

		switch op {
		case vm.OpEq:
			return vm.NewBool(t.time.Equal(operandTime.time)), nil
		case vm.OpSub:
			return Duration{t.time.Sub(operandTime.time)}, nil
		case vm.OpLT:
			return vm.NewBool(t.time.Before(operandTime.time)), nil
		case vm.OpLE:
			return vm.NewBool(t.time.Before(operandTime.time) ||
				t.time.Equal(operandTime.time)), nil
		case vm.OpGT:
			return vm.NewBool(t.time.After(operandTime.time)), nil
		case vm.OpGE:
			return vm.NewBool(t.time.After(operandTime.time) ||
				t.time.Equal(operandTime.time)), nil
		default:
			panic("unreachable")
		}

	case vm.OpAdd:
		operandDur, ok := operand.(Duration)
		if !ok {
			return nil, fmt.Errorf(
				"invalid operation between time and %v",
				vm.TypeName(operand))
		}
		return Time{t.time.Add(operandDur.dur)}, nil

	default:
		return nil, vm.ErrOperationNotSupported
	}
}

func utc0(m *vm.VM, t Time) (Time, error) {
	t.time = t.time.UTC()
	return t, nil
}

func local0(m *vm.VM, t Time) (Time, error) {
	t.time = t.time.Local()
	return t, nil
}

func in0(m *vm.VM, t Time, loc *Location) (Time, error) {
	return NewTime(t.Time().In(loc.Location)), nil
}

func in1(m *vm.VM, t Time, loc string) (Time, error) {
	cache := getLocationCache(m)
	locObj, err := cache.GetLocation(loc)
	if err != nil {
		return Time{}, err
	}
	return in0(m, t, locObj)
}

func format0(m *vm.VM, t Time, layout string) (string, error) {
	return t.time.Format(layout), nil
}

func unix0(m *vm.VM, t Time) (int64, error) {
	return t.time.Unix(), nil
}

func unix_nano0(m *vm.VM, t Time) (int64, error) {
	return t.time.UnixNano(), nil
}

func now0(m *vm.VM) (Time, error) {
	return NewTime(time.Now()), nil
}

func parse0(m *vm.VM, v string, layout string) (Time, error) {
	t, err := time.Parse(layout, v)
	if err != nil {
		return Time{}, err
	}
	return NewTime(t), nil
}

func from_unix0(m *vm.VM, sec, nano int64) (Time, error) {
	t := time.Unix(sec, nano)
	return NewTime(t), nil
}

func to_map0(m *vm.VM, t Time) (*vm.Map, error) {
	mp := vm.NewMap()
	mp.Put(vm.NewString("year"), vm.NewInt(int64(t.time.Year())))
	mp.Put(vm.NewString("month"), vm.NewInt(int64(t.time.Month())))
	mp.Put(vm.NewString("day"), vm.NewInt(int64(t.time.Day())))
	mp.Put(vm.NewString("hour"), vm.NewInt(int64(t.time.Hour())))
	mp.Put(vm.NewString("minute"), vm.NewInt(int64(t.time.Minute())))
	mp.Put(vm.NewString("second"), vm.NewInt(int64(t.time.Second())))
	mp.Put(vm.NewString("nanosecond"), vm.NewInt(int64(t.time.Nanosecond())))
	return mp, nil
}

type Duration struct {
	dur time.Duration
}

func NewDuration(dur time.Duration) Duration {
	return Duration{dur: dur}
}

func (d Duration) String() string    { return d.dur.String() }
func (d Duration) Type() string      { return "duration" }
func (d Duration) Traits() vm.Traits { return vm.TraitEq }

func (d Duration) EvalOp(op vm.Op, operand vm.Value) (vm.Value, error) {
	if op == vm.OpUMinus {
		return Duration{d.dur * -1}, nil
	}

	switch op {
	case vm.OpAdd, vm.OpSub, vm.OpLT, vm.OpLE,
		vm.OpGT, vm.OpGE, vm.OpEq, vm.OpMod:

		otherDur := time.Duration(0)
		operandDur, ok := operand.(Duration)
		if ok {
			otherDur = operandDur.dur
		} else if operandInt, ok := operand.(vm.Int); ok && operandInt.Int64() == 0 {
			// Zero is a special case. It is useful to express `dur < 0` without
			// having to create a duration value for the right side.
			otherDur = 0
		} else if op == vm.OpEq {
			return vm.NewBool(false), nil
		} else {
			return nil, vm.ErrOperationNotSupported
		}

		switch op {
		case vm.OpAdd:
			return Duration{d.dur + otherDur}, nil
		case vm.OpSub:
			return Duration{d.dur - otherDur}, nil
		case vm.OpLT:
			return vm.NewBool(d.dur < otherDur), nil
		case vm.OpLE:
			return vm.NewBool(d.dur <= otherDur), nil
		case vm.OpGT:
			return vm.NewBool(d.dur > otherDur), nil
		case vm.OpGE:
			return vm.NewBool(d.dur >= otherDur), nil
		case vm.OpEq:
			return vm.NewBool(d == operand), nil
		case vm.OpMod:
			if otherDur == 0 {
				return nil, vm.ErrDivByZero
			}
			return Duration{d.dur % otherDur}, nil
		}

	case vm.OpMult:
		operandInt, ok := operand.(vm.Int)
		if !ok {
			return nil, vm.ErrOperationNotSupported
		}
		return Duration{d.dur * time.Duration(operandInt.Int64())}, nil

	case vm.OpDiv:
		switch operand := operand.(type) {
		case Duration:
			if operand.dur == 0 {
				return nil, vm.ErrDivByZero
			}
			return vm.NewInt(int64(d.dur / operand.dur)), nil

		case vm.Int:
			if operand.Int64() == 0 {
				return nil, vm.ErrDivByZero
			}
			return Duration{d.dur / time.Duration(operand.Int64())}, nil

		default:
			return nil, vm.ErrOperationNotSupported
		}
	}

	return nil, vm.ErrOperationNotSupported
}

func (d Duration) FallbackEvalOp(op vm.Op, left vm.Value) (vm.Value, error) {
	if op == vm.OpMult {
		if left, ok := left.(vm.Int); ok {
			return NewDuration(time.Duration(left.Int64()) * d.dur), nil
		}
	}
	return nil, vm.ErrOperationNotSupported
}

func (d Duration) Duration() time.Duration {
	return d.dur
}

func truncate0(m *vm.VM, d Duration, mod Duration) (Duration, error) {
	v := d.dur.Truncate(mod.dur)
	return NewDuration(v), nil
}

type Location struct {
	Location *time.Location
}

func (l *Location) String() string    { return l.Location.String() }
func (l *Location) Type() string      { return "location" }
func (l *Location) Traits() vm.Traits { return vm.TraitNone }

func fixed_zone0(m *vm.VM, name string, offset int64) (*Location, error) {
	return &Location{
		Location: time.FixedZone(name, int(offset)),
	}, nil
}

type locationCache struct {
	locations map[string]*Location
}

const locationCacheUserDataKey = "locationCache"

func getLocationCache(m *vm.VM) *locationCache {
	cache, ok := m.GetUserData(locationCacheUserDataKey).(*locationCache)
	if !ok {
		cache = &locationCache{
			locations: make(map[string]*Location),
		}
		m.SetUserData(locationCacheUserDataKey, cache)
	}
	return cache
}

func (c *locationCache) GetLocation(name string) (*Location, error) {
	loc := c.locations[name]
	if loc != nil {
		return loc, nil
	}

	timeLocation, err := time.LoadLocation(name)
	if err != nil {
		return nil, err
	}

	location := &Location{
		Location: timeLocation,
	}

	c.locations[name] = location

	return location, nil
}

func load_location0(m *vm.VM, name string) (*Location, error) {
	loc, err := getLocationCache(m).GetLocation(name)
	if err != nil {
		return nil, err
	}
	return loc, nil
}

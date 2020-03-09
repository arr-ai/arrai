package rel

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/arr-ai/frozen"
)

const DictValueAttr = "@value"

func NewDictTuple(key, value Value) Tuple {
	return NewTuple(NewAttr("@", key), NewAttr(DictValueAttr, value))
}

// Dict is a map from keys to values.
type Dict struct {
	m frozen.Map
}

// NewDict constructs a dict as a relation.
func NewDict(kvs ...[2]Value) Set {
	if len(kvs) == 0 {
		return None
	}
	var mb frozen.MapBuilder
	for _, kv := range kvs {
		mb.Put(kv[0], kv[1])
	}
	return Dict{m: mb.Finish()}
}

func AsDict(s Set) (Dict, bool) {
	if d, ok := s.(Dict); ok {
		return d, true
	}
	if !s.IsTrue() {
		return Dict{}, true
	}
	var mb frozen.MapBuilder
	match := DictTupleMatcher()
	for i := s.Enumerator(); i.MoveNext(); {
		key, value, matched := match(i.Current())
		if !matched {
			return Dict{}, false
		}
		mb.Put(key, value)
	}

	return Dict{m: mb.Finish()}, true
}

func (d Dict) Hash(seed uintptr) uintptr {
	// TODO: Optimize.
	h := seed
	for e := d.Enumerator(); e.MoveNext(); {
		h ^= e.Current().Hash(seed)
	}
	return h
}

func (d Dict) Equal(v interface{}) bool {
	switch v := v.(type) {
	case Dict:
		return d.m.Equal(v.m)
	case Set:
		if d.IsTrue() != v.IsTrue() || d.Count() != v.Count() {
			return false
		}
		match := DictTupleMatcher()
		for e := v.Enumerator(); e.MoveNext(); {
			if key, value, matches := match(e.Current()); matches {
				if dvalue, has := d.m.Get(key); !(has && value.Equal(dvalue)) {
					return false
				}
			}
		}
		return true
	}
	return false
}

func (d Dict) String() string {
	var sb strings.Builder
	sb.WriteString("{")
	for n, i := 0, d.m.Range(); i.Next(); n++ {
		format := ", %v: %v"
		if n == 0 {
			format = format[2:]
		}
		key, value := i.Entry()
		fmt.Fprintf(&sb, format, key, value)
	}
	sb.WriteString("}")
	return sb.String()
}

func (d Dict) Eval(local Scope) (Value, error) {
	return d, nil
}

var dictKind = registerKind(208, reflect.TypeOf(String{}))

// Kind returns a number that is unique for each major kind of Value.
func (d Dict) Kind() int {
	return dictKind
}

func (d Dict) IsTrue() bool {
	return !d.m.IsEmpty()
}

func (d Dict) Less(v Value) bool {
	if d.Kind() != v.Kind() {
		return d.Kind() < v.Kind()
	}
	panic("unfinished")
}

func (d Dict) Negate() Value {
	return NewTuple(NewAttr(negateTag, d))
}

func (d Dict) Export() interface{} {
	var mb frozen.MapBuilder
	for i := d.m.Range(); i.Next(); {
		k, v := i.Entry()
		mb.Put(k, v)
	}
	return mb.Finish()
}

func (d Dict) Count() int {
	return d.m.Count()
}

func (d Dict) Has(v Value) bool {
	if key, value, matched := DictTupleMatcher()(v); matched {
		if v, has := d.m.Get(key); has {
			return value.Equal(v)
		}
	}
	return false
}

func (d Dict) Enumerator() ValueEnumerator {
	return &dictEnumerator{i: d.m.Range()}
}

func (d Dict) With(v Value) Set {
	if key, value, matched := DictTupleMatcher()(v); matched {
		return Dict{m: d.m.With(key, value)}
	}
	return d
}

func (d Dict) Without(v Value) Set {
	if key, value, matched := DictTupleMatcher()(v); matched {
		if v, has := d.m.Get(key); has {
			if value.Equal(v) {
				return Dict{m: d.m.Without(frozen.NewSet(key))}
			}
		}
	}
	return d
}

func (d Dict) Map(m func(Value) Value) Set {
	var sb frozen.SetBuilder
	for e := d.Enumerator(); e.MoveNext(); {
		sb.Add(m(e.Current()))
	}
	return &genericSet{set: sb.Finish()}
}

func (d Dict) Where(pred func(Value) bool) Set {
	var mb frozen.MapBuilder
	for e := d.m.Range(); e.Next(); {
		k, v := e.Entry()
		key := k.(Value)
		value := v.(Value)
		if pred(value) {
			mb.Put(key, value)
		}
	}
	return Dict{m: mb.Finish()}
}

func (d Dict) Call(arg Value) Value {
	return d.m.MustGet(arg).(Value)
}

func (d Dict) ArrayEnumerator() (ValueEnumerator, bool) {
	return nil, false
}

func (d Dict) DictEnumerator() *DictEnumerator {
	return &DictEnumerator{i: d.m.Range()}
}

func DictTupleMatcher() func(v Value) (key, value Value, matches bool) {
	var key, value Value
	m := NewTupleMatcher(
		map[string]Matcher{
			"@":           Let(func(k Value) { key = k }),
			DictValueAttr: Let(func(v Value) { value = v }),
		},
		Lit(EmptyTuple),
	)
	return func(v Value) (Value, Value, bool) {
		matches := m.Match(v)
		return key, value, matches
	}
}

type dictEnumerator struct {
	i *frozen.MapIterator
}

func (a *dictEnumerator) MoveNext() bool {
	return a.i.Next()
}

func (a *dictEnumerator) Current() Value {
	return NewDictTuple(a.i.Key().(Value), a.i.Value().(Value))
}

type DictEnumerator struct {
	i *frozen.MapIterator
}

func (a *DictEnumerator) MoveNext() bool {
	return a.i.Next()
}

func (a *DictEnumerator) Current() (key, value Value) {
	return a.i.Key().(Value), a.i.Value().(Value)
}

package rel

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/arr-ai/frozen"
)

type multipleValues frozen.Set

// Dict is a map from keys to values.
type Dict struct {
	m frozen.Map
}

// NewDict constructs a dict as a relation {|@, @value|...}.
func NewDict(allowDupKeys bool, entries ...DictEntryTuple) Set {
	if len(entries) == 0 {
		return None
	}
	var mb frozen.MapBuilder
	for _, entry := range entries {
		if v, has := mb.Get(entry.at); has {
			if !allowDupKeys {
				panic(fmt.Errorf("duplicate key: %v", entry.at))
			}
			switch v := v.(type) {
			case multipleValues:
				mb.Put(entry.at, multipleValues(frozen.Set(v).With(entry.value)))
			default:
				mb.Put(entry.at, multipleValues(frozen.NewSet(v, entry.value)))
			}
		} else {
			mb.Put(entry.at, entry.value)
		}
	}
	return Dict{m: mb.Finish()}
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

func (d Dict) OrderedEntries() []DictEntryTuple {
	result := make(dictEntryTupleSort, 0, d.Count())
	for e := d.Enumerator(); e.MoveNext(); {
		result = append(result, e.Current().(DictEntryTuple))
	}
	sort.Sort(result)
	return result
}

func (d Dict) Eval(local Scope) (Value, error) {
	return d, nil
}

var dictKind = registerKind(209, reflect.TypeOf(String{}))

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
	return &dictEnumerator{i: d.m.Range(), j: frozen.Set{}.Range()}
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
	return GenericSet{set: sb.Finish()}
}

func (d Dict) Where(pred func(Value) bool) Set {
	var mb frozen.MapBuilder
	for e := d.Enumerator(); e.MoveNext(); {
		t := e.Current().(DictEntryTuple)
		if pred(t) {
			mb.Put(t.at, t.value)
		}
	}
	m := mb.Finish()
	if m.IsEmpty() {
		return None
	}
	return Dict{m: m}
}

func (d Dict) Call(arg Value) Value {
	return d.m.MustGet(arg).(Value)
}

func (d Dict) ArrayEnumerator() (OffsetValueEnumerator, bool) {
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
	j frozen.Iterator
	v Value
}

func (a *dictEnumerator) MoveNext() bool {
	if !a.j.Next() {
		if !a.i.Next() {
			return false
		}
		switch entry := a.i.Value().(type) {
		case multipleValues:
			a.j = frozen.Set(entry).Range()
			if !a.j.Next() {
				return false
			}
		default:
			a.j = frozen.Set{}.Range()
			a.v = NewDictEntryTuple(a.i.Key().(Value), entry.(Value))
			return true
		}
	}
	a.v = NewDictEntryTuple(a.i.Key().(Value), a.j.Value().(Value))
	return true
}

func (a *dictEnumerator) Current() Value {
	return a.v
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

type dictEntryTupleSort []DictEntryTuple

func (s dictEntryTupleSort) Len() int {
	return len(s)
}

func (s dictEntryTupleSort) Less(a, b int) bool {
	x := s[a]
	y := s[b]
	if !x.at.Equal(y.at) {
		return x.at.Less(y.at)
	}
	return x.value.Less(y.value)
}

func (s dictEntryTupleSort) Swap(a, b int) {
	s[a], s[b] = s[b], s[a]
}

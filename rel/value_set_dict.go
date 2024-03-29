package rel

import (
	"context"
	"fmt"
	"reflect"
	"sort"

	"github.com/go-errors/errors"

	"github.com/arr-ai/frozen"
	"github.com/arr-ai/wbnf/parser"

	"github.com/arr-ai/arrai/pkg/fu"
)

type multipleValues frozen.Set

var _ frozen.Key = multipleValues(frozen.Set{})

func newMultipleValues(values ...interface{}) interface{} {
	s := frozen.NewSet(values...)
	if s.Count() == 1 {
		return s.Any()
	}
	return multipleValues(s)
}

func (m multipleValues) Equal(n interface{}) bool {
	if n, is := n.(multipleValues); is {
		return frozen.Set(m).EqualSet(frozen.Set(n))
	}
	return frozen.Set(m).Equal(n)
}

func (m multipleValues) Hash(seed uintptr) uintptr {
	return frozen.Set(m).Hash(seed)
}

func (m multipleValues) String() string {
	return frozen.Set(m).String()
}

// Dict is a map from keys to values.
type Dict struct {
	m frozen.Map
}

// AsDict checks whether a Value is a valid dictionary.
func AsDict(v Value) (Dict, bool) {
	switch v := v.(type) {
	case Dict:
		return v, true
	case EmptySet:
		return Dict{}, true
	}
	return Dict{}, false
}

// MustNewDict constructs a dict as a relation {|@, @value|...}, or panics if construction fails.
func MustNewDict(allowDupKeys bool, entries ...DictEntryTuple) Set {
	d, err := NewDict(allowDupKeys, entries...)
	if err != nil {
		panic(err)
	}
	return d
}

// NewDict constructs a dict as a relation {|@, @value|...}.
func NewDict(allowDupKeys bool, entries ...DictEntryTuple) (Set, error) {
	if len(entries) == 0 {
		return None, nil
	}
	var mb frozen.MapBuilder
	for _, entry := range entries {
		if v, has := mb.Get(entry.at); has {
			if !allowDupKeys {
				return nil, errors.Errorf("duplicate key: %v", entry.at)
			}
			switch v := v.(type) {
			case multipleValues:
				mb.Put(entry.at, multipleValues(frozen.Set(v).With(entry.value)))
			default:
				mb.Put(entry.at, newMultipleValues(v, entry.value))
			}
		} else {
			mb.Put(entry.at, entry.value)
		}
	}
	return Dict{m: mb.Finish()}, nil
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
			} else {
				// does not match the dictionary tuple format
				return false
			}
		}
		return true
	}
	return false
}

func (d Dict) String() string {
	return fu.String(d)
}

func (d Dict) Format(f fmt.State, verb rune) {
	fu.WriteString(f, "{")
	for i, t := range d.OrderedEntries() {
		writeSep(f, i, ", ")
		fu.Fprintf(f, "%v: %v", t.at, t.value)
	}
	fu.WriteString(f, "}")
}

func (d Dict) OrderedEntries() []DictEntryTuple {
	result := make(dictEntryTupleSort, 0, d.Count())
	for e := d.Enumerator(); e.MoveNext(); {
		result = append(result, e.Current().(DictEntryTuple))
	}
	sort.Sort(result)
	return result
}

func (d Dict) Eval(ctx context.Context, local Scope) (Value, error) {
	return d, nil
}

// Source returns a scanner locating the Dict's source code.
func (d Dict) Source() parser.Scanner {
	return *parser.NewScanner("")
}

var dictKind = registerKind(209, reflect.TypeOf(Dict{}))

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
	dKeys := d.m.Keys().OrderedElements(intfValueLess)
	vDict := v.(Dict)
	vKeys := vDict.m.Keys().OrderedElements(intfValueLess)
	n := len(dKeys)
	if n > len(vKeys) {
		n = len(vKeys)
	}
	for i, k := range dKeys[:n] {
		dKey := k.(Value)
		vKey := vKeys[i].(Value)
		if !dKey.Equal(vKey) {
			return dKey.Less(vKey)
		}

		// TODO: Implement Less directly in frozen.
		var dValues []interface{}
		switch dValue := d.m.MustGet(dKey).(type) {
		case multipleValues:
			dValues = frozen.Set(dValue).OrderedElements(intfValueLess)
		case Value:
			dValues = []interface{}{dValue}
		default:
			panic("wtf?")
		}

		var vValues []interface{}
		switch dValue := vDict.m.MustGet(vKey).(type) {
		case multipleValues:
			vValues = frozen.Set(dValue).OrderedElements(intfValueLess)
		case Value:
			vValues = []interface{}{dValue}
		default:
			panic("wtf?")
		}

		n := len(dValues)
		if n > len(vValues) {
			n = len(vValues)
		}
		for i, dItem := range dValues[:n] {
			dValue := dItem.(Value)
			vValue := vValues[i].(Value)
			if !dValue.Equal(vValue) {
				return dValue.Less(vValue)
			}
		}
		if len(dValues) != len(vValues) {
			return len(dValues) < len(vValues)
		}
	}
	return len(dKeys) < len(vKeys)
}

func (d Dict) Negate() Value {
	return NewTuple(NewAttr(negateTag, d))
}

func (d Dict) Export(_ context.Context) interface{} {
	var mb frozen.MapBuilder
	for i := d.m.Range(); i.Next(); {
		k, v := i.Entry()
		mb.Put(k, v)
	}
	return mb.Finish()
}

func (Dict) getSetBuilder() setBuilder {
	return newGenericTypeSetBuilder()
}

func (Dict) getBucket() fmt.Stringer {
	return genericType
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
	if t, is := v.(DictEntryTuple); is {
		if u, has := d.m.Get(t.at); has {
			switch u := u.(type) {
			case multipleValues:
				return Dict{m: d.m.With(t.at, multipleValues(frozen.Set(u).With(t.value)))}
			default:
				return Dict{m: d.m.With(t.at, newMultipleValues(u, t.value))}
			}
		}
		return Dict{m: d.m.With(t.at, t.value)}
	}
	return toUnionSetWithItem(d, v)
}

func (d Dict) Without(v Value) Set {
	if key, value, matched := DictTupleMatcher()(v); matched {
		if v, has := d.m.Get(key); has {
			if value.Equal(v) {
				m := d.m.Without(frozen.NewSet(key))
				if m.IsEmpty() {
					return None
				}
				return Dict{m: m}
			}
		}
	}
	return d
}

func (d Dict) Map(f func(v Value) (Value, error)) (Set, error) {
	sb := NewSetBuilder()
	for e := d.Enumerator(); e.MoveNext(); {
		v, err := f(e.Current())
		if err != nil {
			return nil, err
		}
		sb.Add(v)
	}
	return sb.Finish()
}

func (d Dict) Where(p func(v Value) (bool, error)) (Set, error) {
	var mb frozen.MapBuilder
	for e := d.Enumerator(); e.MoveNext(); {
		t := e.Current().(DictEntryTuple)
		match, err := p(t)
		if err != nil {
			return nil, err
		}
		if match {
			mb.Put(t.at, t.value)
		}
	}
	m := mb.Finish()
	if m.IsEmpty() {
		return None, nil
	}
	return Dict{m: m}, nil
}

func (d Dict) CallAll(_ context.Context, arg Value, b SetBuilder) error {
	if val, has := d.m.Get(arg); has {
		switch v := val.(type) {
		case Value:
			b.Add(v)
		case multipleValues:
			for e := frozen.Set(v).Range(); e.Next(); {
				b.Add(e.Value().(Value))
			}
		}
	}
	return nil
}

func (Dict) unionSetSubsetBucket() string {
	return DictEntryTuple{}.getBucket().String()
}

type emptyEnumerator struct{}

func (emptyEnumerator) Current() Value {
	panic("wtf")
}

func (emptyEnumerator) MoveNext() bool {
	return false
}

func (d Dict) ArrayEnumerator() ValueEnumerator {
	return &emptyEnumerator{}
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
	k, ok := a.i.Key().(Value)
	if !ok {
		panic(fmt.Errorf("key is not a Value: %s %[1]s", a.i.Key()))
	}
	v, ok := a.i.Value().(Value)
	if !ok {
		panic(fmt.Errorf("dict value for key %s is not a Value (type: %T): %[2]v", k, a.i.Value()))
	}
	return k, v
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

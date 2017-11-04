package rel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"

	"github.com/OneOfOne/xxhash"
	"github.com/mediocregopher/seq"
)

// GenericTuple is the default implementation of Tuple.
type GenericTuple struct {
	tuple *seq.HashMap
	names []string
}

var (
	// EmptyTuple is the tuple with no attributes.
	EmptyTuple Tuple = &GenericTuple{tuple: seq.NewHashMap()}

	hashDelim = []byte{0}
	negateTag = "@neg"
)

// NewAttr returns an Attr with the given name and value.
func NewAttr(name string, value Value) Attr {
	return Attr{Name: name, Value: value}
}

// NewBoolAttr return an attr with a bool value.
func NewBoolAttr(name string, value bool) Attr {
	return NewAttr(name, NewBool(value))
}

// NewFloatAttr return an attr with a float value.
func NewFloatAttr(name string, value float64) Attr {
	return NewAttr(name, NewNumber(value))
}

// NewIntAttr return an attr with an int value.
func NewIntAttr(name string, value int) Attr {
	return NewFloatAttr(name, float64(value))
}

// NewStringAttr return an attr with a string value.
func NewStringAttr(name string, value []rune) Attr {
	return NewAttr(name, NewString(value))
}

// NewTupleAttr return an attr with a new tuple value.
func NewTupleAttr(name string, attrs ...Attr) Attr {
	return NewAttr(name, NewTuple(attrs...))
}

// NewTuple constructs a Tuple from attrs. Passes each Val to NewValue().
func NewTuple(attrs ...Attr) Tuple {
	tuple := EmptyTuple
	for _, kv := range attrs {
		tuple, _ = tuple.With(kv.Name, kv.Value)
	}
	return tuple
}

// NewTupleFromMap constructs a Tuple from a map of strings to Go values.
func NewTupleFromMap(m map[string]interface{}) (Tuple, error) {
	tuple := EmptyTuple
	for name, intf := range m {
		value, err := NewValue(intf)
		if err != nil {
			return nil, err
		}
		tuple, _ = tuple.With(name, value)
	}
	return tuple, nil
}

// NewXML constructs an XML Tuple from the given data
func NewXML(tag []rune, attrs []Attr, children ...Value) Tuple {
	t, _ := EmptyTuple.With("tag", NewString(tag))
	if len(attrs) != 0 {
		t, _ = t.With("attributes", NewTuple(attrs...))
	}
	if len(children) != 0 {
		t, _ = t.With("children", NewArray(children...))
	}
	xml, _ := EmptyTuple.With("@xml", t)
	return xml
}

// Hash computes a hash for a GenericTuple.
func (t *GenericTuple) Hash(seed uint32) uint32 {
	xx := xxhash.NewS32(seed)
	h1, h2 := seed+0x6b783347, seed+0x7b4da23d
	for e := t.Enumerator(); e.MoveNext(); {
		name, value := e.Current()
		h1 ^= value.Hash(seed)
		xx.Write([]byte(name))
		h2 ^= xx.Sum32()
		xx.Reset()
	}
	return h1 + 3*h2
}

// Equal tests two Tuples for equality. Any other type returns false.
func (t *GenericTuple) Equal(v interface{}) bool {
	if b, ok := v.(Tuple); ok {
		for e := t.Enumerator(); e.MoveNext(); {
			aName, aValue := e.Current()
			if bVal, found := b.Get(aName); found {
				if !aValue.Equal(bVal) {
					return false
				}
			} else {
				return false
			}
		}
		for e := b.Enumerator(); e.MoveNext(); {
			name, _ := e.Current()
			if _, found := t.Get(name); !found {
				return false
			}
		}
		return true
	}
	return false
}

// LexerNamePat defines valid unquoted identifiers.
// This really belongs in rel/syntax/lex.go, but that creates a dep cycle.
var LexerNamePat = `([$@A-Za-z_][0-9$@A-Za-z_]*)`

var identRE = regexp.MustCompile(`\A` + LexerNamePat + `\z`)

// String returns a string representation of a Tuple.
func (t *GenericTuple) String() string {
	var buf bytes.Buffer
	buf.WriteRune('{')
	for i, name := range tupleOrderedNames(t) {
		if i != 0 {
			buf.WriteString(", ")
		}
		if identRE.Match([]byte(name)) {
			buf.WriteString(name)
		} else {
			data, err := json.Marshal(name)
			if err != nil {
				panic(err)
			}
			buf.Write(data)
		}
		buf.WriteString(": ")
		value, found := t.Get(name)
		if !found {
			panic(fmt.Sprintf(
				"walk() produced name, %v, which fails lookup", name))
		}
		buf.WriteString(value.String())
	}
	buf.WriteRune('}')
	return buf.String()
}

// Eval returns the tuple.
func (t *GenericTuple) Eval(local, global *Scope) (Value, error) {
	return t, nil
}

// Kind returns a number that is unique for each major kind of Value.
func (t *GenericTuple) Kind() int {
	if t.Count() == 1 {
		if x, ok := t.Get(negateTag); ok {
			return -x.Kind()
		}
	}
	return 300
}

// Bool returns true iff the tuple has attributes.
func (t *GenericTuple) Bool() bool {
	return t.Count() > 0
}

// Less returns true iff v is not a number or tuple, or v is a tuple and t
// precedes v in a lexicographical comparison of their name/value pairs.
func (t *GenericTuple) Less(v Value) bool {
	if t.Kind() != v.Kind() {
		return t.Kind() < v.Kind()
	}
	if t.Count() == 1 {
		if x, ok := t.Get(negateTag); ok {
			u := v.(Tuple)
			if u.Count() != 1 {
				panic(negateTag + " kind not single-attr tuple")
			}
			if y, ok := v.(Tuple).Get(negateTag); ok {
				return y.Less(x)
			}
			panic(negateTag + " kind missing " + negateTag + " attr")
		}
	}

	x := v.(*GenericTuple)
	a := tupleOrderedNames(t)
	b := tupleOrderedNames(x)
	n := len(a)
	if n > len(b) {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
		va, _ := t.Get(a[i])
		vb, _ := x.Get(b[i])
		if va.Less(vb) {
			return true
		}
		if vb.Less(va) {
			return false
		}
	}
	return len(a) < len(b)
}

// Negate returns x if t matches {(negateTag): x} else {(negateTag): t}.
func (t *GenericTuple) Negate() Value {
	if t.Count() == 1 {
		if x, ok := t.Get(negateTag); ok {
			return x
		}
	}
	if !t.Bool() {
		return t
	}
	return NewTuple(NewAttr(negateTag, t))
}

// Export exports a Tuple.
func (t *GenericTuple) Export() interface{} {
	result := make(map[string]interface{}, t.Count())
	for e := t.Enumerator(); e.MoveNext(); {
		name, value := e.Current()
		result[name] = value.Export()
	}
	return result
}

// Count returns how many attributes are in the Tuple.
func (t *GenericTuple) Count() uint64 {
	return t.tuple.Size()
}

// Get returns the Value associated with a name, and true iff it was found.
func (t *GenericTuple) Get(name string) (Value, bool) {
	if v, found := t.tuple.Get(name); found {
		return v.(Value), true
	}
	return nil, false
}

// With returns a Tuple with all name/Value pairs in t (except the one for the
// given name, if present) with the addition of the given name/Value pair, and
// true iff it was newly added.
func (t *GenericTuple) With(name string, value Value) (Tuple, bool) {
	// Strip view/non-view counterpart.
	if name[:1] == "&" {
		u, _ := t.Without(name[1:])
		t = u.(*GenericTuple)
	} else {
		u, _ := t.Without("&" + name)
		t = u.(*GenericTuple)
	}
	hm, added := t.tuple.Set(name, value)
	return &GenericTuple{tuple: hm}, added
}

// Without returns a Tuple with all name/Value pairs in t exception the one of
// the given name, and true iff it was present.
func (t *GenericTuple) Without(name string) (Tuple, bool) {
	hm, removed := t.tuple.Del(name)
	return &GenericTuple{tuple: hm}, removed
}

// HasName returns true iff the Tuple has an attribute with the given name.
func (t *GenericTuple) HasName(name string) bool {
	_, found := t.tuple.Get(name)
	return found
}

// Attributes returns attributes as a map.
func (t *GenericTuple) Attributes() map[string]Value {
	attrs := make(map[string]Value, t.Count())
	for e := t.Enumerator(); e.MoveNext(); {
		name, value := e.Current()
		attrs[name] = value
	}
	return attrs
}

// Names returns the attribute names as a slice.
func (t *GenericTuple) Names() *Names {
	names := EmptyNames
	for e := t.Enumerator(); e.MoveNext(); {
		name, _ := e.Current()
		names = names.With(name)
	}
	return names
}

// Project returns a tuple with the given names from this tuple, or nil if any
// name wasn't found.
func (t *GenericTuple) Project(names *Names) Tuple {
	result := NewTuple()
	for e := names.Enumerator(); e.MoveNext(); {
		name := e.Current()
		value, found := t.Get(name)
		if !found {
			return nil
		}
		result, _ = result.With(name, value)
	}
	return result
}

// GenericTupleEnumerator represents an enumerator over a GenericTuple.
type GenericTupleEnumerator struct {
	ts      tupleSeq
	current Attr
}

// MoveNext moves the enumerator to the next Value.
func (e *GenericTupleEnumerator) MoveNext() bool {
	name, value, ts, ok := e.ts.walk()
	e.ts = ts
	e.current = Attr{name, value}
	return ok
}

// Current returns the enumerator's current Value.
func (e *GenericTupleEnumerator) Current() (string, Value) {
	return e.current.Name, e.current.Value
}

// Enumerator returns an enumerator over the Values in the GenericTuple.
func (t *GenericTuple) Enumerator() AttrEnumerator {
	return &GenericTupleEnumerator{ts: tupleSeq{t.tuple}}
}

// orderedNames returns the names of this tuple in sorted order.
func tupleOrderedNames(t *GenericTuple) []string {
	if t.names == nil {
		t.names = t.Names().ToSlice()
		sort.Strings(t.names)
	}
	return t.names
}

type tupleSeq struct {
	s seq.Seq
}

// walk returns another name/Value pair and access to the rest.
func (tr tupleSeq) walk() (string, Value, tupleSeq, bool) {
	v, s, ok := tr.s.FirstRest()
	if !ok {
		return "", nil, tupleSeq{}, false
	}
	kv := v.(*seq.KV)
	return kv.Key.(string), kv.Val.(Value), tupleSeq{s}, true
}

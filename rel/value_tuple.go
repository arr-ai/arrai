package rel

import (
	"github.com/arr-ai/hash/hash128"

	"context"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/arr-ai/wbnf/parser"

	"github.com/arr-ai/arrai/pkg/fu"
)

// GenericTuple is the default implementation of Tuple: an interned Shape
// (the sorted attribute names and everything derived from them) plus the
// values in shape order. See Shape.
type GenericTuple struct {
	shape *Shape
	vals  []Value
	hash  hashCell
}

var (
	// EmptyTuple is the tuple with no attributes.
	EmptyTuple Tuple = &GenericTuple{shape: emptyShape}

	negateTag = "@neg"
)

// sh returns the tuple's shape, treating a zero GenericTuple as empty.
func (t *GenericTuple) sh() *Shape {
	if t.shape == nil {
		return emptyShape
	}
	return t.shape
}

// newShapedTuple wraps values already in shape order.
func newShapedTuple(shape *Shape, vals []Value) *GenericTuple {
	return &GenericTuple{shape: shape, vals: vals}
}

// TupleBuilder accumulates attributes for a tuple. A later Put of the same
// name replaces an earlier one.
type TupleBuilder struct {
	attrs []Attr
}

func (b *TupleBuilder) Put(name string, value Value) {
	b.attrs = append(b.attrs, Attr{Name: name, Value: value})
}

// Finish returns the tuple, canonicalised to a specialised kind where one
// applies (index/char, index/byte, index/item and key/value pairs).
func (b *TupleBuilder) Finish() Tuple {
	if len(b.attrs) == 0 {
		return EmptyTuple
	}
	shape, vals := shapeOfAttrs(b.attrs)
	return canonicalTuple(shape, vals)
}

// canonicalTuple returns the specialised tuple kind for the "@"-indexed
// two-attribute shapes, or a GenericTuple otherwise.
func canonicalTuple(shape *Shape, vals []Value) Tuple {
	if len(vals) == 2 && shape.names[0] == "@" {
		i := vals[0]
		switch shape.names[1] {
		case StringCharAttr:
			return NewStringCharTuple(int(i.(Number).Float64()), rune(vals[1].(Number).Float64()))
		case BytesByteAttr:
			return NewBytesByteTuple(int(i.(Number).Float64()), byte(vals[1].(Number).Float64()))
		case ArrayItemAttr:
			return NewArrayItemTuple(int(i.(Number).Float64()), vals[1])
		case DictValueAttr:
			return NewDictEntryTuple(i, vals[1])
		}
	}
	return newShapedTuple(shape, vals)
}

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

// NewIuntAttr return an attr for a uint64 value.
func NewUintAttr(name string, value uint64) Attr {
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
	if len(attrs) == 2 {
		if attrs[1].Name == "@" {
			attrs[0], attrs[1] = attrs[1], attrs[0]
		}
		if attrs[0].Name == "@" && strings.HasPrefix(attrs[1].Name, "@") {
			switch attrs[1].Name {
			case StringCharAttr:
				return NewStringCharTuple(
					int(attrs[0].Value.(Number).Float64()),
					rune(attrs[1].Value.(Number).Float64()),
				)
			case BytesByteAttr:
				return NewBytesByteTuple(
					int(attrs[0].Value.(Number).Float64()),
					byte(attrs[1].Value.(Number).Float64()),
				)
			case ArrayItemAttr:
				return NewArrayItemTuple(int(attrs[0].Value.(Number).Float64()), attrs[1].Value)
			case DictValueAttr:
				return NewDictEntryTuple(attrs[0].Value, attrs[1].Value)
			}
		}
	}
	return newTuple(attrs...)
}

func newTuple(attrs ...Attr) Tuple {
	var b TupleBuilder
	for _, kv := range attrs {
		b.Put(kv.Name, kv.Value)
	}
	return b.Finish()
}

// NewTupleFromMap constructs a Tuple from a map of strings to Go values.
func NewTupleFromMap(m map[string]interface{}) (Tuple, error) {
	var b TupleBuilder
	for name, intf := range m {
		value, err := NewValue(intf)
		if err != nil {
			return nil, err
		}
		b.Put(name, value)
	}
	return b.Finish(), nil
}

//TODO: Expand to handle all types and rely less on default replace.
// Potentially, each type could have a `Merge` func that can be called and the switch statement wouldn't be necessary

// MergeTuples takes two tuples and performs a deep merge
func MergeTuples(tuple Tuple, tuple2 Tuple) Tuple {
	tempTuple := tuple
	for e := tuple2.Enumerator(); e.MoveNext(); {
		name, value := e.Current()
		v, found := tempTuple.Get(name)
		if found && v.Kind() == value.Kind() {
			switch v.(type) {
			case Tuple:
				tempTuple = tempTuple.With(name, MergeTuples(v.(Tuple), value.(Tuple)))
			case Dict:
				tempTuple = tempTuple.With(name, mergeDicts(v.(Dict), value.(Dict)))
			default:
				tempTuple = tempTuple.With(name, value)
			}
		} else {
			tempTuple = tempTuple.With(name, value)
		}
	}
	return tempTuple
}

// NewXML constructs an XML Tuple from the given data
func NewXML(tag []rune, attrs []Attr, children ...Value) Tuple {
	var b TupleBuilder
	b.Put("tag", NewString(tag))
	if len(attrs) != 0 {
		b.Put("attributes", NewTuple(attrs...))
	}
	if len(children) != 0 {
		b.Put("children", NewArray(children...))
	}
	return EmptyTuple.With("@xml", b.Finish())
}

// newGenericTuple always returns a generic tuple no matter the attributes.
func newGenericTuple(attrs ...Attr) Tuple {
	if len(attrs) == 0 {
		return EmptyTuple
	}
	return newShapedTuple(shapeOfAttrs(attrs))
}

func (t *GenericTuple) Canonical() Tuple {
	attrs := make([]Attr, 0, t.Count())
	for e := t.Enumerator(); e.MoveNext(); {
		name, value := e.Current()
		attrs = append(attrs, NewAttr(name, value))
	}
	return NewTuple(attrs...)
}

// Hash computes a hash for a GenericTuple.
func (t *GenericTuple) Hash(seed uintptr) uintptr {
	return t.Hash128().Seeded(seed)
}

// Hash128 computes the 128-bit hash of a GenericTuple once: the xor over its
// attributes of the name hash mixed with the value hash.
func (t *GenericTuple) Hash128() hash128.H128 {
	return t.hash.get(func() hash128.H128 {
		h := tupleSalt
		nameH := t.sh().nameH
		for i, v := range t.vals {
			h = h.Xor(hashAttr(nameH[i], v))
		}
		return h
	})
}

// Equal tests two Tuples for equality. Any other type returns false.
func (t *GenericTuple) Equal(v Value) bool {
	if u, ok := v.(*GenericTuple); ok {
		// Shapes are interned: same attribute set iff same shape.
		if t.sh() != u.sh() {
			return false
		}
		for i, tv := range t.vals {
			if !tv.Equal(u.vals[i]) {
				return false
			}
		}
		return true
	}
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

func TupleNameRepr(name string) string {
	if identRE.Match([]byte(name)) {
		return name
	}
	var sb strings.Builder
	switch {
	case !strings.Contains(name, "'"):
		reprEscape(name, '\'', &sb)
	default:
		reprEscape(name, '"', &sb)
	}
	return sb.String()
}

// String returns a string representation of a Tuple.
func (t *GenericTuple) String() string {
	return fu.String(t)
}

// Format formats a GenericTuple.
func (t *GenericTuple) Format(f fmt.State, verb rune) {
	fu.WriteString(f, "(")
	for i, name := range TupleOrderedNames(t) {
		writeSep(f, i, ", ")
		fu.WriteString(f, TupleNameRepr(name))
		fu.WriteString(f, ": ")
		fu.Fprintf(f, "%v", t.MustGet(name))
	}
	fu.WriteString(f, ")")
}

// Eval returns the tuple.
func (t *GenericTuple) Eval(ctx context.Context, local Scope) (Value, error) {
	return t, nil
}

// Source returns a scanner locating the GenericTuple's source code.
func (t *GenericTuple) Source() parser.Scanner {
	return *parser.NewScanner("")
}

var genericTupleKind = registerKind(300, reflect.TypeOf((*GenericTuple)(nil)))

// Kind returns a number that is unique for each major kind of Value.
func (t *GenericTuple) Kind() int {
	if t.Count() == 1 {
		if x, ok := t.Get(negateTag); ok {
			return -x.Kind()
		}
	}
	return genericTupleKind
}

// Bool returns true iff the tuple has attributes.
func (t *GenericTuple) IsTrue() bool {
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
			if y, ok := u.Get(negateTag); ok {
				return y.Less(x)
			}
			panic(negateTag + " kind missing " + negateTag + " attr")
		}
	}

	u := v.(Tuple)
	a := TupleOrderedNames(t)
	b := orderedTupleNames(u)
	n := len(a)
	if n > len(b) {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
		va, _ := t.Get(a[i])
		vb, _ := u.Get(b[i])
		if va.Less(vb) {
			return true
		}
		if vb.Less(va) {
			return false
		}
	}
	return len(a) < len(b)
}

// orderedTupleNames returns sorted attribute names for any Tuple.
func orderedTupleNames(t Tuple) []string {
	if g, ok := t.(*GenericTuple); ok {
		return TupleOrderedNames(g)
	}
	names := make([]string, 0, t.Count())
	for e := t.Enumerator(); e.MoveNext(); {
		name, _ := e.Current()
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Negate returns x if t matches {(negateTag): x} else {(negateTag): t}.
func (t *GenericTuple) Negate() Value {
	if t.Count() == 1 {
		if x, ok := t.Get(negateTag); ok {
			return x
		}
	}
	if !t.IsTrue() {
		return t
	}
	return NewTuple(NewAttr(negateTag, t))
}

// Export exports a Tuple.
func (t *GenericTuple) Export(ctx context.Context) interface{} {
	result := make(map[string]interface{}, t.Count())
	for e := t.Enumerator(); e.MoveNext(); {
		name, value := e.Current()
		result[name] = value.Export(ctx)
	}
	return result
}

func (t *GenericTuple) getSetBuilder() setBuilder {
	if !t.IsTrue() {
		return newGenericTypeSetBuilder()
	}
	return newRelationBuilder(TupleOrderedNames(t), 0)
}

func (t *GenericTuple) getBucket() fmt.Stringer {
	if !t.IsTrue() {
		return genericType
	}
	return t.sh().bucket
}

type NamesSlice []string

func (n NamesSlice) hasIntersect(n2 NamesSlice) bool {
	if len(n) > len(n2) {
		n, n2 = n2, n
	}
	m := n.intoSet()
	for _, name := range n2 {
		if _, has := m[name]; has {
			return true
		}
	}
	return false
}

func (n NamesSlice) intersect(n2 NamesSlice) NamesSlice {
	if len(n) > len(n2) {
		n, n2 = n2, n
	}
	m := n.intoSet()
	intersects := make(NamesSlice, 0, len(n))
	for _, name := range n2 {
		if _, has := m[name]; has {
			intersects = append(intersects, name)
		}
	}
	return intersects
}

func (n NamesSlice) minus(n2 NamesSlice) NamesSlice {
	names := make(NamesSlice, 0, len(n))
	m := n2.intoSet()
	for _, name := range n {
		if _, has := m[name]; !has {
			names = append(names, name)
		}
	}
	return names
}

func (n NamesSlice) isSubset(n2 NamesSlice) bool {
	m := n2.intoSet()
	for _, name := range n {
		if _, has := m[name]; !has {
			return false
		}
	}
	return true
}

func (n NamesSlice) intoSet() map[string]struct{} {
	m := make(map[string]struct{})
	for _, name := range n {
		m[name] = struct{}{}
	}
	return m
}

func (n NamesSlice) String() string {
	return strings.Join(n, ", ")
}

func (n NamesSlice) EqualNamesSlice(n2 NamesSlice) bool {
	if len(n) != len(n2) {
		return false
	}
	left, right := n.GetSorted(), n2.GetSorted()
	for i, name := range left {
		if name != right[i] {
			return false
		}
	}
	return true
}

func (n NamesSlice) EqualTupleAttrs(t Tuple) bool {
	return NewNames(n...).Equal(t.Names())
}

func (n NamesSlice) LessNamesSlice(n2 NamesSlice) bool {
	if len(n) != len(n2) {
		return len(n) < len(n2)
	}
	left, right := n.GetSorted(), n2.GetSorted()
	for i, attr := range left {
		if attr < right[i] {
			return true
		}
		if attr > right[i] {
			return false
		}
	}
	return false
}

func (n NamesSlice) GetSorted() NamesSlice {
	sorted := make(NamesSlice, len(n))
	copy(sorted, n)
	sort.Strings(sorted)
	return sorted
}

// a helper type so that getBucket can return a hashable fmt.Stringer type from a namesSlice.
type hashableNamesSlice string

func newHashableNamesSlice(n NamesSlice) hashableNamesSlice {
	return hashableNamesSlice(n.String())
}

func (s hashableNamesSlice) String() string {
	return string(s)
}

// Count returns how many attributes are in the Tuple.
func (t *GenericTuple) Count() int {
	return len(t.vals)
}

// Get returns the Value associated with a name, and true iff it was found.
func (t *GenericTuple) Get(name string) (Value, bool) {
	if i, ok := t.sh().Index(name); ok {
		return t.vals[i], true
	}
	return nil, false
}

// MustGet returns e.Get(name) or panics if an error occurs.
func (t *GenericTuple) MustGet(name string) Value {
	if v, has := t.Get(name); has {
		return v
	}
	panic(fmt.Errorf("%q not found", name))
}

// With returns a Tuple with all name/Value pairs in t (except the one for the
// given name, if present) with the addition of the given name/Value pair.
func (t *GenericTuple) With(name string, value Value) Tuple {
	// Strip view/non-view counterpart.
	if strings.HasPrefix(name, "&") {
		t = t.without(name[1:])
	} else {
		t = t.without("&" + name)
	}
	shape := t.sh()
	if i, ok := shape.Index(name); ok {
		vals := make([]Value, len(t.vals))
		copy(vals, t.vals)
		vals[i] = value
		return newShapedTuple(shape, vals)
	}
	next, at := shape.With(name)
	vals := make([]Value, 0, len(t.vals)+1)
	vals = append(vals, t.vals[:at]...)
	vals = append(vals, value)
	vals = append(vals, t.vals[at:]...)
	return newShapedTuple(next, vals)
}

// Without returns a Tuple with all name/Value pairs in t exception the one of
// the given name.
func (t *GenericTuple) Without(name string) Tuple {
	return t.without(name)
}

func (t *GenericTuple) without(name string) *GenericTuple {
	shape := t.sh()
	if _, ok := shape.Index(name); !ok {
		return t
	}
	next, at := shape.Without(name)
	vals := make([]Value, 0, len(t.vals)-1)
	vals = append(vals, t.vals[:at]...)
	vals = append(vals, t.vals[at+1:]...)
	return newShapedTuple(next, vals)
}

func (t *GenericTuple) Map(f func(Value) (Value, error)) (Tuple, error) {
	vals := make([]Value, len(t.vals))
	for i, v := range t.vals {
		mapped, err := f(v)
		if err != nil {
			return nil, err
		}
		vals[i] = mapped
	}
	return newShapedTuple(t.sh(), vals), nil
}

// HasName returns true iff the Tuple has an attribute with the given name.
func (t *GenericTuple) HasName(name string) bool {
	_, found := t.sh().Index(name)
	return found
}

// Names returns the attribute names.
func (t *GenericTuple) Names() Names {
	return t.sh().Names()
}

// Project returns a tuple with the given names from this tuple, or nil if any
// name wasn't found.
func (t *GenericTuple) Project(names Names) Tuple {
	var b TupleBuilder
	for e := names.Enumerator(); e.MoveNext(); {
		name := e.Current()
		value, found := t.Get(name)
		if !found {
			return nil
		}
		b.Put(name, value)
	}
	return b.Finish()
}

// GenericTupleEnumerator represents an enumerator over a GenericTuple, in
// attribute-name order.
type GenericTupleEnumerator struct {
	t *GenericTuple
	i int
}

// MoveNext moves the enumerator to the next Value.
func (e *GenericTupleEnumerator) MoveNext() bool {
	e.i++
	return e.i < len(e.t.vals)
}

// Current returns the enumerator's current Value.
func (e *GenericTupleEnumerator) Current() (string, Value) {
	return e.t.sh().names[e.i], e.t.vals[e.i]
}

// Enumerator returns an enumerator over the Values in the GenericTuple.
func (t *GenericTuple) Enumerator() AttrEnumerator {
	return &GenericTupleEnumerator{t: t, i: -1}
}

// TupleOrderedNames returns the names of this tuple in sorted order. The
// slice is shared with the tuple's shape and must not be modified.
func TupleOrderedNames(t *GenericTuple) []string {
	return t.sh().names
}

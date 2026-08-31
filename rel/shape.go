package rel

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/arr-ai/frozen"
	"github.com/arr-ai/hash/hash128"
)

// Shape is the attribute layout of a tuple: its sorted attribute names and
// everything that depends only on them. Tuples are struct-like — a program
// uses a small number of distinct attribute sets millions of times — so
// shapes are interned per process and a tuple carries a *Shape plus its
// values in shape order. Two tuples have the same attribute set iff they
// share the same *Shape, which makes equality, hashing, attribute lookup and
// relation-row inflation positional.
//
// Shapes are immutable once interned; transitions (the shape reached by
// adding or removing one attribute) are memoised on first use.
type Shape struct {
	names []string
	nameH []hash128.H128
	index map[string]int // nil for small shapes, which scan linearly

	// namesH is the xor of the attribute-name hashes, the half of a tuple's
	// hash that depends only on its shape.
	namesH hash128.H128

	// bucket is the SetBuilder bucket key for tuples of this shape, boxed
	// once so getBucket never allocates.
	bucket    fmt.Stringer
	namesOnce sync.Once
	namesSet  Names

	transitions sync.Map // "+name" / "-name" -> *shapeTransition
}

type shapeTransition struct {
	shape *Shape
	at    int
}

// Below this many attributes a linear scan beats a map lookup.
const shapeLinearScanMax = 8

var (
	shapes     sync.Map // key -> *Shape
	emptyShape = shapeOf(nil)
)

// shapeKey uniquely encodes a sorted name list.
func shapeKey(names []string) string {
	var sb strings.Builder
	for _, n := range names {
		sb.WriteString(strconv.Itoa(len(n)))
		sb.WriteByte(':')
		sb.WriteString(n)
	}
	return sb.String()
}

// shapeOf returns the interned Shape for names, which must be sorted and
// free of duplicates. The slice is retained.
func shapeOf(names []string) *Shape {
	key := shapeKey(names)
	if s, ok := shapes.Load(key); ok {
		return s.(*Shape)
	}
	s := &Shape{
		names:  names,
		nameH:  make([]hash128.H128, len(names)),
		bucket: newHashableNamesSlice(names),
	}
	for i, n := range names {
		s.nameH[i] = hash128.String(n)
		s.namesH = s.namesH.Xor(s.nameH[i])
	}
	if len(names) > shapeLinearScanMax {
		s.index = make(map[string]int, len(names))
		for i, n := range names {
			s.index[n] = i
		}
	}
	actual, _ := shapes.LoadOrStore(key, s)
	return actual.(*Shape)
}

// shapeOfAttrs returns the shape for a set of attributes in any order. When
// a name repeats, the last occurrence wins, matching builder semantics.
// It also returns the values in shape order.
func shapeOfAttrs(attrs []Attr) (*Shape, []Value) {
	switch len(attrs) {
	case 0:
		return emptyShape, nil
	case 1:
		return shapeOf([]string{attrs[0].Name}), []Value{attrs[0].Value}
	}
	sorted := make([]Attr, len(attrs))
	copy(sorted, attrs)
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].Name < sorted[j].Name })
	names := make([]string, 0, len(sorted))
	vals := make([]Value, 0, len(sorted))
	for _, a := range sorted {
		if n := len(names); n > 0 && names[n-1] == a.Name {
			vals[n-1] = a.Value // later Put wins
			continue
		}
		names = append(names, a.Name)
		vals = append(vals, a.Value)
	}
	return shapeOf(names), vals
}

// Count returns the number of attributes.
func (s *Shape) Count() int {
	return len(s.names)
}

// Index returns the position of name in the shape.
func (s *Shape) Index(name string) (int, bool) {
	if s.index != nil {
		i, ok := s.index[name]
		return i, ok
	}
	for i, n := range s.names {
		if n == name {
			return i, true
		}
	}
	return 0, false
}

// With returns the shape with name added and the position it occupies
// there. name must not already be present.
func (s *Shape) With(name string) (*Shape, int) {
	key := "+" + name
	if t, ok := s.transitions.Load(key); ok {
		tr := t.(*shapeTransition)
		return tr.shape, tr.at
	}
	at := sort.SearchStrings(s.names, name)
	names := make([]string, 0, len(s.names)+1)
	names = append(names, s.names[:at]...)
	names = append(names, name)
	names = append(names, s.names[at:]...)
	tr := &shapeTransition{shape: shapeOf(names), at: at}
	s.transitions.Store(key, tr)
	return tr.shape, tr.at
}

// Without returns the shape with name removed and the position it had.
// name must be present.
func (s *Shape) Without(name string) (*Shape, int) {
	key := "-" + name
	if t, ok := s.transitions.Load(key); ok {
		tr := t.(*shapeTransition)
		return tr.shape, tr.at
	}
	at, _ := s.Index(name)
	names := make([]string, 0, len(s.names)-1)
	names = append(names, s.names[:at]...)
	names = append(names, s.names[at+1:]...)
	tr := &shapeTransition{shape: shapeOf(names), at: at}
	s.transitions.Store(key, tr)
	return tr.shape, tr.at
}

// Names returns the attribute names as a frozen set, built once per shape.
func (s *Shape) Names() Names {
	s.namesOnce.Do(func() {
		var b frozen.SetBuilder[string]
		for _, n := range s.names {
			b.Add(n)
		}
		s.namesSet = Names(b.Finish())
	})
	return s.namesSet
}

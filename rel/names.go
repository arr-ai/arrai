package rel

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/arr-ai/hash/hash128"
)

// Names is an interned set of attribute names and everything derived only
// from them. The same names always share the same interned identity, so
// equality is a pointer compare and set algebra returns another interned
// Names. A GenericTuple carries a Names plus values in name order.
type Names struct {
	*namesRep
}

type namesRep struct {
	names []string
	nameH []hash128.H128
	index map[string]int // nil for small sets, which scan linearly

	// namesH is the xor of the attribute-name hashes, the half of a tuple's
	// hash that depends only on its attribute set.
	namesH hash128.H128

	// bucket is the SetBuilder bucket key for tuples of this set, boxed
	// once so getBucket never allocates.
	bucket fmt.Stringer

	transitions sync.Map // "+name" / "-name" -> *namesTransition
}

type namesTransition struct {
	names Names
	at    int
}

// Below this many attributes a linear scan beats a map lookup.
const namesLinearScanMax = 8

var internedNames sync.Map // key -> *namesRep

// EmptyNames is the empty set of names.
var EmptyNames = internNames(nil)

var (
	arrayItemNames  = NewNames("@", ArrayItemAttr)
	stringCharNames = NewNames("@", StringCharAttr)
	bytesByteNames  = NewNames("@", BytesByteAttr)
	dictEntryNames  = NewNames("@", DictValueAttr)
)

func (n Names) canon() Names {
	if n.namesRep == nil {
		return EmptyNames
	}
	return n
}

// internNames returns the interned Names for names, which must be sorted and
// free of duplicates. The slice is retained.
func internNames(names []string) Names {
	key := internNamesKey(names)
	if s, ok := internedNames.Load(key); ok {
		return Names{s.(*namesRep)}
	}
	s := &namesRep{
		names:  names,
		nameH:  make([]hash128.H128, len(names)),
		bucket: newHashableNamesSlice(names),
	}
	for i, n := range names {
		s.nameH[i] = hash128.String(n)
		s.namesH = s.namesH.Xor(s.nameH[i])
	}
	if len(names) > namesLinearScanMax {
		s.index = make(map[string]int, len(names))
		for i, n := range names {
			s.index[n] = i
		}
	}
	actual, _ := internedNames.LoadOrStore(key, s)
	return Names{actual.(*namesRep)}
}

func internNamesKey(names []string) string {
	var sb strings.Builder
	for _, n := range names {
		sb.WriteString(strconv.Itoa(len(n)))
		sb.WriteByte(':')
		sb.WriteString(n)
	}
	return sb.String()
}

// internAttrs returns the interned Names for a set of attributes in any
// order. When a name repeats, the last occurrence wins, matching builder
// semantics. It also returns the values in name order.
func internAttrs(attrs []Attr) (Names, []Value) {
	switch len(attrs) {
	case 0:
		return EmptyNames, nil
	case 1:
		return internNames([]string{attrs[0].Name}), []Value{attrs[0].Value}
	}
	sorted := make([]Attr, len(attrs))
	copy(sorted, attrs)
	slices.SortStableFunc(sorted, func(a, b Attr) int { return strings.Compare(a.Name, b.Name) })
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
	return internNames(names), vals
}

// Index returns the position of name in the set.
func (s *namesRep) Index(name string) (int, bool) {
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

// insert returns the set with name added and the position it occupies
// there. name must not already be present.
func (s *namesRep) insert(name string) (Names, int) {
	key := "+" + name
	if t, ok := s.transitions.Load(key); ok {
		tr := t.(*namesTransition)
		return tr.names, tr.at
	}
	at, _ := slices.BinarySearch(s.names, name)
	names := make([]string, 0, len(s.names)+1)
	names = append(names, s.names[:at]...)
	names = append(names, name)
	names = append(names, s.names[at:]...)
	tr := &namesTransition{names: internNames(names), at: at}
	s.transitions.Store(key, tr)
	return tr.names, tr.at
}

// remove returns the set with name removed and the position it had.
// name must be present.
func (s *namesRep) remove(name string) (Names, int) {
	key := "-" + name
	if t, ok := s.transitions.Load(key); ok {
		tr := t.(*namesTransition)
		return tr.names, tr.at
	}
	at, _ := s.Index(name)
	names := make([]string, 0, len(s.names)-1)
	names = append(names, s.names[:at]...)
	names = append(names, s.names[at+1:]...)
	tr := &namesTransition{names: internNames(names), at: at}
	s.transitions.Store(key, tr)
	return tr.names, tr.at
}

// NewNames returns the interned attribute set for the given names. Order and
// duplicates do not matter.
func NewNames(names ...string) Names {
	switch len(names) {
	case 0:
		return EmptyNames
	case 1:
		return internNames([]string{names[0]})
	}
	sorted := append([]string(nil), names...)
	slices.Sort(sorted)
	w := 1
	for i := 1; i < len(sorted); i++ {
		if sorted[i] != sorted[w-1] {
			sorted[w] = sorted[i]
			w++
		}
	}
	return internNames(sorted[:w])
}

// Names returns a slice of the names in the set, in sorted order.
func (n Names) Names() []string {
	return slices.Clone(n.canon().names)
}

// OrderedNames returns a sorted slice of the names in the set.
func (n Names) OrderedNames() []string {
	return n.Names()
}

// IsTrue returns true iff there are names in the set.
func (n Names) IsTrue() bool {
	return n.Count() != 0
}

// Count returns the number of names in the set.
func (n Names) Count() int {
	return len(n.canon().names)
}

// Hash computes a hash value for the set of names.
func (n Names) Hash(seed uint32) uint32 {
	return uint32(n.canon().namesH.Seeded(uintptr(seed) + 0x4e351c91))
}

// Equal returns true iff the given sets of names are equal.
func (n Names) Equal(i interface{}) bool {
	x, ok := i.(Names)
	return ok && n.canon() == x.canon()
}

// String returns a string representation of the set of names.
func (n Names) String() string {
	return fmt.Sprintf("|%s|", strings.Join(n.canon().names, ", "))
}

// With returns a set with all the input names and the given name.
func (n Names) With(name string) Names {
	n = n.canon()
	if _, ok := n.Index(name); ok {
		return n
	}
	next, _ := n.insert(name)
	return next
}

// Without returns a set with all the input names, excluding the given name.
func (n Names) Without(name string) Names {
	n = n.canon()
	if _, ok := n.Index(name); !ok {
		return n
	}
	next, _ := n.remove(name)
	return next
}

// Has returns true iff the given name is in the set of names.
func (n Names) Has(name string) bool {
	_, ok := n.canon().Index(name)
	return ok
}

// Any returns an arbitrary element from `n`.
func (n Names) Any() string {
	names := n.canon().names
	if len(names) == 0 {
		panic("Names.Any on empty set")
	}
	return names[0]
}

// Enumerator returns an enumerator over a set of names, in sorted order.
func (n Names) Enumerator() *NamesEnumerator {
	return &NamesEnumerator{names: n.canon().names, i: -1}
}

// TheOne return the single name in the set; panics otherwise.
func (n Names) TheOne() string {
	names := n.canon().names
	if len(names) != 1 {
		panic("Names.TheOne expects exactly one name in the set")
	}
	return names[0]
}

// ToSlice returns a slice of the names in the set, in sorted order.
func (n Names) ToSlice() []string {
	return n.Names()
}

// Intersect returns names in both sets.
func (n Names) Intersect(o Names) Names {
	a, b := n.canon(), o.canon()
	if a == b {
		return a
	}
	if len(a.names) == 0 || len(b.names) == 0 {
		return EmptyNames
	}
	names := make([]string, 0, min(len(a.names), len(b.names)))
	i, j := 0, 0
	for i < len(a.names) && j < len(b.names) {
		switch strings.Compare(a.names[i], b.names[j]) {
		case 0:
			names = append(names, a.names[i])
			i++
			j++
		case -1:
			i++
		default:
			j++
		}
	}
	return internNames(names)
}

// Minus returns names in one set not found in the other.
func (n Names) Minus(o Names) Names {
	a, b := n.canon(), o.canon()
	if a == b || len(a.names) == 0 {
		return EmptyNames
	}
	if len(b.names) == 0 {
		return a
	}
	names := make([]string, 0, len(a.names))
	i, j := 0, 0
	for i < len(a.names) && j < len(b.names) {
		switch strings.Compare(a.names[i], b.names[j]) {
		case 0:
			i++
			j++
		case -1:
			names = append(names, a.names[i])
			i++
		default:
			j++
		}
	}
	names = append(names, a.names[i:]...)
	return internNames(names)
}

// IsSubsetOf returns true if `n` is a subset of `o`.
func (n Names) IsSubsetOf(o Names) bool {
	a, b := n.canon(), o.canon()
	if a == b || len(a.names) == 0 {
		return true
	}
	if len(a.names) > len(b.names) {
		return false
	}
	i, j := 0, 0
	for i < len(a.names) && j < len(b.names) {
		switch strings.Compare(a.names[i], b.names[j]) {
		case 0:
			i++
			j++
		case -1:
			return false
		default:
			j++
		}
	}
	return i == len(a.names)
}

// NamesEnumerator represents an enumerator over a Names set, in sorted order.
type NamesEnumerator struct {
	names []string
	i     int
}

// MoveNext moves the enumerator to the next name.
func (e *NamesEnumerator) MoveNext() bool {
	e.i++
	return e.i < len(e.names)
}

// Current returns the enumerator's current name.
func (e *NamesEnumerator) Current() string {
	return e.names[e.i]
}

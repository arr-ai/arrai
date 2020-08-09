package rel

import (
	"fmt"
	"sort"
	"strings"

	"github.com/arr-ai/frozen"
)

// Names represents a set of names.
type Names frozen.Set

// EmptyNames is the empty set of names.
var EmptyNames = Names(frozen.Set{})

// NewNames returns a new set of names with the given names.
func NewNames(names ...string) Names {
	s := frozen.Set{}
	for _, name := range names {
		s = s.With(name)
	}
	return Names(s)
}

// Names returns a slice of the names in the set.
func (n Names) Names() []string {
	names := make([]string, 0, n.Count())
	for e := n.Enumerator(); e.MoveNext(); {
		names = append(names, e.Current())
	}
	return names
}

// Names returns a sorted slice of the names in the set.
func (n Names) OrderedNames() []string {
	names := n.Names()
	sort.Strings(names)
	return names
}

// Bool returns true iff there are names in the set.
func (n Names) IsTrue() bool {
	return n.Count() != 0
}

// Count returns the number of names in a set of names.
func (n Names) Count() int {
	return (frozen.Set(n)).Count()
}

// Hash computes a hash value for the set of names.
func (n Names) Hash(seed uint32) uint32 {
	return uint32((frozen.Set(n)).Hash(uintptr(seed) + 0x4e351c91))
}

// Equal returns true iff the given sets of names are equal.
func (n Names) Equal(i interface{}) bool {
	if x, ok := i.(Names); ok {
		return (frozen.Set(n)).Equal(frozen.Set(x))
	}
	return false
}

// String returns a string representation of the set of names.
func (n Names) String() string {
	return fmt.Sprintf("|%s|", strings.Join(n.OrderedNames(), ", "))
}

// With returns a set with all the input names and the given name.
func (n Names) With(name string) Names {
	return Names((frozen.Set(n)).With(name))
}

// Without returns a set with all the input names, excluding the given name.
func (n Names) Without(name string) Names {
	return Names((frozen.Set(n)).Without(name))
}

// Has returns true iff the given name is in the set of names.
func (n Names) Has(name string) bool {
	return (frozen.Set(n)).Has(name)
}

// Any returns an arbitrary element from `n`.
func (n Names) Any() string {
	return (frozen.Set(n)).Any().(string)
}

// Enumerator returns an enumerator over a set of names.
func (n Names) Enumerator() *NamesEnumerator {
	return &NamesEnumerator{(frozen.Set(n)).Range()}
}

// TheOne return the single name in the set; panics otherwise.
func (n Names) TheOne() string {
	if n.Count() != 1 {
		panic("Names.TheOne expects exactly one name in the set")
	}
	e := n.Enumerator()
	e.MoveNext()
	return e.Current()
}

// ToSlice returns a slice of the names in the set.
func (n Names) ToSlice() []string {
	names := make([]string, n.Count())
	i := 0
	for e := n.Enumerator(); e.MoveNext(); {
		names[i] = e.Current()
		i++
	}
	return names
}

// Intersect returns names in both sets.
func (n Names) Intersect(o Names) Names {
	return Names((frozen.Set(n)).Intersection(frozen.Set(o)))
}

// Minus returns names in one set not found in the other.
func (n Names) Minus(o Names) Names {
	return Names((frozen.Set(n)).Difference(frozen.Set(o)))
}

// IsSubsetOf returns true if `n`` is a subset of `o`.
func (n Names) IsSubsetOf(o Names) bool {
	return (frozen.Set(n)).IsSubsetOf(frozen.Set(o))
}

// NamesEnumerator represents an enumerator over a set of names.
type NamesEnumerator struct {
	i frozen.Iterator
}

// MoveNext moves the enumerator to the next Value.
func (e *NamesEnumerator) MoveNext() bool {
	return e.i.Next()
}

// Current returns the enumerator's current Value.
func (e *NamesEnumerator) Current() string {
	return e.i.Value().(string)
}

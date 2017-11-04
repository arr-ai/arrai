package rel

import (
	"bytes"

	"github.com/mediocregopher/seq"
)

// Names represents a set of names.
type Names seq.Set

// EmptyNames is the empty set of names.
var EmptyNames = (*Names)(seq.NewSet())

// NewNames returns a new set of names with the given names.
func NewNames(names ...string) *Names {
	s := seq.NewSet()
	for _, name := range names {
		s, _ = s.SetVal(name)
	}
	return (*Names)(s)
}

// Bool returns true iff there are names in the set.
func (n *Names) Bool() bool {
	return n.Count() != 0
}

// Count returns the number of names in a set of names.
func (n *Names) Count() uint64 {
	return (*seq.Set)(n).Size()
}

// Hash computes a hash value for the set of names.
func (n *Names) Hash(seed uint32) uint32 {
	return (*seq.Set)(n).Hash(seed + 0x4e351c91)
}

// Equal returns true iff the given sets of names are equal.
func (n *Names) Equal(i interface{}) bool {
	if x, ok := i.(*Names); ok {
		return (*seq.Set)(n).Equal((*seq.Set)(x))
	}
	return false
}

// String returns a string representation of the set of names.
func (n *Names) String() string {
	var buf bytes.Buffer
	buf.WriteRune('|')
	i := 0
	for e := n.Enumerator(); e.MoveNext(); {
		if i != 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(e.Current())
	}
	buf.WriteRune('|')
	return buf.String()
}

// With returns a set with all the input names and the given name.
func (n *Names) With(name string) *Names {
	if s, added := (*seq.Set)(n).SetVal(name); added {
		return (*Names)(s)
	}
	return n
}

// Without returns a set with all the input names, excluding the given name.
func (n *Names) Without(name string) *Names {
	if s, deleted := (*seq.Set)(n).DelVal(name); deleted {
		return (*Names)(s)
	}
	return n
}

// Has returns true iff the given name is in the set of names.
func (n *Names) Has(name string) bool {
	_, found := (*seq.Set)(n).GetVal(name)
	return found
}

// Enumerator returns an enumerator over a set of names.
func (n *Names) Enumerator() *NamesEnumerator {
	return &NamesEnumerator{seq.Seq((*seq.Set)(n)), ""}
}

// TheOne return the single name in the set; panics otherwise.
func (n *Names) TheOne() string {
	if n.Count() != 1 {
		panic("Names.TheOne expects exactly one name in the set")
	}
	e := n.Enumerator()
	e.MoveNext()
	return e.Current()
}

// ToSlice returns a slice of the names in the set.
func (n *Names) ToSlice() []string {
	names := make([]string, n.Count())
	i := 0
	for e := n.Enumerator(); e.MoveNext(); {
		names[i] = e.Current()
		i++
	}
	return names
}

// Intersect returns names in both sets.
func (n *Names) Intersect(rhs *Names) *Names {
	result := EmptyNames
	for e := n.Enumerator(); e.MoveNext(); {
		name := e.Current()
		if rhs.Has(name) {
			result = result.With(name)
		}
	}
	return result
}

// Minus returns names in one set not found in the other.
func (n *Names) Minus(rhs *Names) *Names {
	result := EmptyNames
	for e := n.Enumerator(); e.MoveNext(); {
		name := e.Current()
		if !rhs.Has(name) {
			result = result.With(name)
		}
	}
	return result
}

// IsSubsetOf returns true iff the set is a subset of another.
func (n *Names) IsSubsetOf(rhs *Names) bool {
	return !n.Minus(rhs).Bool()
}

// NamesEnumerator represents an enumerator over a set of names.
type NamesEnumerator struct {
	seq     seq.Seq
	current string
}

// MoveNext moves the enumerator to the next Value.
func (e *NamesEnumerator) MoveNext() bool {
	var first interface{}
	var ok bool
	first, e.seq, ok = e.seq.FirstRest()
	if ok {
		e.current = first.(string)
	} else {
		e.current = ""
	}
	return ok
}

// Current returns the enumerator's current Value.
func (e *NamesEnumerator) Current() string {
	return e.current
}

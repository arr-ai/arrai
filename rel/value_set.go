package rel

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sort"

	"github.com/arr-ai/frozen"
)

type setFlavor int

// Flavors of Set
const (
	setFlavorNormal setFlavor = iota
	setFlavorArray
	setFlavorString
)

// genericSet is a set of Values.
// It is implemented as an immutable map[uintptr]set[Value].
type genericSet struct {
	set    frozen.Set
	flavor setFlavor
}

// genericSet equivalents for Boolean true and false
var (
	None  = Set(&genericSet{frozen.Set{}, setFlavorNormal})
	False = None
	True  = None.With(EmptyTuple)
)

// NewSet constructs a genericSet from a set of Values.
func NewSet(values ...Value) Set {
	set := None
	for _, value := range values {
		set = set.With(value)
	}
	return set
}

// NewSetFrom constructs a genericSet from interfaces.
func NewSetFrom(intfs ...interface{}) (Set, error) {
	set := None
	for _, intf := range intfs {
		value, err := NewValue(intf)
		if err != nil {
			return nil, err
		}
		set = set.With(value)
	}
	return set, nil
}

func newSetFromSet(s Set) Set {
	set := None
	for e := s.Enumerator(); e.MoveNext(); {
		set = set.With(e.Current())
	}
	return set
}

// NewBool constructs a bool as a relation.
func NewBool(b bool) Set {
	if b {
		return True
	}
	return False
}

// Hash computes a hash for a genericSet.
func (s *genericSet) Hash(seed uintptr) uintptr {
	h := seed
	for e := s.Enumerator(); e.MoveNext(); {
		h ^= e.Current().Hash(0)
	}
	return h
}

// Equal tests two Sets for equality. Any other type returns false.
func (s *genericSet) Equal(v interface{}) bool {
	if t, ok := v.(*genericSet); ok {
		return s.set.Equal(t.set)
	}
	return false
}

// String returns a string representation of a genericSet.
func (s *genericSet) String() string {
	switch s.flavor {
	case setFlavorNormal:
		// {} == none
		if !s.Bool() {
			return "{}"
		}

		// {()} == true
		if s.Count() == 1 {
			e := s.Enumerator()
			e.MoveNext()
			if tuple, ok := e.Current().(Tuple); ok && !tuple.Bool() {
				return "true"
			}
		}

		var buf bytes.Buffer
		buf.WriteString("{")
		for i, value := range s.OrderedValues() {
			if i != 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(value.String())
		}
		buf.WriteString("}")
		return buf.String()

	case setFlavorArray:
		var buf bytes.Buffer
		buf.WriteString("[")
		for i, tuple := range s.OrderedValues() {
			if i > 0 {
				buf.WriteString(", ")
			}
			if tuple, ok := tuple.(Tuple); ok {
				if value, found := tuple.Get(ArrayItemAttr); found {
					buf.WriteString(value.String())
				} else {
					panic("Array item tuple must have @item attr")
				}
			} else {
				panic("Array item must be a Tuple")
			}
		}
		buf.WriteString("]")
		return buf.String()

	case setFlavorString:
		runes := make([]rune, s.Count())
		for i, tuple := range s.OrderedValues() {
			if tuple, ok := tuple.(Tuple); ok {
				if value, found := tuple.Get(CharAttr); found {
					if n, ok := value.(Number); ok {
						runes[i] = rune(n.Float64())
					} else {
						panic("String tuple " + CharAttr +
							" attr must be a number")
					}
				} else {
					panic("String char tuple must have a " + CharAttr +
						" attr")
				}
			} else {
				panic("String char must be a Tuple")
			}
		}
		j, err := json.Marshal(string(runes))
		if err != nil {
			panic(err)
		}
		return string(j)
	default:
		panic("Unhandled set flavor")
	}
}

// Eval returns the set.
func (s *genericSet) Eval(local Scope) (Value, error) {
	return s, nil
}

var genericSetKind = registerKind(200, reflect.TypeOf(Function{}))

// Kind returns a number that is unique for each major kind of Value.
func (s *genericSet) Kind() int {
	return genericSetKind
}

// Bool returns true iff the tuple has attributes.
func (s *genericSet) Bool() bool {
	return s.Count() > 0
}

// Less returns true iff v.Kind() < genericSet.Kind() or v is a
// genericSet and t precedes v in a lexicographical comparison of their
// sorted values.
func (s *genericSet) Less(v Value) bool {
	if s.Kind() != v.Kind() {
		return s.Kind() < v.Kind()
	}

	x := v.(*genericSet)
	a := s.OrderedValues()
	b := x.OrderedValues()
	n := len(a)
	if n > len(b) {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		if a[i].Less(b[i]) {
			return true
		}
		if b[i].Less(a[i]) {
			return false
		}
	}
	return len(a) < len(b)
}

// Negate returns {(negateTag): s}.
func (s *genericSet) Negate() Value {
	if !s.Bool() {
		return s
	}
	return NewTuple(NewAttr(negateTag, s))
}

// Export exports a genericSet as an array of exported Values.
func (s *genericSet) Export() interface{} {
	result := make([]interface{}, 0, s.set.Count())
	for e := s.Enumerator(); e.MoveNext(); {
		result = append(result, e.Current().Export())
	}
	return result
}

// Count returns the number of elements in the genericSet.
func (s *genericSet) Count() int {
	return s.set.Count()
}

// Has returns true iff the given Value is in the genericSet.
func (s *genericSet) Has(value Value) bool {
	return s.set.Has(value)
}

// With returns the original genericSet with given value added. Iff the value was
// already present, the original genericSet is returned.
func (s *genericSet) With(value Value) Set {
	set := s.set.With(value)

	var isArrayAttr bool
	var isStringAttr bool
	if tuple, ok := value.(Tuple); ok {
		if tuple.Count() == 2 {
			if at, found := tuple.Get("@"); found {
				if _, ok := at.(Number); ok {
					if _, found := tuple.Get(ArrayItemAttr); found {
						isArrayAttr = true
					} else if char, found := tuple.Get(CharAttr); found {
						if number, ok := char.(Number); ok {
							n := number.Float64()
							if n == float64(int64(n)) {
								isStringAttr = true
							}
						}
					}
				}
			}
		}
	}
	flavor := setFlavorNormal
	if (s.flavor == setFlavorArray || !s.Bool()) && isArrayAttr {
		flavor = setFlavorArray
	} else if (s.flavor == setFlavorString || !s.Bool()) && isStringAttr {
		flavor = setFlavorString
	}

	return &genericSet{set, flavor}
}

// Without returns the original genericSet without the given value. Iff the value was
// already absent, the original genericSet is returned.
func (s *genericSet) Without(value Value) Set {
	set := s.set.Without(value)
	if set == s.set {
		return s
	}
	flavor := s.flavor
	if set.Count() == 0 {
		flavor = setFlavorNormal
	}
	return &genericSet{set, flavor}
}

// Map maps values per f.
func (s *genericSet) Map(f func(v Value) Value) Set {
	result := NewSet()
	for e := s.Enumerator(); e.MoveNext(); {
		result = result.With(f(e.Current()))
	}
	return result
}

// Where returns a new genericSet with all the Values satisfying predicate p.
func (s genericSet) Where(p func(v Value) bool) Set {
	s.set = s.set.Where(func(elem interface{}) bool { return p(elem.(Value)) })
	return &s
}

// Call ...
func (s *genericSet) Call(arg Value) Value {
	for e := s.Enumerator(); e.MoveNext(); {
		var at Value
		var t Tuple
		if NewTupleMatcher(map[string]Matcher{"@": Bind(&at)}, Bind(&t)).Match(e.Current()) && at.Equal(arg) {
			for attr := t.Enumerator(); attr.MoveNext(); {
				_, value := attr.Current()
				return value
			}
		}
	}
	return nil
}

// Enumerator returns an enumerator over the Values in the genericSet.
func (s *genericSet) Enumerator() ValueEnumerator {
	return &genericSetEnumerator{s.set.Range()}
}

// Any return any value from the set.
func (s *genericSet) Any() Value {
	for e := s.Enumerator(); e.MoveNext(); {
		return e.Current()
	}
	panic("Any(): empty set")
}

// genericSetEnumerator represents an enumerator over a genericSet.
type genericSetEnumerator struct {
	i frozen.Iterator
}

// MoveNext moves the enumerator to the next Value.
func (e *genericSetEnumerator) MoveNext() bool {
	return e.i.Next()
}

// Current returns the enumerator's current Value.
func (e *genericSetEnumerator) Current() Value {
	return e.i.Value().(Value)
}

// ValueList represents a []Value for use in sort.Sort().
type ValueList []Value

func (vl ValueList) Len() int {
	return len(vl)
}

func (vl ValueList) Less(i, j int) bool {
	return vl[i].Less(vl[j])
}

func (vl ValueList) Swap(i, j int) {
	vl[i], vl[j] = vl[j], vl[i]
}

// OrderedValues returns Values in a canonical ordering.
func (s *genericSet) OrderedValues() []Value {
	a := make([]Value, s.Count())
	i := 0
	for e := s.Enumerator(); e.MoveNext(); {
		a[i] = e.Current()
		i++
	}
	sort.Sort(ValueList(a))
	return a
}

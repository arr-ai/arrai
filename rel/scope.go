package rel

import (
	"bytes"
	"sort"

	"github.com/go-errors/errors"
	"github.com/mediocregopher/seq"
)

// EmptyScope is the scope with no variables.
var EmptyScope = &Scope{seq.NewHashMap()}

// Scope represents an expression scope.
type Scope struct {
	hmap *seq.HashMap
}

func (s *Scope) String() string {
	var buf bytes.Buffer
	buf.WriteRune('{')
	for i, name := range s.orderedNames() {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(name)
		buf.WriteString(": ")
		expr, found := s.Get(name)
		if !found {
			panic(errors.Errorf(
				"Scope iteration produced name %v, which fails lookup", name))
		}
		if expr != nil {
			buf.WriteString(expr.String())
		} else {
			buf.WriteString("nil")
		}
	}
	buf.WriteRune('}')
	return buf.String()
}

// Eval evaluates an expression in a given symbol environment.
func (s *Scope) Eval(local, global *Scope) (Value, error) {
	tuple := NewTuple()
	for e := s.Enumerator(); e.MoveNext(); {
		name, expr := e.Current()
		value, err := expr.Eval(local, global)
		if err != nil {
			return nil, err
		}
		tuple, _ = tuple.With(name, value)
	}
	return tuple, nil
}

// Count returns the number of variables in this Scope.
func (s *Scope) Count() uint64 {
	return s.hmap.Size()
}

// Get returns the Expr for the given name or nil.
func (s *Scope) Get(name string) (Expr, bool) {
	if expr, found := s.hmap.Get(name); found {
		if expr != nil {
			return expr.(Expr), true
		}
		return nil, true
	}
	return nil, false
}

// With returns a new scope with all the old bindings and a new or replacement
// binding for the given name to the given Expr.
func (s *Scope) With(name string, expr Expr) *Scope {
	hmap, _ := s.hmap.Set(name, expr)
	return &Scope{hmap}
}

// Project returns a new scope with just names from the input scope.
func (s *Scope) Project(names *Names) (*Scope, error) {
	result := EmptyScope
	for e := names.Enumerator(); e.MoveNext(); {
		name := e.Current()
		if expr, found := s.Get(name); found {
			result = result.With(name, expr)
		} else {
			return nil, errors.Errorf(
				"name %q not found in scope.Project", name)
		}
	}
	return result, nil
}

// Names returns the attribute names as a slice.
func (s *Scope) Names() []string {
	names := make([]string, s.Count())
	i := 0
	for e := s.Enumerator(); e.MoveNext(); {
		names[i], _ = e.Current()
		i++
	}
	return names
}

// Enumerator returns an enumerator over the Values in the Scope.
func (s *Scope) Enumerator() *ScopeEnumerator {
	return &ScopeEnumerator{s: s.hmap}
}

// orderedNames returns the names of this tuple in sorted order.
func (s *Scope) orderedNames() []string {
	names := s.Names()
	sort.Strings(names)
	return names
}

// ScopeEnumerator represents an enumerator over a Scope.
type ScopeEnumerator struct {
	s    seq.Seq
	name string
	expr Expr
}

// MoveNext moves the enumerator to the next Value.
func (e *ScopeEnumerator) MoveNext() bool {
	item, s, ok := e.s.FirstRest()
	e.s = s
	if !ok {
		return false
	}
	kv := item.(*seq.KV)
	e.name = kv.Key.(string)
	e.expr = kv.Val.(Expr)
	return true
}

// Current returns the enumerator's current Value.
func (e *ScopeEnumerator) Current() (string, Expr) {
	return e.name, e.expr
}

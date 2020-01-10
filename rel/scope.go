package rel

import (
	"bytes"
	"sort"

	"github.com/arr-ai/frozen"
	"github.com/go-errors/errors"
)

// EmptyScope is the scope with no variables.
var EmptyScope = &Scope{}

// Scope represents an expression scope.
type Scope struct {
	m frozen.Map
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
		tuple = tuple.With(name, value)
	}
	return tuple, nil
}

// Count returns the number of variables in this Scope.
func (s *Scope) Count() int {
	return s.m.Count()
}

// Get returns the Expr for the given name or nil.
func (s *Scope) Get(name string) (Expr, bool) {
	if expr, found := s.m.Get(name); found {
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
	return &Scope{s.m.With(name, expr)}
}

// Project returns a new scope with just names from the input scope.
func (s *Scope) Project(names Names) (*Scope, error) {
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
	return &ScopeEnumerator{i: s.m.Range()}
}

// orderedNames returns the names of this tuple in sorted order.
func (s *Scope) orderedNames() []string {
	names := s.Names()
	sort.Strings(names)
	return names
}

// ScopeEnumerator represents an enumerator over a Scope.
type ScopeEnumerator struct {
	i *frozen.MapIterator
}

// MoveNext moves the enumerator to the next Value.
func (e *ScopeEnumerator) MoveNext() bool {
	return e.i.Next()
}

// Current returns the enumerator's current Value.
func (e *ScopeEnumerator) Current() (string, Expr) {
	name, expr := e.i.Entry()
	return name.(string), expr.(Expr)
}

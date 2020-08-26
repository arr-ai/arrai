package rel

import (
	"bytes"
	"context"
	"fmt"
	"sort"

	"github.com/arr-ai/frozen"
	"github.com/go-errors/errors"
)

// EmptyScope is the scope with no variables.
var EmptyScope Scope

// Scope represents an expression scope.
type Scope struct {
	m frozen.Map
}

func (s Scope) String() string {
	var buf bytes.Buffer
	buf.WriteRune('{')
	for i, name := range s.OrderedNames() {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(name)
		buf.WriteString(": ")
		expr, found := s.Get(name)
		if !found {
			panic(errors.Errorf("Scope iteration produced name %v, which fails lookup", name))
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
func (s Scope) Eval(ctx context.Context, local Scope) (Value, error) {
	tuple := NewTuple()
	for e := s.Enumerator(); e.MoveNext(); {
		name, expr := e.Current()
		value, err := expr.Eval(ctx, local)
		if err != nil {
			return nil, WrapContextErr(err, expr, local)
		}
		tuple = tuple.With(name, value)
	}
	return tuple, nil
}

// Count returns the number of variables in this Scope.
func (s Scope) Count() int {
	return s.m.Count()
}

// Get returns the Expr for the given name or nil.
func (s Scope) Get(name string) (Expr, bool) {
	if expr, found := s.m.Get(name); found {
		if expr != nil {
			return expr.(Expr), true
		}
		return nil, true
	}
	return nil, false
}

// MustGet returns the Expr for the given name or panics if not found.
func (s Scope) MustGet(name string) Expr {
	if expr, has := s.Get(name); has {
		return expr
	}
	panic(fmt.Errorf("name not found: %q", name))
}

// With returns a new scope with all the old bindings and a new or replacement
// binding for the given name to the given Expr.
func (s Scope) With(name string, expr Expr) Scope {
	if name == "_" {
		return s
	}
	return Scope{s.m.With(name, expr)}
}

// MatchedWith returns a new scope. New keys are added as With,
// but existing keys fail unless the new value equals the existing value
func (s Scope) MatchedWith(name string, expr Expr) (Scope, error) {
	if name == "_" {
		return s, nil
	}

	if v, exists := s.Get(name); exists {
		if v.String() != expr.String() {
			return Scope{}, fmt.Errorf("%s is redefined differently %s vs %s", name, v, expr)
		}
	}

	return s.With(name, expr), nil
}

// Without returns a new scope with with all the old bindings except the ones
// that correspond to the provided names.
func (s Scope) Without(name ...string) Scope {
	return Scope{s.m.Without(frozen.NewSetFromStrings(name...))}
}

// s.Update(t) merges s and t, choosing t's binding in the event of a name clash.
// It's like calling s.With(t0).With(t1).With(t2)... for each element of t
func (s Scope) Update(t Scope) Scope {
	return Scope{m: s.m.Update(t.m)}
}

// MatchedUpdate merges s and t. New keys are added as Update,
// but existing keys fail unless the new value equals the existing value
func (s Scope) MatchedUpdate(t Scope) (Scope, error) {
	t = t.Without("_")
	for e := s.Enumerator(); e.MoveNext(); {
		name, v := e.Current()
		if expr, exists := t.Get(name); exists {
			if expr.String() != v.String() {
				return Scope{}, fmt.Errorf("the value of %s is different in both scopes", name)
			}
		}
	}

	return s.Update(t), nil
}

// Project returns a new scope with just names from the input scope.
func (s Scope) Project(names Names) (Scope, error) {
	result := EmptyScope
	for e := names.Enumerator(); e.MoveNext(); {
		name := e.Current()
		if expr, found := s.Get(name); found {
			result = result.With(name, expr)
		} else {
			return Scope{}, errors.Errorf(
				"name %q not found in scope.Project", name)
		}
	}
	return result, nil
}

// Names returns the attribute names as a slice.
func (s Scope) Names() []string {
	names := make([]string, s.Count())
	i := 0
	for e := s.Enumerator(); e.MoveNext(); {
		names[i], _ = e.Current()
		i++
	}
	return names
}

// Enumerator returns an enumerator over the Values in the Scope.
func (s Scope) Enumerator() *ScopeEnumerator {
	return &ScopeEnumerator{i: s.m.Range()}
}

// OrderedNames returns the names of this tuple in sorted order.
func (s Scope) OrderedNames() []string {
	names := s.Names()
	sort.Strings(names)
	return names
}

// ScopeEnumerator represents an enumerator over a Scope.
type ScopeEnumerator struct {
	i *frozen.MapIterator
}

// MoveNext moves the enumerator to the next Value.
func (e ScopeEnumerator) MoveNext() bool {
	return e.i.Next()
}

// Current returns the enumerator's current Value.
func (e ScopeEnumerator) Current() (string, Expr) {
	name, expr := e.i.Entry()
	return name.(string), expr.(Expr)
}

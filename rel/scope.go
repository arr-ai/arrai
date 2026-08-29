package rel

import (
	"bytes"
	"context"
	"fmt"
	"sort"

	"github.com/go-errors/errors"
)

// EmptyScope is the scope with no variables.
var EmptyScope Scope

// Scope represents an expression scope: a chain of immutable frames, newest
// first. Arr.ai scoping is lexical, so every binding construct (a call, a
// let, a pattern match, a cond arm) pushes exactly one frame of a fixed
// shape, and a given identifier always finds its binding the same number of
// frames up in the same slot. IdentExpr exploits that with an inline cache;
// the chain itself is the only source of truth.
//
// Frames are never mutated after construction. Closures capture a Scope by
// value (a frame pointer), so chains are shared freely between the scope a
// closure was created in and the scopes derived from it.
type Scope struct {
	f *frame
}

// frame is one link of a scope chain. names and vals are parallel; within a
// frame a later entry shadows an earlier one with the same name, and any
// entry shadows the whole parent chain.
type frame struct {
	parent *frame
	names  []string
	vals   []Expr
}

// lookup finds the newest binding of name, returning the frame, the slot and
// the number of parent hops taken.
func (s Scope) lookup(name string) (f *frame, slot, hops int) {
	for f = s.f; f != nil; f, hops = f.parent, hops+1 {
		for i := len(f.names) - 1; i >= 0; i-- {
			if f.names[i] == name {
				return f, i, hops
			}
		}
	}
	return nil, 0, 0
}

// at returns the frame hops links up the chain, or nil.
func (s Scope) at(hops int) *frame {
	f := s.f
	for ; hops > 0 && f != nil; hops-- {
		f = f.parent
	}
	return f
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
	return len(s.Names())
}

// Get returns the Expr for the given name or nil.
func (s Scope) Get(name string) (Expr, bool) {
	if f, slot, _ := s.lookup(name); f != nil {
		return f.vals[slot], true
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
	return Scope{&frame{parent: s.f, names: []string{name}, vals: []Expr{expr}}}
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
	drop := make(map[string]struct{}, len(name))
	for _, n := range name {
		drop[n] = struct{}{}
	}
	var names []string
	var vals []Expr
	for e := s.Enumerator(); e.MoveNext(); {
		n, v := e.Current()
		if _, dropped := drop[n]; !dropped {
			names = append(names, n)
			vals = append(vals, v)
		}
	}
	if len(names) == 0 {
		return EmptyScope
	}
	return Scope{&frame{names: names, vals: vals}}
}

// s.Update(t) merges s and t, choosing t's binding in the event of a name clash.
// It's like calling s.With(t0).With(t1).With(t2)... for each element of t
func (s Scope) Update(t Scope) Scope {
	if t.f == nil {
		return s
	}
	if s.f == nil {
		return t
	}
	// t is typically a small scope built by a pattern; flatten it into one
	// frame so its bindings occupy fixed slots above s.
	names, vals := t.flatten()
	return Scope{&frame{parent: s.f, names: names, vals: vals}}
}

// flatten returns the scope's bindings oldest-first, one entry per name.
func (s Scope) flatten() ([]string, []Expr) {
	if s.f != nil && s.f.parent == nil {
		return s.f.names, s.f.vals
	}
	var names []string
	var vals []Expr
	for e := s.Enumerator(); e.MoveNext(); {
		n, v := e.Current()
		names = append(names, n)
		vals = append(vals, v)
	}
	return names, vals
}

// MatchedUpdate merges s and t. New keys are added as Update,
// but existing keys fail unless the new value equals the existing value
func (s Scope) MatchedUpdate(t Scope) (Scope, error) {
	// With never binds "_", so no Without("_") pass is needed. t is usually
	// a single small frame from a pattern; check its bindings against s
	// directly rather than enumerating s.
	names, vals := t.flatten()
	for i, name := range names {
		if v, exists := s.Get(name); exists && v.String() != vals[i].String() {
			return Scope{}, fmt.Errorf("the value of %s is different in both scopes", name)
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
	var names []string
	for e := s.Enumerator(); e.MoveNext(); {
		name, _ := e.Current()
		names = append(names, name)
	}
	return names
}

// Enumerator returns an enumerator over the Values in the Scope. Bindings
// are visited oldest-first; a shadowed binding is not visited.
func (s Scope) Enumerator() *ScopeEnumerator {
	// Collect newest-first, dropping shadowed names, then reverse. Scopes
	// are small, so a linear duplicate check beats a map.
	var names []string
	var vals []Expr
	for f := s.f; f != nil; f = f.parent {
	entries:
		for i := len(f.names) - 1; i >= 0; i-- {
			n := f.names[i]
			for _, seen := range names {
				if seen == n {
					continue entries
				}
			}
			names = append(names, n)
			vals = append(vals, f.vals[i])
		}
	}
	for i, j := 0, len(names)-1; i < j; i, j = i+1, j-1 {
		names[i], names[j] = names[j], names[i]
		vals[i], vals[j] = vals[j], vals[i]
	}
	return &ScopeEnumerator{names: names, vals: vals, i: -1}
}

// OrderedNames returns the names of this tuple in sorted order.
func (s Scope) OrderedNames() []string {
	names := s.Names()
	sort.Strings(names)
	return names
}

// ScopeEnumerator represents an enumerator over a Scope.
type ScopeEnumerator struct {
	names []string
	vals  []Expr
	i     int
}

// MoveNext moves the enumerator to the next Value.
func (e *ScopeEnumerator) MoveNext() bool {
	e.i++
	return e.i < len(e.names)
}

// Current returns the enumerator's current Value.
func (e *ScopeEnumerator) Current() (string, Expr) {
	return e.names[e.i], e.vals[e.i]
}

// scopeBuilder accumulates pattern bindings into a single frame. Patterns
// bind their parts one at a time; collecting them here (rather than merging
// a scope per part) gives the enclosing Update one flat frame to push and
// keeps every binding of a pattern in a fixed slot.
type scopeBuilder struct {
	names []string
	vals  []Expr
}

// matchedAdd adds t's bindings, with MatchedUpdate's rule: a name bound
// twice must be bound to the same value.
func (b *scopeBuilder) matchedAdd(t Scope) error {
	names, vals := t.flatten()
	for i, name := range names {
		if j := b.index(name); j >= 0 {
			if b.vals[j].String() != vals[i].String() {
				return fmt.Errorf("the value of %s is different in both scopes", name)
			}
			continue
		}
		b.names = append(b.names, name)
		b.vals = append(b.vals, vals[i])
	}
	return nil
}

func (b *scopeBuilder) index(name string) int {
	for i, n := range b.names {
		if n == name {
			return i
		}
	}
	return -1
}

func (b *scopeBuilder) finish() Scope {
	if len(b.names) == 0 {
		return EmptyScope
	}
	return Scope{&frame{names: b.names, vals: b.vals}}
}

package rel

import (
	"context"
	"fmt"
	"reflect"

	"github.com/arr-ai/arrai/pkg/fu"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

// Closure represents the closure of a function over a scope.
type Closure struct {
	scope Scope
	f     *Function
}

// NewFunction returns a new function.
func NewClosure(scope Scope, f *Function) Closure {
	return Closure{scope: scope, f: f}
}

// Closure equality (#777). Function equality is undecidable in general, so a
// closure is compared by its compiled function node (identity, as
// EqualFunction) together with the bindings it captured. "Captured" means the
// function's free identifiers — the names its body refers to without binding
// them — not the whole scope chain. With `let mk = \k \a a + k`, `mk(1)` and
// `mk(1)` are equal although each call built its own frame, `mk(1)` and
// `mk(2)` are not; with `let mk = \k \a a`, `mk(1) = mk(2)` holds because `k`
// is never used.
// Bindings compare by value (see sameBinding); a free identifier bound in
// neither scope counts as equal. When the body holds a node the tree walker
// cannot see into (children returns ok=false) the free identifiers are
// unknown, and every visible binding is compared instead. Hash follows the
// same rule exactly: the function node and the captured bindings, in an
// order-independent fold, so equal closures always hash alike.

// capturedNames returns the names whose bindings the closure's behaviour can
// depend on: the function's free identifiers when known, otherwise every
// name visible in the captured scope.
func (c Closure) capturedNames() (names []string, exact bool) {
	if names, known := c.f.freeIdents(); known {
		return names, true
	}
	return c.scope.Names(), false
}

// Hash computes a hash for a Closure: the function node folded with each
// captured binding. The fold is order-independent (xor) because the fallback
// name order depends on the shape of the frame chain, not just its contents.
func (c Closure) Hash() uintptr {
	h := mix(closureSalt, c.f.Hash())
	names, _ := c.capturedNames()
	var captured uintptr
	for _, name := range names {
		if v, bound := c.scope.Get(name); bound {
			captured = xor(captured, mix(hashString(name), hashBinding(v)))
		}
	}
	return mix(h, captured)
}

// Equal tests two Values for equality. Any other type returns false.
func (c Closure) Equal(i Value) bool {
	if d, ok := i.(Closure); ok {
		var pairs recPairs
		return c.equalClosure(d, &pairs)
	}
	return false
}

// equalClosure is Equal with the recursive-binding pairs already assumed
// equal by an enclosing comparison; see sameBindingRec.
func (c Closure) equalClosure(d Closure, pairs *recPairs) bool {
	if !c.f.EqualFunction(d.f) {
		return false
	}
	if c.scope.f == d.scope.f {
		// Same frame: identical captures.
		return true
	}
	names, exact := c.capturedNames()
	if !exact && len(d.scope.Names()) != len(names) {
		return false
	}
	for _, name := range names {
		a, aBound := c.scope.Get(name)
		b, bBound := d.scope.Get(name)
		if aBound != bBound {
			return false
		}
		if aBound && !sameBindingRec(a, b, pairs) {
			return false
		}
	}
	return true
}

// recPairs records pairs of recursive-binding cells whose values are being
// compared further up the stack. A self-referential closure captures the
// cell that holds it, so comparing two such closures meets the same pair of
// cells again; taking the pair as equal at that point (a bisimulation) is
// what lets the comparison terminate, and it is sound: the cells are equal
// iff everything else about the two closures agrees.
type recPairs map[[2]*recCell]bool

// sameBinding reports whether two scope entries bind the same thing. It
// serves a pattern variable bound twice (`let [x, x] = ...`) and closure
// captures. Values compare by value; a recursive binding's cell by the value
// it holds; anything else (nothing today) by its printed form, so the
// fallback is the rule the scope matchers always had.
func sameBinding(a, b Expr) bool {
	var pairs recPairs
	return sameBindingRec(a, b, &pairs)
}

func sameBindingRec(a, b Expr, pairs *recPairs) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	ca, aIsCell := a.(*recCell)
	cb, bIsCell := b.(*recCell)
	if aIsCell || bIsCell {
		if !aIsCell || !bIsCell {
			return false
		}
		return sameRecCell(ca, cb, pairs)
	}
	va, aIsValue := a.(Value)
	vb, bIsValue := b.(Value)
	if aIsValue || bIsValue {
		if !aIsValue || !bIsValue {
			return false
		}
		return sameValueRec(va, vb, pairs)
	}
	return a.String() == b.String()
}

// sameRecCell compares two recursive bindings by the values they hold. An
// unset cell (the binding is still being evaluated) equals only itself.
func sameRecCell(a, b *recCell, pairs *recPairs) bool {
	if a == b {
		return true
	}
	if !a.set || !b.set || a.name != b.name {
		return false
	}
	pair := [2]*recCell{a, b}
	if (*pairs)[pair] {
		return true
	}
	if *pairs == nil {
		*pairs = recPairs{}
	}
	(*pairs)[pair] = true
	return sameValueRec(a.val, b.val, pairs)
}

// sameValueRec is Value.Equal threaded through the closures a recursive
// binding can hold (a closure or a tuple of closures), so that the assumed
// pairs reach them. Any other value compares by Equal.
func sameValueRec(a, b Value, pairs *recPairs) bool {
	switch a := a.(type) {
	case Closure:
		b, ok := b.(Closure)
		return ok && a.equalClosure(b, pairs)
	case Tuple:
		b, ok := b.(Tuple)
		if !ok || a.Count() != b.Count() {
			return false
		}
		for e := a.Enumerator(); e.MoveNext(); {
			name, av := e.Current()
			bv, found := b.Get(name)
			if !found || !sameValueRec(av, bv, pairs) {
				return false
			}
		}
		return true
	}
	return a.Equal(b)
}

// hashBinding hashes a scope entry consistently with sameBinding. A
// recursive cell hashes by name alone: its value holds the cell, so hashing
// through it would not terminate, and equal cells always share a name.
// Anything that is neither a Value nor a cell hashes to a constant, which is
// trivially consistent.
func hashBinding(e Expr) uintptr {
	switch e := e.(type) {
	case *recCell:
		return mix(recCellSalt, hashString(e.name))
	case Value:
		return e.Hash()
	}
	return 0
}

// String returns a string representation of the expression.
func (c Closure) String() string {
	return c.f.String()
}

// Format formats the expression.
func (c Closure) Format(f fmt.State, verb rune) {
	fu.FRepr(f, c.f)
}

// Eval returns the Value
func (c Closure) Eval(ctx context.Context, local Scope) (Value, error) {
	return c, nil
}

// Source returns a scanner locating the Closure's source code.
func (c Closure) Source() parser.Scanner {
	return *parser.NewScanner("")
}

var closureKind = registerKind(205, reflect.TypeOf(Closure{}))

// Kind returns a number that is unique for each major kind of Value.
func (c Closure) Kind() int {
	return closureKind
}

// Bool returns true iff the tuple has attributes.
func (c Closure) IsTrue() bool {
	return true
}

// Less returns true iff g is not a number or f.number < g.number.
func (c Closure) Less(d Value) bool {
	if c.Kind() != d.Kind() {
		return c.Kind() < d.Kind()
	}
	return c.String() < d.String()
}

// Negate returns {(negateTag): f}.
func (c Closure) Negate() Value {
	return NewTuple(NewAttr(negateTag, c))
}

// Export exports a Closure.
func (c Closure) Export(ctx context.Context) interface{} {
	if c.f.Arg() == "-" {
		result, err := SetCall(ctx, c, None)
		if err != nil {
			panic(err)
		}
		return result.Export(ctx)
	}
	return func(arg Value) Value {
		result, err := SetCall(ctx, c, None)
		if err != nil {
			panic(err)
		}
		return result
	}
}

func (Closure) getSetBuilder() setBuilder {
	return newGenericTypeSetBuilder()
}

func (Closure) getBucket() fmt.Stringer {
	return genericType
}

// A Closure is the one-element set {c}; see value_set_singleton.go.

func (Closure) Count() int {
	return 1
}

func (c Closure) Has(v Value) bool {
	return c.Equal(v)
}

func (c Closure) Enumerator() ValueEnumerator {
	return &singletonEnumerator{v: c}
}

func (c Closure) With(v Value) Set {
	return singletonWith(c, v)
}

func (c Closure) Without(v Value) Set {
	return singletonWithout(c, v)
}

func (c Closure) Map(f func(v Value) (Value, error)) (Set, error) {
	return singletonMap(c, f)
}

func (c Closure) Where(p func(v Value) (bool, error)) (Set, error) {
	return singletonWhere(c, p)
}

// FIXME: context not used properly
func (c Closure) CallAll(ctx context.Context, arg Value, b SetBuilder) error {
	val, err := c.call(ctx, arg)
	if err != nil {
		return err
	}
	b.Add(val)
	return nil
}

// call applies the closure to arg and returns its single result.
func (c Closure) call(ctx context.Context, arg Value) (Value, error) {
	niladic := c.f.Arg() == "-"
	noArg := arg == nil
	if niladic != noArg {
		panic(errors.Errorf(
			"nullary-vs-unary function arg mismatch (%s vs %s)", c.f.Arg(), arg))
	}
	if niladic {
		return c.f.body.Eval(ctx, c.scope)
	}
	// Fast path for the common `\x ...` shape: IdentPattern.Bind never fails
	// and only produces a one-entry scope, so bind directly instead of
	// building and merging a throwaway Scope. With ignores "_" exactly as
	// Bind+Update would.
	if ident, is := c.f.arg.(IdentPattern); is {
		return c.f.body.Eval(ctx, c.scope.With(string(ident), arg))
	}
	var b scopeBuilder
	ctx, err := c.f.arg.Bind(ctx, c.scope, arg, &b)
	if err != nil {
		if err == errPatternMismatch {
			err = explainBind(ctx, c.f.arg, c.scope, arg)
		}
		return nil, err
	}
	return c.f.body.Eval(ctx, c.scope.updateWith(&b))
}

func (Closure) unionSetSubsetBucket() string {
	// TODO: create its own subset bucket in union set
	return genericType.String()
}

func (c Closure) ArrayEnumerator() ValueEnumerator {
	return &singletonEnumerator{v: c}
}

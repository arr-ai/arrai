package rel

// A function value (*NativeFunction, Closure, ExprClosure) is arr.ai's
// one-element set {f}: it already reports Count() == 1 and IsTrue() == true,
// and hashes as an element of larger sets. The helpers below complete that
// picture so the structural Set operations on a bare function behave exactly
// as they would on the literal set {f}, instead of crashing with a Go stack
// trace. Membership is decided by the function's own Equal, which every
// function kind defines as node identity.

// singletonEnumerator yields one value and then stops.
type singletonEnumerator struct {
	v    Value
	done bool
}

func (e *singletonEnumerator) MoveNext() bool {
	if e.done {
		return false
	}
	e.done = true
	return true
}

func (e *singletonEnumerator) Current() Value {
	return e.v
}

// singletonWith returns {f, v}, or f itself when v is already f.
func singletonWith(f Set, v Value) Set {
	if f.Equal(v) {
		return f
	}
	return MustNewSet(f, v)
}

// singletonWithout returns {} when v is f, otherwise f unchanged.
func singletonWithout(f Set, v Value) Set {
	if f.Equal(v) {
		return None
	}
	return f
}

// singletonMap returns the one-element set of fn(f).
func singletonMap(f Set, fn func(Value) (Value, error)) (Set, error) {
	v, err := fn(f)
	if err != nil {
		return nil, err
	}
	return NewSet(v)
}

// singletonWhere returns f when p(f) holds, otherwise {}.
func singletonWhere(f Set, p func(Value) (bool, error)) (Set, error) {
	keep, err := p(f)
	if err != nil {
		return nil, err
	}
	if keep {
		return f, nil
	}
	return None, nil
}

package rel

// Force returns v with every suspended >> pipeline inside it materialised,
// so the result, and everything reachable from it through arrays, dicts,
// tuples and sets, has a concrete representation: an Array where arr.ai sees
// an array, a Dict where it sees a dict (arr-ai/arrai#789). Errors that a
// pipeline deferred are returned here instead of surfacing later as a panic
// from Export or Hash.
//
// Force is for Go code that receives values from arr.ai and inspects their
// concrete types. syntax.EvaluateExpr, EvalWithScope and EvaluateBundleCtx
// already force their results. Values captured inside closures are not
// visited, so a value obtained by calling a closure from Go may itself need
// Force.
//
// Force is O(size of v). Containers with nothing to force are returned
// unchanged, so the result shares structure with v.
func Force(v Value) (Value, error) {
	f, _, err := force(v)
	return f, err
}

// force reports whether it changed v, so containers rebuild only when a
// descendant was materialised.
func force(v Value) (Value, bool, error) {
	switch x := v.(type) {
	case seqPipeline:
		a, err := x.force()
		if err != nil {
			return nil, false, err
		}
		f, _, err := force(a)
		return f, true, err
	case dictPipeline:
		d, err := x.force()
		if err != nil {
			return nil, false, err
		}
		f, _, err := force(d)
		return f, true, err
	case nil, Number, String, Bytes, EmptySet, TrueSet, StringCharTuple, BytesByteTuple,
		Closure, ExprClosure, *NativeFunction:
		return v, false, nil
	case Array:
		return forceArray(x)
	case Dict:
		return forceDict(x)
	case *GenericTuple:
		return forceTuple(x)
	case ArrayItemTuple:
		item, changed, err := force(x.item)
		if err != nil || !changed {
			return v, false, err
		}
		return NewArrayItemTuple(x.at, item), true, nil
	case DictEntryTuple:
		e, changed, err := forceEntry(x)
		return e, changed, err
	case Set:
		return forceSet(x)
	}
	return v, false, nil
}

func forceEntry(e DictEntryTuple) (DictEntryTuple, bool, error) {
	at, c1, err := force(e.at)
	if err != nil {
		return e, false, err
	}
	value, c2, err := force(e.value)
	if err != nil {
		return e, false, err
	}
	if !c1 && !c2 {
		return e, false, nil
	}
	return NewDictEntryTuple(at, value), true, nil
}

func forceArray(a Array) (Value, bool, error) {
	var values []Value
	for i, item := range a.values {
		if item == nil {
			continue
		}
		f, changed, err := force(item)
		if err != nil {
			return nil, false, err
		}
		if changed && values == nil {
			values = make([]Value, len(a.values))
			copy(values, a.values[:i])
		}
		if values != nil {
			values[i] = f
		}
	}
	if values == nil {
		return a, false, nil
	}
	b := a.derive()
	b.values = values
	b.abuf = nil
	return b, true, nil
}

func forceDict(d Dict) (Value, bool, error) {
	entries := make([]DictEntryTuple, 0, d.Count())
	dirty := false
	for e := d.Enumerator(); e.MoveNext(); {
		f, changed, err := forceEntry(e.Current().(DictEntryTuple))
		if err != nil {
			return nil, false, err
		}
		dirty = dirty || changed
		entries = append(entries, f)
	}
	if !dirty {
		return d, false, nil
	}
	s, err := NewDict(true, entries...)
	return s, err == nil, err
}

func forceTuple(t *GenericTuple) (Value, bool, error) {
	attrs := make([]Attr, 0, t.Count())
	dirty := false
	for e := t.Enumerator(); e.MoveNext(); {
		name, value := e.Current()
		f, changed, err := force(value)
		if err != nil {
			return nil, false, err
		}
		dirty = dirty || changed
		attrs = append(attrs, Attr{name, f})
	}
	if !dirty {
		return t, false, nil
	}
	return NewTuple(attrs...), true, nil
}

func forceSet(s Set) (Value, bool, error) {
	values := make([]Value, 0, s.Count())
	dirty := false
	for e := s.Enumerator(); e.MoveNext(); {
		f, changed, err := force(e.Current())
		if err != nil {
			return nil, false, err
		}
		dirty = dirty || changed
		values = append(values, f)
	}
	if !dirty {
		return s, false, nil
	}
	r, err := NewSet(values...)
	return r, err == nil, err
}

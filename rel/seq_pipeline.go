package rel

import (
	"context"
	"fmt"
	"sync"

	"github.com/arr-ai/hash/hash128"
	"github.com/arr-ai/wbnf/parser"
)

// seqPipeline is a suspended >> chain over an Array (🎯T19). Count is the
// base length (maps are 1-1). Enumeration, hashing and equality force the
// chain in one pass, so n stages materialise once.
type seqPipeline struct {
	base Array
	ops  []seqOp
	cell *seqForceCell
}

type seqOp struct {
	call func(at, v Value) (Value, error)
}

type seqForceCell struct {
	once sync.Once
	arr  Array
	err  error
}

// Observe forces suspended >> pipelines so errors and concrete Array/Dict
// representations surface. Count and further >> do not call this.
func Observe(v Value) (Value, error) {
	switch v := v.(type) {
	case seqPipeline:
		return v.force()
	case dictPipeline:
		return v.force()
	default:
		return v, nil
	}
}

func newSeqPipeline(base Array, call func(at, v Value) (Value, error)) seqPipeline {
	return seqPipeline{base: base, ops: []seqOp{{call: call}}, cell: &seqForceCell{}}
}

func (p seqPipeline) then(call func(at, v Value) (Value, error)) seqPipeline {
	ops := make([]seqOp, len(p.ops)+1)
	copy(ops, p.ops)
	ops[len(p.ops)] = seqOp{call: call}
	return seqPipeline{base: p.base, ops: ops, cell: &seqForceCell{}}
}

func (p seqPipeline) apply(at int, item Value) (Value, error) {
	v := item
	atv := NewNumber(float64(p.base.offset + at))
	for _, op := range p.ops {
		nv, err := op.call(atv, v)
		if err != nil {
			return nil, err
		}
		v = nv
	}
	return v, nil
}

func (p seqPipeline) force() (Array, error) {
	p.cell.once.Do(func() {
		items := make([]Value, len(p.base.values))
		if fastPaths {
			if ranges := parallelRanges(len(p.base.values)); ranges != nil {
				errs := make([]error, len(ranges))
				runRanges(ranges, func(w, lo, hi int) {
					for at := lo; at < hi; at++ {
						if item := p.base.values[at]; item != nil {
							v, err := p.apply(at, item)
							if err != nil {
								errs[w] = err
								return
							}
							items[at] = v
						}
					}
				})
				if err := firstErr(errs); err != nil {
					p.cell.err = err
					return
				}
				s := NewOffsetArray(p.base.offset, items...)
				if a, ok := s.(Array); ok {
					p.cell.arr = a
				}
				return
			}
		}
		for at, item := range p.base.values {
			if item == nil {
				continue
			}
			v, err := p.apply(at, item)
			if err != nil {
				p.cell.err = err
				return
			}
			items[at] = v
		}
		s := NewOffsetArray(p.base.offset, items...)
		if a, ok := s.(Array); ok {
			p.cell.arr = a
		}
	})
	return p.cell.arr, p.cell.err
}

func (p seqPipeline) mustArray() Array {
	a, err := p.force()
	if err != nil {
		panic(err)
	}
	return a
}

func (p seqPipeline) Count() int { return p.base.Count() }

func (p seqPipeline) Hash(seed uintptr) uintptr { return p.mustArray().Hash(seed) }
func (p seqPipeline) Hash128() hash128.H128     { return p.mustArray().Hash128() }
func (p seqPipeline) Equal(v Value) bool {
	if q, ok := v.(seqPipeline); ok {
		a, err1 := p.force()
		b, err2 := q.force()
		if err1 != nil || err2 != nil {
			return false
		}
		return a.Equal(b)
	}
	a, err := p.force()
	return err == nil && a.Equal(v)
}
func (p seqPipeline) String() string {
	a, err := p.force()
	if err != nil {
		return fmt.Sprintf(">>(%v)", err)
	}
	return a.String()
}
func (p seqPipeline) Format(f fmt.State, verb rune) {
	a, err := p.force()
	if err != nil {
		fmt.Fprintf(f, ">>(%v)", err)
		return
	}
	a.Format(f, verb)
}
func (p seqPipeline) Eval(ctx context.Context, local Scope) (Value, error) {
	return p.force()
}
func (p seqPipeline) Source() parser.Scanner { return p.base.Source() }
func (p seqPipeline) Kind() int              { return p.base.Kind() }
func (p seqPipeline) IsTrue() bool           { return p.base.IsTrue() }
func (p seqPipeline) Less(v Value) bool {
	a, err := p.force()
	if err != nil {
		panic(err)
	}
	ov, err := Observe(v)
	if err != nil {
		panic(err)
	}
	return a.Less(ov)
}
func (p seqPipeline) Negate() Value { return p.mustArray().Negate() }
func (p seqPipeline) Export(ctx context.Context) interface{} {
	return p.mustArray().Export(ctx)
}
func (p seqPipeline) getSetBuilder() setBuilder { return p.base.getSetBuilder() }
func (p seqPipeline) getBucket() fmt.Stringer   { return p.base.getBucket() }
func (p seqPipeline) Has(value Value) bool      { return p.mustArray().Has(value) }
func (p seqPipeline) With(value Value) Set      { return p.mustArray().With(value) }
func (p seqPipeline) Without(value Value) Set   { return p.mustArray().Without(value) }
func (p seqPipeline) Map(f func(Value) (Value, error)) (Set, error) {
	return p.mustArray().Map(f)
}
func (p seqPipeline) Where(pred func(Value) (bool, error)) (Set, error) {
	return p.mustArray().Where(pred)
}
func (p seqPipeline) CallAll(ctx context.Context, arg Value, b SetBuilder) error {
	return p.mustArray().CallAll(ctx, arg, b)
}
func (p seqPipeline) Enumerator() ValueEnumerator      { return p.mustArray().Enumerator() }
func (p seqPipeline) ArrayEnumerator() ValueEnumerator { return p.mustArray().ArrayEnumerator() }
func (p seqPipeline) unionSetSubsetBucket() string     { return p.base.unionSetSubsetBucket() }

// dictPipeline is a suspended >> chain over a Dict (🎯T19). Count is the
// base size (maps are 1-1 on entries). Enumeration, hashing and equality
// force the chain in one pass.
type dictPipeline struct {
	base Dict
	ops  []seqOp
	cell *dictForceCell
}

type dictForceCell struct {
	once sync.Once
	dict Dict
	err  error
}

func newDictPipeline(base Dict, call func(at, v Value) (Value, error)) dictPipeline {
	return dictPipeline{base: base, ops: []seqOp{{call: call}}, cell: &dictForceCell{}}
}

func (p dictPipeline) then(call func(at, v Value) (Value, error)) dictPipeline {
	ops := make([]seqOp, len(p.ops)+1)
	copy(ops, p.ops)
	ops[len(p.ops)] = seqOp{call: call}
	return dictPipeline{base: p.base, ops: ops, cell: &dictForceCell{}}
}

func (p dictPipeline) apply(at, v Value) (Value, error) {
	for _, op := range p.ops {
		nv, err := op.call(at, v)
		if err != nil {
			return nil, err
		}
		v = nv
	}
	return v, nil
}

func (p dictPipeline) force() (Dict, error) {
	p.cell.once.Do(func() {
		entries := make([]DictEntryTuple, 0, p.base.Count())
		for i := p.base.Enumerator(); i.MoveNext(); {
			entries = append(entries, i.Current().(DictEntryTuple))
		}
		mapEntries := func(lo, hi int) error {
			for j := lo; j < hi; j++ {
				v, err := p.apply(entries[j].at, entries[j].value)
				if err != nil {
					return err
				}
				entries[j] = NewDictEntryTuple(entries[j].at, v)
			}
			return nil
		}
		if ranges := parallelRanges(len(entries)); fastPaths && ranges != nil {
			errs := make([]error, len(ranges))
			runRanges(ranges, func(w, lo, hi int) {
				errs[w] = mapEntries(lo, hi)
			})
			if err := firstErr(errs); err != nil {
				p.cell.err = err
				return
			}
		} else if err := mapEntries(0, len(entries)); err != nil {
			p.cell.err = err
			return
		}
		d, err := NewDict(true, entries...)
		if err != nil {
			p.cell.err = err
			return
		}
		if dd, ok := d.(Dict); ok {
			p.cell.dict = dd
		}
	})
	return p.cell.dict, p.cell.err
}

func (p dictPipeline) mustDict() Dict {
	d, err := p.force()
	if err != nil {
		panic(err)
	}
	return d
}

func (p dictPipeline) Count() int { return p.base.Count() }

func (p dictPipeline) Hash(seed uintptr) uintptr { return p.mustDict().Hash(seed) }
func (p dictPipeline) Hash128() hash128.H128     { return p.mustDict().Hash128() }
func (p dictPipeline) Equal(v Value) bool {
	if q, ok := v.(dictPipeline); ok {
		a, err1 := p.force()
		b, err2 := q.force()
		if err1 != nil || err2 != nil {
			return false
		}
		return a.Equal(b)
	}
	a, err := p.force()
	return err == nil && a.Equal(v)
}
func (p dictPipeline) String() string {
	d, err := p.force()
	if err != nil {
		return fmt.Sprintf(">>(%v)", err)
	}
	return d.String()
}
func (p dictPipeline) Format(f fmt.State, verb rune) {
	d, err := p.force()
	if err != nil {
		fmt.Fprintf(f, ">>(%v)", err)
		return
	}
	d.Format(f, verb)
}
func (p dictPipeline) Eval(_ context.Context, _ Scope) (Value, error) {
	return p.force()
}
func (p dictPipeline) Source() parser.Scanner { return p.base.Source() }
func (p dictPipeline) Kind() int              { return p.base.Kind() }
func (p dictPipeline) IsTrue() bool           { return p.base.IsTrue() }
func (p dictPipeline) Less(v Value) bool {
	d, err := p.force()
	if err != nil {
		panic(err)
	}
	ov, err := Observe(v)
	if err != nil {
		panic(err)
	}
	return d.Less(ov)
}
func (p dictPipeline) Negate() Value { return p.mustDict().Negate() }
func (p dictPipeline) Export(ctx context.Context) interface{} {
	return p.mustDict().Export(ctx)
}
func (p dictPipeline) getSetBuilder() setBuilder { return p.base.getSetBuilder() }
func (p dictPipeline) getBucket() fmt.Stringer   { return p.base.getBucket() }
func (p dictPipeline) Has(value Value) bool      { return p.mustDict().Has(value) }
func (p dictPipeline) With(value Value) Set      { return p.mustDict().With(value) }
func (p dictPipeline) Without(value Value) Set   { return p.mustDict().Without(value) }
func (p dictPipeline) Map(f func(Value) (Value, error)) (Set, error) {
	return p.mustDict().Map(f)
}
func (p dictPipeline) Where(pred func(Value) (bool, error)) (Set, error) {
	return p.mustDict().Where(pred)
}
func (p dictPipeline) CallAll(ctx context.Context, arg Value, b SetBuilder) error {
	return p.mustDict().CallAll(ctx, arg, b)
}
func (p dictPipeline) Enumerator() ValueEnumerator      { return p.mustDict().Enumerator() }
func (p dictPipeline) ArrayEnumerator() ValueEnumerator { return p.mustDict().ArrayEnumerator() }
func (p dictPipeline) unionSetSubsetBucket() string     { return p.base.unionSetSubsetBucket() }

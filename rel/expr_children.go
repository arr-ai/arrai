package rel

// Tree walking for rewrites (🎯T28). children enumerates an Expr's direct
// children with the facts a rewrite needs about each: how often the child
// is evaluated relative to its parent, and which identifiers the parent
// binds for it. It also returns a rebuild function that produces the same
// kind of node with replacement children, going through the constructor
// wherever one analyses its operands (where, orderby, =>, tuples, …), so a
// rewritten tree carries the same derived state as a freshly compiled one.

// ChildKind is how often a child evaluates per evaluation of its parent.
type ChildKind uint8

const (
	// ChildOnce: exactly once, unconditionally.
	ChildOnce ChildKind = iota
	// ChildMaybe: at most once (a conditional branch, a short-circuit rhs).
	ChildMaybe
	// ChildMany: any number of times (a function body, a mapped body).
	ChildMany
)

func (k ChildKind) join(o ChildKind) ChildKind {
	if o > k {
		return o
	}
	return k
}

// Child is one direct child of an Expr.
type Child struct {
	Expr Expr
	Kind ChildKind
	// Binds, when non-nil, is the pattern whose identifiers are in scope
	// inside Expr and nowhere else in the parent.
	Binds Pattern
}

// Rewritable lets an Expr defined outside this package take part in tree
// rewrites. WithChildren receives replacement children in Children() order
// and returns the rebuilt node.
type Rewritable interface {
	Children() []Child
	WithChildren(kids []Expr) Expr
}

// Effectful marks an Expr whose evaluation has observable side effects
// (logging, file or network I/O), so a rewrite must not move it.
type Effectful interface {
	HasSideEffects() bool
}

func once(es ...Expr) []Child {
	kids := make([]Child, len(es))
	for i, e := range es {
		kids[i] = Child{Expr: e, Kind: ChildOnce}
	}
	return kids
}

// fnChildren lists a Function's children: expressions embedded in its
// pattern (fallbacks, literal sub-patterns), evaluated before the pattern
// binds, then the body under the pattern's bindings.
func fnChildren(f *Function, body ChildKind) []Child {
	kids := make([]Child, 0, 2)
	for _, e := range patternExprs(f.arg) {
		kids = append(kids, Child{Expr: e, Kind: ChildMany})
	}
	return append(kids, Child{Expr: f.body, Kind: body, Binds: f.arg})
}

// rebuildFn takes a Function's children back in fnChildren order.
func rebuildFn(f *Function, kids []Expr) *Function {
	return NewFunction(f.Src, f.arg, kids[len(kids)-1]).(*Function)
}

// children returns e's direct children and a rebuild function, or ok=false
// for a node the walker does not know, which a rewrite must treat as opaque.
//
//nolint:funlen,gocyclo,cyclop,gocognit
func children(e Expr) (kids []Child, rebuild func([]Expr) Expr, ok bool) {
	switch e := e.(type) {
	case nil:
		return nil, nil, false
	case Rewritable:
		return e.Children(), e.WithChildren, true
	case IdentExpr, DynIdentExpr, LiteralExpr, Closure, ExprClosure:
		return nil, func([]Expr) Expr { return e }, true
	case *BinExpr:
		if fn, isFn := e.b.(*Function); isFn {
			kids = append(once(e.a), fnChildren(fn, ChildMany)...)
			return kids, func(k []Expr) Expr {
				return rebuildBin(e, k[0], rebuildFn(fn, k[1:]))
			}, true
		}
		return once(e.a, e.b), func(k []Expr) Expr { return rebuildBin(e, k[0], k[1]) }, true
	case *UnaryExpr:
		return once(e.a), func(k []Expr) Expr {
			n := *e
			n.a = k[0]
			return &n
		}, true
	case *DotExpr:
		return once(e.lhs), func(k []Expr) Expr { return NewDotExpr(e.Src, k[0], e.attr) }, true
	case *Function:
		return fnChildren(e, ChildMany), func(k []Expr) Expr { return rebuildFn(e, k) }, true
	case *ArrowExpr:
		// `lhs -> \pat body` (and the desugared let) applies fn exactly once.
		kids = append(once(e.lhs), fnChildren(e.fn, ChildOnce)...)
		return kids, func(k []Expr) Expr { return NewArrowExpr(e.Src, k[0], rebuildFn(e.fn, k[1:])) }, true
	case *DArrowExpr:
		kids = append(once(e.lhs), fnChildren(e.fn, ChildMany)...)
		return kids, func(k []Expr) Expr { return NewDArrowExpr(e.Src, k[0], rebuildFn(e.fn, k[1:])) }, true
	case *SeqArrowExpr:
		kids = append(once(e.lhs), fnChildren(e.fn, ChildMany)...)
		return kids, func(k []Expr) Expr {
			return NewSeqArrowExpr(e.withAt)(e.Src, k[0], rebuildFn(e.fn, k[1:]))
		}, true
	case *TupleMapExpr:
		kids = append(once(e.lhs), fnChildren(e.fn, ChildMany)...)
		return kids, func(k []Expr) Expr { return NewTupleMapExpr(e.Src, k[0], rebuildFn(e.fn, k[1:])) }, true
	case *ReduceExpr:
		kids = append(once(e.a), fnChildren(e.f, ChildMany)...)
		return kids, func(k []Expr) Expr {
			n := *e
			n.a, n.f = k[0], rebuildFn(e.f, k[1:])
			return &n
		}, true
	case *TupleExpr:
		es := make([]Expr, len(e.attrs))
		for i, a := range e.attrs {
			es[i] = a.expr
		}
		return once(es...), func(k []Expr) Expr {
			attrs := make([]AttrExpr, len(e.attrs))
			for i, a := range e.attrs {
				attrs[i] = AttrExpr{ExprScanner: a.ExprScanner, name: a.name, expr: k[i]}
			}
			return NewTupleExpr(e.Src, attrs...)
		}, true
	case ArrayExpr:
		// Sparse arrays hold nil elements; keep their positions.
		var es []Expr
		for _, el := range e.elements {
			if el != nil {
				es = append(es, el)
			}
		}
		return once(es...), func(k []Expr) Expr {
			els := make([]Expr, len(e.elements))
			j := 0
			for i, el := range e.elements {
				if el != nil {
					els[i] = k[j]
					j++
				}
			}
			return NewArrayExpr(e.Src, els...)
		}, true
	case *SetExpr:
		return once(e.elements...), func(k []Expr) Expr {
			s, err := NewSetExpr(e.Src, k...)
			if err != nil {
				return e
			}
			return s
		}, true
	case BytesExpr:
		return once(e.elements...), func(k []Expr) Expr { return NewBytesExpr(e.Src, k...) }, true
	case DictExpr:
		return dictChildren(e, false)
	case *DictExpr:
		return dictChildren(*e, true)
	case CompareExpr:
		// A chained comparison stops at the first false pair.
		kids = make([]Child, len(e.args))
		for i, a := range e.args {
			kind := ChildOnce
			if i >= 2 {
				kind = ChildMaybe
			}
			kids[i] = Child{Expr: a, Kind: kind}
		}
		return kids, func(k []Expr) Expr { return NewCompareExpr(e.Src, k, e.comps, e.ops) }, true
	case AndExpr:
		return []Child{{Expr: e.a, Kind: ChildOnce}, {Expr: e.b, Kind: ChildMaybe}},
			func(k []Expr) Expr { return NewAndExpr(e.Src, k[0], k[1]) }, true
	case OrExpr:
		return []Child{{Expr: e.a, Kind: ChildOnce}, {Expr: e.b, Kind: ChildMaybe}},
			func(k []Expr) Expr { return NewOrExpr(e.Src, k[0], k[1]) }, true
	case *IfElseExpr:
		kids = []Child{
			{Expr: e.cond, Kind: ChildOnce}, {Expr: e.ifTrue, Kind: ChildMaybe}, {Expr: e.ifFalse, Kind: ChildMaybe},
		}
		return kids, func(k []Expr) Expr { return NewIfElseExpr(e.Src, k[1], k[0], k[2]) }, true
	case CondExpr:
		// The arms dictionary evaluates its keys in order until one holds.
		return []Child{{Expr: e.dicExpr, Kind: ChildMaybe}},
			func(k []Expr) Expr { return NewCondExpr(e.Src, k[0]) }, true
	case CondPatternControlVarExpr:
		kids = once(e.controlVarExpr)
		for _, p := range e.conditionPairs {
			for _, pe := range patternExprs(p.pattern) {
				kids = append(kids, Child{Expr: pe, Kind: ChildMaybe})
			}
			kids = append(kids, Child{Expr: p.expr, Kind: ChildMaybe, Binds: p.pattern})
		}
		return kids, func(k []Expr) Expr {
			pairs := make([]PatternExprPair, len(e.conditionPairs))
			j := 1
			for i, p := range e.conditionPairs {
				j += len(patternExprs(p.pattern))
				pairs[i] = PatternExprPair{pattern: p.pattern, expr: k[j]}
				j++
			}
			return NewCondPatternControlVarExpr(e.Src, k[0], pairs...)
		}, true
	case *NestExpr:
		return once(e.lhs), func(k []Expr) Expr { return NewNestExpr(e.Src, e.inverse, k[0], e.attrs, e.attr) }, true
	case *UnnestExpr:
		return once(e.lhs), func(k []Expr) Expr { return NewUnnestExpr(e.Src, k[0], e.attr) }, true
	case *SingleNestExpr:
		return once(e.lhs), func(k []Expr) Expr { return NewSingleNestExpr(e.Src, k[0], e.attr) }, true
	case *OffsetExpr:
		return once(e.offset, e.array), func(k []Expr) Expr { return NewOffsetExpr(e.Src, k[0], k[1]) }, true
	case *TupleProjectExpr:
		return once(e.base), func(k []Expr) Expr {
			return NewTupleProjectExpr(e.Src, k[0], e.inverse, e.attrs.OrderedNames())
		}, true
	case RecursionExpr:
		return []Child{{Expr: e.fn, Kind: ChildMany, Binds: IdentPattern(e.name)}},
			func(k []Expr) Expr { return NewRecursionExpr(e.Src, e.name, k[0]) }, true
	case *RecursionExpr:
		return []Child{{Expr: e.fn, Kind: ChildMany, Binds: IdentPattern(e.name)}},
			func(k []Expr) Expr { return NewRecursionExpr(e.Src, e.name, k[0]) }, true
	case *DynLetExpr:
		// Dynamic lets bind @-identifiers, a namespace apart from patterns.
		return once(e.bindings, e.expr), func(k []Expr) Expr { return NewDynLetExpr(e.Src, k[0], k[1]) }, true
	case *SafeTailExpr:
		kids = []Child{{Expr: e.base, Kind: ChildOnce}, {Expr: e.fallbackValue, Kind: ChildMaybe}}
		for _, s := range e.steps {
			for _, a := range s.args {
				kids = append(kids, Child{Expr: a, Kind: ChildMaybe})
			}
		}
		return kids, func(k []Expr) Expr {
			n := *e
			n.base, n.fallbackValue = k[0], k[1]
			n.steps = make([]SafeTailStep, len(e.steps))
			j := 2
			for i, s := range e.steps {
				ns := s
				ns.args = make([]Expr, len(s.args))
				for a := range s.args {
					ns.args[a] = k[j]
					j++
				}
				n.steps[i] = ns
			}
			return &n
		}, true
	case ExprExpr:
		return []Child{{Expr: e.Expr, Kind: ChildMany}}, func(k []Expr) Expr { return NewExprExpr(e.Src, k[0]) }, true
	case StringCharTupleExpr:
		return once(e.at, e.char), func(k []Expr) Expr { return NewStringCharTupleExpr(e.Src, k[0], k[1]) }, true
	case ArrayItemTupleExpr:
		return once(e.at, e.item), func(k []Expr) Expr { return NewArrayItemTupleExpr(e.Src, k[0], k[1]) }, true
	case DictEntryTupleExpr:
		return once(e.at, e.value), func(k []Expr) Expr { return NewDictEntryTupleExpr(e.Src, k[0], k[1]) }, true
	case Value:
		return nil, func([]Expr) Expr { return e }, true
	}
	return nil, nil, false
}

func dictChildren(e DictExpr, ptr bool) ([]Child, func([]Expr) Expr, bool) {
	kids := make([]Child, 0, 2*len(e.entryExprs))
	for _, ent := range e.entryExprs {
		kids = append(kids, Child{Expr: ent.at, Kind: ChildOnce}, Child{Expr: ent.value, Kind: ChildOnce})
	}
	return kids, func(k []Expr) Expr {
		entries := make([]DictEntryTupleExpr, len(e.entryExprs))
		for i, ent := range e.entryExprs {
			entries[i] = DictEntryTupleExpr{ExprScanner: ent.ExprScanner, at: k[2*i], value: k[2*i+1]}
		}
		d, err := NewDictExpr(e.Src, e.allowDupKeys, true, entries...)
		if err != nil {
			if ptr {
				return &e
			}
			return e
		}
		if ptr {
			if v, is := d.(DictExpr); is {
				return &v
			}
		}
		return d
	}, true
}

// rebuildBin re-runs the constructor for operators that analyse their
// operands when built; every other BinExpr is a closure over its operator
// alone, so a copy with new operands is faithful.
func rebuildBin(e *BinExpr, a, b Expr) Expr {
	switch e.op {
	case opWhere:
		return NewWhereExpr(e.Src, a, b)
	case "orderby":
		return NewOrderByExpr(e.Src, a, b)
	case "order":
		return NewOrderExpr(e.Src, a, b)
	case "rank":
		return NewRankExpr(e.Src, a, b)
	case opCall:
		return NewCallExpr(e.Src, a, b)
	case "+>":
		return NewAddArrowExpr(e.Src, a, b)
	}
	n := *e
	n.a, n.b = a, b
	return &n
}

// patternIdents returns the identifiers a pattern binds.
func patternIdents(p Pattern) []string {
	return p.Bindings()
}

// patternExprs returns the expressions a pattern evaluates while matching:
// fallback values, literal sub-patterns and dictionary keys.
func patternExprs(p Pattern) []Expr {
	switch p := p.(type) {
	case ExprPattern:
		return []Expr{p.Expr}
	case ExprsPattern:
		return p.exprs
	case TuplePattern:
		var es []Expr
		for _, a := range p.attrs {
			es = append(es, fallbackExprs(a.pattern)...)
		}
		return es
	case ArrayPattern:
		var es []Expr
		for _, it := range p.items {
			es = append(es, fallbackExprs(it)...)
		}
		return es
	case SetPattern:
		var es []Expr
		for _, sp := range p.patterns {
			es = append(es, patternExprs(sp)...)
		}
		return es
	case DictPattern:
		var es []Expr
		for _, ent := range p.entries {
			es = append(es, ent.at)
			es = append(es, fallbackExprs(ent.pattern)...)
		}
		return es
	}
	return nil
}

func fallbackExprs(p FallbackPattern) []Expr {
	var es []Expr
	if p.pattern != nil {
		es = patternExprs(p.pattern)
	}
	if p.fallback != nil {
		es = append(es, p.fallback)
	}
	return es
}

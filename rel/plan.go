package rel

// edgeDecision is the S5 per-edge stream vs materialize choice (🎯T23).
type edgeDecision int

const (
	streamEdge edgeDecision = iota
	materializeEdge
)

func (d edgeDecision) String() string {
	if d == materializeEdge {
		return "materialize"
	}
	return "stream"
}

// decideLetFanout records the S5 rule: a let-bound relation consumed more
// than once is materialised; a single consumer may stream.
func decideLetFanout(uses int) edgeDecision {
	if uses > 1 {
		return materializeEdge
	}
	return streamEdge
}

// zoneMap is a lazy min/max over a numeric column, used as a planner
// statistic (index-as-materialization-format).
type zoneMap struct {
	min, max Number
	ok       bool
}

func (r *positionalRelation) zone(col int) zoneMap {
	if col < 0 || col >= r.store.width || r.n == 0 {
		return zoneMap{}
	}
	m := r.getMeta()
	m.Lock()
	if m.zones != nil && col < len(m.zones) && m.zones[col] != nil {
		z := *m.zones[col]
		m.Unlock()
		return z
	}
	m.Unlock()
	z := r.computeZone(col)
	m.Lock()
	if m.zones == nil {
		m.zones = make([]*zoneMap, r.store.width)
	}
	m.zones[col] = &z
	m.Unlock()
	return z
}

func (r *positionalRelation) computeZone(col int) zoneMap {
	first, ok := r.rowAt(0)[col].(Number)
	if !ok {
		return zoneMap{}
	}
	z := zoneMap{min: first, max: first, ok: true}
	for i := 1; i < r.n; i++ {
		n, ok := r.rowAt(i)[col].(Number)
		if !ok {
			return zoneMap{}
		}
		if n.Less(z.min) {
			z.min = n
		}
		if z.max.Less(n) {
			z.max = n
		}
	}
	return z
}

// planCache is a fact-keyed cache: a hit is returned only while guard() is
// true, so stale certificates cannot reuse a plan.
type planCache struct {
	key   string
	guard func() bool
	value Value
}

func (c *planCache) get(key string, guard func() bool) (Value, bool) {
	if c == nil || c.key != key || c.value == nil || c.guard == nil || !c.guard() || !guard() {
		return nil, false
	}
	return c.value, true
}

func (c *planCache) put(key string, guard func() bool, v Value) {
	c.key = key
	c.guard = guard
	c.value = v
}

// recordedFanout is the S5 plan choice for a \ident body: materialize when
// the ident is consumed more than once.
func (f *Function) recordedFanout() edgeDecision {
	ident, ok := f.arg.(IdentPattern)
	if !ok {
		return streamEdge
	}
	return decideLetFanout(countIdentUses(f.body, string(ident)))
}

// demandedEqAttrs returns attributes of ident used in `ident where .attr = k`.
func demandedEqAttrs(e Expr, ident string) []string {
	var attrs []string
	addDemandedEqAttrs(e, ident, &attrs)
	return attrs
}

func addDemandedEqAttrs(e Expr, ident string, attrs *[]string) {
	if e == nil {
		return
	}
	switch e := e.(type) {
	case *BinExpr:
		if e.op == "where" {
			if id, ok := e.a.(IdentExpr); ok && id.ident == ident {
				if f, ok := e.b.(*Function); ok {
					if p := matchEqAttrPredicate(f); p != nil {
						*attrs = append(*attrs, p.attr)
					}
				}
			}
		}
		addDemandedEqAttrs(e.a, ident, attrs)
		addDemandedEqAttrs(e.b, ident, attrs)
	case *UnaryExpr:
		addDemandedEqAttrs(e.a, ident, attrs)
	case *DArrowExpr:
		addDemandedEqAttrs(e.lhs, ident, attrs)
	case *ArrowExpr:
		addDemandedEqAttrs(e.lhs, ident, attrs)
		if !shadowsIdent(e.fn.arg, ident) {
			addDemandedEqAttrs(e.fn.body, ident, attrs)
		}
	case AndExpr:
		addDemandedEqAttrs(e.a, ident, attrs)
		addDemandedEqAttrs(e.b, ident, attrs)
	case OrExpr:
		addDemandedEqAttrs(e.a, ident, attrs)
		addDemandedEqAttrs(e.b, ident, attrs)
	}
}

// materializeValue is the pipeline-breaker: force a suspended sequence/dict
// pipeline, and for a relation freeze the demanded group indexes (index-as-
// materialization-format) plus zone maps.
func materializeValue(v Value, attrs []string) (Value, error) {
	switch v := v.(type) {
	case seqPipeline:
		return v.force()
	case dictPipeline:
		return v.force()
	case Relation:
		for _, attr := range attrs {
			if i := v.getAttrIndex(attr); i >= 0 {
				_ = v.rows.groupBy(valueProjector{i})
				_ = v.rows.zone(i)
			}
		}
		return v, nil
	default:
		return v, nil
	}
}

func shadowsIdent(p Pattern, ident string) bool {
	ip, ok := p.(IdentPattern)
	return ok && string(ip) == ident
}

func countIdentUses(e Expr, ident string) int {
	n := 0
	addIdentUses(e, ident, &n)
	return n
}

func addIdentUsesAll(ident string, n *int, es ...Expr) {
	for _, e := range es {
		addIdentUses(e, ident, n)
	}
}

func addIdentUses(e Expr, ident string, n *int) {
	if e == nil {
		return
	}
	switch e := e.(type) {
	case IdentExpr:
		if e.ident == ident {
			*n++
		}
	case *BinExpr:
		addIdentUsesAll(ident, n, e.a, e.b)
	case *UnaryExpr:
		addIdentUses(e.a, ident, n)
	case *DotExpr:
		addIdentUses(e.lhs, ident, n)
	case AndExpr:
		addIdentUsesAll(ident, n, e.a, e.b)
	case OrExpr:
		addIdentUsesAll(ident, n, e.a, e.b)
	case *IfElseExpr:
		addIdentUsesAll(ident, n, e.ifTrue, e.cond, e.ifFalse)
	case CompareExpr:
		addIdentUsesAll(ident, n, e.args...)
	case *Function:
		if !shadowsIdent(e.arg, ident) {
			addIdentUses(e.body, ident, n)
		}
	case *ArrowExpr:
		addIdentUses(e.lhs, ident, n)
		if !shadowsIdent(e.fn.arg, ident) {
			addIdentUses(e.fn.body, ident, n)
		}
	case *DArrowExpr:
		addIdentUses(e.lhs, ident, n)
		if !shadowsIdent(e.fn.arg, ident) {
			addIdentUses(e.fn.body, ident, n)
		}
	case *SeqArrowExpr:
		addIdentUses(e.lhs, ident, n)
		if !shadowsIdent(e.fn.arg, ident) {
			addIdentUses(e.fn.body, ident, n)
		}
	default:
		addIdentUsesCompound(e, ident, n)
	}
}

func addIdentUsesCompound(e Expr, ident string, n *int) {
	switch e := e.(type) {
	case *TupleExpr:
		for _, a := range e.attrs {
			addIdentUses(a.expr, ident, n)
		}
	case *ArrayExpr:
		addIdentUsesAll(ident, n, e.elements...)
	case *SetExpr:
		addIdentUsesAll(ident, n, e.elements...)
	case DictExpr:
		for i := range e.entryExprs {
			addIdentUsesAll(ident, n, e.entryExprs[i].at, e.entryExprs[i].value)
		}
	case *DictExpr:
		for i := range e.entryExprs {
			addIdentUsesAll(ident, n, e.entryExprs[i].at, e.entryExprs[i].value)
		}
	case CondExpr:
		addIdentUses(e.dicExpr, ident, n)
	case *DynLetExpr:
		addIdentUsesAll(ident, n, e.bindings, e.expr)
	}
}

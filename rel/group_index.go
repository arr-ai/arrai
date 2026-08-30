package rel

import (
	"github.com/arr-ai/frozen"
)

// groupIndex maps a projection of a relation's rows to the rows sharing it,
// and backs joins and index-answered `where`.
//
// A single-column projection is stored in a Map keyed by Value rather than
// by any. That matters because frozen compares an `any` key without
// reflection only if it satisfies Equaler[any]: an arr.ai Value has
// Equal(Value) bool, so as an `any` key it was silently wrong before frozen
// v1.14.0 and reflection-slow after. Keying the map by Value instead lets
// frozen match Equaler[Value] directly, and stores the key inline rather
// than boxing it — one allocation per key rather than three.
//
// Wider projections key on the projected row as before.
type groupIndex struct {
	byValue frozen.Map[Value, frozen.Set[any]]
	byRow   frozen.Map[any, frozen.Set[any]]
	single  bool
}

// groupKey is a key in whichever form its index uses. Only the index
// constructs one, from a projector of the same width as its own, so the two
// can never disagree about the encoding. It is passed by value and never
// stored, so it costs no allocation.
type groupKey struct {
	v      Value
	row    any
	single bool
}

// keyFrom encodes p's projection of row as a key for this index. p must have
// the same width as the projector the index was built from — true of any
// join's two sides — but need not be the same projector: the two sides of a
// join name the shared columns at different positions.
func (g groupIndex) keyFrom(p valueProjector, row Values) groupKey {
	if g.single {
		return groupKey{v: row[p[0]], single: true}
	}
	return groupKey{row: p.rowKey(row)}
}

// rowKey encodes the projection of a row as a row-shaped key.
func (p valueProjector) rowKey(row Values) interface{} {
	switch len(p) {
	case 0:
		// Disjoint join keys: every row shares the one empty key.
		return Values{}
	case 1:
		return Values{row[p[0]]}
	}
	if p.isContiguous() {
		a, b := p[0], p[len(p)-1]+1
		v := make(Values, b-a)
		copy(v, row[a:b])
		return v
	}
	return row.project(p).values()
}

// newGroupIndex groups a row set by p.
func newGroupIndex(rows frozen.Set[any], p valueProjector) groupIndex {
	switch len(p) {
	case 0:
		return groupIndex{
			byRow: frozen.NewMap[any, frozen.Set[any]](
				frozen.KV[any, frozen.Set[any]](Values{}, rows)),
		}
	case 1:
		if fastPaths {
			i := p[0]
			return groupIndex{
				single: true,
				byValue: frozen.SetGroupBy(rows, func(el any) Value {
					return el.(Values)[i] //nolint:forcetypeassert
				}),
			}
		}
	}
	return groupIndex{byRow: frozen.SetGroupBy(rows, func(el any) any {
		return p.rowKey(el.(Values)) //nolint:forcetypeassert
	})}
}

// get returns the rows sharing k.
func (g groupIndex) get(k groupKey) (frozen.Set[any], bool) {
	if g.single {
		return g.byValue.Get(k.v)
	}
	return g.byRow.Get(k.row)
}

// has reports whether any row shares k.
func (g groupIndex) has(k groupKey) bool {
	if g.single {
		return g.byValue.Has(k.v)
	}
	return g.byRow.Has(k.row)
}

// count returns the number of distinct keys.
func (g groupIndex) count() int {
	if g.single {
		return g.byValue.Count()
	}
	return g.byRow.Count()
}

// each calls f for every group, in unspecified order.
func (g groupIndex) each(f func(k groupKey, rows frozen.Set[any])) {
	if g.single {
		for i := g.byValue.Range(); i.Next(); {
			f(groupKey{v: i.Key(), single: true}, i.Value())
		}
		return
	}
	for i := g.byRow.Range(); i.Next(); {
		f(groupKey{row: i.Key()}, i.Value())
	}
}

// commonKeyRows returns, for each key present in both indexes, the row that
// key projects from. Both indexes must have been built over projectors of
// the same width, which is true of any join's two sides.
func (g groupIndex) commonKeyRows(h groupIndex, toRow func(k groupKey) Values) frozen.Set[any] {
	sb := frozen.NewSetBuilder[any](0)
	if g.single && h.single {
		for i := g.byValue.Keys().Intersection(h.byValue.Keys()).Range(); i.Next(); {
			sb.Add(toRow(groupKey{v: i.Value(), single: true}))
		}
		return sb.Finish()
	}
	for i := g.byRow.Keys().Intersection(h.byRow.Keys()).Range(); i.Next(); {
		sb.Add(toRow(groupKey{row: i.Value()}))
	}
	return sb.Finish()
}

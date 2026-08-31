package rel

import (
	"github.com/arr-ai/hash/hash128"
)

// groupIndex maps each distinct projection of a view's rows to the rows
// sharing it, and backs joins and index-answered `where`. Buckets are keyed
// by the 128-bit hash of the projected values, and each bucket remembers one
// representative row — the key *is* that row's projection, so keys are never
// materialised. Distinct keys colliding on their hash chain through next;
// with 128-bit hashes that is theoretical, but probes still verify the key
// against the representative rather than trusting the hash.
//
// Rows are recorded as arena ids of the view's store, ascending, so a bucket
// can back a view of the same store directly (selView).
type groupIndex struct {
	rel *positionalRelation
	p   valueProjector
	m   map[hash128.H128]*groupBucket
}

type groupBucket struct {
	rep  uint32   // arena id of a representative row
	rows []uint32 // ascending arena ids of the group's rows
	next *groupBucket
}

// hashOf hashes p's projection of row, matching projectedValues.Hash128 and
// the hash of the projected values as a plain row.
func (p valueProjector) hashOf(row Values) hash128.H128 {
	h := valuesSalt
	for _, i := range p {
		h = h.Mix(row[i].Hash128())
	}
	return h
}

// newGroupIndex groups a view's rows by p. Callers go through
// positionalRelation.groupBy, which memoises per view.
func newGroupIndex(r *positionalRelation, p valueProjector) *groupIndex {
	g := &groupIndex{rel: r, p: p, m: make(map[hash128.H128]*groupBucket, r.n)}
	for i := 0; i < r.n; i++ {
		row := r.rowAt(i)
		h := p.hashOf(row)
		b := g.m[h]
		for ; b != nil; b = b.next {
			if projectionsEqual(g.repRow(b), p, row, p) {
				break
			}
		}
		if b == nil {
			b = &groupBucket{rep: r.arenaID(i), next: g.m[h]}
			g.m[h] = b
		}
		b.rows = append(b.rows, r.arenaID(i))
	}
	return g
}

// projectionsEqual reports whether pa's projection of a equals pb's
// projection of b. The projectors must have the same length.
func projectionsEqual(a Values, pa valueProjector, b Values, pb valueProjector) bool {
	for i := range pa {
		if !a[pa[i]].Equal(b[pb[i]]) {
			return false
		}
	}
	return true
}

// row returns the store row at an arena id recorded in this index.
func (g *groupIndex) row(id uint32) Values {
	return rowOf(g.rel.arena, g.rel.store.width, int(id))
}

// repRow returns a bucket's representative row: the full row, of which the
// group key is the projection.
func (g *groupIndex) repRow(b *groupBucket) Values {
	return g.row(b.rep)
}

// bucketOfRow returns the bucket whose key equals q's projection of row, or
// nil. q must have the same length as the index's projector; probing one
// side of a join with the other side's rows uses that side's projector.
func (g *groupIndex) bucketOfRow(row Values, q valueProjector) *groupBucket {
	for b := g.m[q.hashOf(row)]; b != nil; b = b.next {
		if projectionsEqual(g.repRow(b), g.p, row, q) {
			return b
		}
	}
	return nil
}

// getKey returns the arena ids of rows whose projection equals key, given in
// the same order as the index's projector.
func (g *groupIndex) getKey(key ...Value) ([]uint32, bool) {
	h := valuesSalt
	for _, v := range key {
		h = h.Mix(v.Hash128())
	}
	for b := g.m[h]; b != nil; b = b.next {
		rep := g.repRow(b)
		match := true
		for i, v := range key {
			if !rep[g.p[i]].Equal(v) {
				match = false
				break
			}
		}
		if match {
			return b.rows, true
		}
	}
	return nil, false
}

// count returns the number of distinct keys.
func (g *groupIndex) count() int {
	n := 0
	for _, b := range g.m {
		for ; b != nil; b = b.next {
			n++
		}
	}
	return n
}

// each calls f for every group, in unspecified order.
func (g *groupIndex) each(f func(b *groupBucket)) {
	for _, b := range g.m {
		for ; b != nil; b = b.next {
			f(b)
		}
	}
}

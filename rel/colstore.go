package rel

import (
	"sync"

	"github.com/arr-ai/hash/hash128"
)

// colStore holds relation rows in a single flat, append-only arena: row i
// occupies arena[i*width : (i+1)*width]. A positionalRelation is a view of a
// store — a committed prefix, or a selection vector over one — so building a
// relation is one allocation, a `where` shares the arena and allocates only
// the selection, and a row *is* a slice of the arena: inflating it to a
// tuple copies nothing.
//
// The arena only ever grows, and rows below a view's length never change, so
// views read their own arena snapshot without locking. Appends, and the lazy
// hash index over committed rows, are guarded by mu.
//
// Committed rows are always distinct. The index maps a row's 128-bit hash to
// its position; it is built the first time membership is asked for (Has,
// With, dedupe) and extended incrementally as rows are appended, so chains
// of With calls — `acc | {x}` folds — pay for hashing once per row, not once
// per generation. hashOverflow holds rows whose hash collides with an
// earlier row's; with 128-bit hashes that is theoretical, but lookups still
// verify equality rather than trusting the hash.
type colStore struct {
	width int

	mu     sync.Mutex
	arena  []Value
	n      int // committed rows; rows [0, n) are distinct
	hashes []hash128.H128
	index  map[hash128.H128]uint32
	// hashOverflow holds ids of rows whose hash collides with index[h].
	hashOverflow map[hash128.H128][]uint32
}

// rowOf returns row id of an arena snapshot. The result is capped so no
// append through it can ever touch the arena.
func rowOf(arena []Value, width, id int) Values {
	a := id * width
	b := a + width
	return Values(arena[a:b:b])
}

// ensureIndexLocked hashes and indexes any committed rows not yet indexed.
// Caller holds mu.
func (s *colStore) ensureIndexLocked() {
	if s.index == nil {
		s.index = make(map[hash128.H128]uint32, s.n)
	}
	for i := len(s.hashes); i < s.n; i++ {
		h := rowOf(s.arena, s.width, i).Hash128()
		s.hashes = append(s.hashes, h)
		s.insertLocked(h, uint32(i))
	}
}

// insertLocked adds an index entry for a row known not to equal any indexed
// row. Caller holds mu; the row must already be hashed into s.hashes.
func (s *colStore) insertLocked(h hash128.H128, id uint32) {
	if _, exists := s.index[h]; !exists {
		s.index[h] = id
		return
	}
	if s.hashOverflow == nil {
		s.hashOverflow = map[hash128.H128][]uint32{}
	}
	s.hashOverflow[h] = append(s.hashOverflow[h], id)
}

// findLocked returns the id of the row equal to v among the first limit
// committed rows. Caller holds mu and has called ensureIndexLocked.
func (s *colStore) findLocked(v Values, h hash128.H128, limit int) (uint32, bool) {
	if id, ok := s.index[h]; ok && int(id) < limit {
		if hashIdentity || rowOf(s.arena, s.width, int(id)).equalValues(v) {
			return id, true
		}
	}
	for _, id := range s.hashOverflow[h] {
		if int(id) < limit && (hashIdentity || rowOf(s.arena, s.width, int(id)).equalValues(v)) {
			return id, true
		}
	}
	return 0, false
}

// storeBuilder accumulates rows into a fresh store. With dedupe, rows equal
// to one already added are dropped, and the index built for that is kept on
// the store, so membership queries on the result are already paid for.
// Without dedupe the caller guarantees distinctness and rows are appended
// unhashed; the index is built lazily if ever needed. Builders are
// single-owner: no locking until finish publishes the store.
type storeBuilder struct {
	s      *colStore
	dedupe bool
}

func newStoreBuilder(width, capacity int, dedupe bool) storeBuilder {
	s := &colStore{width: width}
	if capacity > 0 {
		s.arena = make([]Value, 0, capacity*width)
	}
	if dedupe {
		s.index = make(map[hash128.H128]uint32, capacity)
	}
	return storeBuilder{s: s, dedupe: dedupe}
}

// add appends a copy of row v.
func (b storeBuilder) add(v Values) {
	s := b.s
	s.arena = append(s.arena, v...)
	b.commit()
}

// addProjection appends row's values at p's positions, without materialising
// the projected row anywhere but the arena.
func (b storeBuilder) addProjection(row Values, p valueProjector) {
	s := b.s
	for _, i := range p {
		s.arena = append(s.arena, row[i])
	}
	b.commit()
}

// addJoined appends left's values at lp's positions followed by right's at
// rp's: one output row of a join.
func (b storeBuilder) addJoined(left Values, lp valueProjector, right Values, rp valueProjector) {
	s := b.s
	for _, i := range lp {
		s.arena = append(s.arena, left[i])
	}
	for _, i := range rp {
		s.arena = append(s.arena, right[i])
	}
	b.commit()
}

// commit accepts or rolls back the row just appended to the arena.
func (b storeBuilder) commit() {
	s := b.s
	if !b.dedupe {
		s.n++
		return
	}
	v := rowOf(s.arena, s.width, s.n)
	h := v.Hash128()
	if _, found := s.findLocked(v, h, s.n); found {
		s.arena = s.arena[:s.n*s.width]
		return
	}
	s.hashes = append(s.hashes, h)
	s.insertLocked(h, uint32(s.n))
	s.n++
}

func (b storeBuilder) finish() *positionalRelation {
	s := b.s
	return &positionalRelation{store: s, arena: s.arena, n: s.n}
}

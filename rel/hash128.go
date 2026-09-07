package rel

import (
	"sync"

	"github.com/arr-ai/hash/hash128"
)

// Values hash through hash128: each Value computes its 128-bit hash once
// (composites cache it) and frozen consumes it in a single pass. The seeded
// Hash(seed) methods remain for callers of the older interface and are
// derived from Hash128, so they never re-walk a value.
//
// Every tuple kind hashes with the same formula, so a GenericTuple and a
// specialised tuple with the same attributes hash the same:
//
//	tupleSalt ^ (hash(name₁) ⋈ hash(value₁)) ^ (hash(name₂) ⋈ hash(value₂)) ^ …
//
// Sets of tuples in turn xor their elements' hashes (with a per-kind salt),
// so a set's hash is independent of element order.

// Salts distinguish empty or structurally similar values of different
// kinds, and name hashes serve the specialised tuple kinds. They are
// computed at package initialisation because hash128 is randomised per
// process; hash128 initialises first because rel imports it.
var (
	tupleSalt    = hash128.String("rel.Tuple")
	valuesSalt   = hash128.String("rel.Values")
	stringSalt   = hash128.String("rel.String")
	bytesSalt    = hash128.String("rel.Bytes")
	emptySetSalt = hash128.String("rel.EmptySet")
	arraySalt    = hash128.String("rel.Array")
	unionSalt    = hash128.String("rel.UnionSet")
	relationSalt = hash128.String("rel.Relation")

	atNameHash    = hash128.String("@")
	itemNameHash  = hash128.String(ArrayItemAttr)
	charNameHash  = hash128.String(StringCharAttr)
	byteNameHash  = hash128.String(BytesByteAttr)
	valueNameHash = hash128.String(DictValueAttr)
)

// hashAttr hashes one attribute of a tuple.
func hashAttr(nameHash hash128.H128, value Value) hash128.H128 {
	return nameHash.Mix(value.Hash128())
}

// hashTuple2 hashes a two-attribute tuple whose attribute-name hashes are
// precomputed.
func hashTuple2(name1 hash128.H128, v1 Value, name2 hash128.H128, v2 Value) hash128.H128 {
	return tupleSalt.Xor(hashAttr(name1, v1)).Xor(hashAttr(name2, v2))
}

// hashCell memoises a composite value's hash. Composites are immutable, so
// one computation is valid for the value's lifetime; the cell is shared by
// pointer between copies of the value.
type hashCell struct {
	once sync.Once
	h    hash128.H128
}

func (c *hashCell) get(compute func() hash128.H128) hash128.H128 {
	c.once.Do(func() { c.h = compute() })
	return c.h
}

package rel

import (
	"hash/maphash"
	"math"
	"sync"
	"unicode/utf16"
)

// Values hash in 64 bits, seedless. Mix chains ordered parts (rows, arrays,
// name⋈value). Non-frozen unordered collections xor their parts and wrap
// with mix (setSalt or tupleSalt). Frozen sets and maps already wrap H0;
// do not wrap those again. {} / {||} / empty Dict use the empty frozen
// set hash; () uses the tuple wrap.

var hashSeed = maphash.MakeSeed()

const (
	golden     = 0x9E3779B97F4A7C15
	mixConst   = 0x2545F4914F6CDD1D
	fmixC1     = 0xff51afd7ed558ccd
	fmixC2     = 0xc4ceb9fe1a85ec53
	floatSalt  = 0x9E3779B97F4A7C15
	intSalt    = 0xC2B2AE3D27D4EB4F
	runesSalt  = 0x165667B19E3779F9
)

func fmix64(k uint64) uint64 {
	k ^= k >> 33
	k *= fmixC1
	k ^= k >> 33
	k *= fmixC2
	k ^= k >> 33
	return k
}

// mix folds o into h asymmetrically. mix(a, b) != mix(b, a) in general.
func mix(h, o uintptr) uintptr {
	x := uint64(h) ^ (uint64(o)*golden + mixConst)
	return uintptr(fmix64(x))
}

func xor(h, o uintptr) uintptr { return h ^ o }

func hashSet(xorElems uintptr) uintptr {
	return mix(setSalt, xorElems)
}

func hashTuple(xorAttrs uintptr) uintptr {
	return mix(tupleSalt, xorAttrs)
}

func hashBytes(b []byte) uintptr {
	return uintptr(maphash.Bytes(hashSeed, b))
}

func hashString(s string) uintptr {
	return uintptr(maphash.String(hashSeed, s))
}

func hashInt(i int) uintptr {
	return mix(intSalt, uintptr(i))
}

func hashFloat64(f float64) uintptr {
	if f == 0 {
		f = 0
	}
	return mix(floatSalt, uintptr(math.Float64bits(f)))
}

func hashRunes(s []rune) uintptr {
	h := uintptr(runesSalt)
	for _, r := range utf16.Encode(s) {
		h = mix(h, uintptr(r))
	}
	return h
}

// Salts distinguish empty or structurally similar values of different kinds.
// Name hashes serve the specialised tuple kinds.
var (
	tupleSalt      = hashString("rel.Tuple")
	setSalt        = hashString("rel.Set")
	valuesSalt     = hashString("rel.Values")
	stringSalt     = hashString("rel.String")
	bytesSalt      = hashString("rel.Bytes")
	arraySalt      = hashString("rel.Array")
	posRelSalt     = hashString("github.com/arr-ai/arrai/rel.positionalRelation")
	funcSalt       = hashString("rel.Function")
	nativeFuncSalt = hashString("rel.NativeFunction")

	atNameHash    = hashString("@")
	itemNameHash  = hashString(ArrayItemAttr)
	charNameHash  = hashString(StringCharAttr)
	byteNameHash  = hashString(BytesByteAttr)
	valueNameHash = hashString(DictValueAttr)
)

func hashAttr(nameHash uintptr, value Value) uintptr {
	return mix(nameHash, value.Hash())
}

func hashTuple2(name1 uintptr, v1 Value, name2 uintptr, v2 Value) uintptr {
	return hashTuple(xor(hashAttr(name1, v1), hashAttr(name2, v2)))
}

// hashCell memoises a composite value's hash. Composites are immutable, so
// one computation is valid for the value's lifetime; the cell is shared by
// pointer between copies of the value.
type hashCell struct {
	once sync.Once
	h    uintptr
}

func (c *hashCell) get(compute func() uintptr) uintptr {
	c.once.Do(func() { c.h = compute() })
	return c.h
}

package rel

import (
	"bytes"
	"sync"
	"unicode/utf8"

	"github.com/arr-ai/hash/hash128"

	"context"
	"fmt"
	"reflect"

	"github.com/arr-ai/wbnf/parser"

	"github.com/arr-ai/arrai/pkg/fu"
)

// StringCharAttr is the standard name for the value-attr of a character tuple.
const StringCharAttr = "@char"

// String is a set of Values.
//
// Two backings, exactly one active. ASCII content — far and away the common
// case — lives in ascii, one byte per rune, so byte index equals rune index
// and hashing, comparison and Go-string conversion all work on a quarter of
// the memory. Everything else (non-ASCII content, and strings carrying
// holes from character-level Without) lives in s as []rune. Constructors
// normalise: contiguous hole-free ASCII content is always in ascii form.
// Character-level edits (with/Without) are rare in practice and take the
// rune form; the hot operations all have byte paths.
type String struct {
	ascii  []byte // active when non-nil; all bytes < 0x80, holes == 0
	s      []rune // active when ascii is nil
	offset int
	holes  int

	// buf/abuf, when non-nil, is the shared append buffer this string is a
	// prefix of. A chain of concatenations extends one buffer in place
	// (amortised O(1) per element) instead of copying the accumulator per
	// step; branching from an older string copies out. Elements below any
	// string's length never change.
	buf  *appendBuf[rune]
	abuf *appendBuf[byte]
}

// appendBuf tracks ownership of a growable buffer shared by the sequences
// concatenation builds (string runes/bytes, array values). n is the
// committed length: only the sequence whose length equals n may extend the
// buffer, everyone else copies.
type appendBuf[E any] struct {
	mu    sync.Mutex
	elems []E
	n     int
}

// extend appends more onto a frontier of length n, returning the grown
// prefix, or nil if this buffer's frontier has moved past n.
func (b *appendBuf[E]) extend(n int, more []E) []E {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.n != n {
		return nil
	}
	b.elems = append(b.elems[:n], more...)
	b.n = n + len(more)
	return b.elems[:b.n]
}

// newAppendBuf starts a buffer owning the concatenation of a and b, with
// room to grow.
func newAppendBuf[E any](a, b []E) (*appendBuf[E], []E) {
	elems := append(append(make([]E, 0, 2*(len(a)+len(b))), a...), b...)
	return &appendBuf[E]{elems: elems, n: len(elems)}, elems
}

// concatStrings concatenates two contiguous zero-offset strings, extending
// a's buffer in place when a is its frontier.
func concatStrings(a, b String) String {
	if a.ascii != nil && b.ascii != nil {
		if a.abuf != nil {
			if s := a.abuf.extend(len(a.ascii), b.ascii); s != nil {
				return String{ascii: s, abuf: a.abuf}
			}
		}
		abuf, s := newAppendBuf(a.ascii, b.ascii)
		return String{ascii: s, abuf: abuf}
	}
	ar, br := a.runes(), b.runes()
	if a.buf != nil {
		if s := a.buf.extend(len(ar), br); s != nil {
			return String{s: s, buf: a.buf}
		}
	}
	buf, s := newAppendBuf(ar, br)
	return String{s: s, buf: buf}
}

// size returns the backing length, holes included.
func (s String) size() int {
	if s.ascii != nil {
		return len(s.ascii)
	}
	return len(s.s)
}

// runeAt returns the rune at backing index i (-1 for a hole).
func (s String) runeAt(i int) rune {
	if s.ascii != nil {
		return rune(s.ascii[i])
	}
	return s.s[i]
}

// runes materialises the backing as []rune.
func (s String) runes() []rune {
	if s.ascii == nil {
		return s.s
	}
	r := make([]rune, len(s.ascii))
	for i, b := range s.ascii {
		r[i] = rune(b)
	}
	return r
}

// goString materialises the content as a Go string. Holes become
// replacement runes, as they always have.
func (s String) goString() string {
	if s.ascii != nil {
		return string(s.ascii)
	}
	return string(s.s)
}

// asRuneForm returns s backed by runes, for the rare character-level edits.
func (s String) asRuneForm() String {
	if s.ascii == nil {
		return s
	}
	return String{s: s.runes(), offset: s.offset}
}

// NewString constructs a string as a relation.
func NewString(s []rune) Set {
	return NewOffsetString(s, 0)
}

// NewGoString constructs a string from a Go string: the cheap path for
// ASCII content, which stays in byte form without a rune round-trip.
func NewGoString(s string) Set {
	if len(s) == 0 {
		return None
	}
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			return NewString([]rune(s))
		}
	}
	return String{ascii: []byte(s)}
}

// internedStrings canonicalises compile-time string literals so every
// occurrence of the same literal shares one backing array, letting
// EqualString's same-backing fast path answer literal-vs-literal
// comparisons without scanning. Bounded by program text, so entries live
// for the process.
var internedStrings sync.Map // string -> Set

// InternedGoString is NewGoString for compile-time literals: the same
// content always returns the same value.
func InternedGoString(s string) Set {
	if v, ok := internedStrings.Load(s); ok {
		return v.(Set)
	}
	v, _ := internedStrings.LoadOrStore(s, NewGoString(s))
	return v.(Set)
}

// NewOffsetString constructs an offset string as a relation.
func NewOffsetString(s []rune, offset int) Set {
	if len(s) == 0 {
		return None
	}
	holes := 0
	ascii := true
	for _, r := range s {
		if r < 0 {
			holes++
		} else if r >= utf8.RuneSelf {
			ascii = false
		}
	}
	if ascii && holes == 0 {
		b := make([]byte, len(s))
		for i, r := range s {
			b[i] = byte(r)
		}
		return String{ascii: b, offset: offset}
	}
	return String{s: s, offset: offset, holes: holes}
}

func asString(values ...Value) String {
	n := len(values)
	tuples := make([]StringCharTuple, 0, n)
	minAt := int(^uint(0) >> 1)
	maxAt := -minAt - 1
	for _, v := range values {
		t := v.(StringCharTuple)
		if minAt > t.at {
			minAt = t.at
		}
		if maxAt < t.at {
			maxAt = t.at
		}
		tuples = append(tuples, t)
	}
	str := make([]rune, maxAt-minAt+1)
	for i := range str {
		str[i] = -1
	}
	for _, t := range tuples {
		str[t.at-minAt] = t.char
	}
	return NewOffsetString(str, minAt).(String)
}

// AsString returns String and the empty set as String or false otherwise.
func AsString(v Value) (String, bool) {
	switch s := v.(type) {
	case String:
		return s, true
	case Set:
		if !s.IsTrue() {
			return String{}, true
		}
	}
	return String{}, false
}

// Hash computes a hash for a String.
func (s String) Hash(seed uintptr) uintptr {
	return s.Hash128().Seeded(seed)
}

// Hash128 computes the 128-bit hash of a String over its content's UTF-8
// encoding, so both backings of the same content hash identically, salted
// so a String never hashes like the Bytes with the same content. Strings
// with holes cannot round-trip through UTF-8 (a hole is not a rune) and
// only ever equal other rune-form strings, so they hash the rune buffer.
func (s String) Hash128() hash128.H128 {
	h := stringSalt.Mix(hash128.Int(s.offset))
	if s.ascii != nil {
		return h.Mix(hash128.Bytes(s.ascii))
	}
	if s.holes != 0 {
		return h.Mix(hash128.Runes(s.s))
	}
	return h.Mix(hash128.String(string(s.s)))
}

// Equal tests two Sets for equality. Any other type returns false.
func (s String) Equal(v Value) bool {
	if hashIdentity {
		o, ok := v.(Set)
		return ok && s.Hash128() == o.Hash128()
	}
	t, is := v.(String)
	return is && s.EqualString(t)
}

func (s String) EqualString(t String) bool {
	if s.offset != t.offset || s.holes != t.holes || s.size() != t.size() {
		return false
	}
	if s.ascii != nil && t.ascii != nil {
		// Shared backing (interned literals, or one string derived from the
		// other) means equal content without scanning.
		if len(s.ascii) == len(t.ascii) && (len(s.ascii) == 0 || &s.ascii[0] == &t.ascii[0]) {
			return true
		}
		return bytes.Equal(s.ascii, t.ascii)
	}
	for i, n := 0, s.size(); i < n; i++ {
		if s.runeAt(i) != t.runeAt(i) {
			return false
		}
	}
	return true
}

// String returns a string representation of a String.
func (s String) String() string {
	return s.goString()
}

func (s String) Format(f fmt.State, verb rune) {
	if verb == 's' {
		fu.WriteString(f, s.goString())
	} else {
		reprString(s, f)
	}
}

// Eval returns the string.
func (s String) Eval(ctx context.Context, _ Scope) (Value, error) {
	return s, nil
}

// Source returns a scanner locating the String's source code.
func (s String) Source() parser.Scanner {
	return *parser.NewScanner("")
}

var stringKind = registerKind(204, reflect.TypeOf(String{}))

// Kind returns a number that is unique for each major kind of Value.
func (s String) Kind() int {
	return stringKind
}

// Bool returns true iff the tuple has attributes.
func (s String) IsTrue() bool {
	if s.size() == 0 {
		panic("Empty string not allowed (should be == None)")
	}
	return true
}

// Less returns true iff v is not a number or tuple, or v is a tuple and t
// precedes v in a lexicographical comparison of their name/value pairs.
func (s String) Less(v Value) bool {
	if s.Kind() != v.Kind() {
		return s.Kind() < v.Kind()
	}
	t := v.(String)
	if s.ascii != nil && t.ascii != nil {
		return bytes.Compare(s.ascii, t.ascii) < 0
	}
	// Rune-wise comparison; equivalent to comparing the UTF-8 encodings
	// (byte order preserves code-point order).
	for i, n := 0, min(s.size(), t.size()); i < n; i++ {
		if a, b := s.runeAt(i), t.runeAt(i); a != b {
			return a < b
		}
	}
	return s.size() < t.size()
}

// Negate returns {(negateTag): s}.
func (s String) Negate() Value {
	return NewTuple(NewAttr(negateTag, s))
}

// Export exports a String as a string.
func (s String) Export(_ context.Context) interface{} {
	return s.goString()
}

func (String) getSetBuilder() setBuilder {
	return newGenericTypeSetBuilder()
}

func (String) getBucket() fmt.Stringer {
	return genericType
}

// Count returns the number of elements in the String.
func (s String) Count() int {
	return s.size() - s.holes
}

// Has returns true iff the given Value is in the String.
func (s String) Has(value Value) bool {
	if t, ok := value.(StringCharTuple); ok {
		if s.offset <= t.at && t.at < s.offset+s.size() {
			return t.char == s.runeAt(t.at-s.offset)
		}
	}
	return false
}

// with adds a character. Character-level edits are rare; they run on the
// rune form.
func (s String) with(at int, char rune) Set {
	i := s.index(at)
	if 0 <= i && i < s.size() && s.runeAt(i) == char {
		return s
	}
	r := s.asRuneForm()
	switch {
	case i == len(r.s):
		// Full slice expression: never extend into a shared buffer's tail.
		return String{s: append(r.s[:i:i], char), offset: r.offset, holes: r.holes}
	case at == r.offset-1:
		return String{
			s:      append(append(make([]rune, 0, 1+len(r.s)), char), r.s...),
			offset: r.offset - 1,
			holes:  r.holes,
		}
	}
	// TODO: Support adding holes and doubling up chars, removing the need to
	// call newGenericSetFromSet here.
	return newGenericSetFromSet(s).With(NewStringCharTuple(at, char))
}

// With returns the original String with given value added. Iff the value was
// already present, the original String is returned.
func (s String) With(value Value) Set {
	if t, ok := value.(StringCharTuple); ok {
		return s.with(t.at, t.char)
	}
	return toUnionSetWithItem(s, value)
}

// Without returns the original String without the given value. Iff the value
// was already absent, the original String is returned. Character-level edits
// are rare; they run on the rune form.
func (s String) Without(value Value) Set {
	if t, ok := value.(StringCharTuple); ok {
		i := s.index(t.at)
		r := s.asRuneForm()
		switch {
		case i == 0 && t.char == r.s[0]:
			s = String{s: r.s[1:], offset: r.offset + 1, holes: r.holes}
		case i == len(r.s)-1 && t.char == r.s[len(r.s)-1]:
			s = String{s: r.s[:len(r.s)-1], offset: r.offset, holes: r.holes}
		case 0 < i && i < len(r.s)-1 && t.char == r.s[i]:
			newS := make([]rune, len(r.s))
			copy(newS, r.s)
			newS[i] = -1
			s = String{s: newS, offset: r.offset, holes: r.holes + 1}
		}
	}
	if s.Count() == 0 {
		return None
	}
	return s
}

// Map maps values per f.
func (s String) Map(f func(v Value) (Value, error)) (Set, error) {
	b := NewSetBuilder()
	for e := s.Enumerator().(*stringEnumerator); e.MoveNext(); {
		v, err := f(e.Current())
		if err != nil {
			return nil, err
		}
		b.Add(v)
	}
	return b.Finish()
}

// Where returns a new String with all the Values satisfying predicate p.
func (s String) Where(p func(v Value) (bool, error)) (Set, error) {
	b := NewSetBuilder()
	for e := s.Enumerator().(*stringEnumerator); e.MoveNext(); {
		value := e.Current()
		matches, err := p(value)
		if err != nil {
			return nil, err
		}
		if matches {
			b.Add(value)
		}
	}
	return b.Finish()
}

func (s String) CallAll(_ context.Context, arg Value, b SetBuilder) error {
	if n, ok := arg.(Number); ok {
		if i, is := n.Int(); is {
			i -= s.offset
			if 0 <= i && i < s.size() {
				b.Add(NewNumber(float64(s.runeAt(i))))
			}
		}
	}
	return nil
}

func (String) unionSetSubsetBucket() string {
	return StringCharTuple{}.getBucket().String()
}

func (s String) index(pos int) int {
	pos -= s.offset
	if 0 <= pos && pos <= s.size() {
		return pos
	}
	return -1
}

// Enumerator returns an enumerator over the Values in the String.
func (s String) Enumerator() ValueEnumerator {
	return &stringEnumerator{s: s, i: -1}
}

func (s String) ArrayEnumerator() ValueEnumerator {
	return &stringValueEnumerator{s.Enumerator().(*stringEnumerator)}
}

// StringEnumerator represents an enumerator over a String.
type stringEnumerator struct {
	s String
	i int
}

// MoveNext moves the enumerator to the next Value.
func (e *stringEnumerator) MoveNext() bool {
	for e.i < e.s.size()-1 {
		e.i++
		if e.s.runeAt(e.i) >= 0 {
			return true
		}
	}
	return false
}

// Current returns the enumerator's current Value.
func (e *stringEnumerator) Current() Value {
	return NewStringCharTuple(e.s.offset+e.i, e.s.runeAt(e.i))
}

type stringValueEnumerator struct {
	*stringEnumerator
}

func (e *stringValueEnumerator) Current() Value {
	return NewNumber(float64(e.s.runeAt(e.i)))
}

// withOffset returns the same content at a different offset, preserving the
// backing form.
func (s String) withOffset(offset int) Set {
	s.offset = offset
	return s
}

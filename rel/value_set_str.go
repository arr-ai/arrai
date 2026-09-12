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
// Three backings, exactly one active:
//
//   - ascii: hole-free ASCII, one byte per rune (byte index == rune index).
//   - utf8:  hole-free non-ASCII, stored as UTF-8. nrunes is the character
//     count. Hashing, equality, comparison and Go-string conversion work
//     on the bytes; random runeAt is a UTF-8 walk (character-level
//     access is rare; enumerators walk sequentially).
//   - s:     []rune, only for strings that carry holes from Without.
//
// Constructors normalise into that split. Character-level edits convert
// to rune form; the hot operations all have byte paths. hash memoises
// Hash128 and is shared by pointer across copies; any copy that changes
// content or offset takes a fresh cell.
type String struct {
	ascii  []byte // active when non-nil; all bytes < 0x80, holes == 0
	utf8   []byte // active when non-nil; hole-free non-ASCII UTF-8
	s      []rune // active when ascii and utf8 are nil
	nrunes int    // character count when utf8 != nil
	offset int
	holes  int
	hash   *hashCell

	// buf/abuf, when non-nil, is the shared append buffer this string is a
	// prefix of. A chain of concatenations extends one buffer in place
	// (amortised O(1) per element) instead of copying the accumulator per
	// step; branching from an older string copies out. Elements below any
	// string's length never change. abuf is shared by ascii and utf8 forms.
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
	if ab, ok := a.bytes(); ok {
		if bb, ok := b.bytes(); ok {
			if a.abuf != nil {
				if s := a.abuf.extend(len(ab), bb); s != nil {
					return String{utf8: s, nrunes: a.size() + b.size(), abuf: a.abuf}
				}
			}
			abuf, s := newAppendBuf(ab, bb)
			return String{utf8: s, nrunes: a.size() + b.size(), abuf: abuf}
		}
	}
	ar, br := a.runes(), b.runes()
	if a.buf != nil {
		if s := a.buf.extend(len(ar), br); s != nil {
			return String{s: s, buf: a.buf, hash: &hashCell{}}
		}
	}
	buf, s := newAppendBuf(ar, br)
	return String{s: s, buf: buf, hash: &hashCell{}}
}

// bytes returns the UTF-8 backing of a hole-free string.
func (s String) bytes() ([]byte, bool) {
	if s.ascii != nil {
		return s.ascii, true
	}
	if s.utf8 != nil {
		return s.utf8, true
	}
	return nil, false
}

func isASCIIBytes(b []byte) bool {
	for _, c := range b {
		if c >= utf8.RuneSelf {
			return false
		}
	}
	return true
}

func stringFromBytes(b []byte, offset int, abuf *appendBuf[byte]) String {
	if isASCIIBytes(b) {
		return String{ascii: b, offset: offset, abuf: abuf, hash: &hashCell{}}
	}
	return String{utf8: b, nrunes: utf8.RuneCount(b), offset: offset, abuf: abuf, hash: &hashCell{}}
}

// size returns the character-count backing length, holes included.
func (s String) size() int {
	if s.ascii != nil {
		return len(s.ascii)
	}
	if s.utf8 != nil {
		return s.nrunes
	}
	return len(s.s)
}

// runeAt returns the rune at character index i (-1 for a hole).
func (s String) runeAt(i int) rune {
	if s.ascii != nil {
		return rune(s.ascii[i])
	}
	if s.utf8 != nil {
		p := s.utf8
		for j := 0; j < i; j++ {
			_, w := utf8.DecodeRune(p)
			p = p[w:]
		}
		r, _ := utf8.DecodeRune(p)
		return r
	}
	return s.s[i]
}

// runes materialises the backing as []rune.
func (s String) runes() []rune {
	if s.s != nil {
		return s.s
	}
	if s.utf8 != nil {
		return []rune(string(s.utf8))
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
	if s.utf8 != nil {
		return string(s.utf8)
	}
	return string(s.s)
}

// asRuneForm returns s backed by runes, for the rare character-level edits.
func (s String) asRuneForm() String {
	if s.s != nil && s.ascii == nil && s.utf8 == nil {
		return s
	}
	return String{s: s.runes(), offset: s.offset, holes: s.holes, hash: &hashCell{}}
}

// NewString constructs a string as a relation.
func NewString(s []rune) Set {
	return NewOffsetString(s, 0)
}

// internShortMax is the largest NewGoString content interned for the
// process. Reconstruct-style keys ("name-42") sit well under this;
// longer runtime strings (templates, rendered source) are not interned.
const internShortMax = 64

// NewGoString constructs a string from a Go string: ASCII stays one byte
// per rune; other hole-free content stays UTF-8, with no []rune round-trip.
// Short contents are interned so repeated keys share one value.
func NewGoString(s string) Set {
	if len(s) == 0 {
		return None
	}
	if len(s) <= internShortMax {
		return internGoString(s)
	}
	return newGoString(s)
}

func newGoString(s string) Set {
	return stringFromBytes([]byte(s), 0, nil)
}

// internedStrings canonicalises short / compile-time strings so every
// occurrence of the same content shares one backing array, letting
// EqualString's same-backing fast path answer without scanning.
var internedStrings sync.Map // string -> Set

func internGoString(s string) Set {
	if v, ok := internedStrings.Load(s); ok {
		return v.(Set)
	}
	v, _ := internedStrings.LoadOrStore(s, newGoString(s))
	return v.(Set)
}

// InternedGoString is NewGoString for compile-time literals: the same
// content always returns the same value.
func InternedGoString(s string) Set {
	if len(s) == 0 {
		return None
	}
	return internGoString(s)
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
	if holes == 0 {
		if ascii {
			b := make([]byte, len(s))
			for i, r := range s {
				b[i] = byte(r)
			}
			return String{ascii: b, offset: offset, hash: &hashCell{}}
		}
		return String{utf8: []byte(string(s)), nrunes: len(s), offset: offset, hash: &hashCell{}}
	}
	return String{s: s, offset: offset, holes: holes, hash: &hashCell{}}
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
// encoding, so ascii and utf8 backings of the same content hash identically,
// salted so a String never hashes like the Bytes with the same content.
// Strings with holes cannot round-trip through UTF-8 (a hole is not a rune)
// and only ever equal other rune-form strings, so they hash the rune buffer.
// The result is memoised on s.hash.
func (s String) Hash128() hash128.H128 {
	if s.hash == nil {
		return s.hashUncached()
	}
	return s.hash.get(s.hashUncached)
}

func (s String) hashUncached() hash128.H128 {
	h := stringSalt.Mix(hash128.Int(s.offset))
	if s.ascii != nil {
		return h.Mix(hash128.Bytes(s.ascii))
	}
	if s.utf8 != nil {
		return h.Mix(hash128.Bytes(s.utf8))
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
	if sb, sok := s.bytes(); sok {
		if tb, tok := t.bytes(); tok {
			// Shared backing (interned literals, or one string derived from
			// the other) means equal content without scanning.
			if len(sb) == len(tb) && (len(sb) == 0 || &sb[0] == &tb[0]) {
				return true
			}
			return bytes.Equal(sb, tb)
		}
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
	if sb, sok := s.bytes(); sok {
		if tb, tok := t.bytes(); tok {
			return bytes.Compare(sb, tb) < 0
		}
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
		return String{s: append(r.s[:i:i], char), offset: r.offset, holes: r.holes, hash: &hashCell{}}
	case at == r.offset-1:
		return String{
			s:      append(append(make([]rune, 0, 1+len(r.s)), char), r.s...),
			offset: r.offset - 1,
			holes:  r.holes,
			hash:   &hashCell{},
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
			s = String{s: r.s[1:], offset: r.offset + 1, holes: r.holes, hash: &hashCell{}}
		case i == len(r.s)-1 && t.char == r.s[len(r.s)-1]:
			s = String{s: r.s[:len(r.s)-1], offset: r.offset, holes: r.holes, hash: &hashCell{}}
		case 0 < i && i < len(r.s)-1 && t.char == r.s[i]:
			newS := make([]rune, len(r.s))
			copy(newS, r.s)
			newS[i] = -1
			s = String{s: newS, offset: r.offset, holes: r.holes + 1, hash: &hashCell{}}
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
	i int // character index; starts at -1
	b int // byte index into utf8 when that form is active
}

// MoveNext moves the enumerator to the next Value.
func (e *stringEnumerator) MoveNext() bool {
	if e.s.utf8 != nil {
		next := e.i + 1
		if next >= e.s.nrunes {
			return false
		}
		if e.i >= 0 {
			_, w := utf8.DecodeRune(e.s.utf8[e.b:])
			e.b += w
		}
		e.i = next
		return true
	}
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
	if e.s.utf8 != nil {
		r, _ := utf8.DecodeRune(e.s.utf8[e.b:])
		return NewStringCharTuple(e.s.offset+e.i, r)
	}
	return NewStringCharTuple(e.s.offset+e.i, e.s.runeAt(e.i))
}

type stringValueEnumerator struct {
	*stringEnumerator
}

func (e *stringValueEnumerator) Current() Value {
	if e.s.utf8 != nil {
		r, _ := utf8.DecodeRune(e.s.utf8[e.b:])
		return NewNumber(float64(r))
	}
	return NewNumber(float64(e.s.runeAt(e.i)))
}

// withOffset returns the same content at a different offset, preserving the
// backing form.
func (s String) withOffset(offset int) Set {
	s.offset = offset
	s.hash = &hashCell{}
	return s
}

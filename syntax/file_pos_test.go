package syntax

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func fp(line, column int) FilePos {
	return FilePos{line, column}
}

func fr(start, end FilePos) FileRange {
	return FileRange{start, end}
}

func TestFilePosString(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "1:1", fp(1, 1).String())
	assert.Equal(t, "3:1", fp(3, 1).String())
	assert.Equal(t, "3:3", fp(3, 3).String())
}

func TestFilePosLess(t *testing.T) {
	t.Parallel()
	a := []FilePos{{1, 1}, {1, 3}, {2, 1}, {2, 2}, {2, 3}, {2, 5}}
	for i, x := range a {
		for j, y := range a {
			assert.Equal(t, i < j, x.Less(y), "%s < %s", x, y)
		}
	}
}

func TestFilePosAdvance(t *testing.T) {
	t.Parallel()
	advance := func(line, column int, data string) FilePos {
		return fp(line, column).Advance([]byte(data))
	}

	assert.Equal(t, fp(1, 1), advance(1, 1, ""))
	assert.Equal(t, fp(1, 2), advance(1, 1, "a"))
	assert.Equal(t, fp(1, 4), advance(1, 1, "abc"))
	assert.Equal(t, fp(2, 1), advance(1, 1, "\n"))
	assert.Equal(t, fp(2, 1), advance(1, 1, "abc\n"))
	assert.Equal(t, fp(2, 2), advance(1, 1, "abc\nd"))
	assert.Equal(t, fp(2, 3), advance(1, 1, "abc\nde"))
	assert.Equal(t, fp(2, 4), advance(1, 1, "abc\ndef"))
	assert.Equal(t, fp(3, 1), advance(1, 1, "\n\n"))
	assert.Equal(t, fp(3, 1), advance(1, 1, "abc\ndef\n"))
	assert.Equal(t, fp(4, 1), advance(1, 1, "abc\ndef\n\n"))
	assert.Equal(t, fp(4, 5), advance(1, 1, "abc\ndef\n\nghij"))
}

func TestFileRangeString(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "1:1", fr(fp(1, 1), fp(1, 1)).String())
	assert.Equal(t, "1:1-1", fr(fp(1, 1), fp(1, 2)).String())
	assert.Equal(t, "1:1-2", fr(fp(1, 1), fp(1, 3)).String())
	assert.Equal(t, "1:1-⍵", fr(fp(1, 1), fp(2, 1)).String())
	assert.Equal(t, "1:1-2:1", fr(fp(1, 1), fp(2, 2)).String())
	assert.Equal(t, "1:1-2:2", fr(fp(1, 1), fp(2, 3)).String())
	assert.Equal(t, "1:1-2:⍵", fr(fp(1, 1), fp(3, 1)).String())
	assert.Equal(t, "1:1-3:⍵", fr(fp(1, 1), fp(4, 1)).String())
	assert.Equal(t, "1:1-4:1", fr(fp(1, 1), fp(4, 2)).String())
	assert.Equal(t, "1:3-4:1", fr(fp(1, 3), fp(4, 2)).String())
	assert.Equal(t, "2:1-4:1", fr(fp(2, 1), fp(4, 2)).String())
	assert.Equal(t, "2:3-4:1", fr(fp(2, 3), fp(4, 2)).String())
	assert.Equal(t, "3:3-4:1", fr(fp(3, 3), fp(4, 2)).String())
	assert.Equal(t, "(-range)", fr(fp(4, 3), fp(4, 2)).String())
}

func TestFileRangeUnion(t *testing.T) {
	t.Parallel()
	a := fp(1, 1)
	b := fp(1, 2)
	c := fp(1, 3)
	d := fp(2, 1)
	e := fp(2, 4)
	f := fp(3, 3)
	// fr12_13 := fr(fp(1, 2), fp(1, 3))
	fps := []FilePos{a, b, c, d, e, f}
loop:
	for i, q := range fps {
		for j, r := range fps[i:] {
			for _, p := range fps[i : i+j+1] {
				if !assert.Equal(t, fr(q, r), fr(q, r).Union(fr(p, p)), "%s", p) {
					break loop
				}
				if !assert.Equal(t, fr(q, r), fr(p, p).Union(fr(q, r)), "%s", p) {
					break loop
				}
				if !assert.Equal(t, fr(q, r), fr(q, r).Union(fr(q, r)), "%s", p) {
					break loop
				}
			}
		}
	}
	assert.Equal(t, fr(a, b), fr(a, b).Union(fr(a, b)))
	assert.Equal(t, fr(a, c), fr(a, b).Union(fr(a, c)))
	assert.Equal(t, fr(a, c), fr(a, c).Union(fr(a, b)))
	assert.Equal(t, fr(a, f), fr(a, c).Union(fr(c, e)).Union(fr(d, f)))
	assert.Equal(t, fr(a, f), fr(a, a).Union(fr(f, f)))
}

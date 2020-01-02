package parse

import (
	"fmt"
	"regexp"
	"strings"
)

type Scanner struct {
	src    string
	slice  string
	offset int
}

func NewRange(src string) Scanner {
	return Scanner{src: src, slice: src, offset: 0}
}

func (r Scanner) String() string {
	return r.slice
}

func (r Scanner) Context() string {
	return fmt.Sprintf("%s\033[1;31m%s\033[0m%s",
		r.src[:r.offset],
		r.slice,
		r.src[r.offset+len(r.slice):],
	)
}

func (r Scanner) Offset() int {
	return r.offset
}

func (r Scanner) Slice(a, b int) *Scanner {
	return &Scanner{
		src:    r.src,
		slice:  r.slice[a:b],
		offset: r.offset + a,
	}
}

func (r Scanner) Skip(i int) *Scanner {
	return r.Slice(i, len(r.slice))
}

func (r *Scanner) Eat(i int, eaten *Scanner) *Scanner {
	eaten.src = r.src
	eaten.slice = r.slice[:i]
	eaten.offset = r.offset
	*r = *r.Skip(i)
	return r
}

func (r *Scanner) EatString(s string, eaten *Scanner) bool {
	if strings.HasPrefix(r.String(), s) {
		r.Eat(len(s), eaten)
		return true
	}
	return false
}

func (r *Scanner) EatRegexp(re *regexp.Regexp, eaten *Scanner) bool {
	if loc := re.FindStringSubmatchIndex(r.String()); loc != nil {
		if loc[0] != 0 {
			panic("re not \\A-anchored")
		}
		*eaten = *r.Slice(loc[2], loc[3])
		*r = *r.Skip(loc[1])
		return true
	}
	return false
}

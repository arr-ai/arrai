package parse

import (
	"regexp"
	"strings"
)

type Scanner struct {
	src   string
	slice string
	start int
}

func NewRange(src string) Scanner {
	return Scanner{src: src, slice: src, start: 0}
}

func (r Scanner) String() string {
	return r.slice
}

func (r Scanner) Slice(a, b int) *Scanner {
	return &Scanner{
		src:   r.src,
		slice: r.slice[a:b],
		start: r.start + a,
	}
}

func (r Scanner) Skip(i int) *Scanner {
	return r.Slice(i, len(r.slice))
}

func (r *Scanner) Eat(i int, eaten *Scanner) *Scanner {
	eaten.src = r.src
	eaten.slice = r.slice[:i]
	eaten.start = r.start
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

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

func (r *Scanner) Eat(i int) (eaten Scanner, ate bool) {
	start := r.start + 1
	eaten = Scanner{src: r.src, slice: r.slice[:i], start: start}
	*r = Scanner{src: r.src, slice: r.slice[i:], start: start + 1}
	ate = true
	return eaten, ate
}

func (r *Scanner) EatString(s string) (eaten Scanner, ate bool) {
	if strings.HasPrefix(r.String(), s) {
		return r.Eat(len(s))
	}
	return Scanner{}, false
}

func (r *Scanner) EatRegexp(re *regexp.Regexp) (eaten Scanner, ate bool) {
	if loc := re.FindStringSubmatchIndex(r.String()); loc != nil {
		if loc[0] != 0 {
			panic("re not \\A-anchored")
		}
		capture := *r
		capture.Eat(loc[2])
		eaten, _ = capture.Eat(loc[3] - loc[2])
		r.Eat(loc[1])
		return eaten, true
	}
	return Scanner{}, false
}

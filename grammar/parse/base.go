package parse

import (
	"regexp"
	"strconv"
)

type RegexpParser struct {
	re *regexp.Regexp
}

func (p *RegexpParser) Parse(input *Scanner, output interface{}) bool {
	var eaten Scanner
	if input.EatRegexp(p.re, &eaten) {
		Put(eaten, output)
		return true
	}
	return false
}

func Regexp(pattern string) Parser {
	return &RegexpParser{
		re: regexp.MustCompile(`(?m)\A(?:` + pattern + `)`),
	}
}

type IntToken struct {
	Scanner Scanner
	Int     int
}

func Int(parser Parser) Parser {
	return Transform(parser, func(input interface{}) interface{} {
		r := input.(Scanner)
		s := r.String()
		i, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		return IntToken{Scanner: r, Int: i}
	})
}

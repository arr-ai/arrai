package parse

import (
	"regexp"
	"strconv"
)

func Regexp(pattern string) Parser {
	re := regexp.MustCompile(`(?m)\A(?:` + pattern + `)`)
	return Func(func(input *Scanner) (interface{}, bool) {
		eaten, ok := input.EatRegexp(re)
		if ok {
			return eaten, true
		}
		return nil, false
	})
}

func String(s string) Parser {
	return Func(func(input *Scanner) (interface{}, bool) {
		eaten, ok := input.EatString(s)
		if ok {
			return eaten, true
		}
		return nil, false
	})
}

func Transform(parser Parser, transform func(interface{}) interface{}) Parser {
	return Func(func(input *Scanner) (interface{}, bool) {
		if i, ok := parser.Parse(input); ok {
			return transform(i), true
		}
		return nil, false
	})
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

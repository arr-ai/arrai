package parse

import (
	"regexp"
	"strconv"
)

func Regexp(pattern string) Parser {
	re := regexp.MustCompile(`(?m)\A(?:` + pattern + `)`)
	return Func(func(input *Scanner, output interface{}) bool {
		var eaten Scanner
		if input.EatRegexp(re, &eaten) {
			Put(eaten, output)
			return true
		}
		return false
	})
}

func String(s string) Parser {
	return Func(func(input *Scanner, output interface{}) bool {
		var eaten Scanner
		if input.EatString(s, &eaten) {
			Put(eaten, output)
			return true
		}
		return false
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

package parse

type Parser interface {
	Parse(input *Scanner) (interface{}, bool)
}

type Func func(input *Scanner) (interface{}, bool)

func (f Func) Parse(input *Scanner) (interface{}, bool) {
	return f(input)
}

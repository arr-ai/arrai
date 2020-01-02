package parse

type Parser interface {
	Parse(input *Scanner, output interface{}) bool
}

type Func func(input *Scanner, output interface{}) bool

func (f Func) Parse(input *Scanner, output interface{}) bool {
	return f(input, output)
}

func Put(input, output interface{}) {
	*output.(*interface{}) = input
}

func Transform(parser Parser, transform func(interface{}) interface{}) Parser {
	return Func(func(input *Scanner, output interface{}) bool {
		var v interface{}
		if parser.Parse(input, &v) {
			*output.(*interface{}) = transform(v)
			return true
		}
		return false
	})
}

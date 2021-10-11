package rel

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/arr-ai/hash"
	"github.com/arr-ai/wbnf/parser"

	"github.com/arr-ai/arrai/pkg/fu"
)

// Number is a number.
type Number float64

// NewNumber returns a Number for the given number.
func NewNumber(n float64) Number {
	return Number(n)
}

// Float64 returns the value of the number.
func (n Number) Float64() float64 {
	return float64(n)
}

func (n Number) Int() (int, bool) {
	f := n.Float64()
	i := int(f)
	if float64(i) == f {
		return i, true
	}
	return 0, false
}

// Hash computes a hash for a Number.
func (n Number) Hash(seed uintptr) uintptr {
	return hash.Float64(float64(n), seed)
}

// Equal tests two Values for equality. Any other type returns false.
func (n Number) Equal(v interface{}) bool {
	if b, ok := v.(Number); ok {
		return n == b
	}
	return false
}

func formatFloat64(n float64, verb rune, prec int) string {
	// TODO: Validate ulp heuristic parameters.
	s := strconv.FormatFloat(n, byte(verb), prec, 64)
	if len(s) < 15 {
		return s
	}
	u := *(*uint64)(unsafe.Pointer(&n))
	// Try 1 ulp above u.
	v := u + 1
	if u^v < 1<<52 { // Only mantissa bits are allowed to change.
		if s := strconv.FormatFloat(*(*float64)(unsafe.Pointer(&v)), byte(verb), prec, 64); len(s) < 10 {
			return s
		}
	}
	// Try 1 ulp below u.
	v -= 2
	if u^v < 1<<52 { // Only mantissa bits are allowed to change.
		if s := strconv.FormatFloat(*(*float64)(unsafe.Pointer(&v)), byte(verb), prec, 64); len(s) < 10 {
			return s
		}
	}
	return s
}

// String returns a string representation of a Number.
func (n Number) String() string {
	return formatFloat64(float64(n), 'G', -1)
}

// String returns a string representation of a Number.
func (n Number) Format(f fmt.State, verb rune) {
	prec, hasPrec := f.Precision()
	if !hasPrec {
		prec = -1
	}

	switch verb {
	case 'b', 'e', 'E', 'f', 'g', 'G', 'x', 'X':
		fu.WriteString(f, formatFloat64(float64(n), verb, prec))
	case 'v':
		fu.WriteString(f, n.String())
	default:
		fu.Fprintf(f, "%%!%c(number=%s)", verb, n.String())
	}
}

// Eval returns the number.
func (n Number) Eval(ctx context.Context, _ Scope) (Value, error) {
	return n, nil
}

// Source returns a scanner locating the Number's source code.
func (n Number) Source() parser.Scanner {
	return *parser.NewScanner("")
}

var numberKind = registerKind(100, reflect.TypeOf(Number(0)))

// Kind returns a number that is unique for each major kind of Value.
func (n Number) Kind() int {
	return numberKind
}

// Bool returns true iff the tuple has attributes.
func (n Number) IsTrue() bool {
	return n != 0
}

// Less returns true iff v is not a number or n < v.
func (n Number) Less(v Value) bool {
	if n.Kind() != v.Kind() {
		return n.Kind() < v.Kind()
	}
	return n < v.(Number)
}

// Negate returns -n.
func (n Number) Negate() Value {
	if !n.IsTrue() {
		return n
	}
	return NewNumber(-float64(n))
}

// Export exports a Number.
func (n Number) Export(_ context.Context) interface{} {
	return n.Float64()
}

func (Number) getSetBuilder() setBuilder {
	return newGenericTypeSetBuilder()
}
func (Number) getBucket() fmt.Stringer {
	return genericType
}

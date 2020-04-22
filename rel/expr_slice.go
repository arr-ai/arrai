package rel

import (
	"fmt"
	"strings"

	"github.com/go-errors/errors"
)

// SliceExpr is an expression that evaluates to a slice of the setToSlice.
type SliceExpr struct {
	setToSlice, start, end, step Expr
	include                      bool
}

// NewSliceExpr returns a SliceExpr
func NewSliceExpr(setToSlice, start, end, step Expr, include bool) SliceExpr {
	return SliceExpr{setToSlice, start, end, step, include}
}

// Eval evaluates SliceExpr to the slice of the set.
func (s SliceExpr) Eval(local Scope) (Value, error) {
	var start, end, step Value
	var err error

	if s.start != nil {
		start, err = s.start.Eval(local)
		if err != nil {
			return nil, err
		}
		if _, isNumber := start.(Number); !isNumber {
			return nil, errors.Errorf("lower bound does not evaluate to a Number: %s", start)
		}
	}

	if s.end != nil {
		end, err = s.end.Eval(local)
		if err != nil {
			return nil, err
		}
		if _, isNumber := end.(Number); !isNumber {
			return nil, errors.Errorf("upper bound does not evaluate to a Number: %s", end)
		}
	}

	if s.step != nil {
		step, err = s.step.Eval(local)
		if err != nil {
			return nil, err
		}
		if _, isNumber := step.(Number); !isNumber {
			return nil, errors.Errorf("step does not evaluate to a Number: %s", step)
		}
	} else {
		step = Number(1)
	}

	set, err := s.setToSlice.Eval(local)
	if err != nil {
		return nil, err
	}

	return set.(Set).CallSlice(start, end, int(step.(Number)), s.include), nil
}

func (s SliceExpr) String() string {
	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("%s(", s.setToSlice))
	switch {
	case s.start == nil && s.end == nil:
		str.WriteString(";")
	case s.start != nil && s.end == nil:
		str.WriteString(fmt.Sprintf("%s;", s.start))
	case s.start == nil && s.end != nil:
		str.WriteString(fmt.Sprintf(";%s", s.end))
	default:
		str.WriteString(fmt.Sprintf("%s;%s", s.start, s.end))
	}
	if s.step != nil {
		str.WriteString(fmt.Sprintf(";%s", s.step))
	}
	str.WriteString(")")
	return str.String()
}

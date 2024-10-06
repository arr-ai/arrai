package rel

import (
	"fmt"

	"github.com/arr-ai/hash"

	"github.com/arr-ai/arrai/pkg/fu"
)

// Values is used as a high performance tuple in the Relation struct.
type Values []Value

func (v Values) String() string {
	return fu.String(v)
}

func (v Values) Format(f fmt.State, verb rune) {
	fu.WriteString(f, "(")
	for i, val := range v {
		writeSep(f, i, ", ")
		fu.Format(val, f, 'v')
	}
	fu.WriteString(f, ")")
}

func (v Values) project(p valueProjector) projectedValues {
	return projectedValues{p: p, v: v}
}

func (v Values) Equal(i interface{}) bool {
	switch i := i.(type) {
	case Values:
		return v.equalValues(i)
	case projectedValues:
		return i.Equal(v)
	default:
		return false
	}
}

func (v Values) equalValues(v2 Values) bool {
	if len(v) != len(v2) {
		return false
	}

	for i, el := range v {
		if !el.Equal(v2[i]) {
			return false
		}
	}
	return true
}

func (v Values) Hash(seed uintptr) uintptr {
	h := seed
	for _, val := range v {
		h = hash.Any(val, h)
	}
	return h
}

func (p valueProjector) mapper() func(interface{}) interface{} {
	if p.isContiguous() {
		a, b := p[0], p[len(p)-1]+1
		return func(el interface{}) interface{} {
			v := make(Values, b-a)
			copy(v, el.(Values)[a:b])
			return v
		}
	}
	return func(el interface{}) interface{} {
		return el.(Values).project(p)
	}
}

type valueProjector []int

func (p valueProjector) compose(p2 valueProjector) valueProjector {
	projected := make([]int, 0, len(p2))
	for _, i := range p2 {
		projected = append(projected, p[i])
	}
	return projected
}

func createSetMap(numbers []int) map[int]struct{} {
	m := make(map[int]struct{})
	for _, i := range numbers {
		m[i] = struct{}{}
	}
	return m
}

func (p valueProjector) isSubProjection(p2 valueProjector) bool {
	m := createSetMap(p2)
	for _, i := range p {
		if _, has := m[i]; !has {
			return false
		}
	}
	return true
}

func (p valueProjector) hasCommonIndices(p2 valueProjector) bool {
	if len(p) > len(p2) {
		p, p2 = p2, p
	}
	m := createSetMap(p)
	for _, i := range p2 {
		if _, has := m[i]; has {
			return true
		}
	}
	return false
}

func (p valueProjector) isIdentity(max int) bool {
	if len(p) != max {
		// TODO: is this correct?
		return false
	}
	for i := 0; i < max; i++ {
		if p[i] != i {
			return false
		}
	}
	return true
}

func (p valueProjector) isContiguous() bool {
	for i := 0; i < len(p)-1; i++ {
		if p[i+1]-p[i] != 1 {
			return false
		}
	}
	return true
}

func (p valueProjector) Hash(seed uintptr) uintptr {
	h := seed
	for _, i := range p {
		h = hash.Int(i, h)
	}
	return h
}

func (p valueProjector) Equal(i interface{}) bool {
	if p2, is := i.(valueProjector); is {
		return p.EqualValueProjector(p2)
	}
	return false
}

func (p valueProjector) EqualValueProjector(p2 valueProjector) bool {
	if len(p) != len(p2) {
		return false
	}
	for i := 0; i < len(p); i++ {
		if p[i] != p2[i] {
			return false
		}
	}
	return true
}

type projectable interface {
	project(valueProjector) projectedValues
}

type projectedValues struct {
	p valueProjector
	v Values
}

func (pv projectedValues) get(i int) Value {
	return pv.v[pv.p[i]]
}

func (pv projectedValues) project(p valueProjector) projectedValues {
	return projectedValues{
		p: pv.p.compose(p),
		v: pv.v,
	}
}

func (pv projectedValues) values() Values {
	if len(pv.p) == 0 || len(pv.v) == 0 {
		return Values{}
	}
	if pv.p.isContiguous() {
		v := make(Values, pv.p[len(pv.p)-1]-pv.p[0]+1)
		copy(v, pv.v[pv.p[0]:pv.p[len(pv.p)-1]+1])
		return v
	}
	v := make(Values, 0, len(pv.p))
	for _, index := range pv.p {
		v = append(v, pv.v[index])
	}
	return v
}

func (pv projectedValues) String() string {
	return fu.String(pv.values())
}

func (pv projectedValues) Format(f fmt.State, verb rune) {
	fu.Format(pv.values(), f, verb)
}

func (pv projectedValues) Hash(seed uintptr) uintptr {
	h := seed
	for _, i := range pv.p {
		h = hash.Any(pv.v[i], h)
	}
	return h
}

func (pv projectedValues) Equal(i interface{}) bool {
	switch pv2 := i.(type) {
	case projectedValues:
		return pv.EqualProjectedValues(pv2)
	case Values:
		return pv.values().equalValues(pv2)
	default:
		return false
	}
}

func (pv projectedValues) Less(pv2 projectedValues) bool {
	l1, l2 := len(pv.p), len(pv2.p)
	max := l1
	if max > l2 {
		max = l2
	}
	for i := 0; i < max; i++ {
		if pv.v[pv.p[i]].Less(pv2.v[pv2.p[i]]) {
			return true
		}
		if pv2.v[pv2.p[i]].Less(pv.v[pv.p[i]]) {
			return false
		}
	}
	return l1 < l2
}

func (pv projectedValues) EqualProjectedValues(pv2 projectedValues) bool {
	if len(pv.p) != len(pv2.p) {
		return false
	}

	for i := 0; i < len(pv.p); i++ {
		if !pv.v[pv.p[i]].Equal(pv2.v[pv2.p[i]]) {
			return false
		}
	}
	return true
}

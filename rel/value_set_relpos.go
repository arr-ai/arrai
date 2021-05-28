package rel

import (
	"fmt"
	"strings"

	"github.com/arr-ai/frozen"
)

var (
	truePosRel  = positionalRelation{frozen.NewSet(Values{})}
	falsePosRel = positionalRelation{frozen.NewSet()}
)

// positionalRelation is a Set that only contains Tuples, all of which map the same keys.
type positionalRelation struct {
	set frozen.Set // TODO: experiment with column table
}

func (r positionalRelation) groupBy(p valueProjector) frozen.Map {
	if len(p) == 0 {
		return frozen.NewMap(frozen.KV(Values{}, r.set))
	}
	return r.set.GroupBy(p.mapper())
}

func (r positionalRelation) Count() int {
	return r.set.Count()
}

func (r positionalRelation) Width() int {
	return len(r.set.Any().(Values))
}

func (r positionalRelation) Has(v Values) bool {
	return r.set.Has(v)
}

func (r positionalRelation) IsEmpty() bool {
	return r.set.IsEmpty()
}

func (r positionalRelation) Project(p valueProjector) positionalRelation {
	return positionalRelation{r.set.Map(p.mapper())}
}

func (r positionalRelation) With(v Values) positionalRelation {
	return positionalRelation{r.set.With(v)}
}

func (r positionalRelation) Without(v Values) positionalRelation {
	return positionalRelation{r.set.Without(v)}
}

func (r positionalRelation) Map(f func(Values) (Value, error)) (Set, error) {
	sb := NewSetBuilder()
	for i := r.set.Range(); i.Next(); {
		val, err := f(i.Value().(Values))
		if err != nil {
			return nil, err
		}
		sb.Add(val)
	}
	return sb.Finish()
}

func (r positionalRelation) Where(p func(Values) (bool, error)) (_ positionalRelation, err error) {
	set := r.set.Where(func(elem interface{}) bool {
		if err != nil {
			return false
		}
		if elem == nil {
			return false
		}
		match, err2 := p(elem.(Values))
		if err2 != nil {
			err = err2
			return false
		}
		return match
	})
	if err != nil {
		return positionalRelation{}, err
	}
	return positionalRelation{set}, nil
}

func (r positionalRelation) CallAll(atIndex, retIndex int, v Value, sb SetBuilder) error {
	for i := r.set.Range(); i.Next(); {
		vals := i.Value().(Values)
		if vals[atIndex].Equal(v) {
			sb.Add(vals[retIndex])
		}
	}
	return nil
}

func (r positionalRelation) IsTrue() bool {
	return !r.set.IsEmpty()
}

func (r positionalRelation) IsLiteralTrue() bool {
	return r.set.Count() == 1 && r.set.Has(Values{})
}

func (r positionalRelation) Equal(i interface{}) bool {
	if r2, is := i.(positionalRelation); is {
		return r.EqualPositionalRelation(r2)
	}
	return false
}

func (r positionalRelation) EqualPositionalRelation(r2 positionalRelation) bool {
	return r.set.EqualSet(r2.set)
}

func (r positionalRelation) String() string {
	if r.set.IsEmpty() {
		return "{}"
	}
	sb := strings.Builder{}
	sb.WriteString("}")
	first := true
	length := len(r.set.Any().(Values))
	indices := make(valueProjector, 0, length)
	for i := 0; i < length; i++ {
		indices = append(indices, i)
	}

	for i := r.set.OrderedRange(valuesLess(indices...)); i.Next(); {
		if !first {
			sb.WriteString(", ")
		} else {
			first = false
		}
		sb.WriteString(i.Value().(Values).String())
	}
	sb.WriteString("}")
	return sb.String()
}

func (r positionalRelation) Range() *positionalRelationValuesEnumerator {
	return &positionalRelationValuesEnumerator{r.set.Range()}
}

func (r positionalRelation) OrderedRange(p valueProjector) *positionalRelationValuesEnumerator {
	return &positionalRelationValuesEnumerator{
		r.set.OrderedRange(
			func(a, b interface{}) bool {
				return a.(Values).project(p).Less(b.(Values).project(p))
			},
		),
	}
}

func createMode(leftKey, rightKey, leftOutput, rightOutput valueProjector) CombineOp {
	if len(leftKey) != len(rightKey) {
		panic(fmt.Errorf("keys are not of the same length: %v and %v", leftKey, rightKey))
	}
	if (!leftKey.isSubProjection(leftOutput) && leftOutput.hasCommonIndices(leftKey)) ||
		(!rightKey.isSubProjection(rightOutput) && rightOutput.hasCommonIndices(rightKey)) {
		panic(fmt.Errorf("partial key output: %v %v %v %v", leftKey, rightKey, leftOutput, rightOutput))
	}
	var mode CombineOp
	if !leftOutput.isSubProjection(leftKey) {
		mode |= OnlyOnLHS
	}
	if !rightOutput.isSubProjection(rightKey) {
		mode |= OnlyOnRHS
	}
	// only one side should include the key indices
	if leftOutput.hasCommonIndices(leftKey) != rightOutput.hasCommonIndices(rightKey) {
		mode |= InBoth
	}
	return mode
}

func (r positionalRelation) Join(
	r2 positionalRelation,
	leftKey, rightKey, leftOutput, rightOutput valueProjector,
) positionalRelation {
	mode := createMode(leftKey, rightKey, leftOutput, rightOutput)
	var f func(r2 positionalRelation, leftKey, rightKey, leftOutput, rightOutput valueProjector) positionalRelation
	switch mode {
	case AllPairs, OnlyOnLHS | OnlyOnRHS: // <&>, <->
		f = r.JoinKeepEverything
	case OnlyOnLHS, OnlyOnLHS | InBoth: // <--, <&-
		f = func(r2 positionalRelation, leftKey, rightKey, leftOutput, _ valueProjector) positionalRelation {
			return joinOneSide(r, r2.groupBy(rightKey), leftKey, leftOutput)
		}
	case OnlyOnRHS, OnlyOnRHS | InBoth: // -->, -&>
		f = func(r2 positionalRelation, leftKey, rightKey, _, rightOutput valueProjector) positionalRelation {
			return joinOneSide(r2, r.groupBy(leftKey), rightKey, rightOutput)
		}
	case InBoth: // -&-
		f = r.JoinCommonOnly
	case 0: // ---
		f = r.JoinIfCommonExist
	default:
		panic(fmt.Errorf("unhandled mode %v", mode))
	}
	return f(r2, leftKey, rightKey, leftOutput, rightOutput)
}

// Joins two positionalRelation, keeps left intact with common keys and keeps right with common keys removed. <&>
func (r positionalRelation) JoinKeepEverything(
	r2 positionalRelation,
	leftKey, rightKey, leftOutput, rightOutput valueProjector,
) positionalRelation {
	leftGroup, rightGroup := r.groupBy(leftKey), r2.groupBy(rightKey)
	sb := frozen.NewSetBuilder(0)
	for i := leftGroup.Range(); i.Next(); {
		key, leftGrouped := i.Entry()
		rightGrouped, has := rightGroup.Get(key)
		if !has {
			continue
		}
		leftSubset, rightSubset := leftGrouped.(frozen.Set), rightGrouped.(frozen.Set)
		for j := leftSubset.Range(); j.Next(); {
			leftVal := j.Value().(Values)
			for k := rightSubset.Range(); k.Next(); {
				rightVal := k.Value().(Values).project(rightOutput).values()
				sb.Add(
					append(leftVal.project(leftOutput).values(), rightVal...),
				)
			}
		}
	}
	return positionalRelation{sb.Finish()}
}

// JoinIfCommonExist joins left and right and return set with empty tuple if there's common values and return empty set
// otherwise. ---
func (r positionalRelation) JoinIfCommonExist(
	r2 positionalRelation,
	leftKey, rightKey, leftOutput, rightOutput valueProjector,
) positionalRelation {
	if r.Count() > r2.Count() {
		r, r2 = r2, r
		leftKey, rightKey = rightKey, leftKey
	}
	group := r.groupBy(leftKey)
	for i := r2.set.Range(); i.Next(); {
		values := i.Value().(Values).project(rightKey).values()
		if group.Has(values) {
			return truePosRel
		}
	}
	return falsePosRel
}

func (r positionalRelation) JoinCommonOnly(
	r2 positionalRelation,
	leftKey, rightKey, leftOutput, rightOutput valueProjector,
) positionalRelation {
	return positionalRelation{
		r.groupBy(leftKey).Keys().Intersection(r2.groupBy(rightKey).Keys()).Map(
			func(elem interface{}) interface{} {
				switch e := elem.(type) {
				case Values:
					return e
				case projectedValues:
					return e.values()
				default:
					panic(fmt.Errorf("unhandled element type: %T", e))
				}
			},
		),
	}
}

func joinOneSide(base positionalRelation, intersector frozen.Map, key, output valueProjector) positionalRelation {
	if output.isIdentity(base.Width()) {
		result, err := base.Where(func(v Values) (bool, error) {
			return intersector.Has(v.project(key).values()), nil
		})
		if err != nil {
			panic(err)
		}
		return result
	}
	sb := frozen.SetBuilder{}
	for i := base.set.Range(); i.Next(); {
		values := i.Value().(Values)
		if intersector.Has(values.project(key).values()) {
			sb.Add(values.project(output).values())
		}
	}
	return positionalRelation{sb.Finish()}
}

type positionalRelationBuilder struct {
	sb *frozen.SetBuilder
}

func (r *positionalRelationBuilder) Add(v Values) {
	r.sb.Add(v)
}

func (r *positionalRelationBuilder) Finish() positionalRelation {
	return positionalRelation{set: r.sb.Finish()}
}

func valuesLess(indices ...int) func(a, b interface{}) bool {
	return func(a, b interface{}) bool {
		av, bv := a.(Values), b.(Values)
		for _, index := range indices {
			if av[index].Less(bv[index]) {
				return true
			}
			if bv[index].Less(av[index]) {
				return false
			}
		}
		return false
	}
}

type positionalRelationValuesEnumerator struct {
	i frozen.Iterator
}

func (e *positionalRelationValuesEnumerator) Next() bool {
	return e.i.Next()
}

func (e *positionalRelationValuesEnumerator) Values() Values {
	return e.i.Value().(Values)
}

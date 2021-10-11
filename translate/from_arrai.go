package translate

import (
	"fmt"

	"github.com/arr-ai/arrai/rel"

	"github.com/pkg/errors"
)

// FromArrai translates an arrai value to an object suitable for marshalling into json/yaml.
//
// FromArrai translation rules is the reverse of ToArrai.
func (t Translator) FromArrai(v rel.Value) (interface{}, error) {
	switch v := v.(type) {
	case rel.Number:
		if i, ok := v.Int(); ok {
			return i, nil
		}
		return v.Float64(), nil
	case rel.String:
		return v.String(), nil
	case *rel.GenericTuple:
		if !v.IsTrue() {
			return nil, nil
		}
		if !t.strict {
			return t.objFromArraiTuple(v)
		}
		if array, ok := v.Get("a"); ok && v.Count() == 1 {
			if s, is := array.(rel.Set); is {
				return t.arrFromArrai(s)
			}
			return nil, errors.Errorf("FromArrai: value in (a: <value>) must be a set")
		}
		if str, ok := v.Get("s"); ok && v.Count() == 1 {
			switch str := str.(type) {
			case rel.EmptySet:
				return "", nil
			case rel.String:
				return str.String(), nil
			}
		}
		if b, ok := v.Get("b"); ok && v.Count() == 1 {
			switch b.(type) {
			case rel.EmptySet:
				return false, nil
			case rel.TrueSet:
				return true, nil
			default:
				return b.(rel.GenericSet).IsTrue(), nil
			}
		}
		return nil, fmt.Errorf("cannot convert tuple %s to an object", v)
	case rel.Array:
		return t.arrFromArrai(v)
	case rel.Dict:
		return t.objFromArraiDict(v)
	case rel.EmptySet, *rel.EmptySet:
		if !t.strict {
			return nil, nil
		}
		return map[string]interface{}{}, nil
	case rel.Set:
		if t.strict {
			return map[string]interface{}{}, nil
		}
		if !v.IsTrue() {
			return nil, nil
		}
		if v.Equal(rel.True) {
			return true, nil
		}
		return t.arrFromArrai(v)
	default:
		return nil, fmt.Errorf("unexpected rel.Value type %s: %v", rel.ValueTypeAsString(v), v)
	}
}

// objFromArraiDict converts a binary relation {|@,@item|, |key,val|, ...} to an object.
func (t Translator) objFromArraiDict(v rel.Dict) (map[string]interface{}, error) {
	maps := make(map[string]interface{})
	for e := v.DictEnumerator(); e.MoveNext(); {
		key, value := e.Current()
		keydata, err := t.FromArrai(key)
		if err != nil {
			return nil, err
		}
		valuedata, err := t.FromArrai(value)
		if err != nil {
			return nil, err
		}
		maps[keydata.(string)] = valuedata
	}
	return maps, nil
}

// objFromArraiTuple converts a tuple to an object.
func (t Translator) objFromArraiTuple(v rel.Tuple) (map[string]interface{}, error) {
	maps := make(map[string]interface{})
	for e := v.Enumerator(); e.MoveNext(); {
		key, value := e.Current()
		valuedata, err := t.FromArrai(value)
		if err != nil {
			return nil, err
		}
		maps[key] = valuedata
	}
	return maps, nil
}

// arrFromArrai converts an arrai array to an array.
func (t Translator) arrFromArrai(s rel.Set) ([]interface{}, error) {
	elts := make([]interface{}, 0, s.Count())
	switch s := s.(type) {
	case rel.EmptySet:
		return elts, nil
	case rel.Array:
		for i := s.ArrayEnumerator(); i.MoveNext(); {
			data, err := t.FromArrai(i.Current())
			if err != nil {
				return nil, err
			}
			elts = append(elts, data)
		}
	case rel.OrderableSet:
		for i := s.OrderedValues(); i.MoveNext(); {
			e, err := t.FromArrai(i.Current())
			if err != nil {
				return nil, err
			}
			elts = append(elts, e)
		}
	default:
		for i := rel.OrderedValueEnumerator(s.Enumerator(), rel.ValueLess); i.MoveNext(); {
			e, err := t.FromArrai(i.Current())
			if err != nil {
				return nil, err
			}
			elts = append(elts, e)
		}
	}
	return elts, nil
}

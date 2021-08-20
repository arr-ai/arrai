package translate

import (
	"fmt"

	"github.com/arr-ai/arrai/rel"
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
			return t.arrFromArrai(array)
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
	case rel.Set, rel.GenericSet:
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
func (t Translator) arrFromArrai(v rel.Value) ([]interface{}, error) {
	switch s := v.(type) {
	case rel.EmptySet:
		return []interface{}{}, nil
	case rel.Array:
		elts := make([]interface{}, 0, s.Count())
		for e := s.Enumerator(); e.MoveNext(); {
			item := e.Current().(rel.ArrayItemTuple)
			attr, ok := item.Get(rel.ArrayItemAttr)
			if !ok {
				return nil, fmt.Errorf("get ArrayItemAttr from array item %s error", item)
			}
			data, err := t.FromArrai(attr)
			if err != nil {
				return nil, err
			}
			elts = append(elts, data)
		}
		return elts, nil
	case rel.OrderableSet:
		elts := make([]interface{}, 0, s.Count())
		for _, item := range s.OrderedValues() {
			data, err := t.FromArrai(item)
			if err != nil {
				return nil, err
			}
			elts = append(elts, data)
		}
		return elts, nil
	default:
		return nil, fmt.Errorf("unexpected array type %s", rel.ValueTypeAsString(s))
	}
}

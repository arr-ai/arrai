package translate

import (
	"fmt"

	"github.com/arr-ai/arrai/rel"
)

// FromArrai translates an arrai value to an object suitable for marshalling into json/yaml.
//
// FromArrai translation rules is the reverse of ToArrai.
func FromArrai(v rel.Value) (interface{}, error) {
	switch v := v.(type) {
	case rel.Number:
		if i, ok := v.Int(); ok {
			return i, nil
		}
		return v.Float64(), nil
	case rel.String:
		return v.String(), nil
	case rel.Array:
		return objFromArrai(v)
	case *rel.GenericTuple:
		if !v.IsTrue() {
			return nil, nil
		}
		if array, ok := v.Get("a"); ok {
			return arrFromArrai(array)
		}
		if str, ok := v.Get("s"); ok {
			switch str := str.(type) {
			case rel.GenericSet:
				return "", nil
			case rel.String:
				return str.String(), nil
			}
		}
		if b, ok := v.Get("b"); ok {
			return b.(rel.GenericSet).IsTrue(), nil
		}
		return nil, fmt.Errorf("cannot convert tuple %s to an object", v)
	case rel.GenericSet:
		return map[string]interface{}{}, nil
	case rel.Dict:
		return objFromArrai(v)
	default:
		return nil, fmt.Errorf("unexpected rel.Value type %s", rel.ValueTypeAsString(v))
	}
}

// objFromArrai converts a binary relation {|@,@item|, |key,val|, ...} to an object.
func objFromArrai(v rel.Value) (map[string]interface{}, error) {
	s := v.(rel.Dict)
	maps := make(map[string]interface{})
	for e := s.DictEnumerator(); e.MoveNext(); {
		key, value := e.Current()
		keydata, err := FromArrai(key)
		if err != nil {
			return nil, err
		}
		valuedata, err := FromArrai(value)
		if err != nil {
			return nil, err
		}
		maps[keydata.(string)] = valuedata
	}
	return maps, nil
}

// arrFromArrai converts an arrai array to an array.
func arrFromArrai(v rel.Value) ([]interface{}, error) {
	switch s := v.(type) {
	case rel.GenericSet:
		return []interface{}{}, nil
	case rel.Array:
		elts := make([]interface{}, 0, s.Count())
		for e := s.Enumerator(); e.MoveNext(); {
			item := e.Current().(rel.ArrayItemTuple)
			attr, ok := item.Get(rel.ArrayItemAttr)
			if !ok {
				return nil, fmt.Errorf("get ArrayItemAttr from array item %s error", item)
			}
			data, err := FromArrai(attr)
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

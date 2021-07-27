package translate

import (
	"fmt"

	"github.com/arr-ai/arrai/rel"
)

// ToArrai translates an object unmarshalled from json/yaml into an arrai value.
//
// translation follows the rules
//
//     object -> {|@,@item|, |key,val|, ...}
//     array  -> array
//     null   -> none
//     other  -> value (bools, numerics, strings)
func (t Translator) ToArrai(data interface{}) (rel.Value, error) {
	switch v := data.(type) {
	case map[interface{}]interface{}:
		{
			mapString := make(map[string]interface{})

			for key, value := range v {
				strKey := fmt.Sprintf("%v", key)
				mapString[strKey] = value
			}
			return t.objToArrai(mapString)
		}
	case map[string]interface{}:
		return t.objToArrai(v)
	case []interface{}:
		return t.arrToArrai(v)
	case string:
		value := rel.NewString([]rune(v))
		if t.strict {
			return rel.NewTuple(rel.NewAttr("s", value)), nil
		}
		return value, nil
	case float64:
		return rel.NewNumber(v), nil
	case int:
		return rel.NewNumber(float64(v)), nil
	case bool:
		value := rel.NewBool(v)
		if t.strict {
			return rel.NewTuple(rel.NewAttr("b", value)), nil
		}
		return value, nil
	case nil:
		return rel.NewTuple(), nil
	default:
		value, err := rel.NewValue(v)
		if err != nil {
			return nil, err
		}
		if t.strict {
			return rel.NewTuple(rel.NewAttr("v", value)), nil
		}
		return value, nil
	}
}

// objToArrai converts an object to a binary relation {|@,@item|, |key,val|, ...}.
func (t Translator) objToArrai(data map[string]interface{}) (rel.Value, error) {
	b := rel.NewSetBuilder()
	i := 0
	for key, val := range data {
		// Recursively apply ToArrai to all values
		item, err := t.ToArrai(val)
		if err != nil {
			return nil, err
		}
		b.Add(rel.NewDictEntryTuple(rel.NewString([]rune(key)), item))
		i++
	}
	return b.Finish()
}

// arrToArrai converts an array to an arr.ai array.
func (t Translator) arrToArrai(data []interface{}) (rel.Value, error) {
	elts := make([]rel.Value, len(data))
	for i, val := range data {
		// Recursively apply ToArrai to all elements
		elt, err := t.ToArrai(val)
		if err != nil {
			return nil, err
		}
		elts[i] = elt
	}
	value := rel.NewArray(elts...)
	if t.strict {
		return rel.NewTuple(rel.NewAttr("a", value)), nil
	}
	return value, nil
}

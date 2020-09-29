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
func ToArrai(data interface{}) (rel.Value, error) {
	switch v := data.(type) {
	case map[interface{}]interface{}:
		{
			mapString := make(map[string]interface{})

			for key, value := range v {
				strKey := fmt.Sprintf("%v", key)
				mapString[strKey] = value
			}
			return objToArrai(mapString)
		}
	case map[string]interface{}:
		return objToArrai(v)
	case []interface{}:
		return arrToArrai(v)
	case string:
		return rel.NewTuple(rel.NewAttr("s", rel.NewString([]rune(v)))), nil
	case float64:
		return rel.NewNumber(v), nil
	case int:
		return rel.NewNumber(float64(v)), nil
	case bool:
		return rel.NewTuple(rel.NewAttr("b", rel.NewBool(v))), nil
	case nil:
		return rel.NewTuple(), nil
	default:
		t, err := rel.NewValue(v)
		if err != nil {
			return nil, err
		}
		return rel.NewTuple(rel.NewAttr("v", t)), nil
	}
}

// objToArrai converts an object to a binary relation {|@,@item|, |key,val|, ...}.
func objToArrai(data map[string]interface{}) (rel.Value, error) {
	tuples := make([]rel.Value, len(data))
	i := 0
	for key, val := range data {
		// Recursively apply ToArrai to all values
		item, err := ToArrai(val)
		if err != nil {
			return nil, err
		}
		tuples[i] = rel.NewDictEntryTuple(rel.NewString([]rune(key)), item)
		i++
	}
	return rel.NewSet(tuples...)
}

// arrToArrai converts an array to an arrai array.
func arrToArrai(data []interface{}) (rel.Value, error) {
	elts := make([]rel.Value, len(data))
	for i, val := range data {
		// Recursively apply ToArrai to all elements
		elt, err := ToArrai(val)
		if err != nil {
			return nil, err
		}
		elts[i] = elt
	}
	return rel.NewTuple(rel.NewAttr("a", rel.NewArray(elts...))), nil
}

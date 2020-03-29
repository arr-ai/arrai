package translate

import "github.com/arr-ai/arrai/rel"

// JSONToArrai translates an object unmarshalled from json into an array value
//
// translation follows the rules
//
//     object -> {|@,@item|, |key,val|, ...}
//     array  -> array
//     null   -> none
//     other  -> value (bools, numerics, strings)
func JSONToArrai(data interface{}) (rel.Value, error) {
	switch v := data.(type) {
	case map[string]interface{}:
		return jsonObjToArrai(v)
	case []interface{}:
		return jsonArrToArrai(v)
	case string:
		return rel.NewTuple(rel.NewAttr("s", rel.NewString([]rune(v)))), nil
	case float64:
		return rel.NewNumber(v), nil
	case bool:
		return rel.NewTuple(rel.NewAttr("b", rel.NewBool(v))), nil
	case nil:
		return rel.NewTuple(rel.NewAttr("null", rel.None)), nil
	default:
		t, err := rel.NewValue(v)
		if err != nil {
			return nil, err
		}
		return rel.NewTuple(rel.NewAttr("v", t)), nil
	}
}

// Converts a json object to a binary relation {|@,@item|, |key,val|, ...}
func jsonObjToArrai(data map[string]interface{}) (rel.Value, error) {
	tuples := make([]rel.Value, len(data))
	i := 0
	for key, val := range data {
		// Recursively apply ToArrai to all values
		item, err := JSONToArrai(val)
		if err != nil {
			return nil, err
		}
		tuples[i] = rel.NewDictEntryTuple(rel.NewString([]rune(key)), item)
		i++
	}
	return rel.NewSet(tuples...), nil
}

// Converts a json array to an arrai array
func jsonArrToArrai(data []interface{}) (rel.Value, error) {
	elts := make([]rel.Value, len(data))
	for i, val := range data {
		// Recursively apply ToArrai to all elements
		elt, err := JSONToArrai(val)
		if err != nil {
			return nil, err
		}
		elts[i] = elt
	}
	return rel.NewTuple(rel.NewAttr("a", rel.NewArray(elts...))), nil
}

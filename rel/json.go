package rel

import (
	"encoding/json"
	"fmt"

	"github.com/go-errors/errors"
)

// MarshalToJSON marshals the given value to JSON with sets escaped as
// s -> {"{||}": s}.
func MarshalToJSON(value Value) []byte {
	j, err := json.Marshal(jsonEscape(value))
	if err != nil {
		panic(err)
	}
	return j
}

// UnmarshalFromJSON unmarshals the given value from JSON interpreting escaped
// sets ({"{||}": s}) accordingly.
func UnmarshalFromJSON(data []byte) (Value, error) {
	var i interface{}
	if err := json.Unmarshal(data, &i); err != nil {
		return nil, err
	}
	return jsonUnescape(i)
}

func array(intfs ...interface{}) []interface{} {
	args := make([]interface{}, len(intfs))
	for i, expr := range intfs {
		switch x := expr.(type) {
		case Expr:
			args[i] = jsonEscapeExpr(x)
		case nil:
		default:
			args[i] = x
		}
	}
	return args
}

func jsonEscapeExpr(expr Expr) interface{} {
	switch x := expr.(type) {
	case *Function:
		if x.Arg() != "-" {
			return array("\\", x.arg, x.body)
		}
		return array("\\", nil, x.body)
	case *IdentExpr:
		return array("ident", x.ident)
	case *BinExpr:
		return array(x.op, x.a, x.b)
	case *UnaryExpr:
		return array(x.op+"1", x.a)
	case *DotExpr:
		return array(".", x.lhs, x.attr)
	case *ArrowExpr:
		return array("->", x.lhs, x.fn)
	case *DArrowExpr:
		return array("=>", x.lhs, x.fn)
	case *SeqArrowExpr:
		return array(">>", x.lhs, x.fn)
	case *TupleMapExpr:
		return array(":>", x.lhs, x.fn)
	case *NestExpr:
		return array("nest", x.lhs, x.attrs, x.attr)
	case *UnnestExpr:
		return array("unnest", x.lhs, x.attr)
	case *IfElseExpr:
		return array("if", x.cond, x.ifTrue, x.ifFalse)
	case *TupleExpr:
		result := make([][2]interface{}, len(x.attrMap))
		i := 0
		for name, expr := range x.attrMap {
			result[i] = [2]interface{}{name, jsonEscapeExpr(expr)}
			i++
		}
		return array("tuple", result)
	case *SetExpr:
		result := make([]interface{}, len(x.elements))
		for i, expr := range x.elements {
			result[i] = jsonEscapeExpr(expr)
		}
		return array("set", result)
	case Value:
		v := jsonEscape(x)
		switch x := v.(type) {
		case nil, bool, float64, string, map[string]interface{}:
			return x
		case []interface{}:
			return array("array", x)
		default:
			return array("value", x)
		}
	default:
		panic(fmt.Sprintf("Not implemented for %#v", expr))
	}
}

func jsonEscape(value Expr) interface{} {
	switch x := value.(type) {
	case *Function:
		return map[string]interface{}{
			"{||}": map[string]interface{}{
				"\\": jsonEscapeExpr(x),
			},
		}
	case Number:
		return x.Float64()
	case Tuple:
		result := make(map[string]interface{}, x.Count())
		for e := x.Enumerator(); e.MoveNext(); {
			name, value := e.Current()
			result[name] = jsonEscape(value)
		}
		return result
	case GenericSet:
		if x.Equal(False) {
			return false
		}
		if x.Equal(True) {
			return true
		}
		array := make([]interface{}, 0, x.Count())
		for e := x.Enumerator(); e.MoveNext(); {
			array = append(array, jsonEscape(e.Current()))
		}
		return map[string]interface{}{"{||}": array}
	case String:
		return string(x.s)
	case Expr:
		panic(fmt.Sprintf("Bare expressions cannot be JSON escaped: %#v", x))
	default:
		panic(fmt.Sprintf("Unrecognised value: %v (%[1]T)", value))
	}
}

func jsonUnescape(i interface{}) (Value, error) {
	switch x := i.(type) {
	case bool:
		if x {
			return True, nil
		}
		return False, nil
	case float64:
		return NewNumber(x), nil
	case string:
		return NewString([]rune(x)), nil
	case []interface{}:
		values := make([]Value, len(x))
		for i, intf := range x {
			value, err := jsonUnescape(intf)
			if err != nil {
				return nil, err
			}
			values[i] = value
		}
		return NewArray(values...), nil
	case map[string]interface{}:
		if len(x) == 1 {
			for name, value := range x {
				if name == "{||}" {
					if array, ok := value.([]interface{}); ok {
						result := None
						for _, v := range array {
							value, err := jsonUnescape(v)
							if err != nil {
								return nil, err
							}
							result = result.With(value)
						}
						return result, nil
					}
					return nil, errors.Errorf(
						`x must be array in {"{||}": x}, not %T`, value)
				}
			}
		}
		var b TupleBuilder
		for name, v := range x {
			if name == "{||}" {
				return nil, errors.Errorf(`"{||}" is a reserved name`)
			}
			value, err := jsonUnescape(v)
			if err != nil {
				return nil, err
			}
			b.Put(name, value)
		}
		return b.Finish(), nil
	default:
		panic(fmt.Sprintf("Unrecognised value: %v (%[1]T)", i))
	}
}

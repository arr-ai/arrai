package syntax

import "github.com/arr-ai/arrai/rel"

func strArrToRelArr(s []string) rel.Value {
	values := make([]rel.Value, 0, len(s))
	for _, a := range s {
		values = append(values, rel.NewString([]rune(a)))
	}
	return rel.NewArray(values...)
}

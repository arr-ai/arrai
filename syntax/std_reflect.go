package syntax

import (
	"github.com/arr-ai/arrai/rel"
)

func stdReflect() rel.Attr {
	return rel.NewAttr("reflect", rel.NewTuple(
	// rel.NewNativeFunctionAttr("tuple", func(v rel.Value) rel.Value {
	// 	s := v.(rel.Set)
	// 	sets := make([]rel.Set, 0, s.Count()-1)
	// 	for e, ok := s.ArrayEnumerator(); ok && e.MoveNext(); {
	// 		sets = append(sets, e.Current().(rel.Set))
	// 	}
	// 	return rel.NUnion(sets...)
	// }),
	))
}

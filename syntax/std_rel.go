package syntax

import (
	"github.com/arr-ai/arrai/rel"
)

func stdRel() rel.Attr {
	return rel.NewAttr("rel", rel.NewTuple(
		rel.NewNativeFunctionAttr("union", func(v rel.Value) rel.Value {
			s := v.(rel.Set)
			sets := make([]rel.Set, 0, s.Count()-1)
			for e := s.ArrayEnumerator(); e.MoveNext(); {
				sets = append(sets, e.Current().(rel.Set))
			}
			return rel.NUnion(sets...)
		}),
	))
}

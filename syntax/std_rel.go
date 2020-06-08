package syntax

import (
	"github.com/arr-ai/arrai/rel"
)

func stdRel() rel.Attr {
	return rel.NewTupleAttr("rel",
		rel.NewNativeFunctionAttr("union", func(v rel.Value) (rel.Value, error) {
			s := v.(rel.Set)
			sets := make([]rel.Set, 0, s.Count())
			for e := s.Enumerator(); e.MoveNext(); {
				sets = append(sets, e.Current().(rel.Set))
			}
			return rel.NUnion(sets...), nil
		}),
	)
}

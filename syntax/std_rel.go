package syntax

import (
	"github.com/arr-ai/arrai/rel"
)

func stdRel() rel.Attr {
	return rel.NewTupleAttr("rel",
		rel.NewNativeFunctionAttr("union", func(v rel.Value) rel.Value {
			s := v.(rel.Set)
			sets := make([]rel.Set, 0, s.Count())
			var e rel.ValueEnumerator
			switch u := v.(type) {
			case rel.Array:
				var ok bool
				e, ok = u.ArrayEnumerator()
				if !ok {
					panic("wat")
				}
			case rel.Set:
				e = u.Enumerator()
			}
			for e.MoveNext() {
				sets = append(sets, e.Current().(rel.Set))
			}
			return rel.NUnion(sets...)
		}),
	)
}

package syntax

import (
	"context"

	"github.com/go-errors/errors"

	"github.com/arr-ai/arrai/rel"
)

func stdRel() rel.Attr {
	return rel.NewTupleAttr("rel",
		rel.NewNativeFunctionAttr("union", func(_ context.Context, v rel.Value) (rel.Value, error) {
			s, ok := v.(rel.Set)
			if !ok {
				return nil, errors.Errorf("arg to //rel.union must be set, not %s", rel.ValueTypeAsString(v))
			}
			sets := make([]rel.Set, 0, s.Count())
			for e := s.Enumerator(); e.MoveNext(); {
				c, ok := e.Current().(rel.Set)
				if !ok {
					return nil, errors.Errorf("elems of set arg to //rel.union must be sets, not %s",
						rel.ValueTypeAsString(e.Current()))
				}
				sets = append(sets, c)
			}
			return rel.NUnion(sets...), nil
		}),
	)
}

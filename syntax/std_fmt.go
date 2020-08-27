package syntax

import (
	"context"

	"github.com/arr-ai/arrai/rel"
)

func stdFmt() rel.Attr {
	return rel.NewTupleAttr(
		"fmt",
		rel.NewNativeFunctionAttr("pretty", func(_ context.Context, value rel.Value) (rel.Value, error) {
			prettifiedString, err := PrettifyString(value, 0)
			if err != nil {
				return nil, err
			}

			return rel.NewString([]rune(prettifiedString)), nil
		}),
	)
}

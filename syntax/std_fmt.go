package syntax

import (
	"fmt"

	"github.com/arr-ai/arrai/rel"
)

func stdFmt() rel.Attr {
	return rel.NewTupleAttr(
		"fmt",
		rel.NewNativeFunctionAttr("pretty", func(value rel.Value) (rel.Value, error) {
			formattedStr, err := FormatString(value, 0)
			if err != nil {
				return nil, err
			}

			fmt.Println(formattedStr)
			return rel.NewString([]rune("")), nil
		}),
	)
}

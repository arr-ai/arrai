package syntax

import (
	"fmt"
	"strings"

	"github.com/arr-ai/arrai/rel"
)

var stdSeqConcat = createNestedFunc("concat", 1, func(args ...rel.Value) rel.Value {
	var sb strings.Builder
	for i, ok := args[0].(rel.Set).ArrayEnumerator(); ok && i.MoveNext(); {
		sb.WriteString(mustAsString(i.Current()))
	}
	return rel.NewString([]rune(sb.String()))
})

func stdSeqRepeat(arg rel.Value) rel.Value {
	n := int(arg.(rel.Number))
	return rel.NewNativeFunction("repeat(n)", func(arg rel.Value) rel.Value {
		switch seq := arg.(type) {
		case rel.String:
			return rel.NewString([]rune(strings.Repeat(seq.String(), n)))
		case rel.Array:
			values := []rel.Value{}
			seqValues := seq.Values()
			for i := 0; i < n; i++ {
				values = append(values, seqValues...)
			}
			return rel.NewArray(values...)
		case rel.Set:
			if !seq.IsTrue() {
				return rel.None
			}
		}
		panic(fmt.Errorf("repeat: unsupported value: %v", arg))
	})
}

func stdSeq() rel.Attr {
	return rel.NewTupleAttr("seq",
		rel.NewAttr("concat", stdSeqConcat),
		rel.NewNativeFunctionAttr("repeat", stdSeqRepeat),
	)
}

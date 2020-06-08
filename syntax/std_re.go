package syntax

import (
	"regexp"

	"github.com/arr-ai/arrai/rel"
)

var (
	stdReMatch = rel.NewNativeFunction("compile", func(re rel.Value) (rel.Value, error) {
		regex := regexp.MustCompile(mustValueAsString(re))
		return rel.NewTuple(
			rel.NewNativeFunctionAttr("match", func(str rel.Value) (rel.Value, error) {
				s := mustValueAsString(str)
				matches := []rel.Value{}
				for _, m := range regex.FindAllStringSubmatchIndex(s, -1) {
					submatches := []rel.Value{}
					for i := 0; i < len(m); i += 2 {
						if m[i] >= 0 {
							submatches = append(submatches, rel.NewOffsetString([]rune(s[m[i]:m[i+1]]), m[i]))
						} else {
							submatches = append(submatches, nil)
						}
					}
					matches = append(matches, rel.NewArray(submatches...))
				}
				return rel.NewArray(matches...), nil
			}),
			createNestedFuncAttr("sub", 2, func(args ...rel.Value) (rel.Value, error) {
				r := mustValueAsString(args[0])
				s := mustValueAsString(args[1])
				return rel.NewString([]rune(regex.ReplaceAllString(s, r))), nil
			}),
			createNestedFuncAttr("subf", 2, func(args ...rel.Value) (rel.Value, error) {
				r := args[0]
				s := mustValueAsString(args[1])
				return rel.NewString([]rune(regex.ReplaceAllStringFunc(s, func(match string) string {
					result, err := rel.Call(r, rel.NewString([]rune(match)), rel.EmptyScope)
					if err != nil {
						panic(err)
					}
					return result.(rel.String).String()
				}))), nil
			}),
		), nil
	})
)

func stdRe() rel.Attr {
	return rel.NewTupleAttr("re",
		rel.NewAttr("compile", stdReMatch),
	)
}

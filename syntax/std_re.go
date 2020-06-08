package syntax

import (
	"fmt"
	"regexp"

	"github.com/arr-ai/arrai/rel"
)

var (
	stdReMatch = rel.NewNativeFunction("compile", func(re rel.Value) (rel.Value, error) {
		reStr, is := valueAsString(re)
		if !is {
			return nil, fmt.Errorf("//re.compile: re not a string: %v", re)
		}
		regex := regexp.MustCompile(reStr)
		return rel.NewTuple(
			rel.NewNativeFunctionAttr("match", func(str rel.Value) (rel.Value, error) {
				s, is := valueAsString(str)
				if !is {
					return nil, fmt.Errorf("//re.compile(re).match: s not a string: %v", str)
				}
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
				new, is := valueAsString(args[0])
				if !is {
					return nil, fmt.Errorf("//re.compile(re).sub: new not a string: %v", args[1])
				}
				s, is := valueAsString(args[1])
				if !is {
					return nil, fmt.Errorf("//re.compile(re).sub: s not a string: %v", args[0])
				}
				return rel.NewString([]rune(regex.ReplaceAllString(s, new))), nil
			}),
			createNestedFuncAttr("subf", 2, func(args ...rel.Value) (rel.Value, error) {
				f := args[0]
				s, is := valueAsString(args[1])
				if !is {
					return nil, fmt.Errorf("//re.compile(re).subf: s not a string: %v", args[0])
				}
				return rel.NewString([]rune(regex.ReplaceAllStringFunc(s, func(match string) string {
					result, err := rel.Call(f, rel.NewString([]rune(match)), rel.EmptyScope)
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

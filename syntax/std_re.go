package syntax

import (
	"context"
	"fmt"
	"regexp"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/tools"
)

var (
	stdReMatch = rel.NewNativeFunction("compile", func(_ context.Context, re rel.Value) (rel.Value, error) {
		reStr, is := tools.ValueAsString(re)
		if !is {
			return nil, fmt.Errorf("//re.compile: re not a string: %v", re)
		}
		regex, err := regexp.Compile(reStr)
		if err != nil {
			return nil, fmt.Errorf("//re.compile: %s", err)
		}
		return rel.NewTuple(
			rel.NewNativeFunctionAttr("match", func(_ context.Context, str rel.Value) (rel.Value, error) {
				s, is := tools.ValueAsString(str)
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
			createFunc2Attr("sub", func(_ context.Context, a, b rel.Value) (rel.Value, error) {
				new, is := tools.ValueAsString(a)
				if !is {
					return nil, fmt.Errorf("//re.compile(re).sub: new not a string: %v", b)
				}
				s, is := tools.ValueAsString(b)
				if !is {
					return nil, fmt.Errorf("//re.compile(re).sub: s not a string: %v", a)
				}
				return rel.NewString([]rune(regex.ReplaceAllString(s, new))), nil
			}),
			createFunc2Attr("subf", func(ctx context.Context, a, b rel.Value) (rel.Value, error) {
				s, is := tools.ValueAsString(b)
				if !is {
					return nil, fmt.Errorf("//re.compile(re).subf: s not a string: %v", a)
				}
				return rel.NewString([]rune(regex.ReplaceAllStringFunc(s, func(match string) string {
					result, err := rel.Call(ctx, a, rel.NewString([]rune(match)), rel.EmptyScope)
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

package syntax

import (
	"context"
	"regexp"
	"strings"

	"github.com/arr-ai/wbnf/ast"

	"github.com/arr-ai/arrai/rel"
)

var (
	leadingWSRE   = regexp.MustCompile(`\A[\t ]*`)
	lastWSRE      = regexp.MustCompile(`\n[\t ]*\z`)
	expansionRE   = regexp.MustCompile(`(?::([-+#*\.\_0-9a-z]*))(:(?:\\.|[^\\:}])*)?(?::((?:\\.|[^\\:}])*))?`)
	indentRE      = regexp.MustCompile(`(\n[\t ]*)\z`)
	firstIndentRE = regexp.MustCompile(`\A((\n[\t ]+)(?:\n)|(\n))`)
	lastSpacesRE  = regexp.MustCompile(`\n([ \t]*)\z`)
)

func (pc ParseContext) compileExpandableString(ctx context.Context, b ast.Branch, c ast.Children) (rel.Expr, error) {
	scanner := c.(ast.One).Node.One("quote").Scanner()
	quote := scanner.String()
	parts := []interface{}{}

	ws := quote[2:]
	trim := ""
	trimIndent := func(s string) {
		s = ws + s
		ws = ""
		if trim == "" {
			s = strings.TrimPrefix(s, "\n")
			i := leadingWSRE.FindStringIndex(s)
			trim = "\n" + s[:i[1]]
			s = s[i[1]:]
		}
		if trim != "\n" {
			s = strings.ReplaceAll(s, trim, "\n")
		}
		if s != "" {
			parts = append(parts, s)
		}
	}
	for i, part := range c.(ast.One).Node.Many("part") {
		p, part := which(part.(ast.Branch), "sexpr", "fragment")
		switch p {
		case "sexpr":
			if i == 0 || ws != "" {
				trimIndent("")
			}
			sexpr := part.(ast.One).Node.(ast.Branch)
			ws = sexpr.One("close").One("").(ast.Leaf).Scanner().String()[1:]
			parts = append(parts, sexpr)
		case "fragment":
			s := part.(ast.One).Node.One("").Scanner().String()
			s = parseArraiStringFragment(s, quote[1:2]+":", "")
			trimIndent(s)
		}
	}

	if len(parts) == 0 {
		return rel.None, nil
	}

	if last, is := parts[len(parts)-1].(string); is {
		if loc := lastWSRE.FindStringIndex(last); loc != nil {
			parts[len(parts)-1] = last[:loc[0]]
		}
	}

	next := ""
	exprs := make([]rel.Expr, len(parts))
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		switch part := part.(type) {
		case ast.Branch:
			format := ""
			delim := ""
			appendIfNotEmpty := ""
			if control := part.One("control").One("").(ast.Leaf).Scanner().String(); control != "" {
				m := expansionRE.FindStringSubmatchIndex(control)
				if m[2] >= 0 {
					format = control[m[2]:m[3]]
				}
				if m[4] >= 0 {
					delim = parseArraiStringFragment(control[m[4]:m[5]], ":}", "\n")
				}
				if m[6] >= 0 {
					appendIfNotEmpty = parseArraiStringFragment(control[m[6]:m[7]], ":}", "\n")
				}
			}
			if strings.HasPrefix(next, "\n") {
				if i > 0 {
					if s, ok := parts[i-1].(string); ok {
						if strings.HasSuffix(s, "\n") {
							appendIfNotEmpty += "\n"
							parts[i+1] = next[1:]
						}
					}
				} else {
					appendIfNotEmpty += "\n"
					parts[i+1] = next[1:]
				}
				next = ""
			}
			expr, err := pc.CompileExpr(ctx, part.One("expr").(ast.Branch))
			if err != nil {
				return nil, err
			}
			exprs[i] = rel.NewCallExprCurry(part.Scanner(), stdStrExpand,
				rel.NewString([]rune(format)), expr,
				rel.NewString([]rune(delim)),
				rel.NewString([]rune(appendIfNotEmpty)),
			)
		case string:
			next = part
		}
	}
	for i, part := range parts {
		if s, ok := part.(string); ok {
			exprs[i] = rel.NewTuple(rel.NewAttr("s", rel.NewString([]rune(s))))
		}
	}
	// TODO: Use a more direct approach to invoke concat implementation.
	return rel.NewCallExpr(b.Scanner(),
		rel.NewNativeFunction("xstr_concat", xstrConcat),
		rel.NewArrayExpr(b.Scanner(), exprs...)), nil
}

func xstrConcat(_ context.Context, seq rel.Value) (rel.Value, error) {
	// this is always a sequence of values between bare string and computed expressions
	// all bare strings are wrapped in a tuple of one attribute "s"
	//
	// bare strings are wrapped in a tuple to differentiate between
	// regular string and computed expressions
	values := cleanEmptyVal(seq.(rel.Array))
	recentIndent := "\n"
	if len(values) == 0 {
		return rel.None, nil
	}
	var sb strings.Builder
	for _, i := range values {
		// suppress empty string
		if !i.IsTrue() {
			continue
		}
		switch i := i.(type) {
		// handle computed expressions
		case rel.String:
			sb.WriteString(strings.ReplaceAll(i.String(), "\n", recentIndent))

		// handle bare string
		case rel.Tuple:
			v := i.MustGet("s")
			if !v.IsTrue() {
				continue
			}
			s := v.String()
			sb.WriteString(s)
			if m := indentRE.FindStringSubmatch(s); m != nil {
				recentIndent = m[1]
			}
		default:
			panic("xstrConcat: not receiving a string")
		}
	}
	return rel.NewString([]rune(sb.String())), nil
}

// cleanEmptyVal cleans whitespaces of bare strings before and after a computed empty string.
func cleanEmptyVal(values rel.Array) []rel.Value {
	arr := values.Values()
	length := len(arr)
	if length == 1 {
		return arr
	}

	getStr := func(i int) string {
		if t, isBareString := arr[i].(rel.Tuple); isBareString {
			if s := t.MustGet("s"); s.IsTrue() {
				return s.String()
			}
		}
		return ""
	}
	clean := func(i int) {
		if i < 0 || i >= length {
			return
		}

		switch {
		// e.g.
		// $`
		//     ${''}
		//         a
		// `
		case i == 0 && i < length-1:
			if s := getStr(i + 1); s != "" {
				if m := firstIndentRE.FindStringSubmatch(s); m != nil && m[1] != "" {
					match := m[1]
					arr[i+1] = arr[i+1].(rel.Tuple).With(
						"s",
						rel.NewString([]rune(strings.TrimPrefix(s, match))),
					)
				}
			}
		// e.g.
		// $`
		//     a:
		//         ${''}
		// `
		case i == length-1 && i > 0:
			if s := getStr(i - 1); s != "" {
				if m := lastSpacesRE.FindStringSubmatch(s); m != nil && m[1] != "" {
					match := m[1]
					arr[i-1] = arr[i-1].(rel.Tuple).With(
						"s",
						rel.NewString([]rune(strings.TrimSuffix(s, match))),
					)
				} else if trimmed := strings.TrimLeft(s, " "); trimmed == "" {
					// this is to remove any whitespace to the left the last empty evaluated str
					arr[i-1] = arr[i-1].(rel.Tuple).With(
						"s",
						rel.NewString([]rune("")),
					)
				}
			}
		case i > 0 && i < length-1:
			left, right := getStr(i-1), getStr(i+1)
			leftMatch, rightMatch := "", ""
			if m := lastSpacesRE.FindStringSubmatch(left); m != nil {
				leftMatch = m[1]
			}
			if m := firstIndentRE.FindStringSubmatch(right); m != nil {
				rightMatch = m[1]
			}

			// left and right needs to be cleaned
			// e.g.
			// $`
			//     a
			//         ${''}
			//         b
			// `
			if leftMatch != "" && rightMatch != "" {
				rightStr := strings.TrimPrefix(right, rightMatch)
				leftStr := strings.TrimSuffix(left, leftMatch)
				// Ensures indentation spaces are on the left string.
				// This is done because indentation processing in xstrConcat
				// is done from left to right.
				if m := leadingWSRE.FindStringSubmatch(rightStr); m != nil {
					newIndent := m[0]
					rightStr = strings.TrimPrefix(rightStr, newIndent)
					leftStr += newIndent
				}
				arr[i+1] = arr[i+1].(rel.Tuple).With("s", rel.NewString([]rune(rightStr)))
				arr[i-1] = arr[i-1].(rel.Tuple).With("s", rel.NewString([]rune(leftStr)))
			}
		}
	}
	shorten := func(i int) {
		arr = append(arr[:i], arr[i+1:]...)
		length--
	}
	for i := 0; i < length; {
		switch v := arr[i].(type) {
		case rel.Set:
			if !v.IsTrue() {
				clean(i)
				shorten(i)
				continue
			}
		case rel.Tuple:
			if getStr(i) == "" {
				shorten(i)
				continue
			}
		}
		i++
	}
	return arr
}

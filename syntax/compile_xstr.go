package syntax

import (
	"regexp"
	"strings"

	"github.com/arr-ai/wbnf/ast"

	"github.com/arr-ai/arrai/rel"
)

var (
	leadingWSRE   = regexp.MustCompile(`\A[\t ]*`)
	lastWSRE      = regexp.MustCompile(`\n[\t ]+\z`)
	expansionRE   = regexp.MustCompile(`(?::([-+#*\.\_0-9a-z]*))(:(?:\\.|[^\\:}])*)?(?::((?:\\.|[^\\:}])*))?`)
	indentRE      = regexp.MustCompile(`(\n[\t ]*)(?:[^\t ]|\z)[^\n]*\z`)
	firstIndentRE = regexp.MustCompile(`\A((\n[\t ]+)(?:\n)|(\n))`)
	lastSpacesRE  = regexp.MustCompile(`([ \t]*)\z`)
)

func (pc ParseContext) compileExpandableString(b ast.Branch, c ast.Children) rel.Expr {
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
		return rel.None
	}

	if last, is := parts[len(parts)-1].(string); is {
		if loc := lastWSRE.FindStringIndex(last); loc != nil {
			parts[len(parts)-1] = last[:loc[0]+1]
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
			exprs[i] = rel.NewCallExprCurry(part.Scanner(), stdStrExpand,
				rel.NewString([]rune(format)),
				pc.CompileExpr(part.One("expr").(ast.Branch)),
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
		rel.NewArrayExpr(b.Scanner(), exprs...))
}

func xstrConcat(seq rel.Value) (rel.Value, error) {
	values := cleanEmptyVal(seq)
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
		// handle sexpr
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

// this function cleans whitespaces of bare strings before and after a computed emptyt string
func cleanEmptyVal(values rel.Value) []rel.Value {
	arr := values.(rel.Array).Values()
	length := len(arr)
	cleanRE := func(re *regexp.Regexp, index int, cleaner func(string, string) string) {
		if index >= 0 && index < length {
			if t, isBareString := arr[index].(rel.Tuple); isBareString {
				if s := t.MustGet("s"); s.IsTrue() {
					match := ""
					if m := re.FindStringSubmatch(s.String()); m != nil {
						match = m[1]
					}
					arr[index] = t.With("s", rel.NewString([]rune(cleaner(match, s.String()))))
				}
			}
		}
	}
	clean := func(i int) {
		// cleans bare string after the empty computed string
		cleanRE(firstIndentRE, i+1, func(match, toReplace string) string {
			if match != "" {
				// cleans bare string before the empty computed string
				//
				// only does this if i+1 will be changed, this is meant to retain
				// any whitespaces in the bare string of arr[i-1].
				//
				// Meant to handle
				// $`
				//   abc
				//   ${''}def
				//  `
				cleanRE(lastSpacesRE, i-1, func(match, toReplace string) string {
					return strings.TrimSuffix(toReplace, match)
				})
				return strings.TrimPrefix(toReplace, match)
			}
			return toReplace
		})
	}
	for i := 0; i < length; i++ {
		switch v := arr[i].(type) {
		case rel.Set:
			if !v.IsTrue() {
				clean(i)
			}
		}
	}
	return arr
}

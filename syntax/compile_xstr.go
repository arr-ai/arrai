package syntax

import (
	"regexp"
	"strings"

	"github.com/arr-ai/wbnf/ast"

	"github.com/arr-ai/arrai/rel"
)

var leadingWSRE = regexp.MustCompile(`\A[\t ]*`)
var trailingWSRE = regexp.MustCompile(`[\t ]*\z`)
var lastWSRE = regexp.MustCompile(`\n[\t ]+\z`)
var expansionRE = regexp.MustCompile(`(?::([-+#*\.\_0-9a-z]*))(:(?:\\.|[^\\:}])*)?(?::((?:\\.|[^\\:}])*))?`)

func (pc ParseContext) compileExpandableString(c ast.Children) rel.Expr {
	quote := c.(ast.One).Node.One("quote").Scanner().String()
	parts := []interface{}{}
	{
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
			indent := ""
			if i > 0 {
				if s, ok := parts[i-1].(string); ok {
					indent = trailingWSRE.FindString(s)
				}
			}

			format := ""
			delim := ""
			appendIfNotEmpty := ""
			if control := part.One("control").One("").(ast.Leaf).Scanner().String(); control != "" {
				m := expansionRE.FindStringSubmatchIndex(control)
				if m[2] >= 0 {
					format = control[m[2]:m[3]]
				}
				if m[4] >= 0 {
					delim = parseArraiStringFragment(control[m[4]:m[5]], ":}", "\n"+indent)
				}
				if m[6] >= 0 {
					appendIfNotEmpty = parseArraiStringFragment(control[m[6]:m[7]], ":}", "\n"+indent)
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
			exprs[i] = rel.NewCallExprCurry(stdStrExpand,
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
			exprs[i] = rel.NewString([]rune(s))
		}
	}
	return rel.NewCallExpr(stdSeqConcat, rel.NewArrayExpr(exprs...))
}

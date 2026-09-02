package syntax

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/arr-ai/wbnf/ast"
	"github.com/arr-ai/wbnf/parser"

	"github.com/arr-ai/arrai/rel"
)

var (
	leadingWSRE   = regexp.MustCompile(`\A[\t ]*`)
	lastWSRE      = regexp.MustCompile(`\n[\t ]*\z`)
	expansionRE   = regexp.MustCompile(`(?::([-+#*\.\_0-9a-z]*))(:(?:\\.|[^\\:}])*)?(?::((?:\\.|[^\\:}])*))?`)
	indentRE      = regexp.MustCompile(`(\n[\t ]*)\z`)
	firstIndentRE = regexp.MustCompile(`\A((\n[\t ]+)(?:\n)|(\n))`)
	lastSpacesRE  = regexp.MustCompile(`\n([ \t]*)\z`)
	spaces        = regexp.MustCompile(`\A[\t ]+`)
)

func (pc ParseContext) compileExpandableString(ctx context.Context, b ast.Branch, c ast.Children) (rel.Expr, error) {
	scanner := c.(ast.One).Node.One("quote").Scanner()
	quote := scanner.String()
	parts := []interface{}{}

	ws := quote[2:]

	// retain initial spaces
	if spaces.MatchString(ws) {
		parts = append(parts, ws)
		ws = ""
	}

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
	interps := make(map[int]xstrPart, len(parts))
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
			interps[i] = xstrPart{expr: expr, format: format, delim: delim, tail: appendIfNotEmpty}
		case string:
			next = part
		}
	}
	xparts := make([]xstrPart, len(parts))
	for i, part := range parts {
		if s, ok := part.(string); ok {
			xparts[i] = xstrPart{literal: s}
		} else {
			xparts[i] = interps[i]
		}
	}
	return &xstrExpr{ExprScanner: rel.ExprScanner{Src: b.Scanner()}, parts: xparts}, nil
}

// xstrExpr is a compiled $-string: literal segments interleaved with
// interpolations. It renders straight into a byte builder — literal
// segments are plain Go strings interned at compile time, never arr.ai
// values, and interpolations expand without the Array-of-parts and curried
// native calls the old path allocated per evaluation.
type xstrExpr struct {
	rel.ExprScanner
	parts []xstrPart
}

// xstrPart is one segment: a literal (expr == nil) or an interpolation with
// its ${...:format:delim:tail} controls.
type xstrPart struct {
	literal string
	expr    rel.Expr
	format  string
	delim   string
	tail    string // appended iff the expansion is non-empty
}

// xstrPiece is a rendered part: the string it contributes and whether it
// was a bare literal (which drives indent tracking) or computed (which has
// newlines rewritten to the most recent literal indent).
type xstrPiece struct {
	s    string
	bare bool
}

func (e *xstrExpr) Source() parser.Scanner {
	return e.Src
}

func (e *xstrExpr) String() string {
	return "$-string"
}

// Eval renders the template. The semantics — expansion, empty-part
// whitespace cleanup, indent propagation — are exactly the old
// xstrConcat/cleanEmptyVal pipeline's, ported from values to strings.
func (e *xstrExpr) Eval(ctx context.Context, local rel.Scope) (rel.Value, error) {
	pieces := make([]xstrPiece, 0, len(e.parts))
	for _, p := range e.parts {
		if p.expr == nil {
			pieces = append(pieces, xstrPiece{s: p.literal, bare: true})
			continue
		}
		v, err := p.expr.Eval(ctx, local)
		if err != nil {
			return nil, rel.WrapContextErr(err, e, local)
		}
		s, err := xstrExpand(ctx, p.format, v, p.delim, p.tail)
		if err != nil {
			return nil, rel.WrapContextErr(err, e, local)
		}
		pieces = append(pieces, xstrPiece{s: s})
	}
	pieces = cleanEmptyPieces(pieces)
	if len(pieces) == 0 {
		return rel.None, nil
	}
	recentIndent := "\n"
	var sb strings.Builder
	for _, p := range pieces {
		if p.s == "" {
			continue
		}
		if p.bare {
			sb.WriteString(p.s)
			if m := indentRE.FindStringSubmatch(p.s); m != nil {
				recentIndent = m[1]
			}
		} else {
			sb.WriteString(strings.ReplaceAll(p.s, "\n", recentIndent))
		}
	}
	return rel.NewGoString(sb.String()), nil
}

// xstrExpand renders one interpolated value per //str.expand's rules.
func xstrExpand(ctx context.Context, format string, value rel.Value, delim, tail string) (string, error) {
	f := "%v"
	if format != "" {
		f = "%" + format
	}
	var s string
	if strings.HasPrefix(delim, ":") {
		forced, err := rel.Observe(value)
		if err != nil {
			return "", err
		}
		array, is := rel.AsArray(forced.(rel.Set))
		if !is {
			return "", fmt.Errorf("expansion arg not an array in ${arg::}: %v", forced)
		}
		var jb strings.Builder
		for i, v := range array.Values() {
			if i > 0 {
				jb.WriteString(delim[1:])
			}
			if v != nil {
				jb.WriteString(formatValue(ctx, f, v))
			}
		}
		s = jb.String()
	} else {
		s = formatValue(ctx, f, value)
	}
	if s != "" {
		s += tail
	}
	return s, nil
}

// cleanEmptyPieces cleans whitespace of bare literals before and after a
// computed empty part, then drops empty parts.
func cleanEmptyPieces(arr []xstrPiece) []xstrPiece {
	length := len(arr)
	if length == 1 {
		return arr
	}

	getStr := func(i int) string {
		if arr[i].bare {
			return arr[i].s
		}
		return ""
	}
	setStr := func(i int, s string) {
		arr[i].s = s
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
					setStr(i+1, strings.TrimPrefix(s, match))
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
					setStr(i-1, strings.TrimSuffix(s, match))
				} else if trimmed := strings.TrimLeft(s, " "); trimmed == "" {
					// this is to remove any whitespace to the left the last empty evaluated str
					setStr(i-1, "")
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
				setStr(i+1, rightStr)
				setStr(i-1, leftStr)
			}
		}
	}
	shorten := func(i int) {
		arr = append(arr[:i], arr[i+1:]...)
		length--
	}
	for i := 0; i < length; {
		if arr[i].s == "" {
			if !arr[i].bare {
				clean(i)
			}
			shorten(i)
			continue
		}
		i++
	}
	return arr
}

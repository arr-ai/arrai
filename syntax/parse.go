package syntax

import (
	"bytes"
	"math"

	"github.com/arr-ai/arrai/rel"
)

type noParseType struct{}

func (*noParseType) Error() string {
	return "No parse"
}

var noParse = &noParseType{}

type parseFunc func(l *Lexer) (rel.Expr, error)

// MustParseString parses input string and returns the parsed Expr or panics.
func MustParseString(s string) rel.Expr {
	return MustParse(NewStringLexer(s))
}

// MustParse parses input and returns the parsed Expr or panics.
func MustParse(l *Lexer) rel.Expr {
	expr, err := Parse(l)
	if err != nil {
		panic(err)
	}
	return expr
}

// ParseString parses input string and returns the parsed Expr or an error.
func ParseString(s string) (rel.Expr, error) {
	return Parse(NewStringLexer(s))
}

// Parse parses input and returns the parsed Expr or an error.
func Parse(l *Lexer) (rel.Expr, error) {
	expr, err := parseExpr(l)
	if err != nil {
		return nil, err
	}
	if !l.Scan(EOF) {
		l.Failf("input not consumed")
		return nil, l.Error()
	}
	return expr, nil
}

func parseExpr(l *Lexer) (rel.Expr, error) {
	return parseNullaryFunc(l)
}

func parseNullaryFunc(l *Lexer) (rel.Expr, error) {
	i := 0
	for l.Scan(Token('&')) {
		i++
	}
	expr, err := parseArrow(l)
	if err != nil {
		return nil, err
	}
	for i > 0 {
		expr = rel.NewFunction("-", expr)
		i--
	}
	return expr, nil
}

var arrowOps = map[Token]newBinOpFunc{
	ARROW:   rel.NewArrowExpr,
	ATARROW: rel.NewAtArrowExpr,
	DARROW:  rel.NewDArrowExpr,
	ORDER:   rel.NewOrderExpr,
	WHERE:   rel.NewWhereExpr,
	SUM:     rel.NewSumExpr,
	MAX:     rel.NewMaxExpr,
	MEAN:    rel.NewMeanExpr,
	MEDIAN:  rel.NewMedianExpr,
	MIN:     rel.NewMinExpr,
}

func parseArrow(l *Lexer) (rel.Expr, error) {
	a, err := parseWith(l)
	if err != nil {
		return nil, err
	}
	parsedTail := false
	for {
		if l.Scan(NEST) {
			names, err := parseNameList(l)
			if err != nil {
				return nil, err
			}
			if names == nil {
				return nil, expecting(l, "after nest", "name list")
			}
			if !l.Scan(IDENT) {
				return nil, expecting(l, "after nest |name,...|", "ident")
			}
			a = rel.NewNestExpr(
				a, rel.NewNames(names...), string(l.Lexeme()),
			)
			parsedTail = false
		} else if l.Scan(UNNEST) {
			if !l.Scan(IDENT) {
				return nil, expecting(l, "after nest |name,...|", "ident")
			}
			a = rel.NewUnnestExpr(a, string(l.Lexeme()))
			parsedTail = false
		} else if !parsedTail {
			a, err = parseBinOpTail(a, l, arrowOps, parseWith)
			if err != nil {
				return nil, err
			}
			parsedTail = true
		} else {
			return a, nil
		}
	}
}

var withOps = map[Token]newBinOpFunc{
	WITH:    rel.NewWithExpr,
	WITHOUT: rel.NewWithoutExpr,
}

func parseWith(l *Lexer) (rel.Expr, error) {
	return parseBinOp(l, withOps, parseLogicalOps)
}

var logicalOps = map[Token]newBinOpFunc{
	AND: rel.MakeBinValExpr("and", func(a, b rel.Value) rel.Value {
		if !a.Bool() {
			return a
		}
		return b
	}),
	OR: rel.MakeBinValExpr("and", func(a, b rel.Value) rel.Value {
		if a.Bool() {
			return a
		}
		return b
	}),
}

func parseLogicalOps(l *Lexer) (rel.Expr, error) {
	return parseBinOp(l, logicalOps, parseEqOps)
}

var eqOps = map[Token]newBinOpFunc{
	'=': rel.MakeEqExpr("=", func(a, b rel.Value) bool { return a.Equal(b) }),
	'<': rel.MakeEqExpr("<", func(a, b rel.Value) bool { return a.Less(b) }),
	'>': rel.MakeEqExpr(">", func(a, b rel.Value) bool { return b.Less(a) }),
	NEQ: rel.MakeEqExpr("!=", func(a, b rel.Value) bool { return !a.Equal(b) }),
	LEQ: rel.MakeEqExpr("<=", func(a, b rel.Value) bool { return !b.Less(a) }),
	GEQ: rel.MakeEqExpr(">=", func(a, b rel.Value) bool { return !a.Less(b) }),
}

func parseEqOps(l *Lexer) (rel.Expr, error) {
	next := parseIfElse
	a, err := next(l)
	if err != nil {
		return nil, err
	}
	e := a
	for op, found := eqOps[l.Peek()]; found; op, found = eqOps[l.Peek()] {
		l.Scan()
		b, err := next(l)
		if err != nil {
			return nil, err
		}
		// Chain eqOps such that, e.g.: a < b < c = (a < b) and (b < c).
		if e == a {
			a = op(a, b)
		} else {
			a = logicalOps[AND](a, op(e, b))
		}
		e = b
	}
	return a, nil
}

func parseIfElse(l *Lexer) (rel.Expr, error) {
	expr, err := parseAdd(l)
	if err != nil {
		return nil, err
	}
	for l.Scan(IF) {
		cond, err := parseExpr(l)
		if err != nil {
			return nil, err
		}
		if !l.Scan(ELSE) {
			return nil, expecting(l, "after `expr if pred`", "'else'")
		}
		ifFalse, err := parseExpr(l)
		if err != nil {
			return nil, err
		}
		expr = rel.NewIfElseExpr(expr, cond, ifFalse)
	}
	return expr, nil
}

var addOps = map[Token]newBinOpFunc{
	Token('+'): rel.NewAddExpr,
	Token('-'): rel.NewSubExpr,
	JOIN:       rel.NewJoinExpr,
}

func parseAdd(l *Lexer) (rel.Expr, error) {
	return parseBinOp(l, addOps, parseMul)
}

var mulOps = map[Token]newBinOpFunc{
	Token('*'): rel.NewMulExpr,
	Token('/'): rel.NewDivExpr,
	Token('%'): rel.NewModExpr,
	SUBMOD:     rel.NewSubModExpr,
	IDIV:       rel.NewIdivExpr,
}

func parseMul(l *Lexer) (rel.Expr, error) {
	return parseBinOp(l, mulOps, parsePow)
}

var powOp = map[Token]newBinOpFunc{
	Token('^'): rel.NewPowExpr,
}

func parsePow(l *Lexer) (rel.Expr, error) {
	return parseBinOp(l, powOp, parseAmbPrefix)
}

// Ambiguous prefixes also have binary counterparts. To avoid unexpected
// interactions with function calls, these must have lower precedence.
var ambPrefixOps = map[Token]newUnOpFunc{
	Token('+'): rel.NewPosExpr,
	Token('-'): rel.NewNegExpr,
	Token('^'): rel.NewPowerSetExpr,
}

func parseAmbPrefix(l *Lexer) (rel.Expr, error) {
	return parsePrefixOp(l, ambPrefixOps, parseCall)
}

func parseCall(l *Lexer) (rel.Expr, error) {
	lhs, err := parsePrefix(l)
	if err != nil {
		return nil, err
	}
	for {
		// Parse `a*b` as `a * b`, not `a (*b)`.
		if l.Peek() == Token('*') {
			return lhs, nil
		}

		rhs, err := parsePrefix(l)
		if err != nil {
			if err == noParse {
				return lhs, nil
			}
			return nil, err
		}
		lhs = rel.NewCallExpr(lhs, rhs)
	}
}

var prefixOps = map[Token]newUnOpFunc{
	Token('!'): rel.NewNotExpr,
	Token('*'): rel.NewEvalExpr,
}

func parsePrefix(l *Lexer) (rel.Expr, error) {
	return parsePrefixOp(l, prefixOps, parseTouch)
}

func parseTouch(l *Lexer) (rel.Expr, error) {
	lhs, err := parseSuffix(l)
	if err != nil {
		return nil, err
	}
	return parseTouchTail(l, lhs)
}

func makeTupleExpr(name string, attr rel.Expr) rel.Expr {
	attrExpr, err := rel.NewAttrExpr(name, attr)
	if err != nil {
		panic(err)
	}
	return rel.NewTupleExpr(
		rel.NewWildcardExpr(rel.NewIdentExpr(".")),
		attrExpr,
	)
}

func parseTouchTail(l *Lexer, expr rel.Expr) (rel.Expr, error) {
	path := make([]string, 0, 4) // A bit of spare buffer
	for l.Scan(ARROWST) {
		if !l.Scan(IDENT, Token('&'), STRING) {
			return nil, expecting(l, "after '.'", "ident", "string", "'*'")
		}
		if l.Token() == STRING {
			path = append(path, ParseArraiString(l.Lexeme()))
		} else if l.Token() == Token('&') {
			if !l.Scan(IDENT) {
				return nil, expecting(l, "after '.&'-prefix", "ident")
			}
			path = append(path, "&"+string(l.Lexeme()))
		} else {
			path = append(path, string(l.Lexeme()))
		}
	}
	if len(path) == 0 {
		return expr, nil
	}
	leaf := path[len(path)-1]
	attrExpr, err := parseAttrExpr(l, leaf)
	tupleExpr := makeTupleExpr(leaf, attrExpr)
	if err != nil {
		return nil, err
	}
	for i := len(path) - 2; i >= 0; i-- {
		tupleFunc := rel.NewFunction(".", tupleExpr)
		dotExpr := rel.NewDotExpr(rel.DotIdent, path[i])
		arrowExpr := rel.NewArrowExpr(dotExpr, tupleFunc)
		tupleExpr = makeTupleExpr(path[i], arrowExpr)
		if err != nil {
			panic(err)
		}
	}
	return rel.NewArrowExpr(expr, rel.NewFunction(".", tupleExpr)), nil
}

func parseSuffix(l *Lexer) (rel.Expr, error) {
	dot, err := parseDot(l)
	if err != nil {
		return nil, err
	}
	if l.Scan(COUNT) {
		return rel.NewCountExpr(dot), nil
	}
	return dot, nil
}

func parseDot(l *Lexer) (rel.Expr, error) {
	atom, err := ParseAtom(l)
	if err != nil {
		return nil, err
	}
	return parseDotTail(l, atom)
}

func parseDotTail(l *Lexer, expr rel.Expr) (rel.Expr, error) {
	for l.Scan(Token('.')) {
		if !l.Scan(IDENT, Token('&'), STRING, Token('*')) {
			return nil, expecting(l, "after '.'", "ident", "string", "'*'")
		}
		if l.Token() == STRING {
			expr = rel.NewDotExpr(expr, ParseArraiString(l.Lexeme()))
		} else if l.Token() == Token('&') {
			if !l.Scan(IDENT) {
				return nil, expecting(l, "after '.&'-prefix", "ident")
			}
			expr = rel.NewDotExpr(expr, "&"+string(l.Lexeme()))
		} else {
			expr = rel.NewDotExpr(expr, string(l.Lexeme()))
		}
	}
	return expr, nil
}

// ParseAtom parses a single arrai value.
func ParseAtom(l *Lexer) (rel.Expr, error) {
	tok := l.Peek()
	switch tok {
	case NUMBER:
		l.Scan()
		return l.Value(), nil

	case IDENT:
		l.Scan()
		ident := string(l.Lexeme())
		switch ident {
		case "false":
			return rel.False, nil
		case "true":
			return rel.True, nil
		case "none":
			return rel.None, nil
		}
		return rel.NewIdentExpr(ident), nil

	case STRING:
		l.Scan()
		return rel.NewString([]rune(ParseArraiString(l.Lexeme()))), nil

	case Token('('):
		l.Scan()
		expr, err := parseExpr(l)
		if err != nil {
			return nil, err
		}
		if !l.Scan(Token(')')) {
			return nil, expecting(l, "after expr", "')'")
		}
		return expr, nil

	case Token('<'):
		return parseXML(l, newXMLContext())

	case Token('.'):
		l.Scan()
		var expr rel.Expr = rel.NewIdentExpr(".")
		if l.Scan(IDENT, Token('*')) {
			expr = rel.NewDotExpr(expr, string(l.Lexeme()))
		}
		return expr, nil

	case PI:
		l.Scan()
		return rel.NewNumber(math.Pi), nil

	case SQRT:
		l.Scan()
		expr, err := parseExpr(l)
		if err != nil {
			return nil, err
		}
		return rel.NewPowExpr(expr, rel.NewNumber(0.5)), nil

	case Token('\\'):
		l.Scan()
		if !l.Scan(IDENT, Token('-'), Token('.')) {
			return nil, expecting(l, "after '\\'", "ident", "'-'")
		}
		arg := string(l.Lexeme())
		body, err := parseExpr(l)
		if err != nil {
			return nil, err
		}
		return rel.NewFunction(arg, body), nil

	case Token('{'):
		l.Scan()
		attrs := []rel.AttrExpr{}
	tokenLoop:
		for {
			var name string
			switch l.Peek() {
			case STRING:
				l.Scan()
				name = ParseArraiString(l.Lexeme())
			case IDENT:
				l.Scan()
				name = string(l.Lexeme())
			case Token('}'):
				break tokenLoop
			case Token('&'):
				l.Scan()
				if !l.Scan(Token(IDENT)) {
					return nil, expecting(l, "after '&'-prefix", "ident")
				}
				name = "&" + string(l.Lexeme())
			}
			expr, err := parseAttrExpr(l, name)
			attr, err := rel.NewAttrExpr(name, expr)
			if err != nil {
				return nil, err
			}
			attrs = append(attrs, attr)
			if !l.Scan(',') {
				break
			}
		}

		if !l.Scan('}') {
			return nil, expecting(l, "after tuple body", "'}'")
		}
		return rel.NewTupleExpr(attrs...), nil

	case OSET, Token('['):
		l.Scan()

		var closer Token
		if tok == OSET {
			closer = CSET
		} else {
			closer = Token(']')
		}

		if closer == CSET {
			names, err := parseNameList(l)
			if err != nil {
				return nil, err
			}
			if names != nil {
				tuples := [][]rel.Expr{}
				for l.Scan(Token('{')) {
					exprs, err := parseExprCommaList(
						l, Token('}'), "'}'", "after relation body")
					if err != nil {
						return nil, err
					}
					tuples = append(tuples, exprs)
					if !l.Scan(Token(',')) {
						break
					}
				}
				if !l.Scan(CSET) {
					return nil, expecting(l, "after set tuple", "'|}'")
				}
				return rel.NewRelationExpr(names, tuples...)
			}
		}

		elts, err := parseExprCommaList(l, closer, "'|}'", "after set body")
		if err != nil {
			return nil, err
		}
		if closer == CSET {
			return rel.NewSetExpr(elts...), nil
		}
		return rel.NewArrayExpr(elts...), nil

	case ERROR:
		l.Failf("syntax error")
		return nil, l.Error()

	default:
		return nil, noParse
	}
}

func parseAttrExpr(l *Lexer, name string) (rel.Expr, error) {
	if name != "" && !l.Scan(Token(':')) {
		return nil, expecting(l, "after tuple name", "':'")
	}
	expr, err := parseExpr(l)
	if err != nil {
		return nil, err
	}
	if name == "" {
		e := expr
		for b, ok := e.(rel.LHSExpr); ok; b, ok = e.(rel.LHSExpr) {
			e = b.LHS()
		}
		if dot, ok := e.(*rel.DotExpr); ok {
			name = dot.Attr()
		} else {
			return nil, expecting(
				l, "after omitted attr ident", "expr with ident LHS")
		}
	}
	if name[:1] == "&" {
		expr = rel.NewFunction("-", expr)
	}
	return expr, nil
}

func parseNameList(l *Lexer) ([]string, error) {
	if !l.Scan(Token('|')) {
		return nil, nil
	}
	// relation shorthand, e.g.: {| |a,b| {1,2}, {3,4} |}
	names := []string{}
	for l.Scan(IDENT) {
		names = append(names, string(l.Lexeme()))
		if !l.Scan(Token(',')) {
			break
		}
	}
	if !l.Scan(Token('|')) {
		return nil, expecting(l, "after name-list", "'|'")
	}
	return names, nil
}

type newBinOpFunc func(a, b rel.Expr) rel.Expr
type newUnOpFunc func(e rel.Expr) rel.Expr

func parseBinOp(
	l *Lexer, ops map[Token]newBinOpFunc, next parseFunc,
) (rel.Expr, error) {
	a, err := next(l)
	if err != nil {
		return nil, err
	}
	return parseBinOpTail(a, l, ops, next)
}

func parseBinOpTail(
	a rel.Expr, l *Lexer, ops map[Token]newBinOpFunc, next parseFunc,
) (rel.Expr, error) {
	for op, found := ops[l.Peek()]; found; op, found = ops[l.Peek()] {
		l.Scan()
		b, err := next(l)
		if err != nil {
			return nil, err
		}
		a = op(a, b)
	}
	return a, nil
}

func parsePrefixOp(
	l *Lexer, ops map[Token]newUnOpFunc, next parseFunc,
) (rel.Expr, error) {
	if op, found := ops[l.Peek()]; found {
		l.Scan()
		rhs, err := parsePrefixOp(l, ops, next)
		if err != nil {
			return nil, err
		}
		return op(rhs), nil
	}
	return next(l)
}

func parseExprCommaList(
	l *Lexer, delim Token, delimStr, context string,
) ([]rel.Expr, error) {
	elts := []rel.Expr{}
	err := parseCommaList(l, delim, delimStr, context, func(l *Lexer) error {
		elt, err := parseExpr(l)
		if err != nil {
			return err
		}
		elts = append(elts, elt)
		return nil
	})
	return elts, err
}

func parseCommaList(
	l *Lexer, delim Token, delimStr, context string, parse func(l *Lexer) error,
) error {
	var err error
	if l.Scan(delim) {
		return nil
	}
	for {
		err := parse(l)
		if err != nil {
			return err
		}
		if !l.Scan(Token(',')) {
			break
		}
		if l.Scan(delim) {
			return nil
		}
	}
	if err != nil {
		return err
	}
	if !l.Scan(delim) {
		return expecting(l, context, delimStr)
	}
	return nil
}

func expecting(l *Lexer, context string, expected ...string) error {
	var b bytes.Buffer
	n := len(expected)
	for i, x := range expected[:n-1] {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(x)
	}
	if n > 0 {
		if n > 2 {
			b.WriteString(", or ") // Oxford comma
		} else if n > 1 {
			b.WriteString(" or ")
		}
		b.WriteString(expected[n-1])
	}
	l.Failf("Expected %s %s, not %q (%s)",
		b.String(), context, l.Lexeme(), TokenRepr(l.Token()))
	return l.Error()
}

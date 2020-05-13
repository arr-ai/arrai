//+build !wasm

package shell

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/arr-ai/arrai/rel"
)

func (s *shellInstance) Do(line []rune, pos int) (newLine [][]rune, length int) {
	l := getLastToken(line[:pos])
	switch {
	case strings.HasSuffix(l, "///"):
		return [][]rune{line}, len(line)
	case strings.HasPrefix(l, "//"):
		var names []string
		var lastName string
		if l == "//" {
			lastName, names = "", []string{}
		} else {
			names = strings.Split(l[2:], ".")
			lastName, names = names[len(names)-1], names[:len(names)-1]
		}
		newLine, length = getScopePredictions(names, lastName, s.scope.MustGet(".").(rel.Tuple))
		if l == "//" {
			newLine = append(newLine, []rune("{"))
		}
	default:
		currentExpr := strings.Join(s.collector.withLine(string(line[:pos])).lines, "\n")
		realExpr, residue := s.trimExpr(currentExpr)
		val, err := tryEval(realExpr, s.scope)
		if err != nil {
			return
		}
		newLine = formatPredictions(predictFromValue(val), residue)
		length = len(residue)
	}
	return
}

//TODO: maybe needed for more advance predictions
// func evalExprWithResidue(s *shellInstance, expr string, scope rel.Scope) (rel.Value, string, error) {
// 	realExpr, residue := s.trimExpr(expr)
// 	val, err := tryEval(realExpr, scope)
// 	if err != nil {
// 		return nil, residue, err
// 	}
// 	return val, residue, nil
// }

func (s *shellInstance) trimExpr(expr string) (string, string) {
	if re := regexp.MustCompile("[(.]('|`|\")?$"); re.MatchString(expr) {
		prefix := re.FindString(expr)
		return expr[:len(expr)-len(prefix)], prefix
	}

	//TODO:
	// string get expression will be handled by the isBalanced operation`
	// getOp := regexp.MustCompile(`\.(\*|\.|[$@A-Za-z_][0-9$@A-Za-z_]*)$`)
	// if getOp.MatchString(expr) {
	// 	residue = getOp.FindString(expr)
	// 	realExpr = strings.TrimPrefix(expr, residue)
	// 	return
	// }

	// for i := len(expr); i > 0; i-- {
	// 	if s.collector.withLine(expr[:i]).isBalanced() {
	// 		return expr[:i], expr[i:]
	// 	}
	// }
	return expr, ""
}

func predictFromValue(val rel.Value) (predictions []string) {
	switch v := val.(type) {
	case rel.Tuple:
		predictions = v.Names().OrderedNames()
		for i := 0; i < len(predictions); i++ {
			predictions[i] = fmt.Sprintf(`.%s`, rel.TupleNameRepr(predictions[i]))
		}
	case rel.Dict:
		predictions = make([]string, 0, v.Count())
		for _, e := range v.OrderedEntries() {
			predictions = append(predictions, fmt.Sprintf(`(%s)`, rel.Repr(e.MustGet("@"))))
		}
	}
	return
}

func formatPredictions(predictions []string, prefix string) [][]rune {
	var p [][]rune
	for i := 0; i < len(predictions); i++ {
		if strings.HasPrefix(predictions[i], prefix) {
			p = append(p, []rune(strings.TrimPrefix(predictions[i], prefix)))
			// fmt.Printf("\nPrediction: %s\n", strings.TrimPrefix(predictions[i], prefix))
		}
	}
	return p
}

func getLastToken(line []rune) string {
	i := len(line) - 1
	for ; i > 0; i-- {
		if !isAlpha(line[i]) && line[i] != '.' {
			if line[i] == '/' {
				switch {
				case strings.HasSuffix(string(line[:i+1]), "///"):
					i -= 2
				case strings.HasSuffix(string(line[:i+1]), "//"):
					i--
				}
			} else {
				i++
			}
			break
		}
	}

	if i == len(line) || i < 0 {
		return ""
	}
	return string(line[i:])
}

func isAlpha(l rune) bool {
	return (l >= 'a' && l <= 'z') || (l >= 'A' && l <= 'Z')
}

func getScopePredictions(tuplePath []string, name string, scope rel.Tuple) ([][]rune, int) {
	var newLine [][]rune
	length := len(name)
	for _, attr := range tuplePath {
		if value, has := scope.Get(attr); has {
			if u, is := value.(rel.Tuple); is {
				scope = u
				continue
			}
		}
		return nil, 0
	}

	for _, attr := range scope.Names().OrderedNames() {
		if strings.HasPrefix(attr, name) {
			newLine = append(newLine, []rune(attr[length:]))
		}
	}
	return newLine, length
}

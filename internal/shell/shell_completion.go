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
		} else if lastName != "" {
			if len(newLine) == 0 {
				length = 0
			}
			names, lastName = append(names, lastName), ""
			predictions, _ := getScopePredictions(names, lastName, s.scope.MustGet(".").(rel.Tuple))
			for i := 0; i < len(predictions); i++ {
				predictions[i] = append([]rune("."), predictions[i]...)
			}
			newLine = append(newLine, predictions...)
		}
	default:
		notAttachedToExpr := pos > 0 && line[pos-1] == ' '

		// checks if it's not attached to any token or starts at pos 0
		if pos == 0 || l != "" || notAttachedToExpr {
			newLine, length = s.globalCompletions(l)
		}

		// no need to check partial expr if not attached to token
		if notAttachedToExpr {
			return
		}
		partialExprPreds, residueLen := s.partialExprPredictions(string(line[:pos]))
		if len(newLine) == 0 || residueLen < length {
			length = residueLen
		}
		newLine = append(newLine, partialExprPreds...)
	}
	return newLine, length
}

func (s *shellInstance) globalCompletions(prefix string) (newLine [][]rune, length int) {
	predictions := make([][]rune, 0)
	for _, p := range formatPredictions(s.scope.OrderedNames(), prefix) {
		if string(p) == "." || string(p) == "" {
			continue
		}
		predictions = append(predictions, p)
	}
	return predictions, len(prefix)
}

func (s *shellInstance) partialExprPredictions(currentLine string) (newLine [][]rune, length int) {
	currentExpr := strings.Join(s.collector.withLine(currentLine).lines, "\n")
	realExpr, residue := s.trimExpr(currentExpr)
	if residue != "" {
		val, err := tryEval(realExpr, s.scope)
		if err != nil {
			return
		}
		toAdd := formatPredictions(predictFromValue(val), residue)
		if !(len(toAdd) == 1 && len(toAdd[0]) == 0) {
			newLine = append(newLine, formatPredictions(predictFromValue(val), residue)...)
			length = len(residue)
		}
	}
	val, err := tryEval(currentExpr, s.scope)
	if err == nil {
		newLine = append(newLine, formatPredictions(predictFromValue(val), "")...)
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

func (s *shellInstance) trimExpr(expr string) (realExpr, residue string) {
	realExpr = expr
	getRe := regexp.MustCompile(`\.(\w*|\"(\\.|[^\\"])*|\'(\\.|[^\\'])*|\x60(\x60\x60|[^\x60])*)$`)
	callRe := regexp.MustCompile(`\((\"([^\\"])*|\'(\\.|[^\\'])*|\x60(\x60\x60|[^\x60])*|[^)]*)$`)

	if callRe.MatchString(expr) {
		residue = callRe.FindString(expr)
	} else if getRe.MatchString(expr) {
		residue = getRe.FindString(expr)
	}

	realExpr = strings.TrimSuffix(expr, residue)
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
	return
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
		if strings.HasPrefix(attr, name) && name != attr {
			newLine = append(newLine, []rune(attr[length:]))
		}
	}
	return newLine, length
}

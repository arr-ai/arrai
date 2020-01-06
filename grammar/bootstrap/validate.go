package bootstrap

import (
	"fmt"

	"github.com/arr-ai/arrai/grammar/parse"
)

func validationErrorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

func validateNode(v interface{}, expectedTag Rule, validate func(parse.Node) error) error {
	if node, ok := v.(parse.Node); ok {
		if node.Tag != string(expectedTag) {
			return validationErrorf("expecting tag `%s`, got `%s`", expectedTag, node.Tag)
		}
		return validate(node)
	}
	return validationErrorf("not a node: %v", v)
}

func validateScanner(
	v interface{},
	validate func(parse.Scanner) error,
) error {
	if scanner, ok := v.(parse.Scanner); ok {
		return validate(scanner)
	}
	return validationErrorf("not a scanner: %v", v)
}

func (t S) ValidateParse(g Grammar, rule Rule, v interface{}) error {
	return validateScanner(v, func(scanner parse.Scanner) error { return nil })
}

func (t RE) ValidateParse(g Grammar, rule Rule, v interface{}) error {
	return validateScanner(v, func(scanner parse.Scanner) error {
		return nil
		// if _, err := regexp.Parse()
	})
}

func (t Seq) ValidateParse(g Grammar, rule Rule, v interface{}) error {
	return validateNode(v, ruleOrAlt(rule, seqTag), func(node parse.Node) error {
		if node.Count() != len(t) {
			return validationErrorf("seq(%d): wrong number of children: %d", len(t), node.Count())
		}
		for i, term := range t {
			if err := term.ValidateParse(g, "", node.Children[i]); err != nil {
				return err
			}
		}
		return nil
	})
}

func (t Oneof) ValidateParse(g Grammar, rule Rule, v interface{}) error {
	return validateNode(v, ruleOrAlt(rule, oneofTag), func(node parse.Node) error {
		if n := node.Count(); n != 1 {
			return validationErrorf("oneof: expecting one child, got %d", n)
		}
		if i, ok := node.Extra.(int); ok {
			return t[i].ValidateParse(g, "", node.Children[0])
		}
		return validationErrorf("oneof: missing selected child")
	})
}

func (t Delim) ValidateParse(g Grammar, rule Rule, v interface{}) error {
	return validateNode(v, ruleOrAlt(rule, delimTag), func(node parse.Node) error {
		n := node.Count()
		if n == 0 {
			return validationErrorf("delim: no children")
		}
		if n%2 != 1 {
			return validationErrorf("delim: expecting odd number of children, not %d", n)
		}
		_, ok := node.Extra.(Associativity)
		if !ok {
			return validationErrorf("delim: missing depth")
		}

		left, right := t.LRTerms(node)

		if err := left.ValidateParse(g, "", node.Children[0]); err != nil {
			return err
		}
		for i := 1; i < n; i += 2 {
			if err := t.Sep.ValidateParse(g, "", node.Children[i]); err != nil {
				return err
			}
			if err := right.ValidateParse(g, "", node.Children[i+1]); err != nil {
				return err
			}
		}
		return nil
	})
}

func (t Quant) ValidateParse(g Grammar, rule Rule, v interface{}) error {
	return validateNode(v, ruleOrAlt(rule, quantTag), func(node parse.Node) error {
		n := node.Count()
		if n < t.Min || (t.Max != 0 && t.Max < n) {
			return validationErrorf("quant(%d..%d): wrong number of children: %d", t.Min, t.Max, n)
		}
		for _, child := range node.Children {
			if err := t.Term.ValidateParse(g, "", child); err != nil {
				return err
			}
		}
		return nil
	})
}

func (t Rule) ValidateParse(g Grammar, rule Rule, v interface{}) error {
	return g[t].ValidateParse(g, t, v)
}

//-----------------------------------------------------------------------------

func (t Stack) ValidateParse(g Grammar, rule Rule, v interface{}) error {
	panic(Inconceivable)
}

//-----------------------------------------------------------------------------

func (t Named) ValidateParse(g Grammar, rule Rule, v interface{}) error {
	// TODO: Be a little more thorough.
	return nil
}

//-----------------------------------------------------------------------------

func (t Diff) ValidateParse(g Grammar, rule Rule, v interface{}) error {
	// TODO: Be a little more thorough.
	return nil
}

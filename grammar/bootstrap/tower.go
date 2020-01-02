package bootstrap

func (t Rule) Resolve(oldRule, newRule Rule) Term {
	if t == oldRule {
		return newRule
	}
	return t
}

func (t RE) Resolve(oldRule, newRule Rule) Term {
	return t
}

func (t Seq) Resolve(oldRule, newRule Rule) Term {
	result := make(Seq, 0, len(t))
	for _, term := range t {
		result = append(result, term.Resolve(oldRule, newRule))
	}
	return result
}

func (t Oneof) Resolve(oldRule, newRule Rule) Term {
	result := make(Oneof, 0, len(t))
	for _, term := range t {
		result = append(result, term.Resolve(oldRule, newRule))
	}
	return result
}

func (t Tower) Resolve(oldRule, newRule Rule) Term {
	panic("should never get here")
}

func (t Delim) Resolve(oldRule, newRule Rule) Term {
	t.Term = t.Term.Resolve(oldRule, newRule)
	t.Sep = t.Sep.Resolve(oldRule, newRule)
	return t
}

func (t Quant) Resolve(oldRule, newRule Rule) Term {
	t.Term = t.Term.Resolve(oldRule, newRule)
	return t
}

func (t NamedTerm) Resolve(oldRule, newRule Rule) Term {
	t.Term = t.Term.Resolve(oldRule, newRule)
	return t
}

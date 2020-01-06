package bootstrap

import (
	"fmt"
	"reflect"
)

type DiffReport interface {
	Equal() bool
}

type InterfaceDiff struct {
	A, B interface{}
}

func (d InterfaceDiff) Equal() bool {
	return d.A == d.B
}

func diffInterfaces(a, b interface{}) InterfaceDiff {
	if diff := (InterfaceDiff{A: a, B: b}); !diff.Equal() {
		return diff
	}
	return InterfaceDiff{}
}

//-----------------------------------------------------------------------------

type GrammarDiff struct {
	OnlyInA []Rule
	OnlyInB []Rule
	Prods   map[Rule]TermDiff
}

func (d GrammarDiff) Equal() bool {
	return len(d.OnlyInA) == 0 && len(d.OnlyInB) == 0 && len(d.Prods) == 0
}

func DiffGrammars(a, b Grammar) GrammarDiff {
	diff := GrammarDiff{
		Prods: map[Rule]TermDiff{},
	}
	for rule, aTerm := range a {
		if bTerm, ok := b[rule]; ok {
			if td := DiffTerms(aTerm, bTerm); !td.Equal() {
				diff.Prods[rule] = td
			}
		} else {
			diff.OnlyInA = append(diff.OnlyInA, rule)
		}
	}
	for rule := range b {
		if _, ok := a[rule]; !ok {
			diff.OnlyInB = append(diff.OnlyInB, rule)
		}
	}
	if diff.Equal() != reflect.DeepEqual(a, b) {
		panic(fmt.Sprintf(
			"diff.Equal() == %v != %v == reflect.DeepEqual(a, b): %#v\n%#v\n%#v",
			diff.Equal(), reflect.DeepEqual(a, b), diff, a, b))
	}
	return diff
}

//-----------------------------------------------------------------------------

type TermDiff interface {
	DiffReport
}

type TypesDiffer struct {
	InterfaceDiff
}

func (d TypesDiffer) Equal() bool {
	return false
}

func DiffTerms(a, b Term) TermDiff {
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return TypesDiffer{
			InterfaceDiff: diffInterfaces(
				reflect.TypeOf(a).String(),
				reflect.TypeOf(b).String(),
			),
		}
	}
	switch a := a.(type) {
	case Rule:
		return diffRules(a, b.(Rule))
	case S:
		return diffSes(a, b.(S))
	case RE:
		return diffREs(a, b.(RE))
	case Seq:
		return diffSeqs(a, b.(Seq))
	case Oneof:
		return diffOneofs(a, b.(Oneof))
	case Stack:
		return diffTowers(a, b.(Stack))
	case Delim:
		return diffDelims(a, b.(Delim))
	case Quant:
		return diffQuants(a, b.(Quant))
	case Named:
		return diffNameds(a, b.(Named))
	case Diff:
		return diffDiff(a, b.(Diff))
	default:
		panic(fmt.Errorf("unknown term type: %v %[1]T", a))
	}
}

//-----------------------------------------------------------------------------

type RuleDiff struct {
	A, B Rule
}

func (d RuleDiff) Equal() bool {
	return d.A == d.B
}

func diffRules(a, b Rule) RuleDiff {
	return RuleDiff{A: a, B: b}
}

//-----------------------------------------------------------------------------

type SDiff struct {
	A, B S
}

func (d SDiff) Equal() bool {
	return d.A == d.B
}

func diffSes(a, b S) SDiff {
	return SDiff{A: a, B: b}
}

//-----------------------------------------------------------------------------

type REDiff struct {
	A, B RE
}

func (d REDiff) Equal() bool {
	return d.A == d.B
}

func diffREs(a, b RE) REDiff {
	return REDiff{A: a, B: b}
}

//-----------------------------------------------------------------------------

type termsesDiff struct {
	Len   InterfaceDiff
	Terms []TermDiff
}

func (d termsesDiff) Equal() bool {
	return d.Len.Equal() && d.Terms == nil
}

func diffTermses(a, b []Term) termsesDiff {
	var tsd termsesDiff
	tsd.Len = diffInterfaces(len(a), len(b))
	lenDiff := len(a) - len(b)
	switch {
	case lenDiff < 0:
		b = b[:len(a)]
	case lenDiff > 0:
		a = a[:len(b)]
	}

	for i, term := range a {
		if td := DiffTerms(term, b[i]); !td.Equal() {
			tsd.Terms = append(tsd.Terms, td)
		}
	}
	return tsd
}

type SeqDiff termsesDiff

func (d SeqDiff) Equal() bool {
	return (termsesDiff(d)).Equal()
}

func diffSeqs(a, b Seq) SeqDiff {
	return SeqDiff(diffTermses(a, b))
}

type OneofDiff termsesDiff

func (d OneofDiff) Equal() bool {
	return (termsesDiff(d)).Equal()
}

func diffOneofs(a, b Oneof) OneofDiff {
	return OneofDiff(diffTermses(a, b))
}

type TowerDiff termsesDiff

func (d TowerDiff) Equal() bool {
	return (termsesDiff(d)).Equal()
}

func diffTowers(a, b Stack) TowerDiff {
	return TowerDiff(diffTermses(a, b))
}

//-----------------------------------------------------------------------------

type DelimDiff struct {
	Term            TermDiff
	Sep             TermDiff
	Assoc           InterfaceDiff
	CanStartWithSep InterfaceDiff
	CanEndWithSep   InterfaceDiff
}

func (d DelimDiff) Equal() bool {
	return d.Term.Equal() &&
		d.Sep.Equal() &&
		d.Assoc.Equal() &&
		d.CanStartWithSep.Equal() &&
		d.CanEndWithSep.Equal()
}

func diffDelims(a, b Delim) DelimDiff {
	return DelimDiff{
		Term:            DiffTerms(a.Term, b.Term),
		Sep:             DiffTerms(a.Sep, b.Sep),
		Assoc:           diffInterfaces(a.Assoc, b.Assoc),
		CanStartWithSep: diffInterfaces(a.CanStartWithSep, b.CanStartWithSep),
		CanEndWithSep:   diffInterfaces(a.CanEndWithSep, b.CanEndWithSep),
	}
}

//-----------------------------------------------------------------------------

type QuantDiff struct {
	Term TermDiff
	Min  InterfaceDiff
	Max  InterfaceDiff
}

func (d QuantDiff) Equal() bool {
	return d.Term.Equal() && d.Min.Equal() && d.Max.Equal()
}

func diffQuants(a, b Quant) QuantDiff {
	return QuantDiff{
		Term: DiffTerms(a.Term, b.Term),
		Min:  diffInterfaces(a.Min, b.Min),
		Max:  diffInterfaces(a.Max, b.Max),
	}
}

//-----------------------------------------------------------------------------

type NamedDiff struct {
	Name InterfaceDiff
	Term TermDiff
}

func (d NamedDiff) Equal() bool {
	return d.Name.Equal() && d.Term.Equal()
}

func diffNameds(a, b Named) NamedDiff {
	return NamedDiff{
		Name: diffInterfaces(a.Name, b.Name),
		Term: DiffTerms(a.Term, b.Term),
	}
}

//-----------------------------------------------------------------------------

type DiffDiff struct {
	A, B TermDiff
}

func (d DiffDiff) Equal() bool {
	return d.A.Equal() && d.B.Equal()
}

func diffDiff(a, b Diff) DiffDiff {
	return DiffDiff{
		A: DiffTerms(a.A, b.A),
		B: DiffTerms(a.B, b.B),
	}
}

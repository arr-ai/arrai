package syntax

import "testing"

func TestJoin(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{}`, `{} <&> {}`)
	AssertCodesEvalToSameValue(t, `{}`, `{()} <&> {}`)
	AssertCodesEvalToSameValue(t, `{}`, `{} <&> {()}`)
	AssertCodesEvalToSameValue(t, `{()}`, `{()} <&> {()}`)

	AssertCodesEvalToSameValue(t, `{}`, `{|a|(1)} <&> {|a|(2)}`)
	AssertCodesEvalToSameValue(t, `{|a|(1)}`, `{|a|(1)} <&> {|a|(1)}`)

	AssertCodesEvalToSameValue(t, `{|a,b|(1,4)}`, `{|a|(1)} <&> {|b|(4)}`)
	AssertCodesEvalToSameValue(t, `{|a,b|(1,4),(2,4)}`, `{|a|(1),(2)} <&> {|b|(4)}`)
	AssertCodesEvalToSameValue(t, `{|a,b|(1,4),(1,5)}`, `{|a|(1)} <&> {|b|(4),(5)}`)
	AssertCodesEvalToSameValue(t, `{|a,b|(1,4),(1,5),(2,4),(2,5)}`, `{|a|(1),(2)} <&> {|b|(4),(5)}`)

	AssertCodesEvalToSameValue(t, `{}`, `{|a,b|(1,4),(1,5)} <&> {|b,c|(6,7)}`)
	AssertCodesEvalToSameValue(t, `{|a,b,c|(2,5,7)}`, `{|a,b|(1,4),(2,5)} <&> {|b,c|(5,7),(6,7)}`)
}

func TestCompose(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{}`, `{} <-> {}`)
	AssertCodesEvalToSameValue(t, `{}`, `{()} <-> {}`)
	AssertCodesEvalToSameValue(t, `{}`, `{} <-> {()}`)
	AssertCodesEvalToSameValue(t, `{()}`, `{()} <-> {()}`)

	AssertCodesEvalToSameValue(t, `{}`, `{|a|(1)} <-> {|a|(2)}`)
	AssertCodesEvalToSameValue(t, `{()}`, `{|a|(1)} <-> {|a|(1)}`)

	AssertCodesEvalToSameValue(t, `{|a,b|(1,4)}`, `{|a|(1)} <-> {|b|(4)}`)
	AssertCodesEvalToSameValue(t, `{|a,b|(1,4),(2,4)}`, `{|a|(1),(2)} <-> {|b|(4)}`)
	AssertCodesEvalToSameValue(t, `{|a,b|(1,4),(1,5)}`, `{|a|(1)} <-> {|b|(4),(5)}`)
	AssertCodesEvalToSameValue(t, `{|a,b|(1,4),(1,5),(2,4),(2,5)}`, `{|a|(1),(2)} <-> {|b|(4),(5)}`)

	AssertCodesEvalToSameValue(t, `{}`, `{|a,b|(1,4),(1,5)} <-> {|b,c|(6,7)}`)
	AssertCodesEvalToSameValue(t, `{|a,c|(2,7)}`, `{|a,b|(1,4),(2,5)} <-> {|b,c|(5,7),(6,7)}`)
}

func TestJoinExists(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{}`, `{} --- {}`)
	AssertCodesEvalToSameValue(t, `{}`, `{()} --- {}`)
	AssertCodesEvalToSameValue(t, `{}`, `{} --- {()}`)
	AssertCodesEvalToSameValue(t, `{()}`, `{()} --- {()}`)

	AssertCodesEvalToSameValue(t, `{}`, `{|a|(1)} --- {|a|(2)}`)
	AssertCodesEvalToSameValue(t, `{()}`, `{|a|(1)} --- {|a|(1)}`)

	AssertCodesEvalToSameValue(t, `{()}`, `{|a|(1)} --- {|b|(4)}`)
	AssertCodesEvalToSameValue(t, `{()}`, `{|a|(1),(2)} --- {|b|(4)}`)
	AssertCodesEvalToSameValue(t, `{()}`, `{|a|(1)} --- {|b|(4),(5)}`)
	AssertCodesEvalToSameValue(t, `{()}`, `{|a|(1),(2)} --- {|b|(4),(5)}`)

	AssertCodesEvalToSameValue(t, `{}`, `{|a,b|(1,4),(1,5)} --- {|b,c|(6,7)}`)
	AssertCodesEvalToSameValue(t, `{()}`, `{|a,b|(1,4),(2,5)} --- {|b,c|(5,7),(6,7)}`)
}

func TestJoinCommon(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{}`, `{} -&- {}`)
	AssertCodesEvalToSameValue(t, `{}`, `{()} -&- {}`)
	AssertCodesEvalToSameValue(t, `{}`, `{} -&- {()}`)
	AssertCodesEvalToSameValue(t, `{()}`, `{()} -&- {()}`)

	AssertCodesEvalToSameValue(t, `{}`, `{|a|(1)} -&- {|a|(2)}`)
	AssertCodesEvalToSameValue(t, `{|a|(1)}`, `{|a|(1)} -&- {|a|(1)}`)

	AssertCodesEvalToSameValue(t, `{()}`, `{|a|(1)} -&- {|b|(4)}`)
	AssertCodesEvalToSameValue(t, `{()}`, `{|a|(1),(2)} -&- {|b|(4)}`)
	AssertCodesEvalToSameValue(t, `{()}`, `{|a|(1)} -&- {|b|(4),(5)}`)
	AssertCodesEvalToSameValue(t, `{()}`, `{|a|(1),(2)} -&- {|b|(4),(5)}`)

	AssertCodesEvalToSameValue(t, `{}`, `{|a,b|(1,4),(1,5)} -&- {|b,c|(6,7)}`)
	AssertCodesEvalToSameValue(t, `{|b|(5)}`, `{|a,b|(1,4),(2,5)} -&- {|b,c|(5,7),(6,7)}`)
}

func TestRightMatchJoin(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{}`, `{} -&> {}`)
	AssertCodesEvalToSameValue(t, `{}`, `{()} -&> {}`)
	AssertCodesEvalToSameValue(t, `{}`, `{} -&> {()}`)
	AssertCodesEvalToSameValue(t, `{()}`, `{()} -&> {()}`)

	AssertCodesEvalToSameValue(t, `{}`, `{|a|(1)} -&> {|a|(2)}`)
	AssertCodesEvalToSameValue(t, `{|a|(1)}`, `{|a|(1)} -&> {|a|(1)}`)

	AssertCodesEvalToSameValue(t, `{|b|(4)}`, `{|a|(1)} -&> {|b|(4)}`)
	AssertCodesEvalToSameValue(t, `{|b|(4),(4)}`, `{|a|(1),(2)} -&> {|b|(4)}`)
	AssertCodesEvalToSameValue(t, `{|b|(4),(5)}`, `{|a|(1)} -&> {|b|(4),(5)}`)
	AssertCodesEvalToSameValue(t, `{|b|(4),(5),(4),(5)}`, `{|a|(1),(2)} -&> {|b|(4),(5)}`)

	AssertCodesEvalToSameValue(t, `{}`, `{|a,b|(1,4),(1,5)} -&> {|b,c|(6,7)}`)
	AssertCodesEvalToSameValue(t, `{|b,c|(5,7)}`, `{|a,b|(1,4),(2,5)} -&> {|b,c|(5,7),(6,7)}`)
}

func TestLeftMatchJoin(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{}`, `{} <&- {}`)
	AssertCodesEvalToSameValue(t, `{}`, `{()} <&- {}`)
	AssertCodesEvalToSameValue(t, `{}`, `{} <&- {()}`)
	AssertCodesEvalToSameValue(t, `{()}`, `{()} <&- {()}`)

	AssertCodesEvalToSameValue(t, `{}`, `{|a|(1)} <&- {|a|(2)}`)
	AssertCodesEvalToSameValue(t, `{|a|(1)}`, `{|a|(1)} <&- {|a|(1)}`)

	AssertCodesEvalToSameValue(t, `{|a|(1)}`, `{|a|(1)} <&- {|b|(4)}`)
	AssertCodesEvalToSameValue(t, `{|a|(1),(2)}`, `{|a|(1),(2)} <&- {|b|(4)}`)
	AssertCodesEvalToSameValue(t, `{|a|(1),(1)}`, `{|a|(1)} <&- {|b|(4),(5)}`)
	AssertCodesEvalToSameValue(t, `{|a|(1),(1),(2),(2)}`, `{|a|(1),(2)} <&- {|b|(4),(5)}`)

	AssertCodesEvalToSameValue(t, `{}`, `{|a,b|(1,4),(1,5)} <&- {|b,c|(6,7)}`)
	AssertCodesEvalToSameValue(t, `{|a,b|(2,5)}`, `{|a,b|(1,4),(2,5)} <&- {|b,c|(5,7),(6,7)}`)
}

func TestRightResidue(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{}`, `{} --> {}`)
	AssertCodesEvalToSameValue(t, `{}`, `{()} --> {}`)
	AssertCodesEvalToSameValue(t, `{}`, `{} --> {()}`)
	AssertCodesEvalToSameValue(t, `{()}`, `{()} --> {()}`)

	AssertCodesEvalToSameValue(t, `{}`, `{|a|(1)} --> {|a|(2)}`)
	AssertCodesEvalToSameValue(t, `{()}`, `{|a|(1)} --> {|a|(1)}`)

	AssertCodesEvalToSameValue(t, `{|b|(4)}`, `{|a|(1)} --> {|b|(4)}`)
	AssertCodesEvalToSameValue(t, `{|b|(4)}`, `{|a|(1),(2)} --> {|b|(4)}`)
	AssertCodesEvalToSameValue(t, `{|b|(4),(5)}`, `{|a|(1)} --> {|b|(4),(5)}`)
	AssertCodesEvalToSameValue(t, `{|b|(4),(5)}`, `{|a|(1),(2)} --> {|b|(4),(5)}`)

	AssertCodesEvalToSameValue(t, `{}`, `{|a,b|(1,4),(1,5)} --> {|b,c|(6,7)}`)
	AssertCodesEvalToSameValue(t, `{|c|(7)}`, `{|a,b|(1,4),(2,5)} --> {|b,c|(5,7),(6,7)}`)
}

func TestLeftResidue(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{}`, `{} <-- {}`)
	AssertCodesEvalToSameValue(t, `{}`, `{()} <-- {}`)
	AssertCodesEvalToSameValue(t, `{}`, `{} <-- {()}`)
	AssertCodesEvalToSameValue(t, `{()}`, `{()} <-- {()}`)

	AssertCodesEvalToSameValue(t, `{}`, `{|a|(1)} <-- {|a|(2)}`)
	AssertCodesEvalToSameValue(t, `{()}`, `{|a|(1)} <-- {|a|(1)}`)

	AssertCodesEvalToSameValue(t, `{|a|(1)}`, `{|a|(1)} <-- {|b|(4)}`)
	AssertCodesEvalToSameValue(t, `{|a|(1),(2)}`, `{|a|(1),(2)} <-- {|b|(4)}`)
	AssertCodesEvalToSameValue(t, `{|a|(1)}`, `{|a|(1)} <-- {|b|(4),(5)}`)
	AssertCodesEvalToSameValue(t, `{|a|(1),(2)}`, `{|a|(1),(2)} <-- {|b|(4),(5)}`)

	AssertCodesEvalToSameValue(t, `{}`, `{|a,b|(1,4),(2,5)} <-- {|b,c|(6,7)}`)
	AssertCodesEvalToSameValue(t, `{|a|(2)}`, `{|a,b|(1,4),(2,5)} <-- {|b,c|(5,7),(6,7)}`)
}

func TestJoinIntoSpecialSet(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{0: 1} `, `{|@, @value, x| (0, 1, 2)}  -&- {|@, @value, y| (0, 1, 3) }`)
	AssertCodesEvalToSameValue(t, `'a'    `, `{|@, @char,  x| (0, 97, 2)} -&- {|@, @char, y | (0, 97, 3)}`)
	AssertCodesEvalToSameValue(t, `<<'a'>>`, `{|@, @byte,  x| (0, 97, 2)} -&- {|@, @byte, y | (0, 97, 3)}`)
	AssertCodesEvalToSameValue(t, `[1]    `, `{|@, @item,  x| (0, 1, 2)}  -&- {|@, @item, y | (0, 1, 3) }`)

	AssertCodesEvalToSameValue(t,
		`{(@: 0, @notspecial: 1)}`,
		`{|@, @notspecial, x| (0, 1, 2)}  -&- {|@, @notspecial, y| (0, 1, 3)}`,
	)
}

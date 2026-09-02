package rel

import (
	"fmt"

	"github.com/arr-ai/wbnf/parser"
)

func decodeExpr(n PlanNode) (Expr, error) {
	if h, ok := planLiftHooks[n.K]; ok {
		return h(n)
	}
	switch n.K {
	case "nil":
		return nil, nil
	case "ident":
		return NewIdentExpr(planSrc, n.Attr), nil
	case "dynident":
		return DynIdentExpr{IdentExpr: IdentExpr{ExprScanner: ExprScanner{Src: planSrc}, ident: n.Attr}}, nil
	case "lit":
		if len(n.Kids) != 1 {
			return nil, fmt.Errorf("plan: lit needs 1 kid")
		}
		return decodeValue(n.Kids[0])
	case "num", "str", "bytes", "none", "true", "arrayval", "setval", "tuple",
		"aitemval", "scharval", "bbyteval", "dentryval", "native", "hole":
		return decodeValue(n)
	case "bin":
		return decodeBin(n)
	case "unary":
		return decodeUnary(n)
	case "dot":
		lhs, err := kidExpr(n, 0)
		if err != nil {
			return nil, err
		}
		return NewDotExpr(planSrc, lhs, n.Attr), nil
	case "fn":
		return decodeFn(n)
	case "arrow":
		return decodeBin(n)
	case "tupexpr":
		return decodeTupleExpr(n)
	case "array":
		es, err := kidExprs(n)
		if err != nil {
			return nil, err
		}
		return NewArrayExpr(planSrc, es...), nil
	case "set":
		es, err := kidExprs(n)
		if err != nil {
			return nil, err
		}
		return NewSetExpr(planSrc, es...)
	case "dictexpr":
		return decodeDictExpr(n)
	case "cmp":
		return decodeCompare(n)
	case "if":
		es, err := kidExprs(n)
		if err != nil {
			return nil, err
		}
		if len(es) != 3 {
			return nil, fmt.Errorf("plan: if needs 3 kids")
		}
		return NewIfElseExpr(planSrc, es[0], es[1], es[2]), nil
	case "cond":
		d, err := kidExpr(n, 0)
		if err != nil {
			return nil, err
		}
		return NewCondExpr(planSrc, d), nil
	case "condpat":
		return decodeCondPat(n)
	case "nest":
		lhs, err := kidExpr(n, 0)
		if err != nil {
			return nil, err
		}
		return NewNestExpr(planSrc, n.Flag, lhs, NewNames(n.Strs...), n.Attr), nil
	case "unnest":
		lhs, err := kidExpr(n, 0)
		if err != nil {
			return nil, err
		}
		return NewUnnestExpr(planSrc, lhs, n.Attr), nil
	case "snest":
		lhs, err := kidExpr(n, 0)
		if err != nil {
			return nil, err
		}
		return NewSingleNestExpr(planSrc, lhs, n.Attr), nil
	case "offset":
		es, err := kidExprs(n)
		if err != nil {
			return nil, err
		}
		if len(es) != 2 {
			return nil, fmt.Errorf("plan: offset needs 2 kids")
		}
		return NewOffsetExpr(planSrc, es[0], es[1]), nil
	case "rec":
		fn, err := kidExpr(n, 0)
		if err != nil {
			return nil, err
		}
		return NewRecursionExpr(planSrc, n.Attr, fn), nil
	case "dynlet":
		es, err := kidExprs(n)
		if err != nil {
			return nil, err
		}
		if len(es) != 2 {
			return nil, fmt.Errorf("plan: dynlet needs 2 kids")
		}
		return NewDynLetExpr(planSrc, es[0], es[1]), nil
	case "project":
		base, err := kidExpr(n, 0)
		if err != nil {
			return nil, err
		}
		return NewTupleProjectExpr(planSrc, base, n.Flag, n.Strs), nil
	case "exprexpr":
		inner, err := kidExpr(n, 0)
		if err != nil {
			return nil, err
		}
		return NewExprExpr(planSrc, inner), nil
	case "safetail":
		return decodeSafeTail(n)
	case "bytesexpr":
		es, err := kidExprs(n)
		if err != nil {
			return nil, err
		}
		return NewBytesExpr(planSrc, es...), nil
	case "schar":
		es, err := kidExprs(n)
		if err != nil {
			return nil, err
		}
		if len(es) != 2 {
			return nil, fmt.Errorf("plan: schar needs 2 kids")
		}
		return NewStringCharTupleExpr(planSrc, es[0], es[1]), nil
	case "aitem":
		es, err := kidExprs(n)
		if err != nil {
			return nil, err
		}
		if len(es) != 2 {
			return nil, fmt.Errorf("plan: aitem needs 2 kids")
		}
		return NewArrayItemTupleExpr(planSrc, es[0], es[1]), nil
	case "dentry":
		es, err := kidExprs(n)
		if err != nil {
			return nil, err
		}
		if len(es) != 2 {
			return nil, fmt.Errorf("plan: dentry needs 2 kids")
		}
		return NewDictEntryTupleExpr(planSrc, es[0], es[1]), nil
	default:
		if stringsHasPrefixP(n.K) {
			p, err := decodePattern(n)
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("plan: pattern %T used as expr", p)
		}
		return nil, fmt.Errorf("plan: cannot lift kind %q", n.K)
	}
}

func stringsHasPrefixP(k string) bool {
	return len(k) >= 2 && k[0] == 'p' && k[1] == '-'
}

func kidExpr(n PlanNode, i int) (Expr, error) {
	if i >= len(n.Kids) {
		return nil, fmt.Errorf("plan: %s missing kid %d", n.K, i)
	}
	return decodeExpr(n.Kids[i])
}

func kidExprs(n PlanNode) ([]Expr, error) {
	es := make([]Expr, len(n.Kids))
	for i := range n.Kids {
		e, err := decodeExpr(n.Kids[i])
		if err != nil {
			return nil, err
		}
		es[i] = e
	}
	return es, nil
}

func decodeBin(n PlanNode) (Expr, error) {
	es, err := kidExprs(n)
	if err != nil {
		return nil, err
	}
	if len(es) != 2 {
		return nil, fmt.Errorf("plan: bin %s needs 2 kids", n.Op)
	}
	f := binCtors[n.Op]
	if f == nil {
		return nil, fmt.Errorf("plan: unknown bin op %q", n.Op)
	}
	return f(planSrc, es[0], es[1]), nil
}

func decodeUnary(n PlanNode) (Expr, error) {
	a, err := kidExpr(n, 0)
	if err != nil {
		return nil, err
	}
	switch n.Op {
	case "+":
		return NewPosExpr(planSrc, a), nil
	case "-":
		return NewNegExpr(planSrc, a), nil
	case "**":
		return NewPowerSetExpr(planSrc, a), nil
	case "!":
		return NewNotExpr(planSrc, a), nil
	case "*":
		return NewEvalExpr(planSrc, a), nil
	case "count":
		return NewCountExpr(planSrc, a), nil
	case "single":
		return NewSingleExpr(planSrc, a), nil
	default:
		return nil, fmt.Errorf("plan: unknown unary op %q", n.Op)
	}
}

var binCtors = map[string]func(parser.Scanner, Expr, Expr) Expr{
	"->":      NewArrowExpr,
	"=>":      NewDArrowExpr,
	">>":      NewSeqArrowExpr(false),
	">>>":     NewSeqArrowExpr(true),
	":>":      NewTupleMapExpr,
	"orderby": NewOrderByExpr,
	"order":   NewOrderExpr,
	"rank":    NewRankExpr,
	"where":   NewWhereExpr,
	"sum":     NewSumExpr,
	"max":     NewMaxExpr,
	"mean":    NewMeanExpr,
	"median":  NewMedianExpr,
	"min":     NewMinExpr,
	"with":    NewWithExpr,
	"without": NewWithoutExpr,
	"&&":      NewAndExpr,
	"||":      NewOrExpr,
	"+":       NewAddExpr,
	"-":       NewSubExpr,
	"++":      NewConcatExpr,
	"&~":      NewDiffExpr,
	"~~":      NewSymmDiffExpr,
	"&":       NewIntersectExpr,
	"|":       NewUnionExpr,
	"<&>":     NewJoinExpr,
	"<->":     NewComposeExpr,
	"-&-":     NewJoinCommonExpr,
	"---":     NewJoinExistsExpr,
	"-&>":     NewRightMatchExpr,
	"<&-":     NewLeftMatchExpr,
	"-->":     NewRightResidueExpr,
	"<--":     NewLeftResidueExpr,
	"*":       NewMulExpr,
	"/":       NewDivExpr,
	"%":       NewModExpr,
	"-%":      NewSubModExpr,
	"//":      NewIdivExpr,
	"^":       NewPowExpr,
	"\\":      NewOffsetExpr,
	"+>":      NewAddArrowExpr,
	"call":    NewCallExpr,
}

func decodeFn(n PlanNode) (Expr, error) {
	if len(n.Kids) != 2 {
		return nil, fmt.Errorf("plan: fn needs 2 kids")
	}
	p, err := decodePattern(n.Kids[0])
	if err != nil {
		return nil, err
	}
	body, err := decodeExpr(n.Kids[1])
	if err != nil {
		return nil, err
	}
	return NewFunction(planSrc, p, body), nil
}

func decodeTupleExpr(n PlanNode) (Expr, error) {
	attrs := make([]AttrExpr, len(n.Kids))
	for i, k := range n.Kids {
		if k.K != "attr" || len(k.Kids) != 1 {
			return nil, fmt.Errorf("plan: tupexpr kid %d not attr", i)
		}
		e, err := decodeExpr(k.Kids[0])
		if err != nil {
			return nil, err
		}
		a, err := NewAttrExpr(planSrc, k.Attr, e)
		if err != nil {
			return nil, err
		}
		attrs[i] = a
	}
	return NewTupleExpr(planSrc, attrs...), nil
}

func decodeDictExpr(n PlanNode) (Expr, error) {
	if len(n.Kids)%2 != 0 {
		return nil, fmt.Errorf("plan: dictexpr odd kids")
	}
	ents := make([]DictEntryTupleExpr, 0, len(n.Kids)/2)
	for i := 0; i < len(n.Kids); i += 2 {
		at, err := decodeExpr(n.Kids[i])
		if err != nil {
			return nil, err
		}
		val, err := decodeExpr(n.Kids[i+1])
		if err != nil {
			return nil, err
		}
		ents = append(ents, NewDictEntryTupleExpr(planSrc, at, val))
	}
	return NewDictExpr(planSrc, n.Flag, true, ents...)
}

func decodeCompare(n PlanNode) (Expr, error) {
	args, err := kidExprs(n)
	if err != nil {
		return nil, err
	}
	if len(n.Strs) != len(args)-1 {
		return nil, fmt.Errorf("plan: cmp ops/args mismatch")
	}
	comps := make([]CompareFunc, len(n.Strs))
	for i, op := range n.Strs {
		f := compareCtor(op)
		if f == nil {
			return nil, fmt.Errorf("plan: unknown compare op %q", op)
		}
		comps[i] = f
	}
	return NewCompareExpr(planSrc, args, comps, n.Strs), nil
}

func decodeCondPat(n PlanNode) (Expr, error) {
	if len(n.Kids) < 1 || (len(n.Kids)-1)%2 != 0 {
		return nil, fmt.Errorf("plan: condpat shape")
	}
	ctrl, err := decodeExpr(n.Kids[0])
	if err != nil {
		return nil, err
	}
	pairs := make([]PatternExprPair, 0, (len(n.Kids)-1)/2)
	for i := 1; i < len(n.Kids); i += 2 {
		p, err := decodePattern(n.Kids[i])
		if err != nil {
			return nil, err
		}
		x, err := decodeExpr(n.Kids[i+1])
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, NewPatternExprPair(p, x))
	}
	return NewCondPatternControlVarExpr(planSrc, ctrl, pairs...), nil
}

func decodeSafeTail(n PlanNode) (Expr, error) {
	if len(n.Kids) < 2 {
		return nil, fmt.Errorf("plan: safetail needs fallback+base")
	}
	fb, err := decodeExpr(n.Kids[0])
	if err != nil {
		return nil, err
	}
	base, err := decodeExpr(n.Kids[1])
	if err != nil {
		return nil, err
	}
	steps := make([]SafeTailStep, 0, len(n.Kids)-2)
	for _, k := range n.Kids[2:] {
		args := make([]Expr, len(k.Kids))
		for i, a := range k.Kids {
			e, err := decodeExpr(a)
			if err != nil {
				return nil, err
			}
			args[i] = e
		}
		if k.Op == "get" {
			steps = append(steps, NewSafeTailGet(k.Flag, k.Attr))
		} else {
			steps = append(steps, NewSafeTailCall(k.Flag, args...))
		}
	}
	return NewSafeTailExpr(planSrc, fb, base, steps), nil
}

func decodePattern(n PlanNode) (Pattern, error) {
	switch n.K {
	case "p-nil":
		return nil, nil
	case "p-ident":
		return NewIdentPattern(n.Attr), nil
	case "p-dyn":
		return DynIdentPattern(n.Attr), nil
	case "p-extra":
		return NewExtraElementPattern(n.Attr), nil
	case "p-expr":
		e, err := kidExpr(n, 0)
		if err != nil {
			return nil, err
		}
		return NewExprPattern(e), nil
	case "p-exprs":
		es, err := kidExprs(n)
		if err != nil {
			return nil, err
		}
		return NewExprsPattern(es...), nil
	case "p-array":
		items := make([]FallbackPattern, len(n.Kids))
		for i, k := range n.Kids {
			p, err := decodeFallback(k)
			if err != nil {
				return nil, err
			}
			items[i] = p
		}
		return NewArrayPattern(items...), nil
	case "p-tuple":
		attrs := make([]TuplePatternAttr, len(n.Kids))
		for i, k := range n.Kids {
			p, err := decodeFallback(k)
			if err != nil {
				return nil, err
			}
			attrs[i] = NewTuplePatternAttr(k.Attr, p)
		}
		return NewTuplePattern(attrs...)
	case "p-dict":
		if len(n.Kids)%2 != 0 {
			return nil, fmt.Errorf("plan: p-dict odd kids")
		}
		ents := make([]DictPatternEntry, 0, len(n.Kids)/2)
		for i := 0; i < len(n.Kids); i += 2 {
			at, err := decodeExpr(n.Kids[i])
			if err != nil {
				return nil, err
			}
			p, err := decodeFallback(n.Kids[i+1])
			if err != nil {
				return nil, err
			}
			ents = append(ents, NewDictPatternEntry(at, p))
		}
		return NewDictPattern(ents...), nil
	case "p-set":
		ps := make([]Pattern, len(n.Kids))
		for i, k := range n.Kids {
			p, err := decodePattern(k)
			if err != nil {
				return nil, err
			}
			ps[i] = p
		}
		return NewSetPattern(ps...), nil
	default:
		return nil, fmt.Errorf("plan: cannot lift pattern kind %q", n.K)
	}
}

func decodeFallback(n PlanNode) (FallbackPattern, error) {
	if n.K != "p-fb" {
		p, err := decodePattern(n)
		if err != nil {
			return FallbackPattern{}, err
		}
		return NewFallbackPattern(p, nil), nil
	}
	if len(n.Kids) < 1 {
		return FallbackPattern{}, fmt.Errorf("plan: p-fb empty")
	}
	p, err := decodePattern(n.Kids[0])
	if err != nil {
		return FallbackPattern{}, err
	}
	var fb Expr
	if n.Flag && len(n.Kids) > 1 {
		fb, err = decodeExpr(n.Kids[1])
		if err != nil {
			return FallbackPattern{}, err
		}
	}
	return NewFallbackPattern(p, fb), nil
}

func decodeValue(n PlanNode) (Value, error) {
	switch n.K {
	case "num":
		return NewNumber(n.Num), nil
	case "str":
		offset := int(n.Num)
		if n.Bytes != nil {
			runes := make([]rune, len(n.Bytes))
			for i, b := range n.Bytes {
				runes[i] = rune(b)
			}
			s := NewOffsetString(runes, offset)
			if v, ok := s.(Value); ok {
				return v, nil
			}
			return nil, fmt.Errorf("plan: empty string became None")
		}
		s := NewOffsetString([]rune(n.Str), offset)
		if v, ok := s.(Value); ok {
			return v, nil
		}
		return None, nil
	case "bytes":
		s := NewOffsetBytes(n.Bytes, int(n.Num))
		if v, ok := s.(Value); ok {
			return v, nil
		}
		return None, nil
	case "none", "hole":
		return None, nil
	case "true":
		return True, nil
	case "arrayval":
		vals := make([]Value, len(n.Kids))
		for i, k := range n.Kids {
			if k.K == "hole" {
				continue
			}
			v, err := decodeValue(k)
			if err != nil {
				return nil, err
			}
			vals[i] = v
		}
		s := NewOffsetArray(int(n.Num), vals...)
		if v, ok := s.(Value); ok {
			return v, nil
		}
		return None, nil
	case "setval":
		b := NewSetBuilder()
		for _, k := range n.Kids {
			v, err := decodeValue(k)
			if err != nil {
				return nil, err
			}
			b.Add(v)
		}
		return b.Finish()
	case "tuple":
		attrs := make([]Attr, len(n.Kids))
		for i, k := range n.Kids {
			v, err := decodeValue(k)
			if err != nil {
				return nil, err
			}
			attrs[i] = NewAttr(k.Attr, v)
		}
		return NewTuple(attrs...), nil
	case "aitemval":
		if len(n.Kids) != 1 {
			return nil, fmt.Errorf("plan: aitemval")
		}
		item, err := decodeValue(n.Kids[0])
		if err != nil {
			return nil, err
		}
		return NewArrayItemTuple(int(n.Num), item), nil
	case "scharval":
		var r rune
		if n.Attr != "" {
			rs := []rune(n.Attr)
			if len(rs) > 0 {
				r = rs[0]
			}
		}
		return NewStringCharTuple(int(n.Num), r), nil
	case "bbyteval":
		var b byte
		if len(n.Bytes) > 0 {
			b = n.Bytes[0]
		}
		return NewBytesByteTuple(int(n.Num), b), nil
	case "dentryval":
		if len(n.Kids) != 2 {
			return nil, fmt.Errorf("plan: dentryval")
		}
		at, err := decodeValue(n.Kids[0])
		if err != nil {
			return nil, err
		}
		val, err := decodeValue(n.Kids[1])
		if err != nil {
			return nil, err
		}
		return NewDictEntryTuple(at, val), nil
	case "native":
		f := lookupNative(n.Op)
		if f == nil {
			return nil, fmt.Errorf("plan: unknown native %q", n.Op)
		}
		return f, nil
	default:
		e, err := decodeExpr(n)
		if err != nil {
			return nil, err
		}
		if v, ok := e.(Value); ok {
			return v, nil
		}
		return nil, fmt.Errorf("plan: kind %q is not a value", n.K)
	}
}

func compareCtor(op string) CompareFunc {
	switch op {
	case "<:":
		return func(a, b Value) (bool, error) {
			set, is := b.(Set)
			if !is {
				return false, fmt.Errorf("<: rhs not a set: %v", b)
			}
			return set.Has(a), nil
		}
	case "!<:":
		return func(a, b Value) (bool, error) {
			set, is := b.(Set)
			if !is {
				return false, fmt.Errorf("!<: rhs not a set: %v", b)
			}
			return !set.Has(a), nil
		}
	case "=":
		return func(a, b Value) (bool, error) { return a.Equal(b), nil }
	case "!=":
		return func(a, b Value) (bool, error) { return !a.Equal(b), nil }
	case "<":
		return func(a, b Value) (bool, error) { return a.Less(b), nil }
	case ">":
		return func(a, b Value) (bool, error) { return b.Less(a), nil }
	case "<=":
		return func(a, b Value) (bool, error) { return !b.Less(a), nil }
	case ">=":
		return func(a, b Value) (bool, error) { return !a.Less(b), nil }
	default:
		return nil
	}
}

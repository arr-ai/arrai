package rel

import (
	"fmt"
	"strings"
)

func encodeExpr(e Expr) (PlanNode, error) {
	if e == nil {
		return PlanNode{K: "nil"}, nil
	}
	for _, h := range planLowerHooks {
		n, ok, err := h(e)
		if err != nil {
			return PlanNode{}, err
		}
		if ok {
			return n, nil
		}
	}
	if v, ok := e.(Value); ok {
		switch v.(type) {
		case Closure, ExprClosure:
			// Fall through: these are Exprs that happen to be Values.
		default:
			return encodeValue(v)
		}
	}
	switch e := e.(type) {
	case IdentExpr:
		return PlanNode{K: "ident", Attr: e.ident}, nil
	case DynIdentExpr:
		return PlanNode{K: "dynident", Attr: e.ident}, nil
	case LiteralExpr:
		n, err := encodeValue(e.literal)
		if err != nil {
			return PlanNode{}, err
		}
		return node("lit", n), nil
	case *BinExpr:
		return encodeBin("bin", e.op, e.a, e.b)
	case *UnaryExpr:
		a, err := encodeExpr(e.a)
		if err != nil {
			return PlanNode{}, err
		}
		return nodeOp("unary", e.op, a), nil
	case *DotExpr:
		lhs, err := encodeExpr(e.lhs)
		if err != nil {
			return PlanNode{}, err
		}
		return PlanNode{K: "dot", Attr: e.attr, Kids: []PlanNode{lhs}}, nil
	case *Function:
		return encodeFn(e)
	case *ArrowExpr:
		return encodeBin("arrow", "->", e.lhs, e.fn)
	case *DArrowExpr:
		return encodeBin("bin", "=>", e.lhs, e.fn)
	case *SeqArrowExpr:
		op := ">>"
		if e.withAt {
			op = ">>>"
		}
		return encodeBin("bin", op, e.lhs, e.fn)
	case *TupleMapExpr:
		return encodeBin("bin", ":>", e.lhs, e.fn)
	case *TupleExpr:
		return encodeTupleExpr(e)
	case ArrayExpr:
		return encodeExprList("array", e.elements)
	case *SetExpr:
		return encodeExprList("set", e.elements)
	case DictExpr:
		return encodeDictExpr(e)
	case *DictExpr:
		return encodeDictExpr(*e)
	case CompareExpr:
		return encodeCompare(e)
	case AndExpr:
		return encodeBin("bin", "&&", e.a, e.b)
	case OrExpr:
		return encodeBin("bin", "||", e.a, e.b)
	case *IfElseExpr:
		return encodeNary("if", e.ifTrue, e.cond, e.ifFalse)
	case CondExpr:
		d, err := encodeExpr(e.dicExpr)
		if err != nil {
			return PlanNode{}, err
		}
		return node("cond", d), nil
	case CondPatternControlVarExpr:
		return encodeCondPat(e)
	case *NestExpr:
		return encodeNest(e)
	case *UnnestExpr:
		lhs, err := encodeExpr(e.lhs)
		if err != nil {
			return PlanNode{}, err
		}
		return PlanNode{K: "unnest", Attr: e.attr, Kids: []PlanNode{lhs}}, nil
	case *SingleNestExpr:
		lhs, err := encodeExpr(e.lhs)
		if err != nil {
			return PlanNode{}, err
		}
		return PlanNode{K: "snest", Attr: e.attr, Kids: []PlanNode{lhs}}, nil
	case *OffsetExpr:
		return encodeBin("offset", "\\", e.offset, e.array)
	case RecursionExpr:
		fn, err := encodeExpr(e.fn)
		if err != nil {
			return PlanNode{}, err
		}
		return PlanNode{K: "rec", Attr: e.name, Kids: []PlanNode{fn}}, nil
	case *RecursionExpr:
		fn, err := encodeExpr(e.fn)
		if err != nil {
			return PlanNode{}, err
		}
		return PlanNode{K: "rec", Attr: e.name, Kids: []PlanNode{fn}}, nil
	case *DynLetExpr:
		return encodeBin("dynlet", "", e.bindings, e.expr)
	case *TupleProjectExpr:
		base, err := encodeExpr(e.base)
		if err != nil {
			return PlanNode{}, err
		}
		return PlanNode{K: "project", Flag: e.inverse, Strs: e.attrs.OrderedNames(), Kids: []PlanNode{base}}, nil
	case ExprExpr:
		inner, err := encodeExpr(e.Expr)
		if err != nil {
			return PlanNode{}, err
		}
		return node("exprexpr", inner), nil
	case *SafeTailExpr:
		return encodeSafeTail(e)
	case *ReduceExpr:
		op := reduceOp(e.format)
		return encodeBin("bin", op, e.a, e.f)
	case BytesExpr:
		return encodeExprList("bytesexpr", e.elements)
	case StringCharTupleExpr:
		return encodeBin("schar", "", e.at, e.char)
	case ArrayItemTupleExpr:
		return encodeBin("aitem", "", e.at, e.item)
	case DictEntryTupleExpr:
		return encodeBin("dentry", "", e.at, e.value)
	default:
		return PlanNode{}, fmt.Errorf("plan: cannot lower %T", e)
	}
}

func encodeBin(k, op string, a, b Expr) (PlanNode, error) {
	na, err := encodeExpr(a)
	if err != nil {
		return PlanNode{}, err
	}
	nb, err := encodeExpr(b)
	if err != nil {
		return PlanNode{}, err
	}
	return nodeOp(k, op, na, nb), nil
}

func encodeNary(k string, es ...Expr) (PlanNode, error) {
	return encodeExprList(k, es)
}

func encodeExprList(k string, es []Expr) (PlanNode, error) {
	kids := make([]PlanNode, len(es))
	for i, e := range es {
		n, err := encodeExpr(e)
		if err != nil {
			return PlanNode{}, err
		}
		kids[i] = n
	}
	return PlanNode{K: k, Kids: kids}, nil
}

func encodeFn(f *Function) (PlanNode, error) {
	p, err := encodePattern(f.arg)
	if err != nil {
		return PlanNode{}, err
	}
	body, err := encodeExpr(f.body)
	if err != nil {
		return PlanNode{}, err
	}
	return node("fn", p, body), nil
}

func encodeTupleExpr(e *TupleExpr) (PlanNode, error) {
	kids := make([]PlanNode, len(e.attrs))
	for i, a := range e.attrs {
		n, err := encodeExpr(a.expr)
		if err != nil {
			return PlanNode{}, err
		}
		kids[i] = PlanNode{K: "attr", Attr: a.name, Kids: []PlanNode{n}}
	}
	return PlanNode{K: "tupexpr", Kids: kids}, nil
}

func encodeDictExpr(e DictExpr) (PlanNode, error) {
	kids := make([]PlanNode, 0, len(e.entryExprs)*2)
	for _, ent := range e.entryExprs {
		at, err := encodeExpr(ent.at)
		if err != nil {
			return PlanNode{}, err
		}
		val, err := encodeExpr(ent.value)
		if err != nil {
			return PlanNode{}, err
		}
		kids = append(kids, at, val)
	}
	return PlanNode{K: "dictexpr", Flag: e.allowDupKeys, Kids: kids}, nil
}

func encodeCompare(e CompareExpr) (PlanNode, error) {
	kids := make([]PlanNode, len(e.args))
	for i, a := range e.args {
		n, err := encodeExpr(a)
		if err != nil {
			return PlanNode{}, err
		}
		kids[i] = n
	}
	return PlanNode{K: "cmp", Strs: e.ops, Kids: kids}, nil
}

func encodeCondPat(e CondPatternControlVarExpr) (PlanNode, error) {
	ctrl, err := encodeExpr(e.controlVarExpr)
	if err != nil {
		return PlanNode{}, err
	}
	kids := []PlanNode{ctrl}
	for _, pair := range e.conditionPairs {
		p, err := encodePattern(pair.pattern)
		if err != nil {
			return PlanNode{}, err
		}
		x, err := encodeExpr(pair.expr)
		if err != nil {
			return PlanNode{}, err
		}
		kids = append(kids, p, x)
	}
	return PlanNode{K: "condpat", Kids: kids}, nil
}

func encodeNest(e *NestExpr) (PlanNode, error) {
	lhs, err := encodeExpr(e.lhs)
	if err != nil {
		return PlanNode{}, err
	}
	return PlanNode{
		K:    "nest",
		Flag: e.inverse,
		Attr: e.attr,
		Strs: e.attrs.OrderedNames(),
		Kids: []PlanNode{lhs},
	}, nil
}

func encodeSafeTail(e *SafeTailExpr) (PlanNode, error) {
	fb, err := encodeExpr(e.fallbackValue)
	if err != nil {
		return PlanNode{}, err
	}
	base, err := encodeExpr(e.base)
	if err != nil {
		return PlanNode{}, err
	}
	kids := []PlanNode{fb, base}
	for _, s := range e.steps {
		step := PlanNode{K: "tail", Flag: s.safe, Attr: s.attr, Op: ""}
		if s.get {
			step.Op = "get"
		}
		step.Kids = make([]PlanNode, len(s.args))
		for i, a := range s.args {
			n, err := encodeExpr(a)
			if err != nil {
				return PlanNode{}, err
			}
			step.Kids[i] = n
		}
		kids = append(kids, step)
	}
	return PlanNode{K: "safetail", Kids: kids}, nil
}

func reduceOp(format string) string {
	for _, op := range []string{"sum", "max", "mean", "median", "min"} {
		if strings.Contains(format, " "+op+" ") {
			return op
		}
	}
	return "sum"
}

func encodePattern(p Pattern) (PlanNode, error) {
	if p == nil {
		return PlanNode{K: "p-nil"}, nil
	}
	switch p := p.(type) {
	case IdentPattern:
		return PlanNode{K: "p-ident", Attr: string(p)}, nil
	case DynIdentPattern:
		return PlanNode{K: "p-dyn", Attr: string(p)}, nil
	case ExtraElementPattern:
		return PlanNode{K: "p-extra", Attr: p.ident}, nil
	case ExprPattern:
		n, err := encodeExpr(p.Expr)
		if err != nil {
			return PlanNode{}, err
		}
		return node("p-expr", n), nil
	case ExprsPattern:
		return encodeExprList("p-exprs", p.exprs)
	case ArrayPattern:
		kids := make([]PlanNode, len(p.items))
		for i, item := range p.items {
			n, err := encodeFallback(item)
			if err != nil {
				return PlanNode{}, err
			}
			kids[i] = n
		}
		return PlanNode{K: "p-array", Kids: kids}, nil
	case TuplePattern:
		kids := make([]PlanNode, len(p.attrs))
		for i, a := range p.attrs {
			n, err := encodeFallback(a.pattern)
			if err != nil {
				return PlanNode{}, err
			}
			n.Attr = a.name
			kids[i] = n
		}
		return PlanNode{K: "p-tuple", Kids: kids}, nil
	case DictPattern:
		kids := make([]PlanNode, 0, len(p.entries)*2)
		for _, ent := range p.entries {
			at, err := encodeExpr(ent.at)
			if err != nil {
				return PlanNode{}, err
			}
			pat, err := encodeFallback(ent.pattern)
			if err != nil {
				return PlanNode{}, err
			}
			kids = append(kids, at, pat)
		}
		return PlanNode{K: "p-dict", Kids: kids}, nil
	case SetPattern:
		kids := make([]PlanNode, len(p.patterns))
		for i, sp := range p.patterns {
			n, err := encodePattern(sp)
			if err != nil {
				return PlanNode{}, err
			}
			kids[i] = n
		}
		return PlanNode{K: "p-set", Kids: kids}, nil
	default:
		return PlanNode{}, fmt.Errorf("plan: cannot lower pattern %T", p)
	}
}

func encodeFallback(p FallbackPattern) (PlanNode, error) {
	pat, err := encodePattern(p.pattern)
	if err != nil {
		return PlanNode{}, err
	}
	n := PlanNode{K: "p-fb", Kids: []PlanNode{pat}}
	if p.fallback != nil {
		fb, err := encodeExpr(p.fallback)
		if err != nil {
			return PlanNode{}, err
		}
		n.Flag = true
		n.Kids = append(n.Kids, fb)
	}
	return n, nil
}

func encodeValue(v Value) (PlanNode, error) {
	if v == nil {
		return PlanNode{K: "nil"}, nil
	}
	switch v := v.(type) {
	case Number:
		return PlanNode{K: "num", Num: v.Float64()}, nil
	case String:
		if v.ascii != nil {
			return PlanNode{K: "str", Bytes: append([]byte(nil), v.ascii...), Num: float64(v.offset)}, nil
		}
		return PlanNode{K: "str", Str: string(v.s), Num: float64(v.offset)}, nil
	case Bytes:
		return PlanNode{K: "bytes", Bytes: append([]byte(nil), v.b...), Num: float64(v.offset)}, nil
	case EmptySet:
		return PlanNode{K: "none"}, nil
	case TrueSet:
		return PlanNode{K: "true"}, nil
	case Array:
		kids := make([]PlanNode, len(v.values))
		for i, x := range v.values {
			if x == nil {
				kids[i] = PlanNode{K: "hole"}
				continue
			}
			n, err := encodeValue(x)
			if err != nil {
				return PlanNode{}, err
			}
			kids[i] = n
		}
		return PlanNode{K: "arrayval", Num: float64(v.offset), Kids: kids}, nil
	case Dict:
		return encodeSetEnum(v)
	case GenericSet:
		return encodeSetEnum(v)
	case UnionSet:
		return encodeSetEnum(v)
	case Relation:
		return encodeSetEnum(v)
	case *GenericTuple:
		return encodeTuple(v)
	case ArrayItemTuple:
		item, err := encodeValue(v.item)
		if err != nil {
			return PlanNode{}, err
		}
		return PlanNode{K: "aitemval", Num: float64(v.at), Kids: []PlanNode{item}}, nil
	case StringCharTuple:
		return PlanNode{K: "scharval", Num: float64(v.at), Attr: string(v.char)}, nil
	case BytesByteTuple:
		return PlanNode{K: "bbyteval", Num: float64(v.at), Bytes: []byte{v.byteval}}, nil
	case DictEntryTuple:
		at, err := encodeValue(v.at)
		if err != nil {
			return PlanNode{}, err
		}
		val, err := encodeValue(v.value)
		if err != nil {
			return PlanNode{}, err
		}
		return node("dentryval", at, val), nil
	case *NativeFunction:
		return PlanNode{K: "native", Op: v.name}, nil
	case seqPipeline:
		forced, err := v.force()
		if err != nil {
			return PlanNode{}, err
		}
		return encodeValue(forced)
	case dictPipeline:
		forced, err := v.force()
		if err != nil {
			return PlanNode{}, err
		}
		return encodeValue(forced)
	default:
		if t, ok := v.(Tuple); ok {
			return encodeTuple(t)
		}
		if s, ok := v.(Set); ok {
			return encodeSetEnum(s)
		}
		return PlanNode{}, fmt.Errorf("plan: cannot encode value %T", v)
	}
}

func encodeTuple(t Tuple) (PlanNode, error) {
	kids := make([]PlanNode, 0, t.Count())
	for e := t.Enumerator(); e.MoveNext(); {
		name, val := e.Current()
		n, err := encodeValue(val)
		if err != nil {
			return PlanNode{}, err
		}
		n.Attr = name
		kids = append(kids, n)
	}
	return PlanNode{K: "tuple", Kids: kids}, nil
}

func encodeSetEnum(s Set) (PlanNode, error) {
	if !s.IsTrue() {
		return PlanNode{K: "none"}, nil
	}
	kids := make([]PlanNode, 0, s.Count())
	for e := s.Enumerator(); e.MoveNext(); {
		n, err := encodeValue(e.Current())
		if err != nil {
			return PlanNode{}, err
		}
		kids = append(kids, n)
	}
	return PlanNode{K: "setval", Kids: kids}, nil
}

package rel

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
)

const (
	planMagic   = "ARPL"
	planVersion = byte(1)
)

// Plan is the S6 dataflow/operator graph: the same artifact that is serialized
// as compiled .arraiz and executed. Nodes are operators; closures are nested
// graphs; kernels rebind from op tags on lift (🎯T25).
type Plan struct {
	Root PlanNode
}

// PlanNode is a serializable operator. K is the node kind; Op/Attr/Num/Str/Bytes
// carry scalars; Kids are child operators.
type PlanNode struct {
	K     string
	Op    string
	Attr  string
	Num   float64
	Str   string
	Bytes []byte
	Flag  bool
	Strs  []string
	Kids  []PlanNode
}

var planSrc = *parser.NewScanner("")

// PlanScanner is an empty source location for reconstituted plan nodes.
func PlanScanner() parser.Scanner { return planSrc }

// LowerHook lets other packages (syntax) lower Expr types rel cannot see.
type LowerHook func(Expr) (PlanNode, bool, error)

// LiftHook reconstitutes a syntax-level Expr from a kind-tagged node.
type LiftHook func(PlanNode) (Expr, error)

var (
	planLowerHooks []LowerHook
	planLiftHooks  = map[string]LiftHook{}
)

func RegisterPlanLower(h LowerHook) {
	planLowerHooks = append(planLowerHooks, h)
}

func RegisterPlanLift(kind string, h LiftHook) {
	planLiftHooks[kind] = h
}

func LowerPlan(e Expr) (*Plan, error) {
	n, err := encodeExpr(e)
	if err != nil {
		return nil, err
	}
	return &Plan{Root: n}, nil
}

func LiftPlan(p *Plan) (Expr, error) {
	if p == nil {
		return nil, fmt.Errorf("plan: nil")
	}
	return decodeExpr(p.Root)
}

func (p *Plan) Eval(ctx context.Context, local Scope) (Value, error) {
	e, err := LiftPlan(p)
	if err != nil {
		return nil, err
	}
	return e.Eval(ctx, local)
}

func EncodePlan(p *Plan) ([]byte, error) {
	if p == nil {
		return nil, fmt.Errorf("plan: nil")
	}
	var buf bytes.Buffer
	buf.WriteString(planMagic)
	buf.WriteByte(planVersion)
	if err := gob.NewEncoder(&buf).Encode(p.Root); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodePlan(b []byte) (*Plan, error) {
	if len(b) < 5 || string(b[:4]) != planMagic {
		return nil, fmt.Errorf("plan: bad magic")
	}
	if b[4] != planVersion {
		return nil, fmt.Errorf("plan: unsupported version %d", b[4])
	}
	var n PlanNode
	if err := gob.NewDecoder(bytes.NewReader(b[5:])).Decode(&n); err != nil {
		return nil, err
	}
	return &Plan{Root: n}, nil
}

func node(k string, kids ...PlanNode) PlanNode {
	return PlanNode{K: k, Kids: kids}
}

func nodeOp(k, op string, kids ...PlanNode) PlanNode {
	return PlanNode{K: k, Op: op, Kids: kids}
}

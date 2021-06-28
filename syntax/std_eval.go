package syntax

import (
	"context"
	"fmt"

	"github.com/go-errors/errors"

	"github.com/arr-ai/arrai/rel"
)

type EvalConfig struct {
	scopes rel.Tuple
	stdlib rel.Tuple
}

func stdEval() rel.Attr {
	return rel.NewTupleAttr("eval",
		//TODO: eval needs to be changed to only evaluate simple expression
		// e.g. no functions, no math operations etc only simple values
		rel.NewNativeFunctionAttr("value", evalExpr),
		//TODO: eval.expr
	)
}

func evalExpr(ctx context.Context, v rel.Value) (rel.Value, error) {
	switch val := v.(type) {
	case rel.String, rel.Bytes:
		evaluated, err := EvaluateExpr(ctx, ".", val.String())
		if err != nil {
			panic(err)
		}
		return evaluated, nil
	}
	return nil, fmt.Errorf("//eval.value: not a byte array or string: %v", v)
}

func evalEval(ctx context.Context, config EvalConfig, v rel.Value) (rel.Value, error) {
	if config == (EvalConfig{}) {
		config.stdlib = SafeStdScope()
	}
	switch val := v.(type) {
	case rel.String, rel.Bytes:
		evaluated, err := EvaluateExpr(ctx, ".", val.String())
		if err != nil {
			panic(err)
		}
		return evaluated, nil
	}
	return nil, fmt.Errorf("//eval.eval: not a byte array or string: %v", v)
}

// parseEvalConfig returns the config arg as a evalConfig.
func parseEvalConfig(configArg rel.Value) (*EvalConfig, error) {
	config, ok := configArg.(*rel.GenericTuple)
	if !ok {
		return nil, errors.Errorf("first arg (config) must be tuple, not %s", rel.ValueTypeAsString(configArg))
	}
	parsedConfig := EvalConfig{}
	scopes := config.Without("stdlib")
	if scopes.IsTrue() {
		parsedConfig.scopes = scopes
	}
	stdlib := config.Without("scopes")
	if stdlib.IsTrue() {
		parsedConfig.stdlib = stdlib
	}
	return &parsedConfig, nil
}

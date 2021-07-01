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

// contextualEval evaluates an expression using only the provided scope and stdlib
func contextualEval(ctx context.Context, config EvalConfig, v rel.Value) (rel.Value, error) {
	var scope = rel.EmptyScope
	if config.stdlib != nil {
		localLib, err := config.stdlib.Eval(ctx, scope)
		if err != nil {
			return nil, err
		}
		scope = scope.With("//", localLib)
	} else {
		scope = scope.With("//", SafeStdScopeTuple())
	}
	for e := config.scopes.Enumerator(); e.MoveNext(); {
		name, value := e.Current()
		scope = scope.With(name, value)
	}
	switch val := v.(type) {
	case rel.String, rel.Bytes:
		evaluated, err := EvalWithScope(ctx, "", val.String(), scope)
		if err != nil {
			return nil, err
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
	scopes, found := config.Get("scope")
	if found {
		parsedConfig.scopes = scopes.(rel.Tuple)
	} else {
		parsedConfig.scopes = rel.EmptyTuple
	}
	stdlib, found := config.Get("stdlib")
	if found {
		parsedConfig.stdlib = stdlib.(rel.Tuple)
	}
	return &parsedConfig, nil
}

package syntax

import (
	"context"
	"os/exec"

	"github.com/arr-ai/arrai/rel"
	"github.com/go-errors/errors"
	"github.com/sirupsen/logrus"
)

func stdDeprecated() rel.Attr {
	return rel.NewTupleAttr("deprecated",
		stdDeprecatedExec(),
	)
}

func stdDeprecatedExec() rel.Attr {
	return rel.NewNativeFunctionAttr("exec", func(_ context.Context, value rel.Value) (rel.Value, error) {
		logrus.Warn("//deprecated is deprecated and may break or disappear at any moment")

		var cmd *exec.Cmd
		switch t := value.(type) {
		case rel.Array:
			if len(t.Values()) == 0 {
				return nil, errors.Errorf("//deprecated.exec arg must not be empty")
			}

			name := t.Values()[0].String()
			args := make([]string, len(t.Values())-1)
			for i, v := range t.Values()[1:] {
				args[i] = v.String()
			}
			cmd = exec.Command(name, args...)
		default:
			return nil, errors.Errorf("//deprecated.exec arg must be an array, not %s", rel.ValueTypeAsString(value))
		}

		var stderr []byte
		status := 0
		stdout, err := cmd.Output()
		if err != nil {
			exit, ok := err.(*exec.ExitError)
			if !ok {
				return nil, err
			}
			stderr = exit.Stderr
			status = exit.ExitCode()
		}
		return rel.NewTuple(
			rel.NewAttr("args", value.(rel.Array)),
			rel.NewAttr("exitCode", rel.NewNumber(float64(status))),
			rel.NewAttr("stdout", rel.NewBytes(stdout)),
			rel.NewAttr("stderr", rel.NewBytes(stderr)),
		), nil
	})
}

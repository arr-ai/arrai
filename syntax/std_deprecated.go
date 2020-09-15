package syntax

import (
	"context"
	"os/exec"

	"github.com/arr-ai/arrai/rel"
	"github.com/go-errors/errors"
	"github.com/google/shlex"
	"github.com/sirupsen/logrus"
)

func stdDeprecated() rel.Attr {
	logrus.Warn("//deprecated is deprecated and may break or disappear at any moment")
	return rel.NewTupleAttr("deprecated",
		stdDeprecatedExec(),
	)
}

func stdDeprecatedExec() rel.Attr {
	return rel.NewNativeFunctionAttr("exec", func(_ context.Context, value rel.Value) (rel.Value, error) {
		var cmd *exec.Cmd
		switch t := value.(type) {
		case rel.Array:
			if len(t.Values()) == 0 {
				return nil, errors.Errorf("//deprecated.exec arg must not be empty")
			}

			name := t.Values()[0].String()
			args := make([]string, len(t.Values()))
			for i, v := range t.Values() {
				if i == 0 {
					continue
				}
				args = append(args, v.String())
			}
			cmd = exec.Command(name, args...)
		case rel.String:
			args, err := shlex.Split(t.String())
			if err != nil {
				return nil, err
			}
			cmd = exec.Command(args[0], args[1:]...)
		default:
			return nil, errors.Errorf("//deprecated.exec arg must be a string or array, not %T", value)
		}
		out, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		return rel.NewBytes(out), nil
	})
}

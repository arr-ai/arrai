package syntax

import (
	"github.com/arr-ai/arrai/rel"
	"github.com/go-errors/errors"
)

// getConfigBool returns the boolean value from a 'config tuple' with the given key, and an "ok"
// flag indicating whether the key was present in the config tuple (true) or not (false). A 'config
// tuple' is a design pattern used to configure the behaviour of functions in the standard library.
func getConfigBool(config rel.Tuple, key string) (bool, bool) {
	vv, ok := config.Get(key)
	return ok && vv.IsTrue(), ok
}

// getConfigInt returns the integer value from a 'config tuple' with the given key, or the default
// value if the key does not exist. A 'config tuple' is a design pattern used to configure the
// behaviour of functions in the standard library.
func getConfigInt(config rel.Tuple, fn string, key string, defaultVal int) (int, error) {
	if vv, ok := config.Get(key); ok {
		vn, ok := vv.(rel.Number)
		if !ok {
			return 0, errors.Errorf("%s config param to %s must be integer, not %s", key, fn, rel.ValueTypeAsString(vv))
		}
		v, ok := vn.Int()
		if !ok {
			return 0, errors.Errorf("%s config param to %s must be integer, not %v", key, fn, vn)
		}
		return v, nil
	}
	return defaultVal, nil
}

// getConfigString returns the string value from a 'config tuple' with the given key, or the default
// value if the key does not exist. A 'config tuple' is a design pattern used to configure the
// behaviour of functions in the standard library.
func getConfigString(config rel.Tuple, fn string, key string, defaultVal string) (string, error) {
	if vv, ok := config.Get(key); ok {
		vn, ok := vv.(rel.String)
		if !ok {
			return "", errors.Errorf("%s config param to %s must be string, not %s", key, fn, rel.ValueTypeAsString(vv))
		}
		return vn.String(), nil
	}
	return defaultVal, nil
}

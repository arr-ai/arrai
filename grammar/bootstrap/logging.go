package bootstrap

import (
	"reflect"
)

var depth = 0

func indentf(format string, args ...interface{}) {
	// logrus.Tracef(strings.Repeat("    ", depth)+format, args...)
}

type enterexit struct{}

func enterf(format string, args ...interface{}) enterexit { //nolint:unparam
	indentf("--> "+format, args...)
	depth++
	return enterexit{}
}

func (enterexit) exitf(format string, ptrs ...interface{}) { //nolint:unparam
	depth--
	args := make([]interface{}, 0, len(ptrs))
	for _, ptr := range ptrs {
		args = append(args, reflect.ValueOf(ptr).Elem().Interface())
	}
	indentf("<-- "+format, args...)
}

package rel

import (
	"fmt"
	"reflect"
)

var kinds = map[int]reflect.Type{}

func registerKind(kind int, t reflect.Type) int {
	if u, has := kinds[kind]; has {
		if t != u {
			panic(fmt.Errorf("kind %d %v already registered for a different type %v", kind, u, t))
		}
	}
	kinds[kind] = t
	return kind
}

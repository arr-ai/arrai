package rel

// var depth int64 = 0

// func ptrsToArgs(ptrs ...interface{}) []interface{} {
// 	args := make([]interface{}, 0, len(ptrs))
// 	for _, ptr := range ptrs {
// 		args = append(args, reflect.ValueOf(ptr).Elem().Interface())
// 	}
// 	return args
// }

// func indentf(format string, args ...interface{}) {
// 	log.Printf(strings.Repeat("    ", int(depth))+format, args...)
// }

// func indentpf(format string, ptrs ...interface{}) {
// 	indentf(format, ptrsToArgs(ptrs...)...)
// }

// type enterexit struct{}

// func enterf(format string, args ...interface{}) enterexit { //nolint:unparam
// 	indentf("--> "+format, args...)
// 	atomic.AddInt64(&depth, 1)
// 	return enterexit{}
// }

// func (enterexit) exitf(format string, ptrs ...interface{}) { //nolint:unparam
// 	atomic.AddInt64(&depth, -1)
// 	indentpf("<-- "+format, ptrs...)
// }

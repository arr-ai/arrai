package syntax

import (
	"context"
	"io"
	"io/ioutil"
	"sync"

	"github.com/arr-ai/arrai/rel"
)

func stdOsSafe() rel.Attr {
	return rel.NewTupleAttr("os",
		stdOsSafeAttrs()...,
	)
}
func stdOsUnsafe() rel.Attr {
	return rel.NewTupleAttr("os",
		stdOsUnsafeAttrs()...,
	)
}

func stdOsSafeAttrs() []rel.Attr {
	a := make([]rel.Attr, 0)
	a = append(a, rel.NewAttr("path_separator", stdOsPathSeparator()))
	a = append(a, rel.NewAttr("path_list_separator", stdOsPathListSeparator()))
	a = append(a, rel.NewAttr("cwd", stdOsCwd()))
	a = append(a, rel.NewNativeFunctionAttr("exists", stdOsExists))
	a = append(a, rel.NewNativeFunctionAttr("tree", stdOsTree))
	a = append(a, rel.NewNativeFunctionAttr("get_env", stdOsGetEnv))
	a = append(a, rel.NewNativeFunctionAttr("&args", stdOsGetArgs))
	a = append(a, rel.NewNativeFunctionAttr("&stdin", stdOsStdinVar.read))
	a = append(a, rel.NewNativeFunctionAttr("isatty", stdOsIsATty))
	return a
}
func stdOsUnsafeAttrs() []rel.Attr {
	a := make([]rel.Attr, 0)
	a = append(a, rel.NewNativeFunctionAttr("file", stdOsFile))
	return a
}

type stdOsStdin struct {
	reader io.Reader
	mutex  sync.Mutex
	bytes  rel.Value
}

func newStdOsStdin(r io.Reader) *stdOsStdin {
	return &stdOsStdin{reader: r}
}

func (d *stdOsStdin) reset(r io.Reader) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.reader = r
	d.bytes = nil
}

func (d *stdOsStdin) read(context.Context, rel.Value) (rel.Value, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if d.bytes != nil {
		return d.bytes, nil
	}
	f, err := ioutil.ReadAll(stdOsStdinVar.reader)
	if err != nil {
		return nil, err
	}
	d.bytes = rel.NewBytes(f)
	return d.bytes, nil
}

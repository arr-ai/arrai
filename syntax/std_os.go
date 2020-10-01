package syntax

import (
	"context"
	"io"
	"io/ioutil"
	"sync"

	"github.com/arr-ai/arrai/rel"
)

func stdOs() rel.Attr {
	return rel.NewTupleAttr("os",
		rel.NewAttr("path_separator", stdOsPathSeparator()),
		rel.NewAttr("path_list_separator", stdOsPathListSeparator()),
		rel.NewAttr("cwd", stdOsCwd()),
		rel.NewNativeFunctionAttr("exists", stdOsExists),
		rel.NewNativeFunctionAttr("file", stdOsFile),
		rel.NewNativeFunctionAttr("tree", stdOsTree),
		rel.NewNativeFunctionAttr("get_env", stdOsGetEnv),
		rel.NewNativeFunctionAttr("&args", stdOsGetArgs),
		rel.NewNativeFunctionAttr("&stdin", stdOsStdinVar.read),
		rel.NewNativeFunctionAttr("isatty", stdOsIsATty),
	)
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

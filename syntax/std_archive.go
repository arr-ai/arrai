package syntax

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"context"
	"io"
	"path"

	"github.com/pkg/errors"

	"github.com/arr-ai/arrai/rel"
)

func stdArchive() rel.Attr {
	return rel.NewTupleAttr("archive",
		rel.NewTupleAttr("tar",
			rel.NewNativeFunctionAttr("tar", func(_ context.Context, v rel.Value) (rel.Value, error) {
				return createArchive(v, func(w io.Writer) (io.Closer, func(string, []byte) (io.Writer, error)) {
					aw := tar.NewWriter(w)
					return aw, func(path string, data []byte) (io.Writer, error) {
						return aw, aw.WriteHeader(&tar.Header{
							Name: path,
							Mode: 0600,
							Size: int64(len(data)),
						})
					}
				})
			}),
		),
		rel.NewTupleAttr("zip",
			rel.NewNativeFunctionAttr("zip", func(_ context.Context, v rel.Value) (rel.Value, error) {
				return createArchive(v, func(w io.Writer) (io.Closer, func(string, []byte) (io.Writer, error)) {
					aw := zip.NewWriter(w)
					return aw, func(path string, _ []byte) (io.Writer, error) {
						return aw.Create(path)
					}
				})
			}),
		),
	)
}

func createArchive(
	v rel.Value,
	creator func(io.Writer) (io.Closer, func(string, []byte) (io.Writer, error)),
) (rel.Set, error) {
	var b bytes.Buffer
	closer, create := creator(&b)
	d, ok := v.(rel.Dict)
	if !ok {
		return nil, errors.Errorf("//archive.zip.zip arg not a dict: %v", v)
	}
	if err := writeDictToArchive(d, create, ""); err != nil {
		return nil, err
	}
	if err := closer.Close(); err != nil {
		return nil, err
	}
	return rel.NewBytes(b.Bytes()), nil
}

func writeDictToArchive(d rel.Dict, w func(string, []byte) (io.Writer, error), parent string) error {
	for e := d.DictEnumerator(); e.MoveNext(); {
		k, v := e.Current()
		name, is := rel.AsString(k.(rel.Set))
		if !is {
			return errors.Errorf("dict key %v not a string", k)
		}
		subpath := path.Join(parent, name.String())
		switch v := v.(type) {
		case rel.Set:
			if s, ok := rel.AsString(v); ok {
				data := []byte(s.String())
				fw, err := w(subpath, data)
				if err != nil {
					return err
				}
				if _, err = fw.Write(data); err != nil {
					return err
				}
			} else if d, ok := v.(rel.Dict); ok {
				if err := writeDictToArchive(d, w, subpath); err != nil {
					return err
				}
			}
		default:
			return errors.Errorf("unsupported entry %q: %s %v", subpath, rel.ValueTypeAsString(v), v)
		}
	}
	return nil
}

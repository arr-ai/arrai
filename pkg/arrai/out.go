package arrai

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/rel"
	"github.com/go-errors/errors"
	"github.com/spf13/afero"
)

const (
	dirField        = "dir"
	fileField       = "file"
	ifExistsConfig  = "ifExists"
	ifExistsRemove  = "remove"
	ifExistsReplace = "replace"
	ifExistsMerge   = "merge"
	ifExistsIgnore  = "ignore"
	ifExistsFail    = "fail"
)

var errFileAndDirMustNotExist = errors.Errorf("%s and %s must not exist", dirField, fileField)
var errFileOrDirMustExist = errors.Errorf("exactly one of %s or %s must exist", dirField, fileField)

// OutputValue handles output writing for evaluated values.
func OutputValue(ctx context.Context, value rel.Value, w io.Writer, out string) error {
	if out != "" {
		return outputValue(ctx, value, out)
	}

	var s string
	switch v := value.(type) {
	case rel.String:
		s = v.String()
	case rel.Bytes:
		s = v.String()
	case rel.Set:
		if !v.IsTrue() {
			s = ""
		} else {
			s = rel.Repr(v)
		}
	default:
		s = rel.Repr(v)
	}
	fmt.Fprintf(w, "%s", s)
	if s != "" && !strings.HasSuffix(s, "\n") {
		if _, err := w.Write([]byte{'\n'}); err != nil {
			return err
		}
	}

	return nil
}

func outputValue(ctx context.Context, value rel.Value, out string) error {
	parts := strings.SplitN(out, ":", 2)
	if len(parts) == 1 {
		parts = []string{"", parts[0]}
	}
	mode := parts[0]
	arg := parts[1]

	fs := ctxfs.RuntimeFsFrom(ctx)
	switch mode {
	case "file", "f", "":
		return outputFile(value, arg, fs, false)
	case "dir", "d":
		if t, is := value.(rel.Dict); is {
			if err := outputTupleDir(t, arg, fs, true); err != nil {
				return err
			}
			return outputTupleDir(t, arg, fs, false)
		}
		return fmt.Errorf("result not a dict: %v", value)
	}
	return fmt.Errorf("invalid --out flag: %s", out)
}

func outputTupleDir(v rel.Value, dir string, fs afero.Fs, dryRun bool) error {
	t, err := getDirField(v)
	if err != nil {
		return err
	}
	if _, err := fs.Stat(dir); os.IsNotExist(err) {
		if err := fs.Mkdir(dir, 0755); err != nil {
			return err
		}
	}

	// this is to allow empty directory
	if !t.IsTrue() {
		return nil
	}

	for e := t.(rel.Dict).DictEnumerator(); e.MoveNext(); {
		k, v := e.Current()
		name, is := k.(rel.String)
		if !is {
			return fmt.Errorf("dir output dict key must be a non-empty string")
		}
		subpath := path.Join(dir, name.String())
		switch content := v.(type) {
		case rel.Tuple:
			if err := configureOutput(content, subpath, fs, dryRun); err != nil {
				return err
			}
		case rel.Dict:
			if err := outputTupleDir(content, subpath, fs, dryRun); err != nil {
				return err
			}
		case rel.Bytes, rel.String:
			if err := outputFile(content, subpath, fs, dryRun); err != nil {
				return err
			}
		case rel.Set:
			if content.IsTrue() {
				return fmt.Errorf("dir output entry must be dict, string or byte array")
			}
			if err := outputFile(content, subpath, fs, dryRun); err != nil {
				return err
			}
		}
	}
	return nil
}

func outputFile(content rel.Value, path string, fs afero.Fs, dryRun bool) error {
	var bytes []byte
	switch content := content.(type) {
	case rel.Bytes:
		bytes = content.Bytes()
	case rel.String:
		bytes = []byte(content.String())
	default:
		if _, is := content.(rel.Set); !(is && !content.IsTrue()) {
			return fmt.Errorf("file output not string or byte array: %v", content)
		}
		bytes = []byte{}
	}

	if dryRun {
		return nil
	}

	f, err := fs.Create(path)
	if err != nil {
		return err
	}
	_, err = f.Write(bytes)
	return err
}

func configureOutput(t rel.Tuple, dir string, fs afero.Fs, dryRun bool) error {
	configNames := []string{ifExistsConfig}

	for _, c := range configNames {
		if t.HasName(c) {
			for _, applier := range getConfigurators() {
				if err := applier(t, dir, fs, dryRun); err != nil {
					return err
				}
			}
			return nil
		}
	}
	return applyFilesFields(t, dir, fs, dryRun)
}

func getConfigurators() []func(rel.Tuple, string, afero.Fs, bool) error {
	// mind the order when adding new configurations
	return []func(rel.Tuple, string, afero.Fs, bool) error{
		applyIfExistsConfig,
	}
}

func applyIfExistsConfig(t rel.Tuple, dir string, fs afero.Fs, dryRun bool) (err error) {
	conf, has := t.Get(ifExistsConfig)
	if !has {
		return nil
	}
	errInvalidConfig := errors.Errorf(
		"%s: value '%s' is not valid value. It has to be one of %s",
		ifExistsConfig, conf,
		strings.Join([]string{ifExistsMerge, ifExistsRemove, ifExistsReplace, ifExistsIgnore, ifExistsFail}, ", "),
	)

	if _, isString := conf.(rel.String); !isString {
		return errInvalidConfig
	}
	switch conf.String() {
	case ifExistsIgnore, ifExistsRemove, ifExistsReplace, ifExistsFail:
	case ifExistsMerge:
		if t.HasName(fileField) {
			return errors.Errorf("%s: '%s' config must not have '%s' field", ifExistsConfig, fileField, ifExistsMerge)
		}
	default:
		return errInvalidConfig
	}

	if _, err := fs.Stat(dir); os.IsNotExist(err) {
		if conf.String() != ifExistsRemove {
			return applyFilesFields(t, dir, fs, dryRun)
		}
	} else if err != nil {
		return err
	}

	switch conf.String() {
	case ifExistsRemove:
		if err := checkNotDirAndNotFileField(t); err != nil {
			return err
		}
		if dryRun {
			return nil
		}
		return fs.RemoveAll(dir)
	case ifExistsReplace:
		if err := checkDirXorFileField(t); err != nil {
			return err
		}
		if dryRun {
			return nil
		}
		if err := fs.RemoveAll(dir); err != nil {
			return err
		}
		return applyFilesFields(t, dir, fs, dryRun)
	case ifExistsMerge:
		if v, has := t.Get(dirField); has {
			d, err := getDirField(v)
			if err != nil {
				return err
			}
			return outputTupleDir(d, dir, fs, dryRun)
		}
		return errors.Errorf("%s: '%s' field must exist", ifExistsConfig, dirField)
	case ifExistsIgnore:
		return nil
	case ifExistsFail:
		return errors.Errorf("%s: '%s' exists", ifExistsConfig, dir)
	}
	panic("impossible")
}

func checkNotDirAndNotFileField(t rel.Tuple) error {
	_, hasDirs := t.Get(dirField)
	_, hasFiles := t.Get(fileField)
	if hasDirs || hasFiles {
		return errFileAndDirMustNotExist
	}
	return nil
}

func checkDirXorFileField(t rel.Tuple) error {
	_, hasDirs := t.Get(dirField)
	_, hasFiles := t.Get(fileField)
	if hasDirs == hasFiles {
		return errFileOrDirMustExist
	}
	return nil
}

func applyFilesFields(t rel.Tuple, path string, fs afero.Fs, dryRun bool) error {
	if dir, has := t.Get(dirField); has {
		d, err := getDirField(dir)
		if err != nil {
			return err
		}
		return outputTupleDir(d, path, fs, dryRun)
	}
	if file, has := t.Get(fileField); has {
		return outputFile(file, path, fs, dryRun)
	}
	return errFileOrDirMustExist
}

func getDirField(v rel.Value) (rel.Set, error) {
	switch k := v.(type) {
	case rel.Dict:
		return k, nil
	case rel.GenericSet:
		if !k.IsTrue() {
			return k, nil
		}
	}
	return nil, errors.Errorf("%s must be of type Dictionary, not %T", dirField, v)
}

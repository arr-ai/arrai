package module

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type Module interface {
	Add(*Mod)
	Load() error
	Get(string) (*Mod, error)
}

type Mod struct {
	Name string
	Dir  string
}

type BaseModule []*Mod

func (m *BaseModule) Add(v *Mod) {
	*m = append(*m, v)
}

func (m *BaseModule) Load() error {
	*m = []*Mod{}
	return nil
}

func (m *BaseModule) Get(filename string) (*Mod, error) {
	for _, mod := range *m {
		if hasPathPrefix(mod.Name, filename) {
			return mod, nil
		}
	}

	return nil, errors.Errorf("module of file %s not found", filename)
}

func hasPathPrefix(prefix, s string) bool {
	prefix = path.Clean(prefix)
	s = path.Clean(s)

	return strings.HasPrefix(s, prefix+string(filepath.Separator)) || s == prefix
}

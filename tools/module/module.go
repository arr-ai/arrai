package module

import (
	"path/filepath"

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

	return nil, errors.Errorf("Module of file %s not found", filename)
}

func hasPathPrefix(prefix, s string) bool {
	prefix = filepath.Clean(prefix)
	s = filepath.Clean(s)

	if len(s) > len(prefix) {
		return s[len(prefix)] == filepath.Separator && s[:len(prefix)] == prefix
	}

	return s == prefix
}

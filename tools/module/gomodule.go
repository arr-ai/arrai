package module

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type GoModule struct {
	BaseModule
}

func NewGoModule() *GoModule {
	return &GoModule{}
}

func (m *GoModule) Load() error {
	err := m.BaseModule.Load()
	if err != nil {
		return err
	}

	out := ioutil.Discard

	err = runGo(context.Background(), out, "mod", "download")
	if err != nil {
		return errors.Wrap(err, "failed to download modules")
	}

	b := &bytes.Buffer{}
	err = runGo(context.Background(), b, "list", "-m", "-json", "all")
	if err != nil {
		return errors.Wrap(err, "failed to list modules")
	}

	dec := json.NewDecoder(b)
	for {
		gm := &goMod{}
		if err := dec.Decode(gm); err != nil {
			if err == io.EOF {
				break
			}
			return errors.Wrap(err, "failed to decode modules list")
		}

		m.Add(&Mod{
			Name: gm.Path,
			Dir:  gm.Dir,
		})
	}

	return nil
}

type goMod struct {
	Path string
	Dir  string
}

func (m *GoModule) Get(filename string) (*Mod, error) {
	err := goGetByFilepath(filename)
	if err != nil {
		return nil, err
	}

	if err := m.Load(); err != nil {
		return nil, err
	}

	return m.BaseModule.Get(filename)
}

func goGetByFilepath(filename string) error {
	names := strings.Split(string(filename), string(os.PathSeparator))
	if len(names) > 0 {
		gogetPath := names[0]

		for i := 1; i < len(names); i++ {
			err := goGet(gogetPath)
			if err == nil {
				return nil
			}
			logrus.Debugf("go get %s error: %s\n", gogetPath, err.Error())

			gogetPath = filepath.Join(gogetPath, names[i])
		}
	}

	return errors.New("No such module")
}

func goGet(args ...string) error {
	if err := runGo(context.Background(), logrus.StandardLogger().Out, append([]string{"get", "-u"}, args...)...); err != nil {
		return errors.Wrapf(err, "failed to get %q", args)
	}
	return nil
}

func runGo(ctx context.Context, out io.Writer, args ...string) error {
	cmd := exec.CommandContext(ctx, "go", args...)

	wd, err := os.Getwd()
	if err != nil {
		return errors.Errorf("get current working directory error: %s\n", err.Error())
	}
	cmd.Dir = wd

	errbuf := new(bytes.Buffer)
	cmd.Stderr = errbuf
	cmd.Stdout = out

	logrus.Debugf("running command `go %v`\n", strings.Join(args, " "))
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.Error); ok && ee.Err == exec.ErrNotFound {
			return nil
		}

		_, ok := err.(*exec.ExitError)
		if !ok {
			return errors.Errorf("failed to execute 'go %v': %s %T", args, err, err)
		}

		// Too old Go version
		if strings.Contains(errbuf.String(), "flag provided but not defined") {
			return nil
		}
		return errors.Errorf("go command failed: %s", errbuf)
	}

	return nil
}

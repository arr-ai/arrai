// Code generated for package arrai by go-bindata DO NOT EDIT. (@generated)
// sources:
// internal/arrai/echo.arraiz
package arrai

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _internalArraiEchoArraiz = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x0a\xf0\x66\x66\x11\x61\xe0\x60\xe0\x60\x40\x07\x3c\x0c\x0c\x0c\xc9\xf9\x79\x69\x99\xe9\x7a\x89\x45\x45\x89\x99\x3d\xa7\x7d\xb9\x98\x0c\x79\x8e\x7e\x5c\x3d\x65\x8d\x48\xc9\xa6\xf9\xa7\x3d\x96\xf5\x9a\x4d\xe0\xef\xd4\xad\xa8\xd9\xf8\xe0\xde\xd5\xea\x0f\xcf\xcb\x8d\x84\xf6\xa7\x1e\xb9\xc1\x6d\x93\xf2\xa5\x32\xc7\x94\x47\xf3\xb3\xd7\x27\x59\xa1\x27\xe6\x37\x83\x7f\xb6\x1e\xff\xb2\x8c\x67\xfa\xa6\x2b\xb1\x53\xbe\x33\x65\x5c\x5c\x66\xf4\x6f\x65\x60\x4d\x81\xd5\x4b\xfd\xeb\x72\xb7\x77\x15\xfc\x63\x61\x60\xf8\xff\x3f\xc0\x9b\x9d\xc3\x27\xb4\xfd\x4b\x3c\x03\x03\x43\x1f\x03\x03\x03\x6e\x87\x99\x32\x30\x30\xe4\xe6\xa7\x94\xe6\xa4\xea\xa7\x67\x96\x64\x94\x26\xe9\x25\xe7\xe7\xea\x27\x16\x15\xe9\x26\x66\xea\x83\x9d\xaa\x9f\x5a\x91\x98\x5b\x90\x93\x5a\xac\x9f\x54\x9a\x97\x02\x52\x97\xaf\x97\x9b\x9f\x72\xea\xac\xbf\xe6\x59\xcf\x50\x8f\xf3\x3a\x27\x35\x7c\x2f\x79\x9f\x3b\x7f\xd5\x5f\x47\xeb\x92\xef\x19\x96\xc0\xce\x99\xbf\x56\x8a\x76\xee\x64\xb5\x5c\x79\xf4\x97\xe7\xd2\x4a\x4d\xcb\x95\x31\x31\x2f\x7f\x72\xb5\x55\xb5\xa9\xc5\xc4\xb8\xae\x6c\x5d\x3a\xd3\xd5\xa2\x54\x6b\x8a\xe4\x0b\xc9\x59\xd1\xaf\x24\x1b\x42\xc4\x9e\xa8\xa8\x3d\x0e\xcb\xcc\x5e\xbe\x60\xf9\xf2\xf5\xcb\x8f\xb3\x72\xc0\x7c\x10\x74\x47\x48\x15\xe4\x83\x1e\xbc\x3e\xf0\x20\xdd\x07\x99\x79\x25\xa9\x45\x79\x89\x39\x30\xe9\xe4\x8c\x7c\x48\xb4\x9c\xf2\xd4\x0d\xec\xe8\xbb\x1e\x78\xe1\xf2\x65\x6f\x1d\x2f\xfd\xbe\xd0\x0d\xa1\x17\xae\x9f\xd7\x67\x2b\x5e\xb3\x73\x47\x88\x50\xc8\x2a\x47\x1f\x6d\xed\xd6\x5b\x2c\x7a\x0e\x0c\x30\x37\xfe\x59\xaa\x62\xa1\xc7\xc0\xc0\xa0\x0b\x76\x23\x23\x93\x08\x03\xc2\x95\xc8\x31\xc0\x83\xe1\x6e\xe4\xe4\x80\xae\x13\xd9\xe7\xa6\x28\xba\x66\x92\x1b\x5f\xe8\x56\x20\x3b\xdc\x03\xc5\x8a\x68\x46\xea\x05\x68\x80\x37\x2b\x1b\xc8\x4c\x66\x06\x66\x06\x61\x46\x06\x86\xff\x8c\x20\x1e\x20\x00\x00\xff\xff\x4b\x7f\x61\x80\x28\x03\x00\x00")

func internalArraiEchoArraizBytes() ([]byte, error) {
	return bindataRead(
		_internalArraiEchoArraiz,
		"internal/arrai/echo.arraiz",
	)
}

func internalArraiEchoArraiz() (*asset, error) {
	bytes, err := internalArraiEchoArraizBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "internal/arrai/echo.arraiz", size: 808, mode: os.FileMode(420), modTime: time.Unix(1601467995, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"internal/arrai/echo.arraiz": internalArraiEchoArraiz,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"internal": {nil, map[string]*bintree{
		"arrai": {nil, map[string]*bintree{
			"echo.arraiz": {internalArraiEchoArraiz, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

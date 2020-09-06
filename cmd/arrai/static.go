// Code generated for package main by go-bindata DO NOT EDIT. (@generated)
// sources:
// internal/build/main.arrai
package main

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

var _internalBuildMainArrai = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x51\xc1\x6e\xdb\x30\x0c\xbd\xeb\x2b\x38\xa3\x07\x1b\x70\xa5\xbb\x87\x1c\xd2\xd9\x18\x7a\xa9\x81\xce\xdb\x65\x1b\x50\x46\x51\x55\x21\xb6\x64\x48\x54\xe7\x22\xf0\xbf\x0f\x96\x1d\xa0\xb7\x44\x17\x91\xe2\x7b\x8f\x7c\xe2\x9f\xc3\x07\xa9\xc0\xee\x5e\x98\x10\xb0\xff\xd9\xb5\xdf\x9b\xa7\xe6\x79\xdf\x35\x35\xdc\x43\xdd\xc2\x53\xdb\x41\x53\x3f\x76\x6c\x44\x79\x42\xad\x60\x40\x63\x19\x33\xc3\xe8\x3c\x41\xce\x00\x00\x32\xe9\x2c\xa9\x89\xb2\x35\x73\x21\x63\x6b\xa4\x0d\xbd\xc5\x03\x97\x6e\x10\xe8\xfd\x3d\x9a\xe5\x42\x23\xc6\x93\x5e\x23\x49\x53\x76\x23\xf4\x0a\x2e\x7c\x58\xc2\x6b\x62\xe4\x5c\x1f\x32\x56\x30\xf6\x1a\xad\x4c\x56\xf2\x02\xce\x89\x24\x04\x74\x6d\xdd\x56\x10\x14\x81\x10\x2e\x70\xf4\x3a\x00\x7a\x1d\x07\x65\x09\xfe\x19\x7a\x73\x91\x40\xf7\xee\x80\x3d\xbc\xa3\x07\x0c\xc1\x68\xbb\x54\x93\x42\x52\xe7\xfb\x8d\x10\x60\x07\x2e\xa5\x21\x55\x25\x4d\x50\xed\xe0\x62\x9b\x3f\x5a\x43\xcf\xd1\x7e\xa3\x29\xdf\xbe\x8f\x3f\xa0\x3c\x69\xef\xa2\x3d\xe6\x45\x91\x48\xef\xd8\x47\x55\x82\xf2\x7e\xe1\xae\x1e\x79\xb3\xbc\x22\xa9\x87\x68\x8f\xbd\xca\x25\x4d\x25\xfc\xfe\xbb\xec\xf1\x7c\x77\x4e\xeb\xac\xaa\x12\xe6\x79\x95\x30\xaf\x89\xfe\x65\x07\xd6\xf4\x9b\xd7\xe5\x8c\x68\x8d\xcc\x95\xf7\x2b\x6c\x66\x9f\xd1\x97\x41\x79\x1b\x69\x8c\xf4\x6b\x19\x63\x6d\xb4\x4d\xe4\x02\xff\x41\x47\x17\xa9\x84\x2c\x2b\xbe\xde\xd8\x62\x66\x2f\xec\x7f\x00\x00\x00\xff\xff\x0b\x77\x36\xa5\x71\x02\x00\x00")

func internalBuildMainArraiBytes() ([]byte, error) {
	return bindataRead(
		_internalBuildMainArrai,
		"internal/build/main.arrai",
	)
}

func internalBuildMainArrai() (*asset, error) {
	bytes, err := internalBuildMainArraiBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "internal/build/main.arrai", size: 625, mode: os.FileMode(420), modTime: time.Unix(1599436354, 0)}
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
	"internal/build/main.arrai": internalBuildMainArrai,
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
	"internal": &bintree{nil, map[string]*bintree{
		"build": &bintree{nil, map[string]*bintree{
			"main.arrai": &bintree{internalBuildMainArrai, map[string]*bintree{}},
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

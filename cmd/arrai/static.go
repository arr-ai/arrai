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

var _internalBuildMainArrai = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x51\xb1\x6e\xdb\x30\x10\xdd\xf9\x15\xaf\x6a\x06\x09\x70\xc4\xdd\x45\x06\xa7\x16\x8a\x2c\x11\x90\xaa\x5d\xda\x02\x39\xd3\x34\x4d\x58\x26\x05\xf2\x98\x2a\x30\xfc\xef\x85\x28\x05\xe8\x96\x68\xba\xd3\xbd\xf7\xee\x1e\xdf\x67\x74\x47\x1b\x11\x55\xb0\x03\xc3\x46\xb0\x87\xd1\x4e\x07\x62\x0d\x82\xf1\x3d\x39\x83\x83\xed\x35\xf8\x48\x8c\x90\x5c\x04\x61\x97\xdc\xbe\xd7\x7b\x50\x08\x64\x17\x76\xac\xc5\xef\xdd\x2b\xeb\x28\x6e\x9e\x85\x94\xd8\xfc\xe8\xda\x6f\xcd\x63\xf3\xb4\xe9\x9a\x2d\x6e\xb1\x6d\xf1\xd8\x76\x68\xb6\x0f\x9d\x18\x48\x9d\xc8\x68\x9c\xc9\x3a\x21\xec\x79\xf0\x81\x51\x0a\x00\x28\x94\x77\xac\x47\x2e\xe6\xce\xc7\x42\xcc\x95\xb1\x7c\x4c\xbb\x5a\xf9\xb3\xa4\x10\x6e\xc9\xca\xbc\x5d\x0e\x27\x33\x57\x8a\xc7\xe2\x83\xd0\x77\x70\xf1\xd5\x31\xbd\x27\xc6\xde\xf7\xb1\x10\x95\x10\x87\xe4\x54\xb6\x52\x56\xb8\x64\x92\x94\xe8\xda\x6d\xbb\x46\xd4\x0c\x29\x7d\xac\x29\x98\x08\x0a\x26\x9d\xb5\x63\xfc\xb5\x7c\xf4\x89\x61\x7a\xbf\xa3\x1e\x2f\x14\x40\x31\x5a\xe3\xa6\x69\x56\xc8\xea\xf5\x66\x21\x44\xdc\xc1\xe7\x36\xe6\xa9\xe2\x11\xeb\x3b\xbc\xd9\xae\x1f\x9c\xe5\xa7\xe4\xbe\xf2\x58\x2e\xcf\x57\xdf\x93\x3a\x99\xe0\x93\xdb\x97\x55\x95\x49\x2f\xd4\x27\xbd\x82\x0e\x61\xe2\xce\x1e\xeb\x66\xfa\x4b\xac\xef\x73\xa2\xa5\xe2\x71\x85\x5f\x7f\xa6\x1c\x2f\x37\x97\x1c\xe7\x7a\xbd\xc2\xf5\x3a\x4b\xd8\x43\xa6\x7f\xba\x83\xb3\xfd\xe2\x75\xfa\x06\x72\x56\x95\x3a\x84\x19\x76\x15\xff\xa3\xdf\x0e\xad\xdb\xc4\x43\xe2\x9f\xd3\x19\xf3\xa2\xe5\x22\x1f\xeb\xef\xbc\xf7\x89\x57\x28\x8a\xea\xcb\x07\x57\x5c\xc5\xb3\xf8\x17\x00\x00\xff\xff\x60\xea\xa5\x67\xbf\x02\x00\x00")

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

	info := bindataFileInfo{name: "internal/build/main.arrai", size: 703, mode: os.FileMode(420), modTime: time.Unix(1599456778, 0)}
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

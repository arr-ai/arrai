// Code generated by go-bindata. (@generated) DO NOT EDIT.

// Package main generated by go-bindata.// sources:
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
		return nil, fmt.Errorf("read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
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

// ModTime return file modify time
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

var _internalBuildMainArrai = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x90\xc1\x6e\xf2\x30\x10\x84\xef\xfb\x14\xfb\x5b\x1c\x12\x29\x24\x77\x7e\x71\x80\x12\x55\x5c\x88\x44\xd3\x5e\xda\x4a\x18\x63\x52\x0b\xb0\x23\x67\x5d\x19\xa1\xbc\x7b\xe5\x38\x48\xbd\x15\x5f\xbc\xb6\x66\x76\x3e\xcd\xc7\xfe\x4a\xb2\x83\xc9\x0e\x8a\x02\x17\xaf\x75\xf5\x5c\x6e\xca\xed\xa2\x2e\x57\x38\xc5\x55\x85\x9b\xaa\xc6\x72\xb5\xae\xa1\xe5\xe2\xc4\x1b\x89\x17\xae\x34\x80\xba\xb4\xc6\x12\x26\x80\x88\xc8\x84\xd1\x24\x3d\xb1\xf8\x32\x1d\x83\x38\x35\x8a\xbe\xdc\x3e\x17\xe6\x52\x70\x6b\xa7\x5c\x85\x8b\xab\xa2\x3d\x35\x71\x12\xe4\xd9\x83\xd2\x3f\x74\xdd\x55\x13\xf7\x0c\x52\x80\xa3\xd3\x62\xc0\x4c\x52\xbc\x0d\x2e\x41\x1e\x67\x73\xbc\x47\xe6\x6b\xad\x68\xeb\xf4\x13\xf9\x64\x44\xcf\x97\x5c\x9c\x1a\x6b\x9c\x3e\x24\x69\x3a\x98\xbe\xf9\xd9\xc9\x0c\xa5\xb5\xc1\x1b\xf7\xe7\x65\xf8\xe5\x24\x97\x4e\x1f\xce\x32\x11\xe4\x33\x7c\xff\x0c\x1d\xde\x26\xb7\xa1\xca\xd9\x2c\xc3\xbe\x8f\x2b\xd4\x71\xb0\xff\x9b\xa3\x56\xe7\x91\x25\x9c\x96\x6b\x25\x12\x69\x6d\x94\xf5\xf0\x5b\x7d\x07\xcd\x2b\x47\xad\xa3\xb7\x80\x11\x83\x46\x22\xd3\xe5\x2f\x74\x30\x8e\x32\x64\x2c\xfd\xff\x60\x44\x0f\xbb\x9f\x00\x00\x00\xff\xff\xd4\xd6\xdf\x36\xec\x01\x00\x00")

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

	info := bindataFileInfo{name: "internal/build/main.arrai", size: 492, mode: os.FileMode(420), modTime: time.Unix(1599394705, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
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
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
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
// AssetDir("foo.txt") and AssetDir("nonexistent") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
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
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}

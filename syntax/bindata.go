// Code generated for package syntax by go-bindata DO NOT EDIT. (@generated)
// sources:
// syntax/implicit_import.arrai
// syntax/stdlib-safe.arraiz
// syntax/stdlib-unsafe.arraiz
package syntax

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

var _syntaxImplicit_importArrai = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x93\x41\x8f\xdb\x2e\x10\xc5\xef\xfe\x14\x23\x65\xb5\x06\xc9\x7f\x7c\x27\xca\xfe\x73\xaa\xd4\x7b\x8f\x2b\xed\x52\x33\x5e\x53\x61\x70\x81\xa4\xb1\x5c\x7f\xf7\x6a\x70\x62\xd7\xad\xea\x93\x19\x7e\x0f\xde\x3c\xe0\x00\x5f\x3a\x13\x21\x36\xc1\x0c\x09\x3a\xe5\xb4\xc5\x08\xa9\x43\xe8\xd5\x30\x18\xf7\x01\xbe\x05\x8d\x8d\xd7\x18\x22\x18\x07\xa6\x1f\xac\x69\x4c\xa2\x1f\x1f\x52\x04\x8d\x03\x3a\x9d\x49\x57\x1c\xb2\x14\x6f\x09\x5d\x34\xde\x91\x98\x0a\x0b\x8b\x1a\x5a\x63\x51\x14\xc5\x61\x5b\xb2\x57\x43\xcc\x65\x68\x7d\xe8\x55\x8a\x90\x7c\xd6\xc4\xa4\x9c\x56\x41\x6f\xe8\x25\xa2\xa6\xd9\xa5\x90\x45\xd9\x52\xea\x54\xba\xab\x45\x61\x31\x6d\x8a\x13\x4c\x05\x00\x40\xd9\xc4\x6b\x29\xa1\xae\xd1\x35\x9e\xbc\x8a\x26\x5e\xc5\x82\x55\x0b\xf1\x2d\x7a\xb7\x47\xa8\xb2\x67\x6e\x36\xde\xf6\x0c\x55\xf6\xcc\xa8\x7a\x5b\x4a\x60\x77\x0f\x3b\x9a\xe6\x1e\x34\x28\x6b\x54\x94\x30\x95\x63\x6f\xcb\x99\x57\xc5\x7c\xa4\x60\x1a\xef\xae\x48\xb9\xe2\x15\xc3\x98\x3a\x0a\x36\x79\x50\xa0\x4d\x93\x8c\x77\x2a\x8c\x30\xad\x01\xcb\x47\xaf\x9f\xdc\xfc\x67\xeb\x2c\x1b\x5a\x0b\xad\xb1\x09\x03\x88\xf3\x55\xd9\x0b\xde\x83\xa1\x8f\xc9\x3b\x53\x81\xcc\xa6\xb8\x84\x89\x49\x71\xae\x60\x61\xd7\x4d\xf8\x0c\x3f\x81\x65\x06\x4e\x2f\xc0\xce\x52\xfc\xcd\xf0\x6a\x5d\xf9\x4d\xc2\x24\xe6\x65\x3c\x17\x1c\xfe\x7b\x81\xba\x0e\x68\xc5\xc5\x19\xef\x98\xe0\xc7\xa2\x78\xdd\x2e\xcb\xeb\xd7\x31\x61\x86\xa9\x93\xad\x7e\xa2\x50\x34\x4c\x75\x1d\xf1\xbb\xe8\x54\x7c\x1b\x02\xb6\xe6\xc6\x4a\x51\x56\x1b\xc7\x29\x6b\x22\x52\x30\xfd\x3f\x90\x8a\x3c\xad\xc3\xf9\xb8\x8b\x88\x6d\xdc\xff\x8c\xbc\x70\xc9\xd6\x5e\xea\xda\xfa\x0f\x31\x04\xe3\x12\x7b\x7a\x5f\x5f\x41\xd6\xd2\x19\xb5\xca\x58\xba\xe0\x3e\xfc\xe6\xfc\x69\x3b\xa9\x59\xc0\xe7\xfc\x0a\x08\xfe\x61\x52\x07\xb4\x43\x7c\x6c\xfe\xce\xe1\xf9\x19\xd6\x00\x78\xf1\x2b\x00\x00\xff\xff\x39\x95\x6f\x16\x99\x03\x00\x00")

func syntaxImplicit_importArraiBytes() ([]byte, error) {
	return bindataRead(
		_syntaxImplicit_importArrai,
		"syntax/implicit_import.arrai",
	)
}

func syntaxImplicit_importArrai() (*asset, error) {
	bytes, err := syntaxImplicit_importArraiBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "syntax/implicit_import.arrai", size: 921, mode: os.FileMode(420), modTime: time.Unix(1, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _syntaxStdlibSafeArraiz = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x57\x67\x54\x53\xfb\x9e\x3d\xd2\x4b\x04\xa4\xd9\x10\x2f\x02\xa1\x4a\x54\x5a\x20\xa2\x80\x80\x22\x4d\x9a\x52\xa4\x04\x92\x40\x08\x45\x12\x10\xe9\xc5\x20\x20\x10\x54\x6a\xa8\xd2\x7b\x93\x00\x52\x04\x29\x22\x28\x10\x20\x20\x45\x34\xd2\xab\x74\xa4\x88\x64\xd6\xdc\x3b\x77\x8d\x3a\x6b\xde\xcc\xbc\x79\xe7\xcb\xf9\xf4\xeb\x7b\xaf\xbd\xff\xb7\x74\xe8\x19\x78\x01\x16\x80\x05\xf8\xfd\x03\x01\x00\x60\xef\xe6\x8a\x42\x3b\xc8\xc0\xb1\x58\x38\x7a\xa0\x53\xaf\xbb\x3b\x12\xf2\x8e\x0c\xa9\x34\xba\x75\x93\x2c\xdd\x25\xae\x37\xa0\xd3\x4b\x1e\x32\x90\x96\x1c\xd0\xeb\x66\x30\x8c\x20\xde\x54\x66\xb3\x35\x45\x13\x95\x0a\x6f\xb3\xdd\xdc\xaf\xdb\x93\xc8\x55\x2a\xdc\x87\x9a\x2c\x1d\xbf\x33\x1b\x2f\x72\x7b\xfe\x78\x62\xce\x34\x02\xac\x26\xf1\x69\xb2\x9d\x64\xa2\x12\x8e\x4e\x0e\x02\x00\x80\x46\xbb\xa5\xc3\xcc\xe2\x4f\x89\x0d\xbb\x03\x00\x80\x07\x00\x00\xff\x7d\x43\xa2\x00\x00\xb8\xb8\x21\x3c\x9d\x91\x10\x07\xb4\x87\xa3\xa7\x9d\x8c\xbd\x9b\x0b\x04\x8e\xc5\x9e\x87\xa3\x21\x7f\xb6\x08\x71\x70\x93\x71\x71\x43\xf8\xc5\xf5\x34\xe5\x80\xb9\xde\x0c\x2f\xe4\xdf\x23\x4a\x4f\x67\x73\x5d\xe5\xef\x92\x69\xbe\xfc\xfa\xf2\x11\x3b\xb4\xec\x98\x1d\xdb\x72\xdd\x37\xca\x9c\x03\x4c\xb6\x76\x5d\x7f\xe5\xf1\x1b\xa5\xd4\x7d\x63\xc6\x71\xbf\x5a\xfd\xa4\xa3\x32\xd6\xc3\x09\xe6\x81\x53\x8f\xc4\x15\x5d\x6d\xa2\x1a\xee\x32\x65\x68\x83\x87\x66\xd4\x5f\xea\x76\x80\x7b\x2b\xb6\x8c\xdc\x29\xc2\x7b\xcb\x32\x44\x3b\xe5\xef\x9b\xe5\x22\x4e\x01\xd3\xa5\xc4\xb7\x5e\x9c\x15\x76\x31\xfb\x5e\x9e\xf8\xd1\xd5\x44\xaa\x75\x77\x11\x6f\x20\x79\x9d\x7d\x2c\xb0\xb3\x5c\xe7\x10\x39\x0e\xbe\x26\x40\x95\xd6\x5f\x09\x3e\xe0\x18\x3d\x93\x64\x77\x6e\x14\x01\x26\xad\x2f\xba\x81\xef\xb5\x86\xbe\x68\xb0\x0d\x89\x10\x60\x60\xc9\x54\x8d\xcc\x02\x0a\x36\xd4\x3f\x84\x32\xf4\x27\x9a\xc2\x65\x7b\x65\x15\xbd\x08\x7d\x1d\x72\x8b\xcc\xcb\x97\x38\xae\x73\xa8\x3a\x41\x4f\x51\x62\x03\xbf\x34\xc2\x54\x53\xd1\x1c\x6f\x87\xf0\x4b\x65\xc3\xcd\x09\xc7\x37\xfa\x9b\x63\xdf\x59\xf0\x02\x59\x92\x92\xa5\xb8\xf9\xef\xfc\x30\xe7\x63\xdb\x90\x9d\x00\xcc\x5c\x00\xb5\xe8\xb1\x71\xcd\xab\x2b\x0a\x22\xfa\x8e\xf5\x5d\x3d\x49\xed\x89\xea\xda\x8c\x33\x82\x54\x93\xbe\x73\xbc\x56\x86\xf6\x2f\xfc\xd7\x94\xb8\xbd\xd9\x0e\xcc\xbe\x99\xa7\xe9\x80\xd2\xc7\xe5\x4e\x1e\x1e\xfd\xf8\xe5\x6b\xf8\xfa\x85\x5c\x3d\x2c\xaa\x26\xa0\x35\x85\xe8\x1d\x8e\xd0\x7b\x90\x75\x23\x8f\x7b\xb8\xa7\x1d\xba\xaa\x69\x75\xe1\xa2\x2b\x73\xa9\x6e\x6c\x82\xfb\x2e\xab\xb0\x59\x1b\x22\x6b\xe0\x0e\x89\x53\xb6\xee\x29\xf9\x4a\x7d\x00\xec\x32\x9a\xae\x54\xe8\x09\x67\x70\x6d\x6c\xa5\xbd\x41\x99\xec\xc7\x1d\x5f\x11\xe2\x2e\x35\x92\xeb\x82\xf5\xea\x23\xa6\x66\x7a\x6e\xad\x31\x82\x70\xce\x27\x16\xc2\x07\xd1\x2f\xdd\x62\x57\xfb\xf6\x9d\xcd\xfc\xed\x0e\x76\xaf\x18\x29\x71\x6b\xf9\x8a\xae\x7c\x2c\x77\x23\x67\x06\x30\x6d\xb7\x09\x5c\x14\x98\x8b\x2f\x51\x16\x2a\x1d\x4c\x81\x61\x61\x75\xab\xf9\xf8\xfc\x16\x03\x94\x60\x52\x7e\x70\xb3\xbb\xce\xc7\x0d\x9f\x3d\xd5\xf0\x53\x0c\x37\xcc\x59\xd3\x30\x96\xc1\x85\x98\x21\x93\x81\x54\x55\x6b\xd9\x59\xff\x8d\xab\xf4\x39\x7d\xcf\x35\x0b\x55\x13\xeb\xd0\x2e\x99\x7c\x9d\x4c\x96\x3f\xb4\x83\x14\x31\x8d\xea\x74\xfd\x2c\x40\x21\x5b\x01\xee\x53\xf0\x4a\x94\x9c\xe7\x9d\x02\x83\x4b\xad\xfc\xdc\xe1\x46\xac\x42\x72\xe2\x4f\x23\x9b\x16\x8e\xc6\xe5\x95\x17\x6f\xda\x07\x1e\x0d\xdf\xe5\x8f\xc9\x4b\x6c\x98\x0c\x63\x1d\x8d\x84\xbf\xde\xdb\xd2\xdf\x39\xe7\x5b\x45\xe2\x37\xfd\x9c\x70\x66\xab\xfe\xab\xe6\x1d\xcf\xb9\x73\x67\xb3\x1c\x2b\xc9\x7b\x67\x69\x2d\xcc\x3f\xcc\xf3\xad\xbe\xcf\xd1\xcc\x04\xd3\x91\x02\x49\x29\x7b\x28\x57\x1e\xdf\x1f\x7b\x46\x6a\x48\xaf\x6d\x7d\x31\x4d\xd1\xd3\x8e\xc0\x46\x96\x84\x67\xa6\x5e\xa1\xc6\x10\xd8\x23\xf0\xee\x21\xfb\xdf\xc0\x96\xca\xbf\x4d\x7d\x44\x07\x00\x6b\x8c\xff\x08\xd8\x8a\xff\x33\xb0\x71\xde\xae\x1e\xf0\x07\x10\x9c\x07\xc2\x19\x6d\x07\x41\x39\xc3\xff\x83\x94\x25\x66\x3a\x31\x29\x9d\x3c\xcd\x73\xc6\x83\x36\x0f\x99\x76\x4f\xc7\xfa\x92\xc7\xa2\x48\xa3\x02\xb1\xef\x9d\x20\xf6\xe5\xe8\x94\x9a\x68\x7c\x8b\xd3\x6b\x6d\x8d\xdd\x95\xe2\x45\x66\x4e\x4a\xd0\xf5\x64\xe5\x91\xcb\x45\xc5\x45\x4b\xeb\xc5\x4e\x0a\xde\xc7\x44\xd7\x0f\x79\x00\xd5\x56\xb6\x7b\x2c\x0e\x82\x5c\x1d\xde\x1d\xeb\x34\x2d\x1e\x4e\x3f\xaa\x17\x9c\xf1\xd9\x76\xee\x44\xb4\x63\x00\x4a\xa9\x24\x48\x96\xf6\xc4\xaa\xcc\xbd\x35\x75\x8d\x3a\x73\x89\xe5\x53\x42\xe1\x97\xf5\x99\x48\x54\x9c\xd4\x41\xe3\x38\xac\xa6\x34\x29\x3d\xa8\x8a\xe7\x3b\xeb\xf5\x74\xab\x20\xac\xbe\x3d\xed\xed\x15\x15\x05\xba\xc9\x03\x4b\xef\x51\xd9\x80\x72\x87\x2a\x50\x6a\x90\x73\x16\xe3\xe1\x7e\xca\x3c\xcd\x75\xed\x34\xed\x78\x19\xcf\x93\xda\xe6\x91\x6f\x2e\x75\x23\x4c\x0f\xce\x9a\x65\x2e\xd9\x38\xba\x43\xbd\xba\xf7\x3d\xf9\xf8\x4d\xef\x5b\xaf\x41\x55\x46\x05\x47\xbe\x4b\xeb\xac\x5d\xe5\xdb\x5f\xb6\xc0\xc2\x26\x8f\x2f\x74\x3f\x93\x0a\xee\x11\xd2\x89\x59\x2f\xf4\x2e\x62\xcd\x7d\x77\xa7\x62\xe0\xf8\xc9\x98\x13\x54\x1a\x1d\x8d\xef\x03\x66\x04\x5a\xce\x40\x60\x55\xf1\xcf\x3e\x3e\xdb\xb2\x14\x0d\x16\x8d\x35\x38\xe4\x3d\x58\x51\x15\x1f\x10\xd2\xc9\xe6\xe0\x31\xd1\x92\x11\x09\xfb\x52\x73\x99\xe5\x06\x91\xc2\x23\x5c\x13\x3d\x2c\x17\x81\x2e\x2b\x13\xbe\x9f\x83\x34\x0a\x87\x45\x49\x1e\x09\x23\xf1\x77\x21\x05\x3b\x45\x05\xc5\xdd\x7d\xe2\x97\xa2\x9d\xb4\x3a\x60\xe6\xd8\x52\x13\x8e\x02\x83\x31\x11\x99\xea\x93\xb6\x1f\x85\xd5\x97\x5a\x28\x89\x2d\x3c\x52\x72\x32\x3f\x02\xce\x4f\x4a\xf7\x74\x15\x54\xde\xce\x6d\xad\xb2\x8e\xe9\xea\x23\xdd\xb5\x88\x2f\x57\x34\xb5\x8c\x87\x7c\xa3\x27\x93\x6d\x43\x05\xa5\xb7\x1f\x2a\xb7\x7b\xbd\xb4\xdf\xa5\xde\x14\xe5\xbe\x34\x6a\xdb\x66\x77\x64\x52\xcf\x25\x9a\x59\x3f\xb6\x5c\x80\x62\x8e\x09\xd5\x9f\x3a\x5f\x3c\x00\xaa\x86\xb9\xec\xf4\xf9\x83\x5f\x48\xda\x30\x20\x93\x61\x22\x08\x1e\x4e\xf3\x56\x85\x7b\x77\xde\xe1\x3a\x7a\x72\x83\x91\x95\xb1\x7c\xc7\x16\x5e\x14\x48\x52\x77\x30\x44\x36\x18\x1d\x9d\xb3\x80\x94\x89\x1b\x04\x69\xfd\x3a\xfd\x48\xa8\xf8\x15\x5b\x71\x09\xda\x8c\x5a\x32\x9f\xde\x97\xc5\xbe\x52\x1b\x26\x0b\x82\x44\x89\x2d\x2e\x9b\xd5\xb9\x37\x4a\xab\x83\xce\x45\x5a\x8d\xa3\xa0\xe2\x0f\x2f\x17\xd1\x95\xf1\x18\x92\xf6\x52\x9c\x04\x85\x82\xd9\x52\x3c\x15\xfb\xed\x49\x4b\x7f\x68\x78\xee\x22\x41\x8c\x04\x4c\x6d\x24\xbf\x66\x73\x95\x7b\x2a\xb1\x10\x6e\xa3\xf8\xc7\xf5\xe4\x69\x67\x47\x22\x18\x3e\xe4\x81\xca\x9a\xff\x20\x8c\x85\xd5\x98\x76\x75\xd7\xce\x1c\x9a\x71\x24\x5e\x92\x98\xc4\xb4\x32\xdc\x86\xf7\x58\xf2\x95\x69\x10\xe4\x5a\x61\xe3\x77\x5d\xc2\x78\x2e\x84\x90\x84\x6f\x49\x0c\x72\xb1\x77\xd4\xf5\xa0\x31\xb3\xd7\x0b\x99\x0c\xb7\xdf\xe7\xb7\x54\x89\x41\x29\x72\x37\x93\x16\x8f\xf1\xdc\xa3\x86\x60\x21\x4d\x4d\x51\x54\xb3\x24\xf9\x54\x85\x13\x85\xc1\x2f\x63\x6a\x04\x53\xd1\x87\x9a\x89\x01\x2d\x3a\x0a\x93\xd0\x7c\x31\xd4\x6e\xa3\x94\x96\x48\x5b\xb8\x40\xc4\xc4\x35\x9f\xf2\xd4\x34\x6d\x16\xd5\x58\xd5\x79\xf9\x68\x06\x77\x9f\x6c\x6d\x01\x21\x52\xe6\x9c\x93\x31\x36\xc9\x6e\xb0\x18\x59\x95\x01\x6b\x31\xec\x3c\x60\x7e\x2a\xe4\x82\x33\xbf\x3c\x5c\x9d\x96\x0b\x7a\x21\x6c\x16\x1f\xe5\x4b\x69\xc0\x58\xcd\xd3\xbd\x69\x83\xc7\xa6\xeb\x9e\xeb\x2d\x8b\xe2\x97\xf6\xb9\x1c\x53\xef\x15\xb9\x3d\x85\xce\x95\xaf\x66\xcd\xcd\x88\x62\x7e\xe7\x1d\xd6\x5c\x79\x0a\x9e\xab\x7b\x71\x0d\x5a\xe4\x74\x9a\x45\x7d\x89\x2f\xc8\xf3\x9b\xc8\xc9\x15\x05\x89\x8a\x4c\x18\xdb\x8c\x27\x95\x7c\xed\xa5\x1b\xaa\x94\x9c\x21\x45\xc3\x99\xb7\x6b\x7f\xc5\xf1\xc7\x33\xcf\x82\x74\x61\x79\xd6\xf6\xd6\xa6\x54\xa2\xf1\x50\x99\xe2\x89\x91\x7d\x57\xf7\xd5\x9a\xef\x0c\x81\x81\xb4\x8d\xb4\xf0\x08\x2e\x19\xe9\xb7\x60\xe5\x0b\xb7\x3c\x34\x9d\xb5\xe3\x32\x36\x88\x2f\x85\x53\x04\x65\xa3\x2a\x93\x83\xd1\xf2\x66\x91\x42\xc3\xfd\x21\x60\x65\x26\x8f\xbc\x4d\xf6\x48\xdb\x21\x22\xe5\xe3\xf6\x94\x15\xfd\xbe\x92\xd0\xab\xfc\x53\x31\xcf\x16\x98\x34\xac\x7a\x22\xeb\x52\x52\x1a\x46\x27\x9e\xb3\x29\xfe\xd1\x1b\x25\x1f\x3d\x54\x1f\x71\x13\x6a\x4c\xb7\x0e\xe6\x32\x2e\x2f\x8b\x98\x8d\x0f\x78\xae\x3c\xa3\xe1\x75\x2a\xf0\x61\xfc\xe0\xe2\x75\x79\x23\x5f\x9b\x89\x38\x3f\x4a\x23\xa6\xfb\x64\x2f\xc7\xa5\xaa\xa6\x73\x51\x45\xcb\xdc\xf6\xe4\x07\x9f\x59\x12\xd6\xa7\xbd\xaf\x51\x18\x9d\x78\xe5\xf4\x6b\x36\xcd\xed\x93\x36\x2d\x4b\x16\x84\x92\xd6\xde\xe9\x4e\xa7\x07\x69\x19\xca\xf0\x07\x6b\xbc\x3c\xf7\x70\xb2\x3b\x48\x90\x63\x26\x57\x97\xc1\xe7\x31\x91\x6e\x03\x44\x88\x23\x8f\xed\x08\xcb\xf4\xce\xa7\x07\xe2\xe5\x94\x78\x13\x66\xeb\xeb\x6d\x7c\x8c\x75\x10\xaf\xaa\x2b\x29\x49\x43\x85\x84\xa0\x46\x17\x77\xf5\x40\x1e\x4c\x8d\x9d\x33\x6f\xe2\x23\xa5\xb5\xf4\x7b\x37\xad\x4a\x90\x03\x1e\x5c\x86\x76\xf0\x14\x6d\xa3\xd4\x3d\xb3\x02\xb4\x05\x71\x24\xcc\x11\xb7\x99\xcf\xe9\xed\x2a\xe4\xe2\xa3\x94\xc4\x69\x4d\x4d\x23\x24\x69\xa0\x50\x50\x21\xbd\x33\x25\xf9\x14\xed\x58\x52\x6c\x47\x17\xfd\x7d\xa7\x67\x84\xe3\x11\xef\xcd\xad\xbb\x45\x56\x54\xaa\x1b\xda\xdf\xec\x34\xf0\x66\x7d\x42\x93\xbd\xd0\x9f\xcf\x26\x2b\x87\xb1\x2f\x08\x11\xf7\x0c\x3b\xb0\x17\x1d\x52\x17\x8e\x92\x4f\x2e\xb2\xd8\x5a\xea\x29\x77\xc7\xf4\x26\x83\x51\x28\xbd\x7d\x61\x66\xef\xda\x39\x66\xeb\x8c\x25\xbd\x9d\x9c\xab\xf3\xe2\x67\xbd\x36\xbf\xc2\x4f\xde\xc8\xab\xb8\x58\x92\x63\xa9\x18\xb2\x59\xd6\xb7\xc5\xf1\x41\x1f\x24\x6a\xdb\x7e\xe9\x3e\x6f\x8e\xa5\xe6\x38\x62\x3d\x00\xd1\xf7\xbc\xdc\x62\x71\x7f\x64\xcb\xfb\xf5\xe3\x4f\x5c\x42\x72\xa9\xf0\x6c\x43\xfa\x3f\x6a\x04\x18\xae\x9b\x34\xa7\xb1\xb7\xed\x32\x87\x77\xa5\xc3\x3e\x6c\x09\xee\x3d\x18\xd4\x75\x84\x1e\x90\xbe\xab\x35\xa8\xbc\x7a\x5a\x5b\x9a\xcf\x59\x65\xfa\xf5\xd4\x58\x92\x12\x3c\x22\x4e\xca\xcf\xb6\xc6\x86\xd5\xd4\x8f\x7b\x2f\x30\x66\xc5\x44\xee\x46\xff\x31\x51\xb1\xbc\xc6\x43\xbc\x82\x80\xb4\xfc\x04\xd3\x28\x19\x0a\x99\xdc\xda\x32\x2d\xb1\xdc\xe6\x2b\xfb\x10\x84\x3b\x18\xf1\x71\x80\x46\xb6\xc4\xc8\x4e\x81\x06\xd5\xda\xdd\x38\xf1\x6b\x07\xa0\xc3\x62\x28\x7c\xdb\xa4\xbe\xc2\x5a\xf1\x6c\x35\x53\xf2\x8a\xd8\xc1\xe7\x52\xe1\x57\xe2\xde\xe3\x06\x34\xc6\xbf\x75\x40\xd2\xff\xfd\x4e\x36\x23\x00\x14\xf0\xff\x23\x1d\xb8\xf2\x7f\xd5\x81\xbf\x7e\xe7\x71\x70\x14\xf2\x2f\x39\xc0\x3f\x31\x20\xa7\x5c\x00\x75\xac\xfa\x3d\xbe\xee\x03\x79\x7e\x7b\x21\x0d\x3d\x54\xd6\xb0\xda\xf9\x86\x8b\xbe\x7f\x80\x41\x34\x2f\x90\xe4\x9a\xe1\xcf\xe6\xc1\x6b\x81\x2e\xe7\x7c\x3a\x35\xb7\xf7\xe6\xfe\x95\xb5\x6c\xc5\x63\x95\xac\x28\x2e\x50\x4b\x46\xf4\x2c\x4b\x6b\xa8\xab\xa3\xbf\xea\x96\x98\x16\xac\x5e\x2b\xca\xeb\x58\x87\x9a\x8f\x16\x2f\x36\x39\xe4\x91\xbc\xe7\xf3\x7b\x0f\x1f\x0c\x46\xe9\x5c\x38\xb7\xfc\x2c\x52\xfd\x16\x82\xb7\x89\x1d\x1f\x80\x1d\xb7\x1d\x7d\x4b\x57\xf6\xa4\xd7\x12\x7b\x46\xbc\x3a\xd1\xbd\x02\x84\x3d\x15\x8c\x04\xc9\x84\xa4\x26\xa6\xb7\x5d\xab\x3b\xb1\xa0\x66\x91\x22\xc7\x3f\x3a\x84\x5f\x95\x45\x98\xf2\xba\xe8\xa5\x1d\x09\x42\xaa\x1d\xf6\x40\x27\xbe\x74\x26\xec\x26\x7e\xcd\x8c\x93\x5c\xe3\xab\xda\x79\x8f\x48\xc7\x6d\xcf\x35\x01\x5e\xcb\x84\x4e\xb1\x8b\x9a\x98\x42\x29\x77\xe8\x81\x51\x58\xdc\x61\xaf\x43\xd6\x08\x84\xe3\x7e\xf4\x70\x87\x67\x83\x7d\xdd\xa4\xd1\x5d\x83\x6f\xd5\x71\xea\xb8\x29\xf9\x30\x6e\xd3\x29\x2b\x8c\xf2\x34\x2e\x57\xc3\x94\x8a\x4f\x0c\x1e\xda\x46\x0c\x49\xcf\x54\xe1\x62\x10\xe2\xb7\x7b\x6a\xdd\x0e\x57\x43\x1b\x1f\xf4\xad\xae\x91\xd4\x42\xa2\x4f\x6b\x8d\x0e\x54\x25\x18\xdc\xac\x14\xca\xa8\xdd\x76\x48\xed\x6a\x0a\xd4\x28\x22\x88\x04\xbc\xc2\x15\xc5\x17\xec\x2a\xe3\x1d\xde\x9a\xa4\x12\x90\x13\x39\xcd\x32\x07\xac\xc9\x0a\x69\x15\x73\x41\xa3\x1e\x07\x47\xff\x3e\x5b\xed\x8e\xba\xa1\xce\x11\x00\x60\xa7\xff\x97\xca\xb7\xa7\x07\xda\xf9\xaf\x7b\x39\x13\x5a\x1e\xe3\x2f\x70\x69\xac\xfe\xe0\x63\xf4\x35\x7c\xb2\x6c\x35\x0a\x31\x57\x21\xa5\x31\xd0\xf5\xab\x8a\x99\xe3\x91\x87\xfe\xe1\x97\x64\x41\xba\x56\xf5\xde\xdb\x1c\xdb\xa1\xe7\x7f\x28\xb1\x97\xa3\x1f\xe6\xbc\x72\x34\x6b\xc6\xb0\xf6\x57\x48\x32\xc7\xe1\x93\xcf\xf8\x19\x4d\xdb\xd9\x9b\xf5\x88\x15\xa3\x12\xcf\xaa\x68\xde\x79\xde\x7b\xdb\x7d\xb7\x65\x99\x15\x90\x01\xf7\x3d\x32\x63\x90\x2d\x4e\xd7\x8b\xcc\xf9\x14\x62\xba\x4e\xd3\x38\xdf\xfc\x9e\x85\xff\xe4\xe7\x89\x8b\xaf\x4f\x18\x43\xdb\x91\xf8\xfb\x4d\x18\x36\x42\x80\xaa\xfd\x8f\xf7\x4d\x69\x54\x8d\x39\xc3\xe4\xf0\xfa\xf5\xb1\x21\x8a\xca\xdf\xd3\x87\x81\x95\xd3\x89\x00\x00\xec\xfc\xe9\xca\x8f\xd0\xf1\x02\xff\x39\xff\xcf\x8e\x1d\xf4\x5f\x36\xf2\xf3\xb3\xe1\xf7\xc8\x9f\x2d\x91\xe8\x2f\x51\x4f\xff\xb7\xfe\xfe\xf7\x94\x3f\xb3\x4b\xf1\x97\x94\x4e\xf4\xff\xb4\xb3\xfa\xbd\xc8\xcf\x58\xb8\xf2\x4b\x11\x2c\xeb\xff\x97\xb6\xbf\xd7\xfa\x79\xf3\xbf\x0e\x24\xc1\xfe\x4f\x63\xed\x96\x0e\x23\xd3\xbf\xa7\x60\x04\x18\x81\xb6\x23\x00\x20\xf6\xe7\xdd\xfe\x2d\x00\x00\xff\xff\x46\xed\x37\x4a\x00\x0e\x00\x00")

func syntaxStdlibSafeArraizBytes() ([]byte, error) {
	return bindataRead(
		_syntaxStdlibSafeArraiz,
		"syntax/stdlib-safe.arraiz",
	)
}

func syntaxStdlibSafeArraiz() (*asset, error) {
	bytes, err := syntaxStdlibSafeArraizBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "syntax/stdlib-safe.arraiz", size: 3584, mode: os.FileMode(420), modTime: time.Unix(1, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _syntaxStdlibUnsafeArraiz = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x0a\xf0\x66\x66\x11\x61\xe0\x60\xe0\x60\x40\x07\x3c\x0c\x0c\x0c\xc9\xf9\x79\x69\x99\xe9\x7a\x89\x45\x45\x89\x99\x97\x4e\xf8\x9e\x39\xd3\xad\x7f\xea\xbc\xfe\xc6\xa0\x00\xaf\xf3\x3a\x27\x35\x7c\x2f\x79\x9f\x3b\x7f\xd5\x5f\x47\xeb\x92\xef\x19\x96\xc0\xce\x99\x5e\x56\x5c\x09\xa1\x99\x33\x2d\x57\x86\x71\x79\xfd\xda\xf9\x53\x73\xa9\xe5\xca\x5f\x16\x21\xaf\xc5\xc3\x9f\x4d\x55\x09\x7b\x21\x3e\x7d\xc9\x93\x14\x35\x47\xcd\xbb\x97\xf6\x3c\x3a\xb2\x25\xc4\xb6\xa3\x72\x56\x03\x03\x03\xc3\xff\xff\x01\xde\xec\x1c\xfa\x49\xa6\xa6\x91\x0c\x0c\x0c\x65\x0c\x0c\x0c\xb8\x9d\xa4\xca\xc0\xc0\x90\x9b\x9f\x52\x9a\x93\xaa\x9f\x9e\x59\x92\x51\x9a\xa4\x97\x9c\x9f\xab\x9f\x58\x54\xa4\x9b\x98\xa9\x0f\x76\xa4\x7e\x7a\xbe\x5e\x6e\x7e\x4a\xcd\x94\xb3\xfb\x96\xa8\x09\x1c\xbd\xfe\x72\x79\xc1\x4c\x9d\x27\x8b\x05\xec\xc5\x4e\xea\xed\xb7\x39\x60\xc3\x98\x94\x69\x7c\x2b\x89\xeb\xcd\xce\x6f\x97\x9f\xa7\x5b\x1b\x6f\xff\xe8\xf7\xae\xeb\xa8\xe5\x9c\x5f\xc1\xac\x77\x6a\xb6\xfb\xcd\xe0\xd5\x8b\xbb\x3e\x2d\xb2\xfe\x71\x9b\x86\x79\x5e\x7c\xcf\xee\x18\xb6\xf9\x9e\x6a\x57\x9f\x3a\xed\xf0\x39\xa6\x76\x6e\xc3\x97\xa0\xc2\xcb\xca\x3f\xdf\xe8\xcd\x4c\xb2\xfa\xfd\x79\xbd\x4a\x56\xdd\x93\xb5\x33\x8f\x97\xf3\x6f\x48\xea\xff\x55\x5e\xda\x72\xf3\xfd\xf4\xfb\x71\x67\x56\x89\xd4\x9f\xff\xc8\x7d\xab\xfe\xc4\x7a\xef\x7f\xa9\x77\xd4\x9c\x65\xee\xeb\xf8\xbd\x6b\xfc\xc3\x77\x53\x76\x46\x92\xd2\xcd\x14\xb5\x2d\x1f\x5f\xe5\xab\x15\x1c\x6a\xdd\xb4\x3b\xa1\xa9\x53\x86\x85\x63\xa1\x43\xf7\x22\x86\x15\x9f\x9c\xae\xb5\xb2\x5c\x9c\x1e\x9a\x68\x7c\xce\xd8\xbc\xbc\xef\xc2\x31\x93\x57\xec\x6f\x8c\xf8\xdc\xf9\x1c\xb2\x2c\xa4\x2e\x4f\xae\x7f\xb0\xc7\xda\x61\x4e\x26\xdf\xf1\xab\x2d\xaf\xd7\x5d\xdf\x3f\x4d\xfc\xd3\xc5\xfd\x93\x4f\x45\x89\x30\x2c\xd2\xd2\x5a\x5b\xfc\xe2\xb7\x98\x75\x8e\xe0\x57\xfd\xef\x75\xd9\xcf\xeb\xee\xaf\xea\x0a\xde\xb6\xd7\xce\x4c\xc5\x2f\x63\xd7\xc9\xb3\x33\x8e\x4c\x77\xf2\x64\x7d\x2a\x77\x3f\xe4\x82\x92\x48\x6c\x60\xf2\xa6\xda\x0f\x96\x42\x95\x5c\x7f\x22\xbe\x45\xce\xf5\xe6\x99\x77\xc7\x44\xf2\x1f\xef\xed\x07\x6f\x3b\x3e\x1a\x2c\xf5\x2d\x4a\xdb\x56\x77\x68\xf6\xcc\xca\x8e\x14\xdf\x8a\x45\x1e\xcb\x84\xae\x9f\x3d\x62\xf1\xde\x35\xd6\xc0\x30\x8f\x7d\xad\xcf\xe4\x69\x85\x3f\x38\x95\x23\x0e\xa7\x2c\xba\x14\xbe\x85\xdf\x78\xe7\xc4\xf3\x76\xbb\xea\xac\x6d\x32\x99\xd6\x2a\x4e\xe0\x6f\xdc\x3e\x79\x63\xb2\xff\x3a\xe3\xdb\xdf\xab\x55\x66\xfe\xb8\xdf\x2d\x60\x10\xf7\xbe\x8d\x6d\x3f\xb3\x90\xdb\xad\x3e\xe5\x25\x77\x39\xfa\xae\xa9\x3e\x38\xa3\x6e\x7f\xe1\x57\x4e\x44\x6d\xd2\x9f\x1f\x76\x41\x96\x42\x6e\xd5\xaa\xef\x6e\xaf\xcf\x3f\xbf\xb0\x8e\xed\xeb\x61\x19\x43\x99\xe7\x53\xd7\x58\x29\xae\xbd\x32\xdb\xba\xc8\x7a\xe7\xfb\xe5\x2d\xcb\x0f\xfa\xa7\xc9\xcd\x58\xde\xb8\xbf\xd0\xfb\xf6\xa7\xaa\x9f\x0e\x1d\x52\x2c\x1e\x91\x9c\x73\xb3\xa3\x1b\x57\x66\x5f\x0d\xb9\x34\xc7\x21\xce\xf8\x59\xed\x27\x7b\xe6\x25\x17\x16\xb8\xae\x74\x98\xbe\x33\x33\x77\xa1\xe8\x09\xb6\xe8\xbf\x9e\x0d\xe6\xd9\x7b\x9c\x98\x2e\x72\x30\xac\xe4\x5a\x51\x7c\xb7\xf1\x5d\x8f\x49\x69\xf8\x0a\x7f\xa3\x43\x62\x42\x1d\x41\x9c\x8a\x26\x1a\x13\xbb\xf7\xbd\xe4\x9d\xb2\x6c\xfd\xea\xcf\xc9\xf5\xbc\x1d\x3f\xc4\xfa\x97\x4d\xdf\xfd\xa8\x9d\xf3\x66\x77\xe2\x81\x9f\x5f\xfc\xbe\x2b\x55\x6f\xde\x22\x16\x7a\x6f\x9a\xec\x97\x5d\x6f\x5d\xc3\x4b\x9f\x2b\xc9\x2f\xca\xd8\x78\xfe\xa7\xfc\xff\x83\xec\x7f\x23\x97\xc7\xfe\x7e\xfe\x3f\x42\x6e\x5e\xaa\xcc\x8c\xd9\x3f\xd3\xf2\x84\xab\xff\xfe\x0c\x72\x4c\x2d\xff\xea\xa7\xee\xaa\x2a\x9d\xc1\xf0\x69\x91\x66\xe9\x42\xdf\x95\x2e\x57\xd5\x4a\xea\x63\xfe\x71\xc3\x12\xb6\xf6\xf2\xb0\xfb\x6d\x4c\x0c\x0c\x1f\x58\xf1\x25\x6c\x07\xc2\x09\xbb\xb8\x32\xaf\x24\xb1\x42\xbf\xb8\x24\x25\x27\x33\x09\x4a\xe9\x96\xe6\x15\x27\xa6\xa5\x42\xf2\x67\x4a\x9f\x63\x57\x4b\x00\x8f\xcb\xf7\x79\x8b\x0e\xec\xfd\x3e\xe3\xaf\x4e\xba\xf0\xdf\x8b\xf2\xcf\xca\x95\x3c\x73\x3e\x7e\xeb\x55\xe8\xf5\x7d\xcf\xb5\x4a\xee\x6b\xd5\x9f\x30\xce\x23\x4f\x36\xa9\x85\x6d\x5f\x2e\xc2\x27\x50\x66\x65\x30\x5b\xf9\xcd\xc1\x19\xe6\x73\x83\x7c\x98\x17\x9d\x38\xd5\xa1\xa2\xfb\x84\xd5\xd7\x4d\xdd\x7e\x42\xdb\xc4\x8d\xbb\x53\x36\xcd\xd3\x2e\x92\xf7\x60\xfb\xa3\x32\xe9\x9f\xa5\xdd\x86\x75\xf5\xcd\xe6\x29\x1b\xdd\x76\x35\x9c\xfb\x7d\xe7\xda\x31\x05\xdd\x82\x73\x67\x4d\x1d\xb2\x22\x6f\xae\x4f\x9b\xbe\x66\x41\xd0\x93\x2b\x8b\xca\x8c\xff\x94\xbd\xf8\xa7\xe4\x3d\x79\x23\x73\x83\xf6\x3e\x26\x58\x10\x4c\x9e\xbb\x62\xe3\x54\x06\x06\x86\x1f\xe0\xbc\xcd\xc8\x24\xc2\x80\x08\x04\xe4\x7c\xcf\x83\x11\x2c\xc8\xc5\x0f\xba\x4e\xe4\x80\x55\x45\xd1\x35\x99\xd8\x52\x02\xdd\x48\x64\x87\x3a\xa0\x18\x99\xc3\x4c\x79\xfc\x04\x78\xb3\xb2\x81\xcc\x62\x66\x60\x66\xf8\xcd\xc0\xc0\x90\xcf\x02\xe2\x01\x02\x00\x00\xff\xff\x4c\x64\x2a\x6e\x80\x05\x00\x00")

func syntaxStdlibUnsafeArraizBytes() ([]byte, error) {
	return bindataRead(
		_syntaxStdlibUnsafeArraiz,
		"syntax/stdlib-unsafe.arraiz",
	)
}

func syntaxStdlibUnsafeArraiz() (*asset, error) {
	bytes, err := syntaxStdlibUnsafeArraizBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "syntax/stdlib-unsafe.arraiz", size: 1408, mode: os.FileMode(420), modTime: time.Unix(1, 0)}
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
	"syntax/implicit_import.arrai": syntaxImplicit_importArrai,
	"syntax/stdlib-safe.arraiz":    syntaxStdlibSafeArraiz,
	"syntax/stdlib-unsafe.arraiz":  syntaxStdlibUnsafeArraiz,
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
	"syntax": {nil, map[string]*bintree{
		"implicit_import.arrai": {syntaxImplicit_importArrai, map[string]*bintree{}},
		"stdlib-safe.arraiz":    {syntaxStdlibSafeArraiz, map[string]*bintree{}},
		"stdlib-unsafe.arraiz":  {syntaxStdlibUnsafeArraiz, map[string]*bintree{}},
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

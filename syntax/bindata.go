// Code generated by go-bindata. DO NOT EDIT.
// sources:
// syntax/implicit_import.arrai (921B)
// syntax/stdlib-safe.arraiz (3.454kB)
// syntax/stdlib-unsafe.arraiz (1.409kB)

package syntax

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
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
		return nil, fmt.Errorf("read %q: %w", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
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

	info := bindataFileInfo{name: "syntax/implicit_import.arrai", size: 921, mode: os.FileMode(0644), modTime: time.Unix(1, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xc5, 0x61, 0xe5, 0x91, 0x8e, 0xa9, 0x63, 0x16, 0x96, 0x63, 0xc7, 0xb9, 0x92, 0x57, 0x54, 0xf6, 0x8, 0xc, 0xca, 0xf0, 0xea, 0x12, 0x5a, 0xfa, 0x4a, 0xb8, 0xfc, 0xa8, 0x37, 0x95, 0x6f, 0x1f}}
	return a, nil
}

var _syntaxStdlibSafeArraiz = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x56\x79\x3c\xd4\x6b\xdf\xfe\x09\xc3\x58\x4e\x64\xa9\x14\x49\x59\x63\x46\x68\x23\x6b\xa8\x66\xc6\x9a\x3d\x2a\xcc\x18\xc3\x30\x32\xd6\x2c\x8d\x65\x90\xac\xd9\xd5\x20\x21\x4b\x0b\xd9\x8a\x89\x83\x42\x86\x19\x59\xb2\x14\xb2\x65\x1b\x0d\xd9\x19\xbc\x9f\x73\xce\x7b\xde\xb7\x7a\x3f\xef\x79\x9e\xe7\x3c\xcf\xfd\xcf\xfd\xb9\xff\xf8\xee\xd7\xf7\xba\x2f\x63\x38\x2b\x9b\x20\xc0\x09\x70\x02\x3f\x1f\x1e\x00\x00\x1c\x71\xee\x4e\x18\x34\xc4\xde\xd3\xd3\x1e\xd3\xd5\x6a\xd0\xde\x7e\x17\xda\x46\x83\x96\x99\x1a\xc3\x68\xf2\xef\x64\x0c\xba\xe0\x54\x5a\xaf\x91\xbc\x5c\x97\x41\x3b\x9b\xc9\x9d\x0c\x98\x2a\xd7\x0d\x73\x4c\xc6\xb9\x22\x0b\x2e\xd8\xe6\xab\x0d\xd9\xfc\x73\x45\x9b\x67\xcd\xe6\xf6\x5b\x7e\x49\x39\x6e\x31\xbd\x3f\x2d\x6f\x02\x29\xa5\x2d\x3b\x34\xf6\xa6\xc2\x4c\x3d\x0a\x93\x49\x00\x00\x60\x77\xd7\x18\xce\xc1\x19\xd4\x9d\x14\x69\x09\x00\x80\x17\x00\x00\xff\x7f\x42\x92\x00\x00\xb8\xe1\x90\xde\x58\x14\x14\x8d\xf1\x72\xf6\x76\x80\x38\xe2\xdc\xa0\xf6\x9e\x9e\x0a\xf6\x18\xe8\xef\x29\x42\xd1\x38\x88\x1b\x0e\x19\x98\x0c\xaf\x19\x97\xe2\xd3\x6f\x59\x94\x53\x41\x80\xc6\x4c\x08\x62\x61\xc8\x54\x28\xf7\xb4\x90\xb8\xd9\x3d\x0c\x45\x47\x36\x61\x6b\xb0\x84\x26\x5d\x52\xfb\xfc\xfd\xfc\x27\xfe\xb0\xcd\x73\x87\xa3\x5b\x4e\x2d\xcf\x85\xbf\x9d\xee\xf6\x9f\x4f\xe7\xb6\x74\xbf\x03\x9b\xab\x6e\xd8\xec\xc6\xf0\x10\x91\xa6\x24\xc3\x70\x0a\xd2\x2e\x35\x86\x22\x2e\x54\xcd\x6a\x2a\xe3\x7e\x4f\xd6\x29\x5a\x73\xb7\xcd\x78\xe9\x11\x9d\x1e\xe8\x4f\x03\xba\x0f\xcd\x8e\xf4\x0a\x08\xd7\x67\x07\x07\x80\x4b\x60\x4c\xf5\xc7\xd0\xe9\x9d\x87\xa5\x90\xc9\x34\x4a\xd0\x8b\x37\x8b\x1d\xa3\x13\x47\x5a\x71\xc6\x9f\x08\xb2\x52\x2f\x2e\xfb\xe8\xa4\x51\x47\x5d\x8e\x49\xdb\x16\xbe\x04\x4b\xa5\xcb\x68\x13\xf8\xb8\x62\x3a\x45\x72\x20\x95\xd8\x0a\xb6\x93\x66\x95\x3e\xec\x53\xc4\xec\x00\xa3\x69\xd3\xb8\xcd\xa3\xb5\xb6\xcd\x24\xd6\x04\x8e\xd5\x66\x1e\x41\xb1\xb3\x38\xf0\x80\x98\xdb\xea\x9b\x81\xe6\x25\x0b\x5c\xb0\xdd\x35\x99\xbc\x05\x77\xb9\x43\x0f\x59\x74\x9f\xba\x3a\xc2\x0e\x44\xfc\x12\x40\xe7\xeb\x80\x4e\x05\xcf\x7d\xb1\xbf\x5f\xbc\x1c\x32\x4f\x8d\x57\xdb\x27\x5d\xb1\x3d\x81\x6e\x5f\x48\x31\x36\xd5\xae\x1c\xfd\x28\xa2\xca\x36\x34\x07\x2a\x9f\x71\x94\x49\x94\x31\x9e\xd8\x2a\x25\xa5\x20\xbf\x25\x1c\xd5\x34\x48\xbd\xdd\x04\x01\xf9\x36\x54\x10\x55\x28\xf7\x27\xe6\xe9\x8d\xe1\x96\x6a\x4b\xe9\xad\x03\xeb\x66\x89\x08\x54\x51\x89\x01\x66\x58\x78\xe5\x7d\x81\xd7\x88\x1b\xca\xd0\x06\x32\x6e\x74\xb0\x61\x30\x8c\x27\x48\x6f\xad\x34\xa9\x7e\xe1\xdc\x71\xcd\x4f\x6b\x08\xc3\x9c\x53\xfa\x5c\x27\xbf\xb9\xbf\xf8\xaa\x94\xc0\x98\xdf\xde\x52\xbf\x87\x60\x32\xf4\xa2\xd2\xa9\x67\xe0\x8e\xfd\x6c\x60\x81\xa1\x35\x85\x09\xde\xf8\xac\x0f\x92\x63\xed\x8f\x15\xbb\x36\xb1\x56\x78\x0f\xe6\xba\xc6\x37\x4f\x09\xc4\x57\xc4\x79\xaf\xeb\x23\x59\x91\x9b\x19\x3d\xf1\x1e\x8d\x2f\x19\x96\x2f\x36\x73\x2a\xba\x6e\xe0\x5a\x87\x26\xea\x19\xfb\xb2\x42\x14\x07\xfc\x91\x21\xfa\x3b\xcd\xa2\x43\x69\x7d\x14\x0e\x1d\x15\x71\xad\x36\x9b\xf4\xae\x97\x22\xd8\xee\xc2\xaa\xd4\x55\xf1\x6f\x6d\xe9\xd5\x49\xbb\xe0\x84\x46\x02\x7f\xca\x9e\x2b\x29\x55\xf8\x41\x99\x04\x05\x9f\x2f\x02\x5a\xa5\xae\xaa\x17\xa4\xde\x73\x01\x45\x2a\x2a\xf9\x33\x0d\xcb\x16\x83\x05\x7d\xb3\x02\xfa\xe4\x16\x44\xae\xa9\x0e\x27\x2b\x2c\xd9\xc1\x99\x83\x3b\x45\x34\x73\x73\xbd\xff\x3a\x0b\x78\xa4\x6b\x78\xd9\xa5\x24\xfc\xa6\xce\xea\xbe\x03\x37\x35\xc8\xdd\x6e\x89\x4f\xad\x9b\xa3\x6a\xe8\xc5\x62\xa0\x6e\xde\xdc\xe7\x5b\xe3\x12\x47\x8c\x2b\x5e\x4e\x6d\x30\x76\x1b\x22\x37\xf1\x24\x77\xd1\xac\x5d\xc7\x4a\x74\x2a\xa5\x95\xd9\xce\x58\x55\x2d\xeb\xcb\x86\x42\x0b\xc7\x0a\x86\xe7\xf3\x0a\x51\xdc\x0a\x5a\x45\x1f\xdb\x9c\x1d\x4f\x48\xce\x42\xca\xb7\xc6\x97\xaf\xff\x0f\xb2\xf5\x1a\xa3\xaf\x45\xee\x01\x80\x05\xf6\xbf\x42\xf6\x99\x7f\x8c\x6c\xbc\xbf\xbb\x97\xbd\x1f\x14\xef\x85\xc4\x62\x1c\xa0\x4e\x58\xfb\xff\xde\xca\x27\x56\x1d\xf8\xdc\x4b\x82\xf5\xb4\xeb\x8b\x05\x11\x5c\x1d\xd1\xa4\x00\x74\x8f\x99\x9b\x4f\x5c\x12\xe5\xc3\x35\x07\x0b\x1b\x4c\x8a\x1c\x5b\x09\x8b\xb2\x6c\xc4\x6e\xd0\xdc\x0a\x47\x0e\x21\xfd\x29\x39\x4c\x82\x1c\x7c\xbd\x7e\x8d\x48\x91\xdf\x3b\xee\xbc\xcd\x0e\x88\x0b\x9f\xe3\x93\x22\xee\x45\x5c\xca\x80\x31\x35\x16\xfb\x39\x60\x45\x1c\xca\x9d\x4f\xf6\x9e\xc2\xb5\x96\xbf\x6e\xbd\x01\x32\xed\x34\xb8\x75\x56\xbf\x4c\xbc\xa1\xbe\x7b\x46\x62\x41\xa4\xc0\x56\x3f\x52\x9a\x52\x68\xd4\xf9\xc4\x4a\xa0\xef\x14\x62\x44\x2b\x9a\xb4\xce\x7e\x35\x43\x8c\x00\x7f\x77\x6c\xb3\xca\xb0\xaa\x92\x20\xb3\x65\x7d\xf8\x83\xce\xb5\x59\x6b\xc3\xfd\xc8\x86\xd4\x51\x36\xcd\x19\x25\x66\xd3\x55\x66\x33\xe3\xce\x45\x61\xfa\x24\x74\x71\x53\x89\x44\x86\x2c\xac\x35\x82\xcf\xdf\xfa\xbc\xe1\x34\xa2\x34\xb7\x5d\x28\x2f\x37\x24\x9f\x66\xf2\x02\x1f\xe4\xd1\x87\x11\x19\xf7\xcd\x5d\x7e\x95\xf2\x45\x6a\x56\x00\x03\x1d\x3b\x16\x36\x29\x65\x50\xf5\xfe\x59\x9d\x2b\x38\x7f\xf6\xa2\x62\xf9\xfe\x22\x1e\xa2\x6d\x3d\x6b\x50\xf8\xc6\xe1\xd4\x89\xc4\x86\xd3\x5f\x6a\xb3\x99\xe2\x81\x7e\xe7\x17\x5f\xf2\x6d\x52\xb3\xa3\x47\xe6\x38\xf5\x06\xf3\x2e\xf7\x3b\x84\xbe\x48\xe3\xbe\xa3\xe8\x45\x89\x64\x3f\xa1\x54\x27\xd0\xee\xbd\x98\xd9\xc2\x4f\x6b\xf7\x96\xa9\x22\x47\xe1\x57\xb8\x4d\x50\x5a\x56\x8e\x89\x1a\x5f\xc2\x34\xc6\xc2\xf8\xd1\x46\xfe\x3b\x39\x4e\x48\x72\x42\x3e\xcd\x65\x92\xb0\x59\x61\x30\xdb\x95\xc9\xf2\x4c\x4d\xbc\x62\xf6\x53\xd0\xf4\xa7\x9b\xa1\xc3\x8f\x68\x2d\xe7\x61\x49\x31\xbc\x22\x73\xc9\xa0\xac\x8d\xc0\x4b\x11\xf9\x33\x3d\xb2\x45\xde\xc5\x9e\xb0\x62\x66\x48\xad\x9a\xc7\x9e\x23\x05\xab\xfa\x42\x9f\x95\x97\x0f\xd4\x51\x11\x12\x47\x93\xb8\x21\x08\xa9\x9c\x2c\x33\xaa\xa2\x47\x6a\x21\x18\x31\xf7\x44\x20\x41\xd5\xc6\x10\x73\x50\x8e\x51\x33\x34\xb4\xb0\x92\x78\xc5\x94\xbd\xcd\xfb\x2b\xa1\x0c\xd5\x83\x02\x8b\x7c\x8e\xb3\x75\xd2\x59\x4d\x03\xa5\xbe\x85\x59\xea\x3a\x53\x72\xe5\x65\xf3\x53\x0c\x89\x42\xac\xbd\xd8\x28\x81\xf4\xc7\x91\x12\x14\x7b\x05\xc2\x85\x03\x6e\x7c\x30\xc8\x48\xd0\xa5\x2e\x32\xe8\x19\x17\x79\xf6\xe1\x23\xd4\xfd\x5f\x79\x52\xa5\xc5\x9a\x3c\xa1\x7e\xfd\x57\x73\x2a\x25\x5b\x59\x24\x85\x27\x8a\x8b\xde\xa4\xac\x65\x9c\xb6\x26\xa9\xa2\xa8\xcb\xab\xcd\x87\x62\x2f\x53\x55\x59\xcb\x88\x51\xf9\xb3\xbc\xd2\x15\x2e\xcd\x6f\x4f\x46\x1c\x7c\xe0\xf7\xe4\xca\xbd\x42\x23\x5c\x43\xd5\x79\xaa\x6e\x9f\x12\xbc\xad\x78\x2e\xe3\x4e\xee\xfd\x28\x67\x5f\xdb\xc7\x3d\x19\x0c\x97\x5d\x5d\x50\xa4\x04\xf5\x91\xa3\xbc\xd9\xf3\xc6\x4c\x5f\xd8\xf9\x3b\x87\x2f\x47\xd5\x77\xbe\x6a\xe3\x7d\x56\xc1\x9f\x27\x14\xd6\xe5\xc5\xeb\x7d\x23\xbd\x56\xd3\x6d\x6b\xc3\xc4\x25\xbe\x90\xa5\x35\xd3\x7d\xeb\x76\x88\xb2\xb7\xcd\x9a\x76\xac\x78\xe1\x81\xd2\x1a\x93\x33\xd8\x80\x2c\xf5\x01\x25\x1a\xea\x28\x99\xa5\x7e\xcb\x8b\x84\x71\x9e\xd6\xb1\xbc\x1a\xd6\x76\xe7\xa3\x79\x0e\x10\x1f\x9b\xcf\x65\xb2\xcd\xa8\x75\xa2\xd5\x50\x87\xb1\x80\x33\x19\xdb\x28\x5b\x56\x05\x74\xeb\x6a\xd5\xb0\x3f\xee\x61\x25\x8b\x95\xa6\x3c\xd5\xef\x75\x89\xf4\xc4\x6a\xd9\xc3\x0f\x4a\xe0\x3e\x9a\x25\x3f\xf0\x3b\x4c\xec\xf7\xe5\xdf\xbe\x56\xd4\xdb\x94\x46\x55\x77\xe8\xaa\x7e\xa6\x97\xde\xde\x60\xe6\x79\xa2\x5c\xe5\x82\x5a\xaf\xd4\xd2\xea\xd2\x92\x5b\xce\x47\xc6\x2c\x6a\x7c\xa7\x2e\x6f\x4d\x4f\xec\xad\x50\x8a\x84\xfc\xe9\x50\x2f\xc4\xa0\x54\xb0\x9b\x85\xb5\x5d\xba\x89\x45\xd6\xc8\x80\x15\xd5\x6e\x05\xf2\x89\x09\xba\xad\xa9\xa9\xae\x5c\xb9\xaf\x41\x1b\x79\x56\x4f\x15\x54\x5c\x46\x50\xcf\xa0\xd8\x9d\x36\x45\xb7\x3a\x4d\x61\xa1\xfa\xc7\xde\xa1\x4f\xc3\x2b\x27\xf9\x4f\x4c\xc3\x28\xe1\x1f\x6e\xdd\xe2\x3b\xda\x6e\x73\xd7\x05\x1f\xfc\xd2\x6d\x70\x29\x7d\x6c\xa5\x0e\xad\xfa\x3c\x71\x46\x21\xfd\x55\x0d\xa4\xc0\x02\x4d\x31\x84\x5e\xd1\x47\x72\xcb\x86\xc6\xbf\x26\x0e\x4a\xa5\xcf\xad\xe2\xee\x2a\xf6\xef\xd1\x46\xdf\x82\x8b\x52\x55\x9f\xd2\x2f\x04\x58\x41\x93\x4d\xb1\x2a\x9f\x35\x4b\x43\x79\xa6\x5d\x0e\xef\x6d\x2a\xb2\x8a\x38\x19\xa4\x0d\x37\x77\xa4\xf8\x0b\xd6\xc7\x85\x0d\x0c\x68\x2e\x3e\xcd\x07\xe7\x7c\x63\x82\x10\x23\xb9\x81\xf1\xd7\x98\x61\x43\x34\xba\x20\xf1\x57\x97\xb7\x5c\x92\x4e\xe4\x81\x31\x73\xd4\xd5\xd3\xda\xe6\xed\xbd\x51\x22\x93\x75\x2e\xd8\x71\xf8\xc3\x6b\xd2\xd6\x02\x76\xa1\xa3\x0b\xb4\x8e\xe3\x2e\xce\x31\x31\x16\x92\x1d\x64\x86\x6d\x7c\x81\xa4\x9b\x74\xfc\xd8\x21\x7a\xc5\x72\xae\x44\x8e\x47\x92\x04\xe5\xa2\x47\x6d\x5e\x33\xfc\x7c\xe0\x4c\xfe\xbe\x79\x5e\xf4\x8e\x64\x6c\x4b\x93\xa9\x2b\x68\x9a\xb3\xe8\x4d\x88\xa2\x50\x46\xa6\x7c\x40\x70\xe3\x0a\x0d\xa1\xdb\x57\xc4\x7f\x64\xc4\x59\xde\xee\xa0\x6f\x2a\xbd\x56\x7a\xcb\x63\x42\xcd\xe4\xc8\x5c\x7a\x1a\x72\xa4\xdc\x20\xe6\xe3\x86\x4f\x39\x4f\x9e\x46\x3f\x9b\x30\xa5\x28\xf6\xb3\x9f\x8c\x7c\x10\xa1\x92\xf8\x10\x32\x70\xc9\xb1\x45\x3b\xd2\x7c\xb4\x60\x6e\x49\x59\x3b\x0c\xdf\x6b\xd3\x37\x55\x40\xe3\xb5\x78\x52\x23\xda\xe8\x23\x4f\x2f\x77\x12\x6d\x88\x53\xcb\x98\xfe\xb0\xa9\x0a\x6a\x5a\x08\x18\xed\x05\x2b\xf7\x2c\xe1\xfc\x1f\x3f\x5f\x03\x1d\x7e\x17\xab\xe9\x5c\x8d\x13\xf5\xf7\x5f\x99\x57\xd3\xb7\xf6\x57\x27\x55\x69\xda\x67\xa4\x96\x38\x74\xe2\xea\xf9\x5d\xd6\xbd\xf0\xbb\xea\xd3\x33\x67\x3e\xdb\x73\x6e\x0c\x93\x17\x0a\x93\xab\x77\xb1\x5e\x7d\xa2\xf0\xe9\x9c\x75\x21\x94\x7f\x32\x3d\xb8\x10\x5f\x8e\x27\x96\x9f\xeb\x6c\x19\x9d\xd1\xae\x09\x60\x21\xcf\x13\x05\xb5\x56\x6a\x0c\xef\xab\x1f\xfa\xfa\x95\xa0\xd1\xd8\xc7\xa0\xdb\xe6\x73\x2c\xf9\x6f\x5f\xbf\x70\xb1\x7b\x72\xd8\x7f\x35\x69\x5a\xc6\x45\x8c\x41\xb9\xcd\xf2\x27\x09\x37\x11\x48\x72\xe2\xec\x00\xa0\x29\xf8\x57\x24\xac\xf1\xaf\x92\xf0\x1f\x97\x02\xde\xde\x09\xf5\x07\x17\x87\x27\x18\xd1\xee\x2b\xf2\x34\x33\x02\xa3\x2f\xde\x82\xe6\x58\xcc\x90\x30\xbd\xcf\x6a\x19\xad\x6f\xf9\x58\xdf\x77\xb1\x49\x16\xdc\xae\x70\xcf\x0e\xe2\xf2\x12\xb4\xc1\x3c\xdf\x9b\x38\x3e\xb5\xf1\xd6\x47\x63\xe1\xd1\x19\xfe\x32\xb0\x13\x1f\x4f\x43\x76\xec\x17\xce\x46\xa2\xbb\x73\x90\xd6\xb2\xb4\xbe\x5a\x8d\x7e\x8c\x2f\x7f\xb3\xf6\x2d\x7d\x41\xcf\xcc\xd0\x88\x53\xde\x39\x1e\x61\x7e\x3d\x31\x70\x45\x09\xfa\xbd\xbb\x3a\xc6\x48\xc1\x3a\xee\xf0\x60\xcf\x4f\x37\x06\x5a\xf6\x3c\x4b\xa0\x5e\xf5\x14\x95\xa9\x4c\xbb\x59\xca\xe3\x29\x12\x82\xe2\x81\x84\x3e\x48\xcb\x6a\xba\xf0\xea\xc0\x8c\xb6\xcd\x7d\x15\xe1\x81\xde\x70\x86\x32\xd2\x5c\xd0\xcd\x80\xc4\x42\x40\x69\xef\x74\x9c\x1d\xfd\xdc\x9a\xba\x9e\x36\xff\x30\x59\x6e\x41\xa8\x7c\x8d\x82\xcc\xc2\xaf\x4c\xd5\x01\xbe\xf4\xb8\x56\xe9\x93\x7a\xae\x45\x27\x6e\x9e\x65\x9a\x46\x26\xef\x50\xd1\xb9\xfd\xd0\x5f\x7c\x62\xfb\x9a\xbd\x6b\x1d\x5f\x8d\x99\xda\x1a\xad\x56\x26\xeb\xe0\xc7\x4f\x45\xee\x33\x1f\xb7\x73\x55\x9d\xc0\xe7\xeb\x9a\x8f\x84\xa7\x85\xf4\xae\x20\x7b\xe5\x27\xcb\xf1\xf1\x48\x19\x8b\x8e\x6a\xdc\x0e\x83\x48\xf6\xeb\x64\x2c\x54\x68\x87\xc6\x1e\xd2\x1f\xe8\x2a\x4f\x35\x82\x95\x1d\xcd\xae\x5e\x41\x3f\x78\x57\x77\x5b\xb7\x38\xee\x78\xf0\x6b\x7c\x71\x4a\xe1\xba\x6a\x38\xba\xc5\xec\x41\x1c\x6a\x34\xaf\x1e\xc2\x04\x67\x9e\x26\x95\x4e\x11\x06\xbc\x98\xbc\x7f\x8e\xad\x7a\x4d\xc7\x04\xce\x02\x00\xdc\xac\xff\xd1\xbf\xd3\xdb\x0b\x83\xfd\x63\x5e\xd8\xb8\x86\xe8\x70\x45\x3e\x5d\xc6\xb6\x10\x7b\x80\x49\x02\xdd\x6e\x00\x6a\xad\x5e\x41\x62\xdb\xf3\x5e\x4b\xda\x3a\x1c\xb5\x13\x14\xa5\xa4\xcc\x83\xb0\xab\xf1\x5f\xf9\x65\x85\xa8\xb0\x7d\x8e\xfb\x39\x26\x2c\xef\xb5\xb3\x55\xbd\x2b\xf8\x7d\xa9\x1c\x47\x72\x78\xa6\x68\xa0\xe9\x84\x83\xa3\x55\x87\x74\x89\x53\xda\x11\x75\x3d\xcb\x1c\xaa\xc5\xcd\xf5\x06\x3a\x18\x80\x48\x75\x46\x58\xb1\x29\x97\x64\x19\xdc\xcd\x1b\x0a\x35\x5f\xdc\xd5\x55\xa8\xa7\x70\x0a\x1f\x1c\x1e\x3d\xf9\xeb\x81\x2b\x67\xdf\xa0\xc2\x7d\xea\x5c\xb9\xe2\x82\xb5\x1c\xb7\x29\x75\xa4\x11\xdd\x29\x93\xcc\xa8\x9a\xc5\xc1\xde\x6e\xf5\x3f\xab\x8f\x94\x52\xcd\xca\x00\x00\x60\xed\x77\x4d\xcc\xb2\x47\x10\xf8\xdf\xfa\xbf\xd7\xcb\x3c\xff\xa7\x23\xdf\x8b\xf6\x9f\x2d\xbf\xd7\x23\x92\x3f\x58\x25\xfe\xb3\xea\xfa\x67\x97\xdf\x6f\xd7\x99\x1f\x5c\xba\xb2\xfe\x6d\x59\xf3\x73\x90\xef\xb1\xa0\xf1\x43\x90\x05\xce\x7f\x77\x6d\x7f\x8e\xf5\x7d\xe7\x7f\x2c\xe8\x31\xd7\xdf\xc6\x9a\x31\x9c\x1d\xf4\x9b\x0b\x76\x80\x1d\x68\x62\x01\x80\x7c\xee\xdf\x5e\xff\x15\x00\x00\xff\xff\xa4\x90\x38\xa8\x7e\x0d\x00\x00")

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

	info := bindataFileInfo{name: "syntax/stdlib-safe.arraiz", size: 3454, mode: os.FileMode(0644), modTime: time.Unix(1, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xf9, 0xfd, 0xd3, 0x22, 0x30, 0xdc, 0x9c, 0x20, 0xb9, 0x79, 0x25, 0x11, 0xe6, 0x87, 0x24, 0xc4, 0x91, 0x29, 0x12, 0xc7, 0x31, 0x49, 0x99, 0x2a, 0x71, 0x86, 0x57, 0xd5, 0x6b, 0xbc, 0x30, 0x7e}}
	return a, nil
}

var _syntaxStdlibUnsafeArraiz = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x0a\xf0\x66\x66\x11\x61\xe0\x60\xe0\x60\x40\x07\x3c\x0c\x0c\x0c\xc9\xf9\x79\x69\x99\xe9\x7a\x89\x45\x45\x89\x99\x97\x4e\xf8\x9e\x39\xd3\xad\x7f\xea\xbc\xfe\xc6\xa0\x00\xaf\xf3\x3a\x27\x35\x7c\x2f\x79\x9f\x3b\x7f\xd5\x5f\x47\xeb\x92\xef\x19\x96\xc0\xce\x99\x5e\x56\x5c\x09\xa1\x99\x33\x2d\x57\x86\x71\x79\xfd\xda\xf9\x53\x73\xa9\xe5\xca\x5f\x16\x21\xaf\xc5\xc3\x9f\x4d\x55\x09\x7b\x21\x3e\x7d\xc9\x93\x14\x35\x47\xcd\xbb\x97\xf6\x3c\x3a\xb2\x25\xc4\xb6\xa3\x72\x56\x03\x03\x03\xc3\xff\xff\x01\xde\xec\x1c\xfa\x49\xa6\xa6\x91\x0c\x0c\x0c\x65\x0c\x0c\x0c\xb8\x9d\xa4\xca\xc0\xc0\x90\x9b\x9f\x52\x9a\x93\xaa\x9f\x9e\x59\x92\x51\x9a\xa4\x97\x9c\x9f\xab\x9f\x58\x54\xa4\x9b\x98\xa9\x0f\x76\xa4\x7e\x7a\xbe\x5e\x6e\x7e\x4a\xcd\x14\xef\x5d\x8f\xd5\x04\xdc\x8e\x7f\xd4\x32\xf1\x61\x7b\x14\xd8\x20\xd7\x9c\x32\x4d\x9f\xfb\x85\xa8\x42\xc8\xa4\xcc\xd3\x4e\x9a\x13\x7e\xdf\x5a\x7d\x5e\x7d\xf5\xee\xf5\x17\xdf\xde\x11\x6c\xfe\x65\x29\xd3\x75\xdc\xf4\xcb\xeb\x96\xa3\x2f\x2e\x57\xbe\x9d\xc1\x1d\x9e\xd7\xe9\xf5\x7a\xfb\xc1\x5f\x97\x33\x79\x5a\x53\x82\xe6\xfa\xb5\x9c\x4e\x89\x9d\xd6\x73\x5a\x41\x74\x3b\x73\x90\x46\xde\x24\xcd\xb4\x2e\xfb\xff\xa7\x02\x3e\x2f\x7e\xf3\xa6\xa6\xf2\x3c\xc3\x65\xe9\x57\xf7\xaf\x0a\x8b\xed\x9f\x5f\x57\xcd\xb9\xda\xeb\x8f\xed\x72\xfd\x17\xff\x16\x6e\xd0\x7b\x3a\xfd\x74\xed\xa6\x23\x1f\xcf\x3e\x7c\x22\x7f\x22\x3f\xe0\x4e\x83\xa6\xda\x26\xcf\x32\xa7\xe9\xe7\x1e\x66\x29\xab\xc7\xac\xd8\xc1\xa9\x36\x43\xc3\xb1\x41\x80\xab\xe7\x82\xd4\x02\xbd\xad\x39\x5b\x58\x0c\x43\xb6\x96\xb1\x3e\x6f\x9d\x5f\xed\xff\x22\xa8\xef\x97\xe2\xee\x98\x63\x73\x99\x27\xb0\x7f\x3b\xc6\x23\x22\x67\x91\xcf\x79\x53\x2e\xf7\xdb\x91\x9b\xc7\x3e\x87\xe5\xd7\xc5\xc6\x69\x2c\xf9\x90\xa7\x25\xbd\x90\xd1\x65\x6d\x76\xb2\x97\x44\x1b\x5f\xf5\x1b\x81\xb3\xfa\xcf\xeb\x5e\x3f\x4b\x9c\xbd\xea\x4b\xe3\xdb\x73\xfd\xd6\x42\xea\x5b\xfe\x3e\x49\x3f\xf3\x61\x6a\x40\x90\xe3\xd6\x87\xb7\xa5\xac\x58\xee\xbe\x66\xdb\xfc\x32\x59\x63\xa2\x46\xc0\x93\xdf\x1b\xe6\x4e\x4d\xf9\x34\x41\xd1\xde\x77\x5a\xfd\x61\x3d\xb6\xf2\x83\x5b\x5a\x4d\x4e\xcf\x7e\xf2\xf6\xcd\xa1\x96\x70\xeb\xcf\x33\x4e\xdc\xfc\x11\x32\xd1\x27\x75\xe5\x6a\xdf\xcc\x7b\x62\x5f\x2f\x2e\x2b\xb9\x9f\x9b\xea\x17\xa5\xf7\xd8\x5f\xf2\xe0\xad\x66\x9e\x5a\xd7\xef\x1b\x26\xef\xff\x60\xa9\x62\x7f\xe7\xbb\x8f\xdf\x02\x53\x37\x2e\xc3\x4f\x79\x9b\xde\x19\x4d\x78\xff\xf6\xef\x6f\xdb\x49\x3e\x7f\xde\xbb\x76\xcc\x38\x67\xee\x9d\x7c\x83\x85\x53\xf8\xee\x77\xdd\x27\xbc\xfd\xf3\xae\xa9\x3e\x3a\xb3\xdc\xe0\xd2\xaf\x9c\x88\xe2\x82\x3f\x3f\xec\x3e\x15\x29\xf9\xbc\xf3\xb1\x29\x89\xbf\x3f\xaf\xfd\xd7\xcc\x2b\xfd\x05\x87\x76\xbc\x0f\xdf\xf4\x6b\xc1\x96\x4b\x09\xf9\x27\xee\x3e\xd9\xff\x5e\x68\x5e\xa3\xc1\xcd\xca\x94\x46\xb7\x7f\xc7\x64\xef\x4e\xbf\x7e\x9a\xdd\xc9\x44\xc1\xe1\x54\xd4\x8c\x4b\x3b\xa4\x72\x2e\xaf\xd8\x36\xed\x9b\xc2\xa7\x53\x33\xb6\x4f\xfe\xcf\x39\xe1\x50\x83\xe0\x54\xa6\xe0\xa9\xdb\x8a\x6f\x69\x4c\xd0\x2d\x7b\x26\xec\xb0\x21\xdb\xca\x59\xed\x22\x17\xc3\x4a\x13\x93\xa5\x2f\x0f\x7e\x09\xbb\xb5\xec\xfa\x2b\x61\xb7\x3d\xc7\x7d\x16\x05\x39\x71\x30\x7b\x4d\x49\xca\x60\xe7\x9e\x2a\x3b\xeb\xd7\x8f\x1b\xf1\x8c\x9c\xf7\x2f\xdd\xfb\x92\xb5\xba\xa5\xd0\xe9\x9b\x90\x44\xa1\xdd\x9e\xcb\xb9\x13\xd7\x46\x1e\xeb\xd8\xf5\x66\x95\x1c\xdb\x65\xde\x45\xeb\x7f\x3f\x56\x92\x0f\xd8\xb2\xe3\xf9\xcf\xf7\xff\x0f\xb6\xff\x2a\x9e\x9b\x27\x3b\xef\x7f\xf2\xd6\xf4\x69\xa7\x4f\xfc\x39\xf3\xfe\x9b\xd5\xc6\xeb\xf3\xf5\xf5\x57\x3c\x5a\x76\xef\xed\x92\x15\xa9\xdc\xba\x0e\x2b\x6f\x9f\xca\x48\xd6\x56\x7d\xa5\xb7\xf9\xf7\xe3\x2f\xf1\xf0\x94\xed\x7a\xa8\x2b\xae\x9d\x89\x81\xe1\x03\x2b\xbe\x94\xed\x40\x38\x65\x17\x57\xe6\x95\x24\x56\xe8\x17\x97\xa4\xe4\x64\x26\x41\x29\xdd\xd2\xbc\xe2\xc4\xb4\x54\x48\x06\x4d\xe9\x73\xec\x6a\x09\xe0\x71\xf9\x3e\x6f\xd1\x81\xbd\xdf\x67\xfc\xd5\x49\x17\xfe\x7b\x51\xfe\x59\xb9\x92\x67\xce\xc7\x6f\xbd\x0a\xbd\xbe\xef\xb9\x56\xc9\x7d\xad\xfa\x13\xc6\x79\xe4\xc9\x26\xb5\xb0\xed\xcb\x45\xf8\x04\xca\xac\x0c\x66\x2b\xbf\x39\x38\xc3\x7c\x6e\x90\x0f\xf3\xa2\x13\xa7\x3a\x54\x74\x9f\xb0\xfa\xba\xa9\xdb\x4f\x68\x9b\xb8\x71\x77\xca\xa6\x79\xda\x45\xf2\x1e\x6c\x7f\x54\x26\xfd\xb3\xb4\xdb\xb0\xae\xbe\xd9\x3c\x65\xa3\xdb\xae\x86\x73\xbf\xef\x5c\x3b\xa6\xa0\x5b\x70\xee\xac\xa9\x43\x56\xe4\xcd\xf5\x69\xd3\xd7\x2c\x08\x7a\x72\x65\x51\x99\xf1\x9f\xb2\x17\xff\x94\xbc\x27\x6f\x64\x6e\xd0\xde\xc7\x04\x0b\x82\xc9\x73\x57\x6c\x9c\xca\xc0\xc0\xf0\x03\x9c\xb9\x19\x99\x44\x18\x10\x81\x80\x9c\xf1\x79\x30\x82\x05\xb9\xfc\x41\xd7\x89\x1c\xb0\xaa\x28\xba\x26\x13\x5b\x4c\xa0\x1b\x89\xec\x50\x07\x14\x23\x73\x99\x29\x8f\x9f\x00\x6f\x56\x36\x90\x59\xcc\x0c\xcc\x0c\xbf\x19\x18\x18\x0a\x58\x40\x3c\x40\x00\x00\x00\xff\xff\x80\x12\x2f\x6d\x81\x05\x00\x00")

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

	info := bindataFileInfo{name: "syntax/stdlib-unsafe.arraiz", size: 1409, mode: os.FileMode(0644), modTime: time.Unix(1, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x56, 0xf7, 0x1b, 0xa9, 0x91, 0x16, 0x3b, 0x12, 0xc4, 0x67, 0xad, 0xc4, 0x43, 0x55, 0x4b, 0x25, 0xb6, 0xd6, 0x99, 0xf1, 0x3, 0xb9, 0xa8, 0x53, 0x7, 0xcc, 0x16, 0x5b, 0x7b, 0x37, 0xd2, 0xad}}
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

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
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

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
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

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
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

// AssetDebug is true if the assets were built with the debug flag enabled.
const AssetDebug = false

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
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
	"syntax": {nil, map[string]*bintree{
		"implicit_import.arrai": {syntaxImplicit_importArrai, map[string]*bintree{}},
		"stdlib-safe.arraiz":    {syntaxStdlibSafeArraiz, map[string]*bintree{}},
		"stdlib-unsafe.arraiz":  {syntaxStdlibUnsafeArraiz, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory.
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
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
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

package syntax

import (
	"embed"
	"io/ioutil"
	"sync"
)

var (
	//go:embed embed/*
	bindata embed.FS

	embeddedFileCache = map[string][]byte{}
	embeddedFileMutex sync.Mutex
)

func implicitImportArrai() []byte {
	return mustReadEmbeddedFile("embed/implicit_import.arrai")
}

func stdlibSafeArraiz() []byte {
	return mustReadEmbeddedFile("embed/stdlib-safe.arraiz")
}

func stdlibUnsafeArraiz() []byte {
	return mustReadEmbeddedFile("embed/stdlib-unsafe.arraiz")
}

func mustReadEmbeddedFile(path string) []byte {
	embeddedFileMutex.Lock()
	defer embeddedFileMutex.Unlock()

	data, has := embeddedFileCache[path]

	if !has {
		f, err := bindata.Open(path)
		if err != nil {
			panic(err)
		}

		data, err = ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}

		embeddedFileCache[path] = data
	}

	return data
}

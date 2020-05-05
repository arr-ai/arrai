package main

import (
	"io/ioutil"
	"os"
)

// Reads ../syntax/arrai.wbnf file
// and encodes them as strings literals in ../syntax/parser.go
func main() {
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	content := append(
		[]byte("// AUTOGENERATED. DO NOT EDIT.\npackage syntax\n\nfunc unfakeBackquote(s string) string {\n  return strings.ReplaceAll(s, \"‵\", \"`\")\n}\n\nvar arraiParsers = wbnf.MustCompile(unfakeBackquote(`\n"), //nolint:lll
		data...)
	content = append(content, []byte("\n`), nil)")...)
	err = ioutil.WriteFile(os.Args[2], content, 0600)
	if err != nil {
		panic(err)
	}
}

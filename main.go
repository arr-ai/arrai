package main

//go:generate go run tools/parser/generate_parser.go syntax/arrai.wbnf syntax/parser.go
//go:generate goimports -w syntax/parser.go

func main() {
}

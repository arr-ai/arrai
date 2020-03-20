package syntax

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/go-errors/errors"

	"github.com/arr-ai/arrai/rel"
)

// Token represents a lexical token.
type Token rune

// Non-character tokens
const (
	NULL  Token = 0
	ERROR Token = 255 + iota
	EOF

	NUMBER
	IDENT
	STRING
	XML

	AND
	AS
	ELSE
	EXCEPT
	FOR
	IF
	IN
	MAX
	MEAN
	MEDIAN
	MIN
	NEST
	OR
	ORDER
	UNNEST
	WHERE
	WITH
	WITHOUT
	COUNT
	SUM

	ARROW   // ->
	ARROWST // ->*
	ATARROW // @>
	CSET    // |}
	DARROW  // =>
	GEQ     // >=
	DSLASH  // //
	JOIN    // <&>
	LEQ     // <=
	NEQ     // !=
	OSET    // {|
	PI      // π
	SQRT    // √
	SUBMOD  // -%
)

// const xmlTokenBase = Token(1000)

var keywordTokens = map[string]Token{
	"and":     AND,
	"as":      AS,
	"count":   COUNT,
	"else":    ELSE,
	"except":  EXCEPT,
	"for":     FOR,
	"if":      IF,
	"in":      IN,
	"max":     MAX,
	"mean":    MEAN,
	"median":  MEDIAN,
	"min":     MIN,
	"nest":    NEST,
	"or":      OR,
	"order":   ORDER,
	"sum":     SUM,
	"unnest":  UNNEST,
	"where":   WHERE,
	"with":    WITH,
	"without": WITHOUT,
}

var wideTokens = map[string]Token{
	"->":  ARROW,
	"->*": ARROWST,
	"@>":  ATARROW,
	"|}":  CSET,
	"=>":  DARROW,
	">=":  GEQ,
	"//":  DSLASH,
	"<&>": JOIN,
	"<=":  LEQ,
	"!=":  NEQ,
	"{|":  OSET,
	"π":   PI,
	"√":   SQRT,
	"-%":  SUBMOD,
}

var tokenReprs = func() map[Token]string {
	reprs := map[Token]string{
		NULL:   "NULL",
		ERROR:  "ERROR",
		NUMBER: "NUMBER",
		IDENT:  "IDENT",
		STRING: "STRING",
		XML:    "XML",
	}
	for repr, token := range keywordTokens {
		reprs[token] = repr
	}
	for repr, token := range wideTokens {
		reprs[token] = repr
	}
	return reprs
}()

// TokenRepr returns a string representation of a token value.
func TokenRepr(token Token) string {
	if token != 0 && token < 256 {
		return fmt.Sprintf("%q", rune(token))
	}
	if repr, found := tokenReprs[token]; found {
		return repr
	}
	return fmt.Sprintf("Token(%d)", token)
}

// ParseArraiString parses an arr.ai string.
func ParseArraiString(lexeme []byte) string {
	var s string
	err := json.Unmarshal(lexeme, &s)
	if err != nil {
		panic(err)
	}
	return s
}

type lexerState func(l *Lexer) (Token, interface{})

// Lexer extracts a stream of tokens from an input file.
type Lexer struct {
	fr     FileRange
	stack  []lexerState
	reader io.Reader
	data   interface{}
	err    error
	buffer *bytes.Buffer
	offset int
	prev   int
	state  lexerState
	lexeme int
	token  Token
	eof    bool
	fresh  bool
}

// NewStringLexer returns a new Lexer for the given input.
func NewStringLexer(input string) *Lexer {
	return NewLexer(bytes.NewBufferString(input))
}

// NewLexer returns a new Lexer for the given input.
func NewLexer(reader io.Reader) *Lexer {
	return NewLexerWithPrefix(bytes.NewBuffer([]byte{}), reader)
}

// NewLexerWithPrefix returns a new Lexer for the given input.
func NewLexerWithPrefix(prefix *bytes.Buffer, reader io.Reader) *Lexer {
	return &Lexer{
		reader: reader,
		buffer: prefix,
		state:  LexerInitState,
		fr:     FileRange{FilePos{1, 1}, FilePos{1, 1}},
	}
}

// func (l *Lexer) copy() *Lexer {
// 	// Copy most fields
// 	newL := *l

// 	if l.stack != nil {
// 		newL.stack = append([]lexerState{}, l.stack...)
// 	}
// 	newL.buffer = bytes.NewBuffer(l.buffer.Bytes())

// 	// We need to duplicate the reader since reading from it is destructive
// 	remBuf, err := ioutil.ReadAll(l.reader)
// 	if err != nil {
// 		panic(err)
// 	}
// 	l.reader = bytes.NewBuffer(remBuf)
// 	newL.reader = bytes.NewBuffer(remBuf)
// 	return &newL
// }

// Reader returns the most recently recognized token.
func (l *Lexer) Reader() io.Reader {
	return l.reader
}

// Offset returns the current scanning position as an offset from the start of
// the input.
func (l *Lexer) Offset() int {
	return l.offset
}

// Tail returns the unconsumed portion of the buffer.
func (l *Lexer) Tail() []byte {
	return l.buffer.Bytes()[l.offset:]
}

// Token returns the most recently recognized token.
func (l *Lexer) Token() Token {
	return l.token
}

// Lexeme returns the lexeme for the most recently recognized token.
func (l *Lexer) Lexeme() []byte {
	return l.buffer.Bytes()[l.lexeme:l.offset]
}

// Value returns the Value for the most recently recognized token.
func (l *Lexer) Value() rel.Value {
	if l.data == nil {
		return nil
	}
	return l.data.(rel.Value)
}

// Data returns the current data.
func (l *Lexer) Data() interface{} {
	return l.data
}

// Value returns the Value for the most recently recognized token.
func (l *Lexer) Error() error {
	return l.err
}

// FileRange returns the FileRange for the most recently recognized token.
func (l *Lexer) FileRange() FileRange {
	return l.fr
}

// PushState pushes the current state onto the stack, and makes the given state
// current.
func (l *Lexer) PushState(state lexerState) {
	l.stack = append(l.stack, l.state)
	l.state = state
	l.unpeek()
}

// PopState pops the current state off the stack and makes it current.
func (l *Lexer) PopState() {
	if len(l.stack) == 0 {
		panic("State stack empty")
	}
	l.state = l.stack[len(l.stack)-1]
	l.stack = l.stack[:len(l.stack)-1]
	l.unpeek()
}

// InState wraps a lambda in PushState()/PopState().
func (l *Lexer) InState(state lexerState, f func()) {
	l.PushState(state)
	defer l.PopState()
	f()
}

func (l *Lexer) unpeek() {
	if l.fresh {
		l.offset = l.prev
		l.lexeme = l.offset
		l.fresh = false
	}
}

var emptyBlock = make([]byte, 4096)

func (l *Lexer) eatRE(re *regexp.Regexp) [][]byte {
	var m [][]byte
	for {
		tail := l.Tail()
		m = re.FindSubmatch(tail)
		if m != nil && len(m[0]) < len(tail) {
			break
		}
		if l.eof {
			break
		}
		size := l.buffer.Len()
		l.buffer.Write(emptyBlock)
		n, err := l.reader.Read(l.buffer.Bytes()[size:])
		if n < len(emptyBlock) {
			l.buffer.Truncate(size + n)
		}
		if err != nil {
			if err == io.EOF {
				l.eof = true
				break
			}
			panic(err)
		}
		if n == 0 {
			break
		}
	}
	if len(m) > 0 {
		l.lexeme = l.offset + len(m[1]) // Skip over whitespace.
		l.prev = l.offset
		l.offset += len(m[0])               // Skip over lexeme.
		l.fr.Start = l.fr.End.Advance(m[1]) // Advance Start past whitespace.
		l.fr.End = l.fr.Start.Advance(m[2]) // Advance End past lexeme.
	}
	return m
}

// Fail sets an error and returns the ERROR token.
func (l *Lexer) Fail(err error) Token {
	l.err = err
	return ERROR
}

// Failf produces a formatted error with a line marker.
func (l *Lexer) Failf(fmtStr string, args ...interface{}) Token {
	return l.Fail(
		errors.Errorf(fmt.Sprintf(fmtStr, args...) + "\n" + l.String()))
}

// String produces a formatted string representation of the lexer with a line
// marker.
func (l *Lexer) String() string {
	input := l.buffer.Bytes()
	sol := bytes.LastIndexByte(input[:l.offset], '\n') + 1
	eol := bytes.IndexByte(l.Tail(), '\n')
	if eol == -1 {
		eol = len(input)
	} else {
		eol += l.offset
	}

	pre := input[sol:l.lexeme]
	lexeme := input[l.lexeme:l.offset]
	post := input[l.offset:eol]

	return fmt.Sprintf("%s🔥 %s🔥 %s", pre, lexeme, post)
}

func found(fmt string, args ...interface{}) {
	// if yyDebug >= 5 {
	//     log.Printf(">>>>>> " + fmt, args...)
	// }
}

// Scan scans the expected tokens, otherwise stays put. If no expected tokens
// are given, scans any token. Returns true iff a token was scanned.
func (l *Lexer) Scan(expected ...Token) bool {
	token := l.lex(false)
	if len(expected) > 0 {
		for _, e := range expected {
			if token == e {
				return true
			}
		}
	}
	l.fresh = len(expected) > 0
	return !l.fresh
}

// Peek peeks at the next token. First scans the next token if Peek() has not
// been called since the last call to Lex().
func (l *Lexer) Peek() Token {
	return l.lex(true)
}

func (l *Lexer) lex(fresh bool) Token {
	if l.fresh {
		l.fresh = fresh
		return l.token
	}
	l.token, l.data = l.state(l)
	if l.token == NULL {
		l.fresh = true
		return l.Failf("Syntax error")
	}
	l.fresh = fresh
	return l.token
}

// LexerSymbol is a convenience structure for defining token recognition regexes
// and handler functions.
type LexerSymbol struct {
	token Token
	name  string
	re    *regexp.Regexp

	data func(tok Token, match [][]byte) (interface{}, Token)
}

// ScanOperator tries to recognise an operator or returns NULL.
func (l *Lexer) ScanOperator(operatorsRe *regexp.Regexp) (Token, interface{}) {
	if m := l.eatRE(operatorsRe); m != nil {
		symbol := m[2]
		found("SYMBOL %q", symbol)
		if len(symbol) == 1 {
			return Token(symbol[0]), nil
		}
		if token, found := wideTokens[string(symbol)]; found {
			return token, nil
		}
		panic("Unknown symbol: " + string(symbol) + " " + l.String())
	}
	return NULL, nil
}

// ScanSymbol tries to scan each given symbol or returns NULL.
func (l *Lexer) ScanSymbol(symbols []LexerSymbol) (Token, interface{}) {
	for _, sym := range symbols {
		if m := l.eatRE(sym.re); m != nil {
			found("%s %s", sym.name, m[2])
			var data interface{}
			t := sym.token
			if sym.data != nil {
				data, t = sym.data(t, m)
			} else {
				data = nil
			}
			return t, data
		}
	}
	return NULL, nil
}

// ScanOperatorOrSymbol tries to scan an operator or a symbol, or returns NULL.
func (l *Lexer) ScanOperatorOrSymbol(
	operatorsRe *regexp.Regexp, symbols []LexerSymbol,
) (Token, interface{}) {
	if token, data := l.ScanOperator(operatorsRe); token != NULL {
		return token, data
	}
	if token, data := l.ScanSymbol(symbols); token != NULL {
		return token, data
	}
	return NULL, nil
}

// tokRE enriches a token regex, anchoring it to the start of its input and
// capturing any leading whitespace.
//   match[0]: full match
//   match[1]: leading whitespace
//   match[2]: lexeme
func tokRE(re string) *regexp.Regexp {
	return regexp.MustCompile(`\A(\s*)(` + re + `)`)
}

var lexerOperatorsRe = tokRE(
	`` +
		`!=|` +
		`&>|` +
		`-%|->\*?|` +
		`//|` +
		`<&>?|<=|` +
		`=>|` +
		`>=|` +
		`@>|` +
		`{[|:]|` +
		`[|:]}|` +
		`π|√|` +
		`[-!%&()*+,./:<=>[\\\]^{|}~]`,
)

var lexerSymbols = []LexerSymbol{
	{IDENT, "IDENT", tokRE(rel.LexerNamePat),
		func(tok Token, match [][]byte) (interface{}, Token) {
			if t, found := keywordTokens[string(match[2])]; found {
				return nil, t
			}
			return nil, tok
		},
	},
	{NUMBER, "NUMBER",
		tokRE(`(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+)(?:[Ee][-+]?[0-9]+)?`),
		func(tok Token, match [][]byte) (interface{}, Token) {
			n, err := strconv.ParseFloat(string(match[2]), 64)
			if err != nil {
				panic("Wat?")
			}
			return rel.NewNumber(n), tok
		},
	},
	{STRING, "STRING", tokRE(`"(?:\\.|[^\\"])*"`), nil},
	{EOF, "EOF", tokRE(`\z`), nil},
}

// LexerInitState recognises the next input Token.
func LexerInitState(l *Lexer) (Token, interface{}) {
	return l.ScanOperatorOrSymbol(lexerOperatorsRe, lexerSymbols)
}

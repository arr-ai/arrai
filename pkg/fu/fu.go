// Package fu provides "fmt utilities", helper functions for implementing
// fmt.Stringer and fmt.Formatter.
package fu

import (
	"fmt"
	"io"
)

func Format(i interface{}, f fmt.State, verb rune) {
	switch i := i.(type) {
	case fmt.Formatter:
		i.Format(f, verb)
	default:
		Fprint(f, i)
	}
}

func Fprint(w io.Writer, a ...interface{}) {
	if _, err := fmt.Fprint(w, a...); err != nil {
		panic(err)
	}
}

func Fprintf(w io.Writer, format string, a ...interface{}) {
	if _, err := fmt.Fprintf(w, format, a...); err != nil {
		panic(err)
	}
}

func FRepr(w io.Writer, i interface{}) {
	Fprintf(w, "%v", i)
}

func Repr(i interface{}) string {
	return fmt.Sprintf("%v", i)
}

func String(i interface{}) string {
	return fmt.Sprintf("%s", i)
}

func Write(w io.Writer, b []byte) {
	if _, err := w.Write(b); err != nil {
		panic(err)
	}
}

var stringFormats = map[[3]bool]string{
	{false, false, false}: "%s",
	{false, false, true}:  "%.*s",
	{false, true, false}:  "%*s",
	{false, true, true}:   "%*.*s",
	{true, false, false}:  "%-s",
	{true, false, true}:   "%-.*s",
	{true, true, false}:   "%-*s",
	{true, true, true}:    "%-*.*s",
}

func WriteFormattedValue(f fmt.State, i interface{}) {
	var buf [3]interface{}
	args := buf[:0]

	width, wok := f.Width()
	if wok {
		args = append(args, width)
	}

	prec, pok := f.Precision()
	if pok {
		args = append(args, prec)
	}

	dash := f.Flag('-')

	args = append(args, i)

	Fprintf(f, stringFormats[[3]bool{dash, wok, pok}], args...)
}

func WriteString(w io.Writer, s string) {
	if _, err := io.WriteString(w, s); err != nil {
		panic(err)
	}
}

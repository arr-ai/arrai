package syntax

import (
	"fmt"
	"strconv"
	"strings"
)

func parseArraiStringFragment(s string, validEscapes string, indent string) string {
	if strings.HasPrefix(validEscapes, "`") {
		return strings.ReplaceAll(s, "``", "`")
	}

	var sb strings.Builder

	number := func(i, size, base int) int {
		n, err := strconv.ParseUint(s[i:i+size], base, size*base/4)
		if err != nil {
			panic(err)
		}
		sb.WriteRune(rune(n))
		return i + size
	}

	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '\\':
			i++
			switch s[i] {
			case 'x':
				i = number(i+1, 2, 16)
			case 'u':
				i = number(i+1, 4, 16)
			case 'U':
				i = number(i+1, 8, 16)
			case '0', '1', '2', '3', '4', '5', '6', '7':
				i = number(i, 3, 8)
			case 'a':
				sb.WriteByte('\a')
			case 'b':
				sb.WriteByte('\b')
			case 'e':
				sb.WriteByte('\x1b')
			case 'f':
				sb.WriteByte('\f')
			case 'n':
				sb.WriteByte('\n')
			case 'r':
				sb.WriteByte('\r')
			case 't':
				sb.WriteByte('\t')
			case 'v':
				sb.WriteByte('\v')
			case '\\':
				sb.WriteByte('\\')
			case '\'':
				sb.WriteByte('\'')
			case '"':
				sb.WriteByte('"')
			case 'i':
				sb.WriteString(indent)
			default:
				if strings.ContainsRune(validEscapes, rune(c)) {
					sb.WriteByte(c)
				}
				panic(fmt.Errorf("unrecognized \\-escape: %q", s[i]))
			}
		default:
			sb.WriteByte(c)
		}
	}
	return sb.String()
}

func parseArraiString(s string) string {
	quote, s := s[:1], s[1:len(s)-1]
	return parseArraiStringFragment(s, quote, "")
}

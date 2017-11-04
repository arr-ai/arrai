package syntax

import (
	"bytes"
	"fmt"
)

var newline = []byte{'\n'}

// FilePos represents a position in a line-based file.
type FilePos struct {
	Line   int
	Column int
}

// String returns a string representation of a FilePos.
func (fp FilePos) String() string {
	return fmt.Sprintf("%d:%d", fp.Line, fp.Column)
}

// Less returns true iff fp comes before fp2.
func (fp FilePos) Less(fp2 FilePos) bool {
	if fp.Line != fp2.Line {
		return fp.Line < fp2.Line
	}
	return fp.Column < fp2.Column
}

// Advance returns a FilePos advanced by the given bytes.
func (fp FilePos) Advance(b []byte) FilePos {
	nLines := bytes.Count(b, newline)
	if nLines == 0 {
		return FilePos{fp.Line, fp.Column + len(b)}
	}
	return FilePos{fp.Line + nLines, len(b) - bytes.LastIndexByte(b, '\n')}
}

// FileRange represents a range of characters in a line-based.
type FileRange struct {
	Start FilePos
	End   FilePos
}

// String returns a string representation of a FileRange.
func (fr FileRange) String() string {
	start := fr.Start
	end := fr.End
	if start == end {
		return start.String()
	}
	if !start.Less(end) {
		return "(-range)"
	}

	var endCol string

	end.Column--
	if end.Column == 0 {
		end.Line--
		endCol = "⍵"
	} else {
		endCol = fmt.Sprintf("%d", end.Column)
	}

	if start.Line == end.Line {
		return fmt.Sprintf("%s-%s", start, endCol)
	}
	return fmt.Sprintf("%s-%d:%s", start, end.Line, endCol)
}

// Union returns the minimal range that covers fr and fr2.
func (fr FileRange) Union(fr2 FileRange) FileRange {
	if fr2.Start.Less(fr.Start) {
		fr.Start = fr2.Start
	}
	if fr.End.Less(fr2.End) {
		fr.End = fr2.End
	}
	return fr
}

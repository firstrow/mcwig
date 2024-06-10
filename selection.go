package mcwig

type Selection struct {
	Start Cursor
	End   Cursor
}

func SelectionCursorInRange(sel *Selection, c Cursor) bool {
	s := *sel

	if s.Start.Line > s.End.Line {
		s.Start, s.End = s.End, s.Start
	}

	if s.Start.Line == s.End.Line {
		if s.Start.Char > s.End.Char {
			s.Start, s.End = s.End, s.Start
		}
	}

	if c.Line < s.Start.Line || c.Line > s.End.Line {
		return false
	}

	if c.Line == s.Start.Line {
		if c.Char < s.Start.Char {
			return false
		}
	}

	if c.Line == s.End.Line {
		if c.Char > s.End.Char {
			return false
		}
	}

	return true
}

func SelectionToString(buf *Buffer) string {
	if buf.Selection == nil {
		return ""
	}

	s := *buf.Selection

	if s.Start.Line > s.End.Line {
		s.Start, s.End = s.End, s.Start
	}

	if s.Start.Line == s.End.Line {
		if s.Start.Char > s.End.Char {
			s.Start, s.End = s.End, s.Start
		}
	}

	lineStart := lineByNum(buf, s.Start.Line)
	lineEnd := lineByNum(buf, s.End.Line)

	if lineStart == nil {
		return ""
	}

	endCh := s.End.Char + 1
	if endCh > len(lineEnd.Value) {
		endCh = len(lineEnd.Value)
	}

	if s.Start.Line == s.End.Line {
		return string(lineStart.Value[s.Start.Char:endCh])
	}

	acc := string(lineStart.Value[s.Start.Char:])
	currentLine := lineStart.Next()
	for currentLine != nil {
		if currentLine != lineEnd {
			acc += "\n" + string(currentLine.Value)
		} else {
			acc += "\n" + string(currentLine.Value[:endCh])
			break
		}
		currentLine = currentLine.Next()
	}

	return acc
}

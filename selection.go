package mcwig

type Selection struct {
	Start Cursor
	End   Cursor
}

func SelectionCursorInRange(sel *Selection, c Cursor) bool {
	s := SelectionNormalize(sel)

	if c.Line < s.Start.Line || c.Line > s.End.Line {
		return false
	}

	if c.Line == s.Start.Line && c.Char < s.Start.Char {
		return false
	}

	if c.Line == s.End.Line && c.Char > s.End.Char {
		return false
	}

	return true
}

func SelectionToString(buf *Buffer) string {
	if buf.Selection == nil {
		return ""
	}

	s := SelectionNormalize(buf.Selection)

	lineStart := CursorLineByNum(buf, s.Start.Line)
	lineEnd := CursorLineByNum(buf, s.End.Line)

	if lineStart == nil {
		return ""
	}

	endCh := s.End.Char + 1
	if endCh > len(lineEnd.Value) {
		endCh = len(lineEnd.Value)
	}

	if s.Start.Line == s.End.Line {
		if len(lineStart.Value) == 0 {
			return ""
		}
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

func SelectionNormalize(sel *Selection) Selection {
	if sel == nil {
		return Selection{}
	}

	s := *sel

	if s.Start.Line > s.End.Line {
		s.Start, s.End = s.End, s.Start
	}

	if s.Start.Line == s.End.Line && s.Start.Char > s.End.Char {
		s.Start, s.End = s.End, s.Start
	}

	return s
}

func SelectionStart(buf *Buffer) {
	buf.Selection = &Selection{
		Start: buf.Cursor,
		End:   buf.Cursor,
	}
}

func SelectionStop(buf *Buffer) {
	buf.Selection.End = buf.Cursor
}

func WithSelection(fn func(Context)) func(Context) {
	return func(ctx Context) {
		fn(ctx)
		buf := ctx.Buf
		if buf.Selection == nil {
			// TODO: this is workaround for when selection was deleted but did
			// not exited VIS_LINE_MODE
			CmdNormalMode(ctx)
			return
		}
		buf.Selection.End = buf.Cursor

		if buf.Mode() == MODE_VISUAL_LINE {
			if buf.Selection.Start.Line > buf.Selection.End.Line {
				lineStart := CursorLineByNum(buf, buf.Selection.Start.Line)
				buf.Selection.Start.Char = len(lineStart.Value) - 1
				buf.Selection.End.Char = 0
			} else {
				lineEnd := CursorLineByNum(buf, buf.Selection.End.Line)
				buf.Selection.Start.Char = 0
				buf.Selection.End.Char = len(lineEnd.Value) - 1
			}
		}
	}
}

func SelectionDelete(ctx Context) {
	defer func() {
		ctx.Buf.Selection = nil
	}()
	if ctx.Buf.Selection == nil {
		return
	}

	sel := SelectionNormalize(ctx.Buf.Selection)

	yankSave(ctx)

	lineStart := CursorLineByNum(ctx.Buf, sel.Start.Line)
	lineEnd := CursorLineByNum(ctx.Buf, sel.End.Line)

	if sel.Start.Line == sel.End.Line {
		if len(lineStart.Value) == 0 || ctx.Buf.Mode() == MODE_VISUAL_LINE {
			ctx.Buf.Lines.Remove(lineStart)
			CmdCursorBeginningOfTheLine(ctx)
			return
		}

		if sel.End.Char < len(lineStart.Value) {
			lineStart.Value = append(lineStart.Value[:sel.Start.Char], lineStart.Value[sel.End.Char+1:]...)
		} else {
			lineStart.Value = lineStart.Value[:sel.Start.Char]
		}

		cursorGotoChar(ctx.Buf, sel.Start.Char)
	} else {
		// delete all lines between start and end line
		for lineStart.Next() != lineEnd {
			ctx.Buf.Lines.Remove(lineStart.Next())
		}

		lineStart.Value = lineStart.Value[:sel.Start.Char]

		if sel.End.Char+1 <= len(lineEnd.Value) {
			lineEnd.Value = lineEnd.Value[sel.End.Char+1:]
		}

		if len(lineEnd.Value) == 0 {
			ctx.Buf.Lines.Remove(lineEnd)
		}

		lineJoinNext(ctx.Buf, lineStart)

		ctx.Buf.Cursor.Line = sel.Start.Line
		if lineStart != nil && sel.Start.Char < len(lineStart.Value) {
			cursorGotoChar(ctx.Buf, sel.Start.Char)
		} else {
			CmdGotoLineEnd(ctx)
		}
	}
}

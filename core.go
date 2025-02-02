package mcwig

import (
	"fmt"
	"strings"
	"unicode"
)

const minVisibleLines = 6

func lineJoinNext(buf *Buffer, line *Element[Line]) {
	next := line.Next()
	if next == nil {
		return
	}
	line.Value = append(line.Value, next.Value...)
	buf.Lines.Remove(next)
}

func CmdEnterInsertMode(ctx Context) {
	line := CursorLine(ctx.Buf)
	if line == nil {
		return
	}

	ctx.Buf.TxStart()

	if len(line.Value) == 0 {
		ctx.Buf.Cursor.Char++
	}

	ctx.Buf.SetMode(MODE_INSERT)
}

func CmdExitInsertMode(ctx Context) {
	defer func() {
		ctx.Buf.SetMode(MODE_NORMAL)
		ctx.Buf.Selection = nil
	}()

	CmdCursorLeft(ctx)
	line := CursorLine(ctx.Buf)
	if ctx.Buf.Cursor.Char >= len(line.Value) {
		CmdGotoLineEnd(ctx)
	}

	ctx.Buf.TxEnd()

	// TODO: this is ugly
	if ctx.Buf.Highlighter != nil {
		ctx.Buf.Highlighter.Build()
	}
}

func CmdInsertModeAfter(ctx Context) {
	ctx.Buf.Cursor.Char++
	CmdEnterInsertMode(ctx)
}

func CmdJoinNextLine(ctx Context) {
	CmdGotoLineEnd(ctx)

	if ctx.Buf.TxStart() {
		defer ctx.Buf.TxEnd()
	}

	lineJoinNext(ctx.Buf, CursorLine(ctx.Buf))
}

func CmdReplaceChar(ctx Context) func(Context) {
	return func(ctx Context) {
		if ctx.Buf.TxStart() {
			defer ctx.Buf.TxEnd()
		}

		c := []rune(ctx.Char)
		line := CursorLine(ctx.Buf)
		line.Value[ctx.Buf.Cursor.Char] = c[0]
	}
}

func CmdDeleteCharForward(ctx Context) {
	if ctx.Buf.TxStart() {
		defer ctx.Buf.TxEnd()
	}

	line := CursorLine(ctx.Buf)
	if len(line.Value) == 0 {
		CmdGotoLineEnd(ctx)
		lineJoinNext(ctx.Buf, line)
		CmdCursorBeginningOfTheLine(ctx)
		return
	}

	line.Value = append(line.Value[:ctx.Buf.Cursor.Char], line.Value[ctx.Buf.Cursor.Char+1:]...)

	if ctx.Buf.Cursor.Char >= len(line.Value) {
		CmdCursorLeft(ctx)
	}
}

func CmdDeleteCharBackward(ctx Context) {
	if ctx.Buf.Cursor.Line == 0 && ctx.Buf.Cursor.Char == 0 {
		return
	}

	if ctx.Buf.Cursor.Char == 0 {
		return
	}

	if ctx.Buf.TxStart() {
		defer ctx.Buf.TxEnd()
	}

	line := CursorLine(ctx.Buf)

	if ctx.Buf.Cursor.Char >= len(line.Value) {
		line.Value = line.Value[:len(line.Value)-1]
		CmdCursorLeft(ctx)
		return
	}

	CmdCursorLeft(ctx)
	CmdDeleteCharForward(ctx)
}

func CmdAppendLine(ctx Context) {
	CmdGotoLineEnd(ctx)
	CmdInsertModeAfter(ctx)
}

func CmdNewLine(ctx Context) {
	if ctx.Buf.TxStart() {
		defer ctx.Buf.TxEnd()
	}
	line := CursorLine(ctx.Buf)

	// EOL
	if (ctx.Buf.Cursor.Char) >= len(line.Value) {
		ctx.Buf.Lines.insertValueAfter(Line{}, line)
		ctx.Buf.Cursor.Line++
		ctx.Buf.Cursor.Char = 1
		ctx.Buf.Cursor.PreserveCharPosition = 0
		return
	}

	// split line
	tmpData := make([]rune, len(line.Value[ctx.Buf.Cursor.Char:]))
	copy(tmpData, line.Value[ctx.Buf.Cursor.Char:])
	line.Value = line.Value[:ctx.Buf.Cursor.Char]
	ctx.Buf.Lines.insertValueAfter(tmpData, line)

	CmdCursorLineDown(ctx)
	CmdCursorBeginningOfTheLine(ctx)
}

func CmdLineOpenBelow(ctx Context) {
	CmdGotoLineEnd(ctx)
	CmdInsertModeAfter(ctx)
	CmdNewLine(ctx)
	CmdIndent(ctx)
}

func CmdLineOpenAbove(ctx Context) {
	if ctx.Buf.Cursor.Line == 0 {
		CmdInsertModeAfter(ctx)
		ctx.Buf.Lines.PushFront(Line{})
		CmdCursorBeginningOfTheLine(ctx)
		return
	}
	CmdCursorLineUp(ctx)
	CmdLineOpenBelow(ctx)
}

func CmdDeleteLine(ctx Context) {
	if ctx.Buf.TxStart() {
		defer ctx.Buf.TxEnd()
	}

	CmdVisualLineMode(ctx)
	ctx.Buf.Selection.End.Line = ctx.Buf.Selection.Start.Line + int(ctx.Count)
	ctx.Buf.Selection.End.Char = len(CursorLineByNum(ctx.Buf, ctx.Buf.Selection.End.Line).Value)
	SelectionDelete(ctx)
	CmdNormalMode(ctx)
}

func CmdDeleteWord(ctx Context) {
	if ctx.Buf.TxStart() {
		defer ctx.Buf.TxEnd()
	}
	_, end := TextObjectWord(ctx.Buf, false)
	ctx.Buf.Selection = &Selection{
		Start: ctx.Buf.Cursor,
		End:   Cursor{Line: ctx.Buf.Cursor.Line, Char: end},
	}
	SelectionDelete(ctx)
}

func CmdChangeWord(ctx Context) {
	_, end := TextObjectWord(ctx.Buf, false)
	ctx.Buf.Selection = &Selection{
		Start: ctx.Buf.Cursor,
		End:   Cursor{Line: ctx.Buf.Cursor.Line, Char: end},
	}
	CmdEnterInsertMode(ctx)
	SelectionDelete(ctx)
}

func CmdChangeWORD(ctx Context) {
	start, end := TextObjectWord(ctx.Buf, true)
	ctx.Buf.Cursor.Char = start
	ctx.Buf.Selection = &Selection{
		Start: ctx.Buf.Cursor,
		End:   Cursor{Line: ctx.Buf.Cursor.Line, Char: end},
	}
	CmdEnterInsertMode(ctx)
	SelectionDelete(ctx)
}

func CmdChangeTo(_ Context) func(Context) {
	return func(ctx Context) {
		SelectionStart(ctx.Buf)
		CmdForwardToChar(ctx)(ctx)
		SelectionStop(ctx.Buf)
		CmdEnterInsertMode(ctx)
		SelectionDelete(ctx)
	}
}

func CmdChangeBefore(_ Context) func(Context) {
	return func(ctx Context) {
		SelectionStart(ctx.Buf)
		CmdForwardBeforeChar(ctx)(ctx)
		SelectionStop(ctx.Buf)
		CmdEnterInsertMode(ctx)
		SelectionDelete(ctx)
	}
}

func CmdChangeEndOfLine(ctx Context) {
	SelectionStart(ctx.Buf)
	CmdGotoLineEnd(ctx)
	SelectionStop(ctx.Buf)
	CmdEnterInsertMode(ctx)
	SelectionDelete(ctx)
}

func CmdChangeLine(ctx Context) {
	CmdInsertModeAfter(ctx)
	line := CursorLine(ctx.Buf)
	line.Value = nil
}

func CmdDeleteTo(_ Context) func(Context) {
	return func(ctx Context) {
		if ctx.Buf.TxStart() {
			defer ctx.Buf.TxEnd()
		}
		SelectionStart(ctx.Buf)
		CmdForwardToChar(ctx)(ctx)
		SelectionStop(ctx.Buf)
		SelectionDelete(ctx)
	}
}

func CmdDeleteBefore(ctx Context) {
	if ctx.Buf.TxStart() {
		defer ctx.Buf.TxEnd()
	}
	SelectionStart(ctx.Buf)
	CmdForwardBeforeChar(ctx)(ctx)
	SelectionDelete(ctx)
}

func CmdSelectionChange(ctx Context) {
	CmdEnterInsertMode(ctx)
	SelectionDelete(ctx)
}

func CmdToggleComment(ctx Context) {
	if ctx.Buf.TxStart() {
		defer ctx.Buf.TxEnd()
	}
	comment := "//"

	cmComment := func(line *Element[Line]) {
		idx := 0
		for i, c := range line.Value {
			if !unicode.IsSpace(c) {
				idx = i
				break
			}
		}
		tmpData := make([]rune, 0, len(line.Value)+len(comment)+1)
		tmpData = append(tmpData, line.Value[:idx]...)
		tmpData = append(tmpData, []rune(comment)...)
		tmpData = append(tmpData, rune(' '))
		tmpData = append(tmpData, line.Value[idx:]...)
		line.Value = tmpData

	}

	cmUncomment := func(line *Element[Line]) {
		var r string
		if strings.HasPrefix(strings.TrimSpace(string(line.Value)), comment+" ") {
			r = strings.Replace(string(line.Value), comment+" ", "", 1)
		} else {
			r = strings.Replace(string(line.Value), comment, "", 1)
		}
		line.Value = []rune(r)
	}

	line := CursorLine(ctx.Buf)
	if strings.HasPrefix(strings.TrimSpace(string(line.Value)), comment) {
		cmUncomment(line)
	} else {
		cmComment(line)
	}
}

func CmdSelectionDelete(ctx Context) {
	if ctx.Buf.TxStart() {
		defer ctx.Buf.TxEnd()
	}
	SelectionDelete(ctx)
}

func CmdSaveFile(ctx Context) {
	err := ctx.Buf.Save()
	var msg string
	if err == nil {
		msg = fmt.Sprintf("Saved file %s. Lines: %d.", ctx.Buf.FilePath, ctx.Buf.Lines.Len)
	} else {
		msg = err.Error()
	}

	ctx.Editor.LogMessage(msg)
	ctx.Editor.EchoMessage(msg)
}

func CmdKillBuffer(ctx Context) {
	buffers := ctx.Editor.Buffers
	if len(buffers) == 0 {
		return
	}

	// creates [No Name] buffer
	defer ctx.Editor.ActiveBuffer()

	// remove from buffers list
	// ands moves to the next buffer
	for i, b := range buffers {
		if b == ctx.Buf {
			buffers = append(buffers[:i], buffers[i+1:]...)
			if len(buffers) > 0 {
				idx := i - 1
				if idx < 0 {
					idx = 0
				}
				ctx.Editor.ActiveWindow().VisitBuffer(buffers[idx])
			}
		}
	}

	// cleanup all nodes
	{
		l := ctx.Buf.Lines.First()
		for l != nil {
			next := l.Next()
			l.Value = nil
			ctx.Buf.Lines.Remove(l)
			l = next
		}
	}

	ctx.Editor.Lsp.DidClose(ctx.Buf)
}

func CmdIndentOrComplete(ctx Context) {
	ctx.Editor.Lsp.Completion(ctx.Buf)
}

func CmdChangeInsideBlock(ctx Context) {
	switch ctx.Char {
	case "w":
		CmdChangeWORD(ctx)
	case "(", "[", "{", "'", "\"":
		found, sel, cur := TextObjectBlock(ctx.Buf, rune(ctx.Char[0]), false) // TODO: handle unicode
		if !found {
			return
		}
		ctx.Buf.Selection = sel
		ctx.Buf.Cursor = cur
		CmdEnterInsertMode(ctx)
		SelectionDelete(ctx)
	}
}

func CmdUndo(ctx Context) {
	ctx.Buf.UndoRedo.Undo()
}

func CmdRedo(ctx Context) {
	ctx.Buf.UndoRedo.Redo()
}

func CmdExit(ctx Context) {
	ctx.Editor.ExitCh <- 1
}

func CmdVisualMode(ctx Context) {
	SelectionStart(ctx.Buf)
	ctx.Buf.SetMode(MODE_VISUAL)
}

func CmdNormalMode(ctx Context) {
	if ctx.Buf.Mode() == MODE_INSERT {
		line := CursorLine(ctx.Buf)
		CmdCursorLeft(ctx)
		if ctx.Buf.Cursor.Char >= len(line.Value) {
			CmdGotoLineEnd(ctx)
		}
	}
	ctx.Buf.SetMode(MODE_NORMAL)
	ctx.Buf.Selection = nil
}

func CmdVisualLineMode(ctx Context) {
	line := CursorLine(ctx.Buf)
	SelectionStart(ctx.Buf)
	ctx.Buf.Selection.Start.Char = 0
	ctx.Buf.Selection.End.Char = len(line.Value) - 1
	ctx.Buf.SetMode(MODE_VISUAL_LINE)
}

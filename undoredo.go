package mcwig

import (
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
)

type Transaction struct {
	buf *Buffer

	before string
}

func NewTx(b *Buffer) *Transaction {
	return &Transaction{
		buf: b,
	}
}

func (tx *Transaction) Start() {
	tx.before = tx.buf.String()
}

func (tx *Transaction) End() {
	aString := tx.before
	bString := tx.buf.String()

	apply := myers.ComputeEdits(span.URIFromPath("a.txt"), aString, bString)
	undo := myers.ComputeEdits(span.URIFromPath("b.txt"), bString, aString)
	tx.buf.UndoRedo.Push(EditDiff{
		apply: apply,
		undo:  undo,
	})
}

type EditDiff struct {
	apply []gotextdiff.TextEdit
	undo  []gotextdiff.TextEdit
}

type UndoRedo struct {
	Buf      *Buffer
	History  []EditDiff
	Position int
}

func NewUndoRedo(buf *Buffer) *UndoRedo {
	return &UndoRedo{
		Buf:      buf,
		Position: -1,
		History:  make([]EditDiff, 0, 256),
	}
}

func (u *UndoRedo) checkPosition() bool {
	if u.Position > len(u.History) || u.Position < 0 {
		return false
	}

	return true
}

func (u *UndoRedo) Push(edits EditDiff) {
	if len(edits.apply) > 0 || len(edits.undo) > 0 {

		// we are back in history. remove all "forward" edits
		if u.Position >= 0 && u.Position != len(u.History)-1 {
			u.History = u.History[:u.Position+1]
		}

		u.History = append(u.History, edits)
		u.Position = len(u.History) - 1
	}
}

func (u *UndoRedo) Undo() {
	if !u.checkPosition() || u.Position < 0 {
		return
	}

	edits := u.History[u.Position].undo
	if len(edits) == 0 {
		return
	}

	res := gotextdiff.ApplyEdits(u.Buf.String(), edits)
	u.Buf.ResetLines()
	u.Buf.Append(res)

	if u.Position >= 0 {
		u.Position--
	}
}

func (u *UndoRedo) Redo() {
	if !u.checkPosition() {
		return
	}

	edits := u.History[u.Position].apply
	if len(edits) == 0 {
		return
	}
	res := gotextdiff.ApplyEdits(u.Buf.String(), edits)
	u.Buf.ResetLines()
	u.Buf.Append(res)

	if u.Position < len(u.History)-1 {
		u.Position++
	}
}

package ui

import (
	"math"

	"github.com/firstrow/mcwig"
)

type AutocompleteWidget struct {
	ctx            mcwig.Context
	triggerPos     mcwig.Cursor
	keymap         *mcwig.KeyHandler
	pos            mcwig.Position
	items          mcwig.CompletionItems
	eventsListener <-chan mcwig.Event
	activeItem     int
}

func (u *AutocompleteWidget) Plane() mcwig.RenderPlane {
	return mcwig.PlaneWin
}

func AutocompleteInit(
	ctx mcwig.Context,
	pos mcwig.Position,
	items mcwig.CompletionItems,
) *AutocompleteWidget {
	if len(items.Items) == 0 {
		return nil
	}

	widget := &AutocompleteWidget{
		ctx:        ctx,
		pos:        pos,
		items:      items,
		activeItem: 0,
	}

	widget.keymap = mcwig.NewKeyHandler(mcwig.ModeKeyMap{
		mcwig.MODE_INSERT: mcwig.KeyMap{
			"Esc": func(ctx mcwig.Context) {
				widget.Close()
			},
			"Tab": func(ctx mcwig.Context) {
				if widget.activeItem < len(widget.items.Items)-1 {
					widget.activeItem++
				}
			},
			"Backtab": func(ctx mcwig.Context) {
				if widget.activeItem > 0 {
					widget.activeItem--
				}
			},
			"Enter": widget.selectItem,
		},
	})

	widget.eventsListener = ctx.Editor.Events.Subscribe()
	go func() {
		for {
			select {
			case event := <-widget.eventsListener:
				event.Wg.Done()

				switch e := event.Msg.(type) {
				case mcwig.EventTextChange:
					widget.activeItem = 0
					widget.items = ctx.Editor.Lsp.Completion(e.Buf)
					if len(widget.items.Items) == 0 {
						widget.Close()
					}

					ctx.Editor.Redraw()
				}
			}
		}
	}()

	ctx.Editor.PushUi(widget)

	return widget
}

func (w *AutocompleteWidget) Close() {
	w.ctx.Editor.PopUi()
	w.ctx.Editor.Events.Unsubscribe(w.eventsListener)
}

func (w *AutocompleteWidget) Mode() mcwig.Mode {
	return mcwig.MODE_INSERT
}

func (w *AutocompleteWidget) Keymap() *mcwig.KeyHandler {
	return w.keymap
}

func (w *AutocompleteWidget) selectItem(ctx mcwig.Context) {
	defer w.Close()

	line := mcwig.CursorLine(ctx.Buf)
	item := w.items.Items[w.activeItem]
	text := item.TextEdit.NewText
	pos := item.TextEdit.Insert.Start.Character

	mcwig.TextDelete(ctx.Buf, &mcwig.Selection{
		Start: mcwig.Cursor{
			Line: item.TextEdit.Replace.Start.Line,
			Char: item.TextEdit.Replace.Start.Character,
		},
		End: mcwig.Cursor{
			Line: item.TextEdit.Replace.End.Line,
			Char: item.TextEdit.Replace.End.Character,
		},
	})

	chpos := len(text)
	if item.InsertTextFormat == 2 {
		text, chpos = mcwig.SnippetProcessString(text)
	}

	mcwig.TextInsert(ctx.Buf, line, int(pos), text)
	ctx.Buf.Cursor.Char = item.TextEdit.Replace.Start.Character + chpos
}

func (w *AutocompleteWidget) Render(view mcwig.View) {
	x := w.pos.Char + 2
	y := w.pos.Line - w.ctx.Buf.ScrollOffset + 1

	maxItems := min(10, len(w.items.Items))

	_, winHeight := view.Size()
	if y+maxItems >= winHeight {
		y -= maxItems + 2
	}

	drawBoxNoBorder(view, w.pos.Char, y, 50, maxItems, mcwig.Color("ui.menu"))

	// pagination
	pageSize := maxItems
	pageNumber := math.Ceil(float64(w.activeItem+1)/float64(pageSize)) - 1
	startIndex := int(pageNumber) * pageSize
	endIndex := startIndex + pageSize
	if endIndex > len(w.items.Items) {
		endIndex = len(w.items.Items)
	}
	dataset := w.items.Items[startIndex:endIndex]

	for i, row := range dataset {
		st := mcwig.Color("ui.menu")
		if i+startIndex == w.activeItem {
			st = mcwig.Color("ui.menu.selected")
		}

		label := row.Label
		view.SetContent(x, y, label, st)
		if i >= maxItems {
			return
		}
		y++
	}
}


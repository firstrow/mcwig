package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/gdamore/tcell/v2"

	"github.com/firstrow/mcwig"
	"github.com/firstrow/mcwig/autocomplete"
	"github.com/firstrow/mcwig/config"
	"github.com/firstrow/mcwig/metrics"
	"github.com/firstrow/mcwig/render"
)

func main() {
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	tscreen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}

	err = tscreen.Init()
	if err != nil {
		panic(err)
	}
	tscreen.Sync()
	w, h := tscreen.Size()

	editor := mcwig.NewEditor(
		render.NewMView(tscreen, 0, 0, w, h),
		mcwig.NewKeyHandler(config.DefaultKeyMap()),
	)
	editor.AutocompleteTrigger = autocomplete.Register(editor)

	args := os.Args
	if len(args) > 1 {
		buf := editor.OpenFile(args[1])
		editor.ActiveWindow().VisitBuffer(buf)
	} else {
		mcwig.CmdNewBuffer(editor.NewContext())
	}

	renderer := render.New(editor, tscreen)

	go func() {
		for {
			switch ev := tscreen.PollEvent().(type) {
			case *tcell.EventResize:
				tscreen.Sync()
				w, h := tscreen.Size()
				editor.View.Resize(0, 0, w, h)
				renderer.Render()
			case *tcell.EventKey:
				metrics.Track("handler", func() {
					editor.HandleInput(ev)
				})

				metrics.Track("render", func() {
					renderer.Render()
				})

				// renderer.RenderMetrics(metrics.Get())
			case *tcell.EventError:
				fmt.Println("error:", ev)
				return
			}
		}
	}()

	go func() {
		for {
			<-editor.RedrawCh
			renderer.Render()
		}
	}()

	go func() {
		for {
			<-editor.ScreenSyncCh
			tscreen.Sync()
		}
	}()

	<-editor.ExitCh
	tscreen.Clear()
	tscreen.Fini()
}


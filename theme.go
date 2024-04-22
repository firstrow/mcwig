package mcwig

import (
	"github.com/gdamore/tcell/v2"
)

type termColor struct {
	bg int
	fg int
}

var (
	bg  = 0x2b3339
	fg  = 0xd3c6aa
	sep = 0x7a8478
	bg2 = 0x445055
)

var theme = map[string]termColor{
	"bg":                {bg: bg},
	"fg":                {fg: fg},
	"text":              {bg: bg, fg: fg},
	"sepatator":         {fg: sep},
	"cursor":            {bg: fg, fg: bg},
	"statusline.active": {bg: fg, fg: bg},
	"statusline.normal": {bg: bg2, fg: bg},
	"statusline.insert": {bg: fg, fg: bg},
}

var styles = map[string]tcell.Style{}
var defaultStyle tcell.Style

func ThemeInit() {
	defaultStyle = tcell.StyleDefault.Foreground(tcell.NewHexColor(int32(fg))).Background(tcell.NewHexColor(int32(bg)))
}

func color(name string) tcell.Style {
	val, ok := styles[name]
	if ok {
		return val
	}

	val = defaultStyle

	c := theme[name]
	if c.fg > 0 {
		val = val.Foreground(tcell.NewHexColor(int32(c.fg)))
	}
	if c.bg > 0 {
		val = val.Background(tcell.NewHexColor(int32(c.bg)))
	}

	styles[name] = val
	return val
}

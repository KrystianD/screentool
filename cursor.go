package main

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

var currentCursor = ""

func setWindowCursor(win *gtk.Window, name string) {
	if name != currentCursor {
		currentCursor = name
		d, _ := win.GetDisplay()
		c, _ := gdk.CursorNewFromName(d, name)
		dw, _ := win.GetWindow()
		dw.SetCursor(c)
	}
}

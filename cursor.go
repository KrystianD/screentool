package main

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

var currentCursor = ""

func setWindowCursor(win *gtk.Window, name string) {
	if name != currentCursor {
		currentCursor = name
		display, _ := win.GetDisplay()
		cursor, _ := gdk.CursorNewFromName(display, name)
		window, _ := win.GetWindow()
		window.SetCursor(cursor)
	}
}

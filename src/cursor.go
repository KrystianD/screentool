package main

import (
	"log"

	"github.com/gotk3/gotk3/gdk"
)

var currentCursor = ""

func setWindowCursor(name string) {
	if name != currentCursor {
		currentCursor = name

		display, err := gdk.DisplayGetDefault()
		if err != nil {
			log.Fatal(err)
		}

		cursor, err := gdk.CursorNewFromName(display, name)
		if err != nil {
			log.Fatal(err)
		}

		screen, err := gdk.ScreenGetDefault()
		if err != nil {
			log.Fatal(err)
		}

		rootWindow, err := screen.GetRootWindow()
		if err != nil {
			log.Fatal(err)
		}

		rootWindow.SetCursor(cursor)
	}
}

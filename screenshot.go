package main

import (
	"github.com/gotk3/gotk3/gdk"

	. "./utils"
)

func captureScreenshot(rect Rectangle) (*gdk.Pixbuf, error) {
	var err error

	screen, err := gdk.ScreenGetDefault()
	if err != nil {
		return nil, err
	}

	rootWindow, err := screen.GetRootWindow()
	if err != nil {
		return nil, err
	}

	x, y, w, h := rect.GetXYWH()
	pixbuf, err := rootWindow.PixbufGetFromWindow(x, y, w, h)
	if err != nil {
		return nil, err
	}

	return pixbuf, nil
}

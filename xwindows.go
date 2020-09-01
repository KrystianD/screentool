package main

import (
	"log"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/gotk3/gotk3/gdk"

	. "./utils"
)

type DesktopWindow struct {
	Geometry Rectangle
	window   *xwindow.Window
}

func (window *DesktopWindow) RaiseToFront() {
	window.window.Stack(xproto.StackModeAbove)
}

func contains(values []string, value string) bool {
	for _, x := range values {
		if x == value {
			return true
		}
	}
	return false
}

func getCurrentToplevelWindows() []DesktopWindow {
	var windows []DesktopWindow

	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	clientIds, err := ewmh.ClientListStackingGet(X)
	if err != nil {
		log.Fatal(err)
	}

	curDesktop, _ := ewmh.CurrentDesktopGet(X)

	for _, clientId := range clientIds {
		win := xwindow.New(X, clientId)
		geom, err := win.DecorGeometry()
		if err != nil {
			continue
		}

		states, err := ewmh.WmStateGet(X, win.Id)
		if err != nil {
			continue
		}

		if contains(states, "_NET_WM_STATE_HIDDEN") {
			continue
		}

		winDesktop, _ := ewmh.WmDesktopGet(X, win.Id)

		if curDesktop != winDesktop && winDesktop != 0xFFFFFFFF {
			continue
		}

		windows = append(windows, DesktopWindow{
			Geometry: NewRectangleFromXYWH(geom.X(), geom.Y(), geom.Width(), geom.Height()),
			window:   win,
		})
	}

	return windows
}

func getRootWindowRect() Rectangle {
	var err error

	screen, err := gdk.ScreenGetDefault()
	if err != nil {
		log.Fatal(err)
	}

	rootWindow, err := screen.GetRootWindow()
	if err != nil {
		log.Fatal(err)
	}

	return NewRectangleFromXYWH(0, 0, rootWindow.WindowGetWidth(), rootWindow.WindowGetHeight())
}

func getMousePosition() Point {
	var err error

	screen, err := gdk.ScreenGetDefault()
	if err != nil {
		log.Fatal(err)
	}
	device, err := screen.GetDisplay()
	if err != nil {
		log.Fatal(err)
	}

	seat, err := device.GetDefaultSeat()
	if err != nil {
		log.Fatal(err)
	}

	pointer, err := seat.GetPointer()
	if err != nil {
		log.Fatal(err)
	}
	_ = pointer

	var x, y int
	var ss *gdk.Screen
	err = pointer.GetPosition(&ss, &x, &y)
	//if err != nil {
	//	log.Fatal(err)
	//}
	pointer.Unref()

	return Point{X: x, Y: y}
}

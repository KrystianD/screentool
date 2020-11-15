package main

import (
	"log"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xprop"
	"github.com/BurntSushi/xgbutil/xrect"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/gotk3/gotk3/gdk"

	. "./utils"
)

type DesktopWindow struct {
	Geometry Rectangle
	window   *xwindow.Window
}

func (window *DesktopWindow) RaiseToFront() {
	X, err := xgbutil.NewConn()
	if err == nil {
		win := xwindow.New(X, window.window.Id)
		ewmh.RestackWindow(X, win.Id)
	}
}

func contains(values []string, value string) bool {
	for _, x := range values {
		if x == value {
			return true
		}
	}
	return false
}

func NewRectangleFromXRect(rect xrect.Rect) Rectangle {
	return NewRectangleFromXYWH(rect.X(), rect.Y(), rect.Width(), rect.Height())
}

func XPropGetPropertyLRTB(xu *xgbutil.XUtil, win xproto.Window, atomName string) (bool, int, int, int, int) {
	prop, err := xprop.GetProperty(xu, win, atomName)
	if err != nil || prop == nil {
		return false, 0, 0, 0, 0
	}

	nums, err := xprop.PropValNums(prop, err)
	if err != nil || len(nums) != 4 {
		return false, 0, 0, 0, 0
	}

	left := int(nums[0])
	right := int(nums[1])
	top := int(nums[2])
	bottom := int(nums[3])
	return true, left, right, top, bottom
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

		geom, err := xproto.GetGeometry(X.Conn(), xproto.Drawable(win.Id)).Reply()
		if err != nil {
			continue
		}

		var geomRect Rectangle
		coord, err := xproto.TranslateCoordinates(X.Conn(), win.Id, X.RootWin(), 0, 0).Reply()
		if err == nil {
			geomRect = NewRectangleFromXYWH(int(coord.DstX), int(coord.DstY), int(geom.Width), int(geom.Height))

			ok, leftExtent, rightExtent, topExtent, bottomExtent := XPropGetPropertyLRTB(X, win.Id, "_NET_FRAME_EXTENTS")
			if ok {
				geomRect.MoveLTRB(-leftExtent, -topExtent, rightExtent, bottomExtent)
			}

			ok, leftExtent, rightExtent, topExtent, bottomExtent = XPropGetPropertyLRTB(X, win.Id, "_GTK_FRAME_EXTENTS")
			if ok {
				geomRect.MoveLTRB(leftExtent, topExtent, -rightExtent, -bottomExtent)
			}
		} else {
			// fallback
			geomDecor, err := win.DecorGeometry()
			if err != nil {
				continue
			}

			geomRect = NewRectangleFromXRect(geomDecor)
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
			Geometry: geomRect,
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

	device, err := gdk.DisplayGetDefault()
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

	var x, y int
	var s *gdk.Screen
	err = pointer.GetPosition(&s, &x, &y)
	if err != nil {
		log.Fatal(err)
	}
	pointer.Unref()

	return Point{X: x, Y: y}
}

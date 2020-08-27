package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/lestrrat-go/strftime"

	"./annotations"
	. "./utils"
)

type State int

const (
	Hovering        State = 0
	SelectingRegion State = 1
	QuickAnnotating State = 2
)

var toplevelWindows []DesktopWindow
var desktopRect Rectangle

var mainWindow *gtk.Window

var state = Hovering
var startPoint Point
var mousePos Point
var hoveredWindow *DesktopWindow
var hoveredWindowRect Rectangle
var finalRect Rectangle
var capturedPixbuf *gdk.Pixbuf

func NewPointFromEventButton(ev *gdk.EventButton) Point {
	x, y := ev.MotionVal()
	return Point{
		X: int(x),
		Y: int(y),
	}
}
func NewPointFromEventMotion(ev *gdk.EventMotion) Point {
	x, y := ev.MotionVal()
	return Point{
		X: int(x),
		Y: int(y),
	}
}

func saveScreenshot(pixbuf *gdk.Pixbuf) {
	go func() {
		name, _ := strftime.Format("%Y-%m-%d_%H%M%S", time.Now())
		screenshotsDir := os.ExpandEnv("$HOME/screenshots")
		if FileExists(screenshotsDir) {
			_ = pixbuf.SavePNG(fmt.Sprintf("%s/%s.png", screenshotsDir, name), 9)
			fmt.Println("saved")
		}
	}()
}

func captureScreen(rect Rectangle, controlPressed, shiftPressed bool) {
	var err error

	finalRect = rect
	capturedPixbuf, err = captureScreenshot(rect)

	if err != nil {
		log.Fatal(err)
	}

	if controlPressed {
	} else if shiftPressed {
		state = QuickAnnotating
		mainWindow.Show()
		mainWindow.Present()
		mainWindow.QueueDraw()

		annotations.InitAnnotations(finalRect.Size())

		updateCursor()
	} else {
		saveScreenshotAndFinish()
	}
}

func saveScreenshotAndFinish() {
	var pixbuf *gdk.Pixbuf

	if annotations.Has() {
		var finalSurface = cairo.CreateImageSurface(cairo.FORMAT_ARGB32, finalRect.Width(), finalRect.Height())
		var finalCtx = cairo.Create(finalSurface)

		gtk.GdkCairoSetSourcePixBuf(finalCtx, capturedPixbuf, 0, 0)
		finalCtx.Paint()

		annotations.Draw(finalCtx, 0, 0)

		pixbuf, _ = gdk.PixbufGetFromSurface(finalSurface, 0, 0, finalRect.Width(), finalRect.Height())
	} else {
		pixbuf = capturedPixbuf
	}

	mainWindow.Hide()
	saveScreenshot(pixbuf)
	saveToClipboardAndWait(pixbuf, func() {
		gtk.MainQuit()
	})
}

func updateCursor() {
	if state == Hovering {
		setWindowCursor(mainWindow, "crosshair")
	}

	if state == SelectingRegion {
		diff := math.Min(math.Abs(float64(mousePos.X-startPoint.X)), math.Abs(float64(mousePos.Y-startPoint.Y)))

		if diff < 20 {
			setWindowCursor(mainWindow, "crosshair")
		} else {
			if mousePos.X < startPoint.X {
				if mousePos.Y < startPoint.Y {
					setWindowCursor(mainWindow, "ul_angle")
				} else {
					setWindowCursor(mainWindow, "ll_angle")
				}
			} else {
				if mousePos.Y < startPoint.Y {
					setWindowCursor(mainWindow, "ur_angle")
				} else {
					setWindowCursor(mainWindow, "lr_angle")
				}
			}
		}
	}

	if state == QuickAnnotating {
		setWindowCursor(mainWindow, "default")
	}
}

func findWindowUnderCursor() {
	l, t, r, b := desktopRect.GetLTRB()
	var FullscreenEdgeDistance = 10

	if mousePos.X <= l+FullscreenEdgeDistance || mousePos.Y <= t+FullscreenEdgeDistance ||
		mousePos.X >= r-1-FullscreenEdgeDistance || mousePos.Y >= b-1-FullscreenEdgeDistance {
		hoveredWindow = nil
		hoveredWindowRect = desktopRect
	} else {
		for i := range toplevelWindows {
			var desktopWindow = toplevelWindows[len(toplevelWindows)-1-i]
			if desktopWindow.Geometry.Contains(mousePos) {
				hoveredWindow = &desktopWindow
				hoveredWindowRect = desktopWindow.Geometry
				break
			}
		}
	}
}

func onDraw(ctx *cairo.Context) {
	ctx.SetOperator(cairo.OPERATOR_SOURCE)

	if state == Hovering {
		ctx.SetSourceRGBA(0.0, 0.0, 0.0, 0)
		ctx.Paint()

		ctx.SetSourceRGB(1.0, 1.0, 0.0)
		ctx.SetLineWidth(5)
		hoveredWindowRect.SetToCairo(ctx)
		ctx.Stroke()
	}

	if state == SelectingRegion {
		ctx.SetSourceRGBA(0.0, 0.0, 0.0, 0.25)
		ctx.Paint()

		ctx.SetSourceRGB(0.0, 0.0, 1.0)
		ctx.SetLineWidth(2)
		finalRect.SetToCairo(ctx)
		ctx.Stroke()

		ctx.SetSourceRGBA(0.0, 0.0, 0.0, 0.0)
		finalRect.SetToCairo(ctx)
		ctx.Fill()
	}

	if state == QuickAnnotating {
		ctx.SetSourceRGBA(0.0, 0.0, 0.0, 0.55)
		ctx.Paint()

		ctx.SetSourceRGB(0.0, 0.0, 0.5)
		ctx.SetLineWidth(1)
		finalRect.SetToCairo(ctx)
		ctx.Stroke()

		finalRect.SetToCairo(ctx)
		ctx.Clip()
		gtk.GdkCairoSetSourcePixBuf(ctx, capturedPixbuf, float64(finalRect.X()), float64(finalRect.Y()))
		ctx.Paint()

		annotations.Draw(ctx, finalRect.X(), finalRect.Y())
	}
}

func onMousePrimaryPressed(event *gdk.EventButton) {
	mousePos = NewPointFromEventButton(event)

	if state == Hovering {
		startPoint = mousePos
		state = SelectingRegion
		finalRect = NewRectangleFromXYWH(0, 0, 0, 0)
	}

	if state == QuickAnnotating {
		mousePosRelative := Point{
			X: mousePos.X - finalRect.X(),
			Y: mousePos.Y - finalRect.Y(),
		}

		annotations.HandleMousePressed(mousePosRelative)
	}

	updateCursor()
	mainWindow.QueueDraw()
}

func onMouseMove(event *gdk.EventMotion) {
	mousePos = NewPointFromEventMotion(event)

	if state == Hovering {
		findWindowUnderCursor()
	}

	if state == SelectingRegion {
		finalRect = NewRectangleFromPoints(startPoint, mousePos)
	}

	if state == QuickAnnotating {
		mousePosRelative := Point{
			X: mousePos.X - finalRect.X(),
			Y: mousePos.Y - finalRect.Y(),
		}

		if (event.State() & gdk.BUTTON1_MASK) > 0 {
			annotations.HandleMouseDrag(mousePosRelative)
		}
	}

	updateCursor()
	mainWindow.QueueDraw()
}

func onMousePrimaryReleased(event *gdk.EventButton) {
	if state == SelectingRegion {
		var controlPressed = (event.State() & uint(gdk.CONTROL_MASK)) > 0
		var shiftPressed = (event.State() & uint(gdk.SHIFT_MASK)) > 0

		if startPoint.ManhattanDistanceTo(mousePos) < 5 {
			mainWindow.Hide()
			if hoveredWindow != nil {
				hoveredWindow.RaiseToFront()
			}

			_, _ = glib.TimeoutAdd(10, func() {
				captureScreen(hoveredWindowRect, controlPressed, shiftPressed)
			})
		} else {
			captureScreen(finalRect, controlPressed, shiftPressed)
		}
	}

	if state == QuickAnnotating {
		annotations.HandleMouseReleased()
	}

	updateCursor()
	mainWindow.QueueDraw()
}

func onMouseSecondaryReleased(event *gdk.EventButton) {
	if state == QuickAnnotating {
		if annotations.Has() {
			annotations.HandleMouseSecondaryReleased()
		} else {
			state = Hovering
		}
	}

	updateCursor()
	mainWindow.QueueDraw()
}

func onKeyReleased(event *gdk.EventKey) {
	if state == QuickAnnotating {
		if event.KeyVal() == gdk.KEY_Shift_L {
			saveScreenshotAndFinish()
		}

		if event.KeyVal() == gdk.KEY_space {
			annotations.CycleTool()
		}
	}

	if event.KeyVal() == gdk.KEY_Escape {
		gtk.MainQuit()
	}
}

func main() {
	var err error

	gtk.Init(nil)

	desktopRect = getRootWindowRect()
	toplevelWindows = getCurrentToplevelWindows()

	mainWindow, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}

	mainWindow.Fullscreen()
	mainWindow.SetKeepAbove(true)
	mainWindow.SetAppPaintable(true)
	mainWindow.SetDecorated(false)
	mainWindow.SetSkipTaskbarHint(true)

	_, _ = mainWindow.Connect("destroy", func() {
		gtk.MainQuit()
	})

	_, _ = mainWindow.Connect("draw", func(window *gtk.Window, context *cairo.Context) {
		onDraw(context)
	})

	_, _ = mainWindow.Connect("button-press-event", func(window *gtk.Window, event *gdk.Event) bool {
		mouseEvent := gdk.EventButtonNewFromEvent(event)
		if mouseEvent.Button() == gdk.BUTTON_PRIMARY {
			onMousePrimaryPressed(mouseEvent)
		}
		return true
	})

	_, _ = mainWindow.Connect("motion-notify-event", func(window *gtk.Window, event *gdk.Event) bool {
		onMouseMove(gdk.EventMotionNewFromEvent(event))
		return true
	})

	_, _ = mainWindow.Connect("button-release-event", func(window *gtk.Window, event *gdk.Event) bool {
		mouseEvent := gdk.EventButtonNewFromEvent(event)
		if mouseEvent.Button() == gdk.BUTTON_PRIMARY {
			onMousePrimaryReleased(mouseEvent)
		}
		if mouseEvent.Button() == gdk.BUTTON_SECONDARY {
			onMouseSecondaryReleased(mouseEvent)
		}
		return true
	})

	_, _ = mainWindow.Connect("key-release-event", func(window *gtk.Window, event *gdk.Event) bool {
		onKeyReleased(gdk.EventKeyNewFromEvent(event))
		return true
	})

	mainWindow.SetEvents(int(gdk.POINTER_MOTION_MASK | gdk.KEY_RELEASE_MASK | gdk.BUTTON_PRESS_MASK))

	// Allow main window to be transparent
	visual, err := mainWindow.GetScreen().GetRGBAVisual()
	if err != nil || visual == nil {
		log.Fatal("Alpha not supported")
	}
	mainWindow.SetVisual(visual)

	// Show and focus main window
	mainWindow.Show()
	mainWindow.Present()

	gtk.Main()
}

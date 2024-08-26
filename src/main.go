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

	"github.com/KrystianD/screentool/src/annotations"
	. "github.com/KrystianD/screentool/src/utils"
)

type State int

const (
	Hovering        State = 0
	SelectingRegion State = 1
	QuickAnnotating State = 2
)

var toplevelWindows []DesktopWindow
var frozenScreen *gdk.Pixbuf
var desktopRect Rectangle

var windows []*gtk.Window

var state = Hovering
var startPoint Point
var mousePos Point
var absoluteMousePos Point
var hoveredWindow *DesktopWindow
var hoveredWindowRect Rectangle
var selectedRegionRect Rectangle
var capturedRect Rectangle
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

func restoreCursor() {
	setWindowCursor("default")
}

func quitApp() {
	restoreCursor()
	gtk.MainQuit()
}

func saveScreenshot(pixbuf *gdk.Pixbuf) {
	name, _ := strftime.Format("%Y-%m-%d_%H%M%S", time.Now())
	screenshotsDir := os.ExpandEnv("$HOME/screenshots")
	if FileExists(screenshotsDir) {
		filePath := fmt.Sprintf("%s/%s.jpg", screenshotsDir, name)

		_ = pixbuf.SaveJPEG(filePath, 90)

		recentManager, err := gtk.RecentManagerGetDefault()
		if err == nil {
			recentManager.AddItem("file://" + filePath)
		}
	}
}

func captureScreen(rect Rectangle, controlPressed, shiftPressed bool) {
	var err error

	if frozenScreen == nil {
		capturedPixbuf, err = captureScreenshot(rect)
		capturedRect = rect
		if err != nil {
			log.Fatal(err)
		}
	} else {
		capturedPixbuf = CropPixbuf(frozenScreen, rect)
		capturedRect = rect
	}

	if controlPressed {
	} else if shiftPressed {
		state = QuickAnnotating

		for _, window := range windows {
			window.Show()
			window.Present()
			window.QueueDraw()
		}

		annotations.InitAnnotations(capturedRect.Size())

		updateCursor()
	} else {
		saveScreenshotAndFinish()
	}
}

func saveScreenshotAndFinish() {
	var pixbuf *gdk.Pixbuf

	if annotations.Empty() {
		pixbuf = capturedPixbuf
	} else {
		var finalSurface = cairo.CreateImageSurface(cairo.FORMAT_ARGB32, capturedRect.Width(), capturedRect.Height())
		var finalCtx = cairo.Create(finalSurface)

		gtk.GdkCairoSetSourcePixBuf(finalCtx, capturedPixbuf, 0, 0)
		finalCtx.Paint()

		annotations.Draw(finalCtx, 0, 0)

		pixbuf, _ = gdk.PixbufGetFromSurface(finalSurface, 0, 0, capturedRect.Width(), capturedRect.Height())
	}

	for _, window := range windows {
		window.Hide()
	}

	var sem = make(chan int, 1)
	go func() {
		saveScreenshot(pixbuf)
		sem <- 0
	}()
	saveToClipboardAndWait(pixbuf, func() {
		restoreCursor()
		<-sem
		quitApp()
	})
}

func redrawAllWindows() {
	for _, window := range windows {
		window.QueueDraw()
	}
}

func updateCursor() {
	if state == Hovering {
		setWindowCursor("crosshair")
	}

	if state == SelectingRegion {
		diff := math.Min(math.Abs(float64(absoluteMousePos.X-startPoint.X)), math.Abs(float64(absoluteMousePos.Y-startPoint.Y)))
		if diff < 20 {
			setWindowCursor("crosshair")
		} else {
			if absoluteMousePos.X < startPoint.X {
				if absoluteMousePos.Y < startPoint.Y {
					setWindowCursor("ul_angle")
				} else {
					setWindowCursor("ll_angle")
				}
			} else {
				if absoluteMousePos.Y < startPoint.Y {
					setWindowCursor("ur_angle")
				} else {
					setWindowCursor("lr_angle")
				}
			}
		}
	}

	if state == QuickAnnotating {
		setWindowCursor("default")
	}
}

func updateWindowUnderCursor() {
	l, t, r, b := desktopRect.GetLTRB()
	var FullscreenEdgeDistance = 10

	mouseOnEdge := absoluteMousePos.X <= l+FullscreenEdgeDistance || absoluteMousePos.Y <= t+FullscreenEdgeDistance ||
		absoluteMousePos.X >= r-1-FullscreenEdgeDistance || absoluteMousePos.Y >= b-1-FullscreenEdgeDistance

	hoveredWindow = nil
	hoveredWindowRect = desktopRect

	if !mouseOnEdge {
		for i := range toplevelWindows {
			var desktopWindow = toplevelWindows[len(toplevelWindows)-1-i]
			if desktopWindow.Geometry.Contains(absoluteMousePos) {
				hoveredWindow = &desktopWindow
				hoveredWindowRect = desktopWindow.Geometry
				break
			}
		}
	}
}

func onDraw(monitor *gdk.Monitor, ctx *cairo.Context) {
	ctx.SetOperator(cairo.OPERATOR_OVER)

	monitorRect := NewRectangleFromGdkRectangle(monitor.GetGeometry())
	monitorX := monitorRect.X()
	monitorY := monitorRect.Y()

	if frozenScreen != nil {
		gtk.GdkCairoSetSourcePixBuf(ctx, frozenScreen, float64(-monitorX), -float64(monitorY))
		ctx.Paint()
	}

	if state == Hovering {
		ctx.SetSourceRGBA(0.0, 0.0, 0.0, 0)
		ctx.Paint()

		ctx.SetSourceRGB(1.0, 1.0, 0.0)
		ctx.SetLineWidth(4)

		hoveredWindowRectOnMonitor := hoveredWindowRect.TranslatedByXY(-monitorX, -monitorY).Shrinked(2)

		hoveredWindowRectOnMonitor.SetToCairo(ctx)
		ctx.Stroke()
	}

	if state == SelectingRegion {
		ctx.SetSourceRGBA(0.0, 0.0, 0.0, 0.25)
		ctx.Paint()

		ctx.Stroke()

		selectedRegionRectOnMonitor := selectedRegionRect.TranslatedByXY(-monitorX, -monitorY)

		ctx.SetSourceRGB(0.0, 0.0, 1.0)
		ctx.SetLineWidth(2)
		selectedRegionRectOnMonitor.SetToCairo(ctx)
		ctx.Stroke()

		if frozenScreen == nil {
			ctx.SetOperator(cairo.OPERATOR_CLEAR)
			selectedRegionRectOnMonitor.SetToCairo(ctx)
			ctx.Fill()
		} else {
			selectedRegionRectOnMonitor.SetToCairo(ctx)
			ctx.Clip()
			gtk.GdkCairoSetSourcePixBuf(ctx, frozenScreen, float64(-monitorX), float64(-monitorY))
			ctx.Paint()
		}
	}

	if state == QuickAnnotating {
		ctx.SetSourceRGBA(0.0, 0.0, 0.0, 0.55)
		ctx.Paint()

		capturedRectOnMonitor := capturedRect.TranslatedByXY(-monitorX, -monitorY)

		ctx.SetSourceRGB(0.0, 0.0, 0.5)
		ctx.SetLineWidth(1)
		capturedRectOnMonitor.SetToCairo(ctx)
		ctx.Stroke()

		capturedRectOnMonitor.SetToCairo(ctx)
		ctx.Clip()
		gtk.GdkCairoSetSourcePixBuf(ctx, capturedPixbuf, float64(capturedRect.X()-monitorX), float64(capturedRect.Y()-monitorY))
		ctx.Paint()

		annotations.Draw(ctx, capturedRect.X()-monitorX, capturedRect.Y()-monitorY)
	}
}

func onMousePrimaryPressed(monitor *gdk.Monitor, event *gdk.EventButton) {
	mousePos = NewPointFromEventButton(event)
	absoluteMousePos = mousePos.TranslatedByXY(monitor.GetGeometry().GetX(), monitor.GetGeometry().GetY())

	if state == Hovering {
		startPoint = absoluteMousePos
		state = SelectingRegion
		selectedRegionRect = NewRectangleFromXYWH(0, 0, 0, 0)
	}

	if state == QuickAnnotating {
		mousePosRelative := Point{
			X: absoluteMousePos.X - capturedRect.X(),
			Y: absoluteMousePos.Y - capturedRect.Y(),
		}

		annotations.HandleMousePressed(mousePosRelative)
	}

	updateCursor()
	redrawAllWindows()
}

func onMouseMove(monitor *gdk.Monitor, event *gdk.EventMotion) {
	mousePos = NewPointFromEventMotion(event)
	absoluteMousePos = mousePos.TranslatedByXY(monitor.GetGeometry().GetX(), monitor.GetGeometry().GetY())

	if state == Hovering {
		updateWindowUnderCursor()
	}

	if state == SelectingRegion {
		selectedRegionRect = NewRectangleFromPoints(startPoint, absoluteMousePos)
	}

	if state == QuickAnnotating {
		mousePosRelative := Point{
			X: absoluteMousePos.X - capturedRect.X(),
			Y: absoluteMousePos.Y - capturedRect.Y(),
		}

		if (event.State() & gdk.BUTTON1_MASK) > 0 {
			annotations.HandleMouseDrag(mousePosRelative)
		}
	}

	updateCursor()
	redrawAllWindows()
}

func onMousePrimaryReleased(event *gdk.EventButton) {
	setWindowCursor("default")

	if state == SelectingRegion {
		var controlPressed = (event.State() & uint(gdk.CONTROL_MASK)) > 0
		var shiftPressed = (event.State() & uint(gdk.SHIFT_MASK)) > 0

		if startPoint.ManhattanDistanceTo(absoluteMousePos) < 5 {
			if frozenScreen == nil {
				for _, window := range windows {
					window.Hide()
				}
				if hoveredWindow != nil {
					hoveredWindow.RaiseToFront()
				}

				_ = glib.TimeoutAdd(200, func() {
					captureScreen(hoveredWindowRect, controlPressed, shiftPressed)
				})
			} else {
				captureScreen(hoveredWindowRect, controlPressed, shiftPressed)
			}
		} else {
			captureScreen(selectedRegionRect, controlPressed, shiftPressed)
		}
	}

	if state == QuickAnnotating {
		annotations.HandleMouseReleased()
	}

	redrawAllWindows()
}

func onMouseSecondaryReleased() {
	setWindowCursor("default")

	if state == QuickAnnotating {
		if annotations.Empty() {
			state = Hovering
		} else {
			annotations.HandleMouseSecondaryReleased()
		}
	} else if state == SelectingRegion || state == Hovering {
		quitApp()
		return
	}

	redrawAllWindows()
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
		quitApp()
		return
	}

	updateCursor()
	redrawAllWindows()
}

func createWindow(monitor *gdk.Monitor) {
	var err error

	position := monitor.GetGeometry()

	window, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}

	window.Fullscreen()
	window.SetKeepAbove(true)
	window.SetAppPaintable(true)
	window.SetDecorated(false)
	window.SetSkipTaskbarHint(true)
	window.Resize(position.GetWidth(), position.GetHeight())
	window.Move(position.GetX(), position.GetY())

	_ = window.Connect("destroy", func() {
		quitApp()
	})

	_ = window.Connect("draw", func(eventWindow *gtk.Window, context *cairo.Context) {
		onDraw(monitor, context)
	})

	_ = window.Connect("button-press-event", func(eventWindow *gtk.Window, event *gdk.Event) bool {
		mouseEvent := gdk.EventButtonNewFromEvent(event)
		if mouseEvent.Button() == gdk.BUTTON_PRIMARY {
			onMousePrimaryPressed(monitor, mouseEvent)
		}
		return true
	})

	_ = window.Connect("motion-notify-event", func(eventWindow *gtk.Window, event *gdk.Event) bool {
		onMouseMove(monitor, gdk.EventMotionNewFromEvent(event))
		return true
	})

	_ = window.Connect("button-release-event", func(eventWindow *gtk.Window, event *gdk.Event) bool {
		mouseEvent := gdk.EventButtonNewFromEvent(event)
		if mouseEvent.Button() == gdk.BUTTON_PRIMARY {
			onMousePrimaryReleased(mouseEvent)
		}
		if mouseEvent.Button() == gdk.BUTTON_SECONDARY {
			onMouseSecondaryReleased()
		}
		return true
	})

	_ = window.Connect("key-release-event", func(eventWindow *gtk.Window, event *gdk.Event) bool {
		onKeyReleased(gdk.EventKeyNewFromEvent(event))
		return true
	})

	window.SetEvents(int(gdk.POINTER_MOTION_MASK | gdk.KEY_RELEASE_MASK | gdk.BUTTON_PRESS_MASK))

	// Allow window to be transparent
	visual, err := window.GetScreen().GetRGBAVisual()
	if err != nil || visual == nil {
		log.Fatal("Alpha not supported")
	}
	window.SetVisual(visual)

	window.Show()
	window.Present()

	windows = append(windows, window)
}

func main() {
	var err error

	gtk.Init(nil)

	desktopRect = getRootWindowRect()
	toplevelWindows = getCurrentToplevelWindows()

	if len(os.Args) >= 2 && os.Args[1] == "--freeze" {
		frozenScreen, err = captureScreenshot(desktopRect)
		_ = err
	}

	for _, monitor := range getMonitors() {
		createWindow(monitor)
	}

	mousePos = getMousePosition()
	updateWindowUnderCursor()
	updateCursor()
	redrawAllWindows()

	gtk.Main()
}

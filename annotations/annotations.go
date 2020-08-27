package annotations

import (
	"github.com/gotk3/gotk3/cairo"

	. "../utils"
)

var surface *cairo.Surface
var context *cairo.Context
var paths [][]Point
var newPathPoints []Point

func InitAnnotations(size Point) {
	surface = cairo.CreateImageSurface(cairo.FORMAT_ARGB32, size.X, size.Y)
	context = cairo.Create(surface)
}

func Has() bool {
	return surface != nil
}

func HandleMousePressed(point Point) {
	savePathPoint(point)
}

func HandleMouseDrag(point Point) {
	savePathPoint(point)
}

func HandleMouseReleased() {
	if len(newPathPoints) > 0 {
		paths = append(paths, newPathPoints)
		newPathPoints = nil
	}
}

func drawPath(points []Point) {
	for i, point := range points {
		if i == 0 {
			context.MoveTo(float64(point.X), float64(point.Y))
		} else {
			context.LineTo(float64(point.X), float64(point.Y))
		}
	}
}

func Draw(destContext *cairo.Context, x, y int) {
	context.SetAntialias(cairo.ANTIALIAS_SUBPIXEL)
	context.SetSourceRGB(1.0, 0.0, 0.0)
	context.SetLineWidth(2)

	for _, path := range paths {
		drawPath(path)
	}
	drawPath(newPathPoints)

	context.Stroke()

	destContext.SetSourceSurface(surface, float64(x), float64(y))
	destContext.Paint()
}

func savePathPoint(point Point) {
	newPathPoints = append(newPathPoints, point)
}

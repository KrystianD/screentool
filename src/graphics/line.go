package graphics

import (
	"math"

	"github.com/gotk3/gotk3/cairo"

	. "github.com/KrystianD/screentool/src/utils"
)

func DrawLine(context *cairo.Context, startPoint, endPoint Point, color Color3F, width float64) {
	context.SetAntialias(cairo.ANTIALIAS_SUBPIXEL)
	context.SetSourceRGB(color.Red, color.Green, color.Blue)
	context.SetLineWidth(width)

	startPoint.CairoMoveTo(context)
	endPoint.CairoLineTo(context)
	context.Stroke()
}

func SnapLineHV(startPoint, endPoint Point) (newStartPoint Point, newEndPoint Point) {
	var dX = math.Abs(float64(endPoint.X - startPoint.X))
	var dY = math.Abs(float64(endPoint.Y - startPoint.Y))

	if dX > dY {
		return NewPoint(startPoint.X, endPoint.Y), NewPoint(endPoint.X, endPoint.Y)
	} else {
		return NewPoint(endPoint.X, startPoint.Y), NewPoint(endPoint.X, endPoint.Y)
	}
}

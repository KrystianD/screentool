package graphics

import (
	"github.com/gotk3/gotk3/cairo"

	. "github.com/KrystianD/screentool/src/utils"
)

func DrawPath(context *cairo.Context, points []Point, color Color3F, width float64) {
	context.SetAntialias(cairo.ANTIALIAS_SUBPIXEL)
	context.SetSourceRGB(color.Red, color.Green, color.Blue)
	context.SetLineWidth(width)

	for i, point := range points {
		if i == 0 {
			context.MoveTo(float64(point.X), float64(point.Y))
		} else {
			context.LineTo(float64(point.X), float64(point.Y))
		}
	}
	context.Stroke()
}

package graphics

import (
	"github.com/gotk3/gotk3/cairo"

	. "github.com/KrystianD/screentool/src/utils"
)

func DrawPath(context *cairo.Context, points []Point) {
	for i, point := range points {
		if i == 0 {
			context.MoveTo(float64(point.X), float64(point.Y))
		} else {
			context.LineTo(float64(point.X), float64(point.Y))
		}
	}
	context.Stroke()
}

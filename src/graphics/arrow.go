package graphics

import (
	"math"

	"github.com/gotk3/gotk3/cairo"

	. "github.com/KrystianD/screentool/src/utils"
)

func DrawArrow(context *cairo.Context, startPoint, endPoint Point, color Color3F, width, arrowSize, arrowAngleDeg float64) {
	context.SetAntialias(cairo.ANTIALIAS_SUBPIXEL)
	context.SetSourceRGB(color.Red, color.Green, color.Blue)
	context.SetLineWidth(width)

	var distance = startPoint.DistanceTo(endPoint)

	var arrowLen = math.Min(distance, arrowSize)

	var normalized = NewPointFFromPoint(startPoint.TranslatedBy(endPoint.Negated())).Normalized()

	startPoint.CairoMoveTo(context)
	var lineLen = distance - 2
	if lineLen > 0 {
		startPoint.TranslatedBy(normalized.MultipliedBy(lineLen).Negated().ToPoint()).CairoLineTo(context)
	}
	context.Stroke()

	var a1 = normalized.RotatedDegree(arrowAngleDeg).MultipliedBy(arrowLen)
	var a2 = normalized.RotatedDegree(-arrowAngleDeg).MultipliedBy(arrowLen)

	context.NewPath()
	endPoint.CairoMoveTo(context)
	endPoint.TranslatedBy(a1.ToPoint()).CairoLineTo(context)
	endPoint.TranslatedBy(a2.ToPoint()).CairoLineTo(context)
	context.ClosePath()
	context.Fill()
}

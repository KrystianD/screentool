package graphics

import (
	"github.com/gotk3/gotk3/cairo"

	. "github.com/KrystianD/screentool/src/utils"
)

func DrawRectangle(context *cairo.Context, rect Rectangle, color Color3F, width float64, fill bool) {
	context.SetAntialias(cairo.ANTIALIAS_SUBPIXEL)
	context.SetSourceRGB(color.Red, color.Green, color.Blue)
	context.SetLineWidth(width)

	l, t, r, b := rect.GetLTRB()
	context.MoveTo(float64(l), float64(t))
	context.LineTo(float64(r), float64(t))
	context.LineTo(float64(r), float64(b))
	context.LineTo(float64(l), float64(b))

	if fill {
		context.Fill()
	} else {
		context.Stroke()
	}
}

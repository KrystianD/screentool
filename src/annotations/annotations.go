package annotations

import (
	"math"

	"github.com/gotk3/gotk3/cairo"

	. "github.com/KrystianD/screentool/src/utils"
)

type Freehand struct {
	Points []Point
}
type Arrow struct {
	Start, End Point
}

var surface *cairo.Surface
var context *cairo.Context
var isDrawing = false

var startPoint Point
var endPoint Point
var newPathPoints []Point

var objects []interface{}

var tool = 1

func InitAnnotations(size Point) {
	surface = cairo.CreateImageSurface(cairo.FORMAT_ARGB32, size.X, size.Y)
	context = cairo.Create(surface)
}

func Empty() bool {
	return len(objects) == 0
}

func HandleMousePressed(point Point) {
	if tool == 0 {
		isDrawing = true
		savePathPoint(point)
	} else if tool == 1 {
		isDrawing = true
		startPoint = point
		endPoint = point
	}
}

func HandleMouseDrag(point Point) {
	if !isDrawing {
		return
	}

	if tool == 0 {
		savePathPoint(point)
	} else if tool == 1 {
		endPoint = point
	}
}

func HandleMouseReleased() {
	finalizeCurrentDrawing()
}

func HandleMouseSecondaryReleased() {
	if len(objects) > 0 {
		objects = objects[:len(objects)-1]
	}
}

func CycleTool() {
	finalizeCurrentDrawing()

	tool += 1

	if tool > 1 {
		tool = 0
	}
}

func finalizeCurrentDrawing() {
	if !isDrawing {
		return
	}

	isDrawing = false

	if tool == 0 {
		objects = append(objects, Freehand{Points: newPathPoints})
		newPathPoints = nil
	} else if tool == 1 {
		objects = append(objects, Arrow{
			Start: startPoint,
			End:   endPoint,
		})
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
	context.Stroke()
}

func drawArrow(startPoint, endPoint Point) {
	const ArrowSize = 20
	const ArrowAngle = 25

	var distance = startPoint.DistanceTo(endPoint)

	var arrowLen = math.Min(distance, ArrowSize)

	var normalized = NewPointFFromPoint(startPoint.TranslatedBy(endPoint.Negated())).Normalized()

	startPoint.CairoMoveTo(context)
	var lineLen = distance - 2
	if lineLen > 0 {
		startPoint.TranslatedBy(normalized.MultipliedBy(lineLen).Negated().ToPoint()).CairoLineTo(context)
	}
	context.Stroke()

	var a1 = normalized.RotatedDegree(ArrowAngle).MultipliedBy(arrowLen)
	var a2 = normalized.RotatedDegree(-ArrowAngle).MultipliedBy(arrowLen)

	context.NewPath()
	endPoint.CairoMoveTo(context)
	endPoint.TranslatedBy(a1.ToPoint()).CairoLineTo(context)
	endPoint.TranslatedBy(a2.ToPoint()).CairoLineTo(context)
	context.ClosePath()
	context.Fill()
}

func Draw(destContext *cairo.Context, x, y int) {
	context.SetSourceRGBA(0.0, 1.0, 0.0, 0)
	context.SetOperator(cairo.OPERATOR_SOURCE)
	context.Paint()

	context.SetAntialias(cairo.ANTIALIAS_SUBPIXEL)
	context.SetSourceRGB(1.0, 0.0, 0.0)
	context.SetLineWidth(2)

	for _, object := range objects {
		if freehand, ok := object.(Freehand); ok {
			drawPath(freehand.Points)
		}
		if arrow, ok := object.(Arrow); ok {
			drawArrow(arrow.Start, arrow.End)
		}
	}

	if isDrawing {
		if tool == 0 {
			drawPath(newPathPoints)
		} else if tool == 1 {
			drawArrow(startPoint, endPoint)
		}
	}

	destContext.SetSourceSurface(surface, float64(x), float64(y))
	destContext.Paint()
}

func savePathPoint(point Point) {
	newPathPoints = append(newPathPoints, point)
}

package annotations

import (
	"github.com/gotk3/gotk3/cairo"

	"github.com/KrystianD/screentool/src/graphics"
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

func Draw(destContext *cairo.Context, x, y int) {
	var LineColor = NewColor3F(1., 0., 0.)
	const LineWidth = 2.
	const ArrowSize = 20
	const ArrowAngle = 25

	context.SetSourceRGBA(0.0, 0.0, 0.0, 0.0)
	context.SetOperator(cairo.OPERATOR_SOURCE)
	context.Paint()

	for _, object := range objects {
		if freehand, ok := object.(Freehand); ok {
			graphics.DrawPath(context, freehand.Points, LineColor, LineWidth)
		}
		if arrow, ok := object.(Arrow); ok {
			graphics.DrawArrow(context, arrow.Start, arrow.End, LineColor, LineWidth, ArrowSize, ArrowAngle)
		}
	}

	if isDrawing {
		if tool == 0 {
			graphics.DrawPath(context, newPathPoints, LineColor, LineWidth)
		} else if tool == 1 {
			graphics.DrawArrow(context, startPoint, endPoint, LineColor, LineWidth, ArrowSize, ArrowAngle)
		}
	}

	destContext.SetSourceSurface(surface, float64(x), float64(y))
	destContext.Paint()
}

func savePathPoint(point Point) {
	newPathPoints = append(newPathPoints, point)
}

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
type Eraser struct {
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
	tool = 1
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
	} else if tool == 2 {
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
	} else if tool == 2 {
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

	if tool > 2 {
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
	} else if tool == 2 {
		objects = append(objects, Eraser{
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
			graphics.DrawPath(context, freehand.Points, LineColor, LineWidth, false)
		}
		if arrow, ok := object.(Arrow); ok {
			graphics.DrawArrow(context, arrow.Start, arrow.End, LineColor, LineWidth, ArrowSize, ArrowAngle)
		}
		if eraser, ok := object.(Eraser); ok {
			graphics.DrawRectangle(context, NewRectangleFromPoints(eraser.Start, eraser.End), LineColor, LineWidth, true)
		}
	}

	if isDrawing {
		if tool == 0 {
			graphics.DrawPath(context, newPathPoints, LineColor, LineWidth, false)
		} else if tool == 1 {
			graphics.DrawArrow(context, startPoint, endPoint, LineColor, LineWidth, ArrowSize, ArrowAngle)
		} else if tool == 2 {
			graphics.DrawRectangle(context, NewRectangleFromPoints(startPoint, endPoint), LineColor, LineWidth, true)
		}
	}

	destContext.SetSourceSurface(surface, float64(x), float64(y))
	destContext.Paint()

	destContext.ResetClip()
	if tool == 0 {
		graphics.DrawPath(destContext, PathIcon.Scaled(0.5).Moved(x, y).Moved(-20/2, -20/2).Points, LineColor, 1.5, false)
	} else if tool == 1 {
		graphics.DrawArrow(destContext, NewPoint(x, y-20), NewPoint(x-20, y), LineColor, 1, ArrowSize/2, ArrowAngle)
	} else if tool == 2 {
		graphics.DrawRectangle(destContext, NewRectangleFromLTRB(x-10-7, y-10-5, x-10+7, y-10+5), LineColor, 1, true)
	}
}

func savePathPoint(point Point) {
	newPathPoints = append(newPathPoints, point)
}

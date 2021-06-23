package utils

import (
	"math"

	"github.com/gotk3/gotk3/cairo"
)

type Point struct {
	X, Y int
}

func NewPoint(x, y int) Point {
	return Point{
		X: x,
		Y: y,
	}
}

func (point Point) ManhattanDistanceTo(pt Point) int {
	return Abs(pt.X-point.X) + Abs(pt.Y-point.Y)
}

func (point Point) DistanceTo(pt Point) float64 {
	diffX := float64(pt.X - point.X)
	diffY := float64(pt.Y - point.Y)
	return math.Sqrt(diffX*diffX + diffY*diffY)
}

func (point Point) TranslatedByXY(x, y int) Point {
	return Point{
		X: point.X + x,
		Y: point.Y + y,
	}
}

func (point Point) TranslatedBy(pt Point) Point {
	return Point{
		X: point.X + pt.X,
		Y: point.Y + pt.Y,
	}
}

func (point Point) Negated() Point {
	return Point{
		X: -point.X,
		Y: -point.Y,
	}
}

func (point Point) Rotated(angle float64) Point {
	s := math.Sin(angle)
	c := math.Cos(angle)

	return Point{
		X: int(float64(point.X)*c - float64(point.Y)*s),
		Y: int(float64(point.X)*s + float64(point.Y)*c),
	}
}

func (point Point) CairoMoveTo(ctx *cairo.Context) {
	ctx.MoveTo(float64(point.X), float64(point.Y))
}

func (point Point) CairoLineTo(ctx *cairo.Context) {
	ctx.LineTo(float64(point.X), float64(point.Y))
}

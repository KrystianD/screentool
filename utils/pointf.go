package utils

import (
	"math"

	"github.com/gotk3/gotk3/cairo"
)

type PointF struct {
	X, Y float64
}

func NewPointF(x, y float64) PointF {
	return PointF{
		X: x,
		Y: y,
	}
}

func NewPointFFromPoint(point Point) PointF {
	return PointF{
		X: float64(point.X),
		Y: float64(point.Y),
	}
}

//func (point Point) TranslatedBy(pt Point) Point {
//	return Point{
//		X: point.X + pt.X,
//		Y: point.Y + pt.Y,
//	}

func (point PointF) Negated() PointF {
	return PointF{
		X: -point.X,
		Y: -point.Y,
	}
}

func (point PointF) MultipliedBy(value float64) PointF {
	return PointF{
		X: point.X * value,
		Y: point.Y * value,
	}
}

func (point PointF) Length() float64 {
	return math.Sqrt(point.X*point.X + point.Y*point.Y)
}

func (point PointF) Normalized() PointF {
	return point.MultipliedBy(1.0 / point.Length())
}

func (point PointF) Rotated(angle float64) PointF {
	s := math.Sin(angle)
	c := math.Cos(angle)

	return PointF{
		X: point.X*c - point.Y*s,
		Y: point.X*s + point.Y*c,
	}
}

func (point PointF) RotatedDegree(angle float64) PointF {
	return point.Rotated(angle / 180.0 * math.Pi)
}

func (point PointF) ToPoint() Point {
	return Point{
		X: int(point.X),
		Y: int(point.Y),
	}
}

func (point PointF) CairoMoveTo(ctx *cairo.Context) {
	ctx.MoveTo(point.X, point.Y)
}

func (point PointF) CairoLineTo(ctx *cairo.Context) {
	ctx.LineTo(point.X, point.Y)
}

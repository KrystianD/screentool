package utils

import (
	"fmt"

	"github.com/gotk3/gotk3/cairo"
)

type Rectangle struct {
	l, r int
	t, b int
}

func NewRectangleFromPoints(Point1, Point2 Point) Rectangle {
	return NewRectangleFromLTRB(Point1.X, Point1.Y, Point2.X, Point2.Y)
}

func NewRectangleFromXYWH(x, y, w, h int) Rectangle {
	return NewRectangleFromLTRB(x, y, x+w, y+h)
}

func NewRectangleFromLTRB(l, t, r, b int) Rectangle {
	return Rectangle{
		l: Min(l, r),
		r: Max(l, r),
		t: Min(t, b),
		b: Max(t, b),
	}
}

func (rect Rectangle) X() int {
	return rect.l
}

func (rect Rectangle) Y() int {
	return rect.t
}

func (rect Rectangle) Width() int {
	return rect.r - rect.l
}

func (rect Rectangle) Height() int {
	return rect.b - rect.t
}

func (rect Rectangle) Size() Point {
	return Point{
		X: rect.Width(),
		Y: rect.Height(),
	}
}

func (rect *Rectangle) MoveLTRB(leftOffset int, topOffset int, rightOffset int, bottomOffset int) {
	rect.l += leftOffset
	rect.t += topOffset
	rect.r += rightOffset
	rect.b += bottomOffset
}

func (rect Rectangle) GetLTRB() (int, int, int, int) {
	return rect.l, rect.t, rect.r, rect.b
}

func (rect Rectangle) GetXYWH() (int, int, int, int) {
	return rect.l, rect.t, rect.r - rect.l, rect.b - rect.t
}

func (rect Rectangle) Contains(point Point) bool {
	l, t, r, b := rect.GetLTRB()
	return l <= point.X && point.X <= r && t <= point.Y && point.Y <= b
}

func (rect Rectangle) TranslatedByXY(x, y int) Rectangle {
	return Rectangle{
		r: int(float32(rect.r)) + x,
		t: int(float32(rect.t)) + y,
		l: int(float32(rect.l)) + x,
		b: int(float32(rect.b)) + y,
	}
}

func (rect Rectangle) Scaled(scale float32) Rectangle {
	return Rectangle{
		r: int(float32(rect.r) * scale),
		t: int(float32(rect.t) * scale),
		l: int(float32(rect.l) * scale),
		b: int(float32(rect.b) * scale),
	}
}

func (rect Rectangle) Shrinked(value int) Rectangle {
	return Rectangle{
		r: rect.r - value,
		t: rect.t + value,
		l: rect.l + value,
		b: rect.b - value,
	}
}

func (rect Rectangle) IntersectionWith(otherRect Rectangle) (bool, Rectangle) {
	newRect := Rectangle{
		r: Min(rect.r, otherRect.r),
		t: Max(rect.t, otherRect.t),
		l: Max(rect.l, otherRect.l),
		b: Min(rect.b, otherRect.b),
	}
	return newRect.l <= newRect.r && newRect.t <= newRect.b, newRect
}

func (rect Rectangle) SetToCairo(ctx *cairo.Context) {
	x, y, w, h := rect.GetXYWH()
	ctx.Rectangle(float64(x), float64(y), float64(w), float64(h))
}

func (rect Rectangle) Format() string {
	return fmt.Sprintf("%d,%d(%dx%d)", rect.X(), rect.Y(), rect.Width(), rect.Height())
}

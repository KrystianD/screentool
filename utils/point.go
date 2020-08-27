package utils

type Point struct {
	X, Y int
}

func (point Point) ManhattanDistanceTo(pt Point) int {
	return Abs(pt.X-point.X) + Abs(pt.Y-point.Y)
}

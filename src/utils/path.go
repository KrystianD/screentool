package utils

type Path struct {
	Points []Point
}

func NewPath(points []Point) Path {
	return Path{
		Points: points,
	}
}

func (path Path) Centered() Path {
	if len(path.Points) == 0 {
		return NewPath([]Point{})
	}

	l := path.Points[0].X
	r := path.Points[0].X
	t := path.Points[0].Y
	b := path.Points[0].Y

	for _, point := range path.Points {
		l = Min(l, point.X)
		r = Max(r, point.X)
		t = Min(t, point.Y)
		b = Max(b, point.Y)
	}

	cx := (l + r) / 2
	cy := (t + b) / 2

	return path.Moved(-cx, -cy)
}

func (path Path) Scaled(s float32) Path {
	newPath := Path{
		Points: make([]Point, len(path.Points)),
	}

	for i, point := range path.Points {
		newPath.Points[i].X = int(float32(point.X) * s)
		newPath.Points[i].Y = int(float32(point.Y) * s)
	}

	return newPath
}

func (path Path) Moved(x, y int) Path {
	newPath := Path{
		Points: make([]Point, len(path.Points)),
	}

	for i, point := range path.Points {
		newPath.Points[i].X = point.X + x
		newPath.Points[i].Y = point.Y + y
	}

	return newPath
}

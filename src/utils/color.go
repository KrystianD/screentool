package utils

type Color3F struct {
	Red, Green, Blue float64
}

func NewColor3F(red, green, blue float64) Color3F {
	return Color3F{
		Red:   red,
		Green: green,
		Blue:  blue,
	}
}

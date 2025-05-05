package geometry

type Point struct {
	X int
	Y int
}

func (p Point) Add(p2 Point) Point {
	return Point{
		X: p.X + p2.X,
		Y: p.Y + p2.Y,
	}
}

func (p Point) Mod(width, height int) Point {
	return Point{
		X: (p.X + width) % width,
		Y: (p.Y + height) % height,
	}
}

func (p Point) Sub(p2 Point) Point {
	return Point{
		X: p.X - p2.X,
		Y: p.Y - p2.Y,
	}
}

package fractal

type Point8 struct {
	X, Y int
}

type Point64 struct {
	X, Y float64
}

type Rectangle8 struct {
	Min, Max Point8
}

func Rect8(x0, y0, x1, y1 int) Rectangle8 {
	return Rectangle8{Point8{x0, y0}, Point8{x1, y1}}
}

type Rectangle64 struct {
	Min, Max Point64
}

func Rect64(x0, y0, x1, y1 float64) Rectangle64 {
	return Rectangle64{Point64{x0, y0}, Point64{x1, y1}}
}

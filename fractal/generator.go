package fractal

import (
	"math/cmplx"
)

type Generator struct {
	Domain        Rectangle64
	Size          Rectangle8
	Bailout       float64
	MaxIterations uint32
	Function      func(complex128) func() complex128
}

func (gen Generator) At(X, Y int) uint32 {
	x0 := float64(X-gen.Size.Min.X)/
		float64(gen.Size.Max.X-gen.Size.Min.X)*
		(gen.Domain.Max.X-gen.Domain.Min.X) + gen.Domain.Min.X
	y0 := float64(Y-gen.Size.Min.Y)/
		float64(gen.Size.Max.Y-gen.Size.Min.Y)*
		(gen.Domain.Max.Y-gen.Domain.Min.Y) + gen.Domain.Min.Y
	c := complex(x0, y0)

	function_instance := gen.Function(c)
	var z complex128
	var itr uint32
	for ; cmplx.Abs(z) < gen.Bailout && itr < gen.MaxIterations; itr++ {
		z = function_instance()
	}
	return itr
}

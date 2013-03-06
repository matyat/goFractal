package fractal

import (
	"math/cmplx"
)

type Generator struct {
	Bailout       float64
	MaxIterations uint32
	Function      func(complex128) func() complex128
}

// Gets the number of iterations until escape for a complex value Z
func (gen Generator) EscapeAt(Z complex128) uint32 {
	function_instance := gen.Function(Z)
	var C complex128
	var itr uint32
	for ; cmplx.Abs(C) < gen.Bailout && itr < gen.MaxIterations; itr++ {
		C = function_instance()
	}
	return itr
}

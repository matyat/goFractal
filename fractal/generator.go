package fractal

import (
	"math/cmplx"
)

type Generator struct {
	Bailout       float64
	MaxIterations int
	Function      func(complex128) func() complex128
}

// Gets the number of iterations until escape for a complex value Z
func (gen Generator) EscapeAt(Z complex128) int {
	function_instance := gen.Function(Z)
	var C complex128
	var itr int
	for ; cmplx.Abs(C) < gen.Bailout && itr < gen.MaxIterations; itr++ {
		C = function_instance()
	}
	return itr
}

package fractal

import (
	"math"
	"math/cmplx"
)

// Function generator for factals, stores the fractal function,
// the iteration normalisation function, and the bailout and max
// iterations params.
type Generator struct {
	Bailout                float64
	MaxIterations          int
	Function               func(complex128) func() (complex128, float64)
	IterationNormalisation func(float64, complex128) float64
}

// Gets the number of iterations until escape for a complex value Z.
func (gen Generator) EscapeAt(C complex128) float64 {
	function_instance := gen.Function(C)
	var Z complex128
	var col float64
	var itr int

	// Loop until the Z becomes unbounded.
	// If the number of iteratation hits the MaxIterations var
	// we assume Z is bounded and return +infinity.
	for ; cmplx.Abs(Z) < gen.Bailout; itr++ {
		if itr == gen.MaxIterations {
			return math.Inf(1)
		}
		var c float64
		Z , c = function_instance()
		col += c
	}
	return gen.IterationNormalisation(col, Z)
}

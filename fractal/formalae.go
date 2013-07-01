package fractal

import(
	"math"
	"math/cmplx"
)

//Example fractal

//Mandelbrot Set
func Mandelbrot(bailout float64, max_iterations int) Generator {
	return Generator{
		Bailout:       bailout,
		MaxIterations: max_iterations,
		Function: func(c complex128) func() (complex128, float64) {
			C := c
			Z := complex(0, 0)
			return func() (complex128, float64) {
				Z = Z*Z + C
				return Z, 1.0
			}
		},
		IterationNormalisation: func(n float64, z complex128) float64{
			return n + 1 - (math.Log(math.Log(cmplx.Abs(z)))/math.Log(2))
		},
	}
}

//Julia set where f(z) = z^2 + c
func Julia(c complex128, bailout float64, max_iterations int) Generator {
	return Generator{
		Bailout:       bailout,
		MaxIterations: max_iterations,
		Function: func(z complex128) func() (complex128, float64) {
			Z := z
			C := c
			return func() (complex128, float64) {
				P := Z
				Z = Z*Z + C
				return Z, math.Exp(-cmplx.Abs(P))
			}
		},
		IterationNormalisation: func(n float64, z complex128) float64{
			return n + math.Exp(-cmplx.Abs(z))
		},

	}
}

// Newton fractal, where P is a function of z and Pd(z) = P'(z)
func Newton(P, Pd func(complex128) complex128, bailout float64, max_iterations int) Generator {
	return Generator{
		Bailout:       bailout,
		MaxIterations: max_iterations,
		Function: func(z complex128) func() (complex128, float64) {
			Z := z
			return func() (complex128, float64) {
				Z = Z - P(Z)/Pd(Z)
				return 1 / P(Z), 1
			}
		},
		IterationNormalisation: func(n float64, z complex128) float64{
			return n
		},
	}
}

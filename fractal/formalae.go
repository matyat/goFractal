package fractal

import(
	"math"
	"math/cmplx"
)

//Example fractal

//Mandelbrot Set
func Mandelbrot(max_iterations int) Generator {
	return Generator{
		Bailout:       4,
		MaxIterations: max_iterations,
		Function: func(c complex128) func() complex128 {
			C := c
			Z := complex(0, 0)
			return func() complex128 {
				Z = Z*Z + C
				return Z
			}
		},
		IterationNormalisation: func(n int, z complex128) float64{
			return float64(n+1) - (math.Log(math.Log(cmplx.Abs(z)))/math.Log(2))
		},
	}
}

//Julia set where f(z) = z^2 + c
func Julia(c complex128, max_iterations int) Generator {
	return Generator{
		Bailout:       2,
		MaxIterations: max_iterations,
		Function: func(z complex128) func() complex128 {
			Z := z
			C := c
			return func() complex128 {
				Z = Z*Z + C
				return Z
			}
		},
		IterationNormalisation: func(n int, z complex128) float64{
			return float64(n)
		},

	}
}

// Newton fractal, where P is a function of z and Pd(z) = P'(z)
func Newton(P, Pd func(complex128) complex128, max_iterations int) Generator {
	return Generator{
		Bailout:       1e14,
		MaxIterations: max_iterations,
		Function: func(z complex128) func() complex128 {
			Z := z
			return func() complex128 {
				Z = Z - P(Z)/Pd(Z)
				return 1 / P(Z)
			}
		},
		IterationNormalisation: func(n int, z complex128) float64{
			return float64(n)
		},
	}
}

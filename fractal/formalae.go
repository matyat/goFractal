package fractal
//Example fractal

//Mandelbrot Set
func Mandelbrot(max_iterations int) Generator {
	return Generator{
		Bailout: 2,
		MaxIterations: max_iterations,
		Function: func(c complex128) func() complex128 {
			C := c
			Z := complex(0, 0)
			return func() complex128 {
				Z = Z*Z + C
				return Z
			}
		},
	}
}

//Julia set where f(z) = z^2 + c
func JuliaQuad(c complex128, max_iterations int) Generator {
	return Generator{
		Bailout: 2,
		MaxIterations: max_iterations,
		Function: func(z complex128) func() complex128 {
			Z := z
			C := c
			return func() complex128 {
				Z = Z*Z + C
				return Z
			}
		},
	}
}

// Newton fractal, where P is a function of z and Pd(z) = P'(z)
func Newton(P, Pd func(complex128) complex128, max_iterations int) Generator {
	return Generator{
		Bailout: 1e14,
		MaxIterations: max_iterations,
		Function: func(z complex128) func() complex128 {
			Z := z
			return func() complex128 {
				Z = Z - P(Z)/Pd(Z)
				return 1/P(Z)
			}
		},
	}
}


package fractal

import (
	"math"
	"strings"
	"strconv"
)

// convert degree to radians
func degToRad(deg float64) float64 {
	return math.Pi * deg / 180
}

// Parse a string into a complex number type.
// The complex number can have any number of real and imaginary
// components in any order i.e. "0.5 - 2 + 12i - 3 + 0.1i" is 
// perfectly valid and will return -4.5 + 12.1i.
func parseCmplxString(str string) (complex128, error) {
	var R float64
	var I float64
	var val_ptr *float64

	value_strs := strings.Split(str, " ")

	sign := 1.0
	for i := range value_strs {
		// if the string is a '+' or '-', we change the sign
		// and go the the next string
		if value_strs[i] == "+" {
			sign = 1.0
			continue
		} else if value_strs[i] == "-" {
			sign = -1.0
			continue
		}

		// if the last char is 'i', this is an imaginary number, so trim off
		// the i and change the val_ptr to I. Else, we must be dealing with
		// a real number.
		if last := len(value_strs[i]) - 1; last >= 0 && value_strs[i][last] == 'i' {
			value_strs[i] = value_strs[i][:last]
			val_ptr = &I
		} else {
			val_ptr = &R
		}

		val, err := strconv.ParseFloat(value_strs[i], 64)
		if err != nil {
			return complex(0, 0), err
		}
		*val_ptr += sign * val
	}

	return complex(R, I), nil
}


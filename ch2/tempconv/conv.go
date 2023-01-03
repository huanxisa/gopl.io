// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 41.

//!+

//这个包名和文件名一致

package tempconv

import (
	"fmt"
	tempconv "gopl.io/ch2/tempconv0"
)

// CToF converts a Celsius temperature to Fahrenheit.、
//这个tempconv可以直接访问同级目录下的tempconv文件下的变量
func CToF(c Celsius) Fahrenheit {
	fmt.Printf("Brrrr! %v\n", tempconv.AbsoluteZeroC) // "Brrrr! -273.15°C"
	return Fahrenheit(c*9/5 + 32)
}

// FToC converts a Fahrenheit temperature to Celsius.
func FToC(f Fahrenheit) Celsius { return Celsius((f - 32) * 5 / 9) }

//!-

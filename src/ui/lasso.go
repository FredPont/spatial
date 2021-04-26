package ui

import (
	"image/color"
	"math"
)

func (r *interactiveRaster) drawline(x, y, x1, y1 int) {
	if x1 == x {
		if y1 < y {
			y, y1 = swap(y, y1)
		}
		for i := y; i < y1; i++ {
			j := x
			//fmt.Println("i=", i, "j=", j)
			r.edit.SetPixelColor(j, i, color.RGBA{255, 0, 0, 255}) // set pixel x,y to red
		}
		return
	}
	a := (float64(y1) - float64(y)) / (float64(x1) - float64(x))
	b := float64(y) - a*float64(x)
	//fmt.Println("x=", x, "y=", y, "x1=", x1, "y1=", y1, "a=", a, "b=", b)

	if x1 < x {
		x, x1 = swap(x, x1)
	}

	if abs(x1-x) > abs(y1-y) {

		for i := x; i < x1; i++ {
			j := int(math.Round(a*float64(i) + b))
			//fmt.Println("i=", i, "j=", j)
			r.edit.SetPixelColor(i, j, color.RGBA{255, 0, 0, 255}) // set pixel x,y to red
		}
	} else {
		if y1 < y {
			y, y1 = swap(y, y1)
		}

		for i := y; i < y1; i++ {
			j := int(math.Round((float64(i) - b) / a))
			//fmt.Println("i=", i, "j=", j)
			r.edit.SetPixelColor(j, i, color.RGBA{255, 0, 0, 255}) // set pixel x,y to red
		}
	}

}

func abs(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}

}

func swap(a, b int) (int, int) {
	var x2 int = 0
	x2 = a
	a = b
	b = x2
	return a, b
}

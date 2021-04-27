package ui

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
)

func (r *interactiveRaster) drawline(x, y, x1, y1 int) []fyne.Position {
	var alldots []fyne.Position // store all line pixels
	if x1 == x {
		if y1 < y {
			y, y1 = swap(y, y1)
		}
		for i := y; i < y1; i++ {
			j := x
			//fmt.Println("i=", i, "j=", j)
			r.edit.SetPixelColor(j, i, color.RGBA{255, 0, 0, 255}) // set pixel x,y to red
			alldots = append(alldots, fyne.Position{float32(j), float32(i)})
		}
		return alldots
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
			alldots = append(alldots, fyne.Position{float32(i), float32(j)})
		}
	} else {
		if y1 < y {
			y, y1 = swap(y, y1)
		}

		for i := y; i < y1; i++ {
			j := int(math.Round((float64(i) - b) / a))
			//fmt.Println("i=", i, "j=", j)
			r.edit.SetPixelColor(j, i, color.RGBA{255, 0, 0, 255}) // set pixel x,y to red
			alldots = append(alldots, fyne.Position{float32(j), float32(i)})
		}
	}
	return alldots
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

func (r *interactiveRaster) clearPolygon(p [][]fyne.Position) {

	for _, fps := range p {
		for _, fp := range fps {
			r.edit.SetPixelColor(int(fp.X), int(fp.Y), color.RGBA{0, 255, 0, 255}) // set pixel x,y to transparent
		}
	}
}

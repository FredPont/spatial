package ui

import (
	"image/color"
	"lasso/src/filter"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// IRect makes a rectangle centered at x,y
func iRect(x, y, w, h int, color color.RGBA) *canvas.Rectangle {
	fx, fy, fw, fh := float32(x), float32(y), float32(w), float32(h)
	r := &canvas.Rectangle{FillColor: color}
	r.Move(fyne.Position{X: fx - (fw / 2), Y: fy - (fh / 2)})
	r.Resize(fyne.Size{Width: fw, Height: fh})
	return r
}

// iLine draws a line
func iLine(x1, y1, x2, y2 int, size float32, color color.RGBA) *canvas.Line {
	p1 := fyne.Position{X: float32(x1), Y: float32(y1)}
	p2 := fyne.Position{X: float32(x2), Y: float32(y2)}
	l := &canvas.Line{StrokeColor: color, StrokeWidth: size, Position1: p1, Position2: p2}
	return l
}

func (r *interactiveRaster) drawline(x, y, x1, y1 int) []filter.Point {
	var alldots []filter.Point // store all line pixels
	r.edit.gateContainer.AddObject(iLine(x, y, x1, y1, 1., color.RGBA{30, 144, 255, 255}))

	return alldots
}

// func (r *interactiveRaster) drawline2(x, y, x1, y1 int) []filter.Point {
// 	var alldots []filter.Point // store all line pixels
// 	if x1 == x {
// 		if y1 < y {
// 			y, y1 = swap(y, y1)
// 		}
// 		for i := y; i < y1; i++ {
// 			j := x
// 			//fmt.Println("i=", i, "j=", j)
// 			r.edit.SetPixelColor(j, i, color.RGBA{255, 0, 0, 255}) // set pixel x,y to red
// 			alldots = append(alldots, filter.Point{j, i})
// 		}
// 		return alldots
// 	}
// 	a := (float64(y1) - float64(y)) / (float64(x1) - float64(x))
// 	b := float64(y) - a*float64(x)
// 	//fmt.Println("x=", x, "y=", y, "x1=", x1, "y1=", y1, "a=", a, "b=", b)

// 	if x1 < x {
// 		x, x1 = swap(x, x1)
// 	}

// 	if abs(x1-x) > abs(y1-y) {

// 		for i := x; i < x1; i++ {
// 			j := int(math.Round(a*float64(i) + b))
// 			//fmt.Println("i=", i, "j=", j)
// 			r.edit.SetPixelColor(i, j, color.RGBA{255, 0, 0, 255}) // set pixel x,y to red
// 			alldots = append(alldots, filter.Point{i, j})
// 		}
// 	} else {
// 		if y1 < y {
// 			y, y1 = swap(y, y1)
// 		}

// 		for i := y; i < y1; i++ {
// 			j := int(math.Round((float64(i) - b) / a))
// 			//fmt.Println("i=", i, "j=", j)
// 			r.edit.SetPixelColor(j, i, color.RGBA{255, 0, 0, 255}) // set pixel x,y to red
// 			alldots = append(alldots, filter.Point{j, i})
// 		}
// 	}
// 	return alldots
// }

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

func (r *interactiveRaster) clearPolygon(p [][]filter.Point) {

	// for _, fps := range p {
	// 	for _, fp := range fps {
	// 		//r.edit.SetPixelColor(int(fp.X), int(fp.Y), color.RGBA{0, 0, 0, 0}) // set pixel x,y to transparent
	// 	}
	// }
}

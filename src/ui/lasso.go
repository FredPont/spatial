package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// credits : https://github.com/ajstarks/fc
// IRect makes a rectangle centered at x,y
func iRect(x, y, w, h int, color color.RGBA) *canvas.Rectangle {
	fx, fy, fw, fh := float32(x), float32(y), float32(w), float32(h)
	r := &canvas.Rectangle{FillColor: color}
	r.Move(fyne.Position{X: fx - (fw / 2), Y: fy - (fh / 2)})
	r.Resize(fyne.Size{Width: fw, Height: fh})
	return r
}

// credits : https://github.com/ajstarks/fc
// iLine draws a line between 2 points
func iLine(x1, y1, x2, y2 int, size float32, color color.RGBA) *canvas.Line {
	p1 := fyne.Position{X: float32(x1), Y: float32(y1)}
	p2 := fyne.Position{X: float32(x2), Y: float32(y2)}
	l := &canvas.Line{StrokeColor: color, StrokeWidth: size, Position1: p1, Position2: p2}
	return l
}

// drawline draws a line between 2 points to the gate container
func (r *interactiveRaster) drawline(x, y, x1, y1 int) fyne.CanvasObject {
	l := iLine(x, y, x1, y1, 1., color.RGBA{30, 144, 255, 255}) // line between 2 points
	r.edit.gateContainer.AddObject(l)                           // add the line to the gate container
	return l
}

// clear the lines in the gate container
func (r *interactiveRaster) clearPolygon(gatesLines []fyne.CanvasObject) {
	for _, gl := range gatesLines {
		r.edit.gateContainer.Remove(gl)
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

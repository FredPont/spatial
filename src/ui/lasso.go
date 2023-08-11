// some functions here are inspired from : https://github.com/ajstarks/fc
// with the folowing licence :
/*
Copyright (C) 2018 Fyne.io developers (see AUTHORS)
All rights reserved.


Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
    * Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.
    * Neither the name of Fyne.io nor the names of its contributors may be
      used to endorse or promote products derived from this software without
      specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

*/

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

// borderRect makes a rectangle centered at x,y without filling
func borderRect(x, y, w, h int, rcolor color.NRGBA) *canvas.Rectangle {
	fx, fy, fw, fh := float32(x), float32(y), float32(w), float32(h)
	bgcolor := color.NRGBA{uint8(250), uint8(250), uint8(250), 0} // transparent filling
	r := &canvas.Rectangle{FillColor: bgcolor, StrokeColor: rcolor, StrokeWidth: 1}
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
	l := iLine(x, y, x1, y1, 1., color.RGBA{212, 170, 0, 255}) // line between 2 points
	r.edit.gateContainer.Add(l)                                // add the line to the gate container
	return l
}

// drawline a circle at x,y position to the gate dot container
func (e *Editor) drawcircleGateDot(x, y, ray int, color color.NRGBA) fyne.CanvasObject {
	c := iCircle(x, y, ray, color) // draw circle rayon ray
	e.gateDotsContainer.Add(c)     // add the cicle to the cluster container
	return c
}

// drawcircleGateCont draw a circle at x,y position to the gate container
func (e *Editor) drawcircleGateCont(x, y, ray int, color color.NRGBA) fyne.CanvasObject {
	c := iCircle(x, y, ray, color) // draw circle rayon ray
	e.gateContainer.Add(c)         // add the cicle to the cluster container
	return c
}

// clear the lines of the last gate in the gate container
func (r *interactiveRaster) clearPolygon(gatesLines []fyne.CanvasObject) {
	//fmt.Println("gate lines", gatesLines)
	for _, gl := range gatesLines {
		r.edit.gateContainer.Remove(gl)
	}

}

// draw the gate number after double click
func (r *interactiveRaster) drawGateNb(x, y int, gateNB string) {
	offset := 20 // x,y offset from 1st dot of the gate
	//gateNB := strconv.Itoa(len(r.alledges) - 1)
	//gateNB := strconv.Itoa(r.gatesNumbers.nb)
	AbsText(r.edit.gateNumberContainer, x-offset, y+offset, gateNB, 20, color.NRGBA{255, 255, 0, 255})
}

// func abs(x int) int {
// 	if x < 0 {
// 		return -x
// 	}
// 	return x

// }

// func swap(a, b int) (int, int) {
// 	var x2 int = 0
// 	x2 = a
// 	a = b
// 	b = x2
// 	return a, b
// }

package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
)

// drawcircleScattCont draw a circle at x,y position to the Gate container
func (r *plotRaster) drawcircleGateCont(x, y, ray int, color color.NRGBA) fyne.CanvasObject {
	c := iCircle(x, y, ray, color)    // draw circle rayon ray
	r.plot2DEdit.gateContainer.Add(c) // add the cicle to the cluster container
	return c
}

// draw the gate number after double click
func (r *plotRaster) plotGateNb(x, y int, gateNB string) {
	offset := 20 // x,y offset from 1st dot of the gate
	//gateNB := strconv.Itoa(len(r.alledges) - 1)
	//gateNB := strconv.Itoa(r.gatesNumbers.nb)
	AbsText(r.plot2DEdit.gateContainer, x-offset, y+offset, gateNB, 20, color.NRGBA{255, 0, 0, 255})
}

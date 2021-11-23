package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
)

// drawcircleScattCont draw a circle at x,y position to the scatter container
func (r *Interactive2Dsurf) drawcircleScattCont(x, y, ray int, color color.NRGBA) fyne.CanvasObject {
	c := iCircle(x, y, ray, color)  // draw circle rayon ray
	r.scatterContainer.AddObject(c) // add the cicle to the cluster container
	return c
}

// drawcircleScattCont draw a circle at x,y position to the Gate container
func (r *plotRaster) drawcircleGateCont(x, y, ray int, color color.NRGBA) fyne.CanvasObject {
	c := iCircle(x, y, ray, color)          // draw circle rayon ray
	r.plot2DEdit.gateContainer.AddObject(c) // add the cicle to the cluster container
	return c
}

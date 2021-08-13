package ui

import (
	"fmt"
	"image/color"

	//"lasso/src/filter"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

////////////////////////////
// interactive raster code
////////////////////////////

type vulcRaster struct {
	widget.BaseWidget
	edit *Vulcano
}

func (r *vulcRaster) MinSize() fyne.Size {
	//fmt.Println("min size :", r.edit.min)
	return r.edit.min
}

func (r *vulcRaster) CreateRenderer() fyne.WidgetRenderer {
	return &vulcWidgetRender{raster: r, bg: canvas.NewRasterWithPixels(vbgPattern)}
}

// this function draw a selection rectangle around dots
func (r *vulcRaster) Tapped(ev *fyne.PointEvent) {
	r.edit.selectContainer.Objects = nil // clear previous selection

	x := int(ev.Position.X)
	y := int(ev.Position.Y)
	w, h := 20, 20 // selection rectangle size
	R := uint8(250)
	G := uint8(50)
	B := uint8(50)

	rect := borderRect(x, y, w, h, color.NRGBA{R, G, B, 255})
	r.edit.selectContainer.AddObject(rect)

	fmt.Println(x, y)

	r.edit.selectContainer.Refresh() // refresh only the gate container, faster than refresh layer
}

func (r *vulcRaster) TappedSecondary(*fyne.PointEvent) {
	// var line fyne.CanvasObject // store all line objects
	// lp := len(r.points)
	// if lp >= 1 {
	// 	x, y := r.points[lp-1].X, r.points[lp-1].Y
	// 	x2, y2 := r.points[0].X, r.points[0].Y // get first coordinates stored
	// 	line = r.drawline(x2, y2, x, y)
	// 	fmt.Println(r.points)
	// 	r.edit.layer.Refresh()
	// }
	// // avoid to add a void polygon :
	// if len(r.points) > 2 {
	// 	r.alledges = append(r.alledges, r.points) // store new edges
	// }
	// r.points = nil                        // reset polygone coordinates
	// r.tmpLines = append(r.tmpLines, line) // store new line object
	// r.gatesLines = r.tmpLines
	// r.tmpLines = nil // initialisation of gate lines

}

// func (r *vulcRaster) locationForPosition(pos fyne.Position) (int, int) {
// 	c := fyne.CurrentApp().Driver().CanvasForObject(r.img)
// 	x, y := int(pos.X), int(pos.Y)
// 	if c != nil {
// 		x, y = c.PixelCoordinateForPosition(pos)
// 	}
// 	return x, y
// }

func newVulcRaster(edit *Vulcano) *vulcRaster {
	r := &vulcRaster{edit: edit}

	r.ExtendBaseWidget(r)
	return r
}

type vulcWidgetRender struct {
	raster *vulcRaster
	bg     *canvas.Raster
}

func vbgPattern(x, y, _, _ int) color.Color {
	//const boxSize = 25

	// if (x/boxSize)%2 == (y/boxSize)%2 {
	// 	return color.Gray{Y: 58}
	// }

	return color.Gray{Y: 84}
}

func (r *vulcWidgetRender) Layout(size fyne.Size) {
	r.bg.Resize(size)

}

func (r *vulcWidgetRender) MinSize() fyne.Size {
	return r.MinSize()
}

func (r *vulcWidgetRender) Refresh() {
	canvas.Refresh(r.raster)
}

func (r *vulcWidgetRender) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *vulcWidgetRender) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg}
}

func (r *vulcWidgetRender) Destroy() {
}

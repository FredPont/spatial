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
	// points     []filter.Point      // points of current polygone edges
	// alledges   [][]filter.Point    // points of all current polygones edges
	// tmpLines   []fyne.CanvasObject // temporary slice with lines of the last gate
	// gatesLines []fyne.CanvasObject // line (fyne canvas object) of last polygone
}

func (r *vulcRaster) MinSize() fyne.Size {
	//fmt.Println("min size :", r.edit.min)
	return r.edit.min
}

func (r *vulcRaster) CreateRenderer() fyne.WidgetRenderer {
	return &vulcWidgetRender{raster: r, bg: canvas.NewRasterWithPixels(vbgPattern)}
}

// this function draw the lasso and store the lasso coordinates in r.points
func (r *vulcRaster) Tapped(ev *fyne.PointEvent) {
	// var line fyne.CanvasObject // store all line pixels
	x := int(ev.Position.X)
	y := int(ev.Position.Y)
	// lp := len(r.points)
	// if lp >= 1 {
	// 	x2, y2 := r.points[lp-1].X, r.points[lp-1].Y // get last coordinates stored
	// 	line = r.drawline(x2, y2, x, y)              // draw a line between the new pixel cliked and the last one stored in r.points
	// }

	fmt.Println(x, y)
	// r.points = append(r.points, filter.Point{x, y}) // store new edges

	// r.tmpLines = append(r.tmpLines, line) // store new lines objects
	// //r.edit.layer.Refresh() // slow
	// r.edit.gateContainer.Refresh() // refresh only the gate container, faster than refresh layer
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

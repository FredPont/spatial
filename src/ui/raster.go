package ui

import (
	"fmt"
	"image/color"
	"lasso/src/filter"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

////////////////////////////
// interactive raster code
////////////////////////////

// Dot contains the x,y pixel coordinates of a dot
// type Dot struct {
// 	x int
// 	y int
// }

// Point = XY coordinates of a point
// type Point struct {
// 	X, Y int
// }

type interactiveRaster struct {
	widget.BaseWidget
	edit *editor
	//min    fyne.Size
	img       *canvas.Raster
	points    []filter.Point   // points of current polygone edges
	alledges  [][]filter.Point // points of all current polygones  edges
	allpoints [][]filter.Point // points of all current polygones including lines between edges
}

func (r *interactiveRaster) MinSize() fyne.Size {
	//fmt.Println("min size :", r.edit.min)
	return r.edit.min
}

func (r *interactiveRaster) CreateRenderer() fyne.WidgetRenderer {
	return &rasterWidgetRender{raster: r, bg: canvas.NewRasterWithPixels(bgPattern)}
}

// this function draw the lasso and store the lasso coordinates in r.points
func (r *interactiveRaster) Tapped(ev *fyne.PointEvent) {
	var alldots []filter.Point // store all line pixels
	x, y := r.locationForPosition(ev.Position)

	r.edit.SetPixelColor(x, y, color.RGBA{255, 0, 0, 255}) // set pixel x,y to red

	lp := len(r.points)
	if lp >= 1 {
		x2, y2 := r.points[lp-1].X, r.points[lp-1].Y // get last coordinates stored
		alldots = r.drawline(x2, y2, x, y)           // draw a line between the new pixel cliked and the last one stored in r.points
	}

	fmt.Println(x, y)
	r.points = append(r.points, filter.Point{x, y}) // store new edges

	r.allpoints = append(r.allpoints, r.points, alldots) // store new edges and lines pixels
	r.edit.layer.Refresh()
}

func (r *interactiveRaster) TappedSecondary(*fyne.PointEvent) {
	var alldots []filter.Point // store all line pixels
	lp := len(r.points)
	if lp >= 1 {
		x, y := r.points[lp-1].X, r.points[lp-1].Y
		x2, y2 := r.points[0].X, r.points[0].Y // get first coordinates stored
		alldots = r.drawline(x2, y2, x, y)
		fmt.Println(r.points)
		r.edit.layer.Refresh()
	}
	r.alledges = append(r.alledges, r.points)  // store new edges
	r.points = nil                             // reset polygone coordinates
	r.allpoints = append(r.allpoints, alldots) // store new edges and lines pixels
	//fmt.Println(r.allpoints)
	//r.clearPolygon(r.allpoints)

}

func (r *interactiveRaster) locationForPosition(pos fyne.Position) (int, int) {
	c := fyne.CurrentApp().Driver().CanvasForObject(r.img)
	x, y := int(pos.X), int(pos.Y)
	if c != nil {
		x, y = c.PixelCoordinateForPosition(pos)
	}
	return x, y
}

func newInteractiveRaster(edit *editor) *interactiveRaster {
	r := &interactiveRaster{img: canvas.NewRaster(edit.draw), edit: edit}

	r.ExtendBaseWidget(r)
	return r
}

type rasterWidgetRender struct {
	raster *interactiveRaster
	bg     *canvas.Raster
}

func bgPattern(x, y, _, _ int) color.Color {
	const boxSize = 25

	if (x/boxSize)%2 == (y/boxSize)%2 {
		return color.Gray{Y: 58}
	}

	return color.Gray{Y: 84}
}

func (r *rasterWidgetRender) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.raster.img.Resize(size)
}

func (r *rasterWidgetRender) MinSize() fyne.Size {
	return r.MinSize()
}

func (r *rasterWidgetRender) Refresh() {
	canvas.Refresh(r.raster)
}

func (r *rasterWidgetRender) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *rasterWidgetRender) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg, r.raster.img}
}

func (r *rasterWidgetRender) Destroy() {
}

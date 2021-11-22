/*
 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.

 Written by Frederic PONT.
 (c) Frederic Pont 2021
*/

package ui

import (
	"image/color"
	"spatial/src/filter"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

////////////////////////////
// interactive raster code
////////////////////////////

type plotRaster struct {
	widget.BaseWidget
	plot2DEdit *Interactive2Dsurf
	//mouseXY    filter.Point //position of the mouse click
	//selection  []PVrecord   // dots selected by user in vulcano plot
	//selItem    string       // item selected by the user to draw expression
	points     []filter.Point        // points coordinates of the last gate
	tmpLines   []fyne.CanvasObject   // circles of the last gate
	gatesLines [][]fyne.CanvasObject // all the gates dots
	alledges   [][]filter.Point      // points coordinates of all the gate
	plot2DBox  PlotBox
}

func (r *plotRaster) MinSize() fyne.Size {
	//fmt.Println("min size :", r.plot2DEdit.min)
	return r.plot2DEdit.min
}

func (r *plotRaster) CreateRenderer() fyne.WidgetRenderer {
	return &plotWidgetRender{raster: r, bg: canvas.NewRasterWithPixels(plotPattern)}
}

///////////////////////////////
// Dragged functions for brush
///////////////////////////////

func (r *plotRaster) Dragged(ev *fyne.DragEvent) {

	x := int(ev.Position.X)
	y := int(ev.Position.Y)
	r.points = append(r.points, filter.Point{X: x, Y: y}) // store new edges
	// draw a dot at the mouse position
	circle := r.drawcircleScattCont(x, y, 1, color.NRGBA{76, 0, 153, 255})

	//test
	//r.plot2DEdit.imageEditor.drawcircle(x, y, 1, color.NRGBA{76, 0, 153, 255})
	//r.plot2DEdit.imageEditor.clusterContainer.Refresh()

	r.tmpLines = append(r.tmpLines, circle) // store new circles objects in the r.tmpLines slice
	r.plot2DEdit.scatterContainer.Refresh()
	//log.Println(x, y, circle)
}

func (r *plotRaster) DragEnd() {
	r.alledges = append(r.alledges, r.points)       // store new edges
	r.points = nil                                  // reset polygone coordinates
	r.gatesLines = append(r.gatesLines, r.tmpLines) // store new circles objects in the r.gatesLines
	r.tmpLines = nil                                // initialisation of gate lines
	r.plot2DEdit.scatterContainer.Refresh()
}

func (r *plotRaster) Tapped(ev *fyne.PointEvent) {

}

func (r *plotRaster) TappedSecondary(*fyne.PointEvent) {

}

func newInteractive2DRaster(plotEdit *Interactive2Dsurf) *plotRaster {
	r := &plotRaster{plot2DEdit: plotEdit}

	r.ExtendBaseWidget(r)
	return r
}

type plotWidgetRender struct {
	raster *plotRaster
	bg     *canvas.Raster
}

func plotPattern(x, y, _, _ int) color.Color {
	//const boxSize = 25

	// if (x/boxSize)%2 == (y/boxSize)%2 {
	// 	return color.Gray{Y: 58}
	// }

	return color.Gray{Y: 84}
}

func (r *plotWidgetRender) Layout(size fyne.Size) {
	r.bg.Resize(size)

}

func (r *plotWidgetRender) MinSize() fyne.Size {
	return r.MinSize()
}

func (r *plotWidgetRender) Refresh() {
	canvas.Refresh(r.raster)
}

func (r *plotWidgetRender) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *plotWidgetRender) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg}
}

func (r *plotWidgetRender) Destroy() {
}

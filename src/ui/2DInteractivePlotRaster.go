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
	vulcEdit  *Interactive2Dsurf
	mouseXY   filter.Point //position of the mouse click
	selection []PVrecord   // dots selected by user in vulcano plot
	selItem   string       // item selected by the user to draw expression
	vulcBox   PlotBox
}

func (r *plotRaster) MinSize() fyne.Size {
	//fmt.Println("min size :", r.vulcedit.min)
	return r.vulcEdit.min
}

func (r *plotRaster) CreateRenderer() fyne.WidgetRenderer {
	return &plotWidgetRender{raster: r, bg: canvas.NewRasterWithPixels(plotPattern)}
}

// this function draw a selection rectangle around dots
func (r *plotRaster) Tapped(ev *fyne.PointEvent) {
	// r.vulcEdit.selectContainer.Objects = nil // clear previous selection

	// x := int(ev.Position.X)
	// y := int(ev.Position.Y)

	// // read the vulcano selection square size in pixels from preferences
	// pref := fyne.CurrentApp().Preferences()
	// vs := binding.BindPreferenceInt("vulcSelectSize", pref)
	// vsquare, err := vs.Get()
	// if err != nil {
	// 	log.Println("Error reading selection size value !", err)
	// }
	// if vsquare == 0 {
	// 	vsquare = 10
	// }
	// w, h := vsquare, vsquare // selection rectangle size
	// R := uint8(250)
	// G := uint8(50)
	// B := uint8(50)

	// rect := borderRect(x, y, w, h, color.NRGBA{R, G, B, 255})
	// r.vulcEdit.selectContainer.Add(rect)

	// //fmt.Println(x, y)
	// r.mouseXY = filter.Point{X: x, Y: y}

	// r.vulcEdit.selectContainer.Refresh() // refresh only the gate container, faster than refresh layer
}

func (r *plotRaster) TappedSecondary(*fyne.PointEvent) {
	//r.vulcEdit.vulcanoSelect(&r.vulcBox, r.mouseXY)

	//refreshVulanoTools(r.vulcEdit)

}

func newInteractive2DRaster(plotEdit *Interactive2Dsurf) *plotRaster {
	r := &plotRaster{vulcEdit: plotEdit}

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

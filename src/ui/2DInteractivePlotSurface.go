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

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// Interactive2Dsurf contains the drawsurface, containers for scatter and select brush dots
type Interactive2Dsurf struct {
	drawSurface      *plotRaster
	min              fyne.Size // size of the scatter container
	win              fyne.Window
	calc             *canvas.Image   // image calc to display the 2D plot dots
	calcSeldots      *canvas.Image   // image calc to display the 2D plot dots selected by the user
	layer            *fyne.Container // container with plot and interactive drawsurface
	gateContainer    *fyne.Container // container with the gate dot lines
	scatterContainer *fyne.Container // container with the scatter circles
	//tools            fyne.Window     // 2D plot tools windows
	//imageEditor      *Editor         // editor of the microscopie image is embeded to allow expression plots
	//header           []string        // header of the first data table (to allow expression plots)
	//tableName        string          // name of the first data table (to allow expression plots)

}

// NewInterative2D creates a new interactive 2D plot
func NewInterative2D() (*Interactive2Dsurf, int, int) {
	w, h := 800, 800
	calcImg := canvas.NewImageFromFile("temp/2Dplot/2Dplot.png")
	calcDots := canvas.NewImageFromFile("temp/2Dplot/dotsIngGates.png")                      // dots selected by the user
	sel := container.NewWithoutLayout(iRect(w/2, h/2, w, h, color.RGBA{0, 0, 0, 0}))         // select container
	sca := container.NewWithoutLayout(iRect(w/2, h/2, w, h, color.RGBA{255, 255, 255, 255})) // scatter container should be independant of select container for separate initialisaion
	plotEdit := &Interactive2Dsurf{calc: calcImg, calcSeldots: calcDots, min: fyne.Size{Width: float32(w), Height: float32(h)}, gateContainer: sel, scatterContainer: sca}
	plotEdit.drawSurface = newInteractive2DRaster(plotEdit)

	return plotEdit, w, h
}

// buildVulc creates the window of the 2D plot
func (p *Interactive2Dsurf) build2DinterPlot(w fyne.Window) {
	p.win = w
	//e.layer = container.NewMax(e.scatterContainer)
	p.layer = container.NewStack(p.drawSurface, p.scatterContainer, p.calc, p.calcSeldots, p.gateContainer)
	w.SetContent(p.layer)

}

// build2DplotWin creates display 2Dplot window
func build2DplotWin(imageEditor *Editor) (fyne.Window, *Interactive2Dsurf) {
	w := fyne.CurrentApp().NewWindow("2D Plot")
	p, finalWidth, finalHeight := NewInterative2D()
	p.build2DinterPlot(w)
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(float32(finalWidth), float32(finalHeight)))
	w.Show()

	//p.imageEditor = imageEditor // store the image Editor to enable expression display from the 2D plot
	return w, p
}

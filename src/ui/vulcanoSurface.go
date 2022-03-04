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
	//"fmt"

	"image/color"

	//"log"

	//"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/data/binding"
	//"fyne.io/fyne/v2/theme"
	//"fyne.io/fyne/v2/widget"
)

// Vulcano contains the drawsurface, containers for scatter and select lines
type Vulcano struct {
	drawSurface      *vulcRaster
	min              fyne.Size // size of the scatter container
	win              fyne.Window
	layer            *fyne.Container // container with plot and interactive drawsurface
	selectContainer  *fyne.Container // container with the select lines
	scatterContainer *fyne.Container // container with the scatter circles
	tools            fyne.Window     // vulcano tools windows
	imageEditor      *Editor         // editor of the microscopie image is embeded to allow expression plots
	header           []string        // header of the first data table (to allow expression plots)
	tableName        string          // name of the first data table (to allow expression plots)
}

// NewVulcano creates a new interactive vulcano plot
func NewVulcano() (*Vulcano, int, int) {
	w, h := 800, 800

	sel := container.NewWithoutLayout(iRect(w/2, h/2, w, h, color.RGBA{0, 0, 0, 0}))         // select container
	sca := container.NewWithoutLayout(iRect(w/2, h/2, w, h, color.RGBA{255, 255, 255, 255})) // scatter container should be independant of select container for separate initialisaion
	vulcEdit := &Vulcano{min: fyne.Size{Width: float32(w), Height: float32(h)}, selectContainer: sel, scatterContainer: sca}
	vulcEdit.drawSurface = newVulcRaster(vulcEdit)

	return vulcEdit, w, h
}

// buildVulc creates the window of the vulcano plot
func (v *Vulcano) buildVulc(w fyne.Window) {
	v.win = w
	//e.layer = container.NewMax(e.scatterContainer)
	v.layer = container.NewMax(v.drawSurface, v.scatterContainer, v.selectContainer)
	w.SetContent(v.layer)

}

// buildVulWin creates display vulcano window
func buildVulcWin(imageEditor *Editor) *Vulcano {
	w := fyne.CurrentApp().NewWindow("Vulcano Plot")
	v, finalWidth, finalHeight := NewVulcano()
	v.buildVulc(w)
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(float32(finalWidth), float32(finalHeight)))
	w.Show()
	v.imageEditor = imageEditor // store the image Editor to enable expression display from the vulcano plot
	return v
}

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
	"log"
	"spatial/src/filter"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

func (v *Vulcano) vulcanoSelect(vul *PlotBox, mouse filter.Point) {
	var selectedPoints []PVrecord

	// selection size
	pref := fyne.CurrentApp().Preferences()
	vs := binding.BindPreferenceInt("vulcSelectSize", pref)
	vsquare, err := vs.Get()
	if err != nil {
		log.Println("Error reading selection size value !", err)
	}
	if vsquare == 0 {
		vsquare = 10 // selection size by default
	}
	//log.Println("vul.X", vul.X)
	for i, Xscatter := range vul.X {
		x1 := xCoord(vul, Xscatter)
		y1 := yCoord(vul, vul.Y[i])

		//log.Println("dot filtration", mouse.X, mouse.Y, x1, y1, vsquare, inSquare(mouse.X, mouse.Y, x1, y1, vsquare))
		if inSquare(mouse.X, mouse.Y, x1, y1, vsquare) {
			selectedPoints = append(selectedPoints, PVrecord{item: vul.id[i], log2fc: Xscatter, log10pv: vul.Y[i]})
		}
	}
	v.drawSurface.selection = selectedPoints
	//log.Println("selected points", v.drawSurface.selection)
}

// inSquare check if x,y is inside a square centered in x1,y1
func inSquare(x, y, x1, y1, size int) bool {

	xmin := x1 - size/2
	xmax := x1 + size/2

	if x > xmax || x < xmin {
		return false
	}

	ymin := y1 - size/2
	ymax := y1 + size/2

	if y > ymax || y < ymin {
		return false
	}

	return true
}

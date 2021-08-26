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
	"fmt"
	"lasso/src/filter"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

func filterActiveGates(e *Editor, alledges [][]filter.Point, dataFiles []string, gateName string, pref fyne.Preferences, f binding.Float) {
	f.Set(0.3) // progress bar
	// get parameters from preferences
	param := prefToConf(pref)
	fmt.Println("start filtering...", alledges)
	// filter all data files with all active gates
	for gateNumber, polygon := range alledges {
		fmt.Println("polygon ", polygon)
		if len(polygon) < 3 {
			continue
		}
		for _, dataFile := range dataFiles {
			gateName = filter.FormatOutFile("filter", gateName, "") // test if name exist, if not, build a file name with the current time
			outFile := strconv.Itoa(gateNumber) + "_" + gateName + "_" + dataFile
			filter.FilterTable(e.zoom, dataFile, outFile, polygon, param)
		}

	}
	f.Set(0.) // progress bar
}

// PrefToConf retreive conf data from fyne pref
func prefToConf(pref fyne.Preferences) filter.Conf {
	// get // X coordinates
	xcor := binding.BindPreferenceString("xcor", pref) // set the link to preferences for rotation
	x, _ := xcor.Get()

	// get y coordinates
	ycor := binding.BindPreferenceString("ycor", pref) // set the link to preferences for rotation
	y, _ := ycor.Get()

	// get scaling factor
	sf := binding.BindPreferenceFloat("scaleFactor", pref) // set the link to preferences for scaling factor
	scale, _ := sf.Get()

	// get coordinates +90° rotation : necessary for 10x Genomics
	r := binding.BindPreferenceBool("rotate", pref) // set the link to preferences for rotation
	rotate, _ := r.Get()

	return filter.Conf{X: x, Y: y, Scale: scale, Rotate: rotate}

}

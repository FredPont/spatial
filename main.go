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

package main

import (
	"fmt"
	"log"
	"spatial/src/pogrebDB"
	"spatial/src/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
)

func main() {
	a := app.NewWithID("SpatialExplorer")
	title()

	// initPref initialise some user preferences when not set
	ui.InitPref()
	pogrebDB.InitPogreb()

	log.Println("preferences saved in :", a.Storage().RootURI())

	//InitTempDir remove old 2D plots and gates from temp/2Dplot dir
	ui.Init2DTempDir()

	w := a.NewWindow("image")
	w.SetMaster()
	e, imgW, imgH := ui.NewEditor()

	e.BuildUI(w)
	w.SetFixedSize(true) // fix win size

	// set the windows size to at least 500x500 and adjust the windows size
	// to the size of the microscopy image if the pref size are below.
	setImageWinSize(a, w, imgW, imgH)

	ui.BuildTools(a, w, e)

	w.ShowAndRun()

}

// set the windows size of the image. If the image is larger than the user
// preferences, the image is displayed with scroll bars.
// if the size in pref is < 500 at first start, the minimal size will be 500x500
func setImageWinSize(a fyne.App, w fyne.Window, imgW, imgH int) {

	pref := a.Preferences()

	// store the original image size
	pref.SetInt("imgW", imgW)
	pref.SetInt("imgH", imgH)

	// get width preference
	winW := binding.BindPreferenceFloat("winW", pref) // set the link to preferences for win Width
	wW, _ := winW.Get()
	// get height preference
	winH := binding.BindPreferenceFloat("winH", pref) // set the link to preferences for win Height
	wH, _ := winH.Get()

	finalWidth := setMinWindow(wW, imgW)
	finalHeight := setMinWindow(wH, imgH)
	//log.Println("finalSize", finalWidth, finalHeight)
	w.Resize(fyne.NewSize(finalWidth, finalHeight))
	//w.Resize(fyne.NewSize(1000, 1000))
}

// image size must be at least 500x500
// set the minimal windows size to minSize = 500 if the size in pref is < minSize
func setMinWindow(prefSize float64, imgSize int) float32 {
	const minSize = 500
	finalSize := float32(minSize)

	if float64(imgSize) < float64(minSize) { // if image too small
		log.Println("Caution ! image must be at least 500x500 ! results will be wrong !")
	} else if prefSize == 0 { // if pref not set return 500
		return 500.
	} else if float64(imgSize) > prefSize && prefSize >= float64(minSize) {
		finalSize = float32(prefSize)
	} else if float64(imgSize) < prefSize && prefSize >= float64(minSize) {
		finalSize = float32(imgSize)
	}

	return finalSize
}

func title() {

	fmt.Println("   ┌────────────────────────────────────────────────────┐") // unicode U+250C
	fmt.Println("   │ Single Cell Spatial Explorer (c)Frederic PONT 2021 │")
	fmt.Println("   │       Free Software GNU General Public License     │")
	fmt.Println("   └────────────────────────────────────────────────────┘")
}

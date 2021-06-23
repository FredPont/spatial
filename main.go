package main

import (
	"fmt"
	"lasso/src/ui"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
)

func main() {
	title()
	a := app.NewWithID("SpatialGate")

	w := a.NewWindow("image")
	e, imgW, imgH := ui.NewEditor()
	//e, imgW, imgH := ui.NewSmallEditor()
	e.BuildUI(w)
	// set the windows size to at least 500x500 and adjust the windows size
	// to the size of the microscopy image if the pref size are below.
	setImageWinSize(a, w, imgW, imgH)
	w.SetFixedSize(true) // fix win size
	w.Show()

	w2 := a.NewWindow("Tool Box")
	ui.BuildTools(a, w2, w, e)

	w2.Show()

	w2.ShowAndRun()
	w.ShowAndRun()

}

// set the windows size of the image. If the image is larger than the user
// preferences, the image is displayed with scroll bars.
// if the size in pref is < 500 at first start, the minimal size will be 500x500
func setImageWinSize(a fyne.App, w fyne.Window, imgW, imgH int) {
	finalWidth := float32(imgW)
	finalHeight := float32(imgH)

	pref := a.Preferences()
	// get width preference
	winW := binding.BindPreferenceFloat("winW", pref)
	wW, _ := winW.Get()
	// get height preference
	winH := binding.BindPreferenceFloat("winH", pref) // set the link to preferences for win width
	wH, _ := winH.Get()

	finalWidth = setMinWindow(wW, imgW)
	finalHeight = setMinWindow(wH, imgH)
	w.Resize(fyne.NewSize(finalWidth, finalHeight))
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

	fmt.Println("   ┌───────────────────────────────────────────────────┐") // unicode U+250C
	fmt.Println("   │   single cell Spatial Gate (c)Frederic PONT 2021  │")
	fmt.Println("   │       Free Software GNU General Public License    │")
	fmt.Println("   └───────────────────────────────────────────────────┘")
}

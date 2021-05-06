package main

import (
	"fmt"
	"lasso/src/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
)

func main() {
	title()
	a := app.NewWithID("SpatialGate")

	w := a.NewWindow("image")
	e, imgW, imgH := ui.NewEditor()
	e.BuildUI(w)
	setImageWinSize(a, w, imgW, imgH)
	w.Show()

	w2 := a.NewWindow("Tool Box")
	ui.BuildTools(a, w2, w, e)

	w2.Show()

	w2.ShowAndRun()
	w.ShowAndRun()

}

// set the windows size of the image. If the image is larger than the user
// preferences, the image is displayed with scroll bars.
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

	if float64(imgW) > wW {
		finalWidth = float32(wW)
	}
	if float64(imgH) > wH {
		finalHeight = float32(wH)
	}

	w.Resize(fyne.NewSize(finalWidth, finalHeight))
}

func title() {

	fmt.Println("   ┌───────────────────────────────────────────────────┐") // unicode U+250C
	fmt.Println("   │   single cell Spatial Gate (c)Frederic PONT 2021  │")
	fmt.Println("   │       Free Software GNU General Public License    │")
	fmt.Println("   └───────────────────────────────────────────────────┘")
}

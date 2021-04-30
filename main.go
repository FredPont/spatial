package main

import (
	"lasso/src/ui"

	//"fmt"
	//"image"
	//"image/color"

	//"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	//"fyne.io/fyne/v2/canvas"
	//"fyne.io/fyne/v2/theme"
)

func main() {
	a := app.NewWithID("SpatialGate")

	//pref.SetPref(a) // set user preferences for the application
	w := a.NewWindow("image")
	e := ui.NewEditor()
	e.BuildUI(w)
	w.Resize(fyne.NewSize(300, 300))
	w.Show()

	w2 := a.NewWindow("Tool Box")
	ui.BuildTools(a, w2, w, e)

	w2.Show()

	w2.ShowAndRun()
	w.ShowAndRun()

}

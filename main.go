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
	//"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/theme"
	//"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Draw")
	e := ui.NewEditor()
	e.BuildUI(w)
	w.Resize(fyne.NewSize(300, 300))
	w.ShowAndRun()

}

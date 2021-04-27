package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func BuildTools(w2 fyne.Window, e *editor) {
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter text...")

	content := container.NewVBox(input, widget.NewButton("Clear", func() {
		e.drawSurface.clearPolygon(e.drawSurface.allpoints)
		e.layer.Refresh()
	}))

	w2.SetContent(content)
}

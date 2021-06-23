package ui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Zoom objet to record the %zoom and change editor min size
type Zoom struct {
	edit *editor

	zoom *widget.Label
}

func (z *Zoom) updateZoom(val int) {
	log.Println("val=", val)
	if val < 1 {
		val = 10
	} else if val > 100 {
		val = 100
	}
	z.edit.setZoom(val)

	z.zoom.SetText(fmt.Sprintf("%d%%", z.edit.zoom))
}

func newZoom(edit *editor) fyne.CanvasObject {
	z := &Zoom{edit: edit, zoom: widget.NewLabel("100%")}
	zoom := container.NewHBox(
		widget.NewButtonWithIcon("", theme.ZoomOutIcon(), func() {
			z.updateZoom(z.edit.zoom - 10)
		}),
		z.zoom,
		widget.NewButtonWithIcon("", theme.ZoomInIcon(), func() {
			z.updateZoom(z.edit.zoom + 10)
		}))
	return zoom
}

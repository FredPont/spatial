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
}

// func (e *Editor) draw(w, h int) image.Image {
// 	return image.NewRGBA(image.Rect(0, 0, w, h))
// }

// NewVulcano creates a new interactive vulcano plot
func NewVulcano() (*Vulcano, int, int) {
	w, h := 800, 800

	sel := fyne.NewContainer(iRect(w/2, h/2, w, h, color.RGBA{0, 0, 0, 0})) // select container
	sca := fyne.NewContainer(iRect(w/2, h/2, w, h, color.RGBA{0, 0, 0, 0})) // scatter container should be independant of select container for separate initialisaion
	//fgCol := color.Transparent
	//edit := &editor{fg: fgCol, fgPreview: canvas.NewRectangle(fgCol), img: image.NewRGBA(image.Rect(0, 0, 600, 600)), microscop: micro}
	edit := &Vulcano{min: fyne.Size{Width: float32(w), Height: float32(h)}, selectContainer: sel, scatterContainer: sca}
	edit.drawSurface = newVulcRaster(edit)

	return edit, w, h
}

// buildVulc creates the window of the vulcano plot
func (e *Vulcano) buildVulc(w fyne.Window) {
	e.win = w
	//e.layer = container.NewMax(e.scatterContainer)
	e.layer = container.NewMax(e.drawSurface, e.scatterContainer, e.selectContainer)
	w.SetContent(container.NewScroll(e.layer))
}

// buildVulWin creates display vulcano window
func buildVulcWin() *Vulcano {
	w := fyne.CurrentApp().NewWindow("Vulcano Plot")
	v, finalWidth, finalHeight := NewVulcano()
	v.buildVulc(w)
	w.Resize(fyne.NewSize(float32(finalWidth), float32(finalHeight)))
	w.Show()
	return v
}

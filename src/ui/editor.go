package ui

import (
	//"fmt"
	"image"
	"image/color"

	//"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/theme"
	//"fyne.io/fyne/v2/widget"
)

type editor struct {
	drawSurface *interactiveRaster
	fgPreview   *canvas.Rectangle
	img         *image.RGBA
	microscop   *canvas.Image
	fg          color.Color
	layer       *fyne.Container
	win         fyne.Window
}

func (e *editor) draw(w, h int) image.Image {
	return image.NewRGBA(image.Rect(0, 0, w, h))
}

// NewEditor creates a new pixel editor that is ready to have a file loaded
func NewEditor() *editor {
	micro := canvas.NewImageFromFile("tissue_lowres_image.png")
	micro.FillMode = canvas.ImageFillOriginal
	fgCol := color.Transparent
	edit := &editor{fg: fgCol, fgPreview: canvas.NewRectangle(fgCol), img: image.NewRGBA(image.Rect(0, 0, 600, 600)), microscop: micro}
	edit.drawSurface = newInteractiveRaster(edit)

	return edit
}

// BuildUI creates the main window of our pixel edit application
func (e *editor) BuildUI(w fyne.Window) {
	e.win = w
	e.layer = container.NewMax(e.drawSurface, e.microscop, canvas.NewImageFromImage(e.img))

	w.SetContent(container.NewScroll(e.layer))
	//w.SetContent(container.NewMax(e.buildUI(), canvas.NewImageFromImage(e.img)))
}

func (e *editor) buildUI() fyne.CanvasObject {
	return container.NewScroll(e.drawSurface)
}

func (e *editor) SetPixelColor(x, y int, c color.RGBA) {
	e.img.SetRGBA(x, y, c)
}

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
	img         *image.RGBA // image with polygons
	microscop   *canvas.Image
	min         fyne.Size // size of the microscop image
	layer       *fyne.Container
	win         fyne.Window
}

func (e *editor) draw(w, h int) image.Image {
	return image.NewRGBA(image.Rect(0, 0, w, h))
}

// NewEditor creates a new pixel editor that is ready to have a file loaded
func NewEditor() *editor {
	imgFile, w, h := ImgSize()
	micro := canvas.NewImageFromFile(imgFile)
	micro.FillMode = canvas.ImageFillOriginal
	//fgCol := color.Transparent
	//edit := &editor{fg: fgCol, fgPreview: canvas.NewRectangle(fgCol), img: image.NewRGBA(image.Rect(0, 0, 600, 600)), microscop: micro}
	edit := &editor{img: image.NewRGBA(image.Rect(0, 0, w, h)), microscop: micro, min: fyne.Size{float32(w), float32(h)}}
	edit.drawSurface = newInteractiveRaster(edit)

	return edit
}

// BuildUI creates the main window of our application
func (e *editor) BuildUI(w fyne.Window) {
	e.win = w
	e.layer = container.NewMax(e.drawSurface, e.microscop, canvas.NewImageFromImage(e.img))
	w.SetContent(container.NewScroll(e.layer))
}

// func (e *editor) buildUI() fyne.CanvasObject {
// 	return container.NewScroll(e.drawSurface)
// }

func (e *editor) SetPixelColor(x, y int, c color.RGBA) {
	e.img.SetRGBA(x, y, c)
}

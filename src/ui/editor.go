package ui

import (
	//"fmt"
	"image"
	"image/color"
	"log"

	//"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/theme"
	//"fyne.io/fyne/v2/widget"
)

type editor struct {
	drawSurface *interactiveRaster
	//img           *image.RGBA // image with polygons
	microscop               *canvas.Image
	min                     fyne.Size       // size of the microscop image and the gate/clusters containers
	layer                   *fyne.Container // container with image and interactive drawsurface
	win                     fyne.Window
	gateContainer           *fyne.Container // container with the gates lines
	clusterContainer        *fyne.Container // container with the cluster circles
	zoom                    int             // image zoom
	cacheWidth, cacheHeight int             // size of the window not zoomed
}

func (e *editor) draw(w, h int) image.Image {
	return image.NewRGBA(image.Rect(0, 0, w, h))
}

// NewEditor creates a new image interactive editor
func NewEditor() (*editor, int, int) {
	imgFile, w, h := ImgSize()

	micro := canvas.NewImageFromFile(imgDir + "/" + imgFile)
	//micro.FillMode = canvas.ImageFillOriginal

	gc := fyne.NewContainer(iRect(w/2, h/2, w, h, color.RGBA{0, 0, 0, 0})) // gate container
	cc := fyne.NewContainer(iRect(w/2, h/2, w, h, color.RGBA{0, 0, 0, 0})) // cluster container should be independant of gate container for separate initialisaion
	//fgCol := color.Transparent
	//edit := &editor{fg: fgCol, fgPreview: canvas.NewRectangle(fgCol), img: image.NewRGBA(image.Rect(0, 0, 600, 600)), microscop: micro}
	edit := &editor{microscop: micro, min: fyne.Size{Width: float32(w), Height: float32(h)}, gateContainer: cc, clusterContainer: gc, zoom: 100, cacheWidth: w, cacheHeight: h}
	edit.drawSurface = newInteractiveRaster(edit)

	return edit, w, h
}

// BuildUI creates the main window of our application
func (e *editor) BuildUI(w fyne.Window) {
	e.win = w
	e.layer = container.NewMax(e.drawSurface, e.microscop, e.clusterContainer, e.gateContainer)

	w.SetContent(container.NewScroll(e.layer))
}

func (e *editor) setZoom(zoom int) {
	e.zoom = zoom

	h := float32(e.cacheHeight) * float32(zoom) / 100
	w := float32(e.cacheWidth) * float32(zoom) / 100
	size := fyne.Size{Width: float32(w), Height: float32(h)}
	e.min = size
	log.Println("zoom=", zoom, "min=", e.min, "microscope H=", e.cacheHeight)
	//e.updateSizes()
	e.drawSurface.Refresh()
	e.clusterContainer.Refresh()
	e.gateContainer.Refresh()
}

/*
func (e *editor) updateSizes() {
	if e.microscop == nil {
		return
	}
	e.cacheWidth = e.microscop.Bounds().Dx() * e.zoom
	e.cacheHeight = e.microscop.Bounds().Dy() * e.zoom

	c := fyne.CurrentApp().Driver().CanvasForObject(e.status)
	scale := float32(1.0)
	if c != nil {
		scale = c.Scale()
	}
	e.drawSurface.SetMinSize(fyne.NewSize(
		float32(e.cacheWidth)/scale,
		float32(e.cacheHeight)/scale))

	e.renderCache()
}
*/

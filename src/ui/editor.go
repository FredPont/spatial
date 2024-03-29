package ui

import (
	//"fmt"
	"image"
	"image/color"
	"spatial/src/filter"

	//"log"

	//"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/data/binding"
	//"fyne.io/fyne/v2/theme"
	//"fyne.io/fyne/v2/widget"
)

// Editor contains the drawsurface, microscopy image, containers for gates and clusters
type Editor struct {
	drawSurface                     *interactiveRaster
	microscop                       *canvas.Image
	calc                            *canvas.Image   // image calc to display cluster and expression
	min                             fyne.Size       // size of the microscop image and the gate/clusters containers
	layer                           *fyne.Container // container with image and interactive drawsurface
	win                             fyne.Window
	gateContainer                   *fyne.Container // container with the gates lines
	gateNumberContainer             *fyne.Container // container with the gates numbers
	gateDotsContainer               *fyne.Container // container with the gates dots at the polygon edges
	clusterContainer                *fyne.Container // container with the cluster circles
	zoom                            int             // image zoom
	zooMin                          int             // minimal value of zoom to fit the window
	microOrigWidth, microOrigHeight int             // size of the microscop picture not zoomed
	brush                           string          // the brush selected can be "lasso" or "pencil"
}

// GateNB number holds the gate number coordinates and the number of gates starting from 1
type GateNB struct {
	x, y []int // position of the number in the image
	nb   int   // number of gates starting from 1
}

func (e *Editor) draw(w, h int) image.Image {
	return image.NewRGBA(image.Rect(0, 0, w, h))
}

// NewEditor creates a new image interactive editor
func NewEditor() (*Editor, int, int) {
	imgFile, w, h := ImgSize()
	// read the microscopy image
	micro := canvas.NewImageFromFile(imgDir + "/" + imgFile)
	// copy the microscopy image in the temp dir to initialize the calc cluster/express image
	filter.CopyFile(imgDir+"/"+imgFile, "temp/imgOut.png")
	calcImg := canvas.NewImageFromFile("temp/imgOut.png")
	// calcName := "temp/imgOut.png" // calcl image with cluster/expression
	// if filter.FileExist(calcName) {
	// 	calcImg = canvas.NewImageFromFile(calcName)
	// }
	//micro.FillMode = canvas.ImageFillOriginal

	gc := container.NewWithoutLayout(iRect(w/2, h/2, w, h, color.RGBA{0, 0, 0, 0}))  // gate container
	gdc := container.NewWithoutLayout(iRect(w/2, h/2, w, h, color.RGBA{0, 0, 0, 0})) // gate dots container
	gnc := container.NewWithoutLayout(iRect(w/2, h/2, w, h, color.RGBA{0, 0, 0, 0})) // gate number container
	cc := container.NewWithoutLayout(iRect(w/2, h/2, w, h, color.RGBA{0, 0, 0, 0}))  // cluster container should be independant of gate container for separate initialisaion
	//fgCol := color.Transparent
	//edit := &editor{fg: fgCol, fgPreview: canvas.NewRectangle(fgCol), img: image.NewRGBA(image.Rect(0, 0, 600, 600)), microscop: micro}
	edit := &Editor{microscop: micro, calc: calcImg, min: fyne.Size{Width: float32(w), Height: float32(h)}, gateContainer: cc, gateDotsContainer: gdc, gateNumberContainer: gnc, clusterContainer: gc, zoom: 100, microOrigWidth: w, microOrigHeight: h, zooMin: 10}
	edit.drawSurface = newInteractiveRaster(edit)

	return edit, w, h
}

// BuildUI creates the main window of our application
func (e *Editor) BuildUI(w fyne.Window) {
	e.win = w
	e.layer = container.NewStack(e.drawSurface, e.microscop, e.calc, e.clusterContainer, e.gateContainer, e.gateDotsContainer, e.gateNumberContainer)

	w.SetContent(container.NewScroll(e.layer))
}

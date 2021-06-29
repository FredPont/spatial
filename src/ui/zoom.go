package ui

import (
	"fmt"
	"log"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Zoom objet to record the %zoom and change editor min size
type Zoom struct {
	edit *Editor
	zoom *widget.Label
}

// zoom between 10-200%
func (z *Zoom) updateZoom(val int) {
	log.Println("val=", val, "zoom Min=", z.edit.zooMin)
	if val < z.edit.zooMin {
		val = z.edit.zooMin // zoom must be at least the zooMin
	} else if val > 200 {
		val = 200 //zoom Max
	}
	z.edit.setZoom(val)

	z.zoom.SetText(fmt.Sprintf("Zoom : %d%%", z.edit.zoom))
}

// create a zoom widget to increase/decrease size by 10%
// it is not possible to zoom more than 100% of image native size
func newZoom(edit *Editor, a fyne.App) fyne.CanvasObject {
	z := &Zoom{edit: edit, zoom: widget.NewLabel("Zoom : 100%")}
	edit.zoomMin(a) //compute zoom Min
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

func (e *Editor) setZoom(zoom int) {
	initAllLayers(e) // remove clusters and gates 
	e.zoom = zoom

	h := float32(e.microOrigHeight) * float32(zoom) / 100
	w := float32(e.microOrigWidth) * float32(zoom) / 100
	size := fyne.Size{Width: float32(w), Height: float32(h)}
	e.min = size
	log.Println("zoom=", zoom, "min=", e.min, "microscope H=", e.microOrigHeight)

	e.drawSurface.Refresh()
	e.clusterContainer.Refresh()
	e.gateContainer.Refresh()
}

// compute the minimal value of the zoom to fit the windows size
func (e *Editor) zoomMin(a fyne.App) {
	pref := a.Preferences()
	// image native size
	imgH, imgW := e.microOrigHeight, e.microOrigWidth
	// windows size
	winW := binding.BindPreferenceFloat("winW", pref) // set the link to preferences for win width
	wW, _ := winW.Get()
	winH := binding.BindPreferenceFloat("winH", pref) // set the link to preferences for win width
	wH, _ := winH.Get()

	e.zooMin = findMin(imgH, imgW, wH, wW)

}

// find zoom min between 10-200%
// image and zoomed image must be larger than the window
func findMin(imgH, imgW int, wH, wW float64) int {
	zH := 100 * wH / float64(imgH)
	zW := 100 * wW / float64(imgW)
	zMax := math.Max(zH, zW)

	// image and un-zoomed image must be larger than the window

	for i := 10.; i < 200.; i += 10. {
		if zMax == i {
			return int(i)
		} else if zMax > i && zMax < i+10 {
			return int(i + 10)
		}

	}

	return 10
}

// ApplyZoomInt correct the input integer by the current zoom factor
func ApplyZoomInt(e *Editor, val int) int {
	if e.zoom == 100 {
		return val
	}
	return val * e.zoom / 100
}

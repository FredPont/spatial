package ui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Zoom objet to record the %zoom and change editor min size
type Zoom struct {
	edit *editor

	zoom *widget.Label
}

// zoom between 10-100%
func (z *Zoom) updateZoom(val int) {
	log.Println("val=", val)
	if val < 10 {
		val = 10
	} else if val > 100 {
		val = 100
	}
	z.edit.setZoom(val)

	z.zoom.SetText(fmt.Sprintf("%d%%", z.edit.zoom))
}

// create a zoom widget to increase/decrease size by 10%
// it is not possible to zoom more than 100% of image native size
func newZoom(edit *editor, a fyne.App) fyne.CanvasObject {
	z := &Zoom{edit: edit, zoom: widget.NewLabel("100%")}
	edit.zoomMin(a)
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

// compute the minimal value of the zoom to fit the windows size
func (e *editor) zoomMin(a fyne.App) {
	pref := a.Preferences()
	// image native size
	imgH, imgW := e.cacheHeight, e.cacheWidth
	// windows size
	winW := binding.BindPreferenceFloat("winW", pref) // set the link to preferences for win width
	wW, _ := winW.Get()
	winH := binding.BindPreferenceFloat("winH", pref) // set the link to preferences for win width
	wH, _ := winH.Get()

	findMin(imgH, imgW, wH, wW)
}

// find zoom min between 10-100%
func findMin(imgH, imgW int, wH, wW float64) int {
	zH := 100 * wH / float64(imgH)
	zW := 100 * wW / float64(imgW)

	// image and un-zoomed image must be larger than the window

	for i := 10.; i < 100.; i += 10. {
		if zH == i || zW == i {
			return int(i)
		}
		if zH >= zW {
			if zH > i && zH < i+10 {
				return int(i + 10)
			}
		} else {
			if zW > i && zW < i+10 {
				return int(i + 10)
			}
		}
	}

	//log.Println("zoom min", zH, zW)
	return 10
}

package ui

import (
	"fmt"
	"lasso/src/filter"
	"log"
	"math"
	"strconv"

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
func (z *Zoom) updateZoom(val, zoomStep int, f binding.Float) {
	f.Set(0.3) // progress bar
	log.Println("val=", val, "zoom Min=", z.edit.zooMin)
	if val < z.edit.zooMin {
		val = z.edit.zooMin // zoom must be at least the zooMin
		f.Set(0.)           // progress bar
		return
	} else if val > 200 {
		val = 200 //zoom Max
		f.Set(0.) // progress bar
		return
	}
	z.edit.setZoom(val, zoomStep)
	z.zoom.SetText(fmt.Sprintf("Zoom : %d%%", z.edit.zoom))
	f.Set(0.) // progress bar
}

// create a zoom widget to increase/decrease size by 10%
// it is not possible to zoom more than 100% of image native size
func newZoom(edit *Editor, a fyne.App, f binding.Float) fyne.CanvasObject {
	z := &Zoom{edit: edit, zoom: widget.NewLabel("Zoom : 100%")}
	edit.zoomMin(a) //compute zoom Min
	step := 10      // zoom step
	zoom := container.NewHBox(
		widget.NewButtonWithIcon("", theme.ZoomOutIcon(), func() {
			go z.updateZoom(z.edit.zoom-step, -step, f)
		}),
		z.zoom,
		widget.NewButtonWithIcon("", theme.ZoomInIcon(), func() {
			go z.updateZoom(z.edit.zoom+step, step, f)
		}))
	return zoom
}

func (e *Editor) setZoom(zoom, zoomStep int) {
	//initAllLayers(e) // remove clusters and gates
	initCluster(e)
	initGatesContainer(e)
	e.zoom = zoom

	h := float32(e.microOrigHeight) * float32(zoom) / 100
	w := float32(e.microOrigWidth) * float32(zoom) / 100
	size := fyne.Size{Width: float32(w), Height: float32(h)}
	e.min = size
	log.Println("zoom=", zoom, "min=", e.min, "microscope H=", e.microOrigHeight)

	// zoom the gates
	zoomGates(e, zoomStep)
	redrawGates(e)
	redrawGatesNB(e.drawSurface) // redraw gates numbers
	e.drawSurface.Refresh()
	//e.clusterContainer.Refresh()
	//e.gateContainer.Refresh()
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

// apply zoom step to the polygon edges and gates numbers coordinates
func zoomGates(e *Editor, zoomStep int) {

	// update points coordinates is not usefull
	// because points are reset when polygon is closed

	// zoom factor between 2 consecutive scales. ex 90/100 -> 80/90 -> 70/80 ...
	zf := float64(e.zoom) / (float64(e.zoom) - float64(zoomStep))
	// update alledges coordinates
	L := len(e.drawSurface.alledges)
	for i := 0; i < L; i++ {
		zoomPoints(e.drawSurface.alledges[i], zf)
	}
	zoomGateNumbers(e.drawSurface.gatesNumbers, zf)
	//log.Println(e.drawSurface.alledges)
}

// update point coordinates with zf and modify initial point
func zoomPoints(p []filter.Point, zf float64) {
	for i := 0; i < len(p); i++ {
		//(*p)[i].X = int(float64((*p)[i].X) * zf)
		//(*p)[i].Y = int(float64((*p)[i].Y) * zf)
		p[i].X = int(math.Round(float64(p[i].X) * zf))
		p[i].Y = int(math.Round(float64(p[i].Y) * zf))
	}
}

func zoomGateNumbers(gn GateNB, zf float64) {
	for i := 0; i < len(gn.x); i++ {
		gn.x[i] = int(math.Round(float64(gn.x[i]) * zf))
		gn.y[i] = int(math.Round(float64(gn.y[i]) * zf))
	}
}

// redraw the gates
func redrawGates(e *Editor) {
	L := len(e.drawSurface.alledges)
	for i := 0; i < L; i++ {
		if i == L-1 {
			// update the last gate coordinates
			e.drawSurface.gatesLines = redrawlastGate(e.drawSurface, e.drawSurface.alledges[i])
		} else {
			redrawpolygon(e.drawSurface, e.drawSurface.alledges[i])
		}

	}
	//e.gateContainer.Refresh()

}

// redraw one polygon
func redrawpolygon(r *interactiveRaster, p []filter.Point) {
	L := len(p)
	if L < 1 {
		return
	}
	for i := 0; i < L-1; i++ {
		r.drawline(p[i].X, p[i].Y, p[i+1].X, p[i+1].Y)
	}
	r.drawline(p[0].X, p[0].Y, p[L-1].X, p[L-1].Y) // close the polygon
}

// redraw last gate and store it
func redrawlastGate(r *interactiveRaster, p []filter.Point) []fyne.CanvasObject {
	var lastGate []fyne.CanvasObject
	var line fyne.CanvasObject
	L := len(p)
	for i := 0; i < L-1; i++ {
		line = r.drawline(p[i].X, p[i].Y, p[i+1].X, p[i+1].Y)
		lastGate = append(lastGate, line)
	}
	line = r.drawline(p[0].X, p[0].Y, p[L-1].X, p[L-1].Y) // close the polygon
	lastGate = append(lastGate, line)
	return lastGate
}

// redraw the gate numbers for all gates
func redrawGatesNB(r *interactiveRaster) {
	initGatesNBwindow(r.edit) // clear gateNB in window
	gn := r.gatesNumbers
	for i := 0; i < len(gn.x); i++ {
		gateNB := strconv.Itoa(i)
		r.drawGateNb(gn.x[i], gn.y[i], gateNB)
	}
}

// redraw the gate numbers for last gate only
func redrawLastGatesNB(r *interactiveRaster) {
	gn := r.gatesNumbers
	L := len(gn.x) - 1
	gateNB := strconv.Itoa(L)
	r.drawGateNb(gn.x[L], gn.y[L], gateNB)

}

// draw and store the gates numbers coordinates after import gate
func drawImportedGatesNB(r *interactiveRaster) {
	nbGates := len(r.edit.drawSurface.alledges)
	// store the gateNB coordinates and gateNB from imported gates
	for _, gate := range r.edit.drawSurface.alledges {
		p := filter.PopPointItem(gate)
		r.gatesNumbers.x = append(r.gatesNumbers.x, p.X)
		r.gatesNumbers.y = append(r.gatesNumbers.y, p.Y)
	}
	r.gatesNumbers.nb = nbGates
	// draw the gatesNB from the coordinated just stored above
	redrawGatesNB(r)
}

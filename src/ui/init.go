package ui

import (
	"spatial/src/filter"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func initAlledges(e *Editor) {
	e.drawSurface.alledges = nil
	e.drawSurface.points = nil
}

// clear only the last gates edges and point and lines in container
func initLastedges(e *Editor) {
	e.drawSurface.alledges = filter.PopPoints(e.drawSurface.alledges)
	e.drawSurface.points = nil
	e.drawSurface.gatesLines = nil
}

func initCluster(e *Editor) {
	e.clusterContainer.Objects = nil
}

func clearCluster(e *Editor) {
	e.clusterContainer.Objects = nil
	e.clusterContainer.Refresh()
}

func initGates(e *Editor) {
	e.gateContainer.Objects = nil
	e.drawSurface.gatesLines = nil
	e.gateContainer.Refresh()
	initAlledges(e)      // reset alledges
	initGatesNB(e)       // clear all gates numbers in arrays
	initGatesNBwindow(e) // clear all gates numbers displayed
}

func initAllLayers(e *Editor) {
	initCluster(e)
	initGates(e)
}

func initGatesContainer(e *Editor) {
	e.gateContainer.Objects = nil
	e.gateContainer.Refresh()

}

// clear all stored gates numbers from arrays
func initGatesNB(e *Editor) {
	e.drawSurface.gatesNumbers.x = nil
	e.drawSurface.gatesNumbers.y = nil
	e.drawSurface.gatesNumbers.nb = 0
}

// clear all displayed gates numbers in the picture
func initGatesNBwindow(e *Editor) {
	e.gateNumberContainer.Objects = nil
	e.gateNumberContainer.Refresh()
}

// clear last gate number coordinates and decrease gateNB
func initLastGatesNB(e *Editor) {
	e.drawSurface.gatesNumbers.x = filter.PopIntArray(e.drawSurface.gatesNumbers.x)
	e.drawSurface.gatesNumbers.y = filter.PopIntArray(e.drawSurface.gatesNumbers.y)
	if e.drawSurface.gatesNumbers.nb > 0 {
		e.drawSurface.gatesNumbers.nb--
	}

}

// clear the dots at the edges of the currently drawing gate
func initGateDots(e *Editor) {
	e.gateDotsContainer.Objects = nil
	e.gateDotsContainer.Refresh()
}

////////////////////////////
// Expression map INIT
////////////////////////////

// reset expression min max slider
func initSliderExp(MaxExp, MinExp *widget.Slider) {
	preference := fyne.CurrentApp().Preferences()
	preference.SetFloat("userMaxExp", 100)
	preference.SetFloat("userMinExp", 0)
	MaxExp.Value = 100. // reset slider position
	MinExp.Value = 0.   // reset slider position

}

////////////////////////////
// interactive 2D plot INIT
////////////////////////////

// clear gates from screen and initialise inter2D.drawSurface.alledges
func init2DScatterGates(inter2D *Interactive2Dsurf) {
	inter2D.drawSurface.alledges = nil
	inter2D.gateContainer.Objects = nil
	inter2D.drawSurface.gatesNumbers.x = nil
	inter2D.drawSurface.gatesNumbers.y = nil
	inter2D.drawSurface.gatesNumbers.nb = 0
	inter2D.gateContainer.Refresh()
}

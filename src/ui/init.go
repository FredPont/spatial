package ui

import "lasso/src/filter"

func initAlledges(e *Editor) {
	e.drawSurface.alledges = nil
	e.drawSurface.points = nil
}

// clear only the last gates edges and point
func initLastedges(e *Editor) {
	e.drawSurface.alledges = filter.PopPoints(e.drawSurface.alledges)
	e.drawSurface.points = nil
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
	e.gateContainer.Refresh()
	initAlledges(e) // reset alledges
	initGatesNB(e)  // clear all gates numbers
}

func initAllLayers(e *Editor) {
	initCluster(e)
	initGates(e)
}

func initGatesContainer(e *Editor) {
	e.gateContainer.Objects = nil
	e.gateContainer.Refresh()

}

// clear all gates numbers
func initGatesNB(e *Editor) {
	e.drawSurface.gatesNumbers.x = nil
	e.drawSurface.gatesNumbers.y = nil
	e.drawSurface.gatesNumbers.nb = 0
}

// clear last gate number coordinates and decrease gateNB
func initLastGatesNB(e *Editor) {
	e.drawSurface.gatesNumbers.x = filter.PopIntArray(e.drawSurface.gatesNumbers.x)
	e.drawSurface.gatesNumbers.y = filter.PopIntArray(e.drawSurface.gatesNumbers.y)
	if e.drawSurface.gatesNumbers.nb > 0 {
		e.drawSurface.gatesNumbers.nb--
	}

}

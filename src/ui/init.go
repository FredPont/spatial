package ui

import "lasso/src/filter"

func initAlledges(e *editor) {
	e.drawSurface.alledges = nil
	e.drawSurface.points = nil
}

// clear only the last gates edges and point
func initLastedges(e *editor) {
	e.drawSurface.alledges = filter.PopPoints(e.drawSurface.alledges)
	e.drawSurface.points = nil
}

func initCluster(e *editor) {
	e.clusterContainer.Objects = nil
}

func clearCluster(e *editor) {
	e.clusterContainer.Objects = nil
	e.clusterContainer.Refresh()
}

func initGates(e *editor) {
	e.gateContainer.Objects = nil
	e.gateContainer.Refresh()
	initAlledges(e) // reset alledges
}

package ui

func initAlledges(e *editor) {
	e.drawSurface.alledges = nil
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
}

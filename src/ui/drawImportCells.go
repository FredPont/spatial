package ui

import (
	"image/color"
	"spatial/src/filter"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

func buttonImportCells(a fyne.App, e *Editor, preference fyne.Preferences, f binding.Float, impCellFindex int, header []string, firstTable string) {
	f.Set(0.3) // progress bar
	importedCells, files := importCells()
	f.Set(0.5) // progress bar
	iCellFI := binding.BindPreferenceInt("imported file index", preference)
	impCellFindex, _ = iCellFI.Get()
	if impCellFindex > len(importedCells)-1 {
		impCellFindex = 0
	}
	cellImport := filter.StrToMap(importedCells[impCellFindex])
	cellfile, _ := filter.RemExt(files[impCellFindex])
	//drawImportCells(a, e, header, firstTable, f, cellImport, cellfile)
	drawImportCellsImg(a, e, header, firstTable, f, cellImport, cellfile)
	// increment file index by 1
	if impCellFindex < len(importedCells)-1 {
		preference.SetInt("imported file index", impCellFindex+1)
	} else {
		preference.SetInt("imported file index", 0)
	}
	f.Set(0.) // reset progress bar
}

// importCells import all cells in all files in import_cells.
func importCells() ([][]string, []string) {
	var importedCells [][]string
	// read import_cells dir
	files := filter.ListFiles("import_cells")
	for _, f := range files {
		cells := filter.ReadImportedCells(f)
		importedCells = append(importedCells, cells)
	}
	return importedCells, files
}

func getImportedCells(a fyne.App, header []string, filename string, cellImport map[string]bool) map[int][]filter.Point {
	pref := a.Preferences()
	// X coordinates
	xcor := binding.BindPreferenceString("xcor", pref) // set the link to preferences for x coordinates
	xc, _ := xcor.Get()

	// y coordinates
	ycor := binding.BindPreferenceString("ycor", pref) // set the link to preferences for y coordinates
	yc, _ := ycor.Get()

	// cluster column
	clustercolumn := binding.BindPreferenceString("clustcol", pref) // set the link to preferences for cluster cluster column
	clucol, _ := clustercolumn.Get()

	// add the cellnames colums to col indexes
	colIndexes := []int{0}
	nextIndexes := filter.GetColIndex(header, []string{clucol, xc, yc})
	colIndexes = append(colIndexes, nextIndexes...)
	return filter.ClustersByCells(a, filename, colIndexes, cellImport)
}

func drawImportCells(a fyne.App, e *Editor, header []string, filename string, f binding.Float, cellImport map[string]bool, cellfile string) {
	initCluster(e) // remove all dots of the cluster container
	pref := a.Preferences()
	clustOp := binding.BindPreferenceFloat("clustOpacity", pref) // cluster opacity
	opacity, _ := clustOp.Get()
	op := uint8(opacity)
	clustDia := binding.BindPreferenceInt("clustDotDiam", pref) // cluster dot diameter
	diameter, _ := clustDia.Get()
	legendDiameter := ApplyZoomInt(e, diameter)

	clusterMap := getImportedCells(a, header, filename, cellImport) // cluster nb => []Point
	//log.Println(clusterMap, "clusters detected")

	nbCluster := len(clusterMap)
	clustNames := filter.KeysIntPoint(clusterMap)
	//log.Println(clustNames, "Name clusters detected")

	legendPosition := filter.Point{X: 15, Y: 15} // initial legend position for cluster names
	title(e, cellfile)                           // draw title with file name

	spotsToDisplay := make([]fyne.CanvasObject, 0)

	for c := 0; c < nbCluster; c++ {
		f.Set(float64(c) / float64(nbCluster-1)) // % progression for progress bar
		coordinates := clusterMap[clustNames[c]]
		clcolor := ClusterColors(nbCluster, c)
		for i := 0; i < len(coordinates); i++ {
			//e.drawcircle(ApplyZoomInt(e, coordinates[i].X), ApplyZoomInt(e, coordinates[i].Y), diameter, color.NRGBA{clcolor.R, clcolor.G, clcolor.B, op})
			//add the spot to the slice of objects
			spotsToDisplay = append(spotsToDisplay, drawRoundedRect(ApplyZoomInt(e, coordinates[i].X), ApplyZoomInt(e, coordinates[i].Y), diameter, color.NRGBA{clcolor.R, clcolor.G, clcolor.B, op}))

		}
		// draw legend dot and name for the current cluster
		impCellLegend(e, clcolor.R, clcolor.G, clcolor.B, op, legendPosition.X, legendPosition.Y, legendDiameter, clustNames[c])
		legendPosition.Y = legendPosition.Y + 30
	}
	e.clusterContainer.Objects = append(e.clusterContainer.Objects, spotsToDisplay...)
	e.clusterContainer.Refresh()
}

func impCellLegend(e *Editor, R, G, B, op uint8, x, y, diameter, clusterName int) {
	AbsText(e.clusterContainer, x+20, y+10, strconv.Itoa(clusterName), 20, color.NRGBA{50, 50, 50, 255})
	e.drawcircle(x, y, diameter*100/e.zoom, color.NRGBA{R, G, B, op})
}

// pring cell file names on the cluster plot
func title(e *Editor, cellfile string) {
	AbsText(e.clusterContainer, 100, 20, cellfile, 20, color.NRGBA{50, 50, 50, 255})
}

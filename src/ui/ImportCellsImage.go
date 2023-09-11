package ui

import (
	"spatial/src/filter"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/fogleman/gg"
)

func drawImportCellsImg(a fyne.App, e *Editor, header []string, filename string, f binding.Float, cellImport map[string]bool, cellfile string) {
	initCluster(e) // remove all dots of the cluster container
	pref := a.Preferences()

	//calc image size to draw the cluster dots
	H := pref.IntWithFallback("imgH", 500)
	W := pref.IntWithFallback("imgW", 500)
	// clear calc image
	dc := gg.NewContext(W, H)
	dc.SetRGBA(0, 0, 0, 0) // create a transparent image
	dc.Clear()

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

	//drawColorMono(e, dc, nbCluster, clustNames, op, clusterMap, diameter, f)
	for c := 0; c < nbCluster; c++ {
		f.Set(float64(c) / float64(nbCluster-1)) // % progression for progress bar
		coordinates := clusterMap[clustNames[c]]
		clcolor := ClusterColors(nbCluster, c)
		for i := 0; i < len(coordinates); i++ {
			dc.SetRGBA(float64(clcolor.R)/255.0, float64(clcolor.G)/255.0, float64(clcolor.B)/255.0, float64(op)/255.0)
			dc.DrawPoint(float64(coordinates[i].X), float64(coordinates[i].Y), float64(diameter))
			dc.Fill()

		}
		// draw legend dot and name for the current cluster
		impCellLegend(e, clcolor.R, clcolor.G, clcolor.B, op, legendPosition.X, legendPosition.Y, legendDiameter, clustNames[c])
		legendPosition.Y = legendPosition.Y + 30
	}
	dc.SavePNG("temp/imgOut.png")
	e.layer.Refresh()
	f.Set(0.) // reset progress bar
}

// draw the spots colored by clusters
// func drawColorByCluster(e *Editor, dc *gg.Context, legendPosition filter.Point, nbCluster int, clustNames []int, op uint8, clusterMap map[int][]filter.Point, diameter, legendDiameter int, f binding.Float) {
// 	for c := 0; c < nbCluster; c++ {
// 		f.Set(float64(c) / float64(nbCluster-1)) // % progression for progress bar
// 		coordinates := clusterMap[clustNames[c]]
// 		clcolor := ClusterColors(nbCluster, c)
// 		for i := 0; i < len(coordinates); i++ {
// 			dc.SetRGBA(float64(clcolor.R)/255.0, float64(clcolor.G)/255.0, float64(clcolor.B)/255.0, float64(op)/255.0)
// 			dc.DrawPoint(float64(coordinates[i].X), float64(coordinates[i].Y), float64(diameter))
// 			dc.Fill()

// 		}
// 		// draw legend dot and name for the current cluster
// 		impCellLegend(e, clcolor.R, clcolor.G, clcolor.B, op, legendPosition.X, legendPosition.Y, legendDiameter, clustNames[c])
// 		legendPosition.Y = legendPosition.Y + 30
// 	}
// }

// draw the spots with one color
// func drawColorMono(e *Editor, dc *gg.Context, nbCluster int, clustNames []int, op uint8, clusterMap map[int][]filter.Point, diameter int, f binding.Float) {
// 	for c := 0; c < nbCluster; c++ {
// 		f.Set(float64(c) / float64(nbCluster-1)) // % progression for progress bar
// 		coordinates := clusterMap[clustNames[c]]
// 		clcolor := RGB{255, 128, 0}
// 		for i := 0; i < len(coordinates); i++ {
// 			dc.SetRGBA(float64(clcolor.R)/255.0, float64(clcolor.G)/255.0, float64(clcolor.B)/255.0, float64(op)/255.0)
// 			dc.DrawPoint(float64(coordinates[i].X), float64(coordinates[i].Y), float64(diameter))
// 			dc.Fill()

// 		}
// 		// draw legend dot and name for the current cluster
// 		//impCellLegend(e, clcolor.R, clcolor.G, clcolor.B, op, legendPosition.X, legendPosition.Y, legendDiameter, clustNames[c])
// 		//legendPosition.Y = legendPosition.Y + 30
// 	}
// }

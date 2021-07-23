package ui

import (
	"image/color"
	"lasso/src/filter"
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
)

func getClusters(a fyne.App, header []string, filename string) map[int][]filter.Point {
	pref := a.Preferences()
	// X coordinates
	xcor := binding.BindPreferenceString("xcor", pref) // set the link to preferences for rotation
	xc, _ := xcor.Get()

	// y coordinates
	ycor := binding.BindPreferenceString("ycor", pref) // set the link to preferences for rotation
	yc, _ := ycor.Get()

	// cluster column
	clustercolumn := binding.BindPreferenceString("clustcol", pref) // set the link to preferences for rotation
	clucol, _ := clustercolumn.Get()

	colIndexes := filter.GetColIndex(header, []string{clucol, xc, yc})
	return filter.ReadClusters(a, filename, colIndexes)
}

func drawClusters(a fyne.App, e *Editor, header []string, filename string, f binding.Float) {
	initCluster(e) // remove all dots of the cluster container
	pref := a.Preferences()
	clustOp := binding.BindPreferenceFloat("clustOpacity", pref) // cluster opacity
	opacity, _ := clustOp.Get()
	op := uint8(opacity)
	clustDia := binding.BindPreferenceInt("clustDotDiam", pref) // cluster dot diameter
	diameter, _ := clustDia.Get()
	diameter = ApplyZoomInt(e, diameter)

	clusterMap := getClusters(a, header, filename) // cluster nb => []Point
	log.Println(len(clusterMap), "clusters detected")

	nbCluster := len(clusterMap)
	clustNames := filter.KeysIntPoint(clusterMap)

	legendPosition := filter.Point{X: 15, Y: 15} // initial legend position for cluster names

	for c := 0; c < nbCluster; c++ {
		f.Set(float64(c) / float64(nbCluster-1)) // % progression for progress bar
		coordinates := clusterMap[clustNames[c]]
		clcolor := ClusterColors(nbCluster, c)
		for i := 0; i < len(coordinates); i++ {
			e.drawcircle(ApplyZoomInt(e, coordinates[i].X), ApplyZoomInt(e, coordinates[i].Y), diameter, color.NRGBA{clcolor.R, clcolor.G, clcolor.B, op})
		}
		// draw legend dot and name for the current cluster
		drawLegend(e, clcolor.R, clcolor.G, clcolor.B, op, legendPosition.X, legendPosition.Y, diameter, clustNames[c])
		legendPosition.Y = legendPosition.Y + 30
	}

	e.clusterContainer.Refresh()
	f.Set(0.) // reset progress bar
}

func drawLegend(e *Editor, R, G, B, op uint8, x, y, diameter, clusterName int) {
	AbsText(e.clusterContainer, x+20, y+10, strconv.Itoa(clusterName), 20, color.NRGBA{50, 50, 50, 255})
	e.drawcircle(x, y, diameter*100/e.zoom, color.NRGBA{R, G, B, op})
}

// credits : https://github.com/ajstarks/fc
// iCircle draws a circle centered at (x,y)
func iCircle(x, y, r int, color color.NRGBA) *canvas.Circle {
	fx, fy, fr := float32(x), float32(y), float32(r)
	p1 := fyne.Position{X: fx - fr, Y: fy - fr}
	p2 := fyne.Position{X: fx + fr, Y: fy + fr}
	c := &canvas.Circle{FillColor: color, Position1: p1, Position2: p2}
	return c
}

// drawline a circle at x,y position to the cluster container
func (e *Editor) drawcircle(x, y, ray int, color color.NRGBA) fyne.CanvasObject {
	c := iCircle(x, y, ray, color)  // draw circle rayon ray
	e.clusterContainer.AddObject(c) // add the cicle to the cluster container
	return c
}

// AbsText places text within a container
// credits : https://github.com/ajstarks/fc
func AbsText(cont *fyne.Container, x, y int, s string, size int, color color.NRGBA) {
	fx, fy, fsize := float32(x), float32(y), float32(size)
	t := &canvas.Text{Text: s, Color: color, TextSize: fsize}
	adj := fsize / 5
	p := fyne.Position{X: fx, Y: fy - (fsize + adj)}
	t.Move(p)
	cont.AddObject(t)
}

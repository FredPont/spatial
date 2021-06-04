package ui

import (
	"image/color"
	"lasso/src/filter"

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

func drawClusters(a fyne.App, e *editor, header []string, filename string) {
	initCluster(e) // remove all dots of the cluster container
	pref := a.Preferences()
	clustOp := binding.BindPreferenceFloat("clustOpacity", pref) // cluster opacity
	opacity, _ := clustOp.Get()
	op := uint8(opacity)
	clustDia := binding.BindPreferenceInt("clustDotDiam", pref) // cluster dot diameter
	diameter, _ := clustDia.Get()

	clusterMap := getClusters(a, header, filename)

	for clusterNB, coordinates := range clusterMap {
		_ = clusterNB
		for i := 0; i < len(coordinates); i++ {
			e.drawcircle(coordinates[i].X, coordinates[i].Y, diameter, color.NRGBA{30, 144, 255, op})
		}

	}

	// for x := 10; x < 1000; x += 10 {
	// 	e.drawcircle(x, x, diameter, color.NRGBA{30, 144, 255, op})
	// }
	e.clusterContainer.Refresh()
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
func (e *editor) drawcircle(x, y, ray int, color color.NRGBA) fyne.CanvasObject {
	c := iCircle(x, y, ray, color)  // draw circle rayon ray
	e.clusterContainer.AddObject(c) // add the cicle to the cluster container
	return c
}

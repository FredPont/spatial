package ui

import (
	"image/color"
	"spatial/src/filter"
	"spatial/src/plot"
	"spatial/src/pogrebDB"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
)

func getClusters(a fyne.App, header []string, filename string) map[int][]filter.Point {
	pref := a.Preferences()
	// X coordinates
	xcor := binding.BindPreferenceString("xcor", pref) // set the link to preferences for  X coordinates
	xc, _ := xcor.Get()

	// y coordinates
	ycor := binding.BindPreferenceString("ycor", pref) // set the link to preferences for y coordinates
	yc, _ := ycor.Get()

	// cluster column
	clustercolumn := binding.BindPreferenceString("clustcol", pref) // set the link to preferences for cluster column
	clucol, _ := clustercolumn.Get()

	colIndexes := filter.GetColIndex(header, []string{clucol, xc, yc})

	useDBpref := binding.BindPreferenceBool("useDataBase", pref)
	useDB, _ := useDBpref.Get()
	// use the pogreb database instead of CSV if selected by the user in preferences
	if useDB {
		return pogrebDB.ReadClustersDB(a, filename, colIndexes)
	}

	return filter.ReadClusters(a, filename, colIndexes)
}

// draw the cluster legend : color + cluster number
func drawLegend(e *Editor, R, G, B, op uint8, x, y, diameter, clusterName int, colorText color.NRGBA) {
	AbsText(e.clusterContainer, x+20, y+10, strconv.Itoa(clusterName), 20, colorText)
	// compute the spot max diameter to avoid overlap
	spotDiam := diameter * 100 / e.zoom
	if spotDiam >= 15 {
		spotDiam = 15
	} else if spotDiam < 5 {
		spotDiam = 5
	}
	e.drawcircle(x, y, spotDiam, color.NRGBA{R, G, B, op})
}

func getLegendColor(a fyne.App) color.NRGBA {
	R, G, B, _ := plot.GetPrefColorRGBA(a, "legendColR", "legendColG", "legendColB", "legendColA")
	colorText := color.NRGBA{uint8(R), uint8(G), uint8(B), 255}
	return colorText
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

// credits : https://github.com/ajstarks/fc
// roundRect makes a rectangle centered at x,y with rounded angles such as it is displayed like a circle :)
// the reason is that rectangle is hardware accelerated and not (yet) the circle
func drawRoundedRect(x, y, rad int, color color.NRGBA) *canvas.Rectangle {
	w := float32(2 * rad) // w=h, we draw a square with rounded corners
	//h := 2 * rad
	fx, fy, fw, fh := float32(x), float32(y), w, w
	r := &canvas.Rectangle{FillColor: color, CornerRadius: float32(rad)}
	r.Move(fyne.Position{X: fx - (fw / 2), Y: fy - (fh / 2)})
	r.Resize(fyne.Size{Width: fw, Height: fh})
	return r
}

// drawcircle, draws a circle at x,y position to the cluster container
func (e *Editor) drawcircle(x, y, ray int, color color.NRGBA) fyne.CanvasObject {
	c := iCircle(x, y, ray, color) // draw circle rayon ray
	e.clusterContainer.Add(c)      // add the cicle to the cluster container
	//e.clusterContainer.Objects = append(e.clusterContainer.Objects, c)
	return c
}

// drawcircle, draws a circle at x,y position to the cluster container
//func (e *Editor) drawRoundedRect(x, y, ray int, color color.NRGBA) fyne.CanvasObject {
//	return roundRect(x, y, ray, color) // draw rectangle
//}

// AbsText places text within a container
// credits : https://github.com/ajstarks/fc
func AbsText(cont *fyne.Container, x, y int, s string, size int, color color.NRGBA) {
	fx, fy, fsize := float32(x), float32(y), float32(size)
	t := &canvas.Text{Text: s, Color: color, TextSize: fsize}
	adj := fsize / 5
	p := fyne.Position{X: fx, Y: fy - (fsize + adj)}
	t.Move(p)
	cont.Add(t)
}

package ui

import (
	"image/color"
	"log"
	"spatial/src/filter"
	"strconv"

	"fyne.io/fyne/v2"
)

// PlotBox holds the essential data for making a chart
type PlotBox struct {
	Title                    string
	id                       []string
	X                        []float64
	Y                        []float64
	Color                    color.NRGBA
	Top, Bottom, Left, Right float64
	Xmax, Xmin, Ymax, Ymin   float64
	winH, winW               float64
}

// newHisto read the pvalue table and return a PlotBox
func readVulcano(title string, pvTable []PVrecord) PlotBox {
	var x, y []float64
	var items []string
	m := 1e308
	Xmax := -m
	Xmin := m
	Ymax := -m
	Ymin := m

	for i := 0; i < len(pvTable); i++ {
		items = append(items, pvTable[i].item)
		xval := pvTable[i].log2fc
		yval := pvTable[i].log10pv
		x = append(x, xval)
		y = append(y, yval)
		if xval > Xmax {
			Xmax = xval
		}
		if xval < Xmin {
			Xmin = xval
		}
		if yval > Ymax {
			Ymax = yval
		}
		if yval < Ymin {
			Ymin = yval
		}
	}
	return PlotBox{
		Title:  title,
		id:     items,
		X:      x,
		Y:      y,
		Xmax:   Xmax,
		Xmin:   Xmin,
		Ymax:   Ymax,
		Ymin:   Ymin,
		Top:    50.,
		Bottom: 50.,
		Left:   50.,
		Right:  50.,
		winH:   800,
		winW:   800,
	}
}

// MapRange -- given a value between low1 and high1, return the corresponding value between low2 and high2
// credits : https://github.com/ajstarks/fc
func MapRange(value, low1, high1, low2, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

// buildVulanoPlot : create the window and the vulcano plot
func buildVulanoPlot(e *Editor, header []string, fname, firstTable string, pvfcTable []PVrecord) {
	vulcBox := readVulcano(fname, pvfcTable)
	//log.Println(readVulcano(fname, pvfcTable))
	v := buildVulcWin(e)
	v.drawSurface.vulcBox = vulcBox
	drawVulcano(v, vulcBox)
	v.header = header
	v.tableName = firstTable
	buildVulcanoTools(v)

}

// drawline a circle at x,y position to the scatter container
func (e *Vulcano) drawcircle(x, y, ray int, color color.NRGBA) fyne.CanvasObject {
	c := iCircle(x, y, ray, color) // draw circle rayon ray
	e.scatterContainer.Add(c)      // add the cicle to the cluster container
	return c
}

func drawVulcano(v *Vulcano, vulcBox PlotBox) {
	R := uint8(106)
	G := uint8(90)
	B := uint8(250)

	vulcBox.Color = color.NRGBA{R, G, B, 255}

	//draw axes
	vulcBox.XAxis(v)
	vulcBox.YAxis(v)

	// draw scatter plot
	vulcBox.Scatter(v, 3)

	v.scatterContainer.Refresh()
}

// Scatter makes a scatter chart
func (p *PlotBox) Scatter(v *Vulcano, dotsize int) {

	for i, xplot := range p.X {
		//x := int(MapRange(x, p.Xmin, p.Xmax, p.Left, 800-p.Right))
		//y := int(MapRange(p.Y[i], p.Ymin, p.Ymax, p.Bottom, 800-p.Top))
		x, y := xCoord(p, xplot), yCoord(p, p.Y[i])
		v.drawcircle(x, y, dotsize, p.Color)
	}

}

// xCoord compute the x windows coordinates of a dot
// from its x scatter plot coordinate
func xCoord(p *PlotBox, xplot float64) int {
	xwin := MapRange(xplot, p.Xmin, p.Xmax, p.Left, p.winW-p.Right)
	return int(xwin)
}

// yCoord compute the y windows coordinates of a dot
// from its y scatter plot coordinate
func yCoord(p *PlotBox, yplot float64) int {
	ywin := p.winH - (MapRange(yplot, p.Ymin, p.Ymax, p.Bottom, p.winH-p.Top))
	return int(ywin)
}

// XAxis makes the X axis
func (p *PlotBox) XAxis(v *Vulcano) {
	y1 := 0
	lef := xCoord(p, p.Xmin)
	rig := xCoord(p, p.Xmax)

	if yZero(p) {
		y1 = yCoord(p, 0.)
	} else {
		y1 = yCoord(p, p.Ymin)
		log.Println("Y axis does not contain 0 value !")
	}
	//log.Println("y axis:", x1, bot, up)
	c := iLine(lef, y1, rig, y1, 1, color.RGBA{0, 0, 0, 255})
	v.scatterContainer.Add(c) // add the line to the cluster container
	p.Xlabel(v, y1)
	//v.scatterContainer.Refresh()
}

// YAxis makes the Y axis
func (p *PlotBox) YAxis(v *Vulcano) {
	x1 := 0
	bot := yCoord(p, p.Ymin)
	up := yCoord(p, p.Ymax)

	if xZero(p) {
		x1 = xCoord(p, 0.) // if x=0 exists, the Y axes is plot on 0.
	} else {
		if p.Xmax < 0 { // if all x values < 0, the y axes is on the right
			x1 = xCoord(p, p.Xmax)
		} else {
			x1 = xCoord(p, p.Xmin)
		}
	}
	//log.Println("y axis:", x1, bot, up)
	c := iLine(x1, bot, x1, up, 1, color.RGBA{0, 0, 0, 255})
	v.scatterContainer.Add(c) // add the line to the cluster container
	p.Ylabel(v, x1)
	//v.scatterContainer.Refresh()
}

// test if x=0 exists
func xZero(p *PlotBox) bool {
	if p.Xmax >= 0 && p.Xmin <= 0 {
		return true
	}
	return false
}

// test if y=0 exists
func yZero(p *PlotBox) bool {
	if p.Ymax >= 0 && p.Ymin <= 0 {
		return true
	}
	return false
}

// Xlabel makes the x axis scale text
func (p *PlotBox) Xlabel(v *Vulcano, y int) {
	ntic := 10
	//var positions []int                                    // ticks position in pixels
	labels := filter.TicInterval(p.Xmin, p.Xmax, ntic) // ticks labels with decimal
	//log.Println("labels", labels)
	//positions := filter.TicPixelPos(Xmin, Xmax, ntic)      // ticks position in pixels
	for _, po := range labels {
		str := TicksDecimals(po)
		x := xCoord(p, po)
		AbsText(v.scatterContainer, x-10, y+20, str, 10, color.NRGBA{0, 0, 0, 255}) // label
		ti := iLine(x, y, x, y+5, 1, color.RGBA{0, 0, 0, 255})                      // tick
		v.scatterContainer.Add(ti)                                                  // add the tick to the cluster container
	}
	AbsText(v.scatterContainer, xCoord(p, (p.Xmax+p.Xmin)/2), y+35, "log2(FC) (group_2/group_1)", 12, color.NRGBA{0, 0, 0, 255}) // axis title
}

// Ylabel makes the x axis scale text
func (p *PlotBox) Ylabel(v *Vulcano, x int) {
	ntic := 10
	//var positions []int                                    // ticks position in pixels
	labels := filter.TicInterval(p.Ymin, p.Ymax, ntic) // ticks labels with decimal
	//log.Println("labels", labels)
	//positions := filter.TicPixelPos(Xmin, Xmax, ntic)      // ticks position in pixels
	for _, po := range labels {
		str := TicksDecimals(po)
		y := yCoord(p, po)
		AbsText(v.scatterContainer, x+10, y, str, 10, color.NRGBA{0, 0, 0, 255}) // label
		ti := iLine(x, y-5, x+5, y-5, 1, color.RGBA{0, 0, 0, 255})               // tick
		v.scatterContainer.Add(ti)                                               // add the tick to the cluster container
	}
	AbsText(v.scatterContainer, x-35, yCoord(p, p.Ymax)-25, "log10(Pvalue)", 12, color.NRGBA{0, 0, 0, 255}) // axis title
}

// TicksDecimals format the number of decimals for ticks
func TicksDecimals(po float64) string {
	var str string
	neg := false
	// decimal format of ticks
	if po < 0 {
		po = -po // format negative numbers
		neg = true
	}
	if filter.Between(1, 10, po) {
		str = strconv.FormatFloat(po, 'f', 1, 64)
	} else if po == 0 {
		str = " 0"
	} else if filter.Between(0.1, 1, po) {
		str = strconv.FormatFloat(po, 'f', 2, 64)
	} else if filter.Between(0.001, 0.1, po) {
		str = strconv.FormatFloat(po, 'f', 3, 64)
	} else if po < 0.001 {
		str = strconv.FormatFloat(po, 'e', 1, 64)
	} else if filter.Between(10, 1000, po) {
		str = strconv.FormatFloat(po, 'f', 0, 64)
	} else {
		str = strconv.FormatFloat(po, 'e', 1, 64)
	}

	if neg {
		str = "-" + str // add "-" to format negative numbers
	}
	return str
}

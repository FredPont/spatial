package ui

import (
	"image/color"
	"spatial/src/filter"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

func buildPlot(plotMap map[string]filter.Dot) PlotBox {
	var x, y []float64
	//var items []string
	m := 1e308
	Xmax := -m
	Xmin := m
	Ymax := -m
	Ymin := m
	for _, val := range plotMap {
		x = append(x, val.X)
		y = append(y, val.Y)
		if val.X > Xmax {
			Xmax = val.X
		}
		if val.X < Xmin {
			Xmin = val.X
		}
		if val.Y > Ymax {
			Ymax = val.Y
		}
		if val.Y < Ymin {
			Ymin = val.Y
		}
	}
	return PlotBox{
		Title: "title",
		//id:     items,
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

// Scatter makes a scatter chart
func (p *PlotBox) scatterPlot(v *Interactive2Dsurf, dotsize int) {

	for i, xplot := range p.X {
		//x := int(MapRange(x, p.Xmin, p.Xmax, p.Left, 800-p.Right))
		//y := int(MapRange(p.Y[i], p.Ymin, p.Ymax, p.Bottom, 800-p.Top))
		x, y := xCoord(p, xplot), yCoord(p, p.Y[i])
		v.drawcircleScattCont(x, y, dotsize, color.NRGBA{128, 128, 128, 255})
		//log.Println(x, y)
	}

}

// XAxis makes the X axis
func (p *PlotBox) xAxisScat(v *Interactive2Dsurf) {
	y1 := 0
	lef := xCoord(p, p.Xmin)
	rig := xCoord(p, p.Xmax)

	//if yZero(p) {
	//	y1 = yCoord(p, 0.)
	//} else {
	y1 = yCoord(p, p.Ymin)
	//	log.Println("Y axis does not contain 0 value !")
	//}
	//log.Println("y axis:", x1, bot, up)
	c := iLine(lef, y1, rig, y1, 1, color.RGBA{0, 0, 0, 255})
	v.scatterContainer.Add(c) // add the line to the cluster container
	p.xScatlabel(v, y1)
	//v.scatterContainer.Refresh()
}

// YAxis makes the Y axis
func (p *PlotBox) yAxisScat(v *Interactive2Dsurf) {
	x1 := 0
	bot := yCoord(p, p.Ymin)
	up := yCoord(p, p.Ymax)

	// if xZero(p) {
	// 	x1 = xCoord(p, 0.) // if x=0 exists, the Y axes is plot on 0.
	// } else {
	// 	if p.Xmax < 0 { // if all x values < 0, the y axes is on the right
	x1 = xCoord(p, p.Xmax)
	// 	} else {
	// 		x1 = xCoord(p, p.Xmin)
	// 	}
	// }
	//log.Println("y axis:", x1, bot, up)
	c := iLine(x1, bot, x1, up, 1, color.RGBA{0, 0, 0, 255})
	v.scatterContainer.Add(c) // add the line to the cluster container
	p.yScatlabel(v, x1)
	//v.scatterContainer.Refresh()
}

// Xlabel makes the x axis scale text
func (p *PlotBox) xScatlabel(v *Interactive2Dsurf, y int) {
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
	prefs := fyne.CurrentApp().Preferences()
	// x coordinates of the 2D plot
	xplot := binding.BindPreferenceString("2DxPlot", prefs) // set the link to preferences for rotation
	xp, _ := xplot.Get()
	AbsText(v.scatterContainer, xCoord(p, (p.Xmax+p.Xmin)/2), y+35, xp, 12, color.NRGBA{0, 0, 0, 255}) // axis title
}

// Ylabel makes the x axis scale text
func (p *PlotBox) yScatlabel(v *Interactive2Dsurf, x int) {
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
	prefs := fyne.CurrentApp().Preferences()
	// y coordinates of the 2D plot
	yplot := binding.BindPreferenceString("2DyPlot", prefs) // set the link to preferences for rotation
	yp, _ := yplot.Get()

	AbsText(v.scatterContainer, x-35, yCoord(p, p.Ymax)-25, yp, 12, color.NRGBA{0, 0, 0, 255}) // axis title
}

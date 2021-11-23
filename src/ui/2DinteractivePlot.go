package ui

import (
	"image/color"
	"log"
	"spatial/src/filter"
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
func (p *PlotBox) scatterPlot(v *plotRaster, dotsize int) {

	for i, xplot := range p.X {
		//x := int(MapRange(x, p.Xmin, p.Xmax, p.Left, 800-p.Right))
		//y := int(MapRange(p.Y[i], p.Ymin, p.Ymax, p.Bottom, 800-p.Top))
		x, y := xCoord(p, xplot), yCoord(p, p.Y[i])
		v.drawcircleScattCont(x, y, dotsize, color.NRGBA{128, 128, 128, 255})
		log.Println(x, y)
	}

}

package ui

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2"
)

// PlotBox holds the essential data for making a chart
type PlotBox struct {
	Title                    string
	id                       []string
	X                        []float64
	Y                        []float64
	Color                    color.RGBA
	Top, Bottom, Left, Right float64
	Xmax, Xmin, Ymax, Ymin   float64
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
	}
}

// MapRange -- given a value between low1 and high1, return the corresponding value between low2 and high2
// credits : https://github.com/ajstarks/fc
func MapRange(value, low1, high1, low2, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

// buildVulanoPlot : create the window and the vulcano plot
func buildVulanoPlot(fname string, pvfcTable []PVrecord) {
	readVulcano(fname, pvfcTable)
	log.Println(readVulcano(fname, pvfcTable))
	v := buildVulcWin()
	drawVulcano(v)
}

// drawline a circle at x,y position to the scatter container
func (e *Vulcano) drawcircle(x, y, ray int, color color.NRGBA) fyne.CanvasObject {
	c := iCircle(x, y, ray, color)  // draw circle rayon ray
	e.scatterContainer.AddObject(c) // add the cicle to the cluster container
	return c
}

func drawVulcano(e *Vulcano) {
	R := uint8(50)
	G := uint8(150)
	B := uint8(250)
	for i := 0; i < 700; i += 10 {
		e.drawcircle(i, i, 3, color.NRGBA{R, G, B, 255})
	}
	e.scatterContainer.Refresh()
}

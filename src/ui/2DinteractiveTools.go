/*
 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.

 Written by Frederic PONT.
 (c) Frederic Pont 2021
*/

package ui

import (
	"image/color"
	"log"
	"math"
	"spatial/src/filter"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/mazznoer/colorgrad"
)

// show2D show 2Di tools and 2Di window
func show2D(a fyne.App, e *Editor, preference fyne.Preferences, f binding.Float, header []string, firstTable string) {
	f.Set(0.3)
	winplot, inter2D := build2DplotWin(e) // show 2D interactive window

	show2DinterTools(a, e, winplot, inter2D, preference, f, header, firstTable) // show tool box

	f.Set(0.)
}

// show2DinterTools show 2Di tools
func show2DinterTools(a fyne.App, e *Editor, winplot fyne.Window, inter2D *Interactive2Dsurf, preference fyne.Preferences, f binding.Float, header []string, firstTable string) {
	plotbox, dotmap, imageMap := build2DPlot(inter2D, preference, header, firstTable) // build scatter plot
	win2Dtools := a.NewWindow("2D plot tools")

	gatename := widget.NewEntry()
	gatename.SetPlaceHolder("Selection name...")

	content := container.NewVBox(
		gatename,
		widget.NewButton("Show Cells in Gates", func() {
			go searchDotsInGates(e, inter2D, &plotbox, dotmap, imageMap, f)
		}),
		widget.NewButton("Save Gates", func() {
			go save2DGates(gatename.Text, e, inter2D)
		}),
		widget.NewButton("Clear Gates", func() {
			go init2DScatterGates(inter2D)
		}),
		// screenshot
		widget.NewButtonWithIcon("", theme.MediaPhotoIcon(), func() {
			go screenShot(winplot, gatename.Text, f)
		}),
		widget.NewButtonWithIcon("Exit", theme.LogoutIcon(), func() {
			win2Dtools.Close() // close tool window
			winplot.Close()    // close plot window
		}),
	)
	win2Dtools.SetContent(content)
	//win2D.Resize(fyne.Size{Width: 500, Height: 500})
	win2Dtools.Show()
}

// build2DPlot start extracting the plot data and make the plot
func build2DPlot(inter2D *Interactive2Dsurf, prefs fyne.Preferences, header []string, firstTable string) (PlotBox, map[string]filter.Dot, map[string]filter.Point) {
	subtable := extract2DinterData(prefs, header, firstTable)
	imageMap, plotMap := subTableToMap(subtable)

	plotbox := buildPlot(plotMap)
	//get scatter dot size
	ds := binding.BindPreferenceString("2Ddotsize", prefs) // set the link to 2D dot size preferences
	ds2 := binding.StringToInt(ds)
	dotsize, _ := ds2.Get()

	// built scatter plot
	plotbox.scatterPlot(inter2D, dotsize)
	plotbox.xAxisScat(inter2D)
	plotbox.yAxisScat(inter2D)
	inter2D.scatterContainer.Refresh()
	return plotbox, plotMap, imageMap
}

// extract cols index from the first data table :
// cells ID
// x,y coordinates of the microcopy image
// x,y coordinates of the 2D scatter plot
func plotColIndex(prefs fyne.Preferences, header []string) []int {
	// X coordinates of the microcopy image
	xMic := binding.BindPreferenceString("xcor", prefs) // set the link to preferences for x coordinates
	xMi, _ := xMic.Get()
	// y coordinates
	yMic := binding.BindPreferenceString("ycor", prefs) // set the link to preferences for y coordinates
	yMi, _ := yMic.Get()

	// x coordinates of the 2D plot
	xplot := binding.BindPreferenceString("2DxPlot", prefs) // set the link to preferences for rotation
	xp, _ := xplot.Get()

	// y coordinates of the 2D plot
	yplot := binding.BindPreferenceString("2DyPlot", prefs) // set the link to preferences for rotation
	yp, _ := yplot.Get()

	list := []string{xMi, yMi, xp, yp}

	colIndexes := []int{0} // 0 = get first column = cell names
	colIndexes = append(colIndexes, filter.GetColIndex(header, list)...)
	//ReadColumns(filename , colIndexes )
	return colIndexes
}

// extract from the first table :
// cells ID
// x,y coordinates of the microcopy image
// x,y coordinates of the 2D scatter plot
func extract2DinterData(prefs fyne.Preferences, header []string, firstTable string) [][]string {

	colIndexes := plotColIndex(prefs, header)
	cols := filter.ReadColumns(firstTable, colIndexes)

	return cols
}

// convert the plot subtable (cells ID, x,y coordinates of the microcopy image, x,y coordinates of the 2D scatter plot)
// into 2 maps : cellID -> []Point (microcopy) cellID -> []Dot (plot)
func subTableToMap(subtable [][]string) (map[string]filter.Point, map[string]filter.Dot) {
	l := len(subtable)
	imageMap := make(map[string]filter.Point, l)
	plotMap := make(map[string]filter.Dot, l)

	for i := 0; i < l; i++ {
		id := subtable[i][0] //cell names

		imx := filter.StrToInt(subtable[i][1]) // x microscopy
		imy := filter.StrToInt(subtable[i][2]) // x microscopy
		px := filter.StrToF64(subtable[i][3])  // x plot
		py := filter.StrToF64(subtable[i][4])  // x plot

		imageMap[id] = filter.Point{X: imx, Y: imy}
		plotMap[id] = filter.Dot{X: px, Y: py}
	}

	return imageMap, plotMap
}

// convert the cellID -> []Dot (plot) to cellID -> Point
// apply to each dot the conversion to the pixel position in the scatter window
func dotMapToPointMap(p *PlotBox, dotmap map[string]filter.Dot) map[string]filter.Point {
	pointMap := make(map[string]filter.Point, len(dotmap))

	for k, v := range dotmap {
		pointMap[k] = filter.Point{X: xCoord(p, v.X), Y: yCoord(p, v.Y)}
	}
	return pointMap
}

// convert the scatter points to dots position in pixel
// filter the dots that are in the gates
// show the dots in gate in the microscopy image
func searchDotsInGates(e *Editor, inter2D *Interactive2Dsurf, p *PlotBox, dotmap map[string]filter.Dot, imagemap map[string]filter.Point, f binding.Float) {
	f.Set(0.3)
	scatter := dotMapToPointMap(p, dotmap)
	cellsInGates := selectedCells(inter2D, scatter)

	go plotDotsMicrocop(e, cellsInGates, imagemap)
	plotDotsInGates(p, inter2D, cellsInGates)

	inter2D.gateContainer.Refresh()
	f.Set(0.)
}

// extract the cells (map cell ID => XY) in the gates drawn in the 2D plot
func selectedCells(inter2D *Interactive2Dsurf, scatter map[string]filter.Point) []map[string]filter.Point {
	cellsInGates := make([]map[string]filter.Point, 0)
	for _, gate := range inter2D.drawSurface.alledges {
		cells := filter.DotsInGate(gate, scatter)
		cellsInGates = append(cellsInGates, cells)
	}
	return cellsInGates
}

// plot the dots in gates in color in the 2D scatter plot. dots are plotted in the gate container
func plotDotsInGates(p *PlotBox, inter2D *Interactive2Dsurf, cellsInGates []map[string]filter.Point) {
	prefs := fyne.CurrentApp().Preferences()
	//get scatter dot size
	ds := binding.BindPreferenceString("2Ddotsize", prefs) // set the link to 2D dot size preferences
	ds2 := binding.StringToInt(ds)
	dotsize, _ := ds2.Get()

	nbGates := len(cellsInGates)
	for i := 0; i < nbGates; i++ {
		dotcolor := dotColors(nbGates, i)
		p.gatesDotPlot(inter2D, dotsize, cellsInGates[i], dotcolor)
	}
}

// dotColors computes the color of scatter dots
// for a total number of clusters "nbGates"
func dotColors(nbGates, gateIndex int) color.NRGBA {
	grad := colorgrad.Rainbow().Sharp(uint(nbGates+1), 0.2)
	return nrgbaModel(grad.Colors(uint(nbGates + 1))[gateIndex])
}

func nrgbaModel(c color.Color) color.NRGBA {
	r, g, b, a := c.RGBA()
	return color.NRGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

// extract the cells ID of the cells in gates and get their corresponding XY coordinates for the microscopy image
func plotDotsMicrocop(e *Editor, cellsInGates []map[string]filter.Point, imageMap map[string]filter.Point) {
	initCluster(e) // remove all dots of the cluster container
	nbGates := len(cellsInGates)
	// get the image microscopy coordinates of the cells in one gate from the cells names in the 2D plot
	for i, cellsingate := range cellsInGates {
		var cellsXY []filter.Point
		dotcolor := dotColors(nbGates, i)
		for cellid := range cellsingate {
			cellsXY = append(cellsXY, imageMap[cellid])
		}
		drawCells(e, cellsXY, dotcolor)
	}

}

// draw the selected cells on the microscopy image
func drawCells(e *Editor, cellsXY []filter.Point, dotcolor color.NRGBA) {

	pref := fyne.CurrentApp().Preferences()
	clustOp := binding.BindPreferenceFloat("clustOpacity", pref) // cluster opacity
	opacity, _ := clustOp.Get()
	op := uint8(opacity)
	clustDia := binding.BindPreferenceInt("clustDotDiam", pref) // cluster dot diameter
	diameter, _ := clustDia.Get()
	diameter = ApplyZoomInt(e, diameter)
	sf := binding.BindPreferenceFloat("scaleFactor", pref) // set the link to preferences for scaling factor
	scaleFactor, _ := sf.Get()                             // read the preference for scaling factor
	rot := binding.BindPreferenceBool("rotate", pref)      // set the link to preferences for rotation
	rotate, _ := rot.Get()

	for _, xy := range cellsXY {
		xScaled, yScaled := scale(xy.X, xy.Y, scaleFactor, rotate)
		e.drawcircle(ApplyZoomInt(e, xScaled), ApplyZoomInt(e, yScaled), diameter, color.NRGBA{dotcolor.R, dotcolor.G, dotcolor.B, op})
		//log.Println(xy)
	}

	e.clusterContainer.Refresh()
}

// apply the scaling factor and rotation to xy coordinates
func scale(x, y int, scaleFactor float64, rotate bool) (int, int) {

	xScaled := int(math.Round(float64(x) * scaleFactor))
	yScaled := int(math.Round(float64(y) * scaleFactor))

	if rotate {
		xRot := yScaled
		yRot := xScaled
		return xRot, yRot
	}
	return xScaled, yScaled

}

// save the gates to csv files withe ImageJ format and 100% zoom
// X,Y
// 131,150
// 105,189
// 156,187
func save2DGates(gateName string, e *Editor, inter2D *Interactive2Dsurf) {

	gateName = filter.FormatOutFile("gate", gateName, "") // test if name exist, if not, build a file name with the current time

	zoomFactor := 100. / float64(e.zoom)
	for i, poly := range inter2D.drawSurface.alledges {
		if len(poly) < 3 {
			continue
		}

		out := strconv.Itoa(i) + "_" + gateName
		writeCSV(out, filter.ZoomPolygon(poly, zoomFactor))
		log.Println("gate saved in gates/", out)
	}
}

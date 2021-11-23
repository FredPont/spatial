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
	"spatial/src/filter"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// show2D show 2Di tools and 2Di window
func show2D(a fyne.App, e *Editor, preference fyne.Preferences, f binding.Float, header []string, firstTable string) {
	f.Set(0.3)
	winplot, inter2D := build2DplotWin(e) // show 2D interactive window
	show2DinterTools(a, e, winplot, inter2D, preference, f, header, firstTable)
	build2DPlot(inter2D, preference, header, firstTable)
	f.Set(0.)
}

// show2DinterTools show 2Di tools
func show2DinterTools(a fyne.App, e *Editor, winplot fyne.Window, inter2D *Interactive2Dsurf, preference fyne.Preferences, f binding.Float, header []string, firstTable string) {

	win2Dtools := a.NewWindow("2D plot tools")

	content := container.NewVBox(
		widget.NewLabel("Tools"),
		widget.NewButton("Show Cells in Gates", func() {

			//go
		}),
		// screenshot
		widget.NewButtonWithIcon("", theme.MediaPhotoIcon(), func() {
			//go screenShot(w, gatename.Text, f)
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

// build2DPlo start extracting the plot data and make the plot
func build2DPlot(inter2D *Interactive2Dsurf, prefs fyne.Preferences, header []string, firstTable string) {
	subtable := extract2DinterData(prefs, header, firstTable)
	_, plotMap := subTableToMap(subtable)
	plotbox := buildPlot(plotMap)
	//get scatter dot size
	ds := binding.BindPreferenceString("2Ddotsize", prefs) // set the link to 2D dot size preferences
	ds2 := binding.StringToInt(ds)
	dotsize, _ := ds2.Get()
	plotbox.scatterPlot(inter2D, dotsize)
	plotbox.xAxisScat(inter2D)
	plotbox.yAxisScat(inter2D)
	inter2D.scatterContainer.Refresh()
}

// extract cols index from the first table :
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

	colIndexes := []int{0} // get first column = cell names
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

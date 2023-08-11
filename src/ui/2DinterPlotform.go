package ui

import (
	"spatial/src/filter"
	"spatial/src/pref"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Plot2Dform display a form with the plot preferences and parameters
func Plot2Dform(a fyne.App, e *Editor, win fyne.Window, zoom int, header []string, firstTable string, alledges [][]filter.Point, f binding.Float) {
	prefs := a.Preferences()

	// plot name
	plotFileTitle := binding.BindPreferenceString("2DplotName", prefs) // set the link to preferences for plot name
	plotFT, _ := plotFileTitle.Get()
	plotName := widget.NewEntry()
	plotName.SetText(plotFT)

	// x coordinates
	xplot := binding.BindPreferenceString("2DxPlot", prefs) // set the link to preferences for x coordinates
	xp, _ := xplot.Get()
	//x := widget.NewSelectEntry(header)
	//x.SetText(xp)
	var x *widget.Button
	x = widget.NewButton(xp, func() {
		pref.ShowTable(header, xplot, x, xp)
	})

	// y coordinates
	yplot := binding.BindPreferenceString("2DyPlot", prefs) // set the link to preferences for y coordinates
	yp, _ := yplot.Get()
	//y := widget.NewSelectEntry(header)
	//y.SetText(yp)
	var y *widget.Button
	y = widget.NewButton(yp, func() {
		pref.ShowTable(header, yplot, y, yp)
	})

	// dot size
	dotsize := binding.BindPreferenceString("2Ddotsize", prefs) // set the link to preferences for dot size
	ds, _ := dotsize.Get()

	plotdot := widget.NewEntry()
	plotdot.SetText(ds)

	// dots color
	// gateDotscol := widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), func() { GateDotscolor(a, win) })
	// unselcol := widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), func() { unseldcolor(a, win) })

	dialog.ShowForm("Plot parameters", "Enter", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Plot Name", plotName),
			widget.NewFormItem("X", x),
			widget.NewFormItem("Y", y),
			// widget.NewFormItem("Dots in Gate color", gateDotscol),
			// widget.NewFormItem("Dots in Bkgd color", unselcol),
			widget.NewFormItem("dot size", plotdot)},
		func(input bool) {
			//log.Println("input = ", input)
			if input {
				//f.Set(0.3)
				xp, _ = xplot.Get()
				yp, _ = yplot.Get()
				//go makeplot(a, zoom, header, firstTable, xp, yp, plotName.Text, plotdot.Text, alledges, f)
				go show2D(a, e, prefs, f, header, firstTable)
				go save2DPlotPrefs(a, xp, yp, plotName.Text, plotdot.Text)
			}

		}, win)

}

// save preference of the plot
func save2DPlotPrefs(a fyne.App, x, y, plotName, dotsize string) {
	pref := a.Preferences()
	// X coordinates
	pref.SetString("2DxPlot", x)

	// Y coordinates
	pref.SetString("2DyPlot", y)

	// plotName
	pref.SetString("2DplotName", plotName)

	//dot size
	pref.SetString("2Ddotsize", dotsize)
}

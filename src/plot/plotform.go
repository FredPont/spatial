package plot

import (
	"image/color"
	"log"
	"spatial/src/filter"
	"spatial/src/pref"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Plotform display a form with the plot preferences and parameters
func Plotform(a fyne.App, win fyne.Window, zoom int, header []string, firstTable string, alledges [][]filter.Point, f binding.Float) {
	prefs := a.Preferences()

	// plot name
	plotFileTitle := binding.BindPreferenceString("plotName", prefs) // set the link to preferences for rotation
	plotFT, _ := plotFileTitle.Get()
	plotName := widget.NewEntry()
	plotName.SetText(plotFT)

	// x coordinates
	xplot := binding.BindPreferenceString("xPlot", prefs) // set the link to preferences for rotation
	xp, _ := xplot.Get()
	//x := widget.NewSelectEntry(header)
	//x.SetText(xp)
	var x *widget.Button
	x = widget.NewButton(xp, func() {
		pref.ShowTable(header, xplot, x, xp)
	})

	// y coordinates
	yplot := binding.BindPreferenceString("yPlot", prefs) // set the link to preferences for rotation
	yp, _ := yplot.Get()
	//y := widget.NewSelectEntry(header)
	//y.SetText(yp)
	var y *widget.Button
	y = widget.NewButton(yp, func() {
		pref.ShowTable(header, yplot, y, yp)
	})

	// dot size
	dotsize := binding.BindPreferenceString("dotsize", prefs) // set the link to preferences for rotation
	ds, _ := dotsize.Get()
	plotdot := widget.NewEntry()
	plotdot.SetText(ds)

	// dots color
	gateDotscol := widget.NewButton("color", func() { GateDotscolor(a, win) })
	unselcol := widget.NewButton("color", func() { unseldcolor(a, win) })

	dialog.ShowForm("Plot parameters", "Enter", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Plot Name", plotName),
			widget.NewFormItem("X", x),
			widget.NewFormItem("Y", y),
			widget.NewFormItem("Dots in Gate color", gateDotscol),
			widget.NewFormItem("Dots in Bkgd color", unselcol),
			widget.NewFormItem("dot size", plotdot)},
		func(input bool) {
			//log.Println("input = ", input)
			if input {
				//f.Set(0.3)
				xp, _ = xplot.Get()
				yp, _ = yplot.Get()
				go makeplot(a, zoom, header, firstTable, xp, yp, plotName.Text, plotdot.Text, alledges, f)
				go savePlotPrefs(a, xp, yp, plotName.Text, plotdot.Text)
			}

		}, win)

}

// save preference of the plot
func savePlotPrefs(a fyne.App, x, y, plotName, dotsize string) {
	pref := a.Preferences()
	// X coordinates
	pref.SetString("xPlot", x)

	// Y coordinates
	pref.SetString("yPlot", y)

	// plotName
	pref.SetString("plotName", plotName)

	//dot size
	pref.SetString("dotsize", dotsize)
}

// color picker for the plot background (unselected) dots color
func unseldcolor(a fyne.App, win fyne.Window) {
	pref := a.Preferences()

	picker := dialog.NewColorPicker("Pick a Color", "What is your favorite color?", func(c color.Color) {
		log.Println("Color picked:", c)
		R, G, B, A := ColorToRGBA(c)
		log.Println("Color RGBA picked:", R, G, B, A)
		pref.SetInt("unselR", R)
		pref.SetInt("unselG", G)
		pref.SetInt("unselB", B)
		pref.SetInt("unselA", A)
	},
		win)
	picker.Advanced = true
	picker.Show()
}

// GateDotscolor color picker for the plot dots color in gate
func GateDotscolor(a fyne.App, win fyne.Window) {
	pref := a.Preferences()

	picker := dialog.NewColorPicker("Pick a Color", "What is your favorite color?", func(c color.Color) {
		log.Println("Color picked:", c)
		R, G, B, A := ColorToRGBA(c)
		log.Println("Color RGBA picked:", R, G, B, A)
		pref.SetInt("gateDotsR", R)
		pref.SetInt("gateDotsG", G)
		pref.SetInt("gateDotsB", B)
		pref.SetInt("gateDotsA", A)
	},
		win)
	picker.Advanced = true
	picker.Show()
}

// ColorToRGBA return r,g,b,a for a color
// credits : fyne/dialog/color.go
func ColorToRGBA(c color.Color) (r, g, b, a int) {
	switch col := c.(type) {
	case color.NRGBA:
		r = int(col.R)
		g = int(col.G)
		b = int(col.B)
		a = int(col.A)
	case *color.NRGBA:
		r = int(col.R)
		g = int(col.G)
		b = int(col.B)
		a = int(col.A)
	default:
		r, g, b, a = unmultiplyAlpha(c)
	}
	return
}

func unmultiplyAlpha(c color.Color) (r, g, b, a int) {
	red, green, blue, alpha := c.RGBA()
	if alpha != 0 && alpha != 0xffff {
		ratio := float64(alpha) / 0xffff
		red = uint32(float64(red) / ratio)
		green = uint32(float64(green) / ratio)
		blue = uint32(float64(blue) / ratio)
	}
	// Convert from range 0-65535 to range 0-255
	r = int(red >> 8)
	g = int(green >> 8)
	b = int(blue >> 8)
	a = int(alpha >> 8)
	return
}

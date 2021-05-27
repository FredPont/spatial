package plot

import (
	"fmt"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Plotform display a form with the plot preferences and parameters
func Plotform(a fyne.App, win fyne.Window, header []string) {

	// x coordinates
	x := widget.NewSelectEntry(header)
	// y coordinates
	y := widget.NewSelectEntry(header)
	// dot size
	plotdot := widget.NewEntry()
	unselcol := widget.NewButton("color", func() { unseldcolor(a, win) })

	dialog.ShowForm("Form Input", "Enter", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("X", x),
			widget.NewFormItem("Y", y),
			widget.NewFormItem("col", unselcol),
			widget.NewFormItem("dot size", plotdot)},
		func(bool) { fmt.Println("Selected", unselcol) }, win)
}

// color picker for the plot background (unselected) dots color
func unseldcolor(a fyne.App, win fyne.Window) {
	pref := a.Preferences()

	picker := dialog.NewColorPicker("Pick a Color", "What is your favorite color?", func(c color.Color) {
		log.Println("Color picked:", c)
		r, g, b, a := colorToRGBA(c)
		pref.SetInt("unselR", r)
		pref.SetInt("unselG", g)
		pref.SetInt("unselB", b)
		pref.SetInt("unselA", a)
	},
		win)
	picker.Advanced = true
	picker.Show()
}

// credits : fyne/dialog/color.go
func colorToRGBA(c color.Color) (r, g, b, a int) {
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

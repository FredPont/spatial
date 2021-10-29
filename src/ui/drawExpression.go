package ui

import (
	"image/color"
	"lasso/src/filter"
	"lasso/src/plot"
	"lasso/src/pref"
	"log"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func buttonDrawExpress(a fyne.App, e *Editor, preference fyne.Preferences, f binding.Float, header []string, firstTable string) {
	ExpressWindow := a.NewWindow("Expression")

	// select the expression to draw
	// expSel := widget.NewSelectEntry(header) // limited for long list and replaced by custom select
	sel := binding.NewString()
	userSel, _ := sel.Get()
	var expSel *widget.Button
	expSel = widget.NewButton(userSel, func() {
		pref.ShowTable(header[1:], sel, expSel, "Selection")
	})

	// show choice of different gradien
	grad := widget.NewRadioGroup([]string{"Turbo", "Viridis", "Inferno", "White - Red", "Yellow - Red", "Purple - Red", "Blue - Yellow - Red"}, func(s string) {

		//fmt.Println("Selected <", s, ">")
	})

	// Dot opacity
	DotOp := binding.BindPreferenceFloat("dotOpacity", preference) // pref binding for the expression dot opacity
	DotOpacity := widget.NewSliderWithData(0, 255, DotOp)
	DotOpacity.Step = 1
	DotOpacity.OnChanged = func(v float64) {
		preference.SetFloat("dotOpacity", v)
	}

	//legend color - the results is store in preferences
	legendcol := widget.NewButton("Legend Text Color", func() { LegendTxtscolor(a, ExpressWindow) })

	//animation
	anim := binding.NewBool()           // if true start animation
	curPathwayIndex := binding.NewInt() // column index of current pathway displayed by slide show
	curPathwayIndex.Set(1)              // start with column 1 by default
	slideDelay := binding.NewFloat()    // default pause between slides
	slideDelay.Set(1)                   // 1 sec pause
	slidePause := widget.NewEntry()
	slidePause.SetPlaceHolder("Pause between slides (sec)")

	content := container.NewVBox(
		widget.NewLabel("Select the variable"),
		expSel,
		widget.NewLabel("Select your gradient"),
		grad,
		widget.NewLabel("Dots Opacity [0-100%] :"),
		DotOpacity,
		legendcol,
		widget.NewButton("Plot Expression", func() {
			userSel, _ = sel.Get()
			// gradien default
			def := "White - Red"
			if grad.Selected == "" {
				grad.Selected = def
			}
			if userSel == "" {
				return // return if nothing is selected
			}
			//log.Println(expSel.Entry.Text, grad.Selected, op)
			go drawExp(a, e, header, firstTable, userSel, grad.Selected, f, curPathwayIndex)
		}),
		slidePause,
		widget.NewButton("Slide show", func() {
			anim.Set(true)
			setDelay(slidePause, slideDelay)
			go startSlideShow(a, e, header, firstTable, grad.Selected, f, anim, curPathwayIndex, slideDelay)

		}),
		widget.NewButton("Stop/Continue Slide show", func() {
			setDelay(slidePause, slideDelay)
			an, _ := anim.Get()
			if an {
				anim.Set(false)
			} else {
				anim.Set(true)
				go startSlideShow(a, e, header, firstTable, grad.Selected, f, anim, curPathwayIndex, slideDelay)
			}

		}),
		widget.NewButton("Previous Slide", func() {
			startIdx, _ := curPathwayIndex.Get()
			if startIdx < 2 {
				log.Println("Column", startIdx-1, "cannot be accessed !")
				return
			}
			curPathwayIndex.Set(startIdx - 1)
			go drawExp(a, e, header, firstTable, header[startIdx-1], grad.Selected, f, curPathwayIndex)
		}),
		widget.NewButton("Next Slide", func() {
			startIdx, _ := curPathwayIndex.Get()
			if startIdx > len(header)-2 {
				log.Println("Data table have", len(header)-1, "columns ! Column", startIdx+1, "cannot be accessed !")
				return
			}
			curPathwayIndex.Set(startIdx + 1)
			go drawExp(a, e, header, firstTable, header[startIdx+1], grad.Selected, f, curPathwayIndex)
		}),
		widget.NewButton("Close", func() { ExpressWindow.Close() }),
	)
	ExpressWindow.SetContent(content)
	ExpressWindow.Show()
}

func startSlideShow(a fyne.App, e *Editor, header []string, firstTable, grad string, f binding.Float, anim binding.Bool, curPathwayIndex binding.Int, slideDelay binding.Float) {
	m := len(header)
	startIdx, _ := curPathwayIndex.Get()
	for i := startIdx; i < m; i++ {
		// listen to anim
		an, _ := anim.Get()
		if !an {
			curPathwayIndex.Set(i)
			break
		}
		if header[i] != "" {
			drawExp(a, e, header, firstTable, header[i], grad, f, curPathwayIndex)
			pause, _ := slideDelay.Get()
			time.Sleep(time.Duration(1000*pause) * time.Millisecond)
			log.Println(i, "/", m-1, " column :", header[i])
		}

	}

}

// set delay between slides
func setDelay(slidePause *widget.Entry, slideDelay binding.Float) {
	if slidePause != nil {
		p, err := strconv.ParseFloat(slidePause.Text, 64)
		if err != nil {
			log.Println(p, "delay cannot be converted to float !")
		} else {
			slideDelay.Set(p)
		}
	}

}

func getExpress(a fyne.App, header []string, filename string, expcol string, curPathwayIndex binding.Int) ([]float64, []filter.Point) {
	pref := a.Preferences()
	// X coordinates
	xcor := binding.BindPreferenceString("xcor", pref) // set the link to preferences for rotation
	xc, _ := xcor.Get()

	// y coordinates
	ycor := binding.BindPreferenceString("ycor", pref) // set the link to preferences for rotation
	yc, _ := ycor.Get()

	colIndexes := filter.GetColIndex(header, []string{expcol, xc, yc})
	curPathwayIndex.Set(colIndexes[0]) // set the current expression Index to the selected column to enable button next/previous slide
	return filter.ReadExpress(a, filename, colIndexes)
}

func drawExp(a fyne.App, e *Editor, header []string, filename string, expcol, gradien string, f binding.Float, curPathwayIndex binding.Int) {
	f.Set(0.2)     // progress bar set to 20%
	initCluster(e) // remove all dots of the cluster container
	pref := a.Preferences()
	// Dot opacity
	DotOp := binding.BindPreferenceFloat("dotOpacity", pref) // pref binding for the  dot opacity
	opacity, _ := DotOp.Get()
	op := uint8(opacity)
	clustDia := binding.BindPreferenceInt("clustDotDiam", pref) //  dot diameter
	diameter, _ := clustDia.Get()
	diameter = ApplyZoomInt(e, diameter)

	expressions, pts := getExpress(a, header, filename, expcol, curPathwayIndex) // []expressions and []Point
	if len(expressions) < 1 {
		log.Println("Intensities not availble for column", expcol)
		return
	}
	f.Set(0.3) // progress bar set to 30% after data reading
	scaleExp, min, max := filter.ScaleSlice01(expressions)

	//legendPosition := filter.Point{X: 15, Y: 15} // initial legend position for cluster names
	nbPts := len(pts)

	for c := 0; c < nbPts; c++ {
		// progress bar increases when 50% of points are loaded
		if c == int(nbPts/2) {
			f.Set(0.5) // 50 % progression for progress bar
		}

		clcolor := grad(gradien)(scaleExp[c])

		e.drawcircle(ApplyZoomInt(e, pts[c].X), ApplyZoomInt(e, pts[c].Y), diameter, color.NRGBA{clcolor.R, clcolor.G, clcolor.B, op})

	}
	// draw legend titel, dot and value for the current cexpression
	R, G, B, _ := plot.GetPrefColorRGBA(a, "legendColR", "legendColG", "legendColB", "legendColA")
	colorText := color.NRGBA{uint8(R), uint8(G), uint8(B), 255}
	titleLegend(e, expcol, colorText)
	expLegend(e, op, diameter, gradien, min, max, colorText)

	e.clusterContainer.Refresh()
	f.Set(0.) // reset progress bar
}

// print pathway name on top of image
func titleLegend(e *Editor, title string, c color.NRGBA) {
	AbsText(e.clusterContainer, 50, 30, title, 20, c)
}

// draw expression legend with dots and values
func expLegend(e *Editor, op uint8, diameter int, gradien string, min, max float64, c color.NRGBA) {
	x, y := 13, 30
	sp := 25
	//AbsText(e.clusterContainer, x+20, y+10, "toto", 20, color.NRGBA{50, 50, 50, 255})
	for i := 5; i >= 0; i-- {
		//exp := fmt.Sprintf("%.1f", unscale(float64(i)/5., min, max))
		exp := TicksDecimals(unscale(float64(i)/5., min, max))
		AbsText(e.clusterContainer, x+20, y+155-sp*i, exp, 15, c)
		co := grad(gradien)(float64(i) / 5.)
		e.drawcircle(x, y+150-sp*i, diameter*100/e.zoom, color.NRGBA{co.R, co.G, co.B, op})
	}
}

// compute the true expression value from a scaled [0-1] value
func unscale(v, min, max float64) float64 {
	return v*(max-min) + min
}

// grad return the gradien function with name "gradien"
func grad(gradien string) func(float64) RGB {
	switch gradien {
	case "Turbo":
		return func(val float64) RGB { return TurboGradien(val) }
	case "Viridis":
		return func(val float64) RGB { return ViridisGrad(val) }
	case "White - Red":
		return func(val float64) RGB { return WRgradien(val) }
	case "Yellow - Red":
		return func(val float64) RGB { return YlRdGradien(val) }
	case "Purple - Red":
		return func(val float64) RGB { return PuRdGradien(val) }
	case "Inferno":
		return func(val float64) RGB { return InferGrad(val) }
	case "Blue - Yellow - Red":
		return func(val float64) RGB { return BYRGradien(val) }
	default:
		return func(val float64) RGB { return WRgradien(val) }
	}

}

// LegendTxtscolor color picker for the legend text color
func LegendTxtscolor(a fyne.App, win fyne.Window) {
	pref := a.Preferences()

	picker := dialog.NewColorPicker("Pick a Color", "What is your favorite color?", func(c color.Color) {
		log.Println("Color picked:", c)
		R, G, B, A := plot.ColorToRGBA(c)
		log.Println("Color RGBA picked:", R, G, B, A)
		pref.SetInt("legendColR", R)
		pref.SetInt("legendColG", G)
		pref.SetInt("legendColB", B)
		pref.SetInt("legendColA", A)
	},
		win)
	picker.Advanced = true
	picker.Show()
}

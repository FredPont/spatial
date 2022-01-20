package ui

import (
	"image/color"
	"log"
	"spatial/src/filter"
	"spatial/src/plot"
	"spatial/src/pref"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func buttonDrawExpress(a fyne.App, e *Editor, preference fyne.Preferences, f binding.Float, header []string, firstTable string) {
	ExpressWindow := a.NewWindow("Expression")

	// show choice of different gradien
	grad := widget.NewRadioGroup([]string{"Turbo", "Viridis", "Inferno", "White - Red", "Yellow - Red", "Purple - Red", "Blue - Yellow - Red"}, func(s string) {
		//fmt.Println("Selected <", s, ">")
	})

	// Dot opacity
	//DotOp := binding.BindPreferenceFloat("dotOpacity", preference) // pref binding for the expression dot opacity
	//DotOp.Set(255.)
	DotOpacity := widget.NewSlider(0., 255.)
	//DotOp := binding.BindPreferenceFloat("dotOpacity", preference) // pref binding for the expression dot opacity
	//DotOp.Set(255.)
	//DotOpacity := widget.NewSliderWithData(0., 255., DotOp)
	//DotOpacity.Step = 1.
	DotOpacity.Value = 255.
	DotOpacity.OnChanged = func(v float64) {
		preference.SetFloat("dotOpacity", v)
	}

	// opacity gradient : checkbox to enable a gradient of opacity based on expression values
	opacityGradient := widget.NewCheck("Opacity gradient", func(v bool) {
		preference.SetBool("gradOpacity", v)
		//log.Println("Check set to", v)
	})
	// max threshold of the opacity gradient
	gradMaxWdgt := widget.NewEntry()
	gradMaxWdgt.OnChanged = func(str string) {
		v, err := strconv.ParseFloat(str, 64)
		if err != nil {
			log.Println(str, "is not a number !")
			return
		}
		preference.SetFloat("gradMax", v)
	}
	// min threshold of the opacity gradient
	gradMinWdgt := widget.NewEntry()
	gradMinWdgt.OnChanged = func(str string) {
		v, err := strconv.ParseFloat(str, 64)
		if err != nil {
			log.Println(str, "is not a number !")
			return
		}
		preference.SetFloat("gradMin", v)
	}

	// density plot
	initDensityPlot() // clear previous density plot
	densityPlot := plot.DensityPicture()

	// Max expression slider
	MaxExp := widget.NewSlider(0., 100.)
	MaxExp.Value = 100.
	//MaxExp.Step = (eMax - eMin) / 100.
	userMaxExp := binding.BindPreferenceFloat("userMaxExp", preference) // pref binding for user current the expression Max
	MaxExp.OnChanged = func(v float64) {
		preference.SetFloat("userMaxExp", v)
	}

	// Min expression slider
	MinExp := widget.NewSlider(0., 100.)
	MinExp.Value = 0.
	//MaxExp.Step = (eMax - eMin) / 100.
	userMinExp := binding.BindPreferenceFloat("userMinExp", preference) // pref binding for user current the expression Min
	MinExp.OnChanged = func(v float64) {
		preference.SetFloat("userMinExp", v)
	}

	// select the expression to draw
	// expSel := widget.NewSelectEntry(header) // limited for long list and replaced by custom select
	sel := binding.NewString()
	userSel, _ := sel.Get()
	var expSel *widget.Button
	expSel = widget.NewButton(userSel, func() {
		pref.ShowTable(header[1:], sel, expSel, "Selection")

	})

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
		container.NewHBox(
			widget.NewLabel("Select the variable and gradient :"),
			widget.NewButtonWithIcon("Close", theme.LogoutIcon(), func() { ExpressWindow.Close() }),
		),
		expSel,
		container.NewHBox(
			container.NewVBox(
				//widget.NewLabel("Select your gradient"),
				grad,
			),
			container.NewVBox(
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

					initSliderExp(MaxExp, MinExp) // reset expression min max sliders values

					go drawExp(a, e, header, firstTable, userSel, grad.Selected, f, curPathwayIndex, ExpressWindow)

				}),
				legendcol,

				slidePause,
				widget.NewButton("Slide show", func() {
					anim.Set(true)
					setDelay(slidePause, slideDelay)
					go startSlideShow(a, e, header, firstTable, grad.Selected, f, anim, curPathwayIndex, slideDelay, ExpressWindow)

				}),
				widget.NewButton("Stop/Continue Slide show", func() {
					setDelay(slidePause, slideDelay)
					an, _ := anim.Get()
					if an {
						anim.Set(false)
					} else {
						anim.Set(true)
						go startSlideShow(a, e, header, firstTable, grad.Selected, f, anim, curPathwayIndex, slideDelay, ExpressWindow)
					}

				}),
				widget.NewButton("Previous Slide", func() {
					startIdx, _ := curPathwayIndex.Get()
					if startIdx < 2 {
						log.Println("Column", startIdx-1, "cannot be accessed !")
						return
					}
					curPathwayIndex.Set(startIdx - 1)
					go drawExp(a, e, header, firstTable, header[startIdx-1], grad.Selected, f, curPathwayIndex, ExpressWindow)
				}),
				widget.NewButton("Next Slide", func() {
					startIdx, _ := curPathwayIndex.Get()
					if startIdx > len(header)-2 {
						log.Println("Data table have", len(header)-1, "columns ! Column", startIdx+1, "cannot be accessed !")
						return
					}
					curPathwayIndex.Set(startIdx + 1)
					go drawExp(a, e, header, firstTable, header[startIdx+1], grad.Selected, f, curPathwayIndex, ExpressWindow)
				}),
			),
		),
		densityPlot,
		container.NewHBox(
			widget.NewLabel("Max :"),
			widget.NewButton("Apply", func() {
				vmax, _ := userMaxExp.Get()
				vmin, _ := userMinExp.Get()
				//go updateMaxExp(v, a, e, userSel, grad.Selected, f, ExpressWindow)
				go updateMinMaxExp(vmin, vmax, a, e, userSel, grad.Selected, f, ExpressWindow)
			}),
		),
		MaxExp,
		container.NewHBox(
			widget.NewLabel("Min :"),
			widget.NewButton("Apply", func() {
				vmax, _ := userMaxExp.Get()
				vmin, _ := userMinExp.Get()

				go updateMinMaxExp(vmin, vmax, a, e, userSel, grad.Selected, f, ExpressWindow)
			}),
		),
		MinExp,
		widget.NewLabel("Dots Opacity [0-100%] :"),
		DotOpacity,
		container.NewHBox(
			opacityGradient,
			widget.NewLabel("Min"),
			gradMinWdgt,
			widget.NewLabel("Max"),
			gradMaxWdgt,
		),
	)
	ExpressWindow.SetContent(content)
	ExpressWindow.Show()
}

func startSlideShow(a fyne.App, e *Editor, header []string, firstTable, grad string, f binding.Float, anim binding.Bool, curPathwayIndex binding.Int, slideDelay binding.Float, ExpressWindow fyne.Window) {
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
			drawExp(a, e, header, firstTable, header[i], grad, f, curPathwayIndex, ExpressWindow)
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

func drawExp(a fyne.App, e *Editor, header []string, filename string, expcol, gradien string, f binding.Float, curPathwayIndex binding.Int, ExpressWindow fyne.Window) {
	f.Set(0.2)     // progress bar set to 20%
	initCluster(e) // remove all dots of the cluster container
	pref := a.Preferences()
	// Dot opacity
	DotOp := binding.BindPreferenceFloat("dotOpacity", pref) // pref binding for the  dot opacity
	opacity, _ := DotOp.Get()
	op := uint8(opacity)
	// Dot opacity gradient
	gradop := binding.BindPreferenceBool("gradOpacity", pref)
	opGrad, _ := gradop.Get()                               // enabled/disabled opacity gradient
	gradMax := binding.BindPreferenceFloat("gradMax", pref) // pref binding for user Max value for the opacity gradient
	opMax, _ := gradMax.Get()
	gradMin := binding.BindPreferenceFloat("gradMin", pref) // pref binding for user Min value for the opacity gradient
	opMin, _ := gradMin.Get()
	if opMin >= opMax {
		log.Println("Min threshold, ", opMin, " must be <= Max threshold, ", opMax)
	}

	clustDia := binding.BindPreferenceInt("clustDotDiam", pref) //  dot diameter
	diameter, _ := clustDia.Get()
	diameter = ApplyZoomInt(e, diameter)

	expressions, pts := getExpress(a, header, filename, expcol, curPathwayIndex) // []expressions and []Point
	if len(expressions) < 1 {
		log.Println("Intensities not availble for column", expcol)
		return
	}
	f.Set(0.3) // progress bar set to 30% after data reading
	nbPts := len(pts)
	scaleExp, min, max := filter.ScaleSlice01(expressions)

	// density plot of the expression distribution
	go plot.BuildDensity(expressions, 100., filter.TrimString(expcol, 40), ExpressWindow)
	go saveTMPfiles(pts, expressions, min, max, nbPts)

	for c := 0; c < nbPts; c++ {
		// progress bar increases when 50% of points are loaded
		if c == int(nbPts/2) {
			f.Set(0.5) // 50 % progression for progress bar
		}
		// transparency gradient
		if opGrad && opMin < opMax {
			op = gradTransp(expressions[c], min, max, opMin, opMax)
		}
		//op = gradTransp(expressions[c], min, max)
		clcolor := gradUser(gradien)(scaleExp[c])

		e.drawcircle(ApplyZoomInt(e, pts[c].X), ApplyZoomInt(e, pts[c].Y), diameter, color.NRGBA{clcolor.R, clcolor.G, clcolor.B, op})

	}
	// draw legend title, dot and value for the current cexpression
	R, G, B, _ := plot.GetPrefColorRGBA(a, "legendColR", "legendColG", "legendColB", "legendColA")
	colorText := color.NRGBA{uint8(R), uint8(G), uint8(B), 255}
	titleLegend(e, expcol, colorText)
	expLegend(e, op, diameter, gradien, min, max, colorText)

	e.clusterContainer.Refresh()
	f.Set(0.) // reset progress bar
}

// gradTransp compute opacity based on expression score
func gradTransp(exp, min, max, opMin, opMax float64) uint8 {
	if max == min {
		return 255
	}
	if exp <= opMin {
		return 0
	}
	if exp >= opMax {
		return 255
	}
	if opMax < max {
		max = opMax
	}
	if opMin > min {
		min = opMin
	}
	return uint8(255. * (exp - min) / (max - min))
}

// print pathway name on top of image
func titleLegend(e *Editor, title string, c color.NRGBA) {
	AbsText(e.clusterContainer, 50, 30, truncTitle(title), 20, c)
}

// truncTitle return a truncated title
// if the title is larger than the window, there is an artefact and a shift between
// expression dots and the image
func truncTitle(title string) string {
	// get the window width
	prefs := fyne.CurrentApp().Preferences()
	winW := binding.BindPreferenceFloat("winW", prefs) // set the link to preferences for win width
	wW, _ := winW.Get()
	if wW < 700 {
		return filter.TrimString(title, int(40.*wW/800.))
	}
	return filter.TrimString(title, 40)
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
		co := gradUser(gradien)(float64(i) / 5.)
		e.drawcircle(x, y+150-sp*i, diameter*100/e.zoom, color.NRGBA{co.R, co.G, co.B, op})
	}
}

// compute the true expression value from a scaled [0-1] value
func unscale(v, min, max float64) float64 {
	return v*(max-min) + min
}

// grad return the gradien function with name "gradien"
func gradUser(gradien string) func(float64) RGB {
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

// save pts and scaleExp to temporary files to be used by the min/max slider
func saveTMPfiles(pts []filter.Point, expressions []float64, min, max float64, nbpts int) {

	tmp := filter.Record{Pts: pts, Exp: expressions, Min: min, Max: max, NbPts: nbpts}
	filter.DumpJson("temp/expressTMP.json", tmp)

}

// load pts and scaleExp to temporary files to be used by the min/max slider
func loadTMPfiles(fname string) filter.Record {
	return filter.LoadJson(fname)
}

// update the Min and Max expression
func updateMinMaxExp(vmin, vmax float64, a fyne.App, e *Editor, expcol, gradien string, f binding.Float, ExpressWindow fyne.Window) {
	tmp := loadTMPfiles("temp/expressTMP.json")
	newMin := vmin/100.*(tmp.Max-tmp.Min) + tmp.Min
	newMax := vmax/100.*(tmp.Max-tmp.Min) + tmp.Min
	scaleExp := filter.ScaleSliceMinMax(tmp.Exp, newMin, newMax)
	refreshExp(a, e, newMin, newMax, tmp, scaleExp, expcol, gradien, f, ExpressWindow)
}

func refreshExp(a fyne.App, e *Editor, newMin, newMax float64, tmp filter.Record, scaleExp []float64, expcol, gradien string, f binding.Float, ExpressWindow fyne.Window) {
	f.Set(0.2)     // progress bar set to 20%
	initCluster(e) // remove all dots of the cluster container
	pref := a.Preferences()
	// Dot opacity
	DotOp := binding.BindPreferenceFloat("dotOpacity", pref) // pref binding for the  dot opacity
	opacity, _ := DotOp.Get()
	op := uint8(opacity)
	// Dot opacity gradient
	gradop := binding.BindPreferenceBool("gradOpacity", pref)
	opGrad, _ := gradop.Get()                               // enabled/disabled opacity gradient
	gradMax := binding.BindPreferenceFloat("gradMax", pref) // pref binding for user Max value for the opacity gradient
	opMax, _ := gradMax.Get()
	gradMin := binding.BindPreferenceFloat("gradMin", pref) // pref binding for user Min value for the opacity gradient
	opMin, _ := gradMin.Get()
	if opMin >= opMax {
		log.Println("Min threshold, ", opMin, " must be <= Max threshold, ", opMax)
	}

	clustDia := binding.BindPreferenceInt("clustDotDiam", pref) //  dot diameter
	diameter, _ := clustDia.Get()
	diameter = ApplyZoomInt(e, diameter)

	if len(tmp.Exp) < 1 {
		log.Println("Intensities not availble for column", expcol)
		return
	}
	f.Set(0.3) // progress bar set to 30% after data reading

	for c := 0; c < tmp.NbPts; c++ {
		// progress bar increases when 50% of points are loaded
		if c == int(tmp.NbPts/2) {
			f.Set(0.5) // 50 % progression for progress bar
		}
		// transparency gradient
		if opGrad && opMin < opMax {
			op = gradTransp(tmp.Exp[c], newMin, newMax, opMin, opMax)
		}
		clcolor := gradUser(gradien)(scaleExp[c])

		e.drawcircle(ApplyZoomInt(e, tmp.Pts[c].X), ApplyZoomInt(e, tmp.Pts[c].Y), diameter, color.NRGBA{clcolor.R, clcolor.G, clcolor.B, op})

	}
	// draw legend title, dot and value for the current cexpression
	R, G, B, _ := plot.GetPrefColorRGBA(a, "legendColR", "legendColG", "legendColB", "legendColA")
	colorText := color.NRGBA{uint8(R), uint8(G), uint8(B), 255}
	titleLegend(e, expcol, colorText)
	expLegend(e, op, diameter, gradien, newMin, newMax, colorText)

	e.clusterContainer.Refresh()
	ExpressWindow.Content().Refresh()
	f.Set(0.) // reset progress bar
}

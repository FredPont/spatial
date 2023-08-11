package ui

import (
	"image/color"
	"log"
	"spatial/src/filter"
	"spatial/src/plot"
	"spatial/src/pogrebDB"
	"spatial/src/pref"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/mazznoer/colorgrad"
)

func buttonDrawExpress(a fyne.App, e *Editor, preference fyne.Preferences, f binding.Float, header []string, firstTable string) {
	ExpressWindow := a.NewWindow("Expression")

	// show choice of different gradien
	grad := widget.NewRadioGroup([]string{"Turbo", "Viridis", "Inferno", "Plasma", "White - Red", "Yellow - Red", "Purple - Red", "Red - Yellow ", "Custom"}, func(s string) {
		//fmt.Println("Selected <", s, ">")
	})

	// Dot opacity
	DotOpacity := widget.NewSlider(0., 255.)
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
	slideDelay.Set(15)                  // 15 sec pause

	slidePause := widget.NewEntry()
	slidePause.SetPlaceHolder("Pause between slides (sec)")
	slidePause.SetText("15")

	content := container.NewVBox(
		container.NewHBox(
			widget.NewLabel("Select the variable and gradient :"),
			widget.NewButtonWithIcon("Close", theme.LogoutIcon(), func() {
				initOpacityGdt() // remove the opacity preferences before closing the window
				ExpressWindow.Close()
			}),
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
					go startEspressComput(a, e, header, firstTable, userSel, grad.Selected, f, curPathwayIndex, ExpressWindow)
					//go drawImageExp(a, e, header, firstTable, userSel, grad.Selected, f, curPathwayIndex, ExpressWindow)
					//go MTdrawImageExp(a, e, header, firstTable, userSel, grad.Selected, f, curPathwayIndex, ExpressWindow)

				}),
				legendcol,
				widget.NewSeparator(),
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
						startIdx, _ := curPathwayIndex.Get()
						log.Println("curPathwayIndex=", startIdx)
						log.Println("stop slide show")
						return
					} else {
						anim.Set(true)
						log.Println("start slide show")
						startIdx, _ := curPathwayIndex.Get()
						log.Println("curPathwayIndex=", startIdx)
						if startIdx >= len(header)-2 {
							log.Println("Column", startIdx+1, "cannot be accessed, final column reached !")
							return
						}
						//curPathwayIndex.Set(startIdx + 1) // start slide show with the next pathway
						go startSlideShow(a, e, header, firstTable, grad.Selected, f, anim, curPathwayIndex, slideDelay, ExpressWindow)
					}

				}),
				widget.NewButton("Previous Slide", func() {
					startIdx, _ := curPathwayIndex.Get()
					if startIdx < 2 {
						log.Println("Column", startIdx-1, "cannot be accessed !")
						return
					}
					//log.Println("pathway index = ", startIdx)
					//curPathwayIndex.Set(startIdx - 1)
					//startIdx, _ = curPathwayIndex.Get()
					//log.Println("new pathway index = ", startIdx)
					//go drawExp(a, e, header, firstTable, header[startIdx-1], grad.Selected, f, curPathwayIndex, ExpressWindow)
					go startEspressComput(a, e, header, firstTable, header[startIdx-1], grad.Selected, f, curPathwayIndex, ExpressWindow)
				}),
				widget.NewButton("Next Slide", func() {
					startIdx, _ := curPathwayIndex.Get()
					if startIdx > len(header)-2 {
						log.Println("Data table have", len(header)-1, "columns ! Column", startIdx+1, "cannot be accessed !")
						return
					}
					curPathwayIndex.Set(startIdx + 1)
					//go drawExp(a, e, header, firstTable, header[startIdx+1], grad.Selected, f, curPathwayIndex, ExpressWindow)
					go startEspressComput(a, e, header, firstTable, header[startIdx+1], grad.Selected, f, curPathwayIndex, ExpressWindow)
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
	log.Println("start slide show : curPathwayIndex=", startIdx)
	// listen to anim
	an, _ := anim.Get()
	if !an {
		return
	}

	for i := startIdx + 1; i < m; i++ {

		if header[i] != "" {
			log.Println(i, "/", m-1, " column :", header[i])

			//drawExp(a, e, header, firstTable, header[i], grad, f, curPathwayIndex, ExpressWindow)
			startEspressComput(a, e, header, firstTable, header[i], grad, f, curPathwayIndex, ExpressWindow)
			pause, _ := slideDelay.Get()
			time.Sleep(time.Duration(1000*pause) * time.Millisecond)
		}
		// listen to anim
		an, _ := anim.Get()
		if !an {
			//curPathwayIndex.Set(i)
			break
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

// getExpress get the values from the expression column selected by the user and the XY columns
func getExpress(a fyne.App, header []string, filename string, expcol string, curPathwayIndex binding.Int) ([]float64, []filter.Point) {
	// get the user preference for using the database
	pref := fyne.CurrentApp().Preferences()
	useDBpref := binding.BindPreferenceBool("useDataBase", pref)
	useDB, _ := useDBpref.Get()
	// use the pogreb database instead of CSV if selected by the user in preferences
	if useDB {
		return pogrebDB.DBgetExpress(a, header, filename, expcol, curPathwayIndex)
	}
	return getExpressCSV(a, header, filename, expcol, curPathwayIndex)
}

func getExpressCSV(a fyne.App, header []string, filename string, expcol string, curPathwayIndex binding.Int) ([]float64, []filter.Point) {
	pref := a.Preferences()
	// X coordinates
	xcor := binding.BindPreferenceString("xcor", pref) // set the link to preferences for X coordinates
	xc, _ := xcor.Get()

	// y coordinates
	ycor := binding.BindPreferenceString("ycor", pref) // set the link to preferences for y coordinates
	yc, _ := ycor.Get()

	colIndexes := filter.GetColIndex(header, []string{expcol, xc, yc}) // selected column + xy
	curPathwayIndex.Set(colIndexes[0])                                 // set the current expression Index to the selected column to enable button next/previous slide
	return filter.ReadExpress(a, filename, colIndexes)
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
func expLegend(e *Editor, op uint8, diameter int, gradien string, grad colorgrad.Gradient, min, max float64, c color.NRGBA) {
	x, y := 13, 30
	sp := 25

	for i := 5; i >= 0; i-- {
		//exp := fmt.Sprintf("%.1f", unscale(float64(i)/5., min, max))
		exp := TicksDecimals(unscale(float64(i)/5., min, max))
		AbsText(e.clusterContainer, x+20, y+155-sp*i, exp, 15, c)
		co := gradUser(gradien, grad, float64(i)/5.)
		// compute the spot max diameter to avoid overlap
		spotDiam := diameter * 100 / e.zoom
		if spotDiam >= 11 {
			spotDiam = 11
		}
		e.drawcircle(x, y+150-sp*i, spotDiam, color.NRGBA{co.R, co.G, co.B, op})
	}
}

// draw expression legend with dots and values and gradient opacity
func gradOpLegend(e *Editor, diameter int, gradien string, grad colorgrad.Gradient, min, max, opMin, opMax float64, c color.NRGBA) {
	x, y := 13, 30
	sp := 25

	for i := 5; i >= 0; i-- {
		exp := TicksDecimals(unscale(float64(i)/5., min, max))
		AbsText(e.clusterContainer, x+20, y+155-sp*i, exp, 15, c)
		co := gradUser(gradien, grad, float64(i)/5.)
		op := gradTransp(unscale(float64(i)/5., min, max), min, max, opMin, opMax)
		// compute the spot max diameter to avoid overlap
		spotDiam := diameter * 100 / e.zoom
		if spotDiam >= 11 {
			spotDiam = 11
		}
		e.drawcircle(x, y+150-sp*i, spotDiam, color.NRGBA{co.R, co.G, co.B, op})
		//fmt.Println(exp, op)
	}
}

// compute the true expression value from a scaled [0-1] value
func unscale(v, min, max float64) float64 {
	return v*(max-min) + min
}

// gradUser return the gradien (with name "gradien") RGB value at a specific point
func gradUser(gradien string, grad colorgrad.Gradient, val float64) RGB {
	return rgbModel(grad.At(val))
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
	f.Set(0.1)
	tmp := loadTMPfiles("temp/expressTMP.json")
	newMin := vmin/100.*(tmp.Max-tmp.Min) + tmp.Min
	newMax := vmax/100.*(tmp.Max-tmp.Min) + tmp.Min
	scaleExp := filter.ScaleSliceMinMax(tmp.Exp, newMin, newMax)
	//refreshExp(a, e, newMin, newMax, tmp, scaleExp, expcol, gradien, f, ExpressWindow)
	refreshImageExp(a, e, newMin, newMax, tmp, scaleExp, expcol, gradien, f, ExpressWindow)
}

func refreshExp(a fyne.App, e *Editor, newMin, newMax float64, tmp filter.Record, scaleExp []float64, expcol, gradien string, f binding.Float, ExpressWindow fyne.Window) {
	f.Set(0.2)     // progress bar set to 20%
	initCluster(e) // remove all dots of the cluster container
	pref := a.Preferences()
	// Dot opacity
	DotOp := binding.BindPreferenceFloat("dotOpacity", pref) // pref binding for the  dot opacity
	opacity, _ := DotOp.Get()
	op := uint8(opacity)
	// pre-build expression gradient
	grad := preBuildGradient(gradien)
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
	f.Set(0.3)                                            // progress bar set to 30% after data reading
	circlesObjets := make([]fyne.CanvasObject, tmp.NbPts) // store all the circles to add them all in one time
	for c := 0; c < tmp.NbPts; c++ {
		// progress bar increases when 50% of points are loaded
		if c == int(tmp.NbPts/2) {
			f.Set(0.5) // 50 % progression for progress bar
		}
		// transparency gradient
		if opGrad && opMin < opMax {
			op = gradTransp(tmp.Exp[c], newMin, newMax, opMin, opMax)
		}
		clcolor := gradUser(gradien, grad, scaleExp[c])

		//e.drawcircle(ApplyZoomInt(e, tmp.Pts[c].X), ApplyZoomInt(e, tmp.Pts[c].Y), diameter, color.NRGBA{clcolor.R, clcolor.G, clcolor.B, op})
		circle := drawRoundedRect(ApplyZoomInt(e, tmp.Pts[c].X), ApplyZoomInt(e, tmp.Pts[c].Y), diameter, color.NRGBA{clcolor.R, clcolor.G, clcolor.B, op})
		circlesObjets[c] = circle //add the spot to the slice of objects
	}
	// draw legend title, dot and value for the current cexpression
	R, G, B, _ := plot.GetPrefColorRGBA(a, "legendColR", "legendColG", "legendColB", "legendColA")
	colorText := color.NRGBA{uint8(R), uint8(G), uint8(B), 255}
	titleLegend(e, expcol, colorText)
	// transparency gradient
	if opGrad && opMin < opMax {
		//fmt.Println("opGrad", opGrad, " opMin opMax", opMin, opMax)
		gradOpLegend(e, diameter, gradien, grad, newMin, newMax, opMin, opMax, colorText)
	} else {
		expLegend(e, op, diameter, gradien, grad, newMin, newMax, colorText) // not transparency gradien in legend
	}
	//expLegend(e, op, diameter, gradien, grad, newMin, newMax, colorText)
	e.clusterContainer.Objects = append(e.clusterContainer.Objects, circlesObjets...)
	e.clusterContainer.Refresh()
	ExpressWindow.Content().Refresh()
	f.Set(0.) // reset progress bar
}

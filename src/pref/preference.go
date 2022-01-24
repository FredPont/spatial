package pref

import (
	"fmt"
	"spatial/src/filter"

	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

// func SetPref(a fyne.App) {
// 	// scaling factor
// 	a.Preferences().SetFloat("scaleFactor", 1.)
// 	val := a.Preferences().Float("scaleFactor")
// 	fmt.Println("scaleFactor is:", val)
// }

// BuildPref create a form where all preferences can be set
func BuildPref(a fyne.App, head []string) {
	prefs := a.Preferences()
	//myWindow := a.NewWindow("Preferences")
	myWindow := fyne.CurrentApp().NewWindow("Preferences")

	// scaling factor
	scalingFactor := widget.NewEntry()
	sf := binding.BindPreferenceFloat("scaleFactor", prefs) // set the link to preferences for scaling factor
	x, _ := sf.Get()                                        // read the preference for scaling factor
	sftxt := fmt.Sprintf("%.10f", x)                        // convert scaling factor to txt
	scalingFactor.SetPlaceHolder(sftxt)                     // display the preference value for scaling factor

	// coordinates +90° rotation : necessary for 10x Genomics
	r := binding.BindPreferenceBool("rotate", prefs) // set the link to preferences for rotation
	b, _ := r.Get()
	rot := widget.NewCheck("rotate coordinates +90°", func(value bool) {})
	rot.SetChecked(b)

	// X coordinates
	xcor := binding.BindPreferenceString("xcor", prefs) // set the link to preferences for x coordinates
	xc, _ := xcor.Get()
	//xSel := widget.NewSelectEntry(head)
	//xSel.SetText(xc)
	var xSel *widget.Button
	xSel = widget.NewButton(xc, func() {
		ShowTable(head, xcor, xSel, xc)
	})

	// y coordinates
	ycor := binding.BindPreferenceString("ycor", prefs) // set the link to preferences for y coordinates
	yc, _ := ycor.Get()
	//ySel := widget.NewSelectEntry(head)
	//ySel.SetText(yc)
	var ySel *widget.Button
	ySel = widget.NewButton(yc, func() {
		ShowTable(head, ycor, ySel, yc)
	})

	//microscop windows size
	//microscop windows W
	winWidth := widget.NewEntry()
	winW := binding.BindPreferenceFloat("winW", prefs) // set the link to preferences for win width
	wW, _ := winW.Get()
	wWtxt := fmt.Sprintf("%.0f", wW)
	winWidth.SetPlaceHolder(wWtxt)

	//microscop windows Height
	winHeight := widget.NewEntry()
	winH := binding.BindPreferenceFloat("winH", prefs) // set the link to preferences for win width
	wH, _ := winH.Get()
	wHtxt := fmt.Sprintf("%.0f", wH)
	winHeight.SetPlaceHolder(wHtxt)

	// cluster column
	clustercolumn := binding.BindPreferenceString("clustcol", prefs) // set the link to preferences for rotation
	clucol, _ := clustercolumn.Get()
	//clco := widget.NewSelectEntry(head)
	//clco.SetText(clucol)
	var clco *widget.Button
	clco = widget.NewButton(clucol, func() {
		ShowTable(head, clustercolumn, clco, clucol)
	})

	// cluster dot diameter
	clustDotDiam := widget.NewEntry()
	cld := binding.BindPreferenceInt("clustDotDiam", prefs) // set the link to preferences for cluster diameter
	clud, _ := cld.Get()                                    // read the preference for cluster diameter
	cldtxt := fmt.Sprintf("%d", clud)                       // convert cluster diameter to txt
	clustDotDiam.SetPlaceHolder(cldtxt)                     // display the prefsence value for cluster diameter

	// vulcano selection square size in pixels
	// record the vulcano selection square size in pixels in preferences
	vulcSquare := widget.NewEntry()
	vs := binding.BindPreferenceInt("vulcSelectSize", prefs)
	vsquare, err := vs.Get()
	if err != nil {
		log.Println("wrong selection size value !", err)
	}
	vstxt := fmt.Sprintf("%d", vsquare)
	vulcSquare.SetPlaceHolder(vstxt)

	// show choice of different gradien
	// read the custom palettes
	listPalet := []string{"Turbo", "Rainbow", "Sinebow"}
	listPalet = append(listPalet, listPalette()...)

	grad := widget.NewRadioGroup(listPalet, func(s string) {
	})
	cG := binding.BindPreferenceString("clusterGradient", prefs)
	clusterGrad, _ := cG.Get()
	grad.SetSelected(clusterGrad)

	// create form
	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Scaling Factor", Widget: scalingFactor},
			{Text: "Rotate", Widget: rot},
			{Text: "X coordinates", Widget: xSel},
			{Text: "Y coordinates", Widget: ySel},
			{Text: "Image windows Width", Widget: winWidth},
			{Text: "Image windows Height", Widget: winHeight},
			{Text: "Cluster column", Widget: clco},
			{Text: "Cluster dots diameter", Widget: clustDotDiam},
			{Text: "vulcano selection square size in pixels", Widget: vulcSquare},
			{Text: "Clusters color gradient", Widget: grad},
		},
		OnSubmit: func() { // optional, handle form submission

			// scaling factor
			sftxt := scalingFactor.Text
			setPrefToF64(sftxt, "scaleFactor", prefs)

			// coordinates +90° rotation
			prefs.SetBool("rotate", rot.Checked)

			// X coordinates
			//prefs.SetString("xcor", xc)

			// Y coordinates
			//prefs.SetString("ycor", ySel.Entry.Text)

			// microscop windows W
			winWidthTxt := winWidth.Text
			setPrefToF64(winWidthTxt, "winW", prefs)

			// microscop windows H
			winHeightTxt := winHeight.Text
			setPrefToF64(winHeightTxt, "winH", prefs)

			// cluster column
			//prefs.SetString("clustcol", clco.Entry.Text)

			// cluster dot diameter
			cluDiamTXT := clustDotDiam.Text
			setPrefToInt(cluDiamTXT, "clustDotDiam", prefs)
			//prefs.SetInt("clustDotDiam", cluDiam)

			// vulcano selection square size in pixels
			vulSelTXT := vulcSquare.Text
			setPrefToInt(vulSelTXT, "vulcSelectSize", prefs)

			// cluster gradient selection
			// gradien default
			def := "Turbo"
			if grad.Selected == "" {
				grad.Selected = def
			}

			prefs.SetString("clusterGradient", grad.Selected)
			//log.Println("Form submitted:", scalingFactor.Text)

			myWindow.Close()
		},
	}

	myWindow.SetContent(form)
	myWindow.Show()
}

// set prefs of widget.NewEntry() to float64
func setPrefToF64(s, prefID string, prefs fyne.Preferences) {
	if s != "" {
		f64, err := strconv.ParseFloat(s, 64)
		if err != nil {
			log.Printf("unable to convert string to float ! %T, %v\n", f64, f64)
		}
		prefs.SetFloat(prefID, f64)
	}
}

// set prefs of widget.NewEntry() to float64
func setPrefToInt(s, prefID string, prefs fyne.Preferences) {
	if s != "" {
		integ, err := strconv.Atoi(s)
		if err != nil {
			log.Printf("unable to convert string to float ! %T, %v\n", integ, integ)
		}
		prefs.SetInt(prefID, integ)
	}
}

func listPalette() []string {
	return filter.ListFiles("src/palette/")
}

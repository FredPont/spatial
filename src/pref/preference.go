package pref

import (
	"fmt"
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
	pref := a.Preferences()
	//myWindow := a.NewWindow("Preferences")
	myWindow := fyne.CurrentApp().NewWindow("Preferences")

	// scaling factor
	scalingFactor := widget.NewEntry()
	sf := binding.BindPreferenceFloat("scaleFactor", pref) // set the link to preferences for scaling factor
	x, _ := sf.Get()                                       // read the preference for scaling factor
	sftxt := fmt.Sprintf("%.10f", x)                       // convert scaling factor to txt
	scalingFactor.SetPlaceHolder(sftxt)                    // display the prefence value for scaling factor

	// coordinates +90° rotation : necessary for 10x Genomics
	r := binding.BindPreferenceBool("rotate", pref) // set the link to preferences for rotation
	b, _ := r.Get()
	rot := widget.NewCheck("rotate coordinates +90°", func(value bool) {})
	rot.SetChecked(b)

	// X coordinates
	xcor := binding.BindPreferenceString("xcor", pref) // set the link to preferences for rotation
	xc, _ := xcor.Get()
	xSel := widget.NewSelectEntry(head)
	xSel.SetText(xc)

	// y coordinates
	ycor := binding.BindPreferenceString("ycor", pref) // set the link to preferences for rotation
	yc, _ := ycor.Get()
	ySel := widget.NewSelectEntry(head)
	ySel.SetText(yc)

	//microscop windows size
	//microscop windows W
	winWidth := widget.NewEntry()
	winW := binding.BindPreferenceFloat("winW", pref) // set the link to preferences for win width
	wW, _ := winW.Get()
	wWtxt := fmt.Sprintf("%.0f", wW)
	winWidth.SetPlaceHolder(wWtxt)

	//microscop windows Height
	winHeight := widget.NewEntry()
	winH := binding.BindPreferenceFloat("winH", pref) // set the link to preferences for win width
	wH, _ := winH.Get()
	wHtxt := fmt.Sprintf("%.0f", wH)
	winHeight.SetPlaceHolder(wHtxt)

	// cluster column
	clustercolumn := binding.BindPreferenceString("clustcol", pref) // set the link to preferences for rotation
	clucol, _ := clustercolumn.Get()
	clco := widget.NewSelectEntry(head)
	clco.SetText(clucol)

	// cluster dot diameter
	clustDotDiam := widget.NewEntry()
	cld := binding.BindPreferenceInt("clustDotDiam", pref) // set the link to preferences for cluster diameter
	clud, _ := cld.Get()                                   // read the preference for cluster diameter
	cldtxt := fmt.Sprintf("%d", clud)                      // convert cluster diameter to txt
	clustDotDiam.SetPlaceHolder(cldtxt)                    // display the prefence value for cluster diameter

	// vulcano selection square size in pixels
	// record the vulcano selection square size in pixels in preferences
	vulcSquare := widget.NewEntry()
	vs := binding.BindPreferenceInt("vulcSelectSize", pref)
	vsquare, err := vs.Get()
	if err != nil {
		log.Println("wrong selection size value !", err)
	}
	vstxt := fmt.Sprintf("%d", vsquare)
	vulcSquare.SetPlaceHolder(vstxt)

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
		},
		OnSubmit: func() { // optional, handle form submission

			// scaling factor
			sftxt := scalingFactor.Text
			setPrefToF64(sftxt, "scaleFactor", pref)

			// coordinates +90° rotation
			pref.SetBool("rotate", rot.Checked)

			// X coordinates
			pref.SetString("xcor", xSel.Entry.Text)

			// Y coordinates
			pref.SetString("ycor", ySel.Entry.Text)

			// microscop windows W
			winWidthTxt := winWidth.Text
			setPrefToF64(winWidthTxt, "winW", pref)

			// microscop windows H
			winHeightTxt := winHeight.Text
			setPrefToF64(winHeightTxt, "winH", pref)

			// cluster column
			pref.SetString("clustcol", clco.Entry.Text)

			// cluster dot diameter
			cluDiamTXT := clustDotDiam.Text
			setPrefToInt(cluDiamTXT, "clustDotDiam", pref)
			//pref.SetInt("clustDotDiam", cluDiam)

			// vulcano selection square size in pixels
			vulSelTXT := vulcSquare.Text
			setPrefToInt(vulSelTXT, "vulcSelectSize", pref)

			//log.Println("Form submitted:", scalingFactor.Text)

			myWindow.Close()
		},
	}

	myWindow.SetContent(form)
	myWindow.Show()
}

// set pref of widget.NewEntry() to float64
func setPrefToF64(s, prefId string, pref fyne.Preferences) {
	if s != "" {
		f64, err := strconv.ParseFloat(s, 64)
		if err != nil {
			log.Printf("unable to convert string to float ! %T, %v\n", f64, f64)
		}
		pref.SetFloat(prefId, f64)
	}
}

// set pref of widget.NewEntry() to float64
func setPrefToInt(s, prefId string, pref fyne.Preferences) {
	if s != "" {
		integ, err := strconv.Atoi(s)
		if err != nil {
			log.Printf("unable to convert string to float ! %T, %v\n", integ, integ)
		}
		pref.SetInt(prefId, integ)
	}
}

package pref

import (
	"fmt"
	"lasso/src/filter"
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

func BuildPref(a fyne.App) {
	pref := a.Preferences()
	myWindow := a.NewWindow("Preferences")

	// scaling factor
	scalingFactor := widget.NewEntry()
	sf := binding.BindPreferenceFloat("scaleFactor", pref) // set the link to preferences for scaling factor
	x, _ := sf.Get()                                       // read the preference for scaling factor
	sftxt := fmt.Sprintf("%.8f", x)                        // convert scaling factor to txt
	scalingFactor.SetPlaceHolder(sftxt)                    // display the prefence value for scaling factor

	// coordinates +90° rotation : necessary for 10x Genomics
	r := binding.BindPreferenceBool("rotate", pref) // set the link to preferences for rotation
	b, _ := r.Get()
	rot := widget.NewCheck("rotate coordinates +90°", func(value bool) {})
	rot.SetChecked(b)

	// X coordinates
	head := filter.ReadHeader("data/H_Exp_Spatial_seuratv3.tsv")
	xcor := binding.BindPreferenceString("xcor", pref) // set the link to preferences for rotation
	xc, _ := xcor.Get()
	xSel := widget.NewSelectEntry(head)
	xSel.SetText(xc)

	// y coordinates
	ycor := binding.BindPreferenceString("ycor", pref) // set the link to preferences for rotation
	yc, _ := ycor.Get()
	ySel := widget.NewSelectEntry(head)
	ySel.SetText(yc)

	// create form
	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Scaling Factor", Widget: scalingFactor},
			{Text: "Rotate", Widget: rot},
			{Text: "X coordinates", Widget: xSel},
			{Text: "Y coordinates", Widget: ySel},
		},
		OnSubmit: func() { // optional, handle form submission

			// scaling factor
			sftxt := scalingFactor.Text

			if sftxt != "" {
				sffloat, err := strconv.ParseFloat(sftxt, 64)
				if err == nil {
					log.Printf("%T, %v\n", sf, sf)
				}
				pref.SetFloat("scaleFactor", sffloat) // store the new preference for scaling factor
			}

			// coordinates +90° rotation
			pref.SetBool("rotate", rot.Checked)

			// X coordinates
			pref.SetString("xcor", xSel.Entry.Text)

			// Y coordinates
			pref.SetString("ycor", ySel.Entry.Text)

			log.Println("Form submitted:", scalingFactor.Text)

			myWindow.Close()
		},
	}

	myWindow.SetContent(form)
	myWindow.Show()
}

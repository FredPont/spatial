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

func BuildPref(a fyne.App) {
	pref := a.Preferences()

	myWindow := a.NewWindow("Preferences")

	scalingFactor := widget.NewEntry()
	sf := binding.BindPreferenceFloat("scaleFactor", pref) // set the link to preferences for scaling factor
	x, _ := sf.Get()                                       // read the preference for scaling factor
	sftxt := fmt.Sprintf("%.8f", x)                        // convert scaling factor to txt
	scalingFactor.SetPlaceHolder(sftxt)                    // display the prefence value for scaling factor
	textArea := widget.NewMultiLineEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Scaling Factor", Widget: scalingFactor}},
		OnSubmit: func() { // optional, handle form submission
			sftxt := scalingFactor.Text
			sffloat, err := strconv.ParseFloat(sftxt, 64)
			if err == nil {
				log.Printf("%T, %v\n", sf, sf)
			}
			pref.SetFloat("scaleFactor", sffloat) // store the new preference for scaling factor

			log.Println("Form submitted:", scalingFactor.Text)
			log.Println("multiline:", textArea.Text)

			myWindow.Close()
		},
	}

	// we can also append items
	form.Append("Text", textArea)

	myWindow.SetContent(form)
	myWindow.Show()
}

package resampling

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func ResampWin(a fyne.App, Table *string, header []string, f binding.Float) {

	win := a.NewWindow("Resampling")

	firstTable := *Table

	resampRate := 1

	// Create the explaination text
	myText := "It is recommended to reduce very large datasets by resampling.\nCaution: a too large resampling rate can create \"holes\"\n in the spot distribution.\nSelect the resampling rate 1/n and the columns to convert to integer.\nXY coordinates and cluster numbers must be integers"

	explain := widget.NewLabel(myText)
	explain.Alignment = fyne.TextAlignLeading
	// Create a list of strings for the checkbox labels
	myCheckboxLabels := header

	// Create a CheckGroup widget with the checkbox labels
	myCheckGroup := widget.NewCheckGroup(myCheckboxLabels, func(selected []string) {
		// Handle the checkbox selection change event
	})

	// Create a text entry widget for the resampling rate 1/n
	skipRows := widget.NewEntry()
	skipRows.SetText("1")
	//skipRows.SetPlaceHolder("n")
	skipRows.Validator = func(s string) error {
		_, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		return nil
	}

	// Create a button to compute the resampling
	computeButton := widget.NewButtonWithIcon("Start computation", theme.ComputerIcon(), func() {
		// get the resampling rate
		valueInt, err := strconv.Atoi(skipRows.Text)
		if err != nil {
			fmt.Println("Invalid resampling rate:", err)
		} else {
			resampRate = valueInt
			fmt.Println("The resampling rate is:", resampRate)

			// get the column indexes
			colnames := myCheckGroup.Selected
			colIndexes := findAllcolIndexe(colnames, header)

			Resample("data/"+firstTable, "data/"+"0_"+firstTable, resampRate, colIndexes, colnames, "\t", f)
			*Table = "0_" + firstTable

			win.Close() // close tool window
		}
	})

	// Create a button to get the selected checkboxes
	closeButton := widget.NewButtonWithIcon("Cancel", theme.LogoutIcon(), func() {
		win.Close() // close tool window
	})

	// Create a container to hold the CheckGroup widget
	myContainer := container.New(layout.NewGridLayout(2), container.NewVBox(

		explain,
		container.NewHBox(
			widget.NewLabel("resampling rate 1/n"),
			skipRows),

		container.NewHBox(
			computeButton,
			closeButton,
		),
	),
		container.NewScroll(
			myCheckGroup),
		//myCheckGroup,
	)

	//myScrollContent := container.NewScroll(myContainer)

	win.SetContent(myContainer)
	// Set the content of the window to the container
	//win.SetContent(myContainer)

	// Show the window
	win.Resize(fyne.Size{Width: 400, Height: 400})
	win.Show()
}

// find the column index with the colname
func findIndex(s string, slice []string) int {
	for i, v := range slice {
		if v == s {
			return i
		}
	}
	return -1 // not found
}

func findAllcolIndexe(colnames []string, header []string) []int {
	var indexes []int

	for _, col := range colnames {
		idx := findIndex(col, header)
		if idx != -1 {
			indexes = append(indexes, idx)
		}
	}
	return indexes
}

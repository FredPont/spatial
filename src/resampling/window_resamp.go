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
	myText := "It is recommended to reduce very large datasets by resampling.\nCaution: a too large resampling rate can create \"holes\"\n in the spot distribution.\nSelect the resampling rate 1/n and the columns to convert to integer.\nXY coordinates and cluster numbers must be integers.\nTo use database instead of CSV, restart the software after resampling."

	explain := widget.NewLabel(myText)
	explain.Alignment = fyne.TextAlignLeading
	// Create a list of strings for the checkbox labels
	headerMap := make(map[string]bool, len(header))
	myCheckbox := listColums(header, headerMap)

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
			colnames := selectedCols(headerMap)
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
			myCheckbox),
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

func listColums(header []string, headerMap map[string]bool) *widget.List {
	list := widget.NewList(
		func() int {
			return len(header)
		},
		func() fyne.CanvasObject {
			return widget.NewCheck("", nil)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Check).SetText(header[i])
			o.(*widget.Check).SetChecked(o.(*widget.Check).Checked)
			o.(*widget.Check).OnChanged = func(value bool) {
				headerMap[header[i]] = true
			}
		})

	return list
}

// get all selected columns from header map
func selectedCols(headerMap map[string]bool) []string {
	var cols []string
	for k, v := range headerMap {
		if v {
			cols = append(cols, k)
		}

	}
	return cols
}

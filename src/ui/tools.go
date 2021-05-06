package ui

import (
	"encoding/csv"
	"fmt"
	"image/png"
	"lasso/src/filter"
	"lasso/src/pref"
	"log"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

// build tools window with buttons and text entry
func BuildTools(a fyne.App, w2, w fyne.Window, e *editor) {
	// get informations from data files to be used with buttons
	dataFiles := filter.ListFiles("data/")              // list all tables in data dir
	header := filter.ReadHeader("data/" + dataFiles[0]) // header of 1st table found in data

	gatename := widget.NewEntry()
	gatename.SetPlaceHolder("Selection name...")

	// scalingFactor := widget.NewEntry()

	// scalingFactor.SetPlaceHolder(a.Preferences().Float("scaleFactor"))

	content := container.NewVBox(
		gatename,
		widget.NewButton("Filter tables with active gates", func() {
			alledges := e.drawSurface.alledges
			go filterActiveGates(alledges, dataFiles, gatename.Text, a.Preferences())
			go saveGates(gatename.Text, e)
		}),
		widget.NewButton("Clear all gates", func() {
			clearDots(e)
		}),
		widget.NewButton("Screen shot", func() {
			screenShot(w, gatename.Text)
		}),
		widget.NewButton("Preferences", func() {
			pref.BuildPref(a, header)
		}),
		widget.NewButton("Exit", func() {
			os.Exit(0)
		}),
	)

	w2.SetContent(content)

}

// clear all gates on draw surface and init all edges
func clearDots(e *editor) {
	e.drawSurface.clearPolygon(e.drawSurface.allpoints)
	e.layer.Refresh()
	initAlledges(e) // reset alledges
}

// save the gates to csv files
func saveGates(gateName string, e *editor) {
	fmt.Println("save gates")
	for i, poly := range e.drawSurface.alledges {
		if len(poly) < 3 {
			continue
		}
		fmt.Println(i, " ", poly)
		out := strconv.Itoa(i) + "_" + gateName
		writeCSV(out, poly)
	}

}

// save image to file
// credits https://www.devdungeon.com/content/working-images-go
func screenShot(w fyne.Window, filename string) {
	out := w.Canvas().Capture()

	// outputFile is a File type which satisfies Writer interface
	path := "plots/" + filename + ".png"
	outputFile, err := os.Create(path)
	if err != nil {
		// Handle error
		log.Println("The image cannot be saved to the file")
	}

	// Encode takes a writer interface and an image interface
	// We pass it the File and the RGBA
	png.Encode(outputFile, out)
	log.Println("Saving image to ", path)

	// Don't forget to close files
	outputFile.Close()

}

// write polygon edge to file
func writeCSV(filename string, poly []filter.Point) {
	path := "gates/" + filename + ".csv"
	file, err := os.Create(path)
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Comma = '\t'

	for _, value := range poly {
		err := writer.Write(filterPtToStr(value))
		checkError("Cannot write to file", err)
	}
}

// convert one filter.Point to one []string
func filterPtToStr(p filter.Point) []string {
	x := strconv.Itoa(p.X)
	y := strconv.Itoa(p.Y)
	str := []string{x, y}
	return str
}

// convert one fyne position to one []string
func fynPosToStr(p fyne.Position) []string {
	x := fmt.Sprintf("%.0f", p.X)
	y := fmt.Sprintf("%.0f", p.Y)
	str := []string{x, y}
	return str
}

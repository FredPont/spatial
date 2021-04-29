package ui

import (
	"encoding/csv"
	"fmt"
	"image/png"
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

func BuildTools(w2, w fyne.Window, e *editor) {
	gatename := widget.NewEntry()
	gatename.SetPlaceHolder("Selection name...")

	content := container.NewVBox(
		gatename,
		widget.NewButton("Filter tables with active gates", func() {
			saveGates(gatename.Text, e)
		}),
		widget.NewButton("Clear all gates", func() {
			clearDots(e)
		}),
		widget.NewButton("Screen shot", func() {
			screenShot(w, gatename.Text)
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
		fmt.Println("The image cannot be saved to the file")
	}

	// Encode takes a writer interface and an image interface
	// We pass it the File and the RGBA
	png.Encode(outputFile, out)

	// Don't forget to close files
	outputFile.Close()

}

// write polygon edge to file
func writeCSV(filename string, poly []fyne.Position) {
	path := "gates/" + filename + ".csv"
	file, err := os.Create(path)
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Comma = '\t'

	for _, value := range poly {
		err := writer.Write(fynPosToStr(value))
		checkError("Cannot write to file", err)
	}
}

// convert one fyne position to one []string
func fynPosToStr(p fyne.Position) []string {
	x := fmt.Sprintf("%.0f", p.X)
	y := fmt.Sprintf("%.0f", p.Y)
	str := []string{x, y}
	return str
}

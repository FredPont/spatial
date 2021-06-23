package ui

import (
	"encoding/csv"
	"fmt"
	"image/png"
	"lasso/src/filter"
	"lasso/src/plot"
	"lasso/src/pref"
	"time"

	//"lasso/src/plot"

	"log"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/software"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

// BuildTools build tools window with buttons and text entry
func BuildTools(a fyne.App, w2, w fyne.Window, e *Editor) {
	preference := a.Preferences()
	// get informations from data files to be used with buttons
	dataFiles := filter.ListFiles("data/") // list all tables in data dir
	firstTable := dataFiles[0]
	header := filter.ReadHeader("data/" + firstTable) // header of 1st table found in data

	gatename := widget.NewEntry()
	gatename.SetPlaceHolder("Selection name...")

	// cluster opacity
	clustOpacity := binding.BindPreferenceFloat("clustOpacity", preference) // pref binding for the cluster dot opacity
	clusDotOpacity := widget.NewSliderWithData(0, 255, clustOpacity)
	clusDotOpacity.Step = 1
	clusDotOpacity.OnChanged = func(v float64) {
		preference.SetFloat("clustOpacity", v)
	}

	// progress bar binding
	f := binding.NewFloat()

	content := container.NewVBox(
		gatename,
		widget.NewButton("Filter tables with gates", func() {
			// get the edges of all selected polygons
			alledges := e.drawSurface.alledges
			ch := make(chan bool, 2)
			f.Set(0.3) // progress bar
			go filterActiveGates(e, alledges, dataFiles, gatename.Text, a.Preferences(), ch)
			f.Set(0.6) // progress bar
			go saveGates(gatename.Text, e, ch)
			log.Println("plot done :", <-ch)
			log.Println("plot saved :", <-ch)
			time.Sleep(1 * time.Second)
			f.Set(0.) // reset progress bar
		}),
		widget.NewButton("Clear last gate", func() {
			clearLastGate(e)
		}),
		widget.NewButton("Clear all gates", func() {
			initGates(e)
		}),
		widget.NewButton("Screen shot", func() {
			screenShot(w, gatename.Text)
		}),
		widget.NewButton("Save HR image", func() {
			f.Set(0.3) // progress bar
			ch := make(chan bool, 2)
			go saveimage(w, gatename.Text, ch)
			log.Println("image saved :", <-ch)
			f.Set(1) // progress bar
			time.Sleep(1 * time.Second)
			f.Set(0.) // reset progress bar
		}),
		widget.NewButton("plot", func() {
			// get the edges of all selected polygons
			alledges := e.drawSurface.alledges
			plot.Plotform(a, w, e.zoom, header, firstTable, alledges, f)
		}),
		widget.NewButton("Show Clusters", func() {
			drawClusters(a, e, header, firstTable, f)
			time.Sleep(1 * time.Second)
			f.Set(0.) // reset progress bar
		}),
		widget.NewButton("Clear Clusters", func() {
			clearCluster(e)
		}),
		clusDotOpacity,
		widget.NewButton("Preferences", func() {
			pref.BuildPref(a, header)
		}),
		// zoom : very important : never unzom under the window size
		// in that case the image size = window size and zoom factor is wrong !
		newZoom(e, a),
		widget.NewButton("Exit", func() {
			os.Exit(0)
		}),
		widget.NewProgressBarWithData(f),
	)

	w2.SetContent(content)

}

// clear last gate on draw surface and init all edges
func clearLastGate(e *Editor) {
	e.drawSurface.clearPolygon(e.drawSurface.gatesLines)
	e.gateContainer.Refresh()
	initLastedges(e) // reset last edges and all points
	//initAlledges(e) // reset alledges
}

// save the gates to csv files
func saveGates(gateName string, e *Editor, ch chan bool) {
	fmt.Println("save gates")
	for i, poly := range e.drawSurface.alledges {
		if len(poly) < 3 {
			continue
		}
		fmt.Println(i, " ", poly)
		out := strconv.Itoa(i) + "_" + gateName
		writeCSV(out, poly)
	}
	ch <- true
}

// save screenshot of image to file
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

// save HR image to file
// credits https://www.devdungeon.com/content/working-images-go
func saveimage(w fyne.Window, filename string, ch chan bool) {
	c := w.Content().(*container.Scroll).Content
	out := software.Render(c, theme.DarkTheme())

	path := "plots/" + filename + ".png"
	outputFile, err := os.Create(path)
	if err != nil {
		log.Println("The image cannot be saved to the file")
	}
	png.Encode(outputFile, out)
	log.Println("Saving image to ", path)

	outputFile.Close()
	ch <- true
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

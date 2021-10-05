package ui

import (
	"encoding/csv"
	"fmt"
	"image/png"
	"lasso/src/filter"
	"lasso/src/plot"
	"lasso/src/pref"

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
func BuildTools(a fyne.App, w fyne.Window, e *Editor) {
	w2 := fyne.CurrentApp().NewWindow("Tool Box")
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

	// import column file index. This index is =0 at the beginning and then incremented by 1
	iCellFI := binding.BindPreferenceInt("imported file index", preference)
	impCellFindex, _ := iCellFI.Get()

	// progress bar binding
	f := binding.NewFloat()

	content := container.NewVBox(
		gatename,
		widget.NewButton("Filter tables with gates", func() {
			// get the edges of all selected polygons
			alledges := e.drawSurface.alledges
			go filterActiveGates(e, alledges, dataFiles, gatename.Text, a.Preferences(), f)
			go saveGates(gatename.Text, e)
		}),
		widget.NewButton("Save gates", func() {
			go saveGates(gatename.Text, e)
		}),
		widget.NewButton("Clear last gate", func() {
			clearLastGate(e)
		}),
		widget.NewButton("Clear all gates", func() {
			initGates(e)
		}),
		widget.NewButton("Import gates", func() {
			go importGates(e, f)
		}),
		widget.NewButton("Screen shot", func() {
			//f.Set(0.3) // progress bar
			go screenShot(w, gatename.Text, f)
			//f.Set(0.) // reset progress bar
		}),
		widget.NewButton("Save zoomed image", func() {
			go startSaveImage(w, gatename.Text, f)
		}),
		widget.NewButton("plot", func() {
			// get the edges of all selected polygons
			alledges := e.drawSurface.alledges
			plot.Plotform(a, w, e.zoom, header, firstTable, alledges, f)
		}),
		widget.NewButton("Show Clusters", func() {
			go drawClusters(a, e, header, firstTable, f)
		}),
		widget.NewButton("Clear Clusters/Expression", func() {
			clearCluster(e)
		}),
		widget.NewLabel("Dots Opacity [0-100%] :"),
		clusDotOpacity,
		widget.NewButton("Show Expression", func() {
			go buttonDrawExpress(a, e, preference, f, header, firstTable)
			//f.Set(0.) // reset progress bar
		}),
		widget.NewButton("Import cells", func() {
			go buttonImportCells(a, e, preference, iCellFI, f, impCellFindex, header, firstTable)
		}),
		widget.NewButton("Compare gates", func() {
			// map that store the check boxes state
			headerMap := make(map[string]interface{}, len(header[1:]))
			buildMapTrue(header[1:], headerMap)

			buttonCompare(a, e, preference, f, header, headerMap, firstTable)
		}),
		widget.NewButton("Preferences", func() {
			pref.BuildPref(a, header)
		}),
		// zoom : very important : never unzom under the window size
		// in that case the image size = window size and zoom factor is wrong !
		newZoom(e, a, f),
		widget.NewButton("Exit", func() {
			os.Exit(0)
		}),
		widget.NewProgressBarWithData(f),
	)

	w2.SetContent(content)
	w2.Show()
}

// clear last gate on draw surface and init all edges
func clearLastGate(e *Editor) {
	nob := len(e.gateContainer.Objects) - 1                  // nb of objects in the gate container
	e.gateContainer.Objects = e.gateContainer.Objects[0:nob] // remove last object = last gate name

	e.drawSurface.clearPolygon(e.drawSurface.gatesLines)

	e.gateContainer.Refresh()
	initLastedges(e)   // reset last edges and all points
	initLastGatesNB(e) // clear last gate number coordinates and decrease gateNB
	//initAlledges(e) // reset alledges
}

// save the gates to csv files withe ImageJ format and 100% zoom
// X,Y
// 131,150
// 105,189
// 156,187
func saveGates(gateName string, e *Editor) {
	log.Println("save gates")

	gateName = filter.FormatOutFile("gate", gateName, "") // test if name exist, if not, build a file name with the current time

	zoomFactor := 100. / float64(e.zoom)
	for i, poly := range e.drawSurface.alledges {
		if len(poly) < 3 {
			continue
		}

		out := strconv.Itoa(i) + "_" + gateName
		writeCSV(out, filter.ZoomPolygon(poly, zoomFactor))

	}
}

// import the gates in csv files withe ImageJ format and 100% zoom into the e.drawSurface.alledges
func importGates(e *Editor, f binding.Float) {
	f.Set(0.3)
	// clear all gates
	initGates(e)
	dir := "import_gates"
	gateFiles := filter.ListFiles(dir)
	for _, file := range gateFiles {
		gate := filter.ZoomPolygon(filter.ReadGate(dir, file), float64(e.zoom)/100.) // import the gate file and apply current zoom to polygon coordinates
		//fmt.Println("gate zoomed:", gate)
		e.drawSurface.alledges = append(e.drawSurface.alledges, gate)
		redrawpolygon(e.drawSurface, gate)
	}
	e.gateContainer.Refresh()
	f.Set(0.)
}

// save screenshot of image to file
// credits https://www.devdungeon.com/content/working-images-go
func screenShot(w fyne.Window, filename string, f binding.Float) {
	f.Set(0.3)                                                  // progress bar
	filename = filter.FormatOutFile("screenshot", filename, "") // test if name exist, if not, build a file name with the current time
	out := w.Canvas().Capture()
	f.Set(0.5) // progress bar
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

	// close files
	outputFile.Close()
	f.Set(0.) // progress bar
}

// start the save image goroutine
func startSaveImage(w fyne.Window, filename string, f binding.Float) {
	f.Set(0.3) // progress bar
	//ch := make(chan bool, 1)
	//go saveimage(w, filename, ch)
	go saveimage(w, filename, f)
	//log.Println("image saved :", <-ch)
}

// save HR image to file
// credits https://www.devdungeon.com/content/working-images-go
// func saveimage3(w fyne.Window, filename string, ch chan bool) {

// 	c := w.Content().(*container.Scroll).Content
// 	out := software.Render(c, theme.DarkTheme())

// 	path := "plots/" + filename + ".png"
// 	outputFile, err := os.Create(path)
// 	if err != nil {
// 		log.Println("The image cannot be saved to the file")
// 	}
// 	err = png.Encode(outputFile, out)
// 	if err != nil {
// 		log.Println("png encoding error : ", err)
// 	}
// 	log.Println("Saving image to ", path)

// 	outputFile.Close()
// 	ch <- true
// }

// save HR image to file
// credits https://www.devdungeon.com/content/working-images-go
func saveimage(w fyne.Window, filename string, f binding.Float) {
	log.Print("Saving image...")
	c := w.Content().(*container.Scroll).Content
	out := software.Render(c, theme.DarkTheme())

	path := "plots/" + filename + ".png"
	outputFile, err := os.Create(path)
	if err != nil {
		log.Println("The image cannot be saved to the file")
	}
	f.Set(0.5) // progress bar
	err = png.Encode(outputFile, out)
	if err != nil {
		log.Println("png encoding error : ", err)
	}
	log.Println("image saved to ", path)
	outputFile.Close()
	f.Set(0.) // progress bar
}

// write polygon edge to file in ImageJ format :
// X,Y
// 131,150
// 105,189
// 156,187
func writeCSV(filename string, poly []filter.Point) {
	path := "gates/" + filename + ".csv"
	file, err := os.Create(path)
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Comma = ','

	// write header
	err = writer.Write([]string{"X", "Y"})
	checkError("Cannot write to file "+filename, err)

	for _, value := range poly {
		err := writer.Write(filterPtToStr(value))
		checkError("Cannot write to file "+filename, err)
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

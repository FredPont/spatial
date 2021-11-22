package ui

import (
	"encoding/csv"
	"fmt"
	"image/png"
	"spatial/src/filter"
	"spatial/src/plot"
	"spatial/src/pref"

	//"spatial/src/plot"

	"log"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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
	w2 := fyne.CurrentApp().NewWindow("scSpatial Explorer")
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

	// brush buttons
	brush1, brush2 := brushesButtons(e, a)

	// icons
	iconFilt := "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" version=\"1.1\" width=\"24\" height=\"24\" fill=\"#de8159\" viewBox=\"0 0 24 24\"><path d=\"M15,19.88C15.04,20.18 14.94,20.5 14.71,20.71C14.32,21.1 13.69,21.1 13.3,20.71L9.29,16.7C9.06,16.47 8.96,16.16 9,15.87V10.75L4.21,4.62C3.87,4.19 3.95,3.56 4.38,3.22C4.57,3.08 4.78,3 5,3V3H19V3C19.22,3 19.43,3.08 19.62,3.22C20.05,3.56 20.13,4.19 19.79,4.62L15,10.75V19.88M7.04,5L11,10.06V15.58L13,17.58V10.05L16.96,5H7.04Z\" /></svg>"
	iconFilter := fyne.NewStaticResource("filter", []byte(iconFilt))

	content := container.NewVBox(
		logo(),
		gatename,
		container.NewHBox(brush1,
			brush2,
			// screenshot
			widget.NewButtonWithIcon("", theme.MediaPhotoIcon(), func() {
				go screenShot(w, gatename.Text, f)
			}),
			// preferences
			widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
				go pref.BuildPref(a, header)
			}),
			widget.NewButtonWithIcon("Exit", theme.LogoutIcon(), func() {
				a.Quit()
				//os.Exit(0)
			}),
		),
		widget.NewButtonWithIcon("Filter tables by gates", iconFilter, func() {
			// get the edges of all selected polygons
			alledges := e.drawSurface.alledges
			go filterActiveGates(e, alledges, dataFiles, gatename.Text, a.Preferences(), f)
			go saveGates(gatename.Text, e)
		}),
		widget.NewButtonWithIcon("Save Gates", theme.DocumentSaveIcon(), func() {
			go saveGates(gatename.Text, e)
		}),
		widget.NewButtonWithIcon("Import gates", theme.FolderOpenIcon(), func() {
			go importGates(e, f)
		}),
		container.NewHBox(
			widget.NewLabel("Clear :"),
			widget.NewButton("last gate", func() {
				go clearLastGate(e)
			}),
			widget.NewButton("all gates", func() {
				go initGates(e)
			}),
		),
		widget.NewButton("Compare gates", func() {
			go showCompareWindow(a, e, preference, f, header, firstTable)
		}),

		widget.NewButton("Plot gates", func() {
			// get the edges of all selected polygons
			alledges := e.drawSurface.alledges
			go plot.Plotform(a, w, e.zoom, header, firstTable, alledges, f)
		}),
		widget.NewButton("2D Plot", func() {
			// get the edges of all selected polygons
			alledges := e.drawSurface.alledges
			go Plot2Dform(a, e, w, e.zoom, header, firstTable, alledges, f)
		}),
		container.NewHBox(
			//widget.NewLabel("Show :"),
			widget.NewButton("Show Clusters", func() {
				go drawClusters(a, e, header, firstTable, f)
			}),
			widget.NewButton("Expression", func() {
				go buttonDrawExpress(a, e, preference, f, header, firstTable)
				//f.Set(0.) // reset progress bar
			}),
		),
		widget.NewButton("Clear Clusters/Expression", func() {
			go clearCluster(e)
		}),
		widget.NewButton("Save zoomed image", func() {
			go startSaveImage(w, gatename.Text, f)
		}),
		widget.NewLabel("Dots Opacity [0-100%] :"),
		clusDotOpacity,

		widget.NewButton("Import cells", func() {
			go buttonImportCells(a, e, preference, iCellFI, f, impCellFindex, header, firstTable)
		}),

		// zoom : very important : never unzom under the window size
		// in that case the image size = window size and zoom factor is wrong !
		newZoom(e, a, f),

		widget.NewProgressBarWithData(f),
	)

	w2.SetContent(content)
	w2.Show()
}

// logo display a log in tool window
func logo() fyne.CanvasObject {
	img := canvas.NewImageFromFile("src/ui/logo.png")
	img.SetMinSize(fyne.Size{Width: 171, Height: 55})
	img.FillMode = canvas.ImageFillContain
	return img
}

// clear last gate on draw surface and init all edges
func clearLastGate(e *Editor) {
	//log.Println("nb de lignes", e.drawSurface.gatesLines)
	if e.drawSurface.gatesNumbers.nb < 2 || len(e.drawSurface.gatesLines) < 2 { // less than 2 gates
		return
	}
	//log.Println("e.drawSurface.gatesNumbers.nb=", e.drawSurface.gatesNumbers.nb)
	nob := len(e.gateNumberContainer.Objects) - 1 // nb of objects in the gate nb container minus the last one
	log.Println("nb de gate names", nob)

	e.gateNumberContainer.Objects = e.gateNumberContainer.Objects[:nob] // remove last object = last gate name

	e.drawSurface.clearPolygon(e.drawSurface.gatesLines)

	initLastedges(e)   // reset last edges and all points
	initLastGatesNB(e) // clear last gate number coordinates and decrease gateNB
	e.gateContainer.Refresh()
}

// save the gates to csv files withe ImageJ format and 100% zoom
// X,Y
// 131,150
// 105,189
// 156,187
func saveGates(gateName string, e *Editor) {

	gateName = filter.FormatOutFile("gate", gateName, "") // test if name exist, if not, build a file name with the current time

	zoomFactor := 100. / float64(e.zoom)
	for i, poly := range e.drawSurface.alledges {
		if len(poly) < 3 {
			continue
		}

		out := strconv.Itoa(i) + "_" + gateName
		writeCSV(out, filter.ZoomPolygon(poly, zoomFactor))
		log.Println("gate saved in gates/", out)
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
	drawImportedGatesNB(e.drawSurface) // draw and store the gates numbers coordinates after import gate
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

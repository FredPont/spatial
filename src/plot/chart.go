package plot

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"io"
	"lasso/src/filter"
	"log"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func check(e error) {
	if e != nil {
		log.Println("plot error !", e)
	}
}

func makeplot(a fyne.App, header []string, filename, colX, colY, plotName, bkgDotSize string, alledges [][]filter.Point) {
	// index of the 2 columns to plot and the XY columns with the image coordinates (to be able to filter the gates)
	// get parameters from preferences
	param := prefToConf(a.Preferences())
	// index of the 2 columns to plot and the XY columns with the image coordinates (to be able to filter the gates)
	colIndexes := filter.GetColIndex(header, []string{colX, colY, param.X, param.Y})
	mapAndGates := filter.ReadColumns(filename, colIndexes)

	// get the two first columns of mapAndGates to get map coordinates
	scatterData := strToplot(extract2cols(mapAndGates, 0, 1))
	// extract dots in all gates
	alldotsInGates := extractGateDots(a, mapAndGates, alledges, colX, colY)
	fmt.Println(alldotsInGates)
	mapDotSize, _ := vg.ParseLength(bkgDotSize)
	makeScatter(a, alldotsInGates, scatterData, mapDotSize, plotName, colX, colY, plotName)

	// display plot on new window
	plotWindow := a.NewWindow("Plot")
	img := canvas.NewImageFromFile("plots/" + plotName + ".png")
	plotWindow.SetContent(img)
	plotWindow.Resize(fyne.NewSize(800, 800))
	plotWindow.SetFixedSize(true)
	plotWindow.Show()

}

func extractGateDots(a fyne.App, tableXYxy [][]string, alledges [][]filter.Point, colX, colY string) [][][]string {

	var allXY [][][]string // all xy coordinates of dots in all gates
	param := prefToConf(a.Preferences())
	gateNB := len(alledges)
	log.Println("all edges", alledges)
	ch1 := make(chan [][]string, gateNB) // ch1 store the xy coordinates of dots in one gate
	for _, polygon := range alledges {
		go filter.TablePlot(tableXYxy, polygon, param, colX, colY, ch1)
	}
	for i := 0; i < gateNB; i++ {
		msg := <-ch1
		allXY = append(allXY, msg)
	}
	//log.Println("all gates gates", allXY)
	return allXY
}

func makeScatter(a fyne.App, alldotsInGates [][][]string, scatterData plotter.XYs, dotsize vg.Length, title, xaxisName, yaxisName, plotName string) {
	mapR, mapG, mapB, mapA := getPrefColorRGBA(a, "unselR", "unselG", "unselB", "unselA")

	// Create a new plot, set its title and
	// axis labels.
	p := plot.New()

	p.Title.Text = title
	p.X.Label.Text = xaxisName
	p.Y.Label.Text = yaxisName
	//p.HideAxes()
	// Draw a grid behind the data
	//p.Add(plotter.NewGrid())

	// Make a scatter plotter and set its style.
	s, err := plotter.NewScatter(scatterData)
	if err != nil {
		panic(err)
	}
	s.GlyphStyle.Shape = draw.CircleGlyph{}
	s.GlyphStyle.Radius = dotsize
	s.GlyphStyle.Color = color.RGBA{R: uint8(mapR), G: uint8(mapG), B: uint8(mapB), A: uint8(mapA)} // background dots color
	p.Add(s)

	// add new points
	showGates(alldotsInGates, p, dotsize)
	// addPoints(scatterData[30:100], p, 5, color.RGBA{R: 75, G: 0, B: 130, A: 255})
	// addPoints(scatterData[300:500], p, 3, color.RGBA{R: 0, G: 150, B: 255, A: 255})
	// addPoints(scatterData[600:800], p, 4, color.RGBA{R: 255, G: 150, B: 30, A: 255})

	savePlot(p, 800, 800, "plots/"+plotName+".png")
}

func showGates(alldotsInGates [][][]string, p *plot.Plot, dotsize vg.Length) {
	for _, gate := range alldotsInGates {
		scatterData := strToplot(gate)
		addPoints(scatterData, p, dotsize, color.RGBA{R: 75, G: 0, B: 130, A: 255})
	}

}

func getPrefColorRGBA(a fyne.App, R, G, B, A string) (int, int, int, int) {
	pref := a.Preferences()
	// map dots color - read RGBA from preferences
	unselR := binding.BindPreferenceInt("unselR", pref)
	mapR, e := unselR.Get()
	check(e)
	unselG := binding.BindPreferenceInt("unselG", pref)
	mapG, e := unselG.Get()
	check(e)
	unselB := binding.BindPreferenceInt("unselB", pref)
	mapB, e := unselB.Get()
	check(e)
	unselA := binding.BindPreferenceInt("unselA", pref)
	mapA, e := unselA.Get()
	check(e)
	log.Println("color in pref", mapR, mapG, mapB, mapA)
	return mapR, mapG, mapB, mapA
}

func savePlot(p *plot.Plot, w, h vg.Length, filename string) {
	// Save the plot to a PNG file.
	// if err := p.Save(4*vg.Inch, 4*vg.Inch, "scatter.png"); err != nil {
	// 	panic(err)
	// }
	if err := p.Save(w, h, filename); err != nil {
		panic(err)
	}
}

func addPoints(pts plotter.XYs, p *plot.Plot, dotsize vg.Length, clr color.Color) {
	s2, err := plotter.NewScatter(pts)
	if err != nil {
		panic(err)
	}
	s2.GlyphStyle.Shape = draw.CircleGlyph{}
	s2.GlyphStyle.Radius = dotsize
	s2.GlyphStyle.Color = clr
	p.Add(s2)
}

func readscv() [][]string {
	var xy [][]string
	// Open the file
	csvfile, err := os.Open("umap.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)
	r.Comma = '\t'
	//r := csv.NewReader(bufio.NewReader(csvfile))

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		xy = append(xy, []string{record[0], record[1]})
	}
	return xy
}

func strToplot(xy [][]string) plotter.XYs {
	n := len(xy)
	pts := make(plotter.XYs, n)
	for i := 0; i < n; i++ {
		pts[i].X = strFloat(xy[i][0])
		pts[i].Y = strFloat(xy[i][1])

	}
	return pts
}

func strFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func extract2cols(cols [][]string, idx1, idx2 int) [][]string {
	var two [][]string
	for i := 0; i < len(cols); i++ {
		two = append(two, []string{cols[i][idx1], cols[i][idx2]})
	}
	return two
}

// PrefToConf retreive conf data from fyne pref
func prefToConf(pref fyne.Preferences) filter.Conf {
	// get // X coordinates
	xcor := binding.BindPreferenceString("xcor", pref) // set the link to preferences for rotation
	x, _ := xcor.Get()

	// get y coordinates
	ycor := binding.BindPreferenceString("ycor", pref) // set the link to preferences for rotation
	y, _ := ycor.Get()

	// get scaling factor
	sf := binding.BindPreferenceFloat("scaleFactor", pref) // set the link to preferences for scaling factor
	scale, _ := sf.Get()

	// get coordinates +90° rotation : necessary for 10x Genomics
	r := binding.BindPreferenceBool("rotate", pref) // set the link to preferences for rotation
	rotate, _ := r.Get()

	return filter.Conf{x, y, scale, rotate}

}

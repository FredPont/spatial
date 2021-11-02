package plot

import (
	"encoding/csv"
	"image/color"
	"io"
	"lasso/src/filter"
	"log"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"github.com/mazznoer/colorgrad"
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

func makeplot(a fyne.App, zoom int, header []string, filename, colX, colY, plotName, bkgDotSize string, alledges [][]filter.Point, f binding.Float) {
	f.Set(0.3) // progress bar
	// index of the 2 columns to plot and the XY columns with the image coordinates (to be able to filter the gates)
	// get parameters from preferences
	param := prefToConf(a.Preferences())
	// index of the 2 columns to plot and the XY columns with the image coordinates (to be able to filter the gates)
	colIndexes := filter.GetColIndex(header, []string{colX, colY, param.X, param.Y})
	mapAndGates := filter.ReadColumns(filename, colIndexes)

	// get the two first columns of mapAndGates to get map coordinates
	scatterData := strToplot(extract2cols(mapAndGates, 0, 1))
	// extract dots in all gates
	alldotsInGates := extractGateDots(a, zoom, mapAndGates, alledges, colX, colY)
	if len(alldotsInGates) < 1 {
		log.Println("Plot canceled or no dots in gate !")
		return
	}
	//fmt.Println(alldotsInGates)
	mapDotSize, _ := vg.ParseLength(bkgDotSize)
	f.Set(0.5) // progress bar
	makeScatter(a, alldotsInGates, scatterData, mapDotSize, plotName, colX, colY, plotName)
	f.Set(0.8) // progress bar
	// display plot on new window
	//plotWindow := a.NewWindow("Plot")
	plotWindow := fyne.CurrentApp().NewWindow("Plot")
	img := canvas.NewImageFromFile("plots/" + plotName + ".png")
	f.Set(0.) // reset progress bar
	plotWindow.SetContent(img)
	plotWindow.Resize(fyne.NewSize(800, 800))
	plotWindow.SetFixedSize(true)
	plotWindow.Show()

}

func extractGateDots(a fyne.App, zoom int, tableXYxy [][]string, alledges [][]filter.Point, colX, colY string) [][][]string {

	var allXY [][][]string // all xy coordinates of dots in all gates
	param := prefToConf(a.Preferences())
	gateNB := len(alledges)
	//log.Println("all edges", alledges)
	ch1 := make(chan [][]string, gateNB) // ch1 store the xy coordinates of dots in one gate
	for _, polygon := range alledges {
		go filter.TablePlot(zoom, tableXYxy, polygon, param, colX, colY, ch1)
	}
	for i := 0; i < gateNB; i++ {
		msg := <-ch1
		if len(msg) > 0 {
			allXY = append(allXY, msg)
		}
	}
	//log.Println("all gates gates", allXY)
	return allXY
}

func makeScatter(a fyne.App, alldotsInGates [][][]string, scatterData plotter.XYs, dotsize vg.Length, title, xaxisName, yaxisName, plotName string) {
	mapR, mapG, mapB, mapA := GetPrefColorRGBA(a, "unselR", "unselG", "unselB", "unselA")

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
	showGates(a, alldotsInGates, p, dotsize)
	// addPoints(scatterData[30:100], p, 5, color.RGBA{R: 75, G: 0, B: 130, A: 255})
	// addPoints(scatterData[300:500], p, 3, color.RGBA{R: 0, G: 150, B: 255, A: 255})
	// addPoints(scatterData[600:800], p, 4, color.RGBA{R: 255, G: 150, B: 30, A: 255})

	savePlot(p, 800, 800, "plots/"+plotName+".png")
}

func showGates(a fyne.App, alldotsInGates [][][]string, p *plot.Plot, dotsize vg.Length) {
	nbGates := len(alldotsInGates)
	for i, gate := range alldotsInGates {
		scatterData := strToplot(gate)
		if nbGates == 1 {
			gatedotsR, gatedotsG, gatedotsB, gatedotsA := GetPrefColorRGBA(a, "gateDotsR", "gateDotsG", "gateDotsB", "gateDotsA")
			addPoints(scatterData, p, dotsize, color.RGBA{R: uint8(gatedotsR), G: uint8(gatedotsG), B: uint8(gatedotsB), A: uint8(gatedotsA)})
		} else {
			addPoints(scatterData, p, dotsize, dotColors(nbGates, i))
		}

		//
	}

}

// dotColors computes the color of scatter dots
// for a total number of clusters "nbGates"
func dotColors(nbGates, gateIndex int) color.RGBA {
	grad := colorgrad.Rainbow().Sharp(uint(nbGates+1), 0.2)
	return rgbaModel(grad.Colors(uint(nbGates + 1))[gateIndex])
}

func rgbaModel(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

// GetPrefColorRGBA get the R, G, B, A values from preferences
func GetPrefColorRGBA(a fyne.App, R, G, B, A string) (int, int, int, int) {
	pref := a.Preferences()
	// map dots color - read RGBA from preferences
	dotsR := binding.BindPreferenceInt(R, pref)
	mapR, e := dotsR.Get()
	check(e)
	dotsG := binding.BindPreferenceInt(G, pref)
	mapG, e := dotsG.Get()
	check(e)
	dotsB := binding.BindPreferenceInt(B, pref)
	mapB, e := dotsB.Get()
	check(e)
	dotsA := binding.BindPreferenceInt(A, pref)
	mapA, e := dotsA.Get()
	check(e)
	//log.Println("color in pref", mapR, mapG, mapB, mapA)
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

package plot

import (
	"encoding/csv"
	"image/color"
	"io"
	"lasso/src/filter"
	"log"
	"math/rand"
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

func makeplot(a fyne.App, header []string, filename, colX, colY, plotName, bkgDotSize string) {
	colIndexes := filter.GetColIndex(header, []string{colX, colY})
	xy := filter.ReadColumns(filename, colIndexes)
	scatterData := strToplot(xy)
	mapDotSize, _ := vg.ParseLength(bkgDotSize)
	makeScatter(a, scatterData, mapDotSize, plotName, colX, colY, plotName)

	// display plot on new window
	plotWindow := a.NewWindow("Plot")
	img := canvas.NewImageFromFile("plots/" + plotName + ".png")
	plotWindow.SetContent(img)
	plotWindow.Resize(fyne.NewSize(800, 800))
	plotWindow.SetFixedSize(true)
	plotWindow.Show()

}

func makeScatter(a fyne.App, scatterData plotter.XYs, dotsize vg.Length, title, xaxisName, yaxisName, plotName string) {
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
	addPoints(scatterData[30:100], p, 5, color.RGBA{R: 75, G: 0, B: 130, A: 255})
	addPoints(scatterData[300:500], p, 3, color.RGBA{R: 0, G: 150, B: 255, A: 255})
	addPoints(scatterData[600:800], p, 4, color.RGBA{R: 255, G: 150, B: 30, A: 255})

	savePlot(p, 800, 800, "plots/"+plotName+".png")
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

// randomPoints returns some random x, y points.
func randomPoints(n int) plotter.XYs {
	pts := make(plotter.XYs, n)
	for i := range pts {
		if i == 0 {
			pts[i].X = rand.Float64()
		} else {
			pts[i].X = pts[i-1].X + rand.Float64()
		}
		pts[i].Y = pts[i].X + 10*rand.Float64()
	}
	return pts
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

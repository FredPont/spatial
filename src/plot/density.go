package plot

import (
	"spatial/src/filter"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// Density split data into n classes intervals
func Density(data []float64, n float64) plotter.XYs {

	var pts plotter.XYs
	//copy(s2, data)
	min, max := filter.FindMinAndMax(data)
	step := (max - min) / n
	dataLen := len(data)

	sma20 := sma(20) // smoothing by moving average on 20 points

	for i := min + step; i <= max; i += step {
		count := 0 // occurences in intervall
		for _, nb := range data {
			if nb >= i-step && nb < i {
				count++
			}
			if nb == max && i == max {
				count++
			}
		}
		if dataLen > 100 {
			pts = append(pts, plotter.XY{X: i - (step / 2.), Y: sma20(float64(count))}) // smoothing
		} else {
			pts = append(pts, plotter.XY{X: i - (step / 2.), Y: float64(count)}) // no smooting
		}

	}

	return pts
}

func makeDensityPlot(pts plotter.XYs, expcol string) {
	p := plot.New()

	//p.Title.Text = ""
	p.X.Label.Text = expcol
	p.Y.Label.Text = "Abundance"
	//p.BackgroundColor = color.Transparent

	err := plotutil.AddLines(p,
		"", pts)
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(5*vg.Inch, 5*vg.Inch, "temp/density.png"); err != nil {
		panic(err)
	}

}

func BuildDensity(data []float64, n float64, expcol string, ExpressWindow fyne.Window) {
	pts := Density(data, n)
	makeDensityPlot(pts, expcol)
	ExpressWindow.Content().Refresh()
}

// credit : https://rosettacode.org/wiki/Averages/Simple_moving_average#Go
// moving average smoothing function
func sma(period int) func(float64) float64 {
	var i int
	var sum float64
	var storage = make([]float64, 0, period)

	return func(input float64) (avrg float64) {
		if len(storage) < period {
			sum += input
			storage = append(storage, input)
		}

		sum += input - storage[i]
		storage[i], i = input, (i+1)%period
		avrg = sum / float64(len(storage))

		return
	}
}

// DensityPicture display a density plot in expression tool window
func DensityPicture() fyne.CanvasObject {
	img := canvas.NewImageFromFile("temp/density.png")
	img.SetMinSize(fyne.Size{Width: 350, Height: 350})
	img.FillMode = canvas.ImageFillContain
	return img
}

package plot

import (
	"spatial/src/filter"

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

	sma20 := sma(20) // smoothing by moving average on 10 points
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
		pts = append(pts, plotter.XY{X: i - (step / 2.), Y: sma20(float64(count))})

	}

	return pts
}

func makeDensityPlot(pts plotter.XYs) {
	p := plot.New()

	p.Title.Text = "Distribution"
	p.X.Label.Text = "Expression"
	p.Y.Label.Text = "Abundance"

	err := plotutil.AddLines(p,
		"", pts)
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "temp/density.png"); err != nil {
		panic(err)
	}

}

func BuildDensity(data []float64, n float64) {
	pts := Density(data, n)
	makeDensityPlot(pts)
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

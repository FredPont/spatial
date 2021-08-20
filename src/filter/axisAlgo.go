package filter

import "math"

// credit : https://stackoverflow.com/questions/326679/choosing-an-attractive-linear-scale-for-a-graphs-y-axis

// TicInterval creates the Y axis values for a graph.
//
// Calculate Min amd Max graphical labels and graph
// increments.  The number of ticks defaults to
// 10 which is the SUGGESTED value.  Any tick value
// entered is used as a suggested value which is
// adjusted to be a 'pretty' value.
//
// Output will be an array of the Y axis values that
// encompass the Y values.
func TicInterval(yMin, yMax float64, ticks int) []float64 {

	var result []float64
	// If yMin and yMax are identical, then
	// adjust the yMin and yMax values to actually
	// make a graph. Also avoids division by zero errors.
	if yMin == yMax {
		yMin = yMin - 10. // some small value
		yMax = yMax + 10. // some small value
	}
	// Determine Range
	yrange := yMax - yMin
	// Adjust ticks if needed
	if ticks < 2 {
		ticks = 2
	} else if ticks > 2 {
		ticks -= 2
	}

	// Get raw step value
	tempStep := yrange / float64(ticks)
	// Calculate pretty step value
	mag := math.Floor(math.Log10(tempStep))
	magPow := math.Pow(10, mag)
	magMsd := math.Round(tempStep/magPow + 0.5)
	stepSize := magMsd * magPow

	// build Y label array.
	// Lower and upper bounds calculations
	lb := stepSize * math.Floor(yMin/stepSize)
	ub := stepSize * math.Ceil(yMax/stepSize)
	// Build array
	val := lb
	for {
		// exclude value outside Y interval
		if Between(yMin, yMax, val) {
			result = append(result, val)
		}

		val += stepSize
		if val > ub {
			break
		}

	}
	return result
}

// TicPixelPos convert Ticks intervals to integers
func TicPixelPos(yMin, yMax float64, ticks int) []int {
	var pixelPos []int
	res := TicInterval(yMin, yMax, ticks)
	for _, v := range res {
		pixelPos = append(pixelPos, int(math.Round(v)))
	}
	return pixelPos
}

// Between test if Xmin <= x <= Xmax
func Between(Xmin, Xmax, x float64) bool {
	if Xmin <= x && x <= Xmax {
		return true
	}
	return false
}

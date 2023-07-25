package resampling

import (
	"bufio"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne/v2/data/binding"
)

type ColMinMax struct {
	Colname string
	Min     int
	Max     int
}

// Resample write one line / n in another file and  convert to Int the columns corresponding to indexes
func Resample(fileIn, fileOut string, n int, indexes []int, colnames []string, delim string, f binding.Float) {
	log.Println("Resampling ", fileIn)
	f.Set(0.3) // progress bar

	// min max values for columns to convert to int
	colVals := make([]ColMinMax, len(indexes))
	initColMinMax(colVals, colnames)

	// Open input file for reading
	inputFile, err := os.Open(fileIn)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	// Create output file for writing
	outputFile, err := os.Create(fileOut)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	// Create scanner to iterate over lines of input file
	scanner := bufio.NewScanner(inputFile)
	// write header
	scanner.Scan()
	outputFile.WriteString(scanner.Text() + "\n")
	i := 1
	for scanner.Scan() {
		// Write every third line to output file
		if i%n == 0 {
			row := scanner.Text()
			if len(indexes) > 0 {
				row = roundRowToInt(row, indexes, delim, &colVals)
			}
			outputFile.WriteString(row + "\n")
		}
		i++
	}
	log.Println("Resampling done in ", fileOut)
	log.Println("column Min/Max", colVals)
	f.Set(1) // progress bar
}

// roundRowToInt round the columns whith indexes number to the closest int and then convert int to string for file writing
func roundRowToInt(row string, indexes []int, delim string, colVals *[]ColMinMax) string {
	cells := strings.Split(row, delim)
	for j, i := range indexes {
		value := strToInt(cells[i])
		cells[i] = strconv.Itoa(value)
		// register min max for the columns selected for Int conversion
		if value < (*colVals)[j].Min {
			(*colVals)[j].Min = value
		}
		if value > (*colVals)[j].Max {
			(*colVals)[j].Max = value
		}
	}
	return strings.Join(cells, delim)
}

// strToInt convert a string to int and round to the nearest int
func strToInt(str string) int {
	// Convert string to float64
	x, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic(err)
	}

	// Round float to nearest integer
	rounded := math.Round(x)

	// Convert float to int
	return int(rounded)

}

// initialise the columns names and min / max
func initColMinMax(colVals []ColMinMax, colnames []string) {
	for i, name := range colnames {
		colVals[i] = ColMinMax{Colname: name, Min: math.MaxInt64, Max: math.MinInt64}
	}
}

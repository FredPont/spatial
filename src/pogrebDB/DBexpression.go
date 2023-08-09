package pogrebDB

import (
	"fmt"
	"log"
	"spatial/src/filter"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/akrylysov/pogreb"
)

// DBgetExpress get the values of the expression and XY columns from pogreb database
func DBgetExpress(a fyne.App, header []string, filename string, expcol string, curPathwayIndex binding.Int) ([]float64, []filter.Point) {
	pref := a.Preferences()
	// X coordinates
	xcor := binding.BindPreferenceString("xcor", pref) // set the link to preferences for X coordinates
	xc, _ := xcor.Get()

	// y coordinates
	ycor := binding.BindPreferenceString("ycor", pref) // set the link to preferences for y coordinates
	yc, _ := ycor.Get()

	colIndexes := filter.GetColIndex(header, []string{expcol, xc, yc}) // selected column + xy
	curPathwayIndex.Set(colIndexes[0])                                 // set the current expression Index to the selected column to enable button next/previous slide
	return DBReadExpress(a, filename, colIndexes)
}

// ReadExpress read only columns with positions in indexes and fill a a map
// of expression normalized between 0-1 => slice of x,y coordinates
func DBReadExpress(a fyne.App, filename string, colIndexes []int) ([]float64, []filter.Point) {

	// get scaleFactor and rotation from pref
	pref := a.Preferences()

	sf := binding.BindPreferenceFloat("scaleFactor", pref) // set the link to preferences for scaling factor
	scaleFactor, _ := sf.Get()                             // read the preference for scaling factor

	rot := binding.BindPreferenceString("rotate", pref) // set the link to preferences for rotation
	rotate, _ := rot.Get()

	// array of expression valules and xy coordinates scaled
	var expressions []float64
	var pts []filter.Point
	cxy := make([][]string, 3) // expression and XY columns
	// Open the database to read expression and XY columns
	dbname, _ := filter.RemExt(filename)
	db, err := pogreb.Open("temp/pogreb/"+dbname, nil)
	if err != nil {
		log.Println("database ", dbname, "not found in temp/pogreb/")
		log.Fatal(err)

	}
	defer db.Close()

	// read the 3 columns : expresion, x,y
	for i := 0; i < 3; i++ {
		cxy[i] = ReadColumn(db, fmt.Sprint(colIndexes[i]))
		//log.Println("column ", i, " value = ", cxy[i][0:5])
	}

	// Iterate through the records
	for ct := 0; ct < len(cxy[1]); ct++ {
		xScaled, yScaled := filter.ScaleXY(cxy[1][ct], cxy[2][ct], scaleFactor, rotate)
		exp, err := strconv.ParseFloat(cxy[0][ct], 64)
		if err != nil {
			//log.Println("column number", colIndexes[0]+1, "does not contain a number", err)
			continue
		}

		pts = append(pts, filter.Point{X: int(xScaled), Y: int(yScaled)})
		expressions = append(expressions, exp)
	}

	return expressions, pts
}

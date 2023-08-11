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

// ReadClusters read only columns with positions in indexes and fill a map
// cluster NB => slice of x,y coordinates
func ReadClustersDB(a fyne.App, filename string, colIndexes []int) map[int][]filter.Point {

	// get scaleFactor and rotation from pref
	pref := a.Preferences()

	sf := binding.BindPreferenceFloat("scaleFactor", pref) // set the link to preferences for scaling factor
	scaleFactor, _ := sf.Get()                             // read the preference for scaling factor

	rot := binding.BindPreferenceString("rotate", pref) // set the link to preferences for rotation
	rotate, _ := rot.Get()

	// map with cluster number => slice of xy coordinates scaled
	clusterMap := make(map[int][]filter.Point, 0)

	cxy := make([][]string, 3) // clusterNB and XY columns
	// Open the database to read expression and XY columns
	dbname, _ := filter.RemExt(filename)
	db, err := pogreb.Open("temp/pogreb/"+dbname, nil)
	if err != nil {
		log.Println("database ", dbname, "not found in temp/pogreb/")
		log.Fatal(err)

	}
	defer db.Close()
	log.Println("Read clusters values from database...")
	// read the 3 columns : clusterNB, x,y
	for i := 0; i < 3; i++ {
		cxy[i] = ReadColumn(db, fmt.Sprint(colIndexes[i]))
		//log.Println("column ", i, " value = ", cxy[i][0:5])
	}

	// Iterate through the records
	for ct := 0; ct < len(cxy[1]); ct++ {
		xScaled, yScaled := filter.ScaleXY(cxy[1][ct], cxy[2][ct], scaleFactor, rotate)
		clustNB, err := strconv.Atoi(cxy[0][ct])
		check(err)

		clusterMap[clustNB] = append(clusterMap[clustNB], filter.Point{int(xScaled), int(yScaled)})
	}

	return clusterMap
}

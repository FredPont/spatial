// vulcano inside gate vs all the cells outside the gate(s)
// count the number of cells inside / outside

// filter data with gates into 2 blocs : inside gates and outside
// and save the corresponding 2 CSV files

// parse the inside and outside gates csv files simultaneously to calculate FC and Pvalue.
// with the option "low memory" this is done column by column
// with no option, this is done for the whole data set concurently
package ui

import (
	"log"
	"spatial/src/filter"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

// start comparison
func compareGatevsAll(e *Editor, header []string, headerMap map[string]interface{}, pref fyne.Preferences, firstTable string, f binding.Float, g1Map, g2Map map[string]interface{}, outfile string) {
	f.Set(.3)
	param := prefToConf(pref) // get xy rotation zoom factor from pref
	// group1, group2  = polygone coordinates of gates in group 1 and 2
	group1, _ := gatesInGroup(e, g1Map, g2Map)
	// index and colnames of selected colum + XY coordinates for filtering
	index, colnames, XYindex := selColIndex(header, headerMap, pref) // colnames : cellsId, colum selected, X,Y
	//log.Println(index)
	// to save RAM each line of the table is filtered when reading the table
	table1, table2, test := filter.ReadGateVsAll(e.zoom, firstTable, index, XYindex, group1, param)
	if !test {
		log.Println("Error detected in XY coordinates ! Comparison aborted !")
		f.Set(0.)
		return
	}
	if len(table1) == 0 || len(table2) == 0 {
		log.Println("No data points in gate ! computation aborted !")
		f.Set(0.)
		return
	}
	f.Set(0.5)
	pvfcTable := foldChangePV(table1, table2, colnames)
	f.Set(0.8)
	// log.Println(colnames, group1, group2)
	// log.Println(table1)
	// log.Println(table2)
	// log.Println(pvfcTable)

	// save vulcano data to file
	fname := filter.FormatOutFile("comparison", outfile, ".csv")
	go writePV(fname, pvfcTable)
	// vulcano plot window
	go buildVulanoPlot(e, header, fname, firstTable, pvfcTable)
	// vulcano plot

	// readVulcano(fname, pvfcTable)
	// log.Println(readVulcano(fname, pvfcTable))
	// buildVulcWin()
	f.Set(0.)
}

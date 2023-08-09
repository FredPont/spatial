/*
 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.

 Written by Frederic PONT.
 (c) Frederic Pont 2022
*/

package pogrebDB

import (
	"bufio"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"spatial/src/filter"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/akrylysov/pogreb"
	"github.com/schollz/progressbar/v3"
)

func check(e error) {
	if e != nil {
		log.Println(e)
	}
}

////////////////////////////////////////
// read all csv in RAM !
////////////////////////////////////////

// ReadAll read a CSV file in RAM
func ReadAll(db *pogreb.DB, path string) []string {

	//var col []string
	csvFile, err := os.Open(path)
	check(err)
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = '\t'
	// read header
	header, err := reader.Read()

	data, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	loadData(db, data, header)

	fmt.Println(path, " done !")

	return header
}

// insertcol inserts one column in the database
func insertCol(db *pogreb.DB, key, column []byte) {

	err := db.Put([]byte(key), column)
	if err != nil {
		log.Fatal(err)
	}

}

// loadData load a CSV file into the database
func loadData(db *pogreb.DB, data [][]string, header []string) {
	ncol := len(header)
	nrow := len(data)
	column := make([]string, nrow)
	fmt.Print(" col : ", ncol, " rows : ", nrow, "	")
	bar := progressbar.Default(int64(ncol)) // Add a new progress bar
	for i := 0; i < ncol; i++ {
		bar.Add(1) // show progress bar
		for j := 0; j < nrow; j++ {
			column[j] = data[j][i]
		}
		//fmt.Println(column)
		str := []byte(strings.Join(column, "\t"))
		//insertCol(db, []byte(header[i]), str)
		insertCol(db, []byte(fmt.Sprint(i)), str) // the key is the colnumber
	}
}

// ReadColumn read one column from the database
func ReadColumn(db *pogreb.DB, key string) []string {

	val, err := db.Get([]byte(key))
	if err != nil {
		log.Fatal(err)
	}

	return strings.Split(string(val), "\t")

}

// ReadColumns read only columns with positions in indexes
func ReadColumns(filename string, colIndexes []int) [][]string {
	// get the user preference for using the database
	pref := fyne.CurrentApp().Preferences()
	useDBpref := binding.BindPreferenceBool("useDataBase", pref)
	useDB, err := useDBpref.Get()
	check(err)

	// if the user want to use a database, data are read from the database
	if useDB {
		return ReadColumnsDB(filename, colIndexes)
	}
	return ReadColumnsCSV(filename, colIndexes)
}

// ReadColumns read only columns with positions in indexes
func ReadColumnsCSV(filename string, colIndexes []int) [][]string {
	var xy [][]string
	// Open the file
	csvfile, err := os.Open("data/" + filename)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	defer csvfile.Close()
	// Parse the file
	r := csv.NewReader(bufio.NewReader(csvfile))
	//r := csv.NewReader(csvfile)
	r.Comma = '\t'
	r.Read() // skip header

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
		xy = append(xy, filter.SelByIndex(record, colIndexes))
	}
	return xy
}

// ReadColumns read from pogreb Database only columns with positions in indexes
func ReadColumnsDB(filename string, colIndexes []int) [][]string {
	var xy [][]string
	log.Println("Reading coordinates from database...")
	cxy := make([][]string, len(colIndexes)) // expression and XY columns
	// Open the database to read expression and XY columns
	dbname, _ := filter.RemExt(filename)
	db, err := pogreb.Open("temp/pogreb/"+dbname, nil)
	if err != nil {
		log.Println("database ", dbname, "not found in temp/pogreb/")
		log.Fatal(err)

	}
	defer db.Close()

	// read the 3 columns : expresion, x,y
	for i := 0; i < len(colIndexes); i++ {
		cxy[i] = ReadColumn(db, fmt.Sprint(colIndexes[i]))
		//log.Println("column ", i, " value = ", cxy[i][0:5])
	}

	// Iterate through the records
	for ct := 0; ct < len(cxy[1]); ct++ {
		tmp := make([]string, len(colIndexes))
		for i := 0; i < len(colIndexes); i++ {
			tmp[i] = cxy[i][ct]
		}
		xy = append(xy, tmp)
	}
	log.Println("data read !")
	return xy
}

func intToBytes(i int) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(i))
	return b
}

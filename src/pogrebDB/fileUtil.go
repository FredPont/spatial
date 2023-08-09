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
	"log"
	"os"
	"strings"

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

func intToBytes(i int) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(i))
	return b
}

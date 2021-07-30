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
 (c) Frederic Pont 2021
*/

package filter

import (
	"bufio"
	"encoding/csv"
	"io"

	//"fmt"

	"log"
	"os"
)

// ReadCompareTable read only columns with positions in indexes and filter rows by groups of gates
// xy coordinates are the last indexes in colIndexes
func ReadCompareTable(zoom int, filename string, colIndexes []int, XYindex []int, group1, group2 [][]Point, param Conf) ([][]string, [][]string, bool) {
	var table1, table2 [][]string
	// Open the file
	csvfile, err := os.Open("data/" + filename)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

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
		test, g := rowInGroup(zoom, record, XYindex, group1, group2, param)
		if test {
			if g == 1 {
				table1 = append(table1, selByIndex(record, colIndexes))
			} else if g == 2 {
				table2 = append(table2, selByIndex(record, colIndexes))
			}
		} else {
			if g == -1 {
				return [][]string{}, [][]string{}, false
			}
		}
	}
	return table1, table2, true
}

// filter row and put it in group1 or 2 for comparison depending on gates in group 1 and 2
func rowInGroup(zoom int, record []string, XYindex []int, group1, group2 [][]Point, param Conf) (bool, int) {

	for _, polygon1 := range group1 {
		if filterRow(zoom, record, XYindex, polygon1, param) {
			return true, 1
		}
	}
	for _, polygon2 := range group2 {
		if filterRow(zoom, record, XYindex, polygon2, param) {
			return true, 2
		}
	}
	return false, 0
}

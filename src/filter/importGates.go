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
	"log"
	"os"
	"strconv"
	"strings"
)

// ReadGate reads a gate in ImageJ format. The gate must be at scale 100%
// X,Y
// 131,150
// 105,189
// 156,187
func ReadGate(dir, filename string) []Point {

	var pts []Point
	// Open the file
	csvfile, err := os.Open(dir + "/" + filename)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(bufio.NewReader(csvfile))
	//r := csv.NewReader(csvfile)
	r.Comma = ','
	header, _ := r.Read() // skip header
	if !checkFormat(header, filename) {
		return nil
	}

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

		xpt, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Println("cannot import gate from "+filename, err)
			continue
		}
		ypt, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Println("cannot import gate from "+filename, err)
			continue
		}

		pts = append(pts, Point{int(xpt), int(ypt)})

	}
	//fmt.Println("gate readed :", pts)
	return pts
}

// checkFormat verify that the header of the gate file contains "X,Y"
func checkFormat(header []string, file string) bool {
	h := strings.Join(header, ",")
	if h == "X,Y" {
		return true
	}
	log.Println("The gate format of ", file, " is wrong, X,Y is header is missing ! cannot import the gate")
	return false
}

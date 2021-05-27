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
	"math"
	"os"
	"strconv"
	"strings"
)

// Conf store user preferences
type Conf struct {
	X, Y   string
	Scale  float64
	Rotate bool
}

// FilterTable filter the scRNAseq table to extract cells in polygon
func FilterTable(dataFile, outfile string, polygon []Point, param Conf) {

	path := "data/" + dataFile
	// open result file for write filtered table
	fout := "results/filtered_" + outfile
	out, err1 := os.Create(fout)
	check(err1)
	defer out.Close()

	// open result file for write cell names
	fout2 := "results/cells_" + outfile
	out2, err1 := os.Create(fout2)
	check(err1)
	defer out2.Close()

	csvFile, err := os.Open(path)
	check(err)
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1

	// read table header
	header, err := reader.Read()                       //read first line of pathway
	writeOneLine(out, strings.Join(header, "\t")+"\n") // write header in result file
	XYindex := GetColIndex(header, []string{param.X, param.Y})
	//fmt.Println(XYindex)
	for {
		// Read in a row. Check if we are at the end of the file.
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if filterRow(record, XYindex, polygon, param) {
			line := strings.Join(record, "\t") + "\n"
			//fmt.Println(line)
			writeOneLine(out, line)
			writeOneLine(out2, record[0]+"\n")
		}

	}
	return
}

func filterRow(record []string, XYindex []int, polygon []Point, param Conf) bool {
	scaleFactor := param.Scale
	rotate := param.Rotate

	XYstr := selByIndex(record, XYindex) // []string with XY coordinates
	//fmt.Println(XYstr)

	x, err := strconv.ParseFloat(XYstr[0], 64)
	if err != nil {
		log.Fatal(err)
	}
	xScaled := int64(math.Round(x * scaleFactor))
	y, err := strconv.ParseFloat(XYstr[1], 64)
	if err != nil {
		log.Fatal(err)
	}
	yScaled := int64(math.Round(y * scaleFactor))

	if rotate == true {
		xRot := yScaled
		yRot := xScaled
		return inGate(xRot, yRot, polygon)
	} else {
		return inGate(xScaled, yScaled, polygon)
	}

}

func inGate(x, y int64, polygon []Point) bool {
	// convert x,y to int
	return isInside(polygon, Point{int(x), int(y)})
}

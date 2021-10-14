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
func FilterTable(zoom int, dataFile, outfile string, polygon []Point, param Conf, gateNumber int) {

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
	header, err := reader.Read()                                         //read first line of pathway
	WriteOneLine(out, strings.Join(header, "\t")+"\t"+"GateNumber"+"\n") // write header in result file
	XYindex := GetColIndex(header, []string{param.X, param.Y})
	//fmt.Println(XYindex)
	for {
		// Read in a row. Check if we are at the end of the file.
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if filterRow(zoom, record, XYindex, polygon, param) {
			line := strings.Join(record, "\t") + "\t" + strconv.Itoa(gateNumber) + "\n"
			//fmt.Println(line)
			WriteOneLine(out, line)
			WriteOneLine(out2, record[0]+"\n")
		}

	}
	return
}

func filterRow(zoom int, record []string, XYindex []int, polygon []Point, param Conf) bool {
	scaleFactor := param.Scale
	rotate := param.Rotate

	XYstr := selByIndex(record, XYindex) // []string with XY coordinates
	//fmt.Println(XYstr)

	x, err := strconv.ParseFloat(XYstr[0], 64)
	if err != nil {
		log.Println(x, "cannot be converted to XY coordinate", err)
		return false
	}
	xScaled := int64(math.Round(x * scaleFactor))

	y, err := strconv.ParseFloat(XYstr[1], 64)
	if err != nil {
		log.Println(y, "cannot be converted to XY coordinate", err)
		return false
	}
	yScaled := int64(math.Round(y * scaleFactor))

	if rotate == true {
		xRot := yScaled
		yRot := xScaled
		return inGate(zoom, xRot, yRot, polygon)
	}
	return inGate(zoom, xScaled, yScaled, polygon)

}

func inGate(zoom int, x, y int64, polygon []Point) bool {
	// apply zoom to polygon
	if zoom != 100 {
		polygon = ZoomPolygon(polygon, 100./float64(zoom))
		//log.Println("zoomPolygon=", polygon)
	}
	// convert x,y to int
	return isInside(polygon, Point{int(x), int(y)})
}

// TablePlot filter the scRNAseq table XY coordinates to extract cells in polygon to draw a plot and return XY coordinates
func TablePlot(zoom int, tableXYxy [][]string, polygon []Point, param Conf, columnX, columnY string, ch1 chan<- [][]string) {
	// tableXYxy contains the index of the 2 columns to plot and the XY columns with the image coordinates
	// cf plot.makeplot() mapAndGates := filter.ReadColumns(filename, colIndexes)
	var xy [][]string
	scaleFactor := param.Scale
	rotate := param.Rotate
	//log.Println("start extract gates", polygon)
	for _, dot := range tableXYxy {
		if len(dot) < 4 {
			ch1 <- xy
			return
		}
		// x,y = coordinates of the dots in gate
		x, _ := strconv.ParseFloat(dot[2], 64)
		y, _ := strconv.ParseFloat(dot[3], 64)

		xScaled := int64(math.Round(x * scaleFactor))
		yScaled := int64(math.Round(y * scaleFactor))

		if rotate == true {
			xRot := yScaled
			yRot := xScaled
			if inGate(zoom, xRot, yRot, polygon) {
				xy = append(xy, []string{dot[0], dot[1]})
			}
		} else {
			if inGate(zoom, xScaled, yScaled, polygon) {
				xy = append(xy, []string{dot[0], dot[1]})
			}
		}
	}
	//log.Println("xy extract gates", xy)
	ch1 <- xy
	//return xy
}

// zoom polygon with zf without modifying stored polygon. function used to export gate a 100% zoom
func ZoomPolygon(p []Point, zf float64) []Point {
	var zoomedPoly []Point
	for i := 0; i < len(p); i++ {
		x := int(math.Round(float64(p[i].X) * zf))
		y := int(math.Round(float64(p[i].Y) * zf))
		zoomedPoly = append(zoomedPoly, Point{X: x, Y: y})
	}
	return zoomedPoly
}

// func zoomPolygon(zoom int, polygon []Point) []Point {
// 	var zoomPoly []Point
// 	for _, p := range polygon {
// 		zoomPoly = append(zoomPoly, Point{p.X * 100 / zoom, p.Y * 100 / zoom})
// 	}
// 	return zoomPoly
// }

/*
// applyZoomInt64 correct the input integer by the current zoom factor
func applyZoomInt64(zoom int, val int64) int64 {
	if zoom == 100 {
		return val
	}
	return val * int64(zoom) / 100
}
*/
/*
func TablePlot3(dataFile string, polygon []Point, param Conf, columnX, columnY string, ch1 chan<- [][]string) {
	var xy [][]string
	path := "data/" + dataFile

	csvFile, err := os.Open(path)
	check(err)
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1

	// read table header
	header, err := reader.Read() //read first line of pathway

	// get the index of the columns with the XY coordinates of the microscopie image
	XYindex := GetColIndex(header, []string{param.X, param.Y})
	// get the index of the columns with the XY coordinates of the gate dots
	gateidx := GetColIndex(header, []string{columnX, columnY})
	//fmt.Println(XYindex)
	for {
		// Read in a row. Check if we are at the end of the file.
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if filterRow(record, XYindex, polygon, param) {
			xy = append(xy, []string{record[gateidx[0]], record[gateidx[1]]})
		}

	}
	ch1 <- xy
}
*/

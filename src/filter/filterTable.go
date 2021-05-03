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
	"log"
)

// detect columns to select
// test :
// header := []string{"a", "b", "c", "d", "e", "f"}
// list := []string{"e", "b", "z", "c"}
// result : [4 1 2]
func getColIndex(header, list []string) []int {
	var indexes []int
	indDic := make(map[string]int) // dic of items -> column index
	list2 := make([]string, len(list))
	copy(list2, list)

	for i, val := range header {
		for j, l := range list2 {
			if val == l {
				indDic[val] = i
				list2 = append(list2[:j], list2[j+1:]...) // remove found item from list
				break
			}
		}
	}

	//indexes = append(indexes, 0) // append the first column containing cells names

	for _, v := range list {
		value, exist := indDic[v]
		if exist {
			indexes = append(indexes, value)
		}
	}

	if len(indexes) < 2 {
		log.Fatal("XY columns not found in table !")
	}
	return indexes
}

// selByIndex select item in a slice according to indexes
// we use it to select in the table X,Y columns
// corresponding to indexes positions
func selByIndex(row []string, indexes []int) []string {
	var selection []string

	for _, i := range indexes {
		selection = append(selection, row[i])
	}
	return selection
}

/*
func filterTable(dataFile string, polygon []Point, param Conf, gateFile string) {

	path := "data/" + dataFile
	gf, _ := remExt(gateFile)
	// open result file for write filtered table
	fout := "results/filtered_" + gf + "_" + dataFile
	out, err1 := os.Create(fout)
	check(err1)
	defer out.Close()

	// open result file for write cell names
	fout2 := "results/cells_" + gf + "_" + dataFile
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
	XYindex := getColIndex(header, []string{param.x, param.y})
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
	scaleFactor := param.scale
	rotate := param.rotate

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
	return isInside(polygon, Point{x, y})
}
*/

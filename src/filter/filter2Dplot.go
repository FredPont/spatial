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

	"fyne.io/fyne/v2/data/binding"
)

// DotsInGate takes a gate and extract the scatter point that are in that gate
func DotsInGate(gate []Point, scatter map[string]Point) map[string]Point {

	cells := make(map[string]Point) // cells names and XY in the gate

	for cell, point := range scatter {

		if isInside(gate, point) {
			cells[cell] = point
		}
	}
	return cells
}

// Filter2DGates filter the data tables by cell ID
// parallel code , one table/thread
func Filter2DGates(cellsInGates []map[string]Point, alledges [][]Point, gateName string, f binding.Float) {

	log.Print("start filtering...")
	dataFiles := ListFiles("data/") // list all tables in data dir
	// progress bar step
	nbFiles := len(cellsInGates) * len(dataFiles)
	step := 0.8 / float64(nbFiles)
	// filter all data files with all active gates
	ch := make(chan string)

	for gateNumber, cellsOnegate := range cellsInGates {
		for _, dataFile := range dataFiles {
			gateName = FormatOutFile("filter", gateName, "") // test if name exist, if not, build a file name with the current time
			outFile := strconv.Itoa(gateNumber) + "_" + gateName + "_" + dataFile
			go FilterTable2Dplot(cellsOnegate, dataFile, outFile, gateNumber, ch)
		}

	}
	// print the file progression and the progress bar
	for i := 0; i < nbFiles; i++ {
		msg := <-ch
		prog, _ := f.Get()
		f.Set(prog + step) // progress bar
		log.Println(msg, " done !")
	}
	f.Set(0.)
}

// FilterTable filter the scRNAseq table to extract cells id from the 2D plot
func FilterTable2Dplot(cellsInGates map[string]Point, dataFile, outfile string, gateNumber int, ch chan string) chan string {

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
	header, err := reader.Read() //read first line of pathway
	if err != nil {
		log.Println("cannot read table header : ", dataFile)
	}
	WriteOneLine(out, strings.Join(header, "\t")+"\t"+"GateNumber"+"\n") // write header in result file

	for {
		// Read in a row. Check if we are at the end of the file.
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if filterRowbyID(record, record[0], cellsInGates) {
			line := strings.Join(record, "\t") + "\t" + strconv.Itoa(gateNumber) + "\n"
			//fmt.Println(line)
			WriteOneLine(out, line)
			WriteOneLine(out2, record[0]+"\n")
		}

	}
	ch <- fout
	return ch

}

// test if the table row begins by "id"
func filterRowbyID(record []string, id string, cellsInGates map[string]Point) bool {
	_, inKey := cellsInGates[id]
	return inKey
}

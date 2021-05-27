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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//ReadHeader read header of table
func ReadHeader(path string) []string {

	csvFile, err := os.Open(path)
	check(err)
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1

	record, err := reader.Read() // read first line

	return record
}

// detect columns to select
// test :
// header := []string{"a", "b", "c", "d", "e", "f"}
// list := []string{"e", "b", "z", "c"}
// result : [4 1 2]
func GetColIndex(header, list []string) []int {
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

// ReadColumns read only columns with positions in indexes
func ReadColumns(filename string, colIndexes []int) [][]string {
	var xy [][]string
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
		xy = append(xy, selByIndex(record, colIndexes))
	}
	return xy
}

// IndexHeader create a map with column names => row number
// func indexHeader(header []string) map[string]int {
// 	index := make(map[string]int)

// 	for i := 0; i < len(header); i++ {
// 		index[header[i]] = i //
// 	}
// 	return index
// }

// remove file extension
func remExt(filename string) (string, string) {
	var extension = filepath.Ext(filename)
	var name = filename[0 : len(filename)-len(extension)]
	return name, extension
}

// ListFiles lists all files in a directory
func ListFiles(dir string) []string {
	var filesList []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		//fmt.Println(f.Name())
		filesList = append(filesList, f.Name())
	}
	return filesList
}

//###########################################
func writeOneLine(f *os.File, line string) {
	_, err := f.WriteString(line)
	check(err)
}

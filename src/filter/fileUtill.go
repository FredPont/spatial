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
	"fmt"
	"io"
	"math"
	"strconv"
	"time"

	//"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// ReadHeader read header of table
func ReadHeader(path string) []string {

	csvFile, err := os.Open(path)
	check(err)
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1

	record, err := reader.Read() // read first line
	if err != nil {
		log.Println("cannot read header of ", path)
	}

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
	//log.Println("list columns", list, header)

	for _, v := range list {
		value, exist := indDic[v]
		if exist {
			indexes = append(indexes, value)
		}
	}

	if len(indexes) < len(list) {
		log.Fatal("Columns not found in table ! ", list)
	}
	return indexes
}

// SelByIndex select item in a slice according to indexes
// we use it to select in the table X,Y columns
// corresponding to indexes positions
func SelByIndex(row []string, indexes []int) []string {
	var selection []string

	for _, i := range indexes {
		selection = append(selection, row[i])
	}
	return selection
}

// ReadClusters read only columns with positions in indexes and fill a map
// cluster NB => slice of x,y coordinates
func ReadClusters(a fyne.App, filename string, colIndexes []int) map[int][]Point {

	// get scaleFactor and rotation from pref
	pref := a.Preferences()

	sf := binding.BindPreferenceFloat("scaleFactor", pref) // set the link to preferences for scaling factor
	scaleFactor, _ := sf.Get()                             // read the preference for scaling factor

	rot := binding.BindPreferenceString("rotate", pref) // set the link to preferences for rotation
	rotate, _ := rot.Get()

	// map with cluster number => slice of xy coordinates scaled
	clusterMap := make(map[int][]Point, 0)
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
		cxy := SelByIndex(record, colIndexes)
		xScaled, yScaled := ScaleXY(cxy[1], cxy[2], scaleFactor, rotate)
		clustNB, err := strconv.Atoi(cxy[0])
		check(err)
		clusterMap[clustNB] = append(clusterMap[clustNB], Point{int(xScaled), int(yScaled)})
	}
	return clusterMap
}

// ScaleXY apply scaling factor and rotation to xy
func ScaleXY(X, Y string, scaleFactor float64, rotate string) (int64, int64) {

	x, err := strconv.ParseFloat(X, 64)
	if err != nil {
		log.Fatal(err)
	}
	xScaled := int64(math.Round(x * scaleFactor))

	y, err := strconv.ParseFloat(Y, 64)
	if err != nil {
		log.Fatal(err)
	}
	yScaled := int64(math.Round(y * scaleFactor))

	return Rotation(xScaled, yScaled, rotate)

}

// Rotation set the image rotation to + 90 or -90
func Rotation(xScaled, yScaled int64, rotate string) (int64, int64) {
	pref := fyne.CurrentApp().Preferences()

	switch rotate {
	case "no rotation":
		return xScaled, yScaled
	case "+90":
		xRot := yScaled
		yRot := xScaled
		return xRot, yRot
	case "-90":

		// get the image heigth
		imgHeight := binding.BindPreferenceInt("imgH", pref)
		imgH, err := imgHeight.Get()
		if err != nil {
			log.Println("cannot read image width from preferences !", err)
		}

		xRot := yScaled
		yRot := int64(imgH) - xScaled
		return xRot, yRot
	case "Vertical mirror":
		// get the image width
		imgWidth := binding.BindPreferenceInt("imgW", pref)
		imgW, err := imgWidth.Get()
		if err != nil {
			log.Println("cannot read image width from preferences !", err)
		}
		yRot := yScaled
		xRot := int64(imgW) - xScaled
		return xRot, yRot
	case "Horizontal mirror":
		// get the image heigth
		imgHeight := binding.BindPreferenceInt("imgH", pref)
		imgH, err := imgHeight.Get()
		if err != nil {
			log.Println("cannot read image width from preferences !", err)
		}

		xRot := xScaled
		yRot := int64(imgH) - yScaled
		return xRot, yRot
	}
	return xScaled, yScaled
}

// ReadImportedCells read cell names into []string
func ReadImportedCells(filename string) []string {
	var cellnames []string
	// Open the file
	csvfile, err := os.Open("import_cells/" + filename)
	if err != nil {
		log.Fatalln("Couldn't open the file", err)
	}
	defer csvfile.Close()
	// Parse the file
	r := csv.NewReader(bufio.NewReader(csvfile))
	//r := csv.NewReader(csvfile)
	//r.Comma = '\t'

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
		cellnames = append(cellnames, record[0])
	}
	return cellnames
}

// ClustersByCells read only columns with specific cell names and positions in indexes and fill a map
// cluster NB => slice of x,y coordinates
func ClustersByCells(a fyne.App, filename string, colIndexes []int, cellImport map[string]bool) map[int][]Point {

	// get scaleFactor and rotation from pref
	pref := a.Preferences()

	sf := binding.BindPreferenceFloat("scaleFactor", pref) // set the link to preferences for scaling factor
	scaleFactor, _ := sf.Get()                             // read the preference for scaling factor

	rot := binding.BindPreferenceString("rotate", pref) // set the link to preferences for rotation
	rotate, _ := rot.Get()

	// map with cluster number => slice of xy coordinates scaled
	clusterMap := make(map[int][]Point, 0)
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
	defer csvfile.Close()
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
		ncxy := SelByIndex(record, colIndexes)
		// continue if cell name is not in the cellnames map keys
		if !strInMap(ncxy[0], cellImport) {
			continue
		}
		xScaled, yScaled := ScaleXY(ncxy[2], ncxy[3], scaleFactor, rotate)
		clustNB, err := strconv.Atoi(ncxy[1])
		check(err)
		clusterMap[clustNB] = append(clusterMap[clustNB], Point{int(xScaled), int(yScaled)})
	}
	return clusterMap
}

// // Express contains expression value associated with x,y coordinates
// type Express struct {
// 	e float64
// 	p Point
// }

// ReadExpress read only columns with positions in indexes and fill a a map
// of expression normalized between 0-1 => slice of x,y coordinates
func ReadExpress(a fyne.App, filename string, colIndexes []int) ([]float64, []Point) {

	// get scaleFactor and rotation from pref
	pref := a.Preferences()

	sf := binding.BindPreferenceFloat("scaleFactor", pref) // set the link to preferences for scaling factor
	scaleFactor, _ := sf.Get()                             // read the preference for scaling factor

	rot := binding.BindPreferenceString("rotate", pref) // set the link to preferences for rotation
	rotate, _ := rot.Get()

	// array of expression valules and xy coordinates scaled
	var expressions []float64
	var pts []Point
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
		cxy := SelByIndex(record, colIndexes)
		xScaled, yScaled := ScaleXY(cxy[1], cxy[2], scaleFactor, rotate)
		exp, err := strconv.ParseFloat(cxy[0], 64)
		if err != nil {
			//log.Println("column number", colIndexes[0]+1, "does not contain a number", err)
			continue
		}

		pts = append(pts, Point{int(xScaled), int(yScaled)})
		expressions = append(expressions, exp)
	}
	return expressions, pts
}

// ReadGradient import CSV list of hexadecimal colors into []string
func ReadGradient(filename string) []string {
	var str []string
	// Open the file
	csvfile, err := os.Open(filename)
	if err != nil {
		log.Fatalln("Couldn't open the file", err)
	}
	defer csvfile.Close()
	// Parse the file
	r := csv.NewReader(bufio.NewReader(csvfile))
	//r := csv.NewReader(csvfile)
	//r.Comma = '\t'

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
		str = append(str, record[0])
	}
	return str
}

// RemExt remove file extension
func RemExt(filename string) (string, string) {
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

// FileExist check if a file exists
func FileExist(filename string) bool {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("File %s does not exist\n", filename)
		return false
	} else {
		fmt.Printf("File %s exists\n", filename)
		return true
	}
}

// CopyFile copy a file from src to dst
// credits https://github.com/mactsouk/opensource.com/blob/master/cp1.go
func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// ReadDir read dir directory and return the file names
func ReadDir(dir string) []string {
	// Read the contents of the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	fileNames := make([]string, len(files))
	// Print the name of each file and folder in the directory
	for i, file := range files {
		fileNames[i] = file.Name()
		//fmt.Println(file.Name())
	}
	return fileNames
}

// ClearDir removes all files in dir
func ClearDir(dir string) {
	// delete dir
	err := os.RemoveAll(dir)
	if err != nil {
		fmt.Println("Error clearing directory:", err)
		return
	}
	// restore the deleted dir
	err = os.Mkdir(dir, 0777)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	//fmt.Println("Directory cleared successfully.")
}

// RmFiles removes file by name
func RmFile(fileName string) {

	err := os.Remove(fileName)
	if err != nil {
		fmt.Println("Error deleting file:", err)
		return
	}
	//fmt.Println("File deleted successfully")
}

// MvFiles moves a file
func MvFile(oldPath, newPath string) {

	err := os.Rename(oldPath, newPath)
	if err != nil {
		fmt.Println("Error moving file:", err)
		return
	}
	//fmt.Println("File moved successfully")
}

// WriteOneLine write line to file
func WriteOneLine(f *os.File, line string) {
	_, err := f.WriteString(line)
	check(err)
}

// FLstr convert float to string
func FLstr(f float64) string {
	return strconv.FormatFloat(f, 'e', 3, 64)
}

// FormatOutFile add extension csv to file name or build a file name with time string when the filename is not given by the user
func FormatOutFile(prefix, name string, ext string) string {
	var outfile string

	if name == "" {
		currentTime := time.Now()
		outfile = prefix + "_" + currentTime.Format("2006-01-02_150405") + ext
	} else {
		outfile = name + ext
	}
	return outfile
}

// TrimString cut a string to length
// credits : https://dev.to/takakd/go-safe-truncate-string-9h0
func TrimString(str string, length int) string {
	if length <= 0 {
		return ""
	}
	truncated := ""
	count := 0
	for _, char := range str {
		truncated += string(char)
		count++
		if count >= length {
			break
		}
	}
	return truncated
}

// WriteCSV export a [][]string as CSV file
func WriteCSV(data [][]string, path string) {
	file, err := os.Create(path)
	check(err)

	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = '\t' // tab separator
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write(value)
		check(err)
	}
}

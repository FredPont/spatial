package ui

import (
	"lasso/src/filter"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// PVrecord is the gene/pathway name and the corresponding FC and Pvalue, Pvalue corrected
type PVrecord struct {
	item                            string
	fc, pv, pvcorr, log2fc, log10pv float64
}

func buttonCompare(a fyne.App, e *Editor, preference fyne.Preferences, f binding.Float, header []string, headerMap map[string]interface{}, firstTable string) {

	compWindow := a.NewWindow("Compare")

	gates := buildGateNames(e)

	// g1Map : stores the gates selected in group 1
	g1Map := make(map[string]interface{}, len(gates))
	buildMapFalse(gates, g1Map)
	// g2Map : stores the gates selected in group 1
	g2Map := make(map[string]interface{}, len(gates))
	buildMapFalse(gates, g2Map)

	content := container.New(layout.
		NewGridLayoutWithColumns(2), container.NewVBox(
		widget.NewLabel("Groups to compare"),
		widget.NewLabel("Group 1"),
		listGates(gates, g1Map),
		widget.NewLabel("Group 2"),
		listGates(gates, g2Map),
		widget.NewLabel("Columns to compare"),
		widget.NewButton("Select all", func() {
			buildMapTrue(header[1:], headerMap)
			//boxes = binding.BindUntypedMap(&headerMap)
			//boxes.Reload()
			compWindow.Content().Refresh()
		}),
		widget.NewButton("unSelect all", func() {
			buildMapFalse(header[1:], headerMap)
			//boxes = binding.BindUntypedMap(&headerMap)
			//boxes.Reload()
			compWindow.Content().Refresh()
		}),
		widget.NewButton("Compare", func() {
			//log.Println(g1Map, g2Map)
			//log.Println(headerMap)
			if !chkGates(g1Map, g2Map) {
				return
			}
			go startComparison(e, header, headerMap, preference, firstTable, f, g1Map, g2Map)
		}),
		widget.NewButton("Close", func() {
			compWindow.Close()
		}),
	),
		listColums(header[1:], headerMap),
	)
	//content := listColums(header, headerMap, boxes)
	compWindow.SetContent(content)
	compWindow.Resize(fyne.Size{Width: 500, Height: 500})
	compWindow.Show()
}

// list with all columns in header except column 1 with cell names
func listColums(header []string, headerMap map[string]interface{}) fyne.CanvasObject {
	strings := binding.BindStringList(&header)
	l := widget.NewListWithData(strings,
		func() fyne.CanvasObject {
			return widget.NewCheck("text", func(val bool) {
				log.Println(val)
			})

		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			lbl := item.(binding.String)
			label, _ := lbl.Get()

			chk := obj.(*widget.Check)
			chk.Text = label
			chk.Checked = headerMap[chk.Text].(bool)
			chk.OnChanged = func(done bool) {
				//log.Println(chk.Text)
				log.Println(chk.Text, chk.Checked)
				headerMap[chk.Text] = chk.Checked
				//cb, _ := boxes.Get()
				//cb[chk.Text] = chk.Checked
				//boxes.Reload()
				//log.Println(boxes.GetValue("one"))
			}
			chk.Refresh()

		})

	return l
}

// build a map with header/slice item set to true
func buildMapTrue(s []string, sMap map[string]interface{}) {
	for i := 0; i < len(s); i++ {
		sMap[s[i]] = true
	}
}

// build a map with header/slice item set to false
func buildMapFalse(s []string, sMap map[string]interface{}) {
	for i := 0; i < len(s); i++ {
		sMap[s[i]] = false
	}
}

// list with all gates
func listGates(gates []string, gMap map[string]interface{}) fyne.CanvasObject {
	strings := binding.BindStringList(&gates)
	l := widget.NewListWithData(strings,
		func() fyne.CanvasObject {
			return widget.NewCheck("text", func(val bool) {
				log.Println(val)
			})

		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			lbl := item.(binding.String)
			label, _ := lbl.Get()

			chk := obj.(*widget.Check)
			chk.Text = label
			chk.Checked = gMap[chk.Text].(bool)
			chk.OnChanged = func(done bool) {
				//log.Println(chk.Text)
				log.Println(chk.Text, chk.Checked)
				gMap[chk.Text] = chk.Checked
			}
			chk.Refresh()

		})

	return l
}

// verify that the same gate is not selected in the group 1 & 2. Return true if the groups are different
func chkGates(g1Map map[string]interface{}, g2Map map[string]interface{}) bool {
	chk := true
	for k, v := range g1Map {
		if v == true && g2Map[k] == v {
			log.Println(k, "is selected in the two groups ! Deselect", k, "in one group !")
			chk = false
		}
	}
	return chk
}

// create an []string with gates names gate_1, gate_2...
func buildGateNames(e *Editor) []string {
	var gates []string
	alledges := e.drawSurface.alledges
	for i := 0; i < len(alledges); i++ {
		gates = append(gates, "gate_"+strconv.Itoa(i))
	}
	return gates
}

// start comparison
func startComparison(e *Editor, header []string, headerMap map[string]interface{}, pref fyne.Preferences, firstTable string, f binding.Float, g1Map, g2Map map[string]interface{}) {
	f.Set(.3)
	param := prefToConf(pref) // get xy rotation zoom factor from pref
	// group1, group2  = polygone coordinates of gates in group 1 and 2
	group1, group2 := gatesInGroup(e, g1Map, g2Map)
	// index and colnames of selected colum + XY coordinates for filtering
	index, colnames, XYindex := selColIndex(header, headerMap, pref) // colnames : cellsId, colum selected, X,Y
	log.Println(index)
	// to save RAM each line of the table is filtered when reading the table
	table1, table2, test := filter.ReadCompareTable(e.zoom, firstTable, index, XYindex, group1, group2, param)
	if test == false {
		log.Println("Error detected in XY coordinates ! Comparison aborted !")
		return
	}

	f.Set(.5)
	pvfcTable := foldChangePV(table1, table2, colnames)

	// log.Println(colnames, group1, group2)
	// log.Println(table1)
	// log.Println(table2)
	// log.Println(pvfcTable)
	writePV("comparison.tsv", pvfcTable)
	f.Set(0.)
}

// get column indexes of selected column from header map (checkboxes = true)
// The cells names are the first column folowed by the columns selected by the user.
// The XY columns are not included in the table but the XY indexes are caclulated for gate filtration
func selColIndex(header []string, headerMap map[string]interface{}, pref fyne.Preferences) ([]int, []string, []int) {
	// XY indexes in header
	var xidx, yidx int

	// X coordinates
	xcor := binding.BindPreferenceString("xcor", pref) // set the link to preferences for rotation
	xc, _ := xcor.Get()

	// y coordinates
	ycor := binding.BindPreferenceString("ycor", pref) // set the link to preferences for rotation
	yc, _ := ycor.Get()

	index := []int{0} // index initialized with cell names = col0
	colnames := []string{header[0]}
	for i, v := range header {
		if i == 0 {
			continue
		}
		if headerMap[v] == true {
			index = append(index, i)
			colnames = append(colnames, header[i])
		}
		// search xy columns
		if v == xc {
			// index = append(index, i)
			// colnames = append(colnames, header[i])
			xidx = i
		}
		if v == yc {
			// index = append(index, i)
			// colnames = append(colnames, header[i])
			yidx = i
		}
	}
	return index, colnames, []int{xidx, yidx}
}

func gatesInGroup(e *Editor, g1Map, g2Map map[string]interface{}) ([][]filter.Point, [][]filter.Point) {
	var group1, group2 [][]filter.Point
	alledges := e.drawSurface.alledges
	for i, gate := range alledges {
		if g1Map["gate_"+strconv.Itoa(i)] == true {
			group1 = append(group1, gate)
		} else if g2Map["gate_"+strconv.Itoa(i)] == true {
			group2 = append(group2, gate)
		}
	}
	return group1, group2
}

func foldChangePV(table1, table2 [][]string, colnames []string) []PVrecord {
	var pvTable []PVrecord
	nc := len(table1[0]) // col number
	//var fc, pv float64   //foldchange pvalue
	for c := 1; c < nc; c++ {
		v1 := getColum(c, table1)
		v2 := getColum(c, table2)
		fc, t := folchange(v1, v2)
		// if undetermined fc == 0/0  the data is skiped
		if t == false {
			continue
		}
		pv, _ := filter.PvMannWhitney(v1, v2)
		// PV corrected by Bonferroni
		pvBonf := pvBonferroni(pv, float64(nc-1))

		pvTable = append(pvTable, PVrecord{colnames[c], fc, pv, pvBonf, math.Log2(fc), log10pv(pvBonf)})
	}
	return pvTable
}

// extract column c from table and convert it to float
func getColum(c int, table [][]string) []float64 {
	var col []float64
	l := len(table)
	for i := 0; i < l; i++ {
		x, err := strconv.ParseFloat(table[i][c], 64)
		if err != nil {
			log.Println("cell=", table[i][0], "col=", c+1, " ", x, "cannot be converted to float", err)
			return []float64{}
		}
		col = append(col, x)
	}
	return col
}

func folchange(x1, x2 []float64) (float64, bool) {
	s1 := sumFloat(x1)
	s2 := sumFloat(x2)
	if s1 == 0 && s2 == 0 {
		log.Println("fold-change undetermined (0/0) !")
		return 1., false
	} else if s1 != 0 && s2 != 0 {
		return s2 / s1, true
	} else if s2 == 0 {
		return 1e-10, true
	}
	log.Println("division by zero in fold-change caculation !")
	return 1e10, true
}

// -log10(pv)
func log10pv(pv float64) float64 {
	if pv == 0 {
		return 300
	}
	return -math.Log10(pv)

}

// PV corrected by Bonferroni
func pvBonferroni(pv, n float64) float64 {
	pvb := pv * n
	if pvb > 1. {
		return 1.
	}
	return pvb
}

func sumFloat(array []float64) float64 {
	result := 0.
	for _, v := range array {
		result += v
	}
	return result
}

func writePV(filename string, pvTable []PVrecord) {
	path := "comparison/" + filename
	// open result file for write filtered table
	out, err1 := os.Create(path)
	if err1 != nil {
		log.Println(path, "cannot be written ! The file is not saved !")
		return
	}
	defer out.Close()

	// write header to file
	header := []string{"item", "FoldChange", "PV_Wilcoxon", "PV_Bonferroni", "log2(FC)", "-log10(PV)\n"}
	filter.WriteOneLine(out, strings.Join(header, "\t"))

	for _, rec := range pvTable {
		line := []string{rec.item, filter.FLstr(rec.fc), filter.FLstr(rec.pv), filter.FLstr(rec.pvcorr),
			filter.FLstr(rec.log2fc), filter.FLstr(rec.log10pv) + "\n"}
		filter.WriteOneLine(out, strings.Join(line, "\t"))
	}
}

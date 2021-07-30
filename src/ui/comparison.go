package ui

import (
	"lasso/src/filter"
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

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
			log.Println(g1Map, g2Map)
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
	index, colnames, XYindex := selColIndex(header, headerMap, pref) // colnames : cellsId, colum selected, X,Y
	log.Println(index)
	// to save RAM each line of the table is filtered when reading the table
	table1, table2, test := filter.ReadCompareTable(e.zoom, firstTable, index, XYindex, group1, group2, param)
	if test == false {
		log.Println("Error detected in XY coordinates ! Comparison aborted !")
		return
	}
	_ = table1
	_ = table2
	f.Set(.5)

	log.Println(colnames, group1, group2)
	log.Println(table1)
	log.Println(table2)
	f.Set(0.)
}

// get column indexes of selected column from header map (checkboxes = true) . The cells names and XY coordinated are required for gate filtration and included from preferences
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
		// add xy columns
		if v == xc {
			index = append(index, i)
			colnames = append(colnames, header[i])
			xidx = i
		}
		if v == yc {
			index = append(index, i)
			colnames = append(colnames, header[i])
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

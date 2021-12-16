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

package ui

import (
	"log"
	"spatial/src/filter"
	"strconv"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func buildVulcanoTools(v *Vulcano) {
	data := [][]string{{"item", "log2FC", "log10pv"}}
	// progress bar binding
	f := binding.NewFloat()
	table := widget.NewTable(
		func() (int, int) {
			return len(data), len(data[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i.Row][i.Col])
		})

	content := container.New(layout.
		NewGridLayoutWithColumns(2), container.NewVBox(
		widget.NewLabel("Left click and right click on the vulcano plot to select points"),
		widget.NewButton("Save vulcano plot", func() {
			imgName := filter.FormatOutFile("vulcano", "", "")
			go screenShot(v.win, imgName, f)
		}),
		widget.NewButton("Close", func() {
			v.tools.Close()
			v.win.Close()
		}),
	),

		table,
	)
	w := fyne.CurrentApp().NewWindow("Vulcano Tools")

	w.SetContent(content)
	v.tools = w
	w.Show()
}

//refreshVulanoTools reload the table of selected vulcano dots
func refreshVulanoTools(v *Vulcano) {
	a := fyne.CurrentApp()
	pref := a.Preferences()
	// progress bar binding
	f := binding.NewFloat()
	selItem := binding.BindString(&v.drawSurface.selItem)
	data := [][]string{{"item", "X (log2FC)", "Y (log10pv)"}}
	selRecord := v.drawSurface.selection
	data = append(data, PVtoString(selRecord)...)

	colwidth := dataColwidth(data)

	table := widget.NewTable(
		func() (int, int) {
			return len(data), len(data[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content : gene exression, pathways etc..")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i.Row][i.Col])
		})

	table.SetColumnWidth(0, colwidth) // table 1st colunm width adjustment
	table.OnSelected = func(id widget.TableCellID) {
		if len(data) < 1 {
			return
		}
		// if user does not select the first colunm , column 0 is used anyway
		item := data[id.Row][0]
		//log.Println("selec=", item)
		v.drawSurface.selItem = item
	}

	// show choice of different gradien
	gradExpression := binding.BindPreferenceString("gradExpression", pref) // pref binding for the expression gradien to avoid reset for each vulcano dot
	selGrad, _ := gradExpression.Get()
	grad := widget.NewRadioGroup([]string{"Turbo", "Viridis", "Inferno", "White - Red", "Yellow - Red", "Purple - Red", "Blue - Yellow - Red"}, func(s string) {
	})
	if selGrad != "" {
		grad.Selected = selGrad
	}
	grad.OnChanged = func(s string) { pref.SetString("gradExpression", s) }

	// Dot opacity
	DotOp := binding.BindPreferenceFloat("dotOpacity", pref) // pref binding for the expression dot opacity
	DotOpacity := widget.NewSliderWithData(0., 255., DotOp)
	//DotOpacity.Step = 1.
	//DotOpacity.Value = 255.
	// DotOpacity.OnChanged = func(v float64) {
	// 	pref.SetFloat("dotOpacity", v)
	// }

	//legend color - the results is store in preferences
	legendcol := widget.NewButton("Legend Text Color", func() { LegendTxtscolor(a, v.tools) })

	content := container.New(layout.
		NewGridLayoutWithColumns(2), container.NewVBox(
		widget.NewLabel("Left click and right click on the vulcano plot to select points"),
		widget.NewLabel("Select your gradient"),
		grad,
		widget.NewLabel("Dots Opacity [0-100%] :"),
		DotOpacity,
		legendcol,
		widget.NewButton("Show Expression", func() {
			//f := binding.NewFloat() // progress bar binding

			// gradien default
			def := "White - Red"
			if grad.Selected == "" {
				grad.Selected = def
			}

			// PathwayIndex is needed for compatiblity with the function to draw expression but useless here because slideshow is not enable
			PathwayIndex := binding.NewInt() // column index of current pathway displayed by slide show
			PathwayIndex.Set(1)              // start with column 1 by default

			choosedItem, _ := selItem.Get() // item selected by the user on the table
			if choosedItem == "" {
				log.Println("You must select one row !")
				return
			}

			go drawExp(a, v.imageEditor, v.header, v.tableName, choosedItem, grad.Selected, f, PathwayIndex, v.tools)
		}),
		widget.NewButton("Save vulcano plot", func() {
			imgName := filter.FormatOutFile("vulcano", "", "")
			go screenShot(v.win, imgName, f)
		}),
		widget.NewButton("Close", func() {
			v.tools.Close()
			v.win.Close()
		}),
	),

		table,
	)

	v.tools.SetContent(content)
}

// PVtoString convert []PVrecord to [][]string with only item, log2fc , log10pv:
func PVtoString(pv []PVrecord) [][]string {
	var data [][]string
	for _, r := range pv {
		data = append(data, []string{r.item, strconv.FormatFloat(r.log2fc, 'f', 3, 64), strconv.FormatFloat(r.log10pv, 'f', 3, 64)})
	}
	return data
}

// dataColwidth mesures the data item max length assuming 10 pixels/ char
func dataColwidth(data [][]string) float32 {
	lmax := 5
	if len(data) == 0 {
		return float32(10 * lmax)
	}
	for _, d := range data {
		l := utf8.RuneCountInString(d[0])
		if l > lmax {
			lmax = l
		}

	}
	return float32(10 * lmax) // 10 pixels per char
}

package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func buttonCompare(a fyne.App, e *Editor, preference fyne.Preferences, f binding.Float, header []string, boxes binding.ExternalUntypedMap, headerMap map[string]interface{}) {

	compWindow := a.NewWindow("Compare")

	content := container.New(layout.
		NewGridLayoutWithColumns(2),
		widget.NewLabel("Columns to compare"),
		widget.NewButton("Select all", func() {
			buildMapTrue(header[1:], headerMap)
			boxes := binding.BindUntypedMap(&headerMap)
			boxes.Reload()
			compWindow.Content().Refresh()

		}),
		widget.NewButton("unSelect all", func() {
			buildMapFalse(header[1:], headerMap)
			boxes := binding.BindUntypedMap(&headerMap)
			boxes.Reload()
			compWindow.Content().Refresh()
		}),
		listColums(header[1:], headerMap, boxes),
	)
	//content := listColums(header, headerMap, boxes)
	compWindow.SetContent(content)
	compWindow.Resize(fyne.Size{Width: 500, Height: 500})
	compWindow.Show()
}

func listColums(header []string, headerMap map[string]interface{}, boxes binding.ExternalUntypedMap) fyne.CanvasObject {
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
				cb, _ := boxes.Get()
				cb[chk.Text] = chk.Checked
				boxes.Reload()
				log.Println(boxes.GetValue("one"))
			}
			chk.Refresh()

		})

	return l
}

// build a map with header item set to true
func buildMapTrue(s []string, sMap map[string]interface{}) {

	for i := 0; i < len(s); i++ {
		sMap[s[i]] = true
	}

}

// build a map with header item set to false
func buildMapFalse(s []string, sMap map[string]interface{}) {

	for i := 0; i < len(s); i++ {
		sMap[s[i]] = false
	}

}

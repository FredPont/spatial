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

package pref

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// updateMyButton change the label of a button
func updateMyButton(b *widget.Button, label string) {
	b.Text = label
	//b.Icon = someResource
	b.Refresh()
}

// ShowTable display a table in a form and return the item selected in sel binding.String
func ShowTable(data []string, sel binding.String, but *widget.Button, winLabel string) {
	a := fyne.CurrentApp()
	w2 := a.NewWindow(winLabel)

	form := widget.NewForm(
		widget.NewFormItem("Select one item in table", widget.NewLabel(":")),
	)
	form.OnCancel = func() {
		//fmt.Println("Cancelled")
		w2.Close()
	}
	form.OnSubmit = func() {
		userSel, _ := sel.Get()
		updateMyButton(but, userSel)
		//fmt.Println("Form submitted :", userSel)
		w2.Close()

	}

	table := widget.NewList(
		func() int { return len(data) },
		func() fyne.CanvasObject {
			//icon := widget.NewIcon(theme.FileIcon())
			label := widget.NewLabel("wide content : gene exression, pathways etc..")
			return container.NewHBox(label)
		},
		func(index int, template fyne.CanvasObject) {
			cont := template.(*fyne.Container)
			label := cont.Objects[0].(*widget.Label)
			label.SetText(fmt.Sprintf(data[index]))
		})

	table.OnSelected = func(id int) {
		if len(data) < 1 {
			return
		}
		item := data[id]
		//log.Println("selec=", item)
		sel.Set(item)
	}

	content := container.New(layout.NewGridLayoutWithColumns(2),
		form,
		table,
	)

	w2.SetContent(content)
	w2.Resize(fyne.Size{100, 500})
	w2.Show()

}

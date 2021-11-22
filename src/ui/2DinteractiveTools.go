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
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// show2DinterTools show 2Di tools and 2Di window
func show2D(a fyne.App, e *Editor, preference fyne.Preferences, f binding.Float, header []string, firstTable string) {
	winplot := build2DplotWin(e) // show 2D interactive window
	show2DinterTools(a, e, winplot, preference, f, header, firstTable)
	//buttonCompare(a, e, preference, f, header, headerMap, firstTable)
}

func show2DinterTools(a fyne.App, e *Editor, winplot fyne.Window, preference fyne.Preferences, f binding.Float, header []string, firstTable string) {

	win2Dtools := a.NewWindow("2D plot tools")

	content := container.NewVBox(
		widget.NewLabel("Tools"),
		widget.NewButton("Show Cells in Gates", func() {

		}),
		// screenshot
		widget.NewButtonWithIcon("", theme.MediaPhotoIcon(), func() {
			//go screenShot(w, gatename.Text, f)
		}),
		widget.NewButtonWithIcon("Exit", theme.LogoutIcon(), func() {
			win2Dtools.Close() // close tool window
			winplot.Close()    // close plot window
		}),
	)
	win2Dtools.SetContent(content)
	//win2D.Resize(fyne.Size{Width: 500, Height: 500})
	win2Dtools.Show()
}

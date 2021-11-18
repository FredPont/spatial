package ui

import (
	"fyne.io/fyne/v2"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func brushesButtons(edit *Editor, a fyne.App) fyne.CanvasObject {

	var svgbrush = "<svg xmlns=\"http://www.w3.org/2000/svg\" height=\"48px\" viewBox=\"0 0 24 24\" width=\"48px\" fill=\"#de8159\"><path d=\"M0 0h24v24H0z\" fill=\"none\"/><path d=\"M7 14c-1.66 0-3 1.34-3 3 0 1.31-1.16 2-2 2 .92 1.22 2.49 2 4 2 2.21 0 4-1.79 4-4 0-1.66-1.34-3-3-3zm13.71-9.37l-1.34-1.34c-.39-.39-1.02-.39-1.41 0L9 12.25 11.75 15l8.96-8.96c.39-.39.39-1.02 0-1.41z\"/></svg>"
	var svglasso = "<svg xmlns=\"http://www.w3.org/2000/svg\" enable-background=\"new 0 0 20 20\" height=\"48px\" viewBox=\"0 0 20 20\" width=\"48px\" fill=\"#de8159\"><g><rect fill=\"none\" height=\"20\" width=\"20\" x=\"0\"/></g><g><path d=\"M16.27,10l-3.14,5.5H6.87L3.73,10l3.14-5.5h6.26L16.27,10z M6,3l-4,7l4,7h8l4-7l-4-7H6z\"/></g></svg>"
	iconbrush := fyne.NewStaticResource("brush", []byte(svgbrush))
	iconlasso := fyne.NewStaticResource("lasso", []byte(svglasso))

	var b1 *widget.Button
	var b2 *widget.Button

	// pencil button
	b1 = widget.NewButtonWithIcon("", iconbrush, func() {
		//preference.SetString("brush", "pencil")
		edit.brush = "pencil"
		activePencil(b1)
		inactiveLasso(b2, svglasso)
	})

	// lasso button
	b2 = widget.NewButtonWithIcon("", iconlasso, func() {
		//preference.SetString("brush", "lasso")
		edit.brush = "lasso"
		activeLasso(b2)
		inactivePencil(b1, svgbrush)
	})

	// by default, brush = pencil
	if edit.brush == "" {
		edit.brush = "pencil"
		activePencil(b1)
		inactiveLasso(b2, svglasso)
	}

	drawtools := container.NewHBox(
		b1,
		b2,
	)
	return drawtools
}

func activePencil(b *widget.Button) {
	var svgbrush2 = "<svg xmlns=\"http://www.w3.org/2000/svg\" height=\"48px\" viewBox=\"0 0 24 24\" width=\"48px\" fill=\"#59db71\"><path d=\"M0 0h24v24H0z\" fill=\"none\"/><path d=\"M7 14c-1.66 0-3 1.34-3 3 0 1.31-1.16 2-2 2 .92 1.22 2.49 2 4 2 2.21 0 4-1.79 4-4 0-1.66-1.34-3-3-3zm13.71-9.37l-1.34-1.34c-.39-.39-1.02-.39-1.41 0L9 12.25 11.75 15l8.96-8.96c.39-.39.39-1.02 0-1.41z\"/></svg>"
	//b.Text = ""
	b.Icon = fyne.NewStaticResource("brush2", []byte(svgbrush2))
	b.Refresh()
}

func activeLasso(b *widget.Button) {
	var svglasso2 = "<svg xmlns=\"http://www.w3.org/2000/svg\" enable-background=\"new 0 0 20 20\" height=\"48px\" viewBox=\"0 0 20 20\" width=\"48px\" fill=\"#59db71\"><g><rect fill=\"none\" height=\"20\" width=\"20\" x=\"0\"/></g><g><path d=\"M16.27,10l-3.14,5.5H6.87L3.73,10l3.14-5.5h6.26L16.27,10z M6,3l-4,7l4,7h8l4-7l-4-7H6z\"/></g></svg>"
	//b.Text = ""
	b.Icon = fyne.NewStaticResource("lasso2", []byte(svglasso2))
	b.Refresh()
}

func inactivePencil(b *widget.Button, svgbrush string) {
	//b.Text = ""
	b.Icon = fyne.NewStaticResource("brush", []byte(svgbrush))
	b.Refresh()
}

func inactiveLasso(b *widget.Button, svglasso string) {
	//b.Text = ""
	b.Icon = fyne.NewStaticResource("lasso", []byte(svglasso))
	b.Refresh()
}

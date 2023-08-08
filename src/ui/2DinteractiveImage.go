package ui

import (
	"fmt"
	"image/color"
	"log"
	"spatial/src/filter"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/fogleman/gg"
)

// startDraw2DplotImg draws the 2D scatter plot on the image "temp/2Dplot/2Dplot.png" with single or multithread
func (p *PlotBox) startDraw2DplotImg(v *Interactive2Dsurf, dotsize int) {
	pref := fyne.CurrentApp().Preferences()
	Nthreads := pref.IntWithFallback("nbExpressThreads", 1)
	log.Println("start 2D scatter plot")
	if Nthreads > 1 {
		go p.MTdraw2DplotImg(v, dotsize)
	} else {
		go p.draw2DplotImg(v, dotsize)
	}

}

// draw2DplotImg draws the 2D scatter plot on the image "temp/2Dplot/2Dplot.png"
func (p *PlotBox) draw2DplotImg(v *Interactive2Dsurf, dotsize int) {
	W, H := 800, 800
	// clear selected dots image
	init2DselectedDots()

	dc := gg.NewContext(W, H)
	dc.SetRGBA(0, 0, 0, 0) // create a transparent image
	dc.Clear()

	for i, xplot := range p.X {
		//x := int(MapRange(x, p.Xmin, p.Xmax, p.Left, 800-p.Right))
		//y := int(MapRange(p.Y[i], p.Ymin, p.Ymax, p.Bottom, 800-p.Top))
		x, y := xCoord(p, xplot), yCoord(p, p.Y[i])
		//v.drawcircleScattCont(x, y, dotsize, color.NRGBA{128, 128, 128, 255})
		//log.Println(x, y)
		dc.SetRGBA(128.0/255.0, 128.0/255.0, 128.0/255.0, 1.0)
		dc.DrawPoint(float64(x), float64(y), float64(dotsize))
		dc.Fill()
	}
	dc.SavePNG("temp/2Dplot/2Dplot.png")
	log.Println("2D plot done !")
	v.layer.Refresh()
}

// MTdraw2DplotImg draws the 2D scatter plot on the image "temp/2Dplot/2Dplot.png" MultiThread version
func (p *PlotBox) MTdraw2DplotImg(v *Interactive2Dsurf, dotsize int) {
	W, H := 800, 800
	// clear selected dots image
	init2DselectedDots()

	dc := gg.NewContext(W, H)
	dc.SetRGBA(0, 0, 0, 0) // create a transparent image
	dc.Clear()

	pref := fyne.CurrentApp().Preferences()
	// divide the number of points into Nthreads parts
	// number of threads
	Nthreads := pref.IntWithFallback("nbExpressThreads", 1)
	chunks := filter.DivideNB(len(p.X), Nthreads)
	var chunkStart, chunkEnd int
	var wg sync.WaitGroup
	for index := 0; index < len(chunks); index++ {
		wg.Add(1)
		if index == 0 {
			chunkStart = 0
		} else {
			chunkStart = filter.SumSliceInt(index, chunks)
		}
		chunkEnd = filter.SumSliceInt(index+1, chunks) - 1 // -1 because the first index is zero
		//log.Println("chunkStart  ", chunkStart, " chunkEnd  ", chunkEnd, "  ", chunks)
		go draw2DplotOnImage(p, chunkStart, chunkEnd, dotsize, index, &wg)
	}
	wg.Wait()
	MergeIMG("temp/2Dplot/MTdots/", "temp/2Dplot/2Dplot.png") // faster than MTmergeIMG
	log.Println("2D plot done !")
	v.layer.Refresh()
}

// draw2DplotOnImage draw 1/nThtreads of 2D plot on a separate image
func draw2DplotOnImage(p *PlotBox, chunkStart, chunkEnd, dotsize, index int, wg *sync.WaitGroup) {
	W, H := 800, 800
	dc := gg.NewContext(W, H)
	dc.SetRGBA(0, 0, 0, 0) // create a transparent image
	dc.Clear()

	for c := chunkStart; c <= chunkEnd; c++ {

		x, y := xCoord(p, p.X[c]), yCoord(p, p.Y[c])
		dc.SetRGBA(128.0/255.0, 128.0/255.0, 128.0/255.0, 1.0)
		dc.DrawPoint(float64(x), float64(y), float64(dotsize))
		dc.Fill()
	}

	dc.SavePNG("temp/2Dplot/MTdots/" + fmt.Sprint(index) + "out.png")

	wg.Done()

}

// draw the selected cells on the microscopy image
func drawCellsImg(e *Editor, cellsXY []filter.Point, dotcolor color.NRGBA, gateNB int) {
	pref := fyne.CurrentApp().Preferences()

	//calc image size to draw the cluster dots
	H := pref.IntWithFallback("imgH", 500)
	W := pref.IntWithFallback("imgW", 500)
	// clear calc image
	dc := gg.NewContext(W, H)
	dc.SetRGBA(0, 0, 0, 0) // create a transparent image
	dc.Clear()

	clustOp := binding.BindPreferenceFloat("clustOpacity", pref) // cluster opacity
	opacity, _ := clustOp.Get()
	op := uint8(opacity)
	clustDia := binding.BindPreferenceInt("clustDotDiam", pref) // cluster dot diameter
	diameter, _ := clustDia.Get()
	//diameter = ApplyZoomInt(e, diameter)
	sf := binding.BindPreferenceFloat("scaleFactor", pref) // set the link to preferences for scaling factor
	scaleFactor, _ := sf.Get()                             // read the preference for scaling factor
	rot := binding.BindPreferenceString("rotate", pref)    // set the link to preferences for rotation
	rotate, _ := rot.Get()

	for _, xy := range cellsXY {
		xScaled, yScaled := scale(xy.X, xy.Y, scaleFactor, rotate)
		//e.drawcircle(ApplyZoomInt(e, xScaled), ApplyZoomInt(e, yScaled), diameter, color.NRGBA{dotcolor.R, dotcolor.G, dotcolor.B, op})
		//log.Println(xy)
		dc.SetRGBA(float64(dotcolor.R)/255.0, float64(dotcolor.G)/255.0, float64(dotcolor.B)/255.0, float64(op)/255.0)
		dc.DrawPoint(float64(xScaled), float64(yScaled), float64(diameter))
		dc.Fill()
	}
	dc.SavePNG("temp/2Dplot/selectedDots/" + fmt.Sprint(gateNB) + "_gate.png")
}

// plot the dots in gates in color in the 2D scatter plot. dots are plotted in the calcSeldots image container
func plotDotsInGatesImg(p *PlotBox, inter2D *Interactive2Dsurf, cellsInGates []map[string]filter.Point, selectedGradient string) {
	prefs := fyne.CurrentApp().Preferences()
	//get scatter dot size
	ds := binding.BindPreferenceString("2Ddotsize", prefs) // set the link to 2D dot size preferences
	ds2 := binding.StringToInt(ds)
	dotsize, _ := ds2.Get()

	W, H := 800, 800 // 2D plot image size
	// clear calc image
	dc := gg.NewContext(W, H)
	dc.SetRGBA(0, 0, 0, 0) // create a transparent image
	dc.Clear()

	nbGates := len(cellsInGates)
	for i := 0; i < nbGates; i++ {
		dotcolor := dotColors(nbGates, i, selectedGradient)
		//p.gatesDotPlot(inter2D, dotsize, cellsInGates[i], dotcolor)
		dc.SetRGBA(float64(dotcolor.R)/255.0, float64(dotcolor.G)/255.0, float64(dotcolor.B)/255.0, 1.0)
		for _, xy := range cellsInGates[i] {
			dc.DrawPoint(float64(xy.X), float64(xy.Y), float64(dotsize))
			dc.Fill()
		}

	}
	dc.SavePNG("temp/2Dplot/dotsIngGates.png")
	inter2D.layer.Refresh()
}

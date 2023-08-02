package ui

import (
	"fmt"
	"image/color"
	"log"
	"spatial/src/filter"
	"spatial/src/plot"
	"sync"

	"github.com/fogleman/gg"
	"github.com/mazznoer/colorgrad"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

// startClusterComput start cluster computation with single or multithread
func startClusterComput(a fyne.App, e *Editor, header []string, firstTable string, f binding.Float) {
	initTempDir("temp/clusters")
	pref := a.Preferences()
	MT := pref.BoolWithFallback("multithreadCluster", false)
	if MT {
		go MTDrawImageClusters(a, e, header, firstTable, f)
	} else {
		go DrawImageClusters(a, e, header, firstTable, f)
	}

}

// startEspressComput start expression computation with single or multithread
func startEspressComput(a fyne.App, e *Editor, header []string, firstTable string, userSel, gradien string, f binding.Float, curPathwayIndex binding.Int, ExpressWindow fyne.Window) {
	initTempDir("temp/expression")
	pref := a.Preferences()
	Nthreads := pref.IntWithFallback("nbExpressThreads", 1)
	if Nthreads == 1 {
		go drawImageExp(a, e, header, firstTable, userSel, gradien, f, curPathwayIndex, ExpressWindow)

	} else {

		go MTdrawImageExp(a, e, header, firstTable, userSel, gradien, f, curPathwayIndex, ExpressWindow)
	}

}

// MTDrawImageClusters use the library gg (https://pkg.go.dev/github.com/fogleman/gg) to create a transparent png image
// with the size of the original microscopy image. Clusters are ploted as dots in this image
// MT multithread version
func MTDrawImageClusters(a fyne.App, e *Editor, header []string, filename string, f binding.Float) {
	// clear the cluster temp directory
	filter.ClearDir("temp/clusters")
	log.Println("start clusters computation...")
	f.Set(0.1)     // progress bar
	initCluster(e) // remove all dots of the cluster container
	pref := a.Preferences()

	//calc image size to draw the cluster dots
	H := pref.IntWithFallback("imgH", 500)
	W := pref.IntWithFallback("imgW", 500)

	clustOp := binding.BindPreferenceFloat("clustOpacity", pref) // cluster opacity
	opacity, _ := clustOp.Get()
	op := uint8(opacity)
	clustDia := binding.BindPreferenceInt("clustDotDiam", pref) // cluster dot diameter
	diameter, _ := clustDia.Get()
	legendDiameter := SetLegendDiameter(e, diameter)

	clusterMap := getClusters(a, header, filename) // cluster nb => []Point
	log.Println(len(clusterMap), "clusters detected")
	f.Set(0.3) // progress bar

	nbCluster := len(clusterMap)
	clustNames := filter.KeysIntPoint(clusterMap)

	legendPosition := filter.Point{X: 15, Y: 15} // initial legend position for cluster names

	colors := allClustColors(nbCluster)
	R, G, B, _ := plot.GetPrefColorRGBA(a, "legendColR", "legendColG", "legendColB", "legendColA")
	colorText := color.NRGBA{uint8(R), uint8(G), uint8(B), 255}

	// if the hide legend preference is checked, the legend is not drawn
	hideL := binding.BindPreferenceBool("hideLegend", pref)
	hideLgd, _ := hideL.Get()
	var wg sync.WaitGroup
	for c := 0; c < nbCluster; c++ {
		wg.Add(1)

		go drawOnImage(c, W, H, diameter, clusterMap, colors, clustNames, op, &wg)

		// coordinates := clusterMap[clustNames[c]]
		clcolor := colors[c]

		//draw legend dot and name for the current cluster
		// if the hide legend preference is checked, the legend is not drawn
		if !hideLgd {
			drawLegend(e, clcolor.R, clcolor.G, clcolor.B, op, legendPosition.X, legendPosition.Y, legendDiameter, clustNames[c], colorText)
			legendPosition.Y = legendPosition.Y + 30
		}

		// // set progress bar to 50% when half cluster have been computed
		if c == int(nbCluster/2) {
			f.Set(0.5) // progress bar
		}
	}
	wg.Wait()
	// if the hide legend preference is checked, the legend name is not drawn
	if !hideLgd {
		titleLegend(e, "     clusters", getLegendColor(a))
	}

	f.Set(0.75) // progress bar
	MergeIMG("temp/clusters/", "temp/imgOut.png")
	//MTmergeIMG("temp/clusters/", "temp/imgOut.png")
	log.Println("clusters computation done ! Images in temp/clusters/")

	e.layer.Refresh()
	f.Set(0.) // reset progress bar
}

// drawOnImage draw each cluster on a separate image
func drawOnImage(c, W, H, diameter int, clusterMap map[int][]filter.Point, colors []RGB, clustNames []int, op uint8, wg *sync.WaitGroup) {

	dc := gg.NewContext(W, H)
	dc.SetRGBA(0, 0, 0, 0) // create a transparent image
	dc.Clear()
	coordinates := clusterMap[clustNames[c]]
	clcolor := colors[c]
	for i := 0; i < len(coordinates); i++ {
		dc.SetRGBA(float64(clcolor.R)/255.0, float64(clcolor.G)/255.0, float64(clcolor.B)/255.0, float64(op)/255.0)
		dc.DrawPoint(float64(coordinates[i].X), float64(coordinates[i].Y), float64(diameter))
		dc.Fill()
	}

	dc.SavePNG("temp/clusters/" + fmt.Sprint(c) + "out.png")

	wg.Done()

}

// DrawImageClusters use the library gg (https://pkg.go.dev/github.com/fogleman/gg) to create a transparent png image
// with the size of the original microscopy image. Clusters are ploted as dots in this image
func DrawImageClusters(a fyne.App, e *Editor, header []string, filename string, f binding.Float) {
	f.Set(0.1) // progress bar
	log.Println("start clusters computation...")
	initCluster(e) // remove all dots of the cluster container
	pref := a.Preferences()

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
	legendDiameter := SetLegendDiameter(e, diameter)

	clusterMap := getClusters(a, header, filename) // cluster nb => []Point
	log.Println(len(clusterMap), "clusters detected")
	f.Set(0.3) // progress bar

	nbCluster := len(clusterMap)
	clustNames := filter.KeysIntPoint(clusterMap)

	legendPosition := filter.Point{X: 15, Y: 15} // initial legend position for cluster names

	colors := allClustColors(nbCluster)
	R, G, B, _ := plot.GetPrefColorRGBA(a, "legendColR", "legendColG", "legendColB", "legendColA")
	colorText := color.NRGBA{uint8(R), uint8(G), uint8(B), 255}

	// if the hide legend preference is checked, the legend is not drawn
	hideL := binding.BindPreferenceBool("hideLegend", pref)
	hideLgd, _ := hideL.Get()

	for c := 0; c < nbCluster; c++ {
		// f.Set(float64(c) / float64(nbCluster-1)) // % progression for progress bar. This is too fast to be seen
		coordinates := clusterMap[clustNames[c]]
		clcolor := colors[c]
		for i := 0; i < len(coordinates); i++ {
			dc.SetRGBA(float64(clcolor.R)/255.0, float64(clcolor.G)/255.0, float64(clcolor.B)/255.0, float64(op)/255.0)
			dc.DrawPoint(float64(coordinates[i].X), float64(coordinates[i].Y), float64(diameter))
			dc.Fill()

		}
		// draw legend dot and name for the current cluster
		// if the hide legend preference is checked, the legend is not drawn
		if hideLgd {
			continue
		}
		drawLegend(e, clcolor.R, clcolor.G, clcolor.B, op, legendPosition.X, legendPosition.Y, legendDiameter, clustNames[c], colorText)
		legendPosition.Y = legendPosition.Y + 30
		// set progress bar to 50% when half cluster have been computed
		if c == int(nbCluster/2) {
			f.Set(0.5) // progress bar
		}
	}
	// if the hide legend preference is checked, the legend name is not drawn
	if !hideLgd {
		titleLegend(e, "     clusters", getLegendColor(a))
	}

	dc.SavePNG("temp/imgOut.png")
	log.Println("clusters computation done !")
	e.layer.Refresh()
	f.Set(0.) // reset progress bar
}

// drawImageExp use the library gg (https://pkg.go.dev/github.com/fogleman/gg) to create a transparent png image
// with the size of the original microscopy image. Expression spots are ploted as dots in this image
func drawImageExp(a fyne.App, e *Editor, header []string, filename string, expcol, gradien string, f binding.Float, curPathwayIndex binding.Int, ExpressWindow fyne.Window) {
	f.Set(0.1)     // progress bar set to 20%
	initCluster(e) // remove all dots of the cluster container
	pref := a.Preferences()

	//calc image size to draw the cluster dots
	H := pref.IntWithFallback("imgH", 500)
	W := pref.IntWithFallback("imgW", 500)
	// clear calc image
	dc := gg.NewContext(W, H)
	dc.SetRGBA(0, 0, 0, 0) // create a transparent image
	dc.Clear()

	// Dot opacity
	DotOp := binding.BindPreferenceFloat("dotOpacity", pref) // pref binding for the  dot opacity
	opacity, _ := DotOp.Get()
	op := uint8(opacity)

	// pre-build expression gradient
	grad := preBuildGradient(gradien)

	// Dot opacity gradient
	gradop := binding.BindPreferenceBool("gradOpacity", pref)
	opGrad, _ := gradop.Get()                               // enabled/disabled opacity gradient
	gradMax := binding.BindPreferenceFloat("gradMax", pref) // pref binding for user Max value for the opacity gradient
	opMax, _ := gradMax.Get()
	gradMin := binding.BindPreferenceFloat("gradMin", pref) // pref binding for user Min value for the opacity gradient
	opMin, _ := gradMin.Get()
	if opMin >= opMax {
		log.Println("Min threshold, ", opMin, " must be <= Max threshold, ", opMax)
	}

	clustDia := binding.BindPreferenceInt("clustDotDiam", pref) //  dot diameter
	diameter, _ := clustDia.Get()
	legendDiameter := SetLegendDiameter(e, diameter)

	log.Println("start reading data")
	expressions, pts := getExpress(a, header, filename, expcol, curPathwayIndex) // []expressions and []Point
	if len(expressions) < 1 {
		log.Println("Intensities not available for column", expcol)
		return
	}
	log.Println("stop reading data")
	f.Set(0.2) // progress bar set to 30% after data reading
	nbPts := len(pts)
	scaleExp, min, max := filter.ScaleSlice01(expressions)

	// density plot of the expression distribution
	go plot.BuildDensity(expressions, 100., filter.TrimString(expcol, 40), ExpressWindow)
	go saveTMPfiles(pts, expressions, min, max, nbPts)

	for c := 0; c < nbPts; c++ {
		// Calculate the progress as a percentage
		progress := (c + 1) * 100 / nbPts

		// Check if the progress is a multiple of 25%
		if progress%25 == 0 {
			f.Set(float64(progress) / 100.) // 50 % progression for progress bar
		}
		// transparency gradient
		if opGrad && opMin < opMax {
			op = gradTransp(expressions[c], min, max, opMin, opMax)
		}
		//op = gradTransp(expressions[c], min, max)
		clcolor := gradUser(gradien, grad, scaleExp[c])

		dc.SetRGBA(float64(clcolor.R)/255.0, float64(clcolor.G)/255.0, float64(clcolor.B)/255.0, float64(op)/255.0)
		dc.DrawPoint(float64(pts[c].X), float64(pts[c].Y), float64(diameter))
		dc.Fill()
	}
	// draw legend title, dot and value for the current cexpression
	// if the hide legend preference is checked, the legend is not drawn
	hideL := binding.BindPreferenceBool("hideLegend", pref)
	hideLgd, _ := hideL.Get()
	if !hideLgd {
		R, G, B, _ := plot.GetPrefColorRGBA(a, "legendColR", "legendColG", "legendColB", "legendColA")
		colorText := color.NRGBA{uint8(R), uint8(G), uint8(B), 255}
		titleLegend(e, expcol, colorText)
		// transparency gradient
		if opGrad && opMin < opMax {
			//fmt.Println("opGrad", opGrad, " opMin opMax", opMin, opMax)
			gradOpLegend(e, legendDiameter, gradien, grad, min, max, opMin, opMax, colorText)
		} else {
			expLegend(e, op, legendDiameter, gradien, grad, min, max, colorText) // not transparency gradien in legend
		}
	}
	dc.SavePNG("temp/imgOut.png")
	log.Println("expression computation done !")
	e.layer.Refresh()
	f.Set(0.) // reset progress bar
}

// drawImageExp use the library gg (https://pkg.go.dev/github.com/fogleman/gg) to create a transparent png image
// with the size of the original microscopy image. Expression spots are ploted as dots in this image
func refreshImageExp(a fyne.App, e *Editor, newMin, newMax float64, tmp filter.Record, scaleExp []float64, expcol, gradien string, f binding.Float, ExpressWindow fyne.Window) {
	f.Set(0.1)     // progress bar set to 20%
	initCluster(e) // remove all dots of the cluster container
	pref := a.Preferences()

	//calc image size to draw the cluster dots
	H := pref.IntWithFallback("imgH", 500)
	W := pref.IntWithFallback("imgW", 500)
	// clear calc image
	dc := gg.NewContext(W, H)
	dc.SetRGBA(0, 0, 0, 0) // create a transparent image
	dc.Clear()

	// Dot opacity
	DotOp := binding.BindPreferenceFloat("dotOpacity", pref) // pref binding for the  dot opacity
	opacity, _ := DotOp.Get()
	op := uint8(opacity)
	// pre-build expression gradient
	grad := preBuildGradient(gradien)
	// Dot opacity gradient
	gradop := binding.BindPreferenceBool("gradOpacity", pref)
	opGrad, _ := gradop.Get()                               // enabled/disabled opacity gradient
	gradMax := binding.BindPreferenceFloat("gradMax", pref) // pref binding for user Max value for the opacity gradient
	opMax, _ := gradMax.Get()
	gradMin := binding.BindPreferenceFloat("gradMin", pref) // pref binding for user Min value for the opacity gradient
	opMin, _ := gradMin.Get()
	if opMin >= opMax {
		log.Println("Min threshold, ", opMin, " must be <= Max threshold, ", opMax)
	}

	clustDia := binding.BindPreferenceInt("clustDotDiam", pref) //  dot diameter
	diameter, _ := clustDia.Get()
	legendDiameter := SetLegendDiameter(e, diameter)

	if len(tmp.Exp) < 1 {
		log.Println("Intensities not availble for column", expcol)
		return
	}
	f.Set(0.2) // progress bar set to 30% after data reading

	for c := 0; c < tmp.NbPts; c++ {
		// Calculate the progress as a percentage
		progress := (c + 1) * 100 / tmp.NbPts

		// Check if the progress is a multiple of 25%
		if progress%25 == 0 {
			f.Set(float64(progress) / 100.) // 50 % progression for progress bar
		}
		// transparency gradient
		if opGrad && opMin < opMax {
			op = gradTransp(tmp.Exp[c], newMin, newMax, opMin, opMax)
		}
		clcolor := gradUser(gradien, grad, scaleExp[c])

		dc.SetRGBA(float64(clcolor.R)/255.0, float64(clcolor.G)/255.0, float64(clcolor.B)/255.0, float64(op)/255.0)
		dc.DrawPoint(float64(tmp.Pts[c].X), float64(tmp.Pts[c].Y), float64(diameter))
		dc.Fill()
	}
	// draw legend title, dot and value for the current cexpression
	R, G, B, _ := plot.GetPrefColorRGBA(a, "legendColR", "legendColG", "legendColB", "legendColA")
	colorText := color.NRGBA{uint8(R), uint8(G), uint8(B), 255}
	titleLegend(e, expcol, colorText)
	// transparency gradient
	if opGrad && opMin < opMax {
		//fmt.Println("opGrad", opGrad, " opMin opMax", opMin, opMax)
		gradOpLegend(e, legendDiameter, gradien, grad, newMin, newMax, opMin, opMax, colorText)
	} else {
		expLegend(e, op, legendDiameter, gradien, grad, newMin, newMax, colorText) // not transparency gradien in legend
	}

	dc.SavePNG("temp/imgOut.png")
	log.Println("expression computation done !")
	e.layer.Refresh()
	f.Set(0.) // reset progress bar
}

// MTdrawImageExp use the library gg (https://pkg.go.dev/github.com/fogleman/gg) to create a transparent png image
// with the size of the original microscopy image. Expression spots are ploted as dots in this image
// MT multithread version
func MTdrawImageExp(a fyne.App, e *Editor, header []string, filename string, expcol, gradien string, f binding.Float, curPathwayIndex binding.Int, ExpressWindow fyne.Window) {
	f.Set(0.1)     // progress bar set to 20%
	initCluster(e) // remove all dots of the cluster container
	pref := a.Preferences()

	// number of threads
	Nthreads := pref.IntWithFallback("nbExpressThreads", 1)

	//calc image size to draw the cluster dots
	H := pref.IntWithFallback("imgH", 500)
	W := pref.IntWithFallback("imgW", 500)
	// clear calc image
	dc := gg.NewContext(W, H)
	dc.SetRGBA(0, 0, 0, 0) // create a transparent image
	dc.Clear()

	// Dot opacity
	DotOp := binding.BindPreferenceFloat("dotOpacity", pref) // pref binding for the  dot opacity
	opacity, _ := DotOp.Get()
	op := uint8(opacity)

	// pre-build expression gradient
	grad := preBuildGradient(gradien)

	// Dot opacity gradient
	gradop := binding.BindPreferenceBool("gradOpacity", pref)
	opGrad, _ := gradop.Get()                               // enabled/disabled opacity gradient
	gradMax := binding.BindPreferenceFloat("gradMax", pref) // pref binding for user Max value for the opacity gradient
	opMax, _ := gradMax.Get()
	gradMin := binding.BindPreferenceFloat("gradMin", pref) // pref binding for user Min value for the opacity gradient
	opMin, _ := gradMin.Get()
	if opMin >= opMax {
		log.Println("Min threshold, ", opMin, " must be <= Max threshold, ", opMax)
	}

	clustDia := binding.BindPreferenceInt("clustDotDiam", pref) //  dot diameter
	diameter, _ := clustDia.Get()
	legendDiameter := SetLegendDiameter(e, diameter)

	log.Println("start reading data")
	expressions, pts := getExpress(a, header, filename, expcol, curPathwayIndex) // []expressions and []Point
	if len(expressions) < 1 {
		log.Println("Intensities not available for column", expcol)
		return
	}
	log.Println("stop reading data")
	f.Set(0.2) // progress bar set to 30% after data reading
	nbPts := len(pts)
	scaleExp, min, max := filter.ScaleSlice01(expressions)

	// density plot of the expression distribution
	go plot.BuildDensity(expressions, 100., filter.TrimString(expcol, 40), ExpressWindow)
	go saveTMPfiles(pts, expressions, min, max, nbPts)

	// divide the number of points into Nthreads parts
	chunks := filter.DivideNB(nbPts, Nthreads)
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
		go drawExpOnImage(chunkStart, chunkEnd, pts, expressions, scaleExp, min, max, W, H, diameter, index, op, opGrad, opMin, opMax, gradien, grad, &wg)
	}
	wg.Wait()

	// draw legend title, dot and value for the current cexpression
	// if the hide legend preference is checked, the legend is not drawn
	hideL := binding.BindPreferenceBool("hideLegend", pref)
	hideLgd, _ := hideL.Get()
	if !hideLgd {
		R, G, B, _ := plot.GetPrefColorRGBA(a, "legendColR", "legendColG", "legendColB", "legendColA")
		colorText := color.NRGBA{uint8(R), uint8(G), uint8(B), 255}
		titleLegend(e, expcol, colorText)
		// transparency gradient
		if opGrad && opMin < opMax {
			//fmt.Println("opGrad", opGrad, " opMin opMax", opMin, opMax)
			gradOpLegend(e, legendDiameter, gradien, grad, min, max, opMin, opMax, colorText)
		} else {
			expLegend(e, op, legendDiameter, gradien, grad, min, max, colorText) // not transparency gradien in legend
		}
	}
	f.Set(0.75) // progress bar
	//MTmergeIMG("temp/expression/", "temp/imgOut.png")
	MergeIMG("temp/expression/", "temp/imgOut.png") // faster than MTmergeIMG
	log.Println("expression computation done !")
	e.layer.Refresh()
	f.Set(0.) // reset progress bar
}

// drawExpOnImage draw 1/nThtreads of expression on a separate image
func drawExpOnImage(chunkStart, chunkEnd int, pts []filter.Point, expressions, scaleExp []float64, min, max float64, W, H, diameter, index int, op uint8, opGrad bool, opMin, opMax float64, gradien string, grad colorgrad.Gradient, wg *sync.WaitGroup) {

	dc := gg.NewContext(W, H)
	dc.SetRGBA(0, 0, 0, 0) // create a transparent image
	dc.Clear()

	for c := chunkStart; c <= chunkEnd; c++ {

		if opGrad && opMin < opMax {
			op = gradTransp(expressions[c], min, max, opMin, opMax)
		}
		//op = gradTransp(expressions[c], min, max)
		clcolor := gradUser(gradien, grad, scaleExp[c])

		dc.SetRGBA(float64(clcolor.R)/255.0, float64(clcolor.G)/255.0, float64(clcolor.B)/255.0, float64(op)/255.0)
		dc.DrawPoint(float64(pts[c].X), float64(pts[c].Y), float64(diameter))
		dc.Fill()
	}

	dc.SavePNG("temp/expression/" + fmt.Sprint(index) + "out.png")

	wg.Done()

}

// SetLegendDiameter set the diameter for the legend spots
func SetLegendDiameter(e *Editor, diameter int) int {
	legendDiameter := ApplyZoomInt(e, diameter)
	if legendDiameter >= 15 {
		legendDiameter = 15
	} else if legendDiameter < 5 {
		legendDiameter = 5
	}
	return legendDiameter
}

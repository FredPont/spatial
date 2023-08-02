package ui

import (
	"io"
	"log"
	"os"
	"spatial/src/filter"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func initAlledges(e *Editor) {
	e.drawSurface.alledges = nil
	e.drawSurface.points = nil
}

// clear only the last gates edges and point and lines in container
func initLastedges(e *Editor) {
	e.drawSurface.alledges = filter.PopPoints(e.drawSurface.alledges)
	e.drawSurface.points = nil
	e.drawSurface.gatesLines = nil
}

func initCluster(e *Editor) {
	e.clusterContainer.Objects = nil
	// imgs := filter.ReadDir("image/")
	// filter.CopyFile("image/"+imgs[0], "temp/imgOut.png")
	// e.layer.Refresh()
}

func clearCluster(e *Editor) {
	e.clusterContainer.Objects = nil
	imgs := filter.ReadDir("image/")
	filter.CopyFile("image/"+imgs[0], "temp/imgOut.png")
	e.layer.Refresh()
	//e.clusterContainer.Refresh()
}

func initGates(e *Editor) {
	e.gateContainer.Objects = nil
	e.drawSurface.gatesLines = nil
	e.gateContainer.Refresh()
	initAlledges(e)      // reset alledges
	initGatesNB(e)       // clear all gates numbers in arrays
	initGatesNBwindow(e) // clear all gates numbers displayed
}

func initAllLayers(e *Editor) {
	initCluster(e)
	initGates(e)
}

func initGatesContainer(e *Editor) {
	e.gateContainer.Objects = nil
	e.gateContainer.Refresh()

}

// clear all stored gates numbers from arrays
func initGatesNB(e *Editor) {
	e.drawSurface.gatesNumbers.x = nil
	e.drawSurface.gatesNumbers.y = nil
	e.drawSurface.gatesNumbers.nb = 0
}

// clear all displayed gates numbers in the picture
func initGatesNBwindow(e *Editor) {
	e.gateNumberContainer.Objects = nil
	e.gateNumberContainer.Refresh()
}

// clear last gate number coordinates and decrease gateNB
func initLastGatesNB(e *Editor) {
	e.drawSurface.gatesNumbers.x = filter.PopIntArray(e.drawSurface.gatesNumbers.x)
	e.drawSurface.gatesNumbers.y = filter.PopIntArray(e.drawSurface.gatesNumbers.y)
	if e.drawSurface.gatesNumbers.nb > 0 {
		e.drawSurface.gatesNumbers.nb--
	}

}

// clear the dots at the edges of the currently drawing gate
func initGateDots(e *Editor) {
	e.gateDotsContainer.Objects = nil
	e.gateDotsContainer.Refresh()
}

// remove opacity gradient
func initOpacityGdt() {
	prefs := fyne.CurrentApp().Preferences()
	// Dot opacity gradient
	prefs.SetBool("gradOpacity", false)
}

// clear the tempdir for clusters or expression
func initTempDir(tempdir string) {
	filter.ClearDir(tempdir)
}

////////////////////////////
// Expression map INIT
////////////////////////////

// reset expression min max slider
func initSliderExp(MaxExp, MinExp *widget.Slider) {
	preference := fyne.CurrentApp().Preferences()
	preference.SetFloat("userMaxExp", 100)
	preference.SetFloat("userMinExp", 0)
	MaxExp.Value = 100. // reset slider position
	MinExp.Value = 0.   // reset slider position

}

// reset the density plot picture
func initDensityPlot() {
	src := "src/ui/sky.png"
	dst := "temp/density.png"

	fin, err := os.Open(src)
	if err != nil {
		log.Println(err)
	}
	defer fin.Close()

	fout, err := os.Create(dst)
	if err != nil {
		log.Println(err)
	}
	defer fout.Close()

	_, err = io.Copy(fout, fin)

	if err != nil {
		log.Println(err)
	}

}

////////////////////////////
// interactive 2D plot INIT
////////////////////////////

// clear gates from screen and initialise inter2D.drawSurface.alledges
func init2DScatterGates(inter2D *Interactive2Dsurf) {
	inter2D.drawSurface.alledges = nil
	inter2D.gateContainer.Objects = nil
	inter2D.drawSurface.gatesNumbers.x = nil
	inter2D.drawSurface.gatesNumbers.y = nil
	inter2D.drawSurface.gatesNumbers.nb = 0
	inter2D.gateContainer.Refresh()
}

////////////////////////////
//     preferences
////////////////////////////

// InitPref initialise some user preferences when not set
func InitPref() {

	prefs := fyne.CurrentApp().Preferences()

	// cluster dot diameter
	cld := prefs.IntWithFallback("clustDotDiam", 12)
	prefs.SetInt("clustDotDiam", cld)

	// cluster column
	clustercolumn := prefs.StringWithFallback("clustcol", "num_cluster")
	prefs.SetString("clustcol", clustercolumn)

	// chart dot size
	dotsize := prefs.StringWithFallback("dotsize", "3")
	prefs.SetString("dotsize", dotsize)

	// 2D interactive plot dot size
	ds2D := prefs.StringWithFallback("2Ddotsize", "3")
	prefs.SetString("2Ddotsize", ds2D)

	// get scaleFactor and rotation from pref
	sf := prefs.FloatWithFallback("scaleFactor", 0.107869044)
	prefs.SetFloat("scaleFactor", sf)

	// set rotate pref to +90
	rot := prefs.StringWithFallback("rotate", "+90")
	prefs.SetString("rotate", rot)

	// set antirotate pref to false
	//antirot := prefs.BoolWithFallback("rot-90", false)
	//prefs.SetBool("rot-90", antirot)

	// X coordinates
	xcor := prefs.StringWithFallback("xcor", "x_image")
	prefs.SetString("xcor", xcor)

	// y coordinates
	ycor := prefs.StringWithFallback("ycor", "y_image")
	prefs.SetString("ycor", ycor)

	//microscop windows W
	winW := prefs.FloatWithFallback("winW", 500)
	prefs.SetFloat("winW", winW)

	//microscop windows Height
	winH := prefs.FloatWithFallback("winH", 500)
	prefs.SetFloat("winH", winH)

	// vulcano selection square size in pixels
	// record the vulcano selection square size in pixels in preferences
	vs := prefs.IntWithFallback("vulcSelectSize", 20)
	prefs.SetInt("vulcSelectSize", vs)

	// Dot opacity
	prefs.SetFloat("dotOpacity", 255)

	// Dot opacity gradient
	prefs.SetBool("gradOpacity", false)

	// cluster opacity
	prefs.SetFloat("clustOpacity", 255)

	// vulcano default gradien
	gradExpression := prefs.StringWithFallback("gradExpression", "Turbo")
	prefs.SetString("gradExpression", gradExpression)

	// plot background color
	initBCKGColors([]string{"unselR", "unselG", "unselB", "unselA"})
	// plot foreground colors
	initFORGColors([]string{"gateDotsR", "gateDotsG", "gateDotsB", "gateDotsA"})
	// legend text colors
	initFORGColors([]string{"legendColR", "legendColG", "legendColB", "legendColA"})

	// multithread cluster computation
	clustThreads := prefs.BoolWithFallback("multithreadCluster", false)
	prefs.SetBool("multithreadCluster", clustThreads)
	// number of threads for expression
	expThreads := prefs.IntWithFallback("nbExpressThreads", 1)
	prefs.SetInt("nbExpressThreads", expThreads)
}

// init background colors
func initBCKGColors(rgba []string) {
	pref := fyne.CurrentApp().Preferences()
	if !sumRGBA(rgba) {
		//log.Println("background is not white")
		return
	}
	for _, c := range rgba {
		val := pref.IntWithFallback(c, 170)
		pref.SetInt(c, val) // grey
	}
	pref.SetInt(rgba[3], 255) // full opacity
}

// init foreground colors
func initFORGColors(rgba []string) {
	pref := fyne.CurrentApp().Preferences()
	if !sumRGBA(rgba) {
		return
	}
	for _, c := range rgba {
		val := pref.IntWithFallback(c, 1)
		pref.SetInt(c, val) // black
	}
	pref.SetInt(rgba[3], 255) // full opacity
}

func sumRGBA(rgba []string) bool {
	pref := fyne.CurrentApp().Preferences()
	sumRGBA := 0
	for _, c := range rgba[:3] {
		dotsRGBA := binding.BindPreferenceInt(c, pref)
		mapRGBA, _ := dotsRGBA.Get()
		sumRGBA += mapRGBA
	}
	// if the color is white sumRGBA == 3*255
	//log.Println("sumRGBA = ", sumRGBA)
	if sumRGBA == 3*255 || sumRGBA == 0 {
		return true
	}

	return false

}

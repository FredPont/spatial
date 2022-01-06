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
}

func clearCluster(e *Editor) {
	e.clusterContainer.Objects = nil
	e.clusterContainer.Refresh()
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
	cld := binding.BindPreferenceInt("clustDotDiam", prefs) // set the link to preferences for cluster diameter
	clud, _ := cld.Get()
	if clud == 0 {
		prefs.SetInt("clustDotDiam", 12)
	}

	// cluster column
	clustercolumn := binding.BindPreferenceString("clustcol", prefs) // set the link to preferences for rotation
	clucol, _ := clustercolumn.Get()
	if len(clucol) == 0 {
		prefs.SetString("clustcol", "num_cluster")
	}

	// chart dot size
	dotsize := binding.BindPreferenceString("dotsize", prefs) // set the link to preferences for rotation
	ds, _ := dotsize.Get()
	if len(ds) == 0 {
		prefs.SetString("dotsize", "3")
	}

	// 2D interactive plot dot size
	ds2D := binding.BindPreferenceString("2Ddotsize", prefs) // set the link to 2D dot size preferences
	dotsize2D, _ := ds2D.Get()
	if len(dotsize2D) == 0 {
		prefs.SetString("2Ddotsize", "3")
	}

	// get scaleFactor and rotation from pref
	sf := binding.BindPreferenceFloat("scaleFactor", prefs) // set the link to preferences for scaling factor
	scaleFactor, _ := sf.Get()
	if scaleFactor == 0 {
		prefs.SetFloat("scaleFactor", 1.)
	}

	// vulcano selection square size in pixels
	// record the vulcano selection square size in pixels in preferences
	vs := binding.BindPreferenceInt("vulcSelectSize", prefs)
	vsquare, _ := vs.Get()
	if vsquare == 0 {
		prefs.SetInt("vulcSelectSize", 20)
	}

	// Dot opacity
	prefs.SetFloat("dotOpacity", 255)

	// cluster opacity
	prefs.SetFloat("clustOpacity", 255)

	// vulcano default gradien
	gradExpression := binding.BindPreferenceString("gradExpression", prefs) // pref binding for the expression gradien to avoid reset for each vulcano dot
	selGrad, _ := gradExpression.Get()
	if len(selGrad) == 0 {
		prefs.SetString("gradExpression", "Turbo")
	}

	// plot background color
	initBCKGColors([]string{"unselR", "unselG", "unselB", "unselA"})
	// plot foreground colors
	initFORGColors([]string{"gateDotsR", "gateDotsG", "gateDotsB", "gateDotsA"})
	// legend text colors
	initFORGColors([]string{"legendColR", "legendColG", "legendColB", "legendColA"})
}

// init background colors
func initBCKGColors(rgba []string) {
	pref := fyne.CurrentApp().Preferences()
	if sumRGBA(rgba) == false {
		return
	}
	for _, c := range rgba {
		pref.SetInt(c, 170) // grey
	}
	pref.SetInt(rgba[3], 255) // full opacity
}

// init foreground colors
func initFORGColors(rgba []string) {
	pref := fyne.CurrentApp().Preferences()
	if sumRGBA(rgba) == false {
		return
	}
	for _, c := range rgba {
		pref.SetInt(c, 1) // black
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
	if sumRGBA == 3*255 {
		return true
	}
	return false
}

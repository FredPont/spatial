package ui

import (
	"image/color"
	"log"
	"spatial/src/filter"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/mazznoer/colorgrad"
)

// RGB color
type RGB struct {
	R, G, B uint8
}

//////////////////////////////////////////////////////
//			clusters gradients
//////////////////////////////////////////////////////

func allClustColors(nbCluster int) []RGB {
	pref := fyne.CurrentApp().Preferences()
	clustIndex := filter.FillSliceInt(nbCluster) // []int{0,1,...,n} slice of increasing int of nbcluster numbers
	shuf := binding.BindPreferenceBool("shuffClustgrad", pref)
	shuffle, _ := shuf.Get()

	if shuffle {
		filter.ShuffleInt(clustIndex)
		log.Println("shuffled color order", clustIndex)
	}

	cG := binding.BindPreferenceString("clusterGradient", pref)
	clusterGrad, _ := cG.Get()

	if strings.Contains(clusterGrad, "json") {
		return customGrad(nbCluster, clustIndex, clusterGrad)
	}

	switch clusterGrad {
	case "Turbo":
		return clusTurbo(nbCluster, clustIndex)
	case "Rainbow":
		return clusRainbow(nbCluster, clustIndex)
	case "Sinebow":
		return clusSinebow(nbCluster, clustIndex)
	default:
		return clusTurbo(nbCluster, clustIndex)
	}
}

func customGrad(nbCluster int, clustIndex []int, clusterGrad string) []RGB {
	pal := loadPalette("src/palette/" + clusterGrad)

	// check if the palette is too small to plot all clusters
	if len(pal.Pal) < nbCluster {
		log.Println("palette length ", len(pal.Pal), " is too small to plot ", nbCluster, " clusters ! Turbo gradient is used instead.")
		return clusTurbo(nbCluster, clustIndex)
	}
	RGBarray := make([]RGB, nbCluster)
	for i, cluster := range clustIndex {
		R, err := strconv.ParseUint(pal.Pal[cluster].R, 10, 32)
		checkError("cannot convert json file to RGB", err)
		G, err := strconv.ParseUint(pal.Pal[cluster].G, 10, 32)
		checkError("cannot convert json file to RGB", err)
		B, err := strconv.ParseUint(pal.Pal[cluster].B, 10, 32)
		checkError("cannot convert json file to RGB", err)

		RGBarray[i] = RGB{uint8(R), uint8(G), uint8(B)}
	}
	exportPalette(clustIndex, RGBarray)
	return RGBarray
}

// export the clusters # folowed by the RVB palette
func exportPalette(clustIndex []int, RGBarray []RGB) {
	table := [][]string{{"Color_Number", "R", "G", "B", "R_code"}}
	for i, j := range clustIndex {
		R := strconv.FormatUint(uint64(RGBarray[i].R), 10)
		G := strconv.FormatUint(uint64(RGBarray[i].G), 10)
		B := strconv.FormatUint(uint64(RGBarray[i].B), 10)
		nb := strconv.Itoa(j)

		table = append(table, []string{nb, R, G, B, "rgb(" + R + "," + G + "," + B + ", maxColorValue=255)"})
	}
	out := filter.FormatOutFile("palette", "", ".csv")
	// export the palette as CSV in the "palettes" dir
	filter.WriteCSV(table, "palettes/"+out)
	log.Println("palette saved in palettes/", out)
}

func clusRainbow(nbCluster int, clustIndex []int) []RGB {
	RGBarray := make([]RGB, nbCluster)
	for i, cluster := range clustIndex {
		grad := colorgrad.Rainbow().Sharp(uint(nbCluster+1), 0.2)
		color := rgbModel(grad.Colors(uint(nbCluster + 1))[cluster])
		RGBarray[i] = color
	}
	//log.Println("clustIndex", clustIndex, "RGB", RGBarray)
	exportPalette(clustIndex, RGBarray)
	return RGBarray
}

func clusTurbo(nbCluster int, clustIndex []int) []RGB {
	RGBarray := make([]RGB, nbCluster)
	for i, cluster := range clustIndex {
		grad := colorgrad.Turbo().Sharp(uint(nbCluster+1), 0.2)
		color := rgbModel(grad.Colors(uint(nbCluster + 1))[cluster])
		RGBarray[i] = color
	}
	//log.Println("clustIndex", clustIndex, "RGB", RGBarray)
	exportPalette(clustIndex, RGBarray)
	return RGBarray
}

func clusSinebow(nbCluster int, clustIndex []int) []RGB {
	RGBarray := make([]RGB, nbCluster)
	for i, cluster := range clustIndex {
		grad := colorgrad.Sinebow().Sharp(uint(nbCluster+1), 0.2)
		color := rgbModel(grad.Colors(uint(nbCluster + 1))[cluster])
		RGBarray[i] = color
	}
	//log.Println("clustIndex", clustIndex, "RGB", RGBarray)
	exportPalette(clustIndex, RGBarray)
	return RGBarray
}

// ClusterColors computes the color of cluster number "cluster"
// for a total number of clusters "nbCluster"
func ClusterColors(nbCluster, cluster int) RGB {
	grad := colorgrad.Rainbow().Sharp(uint(nbCluster+1), 0.2)
	return rgbModel(grad.Colors(uint(nbCluster + 1))[cluster])
}

//////////////////////////////////////////////////////
// 			Expression gradients
//////////////////////////////////////////////////////
// credits https://github.com/mazznoer/colorgrad

// pre-build the user selected gradient
func preBuildGradient(gradien string) colorgrad.Gradient {
	switch gradien {
	case "Turbo":
		return colorgrad.Turbo()
	case "Viridis":
		return colorgrad.Viridis()
	case "White - Red":
		return colorgrad.Reds()
	case "Yellow - Red":
		return colorgrad.YlOrRd()
	case "Purple - Red":
		return colorgrad.PuRd()
	case "Inferno":
		return colorgrad.Inferno()
	case "Plasma":
		return colorgrad.Plasma()
	case "Red - Yellow ":
		grad, _ := colorgrad.NewGradient().
			HtmlColors("#800026", "#bd0026", "#e31a1c", "#fc4e2a", "#fd8d3c", "#feb24c", "#fed976", "#ffeda0", "#ffffcc").
			Build()
		return grad
	case "Custom":
		colArray := filter.ReadGradient("src/gradient/custom.csv")
		grad, _ := colorgrad.NewGradient().
			HtmlColors(colArray...).
			Build()
		return grad
	default:
		return colorgrad.Reds()
	}
}

// WRgradien white -> red gradien
func WRgradien(val float64) RGB {
	grad := colorgrad.Reds()
	return rgbModel(grad.At(val))
}

// TurboGradien black blue green yellow red gradien
func TurboGradien(val float64) RGB {
	grad := colorgrad.Turbo()
	return rgbModel(grad.At(val))
}

// YlRdGradien yellow red gradien
func YlRdGradien(val float64) RGB {
	grad := colorgrad.YlOrRd()
	return rgbModel(grad.At(val))
}

// InferGrad rainbow gradien
func InferGrad(val float64) RGB {
	grad := colorgrad.Inferno()
	return rgbModel(grad.At(val))
}

// ViridisGrad Viridis rainbow gradien
func ViridisGrad(val float64) RGB {
	grad := colorgrad.Viridis()
	return rgbModel(grad.At(val))
}

// PuRdGradien purple Red
func PuRdGradien(val float64) RGB {
	grad := colorgrad.PuRd()
	return rgbModel(grad.At(val))
}

// PlasmaGradien blue orange yellow - colorgrad.Plasma()
func PlasmaGradien(val float64) RGB {
	grad := colorgrad.Plasma()
	return rgbModel(grad.At(val))
}

// BYRGradien blue yellow red
func BYRGradien(val float64) RGB {
	grad, _ := colorgrad.NewGradient().
		HtmlColors("#1726BD", "03F6FA", "03FA03", "#E8FB02", "FAE403", "FA6803", "#FA0303").
		Build()

	return rgbModel(grad.At(val))
}

// RDYLGradien yellow red
func RDYLGradien(val float64) RGB {
	grad, _ := colorgrad.NewGradient().
		HtmlColors("#800026", "#bd0026", "#e31a1c", "#fc4e2a", "#fd8d3c", "#feb24c", "#fed976", "#ffeda0", "#ffffcc").
		Build()

	return rgbModel(grad.At(val))
}

// CustomGradien from src/gradient/custom.csv
func CustomGradien(val float64) RGB {
	colArray := filter.ReadGradient("src/gradient/custom.csv")
	grad, _ := colorgrad.NewGradient().
		HtmlColors(colArray...).
		Build()

	return rgbModel(grad.At(val))
}

func rgbModel(c color.Color) RGB {
	r, g, b, _ := c.RGBA()
	return RGB{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}
}

//////////////////////
// 2D plot gradients
//////////////////////

// these functions return a gradien, not a RGB value
// the NRGBA value is computed by the dotColors function

// FULLGradien Full color gradient blue to red for 2D plot
func FULLGradien() colorgrad.Gradient {
	grad, _ := colorgrad.NewGradient().
		Colors(
			color.RGBA{0, 206, 209, 255},
			color.RGBA{255, 105, 180, 255},
			colorful.Color{R: 0.274, G: 0.5, B: 0.7},
			colorful.Hsv(50, 1, 1),
			colorful.Hsv(348, 0.9, 0.8),
		).
		Build()
	return grad
}

// YELLBLUEGradien color gradient yellow to rblue for 2D plot
func YELLBLUEGradien() colorgrad.Gradient {
	grad, _ := colorgrad.NewGradient().
		HtmlColors("gold", "hotpink", "darkturquoise").
		Build()
	return grad
}

// RAINBOWGradien color gradient for 2D plot
func RAINBOWGradien() colorgrad.Gradient {
	grad := colorgrad.Rainbow()
	return grad
}

// SINEBOWGradien color gradient for 2D plot
func SINEBOWGradien() colorgrad.Gradient {
	grad := colorgrad.Sinebow()
	return grad
}

// TURBOGradien color gradient for 2D plot
func TURBOGradien() colorgrad.Gradient {
	grad := colorgrad.Turbo()
	return grad
}

// PLASMAGradien color gradient for 2D plot
func PLASMAGradien() colorgrad.Gradient {
	grad := colorgrad.Plasma()
	return grad
}

// WARMGradien color gradient for 2D plot
func WARMGradien() colorgrad.Gradient {
	grad := colorgrad.Warm()
	return grad
}

// HEGradien color gradient for 2D plot - not used
// func HEGradien() colorgrad.Gradient {

// 	grad, _ := colorgrad.NewGradient().
// 		Colors(
// 			color.RGBA{255, 255, 51, 255}, // yellow
// 			color.RGBA{0, 213, 255, 255},  // cyan
// 			//color.RGBA{255, 255, 255, 255}, // white
// 			color.RGBA{255, 0, 0, 255},     // red
// 			color.RGBA{128, 255, 128, 255}, // green
// 			color.RGBA{255, 102, 25, 255},  // orange

// 		).
// 		Build()
// 	return grad
// }

// HEcustom return an array of custom colors
func HEcustom() []color.NRGBA {
	return []color.NRGBA{
		{255, 255, 51, 255},  // yellow
		{255, 0, 0, 255},     // red
		{0, 213, 255, 255},   // cyan
		{0, 204, 0, 255},     // green
		{255, 102, 25, 255},  // orange
		{205, 219, 209, 255}, // grey
		{0, 0, 0, 255},       // black
		//{255, 255, 255, 255}, // white
	}
}

func chooseHE(gateIndex int) color.NRGBA {
	return HEcustom()[gateIndex]
}

// grad return the gradien function with name "gradien"
func grad2D(gradient string) colorgrad.Gradient {
	switch gradient {
	case "Rainbow":
		return RAINBOWGradien()
	case "Sinebow":
		return SINEBOWGradien()
	case "Turbo":
		return TURBOGradien()
	case "Plasma":
		return PLASMAGradien()
	case "Warm":
		return WARMGradien()
	case "FullColor":
		return FULLGradien()
	case "Gold - Turquoise":
		return YELLBLUEGradien()
	// case "Hematoxilin Eosine":
	// 	return HEGradien()
	default:
		return FULLGradien()
	}

}

package ui

import (
	"image/color"

	"github.com/mazznoer/colorgrad"
)

// RGB color
type RGB struct {
	R, G, B uint8
}

// ClusterColors computes the color of cluster number "cluster"
// for a total number of clusters "nbCluster"
func ClusterColors(nbCluster, cluster int) RGB {
	grad := colorgrad.Rainbow().Sharp(uint(nbCluster+1), 0.2)
	return rgbModel(grad.Colors(uint(nbCluster + 1))[cluster])
}

// credits https://github.com/mazznoer/colorgrad

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

// BYRGradien blue yellow red
func BYRGradien(val float64) RGB {
	grad, _ := colorgrad.NewGradient().
		HtmlColors("#1726BD", "03F6FA", "03FA03", "#E8FB02", "FAE403", "FA6803", "#FA0303").
		Build()

	return rgbModel(grad.At(val))
}
func rgbModel(c color.Color) RGB {
	r, g, b, _ := c.RGBA()
	return RGB{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}
}

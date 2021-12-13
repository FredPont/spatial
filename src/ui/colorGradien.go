package ui

import (
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
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

// RAINBOWGradien color gradient yellow to rblue for 2D plot
func RAINBOWGradien() colorgrad.Gradient {
	grad := colorgrad.Rainbow()
	return grad
}

// SINEBOWGradien color gradient yellow to rblue for 2D plot
func SINEBOWGradien() colorgrad.Gradient {
	grad := colorgrad.Sinebow()
	return grad
}

// TURBOGradien color gradient yellow to rblue for 2D plot
func TURBOGradien() colorgrad.Gradient {
	grad := colorgrad.Turbo()
	return grad
}

// PLASMAGradien color gradient yellow to rblue for 2D plot
func PLASMAGradien() colorgrad.Gradient {
	grad := colorgrad.Plasma()
	return grad
}

// WARMGradien color gradient yellow to rblue for 2D plot
func WARMGradien() colorgrad.Gradient {
	grad := colorgrad.Warm()
	return grad
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
	default:
		return FULLGradien()
	}

}

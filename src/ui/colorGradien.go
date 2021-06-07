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

func rgbModel(c color.Color) RGB {
	r, g, b, _ := c.RGBA()
	return RGB{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}
}

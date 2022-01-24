package ui

import (
	"fmt"
	"testing"
)

func TestloadPalette(t *testing.T) {
	t.Run(fmt.Sprintf("palette"), func(t *testing.T) {

		cs := loadPalette("src/palette/palette1.json")
		t.Log(cs)

	})
}

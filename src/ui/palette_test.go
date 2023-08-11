package ui

import (
	"fmt"
	"testing"
)

func TestLoadPalette(t *testing.T) {
	t.Run(fmt.Sprintf("palette"), func(t *testing.T) {

		cs := loadPalette("../palette/palette1.json")
		t.Log(cs)

	})
}

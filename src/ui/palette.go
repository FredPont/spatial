package ui

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// RGBstr color in string not yet converted to uint8
type RGBstr struct {
	R, G, B string
}

// Palette stores a slice of RGB values in string
type Palette struct {
	Pal []RGBstr `json:"pal"`
}

// loadPalette read palette file from json file
func loadPalette(fname string) Palette {
	var plt Palette
	fp, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	bytes, err := ioutil.ReadAll(fp)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, &plt)
	if err != nil {
		panic(err)
	}
	//fmt.Println(plt)
	return plt
}

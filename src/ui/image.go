//credits : https://stackoverflow.com/questions/21741431/get-image-size-with-golang

package ui

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"os"
	"path/filepath"
)

const dirToScan string = "image/"

func ImgSize() (string, int, int) {
	// read files in "image" dir
	fmt.Println("reading", dirToScan)
	files, _ := ioutil.ReadDir(dirToScan)
	if len(files) == 0 {
		fmt.Println("image not found in", dirToScan)
		os.Exit(1)
	}

	for _, imgFile := range files {

		if reader, err := os.Open(filepath.Join(dirToScan, imgFile.Name())); err == nil {
			defer reader.Close()
			im, _, err := image.DecodeConfig(reader)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", imgFile.Name(), err)
				continue
			}
			fmt.Printf("%s %d %d\n", imgFile.Name(), im.Width, im.Height)
			return imgFile.Name(), im.Width, im.Height
		} else {
			fmt.Println("Impossible to open the file:", err)
		}
	}
	return "", 0, 0
}

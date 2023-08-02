package ui

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"spatial/src/filter"
	"sync"
)

// MergeIMG merges png images together
func MergeIMG(dir string, pathOut string) {
	log.Println("merging layers")
	imgs := filter.ReadDir(dir)

	// Open the first image
	file1, err := os.Open(dir + imgs[0])
	if err != nil {
		panic(err)
	}
	defer file1.Close()

	img1, _, err := image.Decode(file1)
	if err != nil {
		panic(err)
	}
	// Create a new draw image with the same dimensions as the original
	finalImg := image.NewRGBA(img1.Bounds())

	// Draw the original image onto the new draw image
	draw.Draw(finalImg, img1.Bounds(), img1, image.Point{0, 0}, draw.Src)

	// final image
	for _, img := range imgs[1:] {

		// Open the second image
		file2, err := os.Open(dir + img)
		if err != nil {
			panic(err)
		}
		defer file2.Close()

		img2, _, err := image.Decode(file2)
		if err != nil {
			panic(err)
		}

		// Draw the original image onto the new draw image
		draw.Draw(finalImg, img1.Bounds(), img2, image.Point{0, 0}, draw.Over)
	}

	// Save the stacked image to a file
	outputFile, err := os.Create(pathOut)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, finalImg)
	if err != nil {
		panic(err)
	}
	log.Println("layers merged !")
}

// Merge2Img stack two images on top of each other
func Merge2Img(img1, img2 image.Image, filename1, filename2 string, fileOut string, wg *sync.WaitGroup) {
	//log.Println("merge ", filename1, " ", filename2, " into ", fileOut)
	// Create a new draw image with the same dimensions as the original
	finalImg := image.NewRGBA(img1.Bounds())

	// Draw the original image onto the new draw image
	draw.Draw(finalImg, img1.Bounds(), img1, image.Point{0, 0}, draw.Src)

	// Draw the  image 2 onto the new draw image
	draw.Draw(finalImg, img1.Bounds(), img2, image.Point{0, 0}, draw.Over)

	// Save the stacked image to a file
	outputFile, err := os.Create(fileOut)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, finalImg)
	if err != nil {
		panic(err)
	}
	//log.Println("layers merged !")
	// remove the 2 initial images
	filter.RmFile(filename1) // delete file1
	filter.RmFile(filename2) // delete file2
	wg.Done()
}

// MergeIMG merges png images together 2 by 2 using multithread in temp rir and then remove all files in temp dir
func MTmergeIMG(dir string, pathOut string) {

	imgs := filter.ReadDir(dir)

	if len(imgs) < 4 {
		// use single thread
		MergeIMG(dir, pathOut)
		return
	}
	log.Println("multithread merging layers")
	//log.Println(imgs)
	var wg sync.WaitGroup
	// while there are files to merge in the temp directory
	ct := 1
	for len(imgs) > 1 {
		imgTomerge := len(imgs)
		for i := 0; i < len(imgs); i += 2 {

			// Open the first image
			file1, err := os.Open(dir + imgs[i])
			if err != nil {
				panic(err)
			}
			defer file1.Close()

			img1, _, err := image.Decode(file1)
			if err != nil {
				panic(err)
			}

			// Open the second image
			file2, err := os.Open(dir + imgs[i+1])
			if err != nil {
				panic(err)
			}
			defer file1.Close()

			img2, _, err := image.Decode(file2)
			if err != nil {
				panic(err)
			}
			wg.Add(1)
			fileOut := dir + fmt.Sprint(ct) + "_temp.png"
			go Merge2Img(img1, img2, dir+imgs[i], dir+imgs[i+1], fileOut, &wg)

			imgTomerge = imgTomerge - 2

			// if the number of image to merge is 1, stop the loop
			if imgTomerge == 1 {
				//log.Println("break merge loop")
				break
			}

			//log.Println(i, " len img ", len(imgs), " ct ", ct)
			ct++
		}
		ct = 2 * ct
		wg.Wait()
		imgs = filter.ReadDir(dir)
		//log.Println(imgs)

	}
	filter.MvFile(dir+imgs[0], pathOut)
}

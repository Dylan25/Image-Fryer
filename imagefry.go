// imagefry applies pseudo random filter to an image
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	ImageFile, timesFry := inputParseAndOpen()
	defer ImageFile.Close()

	imageData, imageType, err := image.Decode(ImageFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "perlin: %v\n", err)
	}

	ImageFile.Seek(0, 0)

	imgCfg, _, err := image.DecodeConfig(ImageFile)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ImageFile.Seek(0, 0)

	newImg := randFilter(imageData, imageType, imgCfg, timesFry)

	outputFile, err := os.Create("fryd" + os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "perlin output error: %s\n", err)
	}

	if imageType == "png" {
		png.Encode(outputFile, newImg)
	} else if imageType == "jpeg" {
		jpeg.Encode(outputFile, newImg, nil)
	} else {
		fmt.Println("ERROR: unrecognized file format")
	}
	outputFile.Close()

	fmt.Printf("output written to %s\n", "fryd" + os.Args[1])
}

func inputParseAndOpen() (*os.File, int) {
	if len(os.Args) <= 1 {
		fmt.Fprint(os.Stderr, "ERROR: please provide a filename\n")
		fmt.Println("USAGE: 'imagefry.exe image.jpg/png #times_fryd'")
		os.Exit(1)
	}
	if strings.HasSuffix(os.Args[1], ".png") || strings.HasSuffix(os.Args[1], ".jpg") {
		imageFile, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not open file, %v\n", err)
			os.Exit(1)
		}

		if len(os.Args) == 3 {
			numfry := os.Args[2]
			intnumfry, err := strconv.Atoi(numfry)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: please enter a number of times to fry, %s\n", err)
				os.Exit(1)
			}
			return imageFile, intnumfry
		} else if len(os.Args) > 3 {
			fmt.Fprint(os.Stderr, "ERROR: too many arguments")
			os.Exit(1)
		} else {
			return imageFile, 1
		}

	} else {
		fmt.Fprint(os.Stderr, "ERROR: please provide a filename\n")
		fmt.Println("USAGE: 'imagefry.exe image.jpg/png #times_fryd'")
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Argument parse error, try again")
	os.Exit(1)
	return nil, 0
}

func randFilter(imageData image.Image, imageType string, imgCfg image.Config, timesFry int) image.Image {
	// copy old image to a new template
	alteredImage := image.NewRGBA(imageData.Bounds())
	draw.Draw(alteredImage, imageData.Bounds(), imageData, image.Point{}, draw.Over)

	width := imgCfg.Width
	height := imgCfg.Height

	// apply random changes to the image
	for i := 0; i < timesFry; i++ {
		for y := 0; y < height; y++ {
			rand.Seed(time.Now().UTC().UnixNano())
			for x := 0; x < width; x++ {
				r, g, b, a := alteredImage.At(x, y).RGBA()
				newColor := color.RGBA{randColor(uint8(r)), randColor(uint8(g)), randColor(uint8(b)), uint8(a)}
				alteredImage.Set(x, y, newColor)
			}
		}
	}
	return alteredImage
}

func randColor(origColor uint8) uint8 {
	key := rand.Intn(1)
	if key == 0 {
		return origColor + uint8(rand.Intn(10))
	} else {
		return origColor - uint8(rand.Intn(10))
	}
}

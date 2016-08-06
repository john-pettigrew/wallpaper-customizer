package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/john-pettigrew/wallpaper-customizer/imgtransforms"
	"golang.org/x/image/draw"
)

func main() {

	var flipImage, invertColors bool
	var outputPath, inputPath, maskPath, maskDrawPath string

	//cmd line arguments
	flag.BoolVar(&flipImage, "f", false, "Flip the image that will be drawn.")
	flag.BoolVar(&invertColors, "i", false, "Invert the colors for image that will be drawn.")
	flag.StringVar(&outputPath, "out", "", "Output file path (required)")
	flag.StringVar(&inputPath, "in", "", "Input file path (required)")
	flag.StringVar(&maskPath, "m", "", "Mask file path")
	flag.StringVar(&maskDrawPath, "mi", "", "File path for an image to draw using a supplied mask")
	flag.Parse()

	//validation
	if outputPath == "" {
		errorAndExit("Output file is required\n")
	}

	if inputPath == "" {
		errorAndExit("Input file is required\n")
	}

	if maskPath != "" || maskDrawPath != "" {
		if maskPath == "" || maskDrawPath == "" {
			errorAndExit("Mask path and mask image to draw path are both required to draw mask\n")
		}
	}

	//Read in input image
	dst, err := readImage(inputPath)
	if err != nil {
		errorAndExit("Error reading input image\n")
	}

	//apply transforms
	if flipImage {
		dst = imgtransforms.Flip(dst)
	}

	if invertColors {
		dst = imgtransforms.InvertColors(dst)
	}

	//Convert dst
	dstB := dst.Bounds()
	finalDst := image.NewRGBA(dstB)
	draw.Draw(finalDst, finalDst.Bounds(), dst, dstB.Min, draw.Src)

	//apply mask
	if maskPath != "" && maskDrawPath != "" {

		mask, err := readImage(maskPath)
		if err != nil {
			errorAndExit("Error reading mask\n")
		}

		maskImage, err := readImage(maskDrawPath)
		if err != nil {
			errorAndExit("Error reading mask image\n")
		}

		//Scale mask
		finalMask := image.NewRGBA(maskImage.Bounds())
		draw.ApproxBiLinear.Scale(finalMask, maskImage.Bounds(), mask, mask.Bounds(), draw.Over, nil)

		//Draw our image
		draw.DrawMask(finalDst, finalDst.Bounds(), maskImage, image.ZP, finalMask, image.ZP, draw.Over)
	}

	//Create output file
	output, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}

	//Save output file
	options := jpeg.Options{Quality: 100}
	err = jpeg.Encode(output, finalDst, &options)
	if err != nil {
		log.Fatal(err)
	}
}

func readImage(file string) (image.Image, error) {
	imageFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	image, _, err := image.Decode(imageFile)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func errorAndExit(errStr string) {
	fmt.Println(errStr)
	flag.Usage()
	os.Exit(1)
}

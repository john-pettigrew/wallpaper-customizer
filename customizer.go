package main

import (
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

	//Check parameters
	if len(os.Args) < 4 {
		fmt.Println("Usage: image-customizer [dst] [mask] [output]")
		return
	}

	dstPath := os.Args[1]
	maskPath := os.Args[2]
	outputPath := os.Args[3]

	//Read in images
	dst, err := readImage(dstPath)
	if err != nil {
		log.Fatal("Error reading dst")
	}

	mask, err := readImage(maskPath)
	if err != nil {
		log.Fatal("Error reading src")
	}

	//Scale mask
	finalMask := image.NewRGBA(dst.Bounds())
	draw.ApproxBiLinear.Scale(finalMask, dst.Bounds(), mask, mask.Bounds(), draw.Over, nil)

	//Create changed dst
	changedDst := imgtransforms.Flip(dst)
	changedDst = imgtransforms.InvertColors(changedDst)

	//Convert dst
	dstB := dst.Bounds()
	finalDst := image.NewRGBA(dstB)
	draw.Draw(finalDst, finalDst.Bounds(), dst, dstB.Min, draw.Src)

	//Draw our image
	draw.DrawMask(finalDst, finalDst.Bounds(), changedDst, image.ZP, finalMask, image.ZP, draw.Over)

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

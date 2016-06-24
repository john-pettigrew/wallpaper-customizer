package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/png"
	"log"
	"os"

	"golang.org/x/image/draw"
)

type MyColor struct {
	r, g, b, a uint32
}

func (c MyColor) RGBA() (uint32, uint32, uint32, uint32) {
	return c.r, c.g, c.b, c.a
}

func main() {

	if len(os.Args) < 4 {
		fmt.Println("Usage: image-customizer [dst] [mask] [output]")
		return
	}

	dstPath := os.Args[1]
	maskPath := os.Args[2]
	outputPath := os.Args[3]

	dst, err := readImage(dstPath)
	if err != nil {
		log.Fatal("Error reading dst")
	}

	mask, err := readImage(maskPath)
	if err != nil {
		log.Fatal("Error reading src")
	}

	// scale mask
	finalMask := image.NewRGBA(dst.Bounds())
	draw.ApproxBiLinear.Scale(finalMask, dst.Bounds(), mask, mask.Bounds(), draw.Over, nil)

	// create changed dst
	changedDst := flipImage(dst)
	changedDst = invertImageColors(changedDst)

	dstB := dst.Bounds()
	finalDst := image.NewRGBA(image.Rect(0, 0, dstB.Dx(), dstB.Dy()))
	draw.Draw(finalDst, finalDst.Bounds(), dst, dstB.Min, draw.Src)

	draw.DrawMask(finalDst, finalDst.Bounds(), changedDst, image.ZP, finalMask, image.ZP, draw.Over)

	output, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}

	options := jpeg.Options{Quality: 100}
	err = jpeg.Encode(output, finalDst, &options)
	if err != nil {
		log.Fatal(err)
	}
}

func readImage(file string) (image.Image, error) {
	//dest
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

func flipImage(input image.Image) image.Image {

	//create new image
	bounds := input.Bounds()
	newImg := image.NewRGBA(bounds)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			newImg.Set(bounds.Max.X-x, bounds.Max.Y-y, input.At(x, y))
		}
	}

	return newImg
}

func invertImageColors(input image.Image) image.Image {

	//create new image
	bounds := input.Bounds()
	newImg := image.NewRGBA(bounds)
	var pixelColor color.Color
	var r, g, b, a uint32
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			r, g, b, a = input.At(x, y).RGBA()
			pixelColor = MyColor{
				r: 255 - r,
				g: 255 - g,
				b: 255 - b,
				a: 255 - a,
			}
			newImg.Set(x, y, pixelColor)
		}
	}

	return newImg
}

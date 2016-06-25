package imgtransforms

import (
	"image"
	"image/color"
)

type pixelColor struct {
	r, g, b, a uint32
}

func (c pixelColor) RGBA() (uint32, uint32, uint32, uint32) {
	return c.r, c.g, c.b, c.a
}

//Flip returns a copy of input that has been flipped horizontally and vertically.
func Flip(input image.Image) image.Image {

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

//InvertColors returns a copy of input that has its colors inverted.
func InvertColors(input image.Image) image.Image {

	//create new image
	bounds := input.Bounds()
	newImg := image.NewRGBA(bounds)

	var currentPixelColor color.Color
	var r, g, b, a uint32
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			r, g, b, a = input.At(x, y).RGBA()
			currentPixelColor = pixelColor{
				r: 0xffff - r,
				g: 0xffff - g,
				b: 0xffff - b,
				a: a,
			}
			newImg.Set(x, y, currentPixelColor)
		}
	}

	return newImg
}

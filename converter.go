package imgedit

import (
	"image"
	"math"
)

// Converter interface for image edit
type Converter interface {
	Resize(x, y int)
	ResizeRatio(ratio float64)
	Trim(left, top, width, height int)
	ReverseX()
	ReverseY()
	Grayscale()
	Convert() image.Image
}

type converter struct {
	image.Image
}

// NewConverter create converter
func NewConverter(image image.Image) Converter {
	return &converter{image}
}

// Resize resize the image
func (c *converter) Resize(resizeX, resizeY int) {
	dst := image.NewRGBA(image.Rect(0, 0, resizeX, resizeY))
	dstSize := dst.Bounds().Size()
	xRate, yRate := float64(c.Bounds().Dx())/float64(dstSize.X), float64(c.Bounds().Dy())/float64(dstSize.Y)
	for x := 0; x < dstSize.X; x++ {
		for y := 0; y < dstSize.Y; y++ {
			srcX, srcY := int(math.Round(float64(x)*xRate)), int(math.Round(float64(y)*yRate))
			dst.Set(x, y, c.Image.At(srcX, srcY))
		}
	}
	c.Image = dst
}

// ResizeRatio resize the image with ratio
func (c *converter) ResizeRatio(ratio float64) {
	dst := image.NewRGBA(image.Rect(0, 0, int(math.Round(float64(c.Bounds().Dx())*ratio)), int(math.Round(float64(c.Bounds().Dy())*ratio))))
	dstSize := dst.Bounds().Size()
	xRate, yRate := 1/ratio, 1/ratio
	for x := 0; x < dstSize.X; x++ {
		for y := 0; y < dstSize.Y; y++ {
			srcX, srcY := int(math.Round(float64(x)*xRate)), int(math.Round(float64(y)*yRate))
			dst.Set(x, y, c.Image.At(srcX, srcY))
		}
	}
	c.Image = dst
}

// Trim trim the image to the specified size
func (c *converter) Trim(left, top, width, height int) {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	dstSize := dst.Bounds().Size()
	for x := 0; x < dstSize.X; x++ {
		for y := 0; y < dstSize.Y; y++ {
			srcX, srcY := x+left, y+top
			dst.Set(x, y, c.Image.At(srcX, srcY))
		}
	}
	c.Image = dst
}

// ReverseX reverse the image about horizon
func (c *converter) ReverseX() {
	dst := image.NewRGBA(image.Rect(0, 0, c.Bounds().Dx(), c.Bounds().Dy()))
	srcSize := c.Bounds().Size()
	dstSize := dst.Bounds().Size()
	for x := 0; x < dstSize.X; x++ {
		for y := 0; y < dstSize.Y; y++ {
			srcX, srcY := srcSize.X-x, y
			dst.Set(x, y, c.Image.At(srcX, srcY))
		}
	}
	c.Image = dst
}

// ReverseY reverse the image about vertical
func (c *converter) ReverseY() {
	dst := image.NewRGBA(image.Rect(0, 0, c.Bounds().Dx(), c.Bounds().Dy()))
	srcSize := c.Bounds().Size()
	dstSize := dst.Bounds().Size()
	for x := 0; x < srcSize.X; x++ {
		for y := 0; y < srcSize.Y; y++ {
			srcX, srcY := x, dstSize.Y-y
			dst.Set(x, y, c.Image.At(srcX, srcY))
		}
	}
	c.Image = dst
}

// Grayscale change the image color to grayscale
func (c *converter) Grayscale() {
	dst := image.NewGray(c.Bounds())
	dstSize := dst.Bounds().Size()
	for x := 0; x < dstSize.X; x++ {
		for y := 0; y < dstSize.Y; y++ {
			dst.Set(x, y, c.Image.At(x, y))
		}
	}
	c.Image = dst
}

// Convert get convert image
func (c *converter) Convert() image.Image {
	return c.Image
}

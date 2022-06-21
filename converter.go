package imgedit

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"math"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	// DefaultTtfFilePath used when font is not specified in Options
	DefaultTtfFilePath = "./assets/font/07LogoTypeGothic7.ttf"

	// DefaultFontSize used when font size is not specified in Options
	DefaultFontSize = 100
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
	dst := image.NewRGBA(image.Rect(0, 0, c.Bounds().Dx(), c.Bounds().Dy()))
	dstSize := dst.Bounds().Size()
	grayModel := color.GrayModel
	alphaModel := color.AlphaModel
	for x := 0; x < dstSize.X; x++ {
		for y := 0; y < dstSize.Y; y++ {
			dstColor := c.Image.At(x, y)
			_, _, _, a := dstColor.RGBA()
			if a <= 0 {
				dst.Set(x, y, alphaModel.Convert(dstColor))
			} else {
				dst.Set(x, y, grayModel.Convert(dstColor))
			}
		}
	}
	c.Image = dst
}

// Options options for AddString
type Options struct {
	// TrueTypeFont use ReadTtf to get font
	TrueTypeFont    *truetype.Font
	TrueTypeOptions *truetype.Options
	// Point left top = (0px, 0px)
	Point     *image.Point
	FontColor *image.Uniform
}

// Face get font.Face
func (o *Options) Face() font.Face {
	return truetype.NewFace(o.TrueTypeFont, o.TrueTypeOptions)
}

// AddString add string on current Image
func (c *converter) AddString(text string, options *Options) {
	if options == nil {
		options = &Options{}
	}
	if options.TrueTypeFont == nil {
		options.TrueTypeFont, _ = ReadTtf(DefaultTtfFilePath)
	}
	if options.TrueTypeOptions == nil {
		options.TrueTypeOptions = &truetype.Options{
			Size: DefaultFontSize,
		}
	}
	if options.FontColor == nil {
		options.FontColor = image.Black
	}

	// copy base image
	dst := image.NewRGBA(image.Rect(0, 0, c.Bounds().Dx(), c.Bounds().Dy()))
	draw.Draw(dst, image.Rect(0, 0, c.Bounds().Dx(), c.Bounds().Dy()), c.Image, image.Point{}, draw.Over)

	// add string on base image
	face := options.Face()
	drawer := &font.Drawer{
		Dst:  dst,
		Src:  options.FontColor,
		Face: face,
	}

	// default center
	centerX, centerY := (fixed.I(dst.Bounds().Dx())-drawer.MeasureString(text))/2, (fixed.I(dst.Bounds().Dy())+(face.Metrics().Ascent+face.Metrics().Descent)/2)/2

	// notice : Dot values are determined by feeling.
	if options.Point == nil {
		drawer.Dot.X, drawer.Dot.Y = centerX, centerY
	} else {
		drawer.Dot.X, drawer.Dot.Y = fixed.I(options.Point.X)-drawer.MeasureString(text)/2, fixed.I(options.Point.Y)+((face.Metrics().Ascent+face.Metrics().Descent)/2)/2
	}

	drawer.DrawString(text)
	c.Image = dst
}

// Convert get convert image
func (c *converter) Convert() image.Image {
	return c.Image
}

// ReadTtf return ttf from file path
func ReadTtf(ttfFilePath string) (*truetype.Font, error) {
	fontFile, err := ioutil.ReadFile(ttfFilePath)
	if err != nil {
		return nil, err
	}
	ttf, err := freetype.ParseFont(fontFile)
	if err != nil {
		return nil, err
	}
	return ttf, nil
}

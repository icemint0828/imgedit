package imgedit

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"sort"
)

// Extension is image file extension
type Extension string

// Png is one of the supported extension
var Png = Extension("png")

// Jpeg is one of the supported extension
var Jpeg = Extension("jpeg")

// Gif is one of the supported extension
var Gif = Extension("gif")

// SupportedExtensions are supported extensions
var SupportedExtensions = []Extension{
	Png,
	Jpeg,
	Gif,
}

// SupportedExtension return true, if extension is in the SupportedExtensions
func SupportedExtension(extension Extension) bool {
	for _, e := range SupportedExtensions {
		if e == extension {
			return true
		}
	}
	return false
}

// ByteConverter interface for image edit
type ByteConverter interface {
	Converter
	WriteAs(io.Writer, Extension) error
}

type byteConverter struct {
	*converter
}

// NewByteConverter create byteConverter
func NewByteConverter(r io.Reader) (ByteConverter, Extension, error) {
	return newByteConverter(r)
}

func newByteConverter(r io.Reader) (*byteConverter, Extension, error) {
	srcImage, format, err := image.Decode(r)
	if err != nil {
		return nil, "", err
	}
	// all formats are supported.
	extension := Extension(format)
	if !SupportedExtension(extension) {
		return nil, "", errors.New(fmt.Sprintf("extension is not supported : %s", format))
	}
	return &byteConverter{converter: &converter{srcImage}}, extension, nil
}

func (b *byteConverter) WriteAs(writer io.Writer, extension Extension) error {
	switch extension {
	case Png:
		return png.Encode(writer, b.Image)
	case Jpeg:
		return jpeg.Encode(writer, b.Image, &jpeg.Options{Quality: 100})
	case Gif:
		return gifEncode(writer, b.Image, &gif.Options{NumColors: 256})
	default:
		return errors.New("extension is unsupported")
	}
}

// gifEncode wrap the original mainly due to transparency color issues.
func gifEncode(w io.Writer, m image.Image, o *gif.Options) error {
	// if m.ColorModel().(color.Palette) is not satisfied, problems occur during image encoding
	// e.g) gif.Encode transparent images.
	if _, ok := m.ColorModel().(color.Palette); ok {
		return gif.Encode(w, m, o)
	}

	// Check for bounds and size restrictions.
	b := m.Bounds()
	if b.Dx() >= 1<<16 || b.Dy() >= 1<<16 {
		return errors.New("gif: image is too large to encode")
	}

	opts := gif.Options{}
	if o != nil {
		opts = *o
	}
	if opts.NumColors < 1 || 256 < opts.NumColors {
		opts.NumColors = 256
	}
	if opts.Drawer == nil {
		opts.Drawer = draw.FloydSteinberg
	}

	dst := image.NewPaletted(b, createMyPalette(m, o.NumColors))
	myDraw(dst, m)
	return gif.EncodeAll(w, &gif.GIF{
		Image: []*image.Paletted{dst},
		Delay: []int{0},
		Config: image.Config{
			ColorModel: dst.Palette,
			Width:      b.Dx(),
			Height:     b.Dy(),
		},
	})
}

type sortedColors map[color.Color]uint

func (s sortedColors) Sort() []color.Color {
	type sortedColor struct {
		color.Color
		Count uint
	}
	var sortedColors []sortedColor
	for c, count := range s {
		sortedColors = append(sortedColors, sortedColor{Color: c, Count: count})
	}
	sort.SliceStable(sortedColors, func(i, j int) bool { return sortedColors[i].Count > sortedColors[j].Count })

	var colors []color.Color
	for _, sortedColor := range sortedColors {
		colors = append(colors, sortedColor.Color)
	}
	return colors
}

// createMyPalette create a palette with efficient colors to represent the image
func createMyPalette(m image.Image, numColors int) []color.Color {
	b := m.Bounds()
	transparentColors := sortedColors{}
	usedColors := sortedColors{}
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			// transparent colors are handled separately. draw.sqDiff would be
			// a meaningless value in transparent colors. e.g(0xffff, 0xffff, 0xffff, 0)
			c := m.At(x, y)
			_, _, _, a := c.RGBA()
			if a == 0 {
				transparentColors[c]++
			}
			if a == math.MaxUint16 {
				usedColors[c]++
			}
		}
	}

	var myPalette []color.Color
	colors := append(transparentColors.Sort(), usedColors.Sort()...)
	for _, c := range colors {
		if len(myPalette) >= numColors {
			break
		}
		myPalette = append(myPalette, c)
	}
	return myPalette
}

func myDraw(dst *image.Paletted, src image.Image) {
	b := src.Bounds()
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			c := src.At(x, y)
			_, _, _, a := c.RGBA()
			if a < math.MaxUint16 {
				c = dst.Palette[0]
			}
			dst.Set(x, y, c)
		}
	}
}

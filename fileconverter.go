package imgedit

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
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

// FileConverter interface for image edit
type FileConverter interface {
	Converter
	SaveAs(string, Extension) error
}

type fileConverter struct {
	*converter
}

// NewFileConverter create fileConverter
func NewFileConverter(srcPath string) (FileConverter, Extension, error) {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return nil, "", err
	}
	srcImage, format, err := image.Decode(srcFile)
	if err != nil {
		return nil, "", err
	}
	// all formats are supported.
	extension := Extension(format)
	if !SupportedExtension(extension) {
		return nil, "", errors.New(fmt.Sprintf("extension is not supported : %s", format))
	}
	return &fileConverter{converter: &converter{srcImage}}, extension, nil
}

func (p *fileConverter) SaveAs(dstPath string, extension Extension) error {
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	switch extension {
	case Png:
		return png.Encode(dstFile, p.Image)
	case Jpeg:
		return jpeg.Encode(dstFile, p.Image, &jpeg.Options{Quality: 100})
	case Gif:
		return gifEncode(dstFile, p.Image, &gif.Options{NumColors: 256})
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

	// replace unused color as transparent color
	myPalette := make([]color.Color, opts.NumColors)
	copy(myPalette, palette.Plan9[:opts.NumColors])
	dst := image.NewPaletted(b, myPalette)

	usedColors := map[color.Color]bool{}
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			usedColor := dst.ColorModel().Convert(m.At(x, y))
			if _, ok := usedColors[usedColor]; !ok {
				usedColors[usedColor] = true
			}
		}
	}
	for i, usedColor := range myPalette {
		if _, ok := usedColors[usedColor]; !ok {
			myPalette[i] = image.Transparent
			break
		}
	}
	opts.Drawer.Draw(dst, b, m, b.Min)
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

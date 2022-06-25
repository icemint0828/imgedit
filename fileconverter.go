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
		// if c.Image.ColorModel().(color.Palette) is not satisfied, problems occur during image encoding
		// e.g) gif.Encode transparent images.
		if _, ok := p.Image.ColorModel().(color.Palette); ok {
			return gif.Encode(dstFile, p.Image, &gif.Options{NumColors: 256})
		}

		// replace unused color as transparent color
		myPalette := make([]color.Color, 256)
		copy(myPalette, palette.Plan9)
		dst := image.NewPaletted(image.Rect(0, 0, p.Bounds().Dx(), p.Bounds().Dy()), myPalette)

		usedColors := map[color.Color]bool{}
		for x := 0; x < dst.Bounds().Dx(); x++ {
			for y := 0; y < dst.Bounds().Dy(); y++ {
				usedColor := dst.ColorModel().Convert(p.At(x, y))
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

		draw.FloydSteinberg.Draw(dst, image.Rect(0, 0, p.Bounds().Dx(), p.Bounds().Dy()), p.Image, image.Point{})

		return gifEncode(dstFile, dst, &gif.Options{NumColors: 256})
	default:
		return errors.New("extension is unsupported")
	}
}

//type TransparentQuantizer struct {
//}
//
//func (t TransparentQuantizer) Quantize(_ color.Palette, _ image.Image) color.Palette {
//	return color.Palette{image.Transparent}
//}

func gifEncode(w io.Writer, m image.Image, o *gif.Options) error {
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

	pm, _ := m.(*image.Paletted)
	if pm == nil {
		fmt.Println("■通過チェック1")
		if cp, ok := m.ColorModel().(color.Palette); ok {
			fmt.Println("■通過チェック2")
			pm = image.NewPaletted(b, cp)
			for y := b.Min.Y; y < b.Max.Y; y++ {
				for x := b.Min.X; x < b.Max.X; x++ {
					pm.Set(x, y, cp.Convert(m.At(x, y)))
				}
			}
		}
	}
	if pm == nil || len(pm.Palette) > opts.NumColors {
		fmt.Println("■通過チェック3")
		// Set pm to be a palettedized copy of m, including its bounds, which
		// might not start at (0, 0).
		//
		// TODO: Pick a better sub-sample of the Plan 9 palette.
		pm = image.NewPaletted(b, palette.Plan9[:opts.NumColors])
		if opts.Quantizer != nil {
			pm.Palette = opts.Quantizer.Quantize(make(color.Palette, 0, opts.NumColors), m)
		}
		opts.Drawer.Draw(pm, b, m, b.Min)
	}

	// When calling Encode instead of EncodeAll, the single-frame image is
	// translated such that its top-left corner is (0, 0), so that the single
	// frame completely fills the overall GIF's bounds.
	if pm.Rect.Min != (image.Point{}) {
		dup := *pm
		dup.Rect = dup.Rect.Sub(dup.Rect.Min)
		pm = &dup
	}

	fmt.Println(len(pm.Palette))

	return gif.EncodeAll(w, &gif.GIF{
		Image: []*image.Paletted{pm},
		Delay: []int{0},
		Config: image.Config{
			ColorModel: pm.Palette,
			Width:      b.Dx(),
			Height:     b.Dy(),
		},
	})
}

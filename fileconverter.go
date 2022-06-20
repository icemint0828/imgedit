package imgedit

import (
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
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
		return gif.Encode(dstFile, p.Image, &gif.Options{NumColors: 256})
	default:
		return errors.New("extension is unsupported")
	}
}

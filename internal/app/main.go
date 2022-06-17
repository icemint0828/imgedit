package app

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/icemint0828/imgedit"
)

var SupportedExtensions = []string{
	"png",
	"jpeg",
	"gif",
}

type App struct {
	subCommand    *SubCommand
	filePath      string
	fileExtension string
	convertFormat string
}

// NewApp create app
func NewApp(subCommand *SubCommand, filePath string) *App {
	return &App{
		subCommand:    subCommand,
		filePath:      filePath,
		fileExtension: filepath.Ext(filePath),
	}
}

// Run edit the image
func (a *App) Run() error {
	// load image
	loadImage, err := a.loadImage()
	if err != nil {
		return err
	}

	// convert image
	c := imgedit.NewConverter(loadImage)
	switch a.subCommand.Name {
	case SubCommandReverse.Name:
		isVertical := flagBool(OptionVertical.Name)
		if isVertical {
			c.ReverseY()
		} else {
			c.ReverseX()
		}
	case SubCommandResize.Name:
		ratio := flagFloat64(OptionRatio.Name)
		width, height := int(flagUint(OptionWidth.Name)), int(flagUint(OptionHeight.Name))
		if ratio != 0 {
			c.ResizeRatio(ratio)
		} else {
			c.Resize(width, height)
		}
	case SubCommandTrim.Name:
		left, right := int(flagUint(OptionLeft.Name)), int(flagUint(OptionRight.Name))
		bottom, top := int(flagUint(OptionBottom.Name)), int(flagUint(OptionTop.Name))
		c.Trim(left, right, bottom, top)
	case SubCommandGrayscale.Name:
		c.Grayscale()
	}

	// save image
	saveImage := c.Convert()
	err = a.saveImage(saveImage)
	if err != nil {
		return err
	}
	return nil
}

func flagUint(name string) uint {
	return flag.Lookup(name).Value.(flag.Getter).Get().(uint)
}

func flagFloat64(name string) float64 {
	return flag.Lookup(name).Value.(flag.Getter).Get().(float64)
}

func flagBool(name string) bool {
	return flag.Lookup(name).Value.(flag.Getter).Get().(bool)
}

func (a *App) loadImage() (image.Image, error) {
	file, err := os.Open(a.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	if !supportedExtension(format) {
		return nil, errors.New(fmt.Sprintf("extension is not supported : %s", format))
	}

	a.convertFormat = format
	return img, err
}

func (a *App) saveImage(img image.Image) error {
	outputDir, err := os.Getwd()
	if err != nil {
		return err
	}
	var outputFileName string
	if a.fileExtension == "" {
		outputFileName = filepath.Base(a.filePath) + "_imgedit"
	} else {
		outputFileName = strings.Replace(filepath.Base(a.filePath), a.fileExtension, "_imgedit"+a.fileExtension, 1)
	}

	outputPath := path.Join(outputDir, outputFileName)
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Printf("save convert file: %s\n", outputPath)

	switch a.convertFormat {
	case "png":
		return png.Encode(file, img)
	case "jpeg":
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 100})
	case "gif":
		return gif.Encode(file, img, &gif.Options{NumColors: 256})
	default:
		return errors.New("extension is not supported")
	}
}

func supportedExtension(convertFormat string) bool {
	for _, extension := range SupportedExtensions {
		if extension == convertFormat {
			return true
		}
	}
	return false
}

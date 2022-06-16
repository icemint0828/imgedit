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
	"strconv"
	"strings"

	"github.com/icemint0828/imgedit"
)

var SupportedExtensions = []string{
	".png",
	".jpg",
	".jpeg",
	".gif",
}

type App struct {
	*SubCommand
	FilePath string
}

// NewApp create app
func NewApp(subCommand *SubCommand, filepath string) *App {
	return &App{
		SubCommand: subCommand,
		FilePath:   filepath,
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
	switch a.SubCommand.Name {
	case SubCommandReverse.Name:
		if flag.Lookup(OptionVertical.Name).Value.String() == strconv.FormatBool(true) {
			c.ReverseY()
		} else {
			c.ReverseX()
		}
	case SubCommandResize.Name:
		ratio := flag.Lookup(OptionRatio.Name).Value.(flag.Getter).Get().(float64)
		width, _ := strconv.Atoi(flag.Lookup(OptionWidth.Name).Value.String())
		height, _ := strconv.Atoi(flag.Lookup(OptionHeight.Name).Value.String())
		if ratio != 0 {
			c.ResizeRatio(ratio)
		} else {
			c.Resize(width, height)
		}
	case SubCommandTrim.Name:
		left, _ := strconv.Atoi(flag.Lookup(OptionLeft.Name).Value.String())
		right, _ := strconv.Atoi(flag.Lookup(OptionRight.Name).Value.String())
		bottom, _ := strconv.Atoi(flag.Lookup(OptionBottom.Name).Value.String())
		top, _ := strconv.Atoi(flag.Lookup(OptionTop.Name).Value.String())
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

func (a *App) loadImage() (image.Image, error) {
	file, err := os.Open(a.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	e := filepath.Ext(a.FilePath)
	switch strings.ToLower(e) {
	case ".png":
		return png.Decode(file)
	case ".jpg", ".jpeg":
		return jpeg.Decode(file)
	case ".gif":
		return gif.Decode(file)
	default:
		return nil, errors.New("extension is not supported")
	}
}

func (a *App) saveImage(img image.Image) error {
	outputDir, err := os.Getwd()
	if err != nil {
		return err
	}
	outputFileName := strings.Replace(filepath.Base(a.FilePath), ".", "_imgedit.", 1)
	outputPath := path.Join(outputDir, outputFileName)
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	e := filepath.Ext(a.FilePath)
	fmt.Printf("save convert file: %s\n", outputPath)

	switch strings.ToLower(e) {
	case ".png":
		return png.Encode(file, img)
	case ".jpg", ".jpeg":
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 100})
	case ".gif":
		return gif.Encode(file, img, &gif.Options{NumColors: 256})
	default:
		return errors.New("extension is not supported")
	}
}

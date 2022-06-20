package app

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/icemint0828/imgedit"
)

type App struct {
	subCommand    *SubCommand
	filePath      string
	fileExtension string
	extension     imgedit.Extension
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
	c, extension, err := imgedit.NewFileConverter(a.filePath)
	if err != nil {
		return err
	}

	// convert image
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
		left, top := int(flagUint(OptionLeft.Name)), int(flagUint(OptionTop.Name))
		width, height := int(flagUint(OptionWidth.Name)), int(flagUint(OptionHeight.Name))
		c.Trim(left, top, width, height)
	case SubCommandGrayscale.Name:
		c.Grayscale()
	case SubCommandPng.Name:
		extension = imgedit.Png
	case SubCommandJpeg.Name:
		extension = imgedit.Jpeg
	case SubCommandGif.Name:
		extension = imgedit.Gif
	}

	// save image
	outputPath, err := a.getOutputPath(extension)
	if err != nil {
		return err
	}
	err = c.SaveAs(outputPath, extension)
	if err != nil {
		return err
	}
	fmt.Printf("save convert file: %s\n", outputPath)
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

func (a *App) getOutputPath(extension imgedit.Extension) (string, error) {
	outputDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	var outputFileName string
	if a.fileExtension == "" {
		outputFileName = filepath.Base(a.filePath) + "_imgedit"
	} else {
		outputFileName = strings.Replace(filepath.Base(a.filePath), a.fileExtension, "_imgedit."+string(extension), 1)
	}
	return path.Join(outputDir, outputFileName), nil
}

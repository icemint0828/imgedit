package app

import (
	"flag"
	"fmt"
	"github.com/golang/freetype/truetype"
	"image"
	"image/color"
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
		isVertical := flagBool(OptionVertical)
		if isVertical {
			c.ReverseY()
		} else {
			c.ReverseX()
		}
	case SubCommandResize.Name:
		ratio := flagFloat64(OptionRatio)
		width, height := int(flagUint(OptionWidth)), int(flagUint(OptionHeight))
		if ratio != 0 {
			c.ResizeRatio(ratio)
		} else {
			c.Resize(width, height)
		}
	case SubCommandTile.Name:
		xLength, yLength := int(flagUint(OptionX)), int(flagUint(OptionY))
		c.Tile(xLength, yLength)
	case SubCommandTrim.Name:
		left, top := int(flagUint(OptionLeft)), int(flagUint(OptionTop))
		width, height := int(flagUint(OptionWidth)), int(flagUint(OptionHeight))
		c.Trim(left, top, width, height)
	case SubCommandGrayscale.Name:
		c.Grayscale()
	case SubCommandAddstring.Name:
		left, top := int(flagUint(OptionLeft)), int(flagUint(OptionTop))
		oTtf, oColor := flagString(OptionTtf), flagString(OptionColor)
		size, text := float64(flagUint(OptionSize)), flagString(OptionText)
		option := &imgedit.StringOptions{
			Point: &image.Point{X: left, Y: top},
			Font:  &imgedit.Font{TrueTypeFont: getTtf(oTtf), Size: size, Color: getColor(oColor)},
		}
		c.AddString(text, option)
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

func getTtf(ttfPath string) *truetype.Font {
	if ttfPath == "" {
		return nil
	}
	ttf, err := imgedit.ReadTtf(ttfPath)
	if err != nil {
		return nil
	}
	return ttf
}

func getColor(colorString string) color.Color {
	if colorString == "" {
		return nil
	}
	switch colorString {
	case "black":
		return color.Black
	case "white":
		return color.White
	case "red":
		return color.RGBA{R: 255, G: 0, B: 0, A: 255}
	case "blue":
		return color.RGBA{R: 0, G: 0, B: 255, A: 255}
	case "green":
		return color.RGBA{R: 0, G: 255, B: 0, A: 255}
	default:
		return nil
	}
}

func flagUint(option Option) uint {
	return flag.Lookup(option.Name()).Value.(flag.Getter).Get().(uint)
}

func flagFloat64(option Option) float64 {
	return flag.Lookup(option.Name()).Value.(flag.Getter).Get().(float64)
}

func flagBool(option Option) bool {
	return flag.Lookup(option.Name()).Value.(flag.Getter).Get().(bool)
}

func flagString(option Option) string {
	return flag.Lookup(option.Name()).Value.(flag.Getter).Get().(string)
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

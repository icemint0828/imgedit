package app

import (
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

type subcommand func(imgedit.FileConverter)

var subcommands = map[string]subcommand{
	"resize":    resize,
	"trim":      trim,
	"tile":      tile,
	"reverse":   reverse,
	"grayscale": grayscale,
	"addstring": addstring,
	"filter":    filter,
}

// Run edit the image
func (a *App) Run() error {
	// load image
	c, extension, err := imgedit.NewFileConverter(a.filePath)
	if err != nil {
		return err
	}
	// convert image
	if subcommand, ok := subcommands[a.subCommand.Name]; ok {
		subcommand(c)
	} else {
		switch a.subCommand.Name {
		case SubCommandPng.Name:
			extension = imgedit.Png
		case SubCommandJpeg.Name:
			extension = imgedit.Jpeg
		case SubCommandGif.Name:
			extension = imgedit.Gif
		}
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

func resize(c imgedit.FileConverter) {
	if OptionRatio.Float64() != 0 {
		c.ResizeRatio(OptionRatio.Float64())
	} else {
		c.Resize(OptionWidth.Int(), OptionHeight.Int())
	}
}

func trim(c imgedit.FileConverter) {
	c.Trim(OptionLeft.Int(), OptionTop.Int(), OptionWidth.Int(), OptionHeight.Int())
}

func reverse(c imgedit.FileConverter) {
	if OptionVertical.Bool() {
		c.ReverseY()
	} else {
		c.ReverseX()
	}
}

func tile(c imgedit.FileConverter) {
	c.Tile(OptionX.Int(), OptionY.Int())
}

func grayscale(c imgedit.FileConverter) {
	c.Grayscale()
}

func filter(c imgedit.FileConverter) {
	c.Filter(getModel(OptionMode.String()))
}

func addstring(c imgedit.FileConverter) {
	option := &imgedit.StringOptions{
		Point: &image.Point{X: OptionLeft.Int(), Y: OptionTop.Int()},
		Font:  &imgedit.Font{TrueTypeFont: getTtf(OptionTtf.String()), Size: OptionSize.Float64(), Color: getColor(OptionColor.String())},
	}
	c.AddString(OptionText.String(), option)
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

func getModel(modeString string) imgedit.FilterModel {
	switch modeString {
	case "gray":
		return imgedit.GrayModel
	case "sepia":
		return imgedit.SepiaModel
	default:
		return nil
	}
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

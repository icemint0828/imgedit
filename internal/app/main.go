package app

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/icemint0828/imgedit"
)

const (
	EnvWd = "WD"
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
	outputPath, displayPath, err := a.getOutputPath(extension)
	if err != nil {
		return err
	}
	err = c.SaveAs(outputPath, extension)
	if err != nil {
		return err
	}
	fmt.Printf("save convert file: %s\n", displayPath)
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
	c.Reverse(!OptionVertical.Bool())
}

func tile(c imgedit.FileConverter) {
	c.Tile(OptionX.Int(), OptionY.Int())
}

func grayscale(c imgedit.FileConverter) {
	c.Filter(imgedit.GrayModel)
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
	// specify by color name
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
	}
	// specify by color code(like #FF0000)
	if string(colorString[0]) != "#" || len(colorString) != 7 {
		return nil
	}
	red, err := getColorBits(colorString[1:3])
	if err != nil {
		return nil
	}
	green, err := getColorBits(colorString[3:5])
	if err != nil {
		return nil
	}
	blue, err := getColorBits(colorString[5:7])
	if err != nil {
		return nil
	}
	return color.RGBA{R: red, G: green, B: blue, A: 255}
}

func getColorBits(colorString string) (uint8, error) {
	v, err := strconv.ParseInt(colorString, 16, 64)
	if err != nil {
		return 0, err
	}
	if v < 0 || 255 < v {
		return 0, errors.New("color string is out of range")
	}
	return uint8(v), nil
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

func (a *App) getOutputPath(extension imgedit.Extension) (string, string, error) {
	// Directory of the host when started by docker
	hostDir := os.Getenv(EnvWd)
	outputDir, err := os.Getwd()
	if err != nil {
		return "", "", err
	}
	var outputFileName string
	if a.fileExtension == "" {
		outputFileName = filepath.Base(a.filePath) + "_imgedit"
	} else {
		outputFileName = strings.Replace(filepath.Base(a.filePath), a.fileExtension, "_imgedit."+string(extension), 1)
	}

	outputPath := path.Join(outputDir, outputFileName)
	displayPath := outputPath
	if hostDir != "" {
		displayPath = path.Join(hostDir, outputFileName)
	}
	return outputPath, displayPath, nil
}

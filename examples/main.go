package main

import (
	"github.com/icemint0828/imgedit"
	"image/png"
	"os"
)

func main() {

	// FileConverter
	fc, _, err := imgedit.NewFileConverter("examples/srcImage.png")
	if err != nil {
		panic(err)
	}
	fc.Grayscale()
	err = fc.SaveAs("examples/dstImage.png", imgedit.Png)
	if err != nil {
		panic(err)
	}

	// Converter
	srcFile, err := os.Open("examples/srcImage.png")
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()
	srcImage, err := png.Decode(srcFile)
	if err != nil {
		panic(err)
	}

	c := imgedit.NewConverter(srcImage)
	c.Resize(500, 500)
	dstImage := c.Convert()

	dstFile, err := os.Create("examples/dstImage.png")
	if err != nil {
		panic(err)
	}
	defer dstFile.Close()
	err = png.Encode(dstFile, dstImage)
	if err != nil {
		panic(err)
	}
}

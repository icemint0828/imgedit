package main

import (
	"image/png"
	"os"

	"github.com/icemint0828/imgedit"
)

func main() {
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

	dstFile, err := os.Create("examples/dstImage.gif")
	if err != nil {
		panic(err)
	}
	defer dstFile.Close()
	err = png.Encode(dstFile, dstImage)
	if err != nil {
		panic(err)
	}
}

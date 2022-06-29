package main

import (
	"bytes"
	"os"

	"github.com/icemint0828/imgedit"
)

func main() {
	// ByteConverter
	srcFile, err := os.Open("examples/srcImage.png")
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()

	bc, _, err := imgedit.NewByteConverter(srcFile)
	bc.ResizeRatio(0.5)

	buffer := bytes.NewBuffer([]byte{})
	_ = bc.WriteAs(buffer, imgedit.Jpeg)

	dstFile, err := os.Create("examples/dstImage.jpeg")
	if err != nil {
		panic(err)
	}
	defer dstFile.Close()
	_, _ = buffer.WriteTo(dstFile)
}

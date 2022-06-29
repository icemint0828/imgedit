package main

import (
	"github.com/icemint0828/imgedit"
)

func main() {
	// FileConverter
	fc, _, err := imgedit.NewFileConverter("examples/srcImage.png")
	if err != nil {
		panic(err)
	}
	fc.Filter(imgedit.GrayModel)
	err = fc.SaveAs("examples/dstImage.png", imgedit.Png)
	if err != nil {
		panic(err)
	}
}

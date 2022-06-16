imgedit
===============

## Overview
Imgedit is a package that performs image processing such as resizing and trimming.  
It,s also work on CLI.

## Feature
* resize
* trim
* reverse
* grayscale

## Supported Extensions
* png
* jpg, jpeg
* gif

## Usage(Package)

``` go
package main

import (
	"image/png"
	"os"

	"github.com/icemint0828/imgedit"
)

func main() {
	srcFile, err := os.Open("srcImage.png")
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

	dstFile, err := os.Create("dstImage.png")
	if err != nil {
		panic(err)
	}
	defer dstFile.Close()
	err = png.Encode(dstFile, dstImage)
	if err != nil {
		panic(err)
	}
}
```

## Usage(CLI)

You can download the executable file from the link below.

- ### [Windows](https://github.com/icemint0828/imgedit/releases/latest/download/imgedit_Windows.zip)

- ### [Linux](https://github.com/icemint0828/imgedit/releases/latest/download/imgedit_Linux.zip)

- ### [macOS](https://github.com/icemint0828/imgedit/releases/latest/download/imgedit_MacOS.zip)


For more information on the executable file, please see the following command

```shell
imgedit -help
```

## License

imgedit is under [MIT license](https://github.com/icemint0828/imgedit/blob/main/LICENSE).
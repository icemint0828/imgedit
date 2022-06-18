
<div align="center" style="width: 100%;">
	<img alt="imgedit" src="assets/image/logo.png" style="width:60%;">
</div>

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/icemint0828/imgedit)
[![Go project version](https://badge.fury.io/go/github.com%2Ficemint0828%2Fimgedit.svg)](https://badge.fury.io/go/github.com%2Ficemint0828%2Fimgedit)
[![Go](https://github.com/icemint0828/imgedit/actions/workflows/go.yml/badge.svg)](https://github.com/icemint0828/imgedit/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/icemint0828/imgedit)](https://goreportcard.com/report/github.com/icemint0828/imgedit)
[![codecov](https://codecov.io/gh/icemint0828/imgedit/branch/main/graph/badge.svg?token=GI2WTY1V5O)](https://codecov.io/gh/icemint0828/imgedit)
[![CodeFactor](https://www.codefactor.io/repository/github/icemint0828/imgedit/badge)](https://www.codefactor.io/repository/github/icemint0828/imgedit)
[![golang example](https://img.shields.io/badge/golang%20example-image%20processing-blue)](https://golangexample.com/imgedit-a-package-that-performs-image-processing-such-as-resizing-and-trimming)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Overview

Imgedit is a package that performs image processing such as resizing and trimming.  
It's also work on CLI.

## Feature

* resize
* trim
* reverse
* grayscale

 <table>
    <tr>
      <td>resize</td>
      <td>trim</td>
      <td>grayscale</td>
    </tr>
    <tr>
      <td><img alt="resize" src="assets/image/resize.gif"></td>
      <td><img alt="trim" src="assets/image/trim.gif"></td>
      <td><img alt="grayscale" src="assets/image/grayscale.gif"> </td>
    </tr>
    <tr>
      <td>reverse horizon </td>
      <td>reverse vertical </td>
      <td></td>
    </tr>
    <tr>
      <td><img alt="reverse-x" src="assets/image/reverse-x.gif"></td>
      <td><img alt="reverse-y" src="assets/image/reverse-y.gif"></td>
      <td></td>
    </tr>
 </table>

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

For more information on the executable file, please see the following command display

```shell
imgedit -help
```

## Supported Extensions(CLI)

* png
* jpg, jpeg
* gif

## Usage(GUI)

You can also use a sample GUI tool that created with wasm by this package.

- ### [image edit tools](https://icemint0828.github.io/imgedit-wasm/)

## License

imgedit is under [MIT license](https://github.com/icemint0828/imgedit/blob/main/LICENSE).

package main

import (
	"bytes"
	"github.com/icemint0828/imgedit"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"syscall/js"
)

// グローバルオブジェクト（Webブラウザはwindow）の取得
var window = js.Global()
var document = window.Get("document")

func main() {
	ch := make(chan interface{})
	window.Set("grayscale", js.FuncOf(grayscale))
	<-ch
}

func grayscale(_ js.Value, _ []js.Value) interface{} {
	fileInput := getElementById("file-input")
	message := getElementById("error-message")
	item := fileInput.Get("files").Call("item", 0)
	if item.IsNull() {
		message.Set("innerHTML", "file not found")
		return nil
	}

	item.Call("arrayBuffer").Call("then", js.FuncOf(func(v js.Value, x []js.Value) any {
		srcData := window.Get("Uint8Array").New(x[0])
		src := make([]byte, srcData.Get("length").Int())
		js.CopyBytesToGo(src, srcData)
		srcImg, fmt, err := image.Decode(bytes.NewBuffer(src))
		if err != nil {
			message.Set("innerHTML", " unsupported file")
			return nil
		}
		bytes.NewBuffer(src)
		c := imgedit.NewConverter(srcImg)
		c.Grayscale()
		dstImg := c.Convert()

		//dst := make([]byte, 1024*1000) // 1M
		dstBuf := &bytes.Buffer{}
		switch fmt {
		case "png":
			png.Encode(dstBuf, dstImg)
		case "jpeg":
			jpeg.Encode(dstBuf, dstImg, &jpeg.Options{Quality: 100})
		case "gif":
			gif.Encode(dstBuf, dstImg, &gif.Options{NumColors: 256})
		}
		var dstData = window.Get("Uint8Array").New(dstBuf.Len())
		js.CopyBytesToJS(dstData, dstBuf.Bytes())
		window.Call("previewBlob", dstData.Get("buffer"))
		return nil
	}))

	return nil
}

func getElementById(id string) js.Value {
	return document.Call("getElementById", id)
}

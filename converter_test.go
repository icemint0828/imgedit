package imgedit

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"reflect"
	"testing"

	"github.com/magiconair/properties/assert"
)

const (
	MissingImagePath  = "assets/image/missingFile"
	WongExtensionPath = "assets/image/wongExtension.txt"
	MissingDirPath    = "assets/image/missingDir/missingFile"
	SrcGifImagePath   = "assets/image/srcImage.gif"
	SrcJpegImagePath  = "assets/image/srcImage.jpeg"
	SrcPngImagePath   = "assets/image/srcImage.png"
	AlphaImagePath    = "assets/image/alphaImage.png"
	DstPngImagePath   = "assets/image/dstImage.png"
	DstJpegImagePath  = "assets/image/dstImage.jpeg"
	DstGifImagePath   = "assets/image/dstImage.gif"
)

func GetPngImage() image.Image {
	p, err := os.Open(SrcPngImagePath)
	if err != nil {
		panic(err)
	}
	defer p.Close()

	img, err := png.Decode(p)
	if err != nil {
		panic(err)
	}
	return img
}

func GetJpegImage() image.Image {
	p, err := os.Open(SrcJpegImagePath)
	if err != nil {
		panic(err)
	}
	defer p.Close()

	img, err := jpeg.Decode(p)
	if err != nil {
		panic(err)
	}
	return img
}

func GetGifImage() image.Image {
	p, err := os.Open(SrcGifImagePath)
	if err != nil {
		panic(err)
	}
	defer p.Close()

	img, err := gif.Decode(p)
	if err != nil {
		panic(err)
	}
	return img
}

func GetAlphaImage() image.Image {
	p, err := os.Open(AlphaImagePath)
	if err != nil {
		panic(err)
	}
	defer p.Close()

	img, err := png.Decode(p)
	if err != nil {
		panic(err)
	}
	return img
}

func SaveTestImage(img image.Image) {
	p, err := os.Create(DstPngImagePath)
	if err != nil {
		panic(err)
	}
	defer p.Close()

	err = png.Encode(p, img)
	if err != nil {
		panic(err)
	}
}

func TestNewConverter(t *testing.T) {
	type args struct {
		image image.Image
	}
	tests := []struct {
		name string
		args args
		want Converter
	}{
		{
			name: "normal",
			args: args{image: GetPngImage()},
			want: &converter{Image: GetPngImage()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConverter(tt.args.image); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConverter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_converter_Resize(t *testing.T) {
	type fields struct {
		Image image.Image
	}
	type args struct {
		resizeX int
		resizeY int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "normal",
			fields: fields{Image: GetPngImage()},
			args:   args{resizeX: 500, resizeY: 500},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &converter{
				Image: tt.fields.Image,
			}
			c.Resize(tt.args.resizeX, tt.args.resizeY)
			img := c.Convert()
			assert.Equal(t, img.Bounds().Dx(), tt.args.resizeX)
			assert.Equal(t, img.Bounds().Dy(), tt.args.resizeX)
			SaveTestImage(img)
		})
	}
}

func Test_converter_ResizeRatio(t *testing.T) {
	type fields struct {
		Image image.Image
	}
	type args struct {
		ratio float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "normal",
			fields: fields{Image: GetPngImage()},
			args:   args{ratio: 0.3},
		},
		{
			name:   "alpha",
			fields: fields{Image: GetAlphaImage()},
			args:   args{ratio: 0.3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &converter{
				Image: tt.fields.Image,
			}
			c.ResizeRatio(tt.args.ratio)
			img := c.Convert()
			assert.Equal(t, img.Bounds().Dx(), int(math.Round(float64(tt.fields.Image.Bounds().Dx())*tt.args.ratio)))
			assert.Equal(t, img.Bounds().Dy(), int(math.Round(float64(tt.fields.Image.Bounds().Dy())*tt.args.ratio)))
			SaveTestImage(img)
		})
	}
}

func Test_converter_Trim(t *testing.T) {
	type fields struct {
		Image image.Image
	}
	type args struct {
		left   int
		top    int
		width  int
		height int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "normal",
			fields: fields{Image: GetPngImage()},
			args:   args{500, 500, 500, 500},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &converter{
				Image: tt.fields.Image,
			}
			c.Trim(tt.args.left, tt.args.top, tt.args.width, tt.args.height)
			img := c.Convert()
			assert.Equal(t, img.Bounds().Dx(), tt.args.width)
			assert.Equal(t, img.Bounds().Dy(), tt.args.height)
			SaveTestImage(img)
		})
	}
}

func Test_converter_ReverseX(t *testing.T) {
	type fields struct {
		Image image.Image
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "normal",
			fields: fields{Image: GetPngImage()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &converter{
				Image: tt.fields.Image,
			}
			c.ReverseX()
			img := c.Convert()
			assert.Equal(t, img.Bounds().Dx(), tt.fields.Image.Bounds().Dx())
			assert.Equal(t, img.Bounds().Dy(), tt.fields.Image.Bounds().Dy())
			SaveTestImage(img)
		})
	}
}

func Test_converter_ReverseY(t *testing.T) {
	type fields struct {
		Image image.Image
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "normal",
			fields: fields{Image: GetPngImage()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &converter{
				Image: tt.fields.Image,
			}
			c.ReverseY()
			img := c.Convert()
			assert.Equal(t, img.Bounds().Dx(), tt.fields.Image.Bounds().Dx())
			assert.Equal(t, img.Bounds().Dy(), tt.fields.Image.Bounds().Dy())
			SaveTestImage(img)
		})
	}
}

func Test_converter_Grayscale(t *testing.T) {
	type fields struct {
		Image image.Image
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "normal",
			fields: fields{Image: GetPngImage()},
		},
		{
			name:   "alpha",
			fields: fields{Image: GetAlphaImage()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &converter{
				Image: tt.fields.Image,
			}
			c.Grayscale()
			img := c.Convert()
			assert.Equal(t, img.Bounds().Dx(), tt.fields.Image.Bounds().Dx())
			assert.Equal(t, img.Bounds().Dy(), tt.fields.Image.Bounds().Dy())
			SaveTestImage(img)
		})
	}
}

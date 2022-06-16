package imgedit

import (
	"image"
	"image/png"
	"os"
	"reflect"
	"testing"

	"github.com/magiconair/properties/assert"
)

const (
	ImagePath  = "assets/image/mint.png"
	OutputPath = "assets/output/test.png"
)

func GetTestImage() image.Image {
	p, err := os.Open(ImagePath)
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
	p, err := os.Create(OutputPath)
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
			args: args{image: GetTestImage()},
			want: &converter{Image: GetTestImage()},
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
			fields: fields{Image: GetTestImage()},
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

func Test_converter_Trim(t *testing.T) {
	type fields struct {
		Image image.Image
	}
	type args struct {
		left   int
		right  int
		bottom int
		top    int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "normal",
			fields: fields{Image: GetTestImage()},
			args:   args{750, 1250, 750, 1250},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &converter{
				Image: tt.fields.Image,
			}
			c.Trim(tt.args.left, tt.args.right, tt.args.bottom, tt.args.top)
			img := c.Convert()
			assert.Equal(t, img.Bounds().Dx(), tt.args.right-tt.args.left)
			assert.Equal(t, img.Bounds().Dy(), tt.args.top-tt.args.bottom)
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
			fields: fields{Image: GetTestImage()},
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
			fields: fields{Image: GetTestImage()},
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
			fields: fields{Image: GetTestImage()},
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

package imgedit

import (
	_ "embed"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"math"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	// DefaultFontSize used when font size is not specified in StringOptions
	DefaultFontSize = 100

	// DefaultOutlineWidth used when outline width is not specified in StringOptions
	DefaultOutlineWidth = 100
)

//go:embed assets/font/07LogoTypeGothic7.ttf
var TtfFile []byte

// Converter interface for image edit
type Converter interface {
	Resize(x, y int)
	ResizeRatio(ratio float64)
	Trim(left, top, width, height int)
	ReverseX()
	ReverseY()
	Filter(filterModel FilterModel)
	// Deprecated: Replace Filter(imgedit.GrayModel).
	Grayscale()
	AddString(text string, options *StringOptions)
	Tile(xLength, yLength int)
	Convert() image.Image
}

type converter struct {
	image.Image
}

// NewConverter create converter
func NewConverter(image image.Image) Converter {
	return &converter{image}
}

// Resize resize the image
func (c *converter) Resize(resizeX, resizeY int) {
	dst := image.NewRGBA(image.Rect(0, 0, resizeX, resizeY))
	dstSize := dst.Bounds().Size()
	xRate, yRate := float64(c.Bounds().Dx())/float64(dstSize.X), float64(c.Bounds().Dy())/float64(dstSize.Y)
	for x := 0; x < dstSize.X; x++ {
		for y := 0; y < dstSize.Y; y++ {
			srcX, srcY := int(math.Round(float64(x)*xRate)), int(math.Round(float64(y)*yRate))
			dst.Set(x, y, c.Image.At(srcX, srcY))
		}
	}
	c.Image = dst
}

// ResizeRatio resize the image with ratio
func (c *converter) ResizeRatio(ratio float64) {
	dst := image.NewRGBA(image.Rect(0, 0, int(math.Round(float64(c.Bounds().Dx())*ratio)), int(math.Round(float64(c.Bounds().Dy())*ratio))))
	dstSize := dst.Bounds().Size()
	xRate, yRate := 1/ratio, 1/ratio
	for x := 0; x < dstSize.X; x++ {
		for y := 0; y < dstSize.Y; y++ {
			srcX, srcY := int(math.Round(float64(x)*xRate)), int(math.Round(float64(y)*yRate))
			dst.Set(x, y, c.Image.At(srcX, srcY))
		}
	}
	c.Image = dst
}

// Trim trim the image to the specified size
func (c *converter) Trim(left, top, width, height int) {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	dstSize := dst.Bounds().Size()
	for x := 0; x < dstSize.X; x++ {
		for y := 0; y < dstSize.Y; y++ {
			srcX, srcY := x+left, y+top
			dst.Set(x, y, c.Image.At(srcX, srcY))
		}
	}
	c.Image = dst
}

// ReverseX reverse the image about horizon
func (c *converter) ReverseX() {
	dst := image.NewRGBA(image.Rect(0, 0, c.Bounds().Dx(), c.Bounds().Dy()))
	srcSize := c.Bounds().Size()
	dstSize := dst.Bounds().Size()
	for x := 0; x < dstSize.X; x++ {
		for y := 0; y < dstSize.Y; y++ {
			srcX, srcY := srcSize.X-x, y
			dst.Set(x, y, c.Image.At(srcX, srcY))
		}
	}
	c.Image = dst
}

// ReverseY reverse the image about vertical
func (c *converter) ReverseY() {
	dst := image.NewRGBA(image.Rect(0, 0, c.Bounds().Dx(), c.Bounds().Dy()))
	srcSize := c.Bounds().Size()
	dstSize := dst.Bounds().Size()
	for x := 0; x < srcSize.X; x++ {
		for y := 0; y < srcSize.Y; y++ {
			srcX, srcY := x, dstSize.Y-y
			dst.Set(x, y, c.Image.At(srcX, srcY))
		}
	}
	c.Image = dst
}

// Grayscale change the image color to grayscale
// Deprecated: Replace Filter(imgedit.GrayModel).
func (c *converter) Grayscale() {
	dst := image.NewRGBA(image.Rect(0, 0, c.Bounds().Dx(), c.Bounds().Dy()))
	dstSize := dst.Bounds().Size()
	grayModel := color.GrayModel
	alphaModel := color.AlphaModel
	for x := 0; x < dstSize.X; x++ {
		for y := 0; y < dstSize.Y; y++ {
			dstColor := c.Image.At(x, y)
			_, _, _, a := dstColor.RGBA()
			if a <= 0 {
				dst.Set(x, y, alphaModel.Convert(dstColor))
			} else {
				dst.Set(x, y, grayModel.Convert(dstColor))
			}
		}
	}
	c.Image = dst
}

// FilterModel wrapper
type FilterModel color.Model

// GrayModel convert image to gray
var GrayModel = FilterModel(color.GrayModel)

// SepiaModel convert image to sepia
var SepiaModel = FilterModel(color.ModelFunc(sepiaModel))

func sepiaModel(c color.Color) color.Color {
	// once converted to GRAY, then to SEPIA, we can get a beautiful conversion.
	grayColor := color.GrayModel.Convert(c)
	r, g, b, a := grayColor.RGBA()

	r = uint32(float64(r) * (float64(240) / float64(255)))
	g = uint32(float64(g) * (float64(200) / float64(255)))
	b = uint32(float64(b) * (float64(148) / float64(255)))
	return color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
}

// Filter change the image color to specified color model
// GrayModel, SepiaModel
func (c *converter) Filter(filterModel FilterModel) {
	if filterModel == nil {
		return
	}

	dst := image.NewRGBA(image.Rect(0, 0, c.Bounds().Dx(), c.Bounds().Dy()))
	dstSize := dst.Bounds().Size()
	alphaModel := color.AlphaModel
	for x := 0; x < dstSize.X; x++ {
		for y := 0; y < dstSize.Y; y++ {
			dstColor := c.Image.At(x, y)
			_, _, _, a := dstColor.RGBA()
			if a <= 0 {
				dst.Set(x, y, alphaModel.Convert(dstColor))
			} else {
				dst.Set(x, y, filterModel.Convert(dstColor))
			}
		}
	}
	c.Image = dst
}

// StringOptions options for AddString
type StringOptions struct {
	// Point left top = (0px, 0px), default center
	Point *image.Point
	// Font
	Font *Font
	// Outline
	Outline *Outline
}

// Font used in the options
type Font struct {
	// TrueTypeFont use ReadTtf to get font
	TrueTypeFont *truetype.Font
	// Size default 100
	Size float64
	// Color default color.Black
	Color color.Color
}

// Outline used in the options
type Outline struct {
	// Color default color.White
	Color color.Color
	// Width from font. 0 <= Width <= 200 recommended
	Width int
}

func (o *StringOptions) setDefault() {
	// font
	if o.Font == nil {
		o.Font = &Font{}
	}
	if o.Font.TrueTypeFont == nil {
		o.Font.TrueTypeFont, _ = ReadTtfFromByte(TtfFile)
	}
	if o.Font.Size == 0 {
		o.Font.Size = DefaultFontSize
	}
	if o.Font.Color == nil {
		o.Font.Color = color.Black
	}

	// outLine
	if o.Outline != nil {
		if o.Outline.Color == nil {
			o.Outline.Color = color.White
		}
		if o.Outline.Width == 0 {
			o.Outline.Width = DefaultOutlineWidth
		}
	}
}

func (o *StringOptions) face() font.Face {
	// use only font size
	return truetype.NewFace(o.Font.TrueTypeFont, &truetype.Options{Size: o.Font.Size})
}

func (o *StringOptions) color() *image.Uniform {
	return image.NewUniform(o.Font.Color)
}

func (o *StringOptions) colorOutLine() *image.Uniform {
	return image.NewUniform(o.Outline.Color)
}

// AddString add string on current Image
func (c *converter) AddString(text string, options *StringOptions) {
	if text == "" {
		return
	}
	if options == nil {
		options = &StringOptions{}
	}
	options.setDefault()

	// copy base image
	dst := image.NewRGBA(image.Rect(0, 0, c.Bounds().Dx(), c.Bounds().Dy()))
	draw.Draw(dst, image.Rect(0, 0, c.Bounds().Dx(), c.Bounds().Dy()), c.Image, image.Point{}, draw.Over)

	var outLinDrawer *font.Drawer
	if options.Outline != nil {
		outLinDrawer = &font.Drawer{
			Dst:  dst,
			Src:  options.colorOutLine(),
			Face: options.face(),
		}
	}

	drawer := &font.Drawer{
		Dst:  dst,
		Src:  options.color(),
		Face: options.face(),
	}
	drawString(dst, drawer, outLinDrawer, text, options)
	c.Image = dst
}

// drawString draw string at adjusted position
func drawString(dst draw.Image, drawer *font.Drawer, outlineDrawer *font.Drawer, text string, options *StringOptions) {
	// notice : Dot values are determined by feeling.
	// set drawer first position
	if options.Point == nil {
		drawer.Dot.X, drawer.Dot.Y = (fixed.I(dst.Bounds().Dx())-drawer.MeasureString(text))/2, (fixed.I(dst.Bounds().Dy())+(drawer.Face.Metrics().Ascent+drawer.Face.Metrics().Descent)/2)/2
	} else {
		drawer.Dot.X, drawer.Dot.Y = fixed.I(options.Point.X)-drawer.MeasureString(text)/2, fixed.I(options.Point.Y)+((drawer.Face.Metrics().Ascent+drawer.Face.Metrics().Descent)/2)/2
	}
	if outlineDrawer == nil {
		drawer.DrawString(text)
		return
	}

	// notice : width values are determined by feeling
	width := drawer.Face.Metrics().Height / 128 / 100 * fixed.Int26_6(options.Outline.Width)

	// draw letter by letter for adjust outline position
	for _, s := range []byte(text) {
		// As a technique, the image is drawn on the four corners,
		// but that is an incomplete implementation of OUTLINE.

		// left top
		outlineDrawer.Dot = fixed.Point26_6{X: drawer.Dot.X - width, Y: drawer.Dot.Y - width}
		outlineDrawer.DrawBytes([]byte{s})

		// left bottom
		outlineDrawer.Dot = fixed.Point26_6{X: drawer.Dot.X - width, Y: drawer.Dot.Y + width}
		outlineDrawer.DrawBytes([]byte{s})

		// right top
		outlineDrawer.Dot = fixed.Point26_6{X: drawer.Dot.X + width, Y: drawer.Dot.Y - width}
		outlineDrawer.DrawBytes([]byte{s})

		// right bottom
		outlineDrawer.Dot = fixed.Point26_6{X: drawer.Dot.X + width, Y: drawer.Dot.Y + width}
		outlineDrawer.DrawBytes([]byte{s})

		drawer.DrawBytes([]byte{s})
	}
}

func (c *converter) Tile(xLength, yLength int) {
	dst := image.NewRGBA(image.Rect(0, 0, c.Bounds().Dx()*xLength, c.Bounds().Dy()*yLength))
	srcSize := c.Bounds().Size()
	for xLen := 0; xLen < xLength; xLen++ {
		for yLen := 0; yLen < yLength; yLen++ {
			for x := 0; x < srcSize.X; x++ {
				for y := 0; y < srcSize.Y; y++ {
					dst.Set(x+srcSize.X*xLen, y+srcSize.Y*yLen, c.Image.At(x, y))
				}
			}
		}
	}
	c.Image = dst
}

// Convert get convert image
func (c *converter) Convert() image.Image {
	return c.Image
}

// ReadTtf return ttf from file path
func ReadTtf(ttfFilePath string) (*truetype.Font, error) {
	fontFile, err := ioutil.ReadFile(ttfFilePath)
	if err != nil {
		return nil, err
	}
	ttf, err := freetype.ParseFont(fontFile)
	if err != nil {
		return nil, err
	}
	return ttf, nil
}

// ReadTtfFromByte return ttf from file byte
func ReadTtfFromByte(ttfFile []byte) (*truetype.Font, error) {
	ttf, err := freetype.ParseFont(ttfFile)
	if err != nil {
		return nil, err
	}
	return ttf, nil
}

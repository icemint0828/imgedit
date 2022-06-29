package imgedit

import (
	"reflect"
	"testing"
)

func TestNewFileConverter(t *testing.T) {
	type args struct {
		srcPath string
	}
	tests := []struct {
		name          string
		args          args
		want          FileConverter
		wantExtension Extension
		wantErr       bool
	}{
		{
			name:          "normal",
			args:          args{srcPath: SrcPngImagePath},
			want:          &fileConverter{byteConverter: &byteConverter{&converter{Image: GetPngImage()}}},
			wantExtension: Png,
			wantErr:       false,
		},
		{
			name:          "missing file",
			args:          args{srcPath: MissingImagePath},
			want:          nil,
			wantExtension: "",
			wantErr:       true,
		},
		{
			name:          "wong extension",
			args:          args{srcPath: WrongExtensionPath},
			want:          nil,
			wantExtension: "",
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, extension, err := NewFileConverter(tt.args.srcPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFileConverter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFileConverter() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(extension, tt.wantExtension) {
				t.Errorf("NewFileConverter() extension = %v, wantExtension %v", got, tt.want)
			}
		})
	}
}

func Test_fileConverter_SaveAs(t *testing.T) {
	type fields struct {
		byteConverter *byteConverter
	}
	type args struct {
		dstPath   string
		extension Extension
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "missing directory",
			fields:  fields{byteConverter: &byteConverter{&converter{Image: GetPngImage()}}},
			args:    args{dstPath: MissingDirPath, extension: Png},
			wantErr: true,
		},
		{
			name:    "unsupported extension",
			fields:  fields{byteConverter: &byteConverter{&converter{Image: GetPngImage()}}},
			args:    args{dstPath: DstPngImagePath, extension: Extension("unsupported")},
			wantErr: true,
		},
		{
			name:    "png to png",
			fields:  fields{byteConverter: &byteConverter{&converter{Image: GetPngImage()}}},
			args:    args{dstPath: DstPngImagePath, extension: Png},
			wantErr: false,
		},
		{
			name:    "jpeg to png",
			fields:  fields{byteConverter: &byteConverter{&converter{Image: GetJpegImage()}}},
			args:    args{dstPath: DstPngImagePath, extension: Png},
			wantErr: false,
		},
		{
			name:    "gif to png",
			fields:  fields{byteConverter: &byteConverter{&converter{Image: GetGifImage()}}},
			args:    args{dstPath: DstPngImagePath, extension: Png},
			wantErr: false,
		},
		{
			name:    "png to jpeg",
			fields:  fields{byteConverter: &byteConverter{&converter{Image: GetPngImage()}}},
			args:    args{dstPath: DstJpegImagePath, extension: Jpeg},
			wantErr: false,
		},
		{
			name:    "jpeg to jpeg",
			fields:  fields{byteConverter: &byteConverter{&converter{Image: GetJpegImage()}}},
			args:    args{dstPath: DstJpegImagePath, extension: Jpeg},
			wantErr: false,
		},
		{
			name:    "gif to jpeg",
			fields:  fields{byteConverter: &byteConverter{&converter{Image: GetGifImage()}}},
			args:    args{dstPath: DstJpegImagePath, extension: Jpeg},
			wantErr: false,
		},
		{
			name:    "png to gif",
			fields:  fields{byteConverter: &byteConverter{&converter{Image: GetPngImage()}}},
			args:    args{dstPath: DstGifImagePath, extension: Gif},
			wantErr: false,
		},
		{
			name:    "jpeg to gif",
			fields:  fields{byteConverter: &byteConverter{&converter{Image: GetJpegImage()}}},
			args:    args{dstPath: DstGifImagePath, extension: Gif},
			wantErr: false,
		},
		{
			name:    "gif to gif",
			fields:  fields{byteConverter: &byteConverter{&converter{Image: GetGifImage()}}},
			args:    args{dstPath: DstGifImagePath, extension: Gif},
			wantErr: false,
		},
		{
			name:    "alpha png to gif",
			fields:  fields{byteConverter: &byteConverter{&converter{Image: GetAlphaPngImage()}}},
			args:    args{dstPath: DstGifImagePath, extension: Gif},
			wantErr: false,
		},
		{
			name:    "alpha gif to gif",
			fields:  fields{byteConverter: &byteConverter{&converter{Image: GetAlphaGifImage()}}},
			args:    args{dstPath: DstGifImagePath, extension: Gif},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &fileConverter{
				byteConverter: tt.fields.byteConverter,
			}
			if err := p.SaveAs(tt.args.dstPath, tt.args.extension); (err != nil) != tt.wantErr {
				t.Errorf("SaveAs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileEdit(t *testing.T) {
	c, _, _ := NewFileConverter(SrcPngImagePath)
	c.Grayscale()
	_ = c.SaveAs(DstPngImagePath, Png)
}

package imgedit

import (
	"bytes"
	"image"
	"image/gif"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestNewByteConverter(t *testing.T) {
	type args struct {
		rc io.ReadCloser
	}
	tests := []struct {
		name    string
		args    args
		want    ByteConverter
		want1   Extension
		wantErr bool
	}{
		{
			name:    "normal",
			args:    args{rc: GetPngImageReadCloser()},
			want:    &byteConverter{&converter{Image: GetPngImage()}},
			want1:   Extension("png"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.args.rc.Close()
			got, got1, err := NewByteConverter(tt.args.rc)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewByteConverter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewByteConverter() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("NewByteConverter() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_byteConverter_WriteAs(t *testing.T) {
	type fields struct {
		converter *converter
	}
	type args struct {
		extension Extension
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantWriter string
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &byteConverter{
				converter: tt.fields.converter,
			}
			writer := &bytes.Buffer{}
			err := b.WriteAs(writer, tt.args.extension)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteAs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("WriteAs() gotWriter = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestSupportedExtension(t *testing.T) {
	type args struct {
		extension Extension
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "normal",
			args: args{extension: Png},
			want: true,
		},
		{
			name: "unsupported extension",
			args: args{extension: Extension("unsupported")},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SupportedExtension(tt.args.extension); got != tt.want {
				t.Errorf("SupportedExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_gifEncode(t *testing.T) {
	type args struct {
		m image.Image
		o *gif.Options
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// fatal error: runtime: out of memory on GitHub
		//{
		//	name:    "out of bounds size",
		//	args:    args{m: image.NewRGBA(image.Rect(0, 0, 1<<16+1, 1<<16+1)), o: nil},
		//	wantErr: true,
		//},
		{
			name:    "over opts.NumColors",
			args:    args{m: image.NewRGBA(image.Rect(0, 0, 100, 100)), o: &gif.Options{NumColors: 257}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			err := gifEncode(w, tt.args.m, tt.args.o)
			if (err != nil) != tt.wantErr {
				t.Errorf("gifEncode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func GetPngImageReadCloser() io.ReadCloser {
	p, err := os.Open(SrcPngImagePath)
	if err != nil {
		panic(err)
	}
	return p
}

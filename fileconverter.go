package imgedit

import (
	"os"
)

// FileConverter interface for image edit
type FileConverter interface {
	Converter
	SaveAs(string, Extension) error
}

type fileConverter struct {
	*byteConverter
}

// NewFileConverter create fileConverter
func NewFileConverter(srcPath string) (FileConverter, Extension, error) {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return nil, "", err
	}
	bc, extension, err := newByteConverter(srcFile)
	if err != nil {
		return nil, "", err
	}
	return &fileConverter{byteConverter: bc}, extension, nil
}

func (p *fileConverter) SaveAs(dstPath string, extension Extension) error {
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	return p.byteConverter.WriteAs(dstFile, extension)
}

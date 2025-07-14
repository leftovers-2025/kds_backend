package entity

import (
	"errors"
	"mime/multipart"
)

type ImageType int

const (
	TYPE_IMAGE_NAME ImageType = 1
	TYPE_IMAGE_FILE ImageType = 2
)

var (
	ErrImageNameRequired = errors.New("image name is required")
	ErrImageFileRequired = errors.New("image file is required")
	ErrImageNoFiles      = errors.New("the image has no files")
)

type Image struct {
	imageType ImageType
	name      string
	file      *multipart.FileHeader
}

func NewNameImage(name string) (*Image, error) {
	if name == "" {
		return nil, ErrImageNameRequired
	}
	return &Image{
		imageType: TYPE_IMAGE_NAME,
		name:      name,
	}, nil
}

func NewFileImage(file *multipart.FileHeader) (*Image, error) {
	if file == nil {
		return nil, ErrImageFileRequired
	}
	return &Image{
		imageType: TYPE_IMAGE_FILE,
		file:      file,
	}, nil
}

func (i *Image) isFile() bool {
	return i.imageType == TYPE_IMAGE_FILE
}

func (i *Image) isName() bool {
	return i.imageType == TYPE_IMAGE_NAME
}

// 名前を取得する。(名前、名前イメージか)
func (i *Image) Name() (string, bool) {
	if i.isName() {
		return i.name, true
	} else {
		return i.file.Filename, false
	}
}

// ファイルを取得する
func (i *Image) File() (*multipart.FileHeader, error) {
	if !i.isFile() {
		return nil, ErrImageNoFiles
	}
	return i.file, nil
}

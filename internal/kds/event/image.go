package event

import (
	"errors"
	"mime/multipart"
)

var (
	ErrImageNameRequired = errors.New("image name is required")
	ErrImageFileRequired = errors.New("image file is required")
)

type ImageUploadEvent struct {
	images []Image
}

type Image struct {
	imageName string
	imageFile *multipart.FileHeader
}

type UploadImage struct {
	Name string `json:"image"`
	File []byte `json:"file"`
}

func NewImage(name string, file *multipart.FileHeader) (*Image, error) {
	if name == "" {
		return nil, ErrImageNameRequired
	}
	if file == nil {
		return nil, ErrImageFileRequired
	}
	return &Image{
		imageName: name,
		imageFile: file,
	}, nil
}

func NewImageUploadEvent(images []Image) (Event, error) {
	if images == nil {
		return nil, ErrImageNameRequired
	}
	return &ImageUploadEvent{
		images: images,
	}, nil
}
func (e *ImageUploadEvent) Event() (any, error) {
	// 	// ファイルを開く
	// 	file, err := e.images.Open()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	defer file.Close()
	//
	// 	// ファイル中身読み込み
	// 	content, err := io.ReadAll(file)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	// 圧縮
	// 	var b bytes.Buffer
	// 	gzWriter := gzip.NewWriter(&b)
	// 	_, err = gzWriter.Write(content)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	err = gzWriter.Close()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	compressedData := b.Bytes()
	// 	return UploadImage{
	// 		Name: e.imageName,
	// 		File: compressedData,
	// 	}, nil
	return nil, nil
}

func (e *ImageUploadEvent) Type() EventType {
	return EVENT_UNKNOWN
}

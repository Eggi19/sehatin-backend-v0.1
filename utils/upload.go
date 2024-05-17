package utils

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type FileUploader interface {
	UploadFile(ctx context.Context, file interface{}) (string, error)
}

type CloudinaryUploadFile struct {
}

func NewCloudinaryUploadFile() *CloudinaryUploadFile {
	return &CloudinaryUploadFile{}
}

func (c *CloudinaryUploadFile) UploadFile(ctx context.Context, file interface{}) (string, error) {
	config, err := ConfigInit()
	if err != nil {
		return "", err
	}

	cld, _ := cloudinary.NewFromParams(config.CloudinaryName, config.CloudinaryKey, config.CloudinarySecret)
	res, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{})
	if err != nil {
		return "", err
	}

	return res.SecureURL, nil
}

package utils

import (
	"context"
	"mime/multipart"
	"strings"

	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
)

func GetFileUrl(ctx context.Context, file *multipart.FileHeader, fileFormat string) (*string, error) {

	if fileFormat == "png" {
		if file.Size > 500000 {
			return nil, custom_errors.FileTooLarge()
		}
		if strings.Split(file.Filename, ".")[1] != "png" {
			return nil, custom_errors.FileNotPng()
		}
	}

	if fileFormat == "pdf" {
		if file.Size > 1000000 {
			return nil, custom_errors.FileTooLarge()
		}
		if strings.Split(file.Filename, ".")[1] != "pdf" {
			return nil, custom_errors.FileNotPdf()
		}
	}

	openFile, _ := file.Open()
	fileUrl, err := NewCloudinaryUploadFile().UploadFile(ctx, openFile)
	if err != nil {
		return nil, custom_errors.UploadFile()
	}

	return &fileUrl, nil
}

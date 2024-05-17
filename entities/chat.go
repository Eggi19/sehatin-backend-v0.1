package entities

import (
	"mime/multipart"
	"time"
)

type Chat struct {
	Id             int64
	IsFromUser     bool
	Content        string
	Type           string
	ConsultationId int64
	CreatedAt      time.Time
	File           *multipart.FileHeader
}

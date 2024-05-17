package dtos

import (
	"mime/multipart"
	"time"

	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
)

type ChatRequest struct {
	Content string `json:"content" binding:"required"`
	Type    string `json:"type" binding:"required"`
}

type ChatFileRequest struct {
	File *multipart.FileHeader `form:"file"`
}

type ChatResponse struct {
	Id         int64     `json:"id"`
	IsFromUser bool      `json:"is_from_user"`
	Content    string    `json:"content"`
	Type       string    `json:"type"`
	CreatedAt  time.Time `json:"created_at"`
}

func ConvertToChatResponse(chat entities.Chat) ChatResponse {
	return ChatResponse{
		Id:         chat.Id,
		IsFromUser: chat.IsFromUser,
		Content:    chat.Content,
		Type:       chat.Type,
		CreatedAt:  chat.CreatedAt,
	}
}

func ConvertToChatResponses(chats []entities.Chat) []ChatResponse {
	chatsResponses := []ChatResponse{}

	for _, c := range chats {
		chatsResponses = append(chatsResponses, ConvertToChatResponse(c))
	}

	return chatsResponses
}

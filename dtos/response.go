package dtos

type ResponseMessage struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

package dtos

type ValidationErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrResponse struct {
	Message string               `json:"message"`
	Details []ValidationErrorMsg `json:"details,omitempty"`
}

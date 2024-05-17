package dtos

import "github.com/tsanaativa/sehatin-backend-v0.1/entities"

type PaginationResponse struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
}

func ConvertToPaginationResponse(pagination entities.PaginationInfo) *PaginationResponse {
	return &PaginationResponse{
		Page:      pagination.Page,
		Limit:     pagination.Limit,
		TotalData: pagination.TotalData,
		TotalPage: pagination.TotalPage,
	}
}

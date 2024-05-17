package dtos

import "github.com/tsanaativa/sehatin-backend-v0.1/entities"

type CategoryResponse struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type CategoryResponses struct {
	Pagination PaginationResponse `json:"pagination_info"`
	Categories []CategoryResponse `json:"categories"`
}

type CategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

func ConvertToCategoryResponse(category entities.Category) *CategoryResponse {
	return &CategoryResponse{
		Id:   category.Id,
		Name: category.Name,
	}
}

func ConvertToCategoryResponses(categories []entities.Category, pagination entities.PaginationInfo) *CategoryResponses {
	categoryResponses := []CategoryResponse{}

	for _, category := range categories {
		categoryResponses = append(categoryResponses, *ConvertToCategoryResponse(category))
	}

	return &CategoryResponses{
		Pagination: *ConvertToPaginationResponse(pagination),
		Categories: categoryResponses,
	}
}

func ConvertToCategoryResponsesWithoutPagination(categories []entities.Category) []CategoryResponse {
	categoryResponses := []CategoryResponse{}

	for _, category := range categories {
		categoryResponses = append(categoryResponses, *ConvertToCategoryResponse(category))
	}

	return categoryResponses
}

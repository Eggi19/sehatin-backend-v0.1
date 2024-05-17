package dtos

import (
	"time"

	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
)

type SalesReportCategoryResponse struct {
	CategoryId   int64  `json:"category_id"`
	CategoryName string `json:"category_name"`
	TotalSold    int    `json:"total_sold"`
	Month        string `json:"month"`
	Year         string `json:"year"`
}

type SalesReportCategoryResponses struct {
	Pagination            PaginationResponse            `json:"pagination_info"`
	SalesReportCategories []SalesReportCategoryResponse `json:"category_reports"`
}

func ConvertSalesReportCategoryResponse(sc *entities.SalesReportCategory) *SalesReportCategoryResponse {
	return &SalesReportCategoryResponse{
		CategoryId:   sc.Category.Id,
		CategoryName: sc.Category.Name,
		TotalSold:    sc.TotalSold,
		Month:        time.Month(sc.Month).String(),
		Year:         sc.Year,
	}
}

func ConvertSalesReportCategoryResponses(scs []entities.SalesReportCategory, pagination entities.PaginationInfo) *SalesReportCategoryResponses {
	scResponeses := []SalesReportCategoryResponse{}

	for _, sc := range scs {
		scResponeses = append(scResponeses, *ConvertSalesReportCategoryResponse(&sc))
	}

	return &SalesReportCategoryResponses{
		Pagination:            *ConvertToPaginationResponse(pagination),
		SalesReportCategories: scResponeses,
	}
}

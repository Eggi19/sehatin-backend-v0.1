package dtos

import (
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/shopspring/decimal"
)

type SalesReportResponse struct {
	Products          ProductCategoryResponse `json:"product_detail"`
	PharmacyId        int64                   `json:"pharmacy_id"`
	PharmacyName      string                  `json:"pharmacy_name"`
	TotalSalesAmount  decimal.Decimal         `json:"total_sales_amount"`
	TotalQuantitySold int                     `json:"total_quantity_sold"`
	Month             string                  `json:"month"`
	Year              int                     `json:"year"`
}

type SalesReportResponses struct {
	Pagination   PaginationResponse    `json:"pagination_info"`
	SalesReports []SalesReportResponse `json:"sales_reports"`
}

func ConvertToSalesResponse(sr *entities.SalesReport) *SalesReportResponse {
	return &SalesReportResponse{
		Products:          *ConvertToProductResponse(sr.PharmacyProduct.Product),
		PharmacyId:        sr.PharmacyProduct.Pharmacy.Id,
		PharmacyName:      sr.PharmacyProduct.Pharmacy.Name,
		TotalSalesAmount:  sr.TotalSalesAmount,
		TotalQuantitySold: sr.TotalQuantitySold,
		Month:             sr.Month.Month().String(),
		Year:              sr.Year.Year(),
	}
}

func ConvertToSalesResponses(srs []entities.SalesReport, pagination entities.PaginationInfo) *SalesReportResponses {
	srResponses := []SalesReportResponse{}

	for _, sr := range srs {
		srResponses = append(srResponses, *ConvertToSalesResponse(&sr))
	}

	return &SalesReportResponses{
		Pagination:   *ConvertToPaginationResponse(pagination),
		SalesReports: srResponses,
	}
}

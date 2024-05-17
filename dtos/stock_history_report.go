package dtos

import (
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
)

type StockHistoryReportResponse struct {
	TotalAddition  int    `json:"total_addition"`
	TotalDeduction int    `json:"total_deduction"`
	FinalStock     int    `json:"final_stock"`
	ProductId      int64  `json:"product_id"`
	ProductName    string `json:"product_name"`
	PharmacyId     int64  `json:"pharmacy_id"`
	PharmacyName   string `json:"pharmacy_name"`
	Month          string `json:"month"`
	Year           int    `json:"year"`
}

type StockHistoryReportResponses struct {
	Pagination          PaginationResponse           `json:"pagination_info"`
	StockHistoryReports []StockHistoryReportResponse `json:"stock_history_reports"`
}

func ConvertToStockHistoryReportResponse(shr *entities.StockHistoryReport) *StockHistoryReportResponse {
	return &StockHistoryReportResponse{
		TotalAddition:  shr.TotalAddition,
		TotalDeduction: shr.TotalDeduction,
		FinalStock:     shr.PharmacyProduct.TotalStock,
		ProductId:      shr.PharmacyProduct.Product.Id,
		ProductName:    shr.PharmacyProduct.Product.Name,
		PharmacyId:     shr.PharmacyProduct.Pharmacy.Id,
		PharmacyName:   shr.PharmacyProduct.Pharmacy.Name,
		Month:          shr.Month.Month().String(),
		Year:           shr.Year.Year(),
	}
}

func ConvertToStockHistoryReportResponses(shr []entities.StockHistoryReport, pagination entities.PaginationInfo) *StockHistoryReportResponses {
	shrResponses := []StockHistoryReportResponse{}

	for _, sh := range shr {
		shrResponses = append(shrResponses, *ConvertToStockHistoryReportResponse(&sh))
	}

	return &StockHistoryReportResponses{Pagination: *ConvertToPaginationResponse(pagination), StockHistoryReports: shrResponses}
}

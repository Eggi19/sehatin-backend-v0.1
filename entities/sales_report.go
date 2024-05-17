package entities

import (
	"time"

	"github.com/shopspring/decimal"
)

type SalesReport struct {
	PharmacyProduct   PharmacyProduct
	TotalSalesAmount  decimal.Decimal
	TotalQuantitySold int
	Month             time.Time
	Year              time.Time
}

type SalesReportParams struct {
	SortBy     string
	Sort       string
	Limit      int
	Page       int
	Keyword    string
	PharmacyId int64
	ProductId  int64
}

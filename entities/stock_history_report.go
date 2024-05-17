package entities

import "time"

type StockHistoryReport struct {
	TotalAddition   int
	TotalDeduction  int
	FinalStock      int
	PharmacyProduct PharmacyProduct
	Month           time.Time
	Year            time.Time
}

type StockHistoryReportParams struct {
	SortBy     string
	Sort       string
	Limit      int
	Page       int
	Keyword    string
	PharmacyId int64
}

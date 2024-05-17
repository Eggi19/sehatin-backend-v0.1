package entities

import (
	"database/sql"
	"time"
)

type StockHistory struct {
	Id              int64
	PharmacyProduct PharmacyProduct
	Pharmacy        Pharmacy
	Quantity        int
	Description     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       sql.NullTime
}

type StockHistoryParams struct {
	SortBy  string
	Sort    string
	Limit   int
	Page    int
	Keyword string
}

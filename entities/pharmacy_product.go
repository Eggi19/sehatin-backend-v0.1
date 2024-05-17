package entities

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type PharmacyProduct struct {
	Id          int64
	Price       decimal.Decimal
	TotalStock  int
	IsAvailable bool
	Product     Product
	Pharmacy    Pharmacy
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   sql.NullTime
}

type PharmacyProductDetailParams struct {
	PharmacyProductId int64
	Coordinat         PharmacyByProductParams
}

type PharmacyProductParams struct {
	SortBy  string
	Sort    string
	Limit   int
	Page    int
	Keyword string
}

type NearestPharmacyProductsParams struct {
	Longitude  float64
	Latitude   float64
	Radius     int
	CategoryId int
	SortBy     string
	Sort       string
	Limit      int
	Page       int
	Keyword    string
}

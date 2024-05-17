package entities

import (
	"database/sql"
	"time"
)

type StockTransfer struct {
	Id               int64
	PharmacySender   Pharmacy
	PharmacyReceiver Pharmacy
	MutationStatus   MutationSatus
	Product          Product
	Quantity         int
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        sql.NullTime
}

type StockTransferParams struct {
	SortBy string
	Sort   string
	Limit  int
	Page   int
}

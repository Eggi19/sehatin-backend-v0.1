package entities

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type CartItem struct {
	Id                int64
	Quantity          int
	UserId            int64
	PharmacyProductId int64
	ProductName       string
	ProductPicture    string
	SellingUnit       string
	Price             decimal.Decimal
	PharmacyId        int64
	PharmacyName      string
	SlugId            string
	TotalStock        int
	Weight            int
	IsAvailable       bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         sql.NullTime
}

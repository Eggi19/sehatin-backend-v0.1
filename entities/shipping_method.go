package entities

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type OfficialShippingMethod struct {
	Id        int64
	Name      string
	Fee       decimal.Decimal
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

type NonOfficialShippingMethod struct {
	Id          int64
	Name        string
	Courier     string
	Service     string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   sql.NullTime
}

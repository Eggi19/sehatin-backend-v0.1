package entities

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	Id              int64
	OrderNumber     string
	TotalPrice      decimal.Decimal
	PaymentProof    *string
	PaymentDeadline time.Time
	ShippingFee     decimal.Decimal
	ShippingMethod  string
	UserAddressId   int64
	UserAddress     UserAddress
	PharmacyAddress PharmacyAddress
	OrderStatusId   int64
	OrderStatus     string
	PharmacyId      int64
	PharmacyName    string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       sql.NullTime
	Total           int
	UserName        string
	UserEmail       string
}

type OrderParams struct {
	Limit      int
	Page       int
	Status     string
	PharmacyId int
}

type UpdateOrderStatus struct {
	OrderId           int64
	OrderStatus       string
	PharmacyManagerId int64
	UserId            int64
}

type UploadPaymentProof struct {
	PaymentProof string
	OrderId      int64
	UserId       int64
}

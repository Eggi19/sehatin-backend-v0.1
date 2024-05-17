package entities

import (
	"database/sql"
	"time"
)

type Pharmacy struct {
	Id                        int64
	PharmacyManager           PharmacyManager
	Name                      string
	OperationalHour           string
	OperationalDay            string
	PharmacistName            string
	PharmacistLicenseNumber   string
	PharmacistPhoneNumber     string
	PharmacyAddress           PharmacyAddress
	OfficialShippingMethod    []OfficialShippingMethod
	NonOfficialShippingMethod []NonOfficialShippingMethod
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
	DeletedAt                 sql.NullTime
}

type PharmacyParams struct {
	SortBy  string
	Sort    string
	Limit   int
	Page    int
	Keyword string
}

type NearestPharmacyParams struct {
	Longitude  float64
	Latitude   float64
	Radius     int
	CategoryId int
}

type PharmacyByProductParams struct {
	Longitude float64
	Latitude  float64
	Radius    int
	ProductId int64
}

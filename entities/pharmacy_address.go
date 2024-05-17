package entities

import (
	"database/sql"
	"time"
)

type PharmacyAddress struct {
	Id          int64
	PharmacyId  int64
	City        string
	CityId      int
	Province    string
	Address     string
	District    string
	SubDistrict string
	PostalCode  string
	Coordinate  string
	Longitude   string
	Latitude    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   sql.NullTime
}

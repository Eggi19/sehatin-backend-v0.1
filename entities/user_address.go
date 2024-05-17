package entities

import (
	"database/sql"
	"time"
)

type UserAddress struct {
	Id          int64
	UserId      int64
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
	IsMain      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   sql.NullTime
}

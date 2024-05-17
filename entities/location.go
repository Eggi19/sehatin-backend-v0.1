package entities

import (
	"database/sql"
	"time"
)

type Province struct {
	Id        int16
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

type City struct {
	Id         int16
	Name       string
	Type       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  sql.NullTime
	ProvinceId int16
}

type District struct {
	Id        int16
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
	CityId    int16
}

type SubDistrict struct {
	Id         int32
	Name       string
	PostalCode int32
	Coordinate string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  sql.NullTime
	DistrictId int16
}

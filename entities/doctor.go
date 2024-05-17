package entities

import (
	"database/sql"
	"time"
)

type Doctor struct {
	Id             int64
	Name           string
	Email          string
	Password       sql.NullString
	IsVerified     bool
	IsGoogle       bool
	IsOnline       bool
	Fee            sql.NullInt64
	Certificate    sql.NullString
	WorkStartYear  sql.NullInt64
	Specialist     *DoctorSpecialist
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      sql.NullTime
	ProfilePicture sql.NullString
}

type DoctorSpecialist struct {
	Id        sql.NullInt64
	Name      sql.NullString
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

type DoctorParams struct {
	SpecialistId int64
	IsOnline     string
	SortBy       string
	Sort         string
	Limit        int
	Page         int
	Keyword      string
}

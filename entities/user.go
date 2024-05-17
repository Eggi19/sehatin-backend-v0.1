package entities

import (
	"database/sql"
	"time"
)

type User struct {
	Id             int64
	Name           string
	Email          string
	Password       sql.NullString
	BirthDate      sql.NullString
	Gender         *Gender
	IsVerified     bool
	IsGoogle       bool
	IsOnline       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      sql.NullTime
	ProfilePicture sql.NullString
	Address        []UserAddress
}

type UserParams struct {
	GenderId   int64
	IsVerified bool
	IsGoogle   bool
	IsOnline   bool
	SortBy     string
	Sort       string
	Limit      int
	Page       int
	Keyword    string
}

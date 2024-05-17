package entities

import (
	"database/sql"
	"time"
)

type PharmacyManager struct {
	Id          int64
	Name        string
	Email       string
	Password    string
	PhoneNumber string
	Logo        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   sql.NullTime
}

type PharmacyManagerParams struct {
	SortBy  string
	Sort    string
	Limit   int
	Page    int
	Keyword string
}

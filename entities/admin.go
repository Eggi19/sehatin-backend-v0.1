package entities

import (
	"database/sql"
	"time"
)

type Admin struct {
	Id        int64
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

type AdminParams struct {
	SortBy  string
	Sort    string
	Limit   int
	Page    int
	Keyword string
}

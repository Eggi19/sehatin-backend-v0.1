package entities

import (
	"database/sql"
	"time"
)

type Category struct {
	Id        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

type CategoryParams struct {
	SortBy  string
	Sort    string
	Limit   int
	Page    int
	Keyword string
}

package entities

import (
	"database/sql"
	"time"
)

type MutationSatus struct {
	Id        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

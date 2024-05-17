package entities

import (
	"database/sql"
	"time"
)

type PrescriptionItem struct {
	Id        int64
	Product   Product
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

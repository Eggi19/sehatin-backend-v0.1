package entities

import (
	"database/sql"
	"time"
)

type UserResetPassword struct {
	Id        int64
	UserId    int64
	Token     string
	ExpiredAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

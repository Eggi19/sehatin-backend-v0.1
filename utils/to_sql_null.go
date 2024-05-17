package utils

import "database/sql"

func ByteToNullString(value []byte) *sql.NullString {
	valueStr := string(value)
	return &sql.NullString{String: valueStr, Valid: true}
}

func StringToNullString(value string) *sql.NullString {
	return &sql.NullString{String: value, Valid: true}
}

func Int64ToNullInt64(value int64) *sql.NullInt64 {
	return &sql.NullInt64{
		Int64: value,
		Valid: true,
	}
}

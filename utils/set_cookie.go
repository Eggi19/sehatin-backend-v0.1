package utils

import (
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Expired struct {
	AccessTokenExp        string
	RefreshTokenExp       string
	ResetPasswordTokenExp time.Time
}

func SetExpire() *Expired {
	env, _ := godotenv.Read()

	accessExpHourInt, _ := strconv.Atoi(env["EXP_HOUR"])
	refreshExpHourInt, _ := strconv.Atoi(env["REFRESH_EXP_HOUR"])
	resetPasswordExpHourInt, _ := strconv.Atoi(env["RESET_PASSWORD_TOKEN_EXP"])

	accessExpHour := time.Now().Add(time.Duration(accessExpHourInt) * time.Hour).Format("2006-01-02T15:04:05Z07:00")
	refreshExpHour := time.Now().Add(time.Duration(refreshExpHourInt) * time.Hour).Format("2006-01-02T15:04:05Z07:00")
	resetPasswordExpHour := time.Now().Add(time.Duration(resetPasswordExpHourInt) * time.Minute)

	return &Expired{
		AccessTokenExp:        accessExpHour,
		RefreshTokenExp:       refreshExpHour,
		ResetPasswordTokenExp: resetPasswordExpHour,
	}
}

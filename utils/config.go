package utils

import (
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DbUrl                 string
	Port                  string
	HashCost              int
	Issuer                string
	ExpDurationHour       int
	RefreshExpDuration    int
	SecretKey             string
	Email                 string
	EmailPassword         string
	SmtpHost              string
	SmtpPort              string
	FrontendUrl           string
	CloudinaryName        string
	CloudinaryKey         string
	CloudinarySecret      string
	GoogleId              string
	GoogleKey             string
	GoogleUri             string
	ResetTokenExpDuration int
	RajaOngkirKey         string
}

func ConfigInit() (Config, error) {
	env, err := godotenv.Read()
	if err != nil {
		return Config{}, err
	}

	hashCost, _ := strconv.Atoi(env["HASH_COST"])
	if err != nil {
		return Config{}, err
	}

	expHour, err := strconv.Atoi(env["EXP_HOUR"])
	if err != nil {
		return Config{}, err
	}

	refreshExpDuration, err := strconv.Atoi(env["REFRESH_EXP_HOUR"])
	if err != nil {
		return Config{}, err
	}

	resetPasswordTokenExp, err := strconv.Atoi(env["RESET_PASSWORD_TOKEN_EXP"])
	if err != nil {
		return Config{}, err
	}

	return Config{
		DbUrl:                 env["DATABASE_URL"],
		Port:                  env["PORT"],
		HashCost:              hashCost,
		Issuer:                env["ISSUER"],
		ExpDurationHour:       expHour,
		RefreshExpDuration:    refreshExpDuration,
		SecretKey:             env["SECRET_KEY"],
		Email:                 env["EMAIL_SENDER"],
		EmailPassword:         env["EMAIL_PASSWORD"],
		SmtpHost:              env["SMTP_HOST"],
		SmtpPort:              env["SMTP_PORT"],
		FrontendUrl:           env["FRONTEND_URL"],
		CloudinaryName:        env["CLOUDINARY_NAME"],
		CloudinaryKey:         env["CLOUDINARY_KEY"],
		CloudinarySecret:      env["CLOUDINARY_SECRET"],
		GoogleId:              env["GOOGLE_ID"],
		GoogleKey:             env["GOOGLE_KEY"],
		GoogleUri:             env["GOOGLE_URI"],
		ResetTokenExpDuration: resetPasswordTokenExp,
		RajaOngkirKey:         env["RAJA_ONGKIR_KEY"],
	}, nil
}

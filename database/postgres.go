package database

import (
	"database/sql"

	"github.com/tsanaativa/sehatin-backend-v0.1/utils"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectDB(config utils.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.DbUrl)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

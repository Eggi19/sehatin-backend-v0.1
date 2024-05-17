package repositories

import (
	"context"
	"database/sql"

	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
)

type GenderRepoOpts struct {
	Db *sql.DB
}

type GenderRepository interface {
	FindById(ctx context.Context, genderId int64) (*entities.Gender, error)
}

type GenderRepositoryPostgres struct {
	db *sql.DB
}

func NewGenderRepositoryPostgres(gOpts *GenderRepoOpts) GenderRepository {
	return &GenderRepositoryPostgres{
		db: gOpts.Db,
	}
}

func (r *GenderRepositoryPostgres) FindById(ctx context.Context, genderId int64) (*entities.Gender, error) {
	g := entities.Gender{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindGenderById, genderId).Scan(&g.Id, &g.Name)
	} else {
		err = r.db.QueryRowContext(ctx, qFindGenderById, genderId).Scan(&g.Id, &g.Name)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &g, nil
}

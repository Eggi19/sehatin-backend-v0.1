package repositories

import (
	"context"
	"database/sql"

	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
)

type SpecialistRepoOpts struct {
	Db *sql.DB
}

type SpecialistRepository interface {
	FindAll(ctx context.Context) ([]entities.DoctorSpecialist, error)
}

type SpecialistRepositoryPostgres struct {
	db *sql.DB
}

func NewSpecialistRepositoryPostgres(spOpt *SpecialistRepoOpts) SpecialistRepository {
	return &SpecialistRepositoryPostgres{
		db: spOpt.Db,
	}
}

func (r *SpecialistRepositoryPostgres) FindAll(ctx context.Context) ([]entities.DoctorSpecialist, error) {
	sps := []entities.DoctorSpecialist{}

	var err error
	var rows *sql.Rows

	tx := extractTx(ctx)
	if tx != nil {
		rows, err = tx.QueryContext(ctx, qFindAllSpecialist)
	} else {
		rows, err = r.db.QueryContext(ctx, qFindAllSpecialist)
	}
	defer rows.Close()

	for rows.Next() {
		sp := entities.DoctorSpecialist{}
		rows.Scan(&sp.Id, &sp.Name)
		sps = append(sps, sp)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return sps, nil
}

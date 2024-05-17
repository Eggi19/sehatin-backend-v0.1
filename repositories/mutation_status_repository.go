package repositories

import (
	"context"
	"database/sql"

	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
)

type MutationStatusRepoOpts struct {
	Db *sql.DB
}

type MutationStatusRepository interface {
	FindAll(ctx context.Context) ([]entities.MutationSatus, error)
	FindOneById(ctx context.Context, id int64) (*entities.MutationSatus, error)
}

type MutationStatusRepositoryPostgres struct {
	db *sql.DB
}

func NewMutationStatusRepositoryPostgres(msrOpts *MutationStatusRepoOpts) MutationStatusRepository {
	return &MutationStatusRepositoryPostgres{
		db: msrOpts.Db,
	}
}

func (r *MutationStatusRepositoryPostgres) FindAll(ctx context.Context) ([]entities.MutationSatus, error) {
	mutations := []entities.MutationSatus{}

	var err error
	var rows *sql.Rows

	tx := extractTx(ctx)
	if tx != nil {
		rows, err = tx.QueryContext(ctx, qFindAllMutationStatus)
	} else {
		rows, err = r.db.QueryContext(ctx, qFindAllMutationStatus)
	}
	defer rows.Close()

	for rows.Next() {
		mutation := entities.MutationSatus{}
		rows.Scan(&mutation.Id, &mutation.Name)
		mutations = append(mutations, mutation)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return mutations, nil
}

func (r *MutationStatusRepositoryPostgres) FindOneById(ctx context.Context, id int64) (*entities.MutationSatus, error) {
	mutation := entities.MutationSatus{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindOneUserById, id).Scan(&mutation.Id, &mutation.Name)
	} else {
		err = r.db.QueryRowContext(ctx, qFindOneUserById, id).Scan(&mutation.Id, &mutation.Name)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &mutation, nil
}

package repositories

import (
	"context"
	"database/sql"

	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
)

type ProductFieldRepoOpts struct {
	Db *sql.DB
}

type ProductFieldRepository interface {
	FindAllForm(ctx context.Context) ([]entities.ProductForm, error)
	FindAllClassification(ctx context.Context) ([]entities.ProductClassification, error)
	FindAllManufacture(ctx context.Context) ([]entities.Manufacture, error)
}

type ProductFieldRepositoryPostgres struct {
	db *sql.DB
}

func NewProductFieldRepositoryPostgres(pcOpts *ProductCategoryRepoOpts) ProductFieldRepository {
	return &ProductFieldRepositoryPostgres{db: pcOpts.Db}
}

func (r *ProductFieldRepositoryPostgres) FindAllForm(ctx context.Context) ([]entities.ProductForm, error) {
	pfs := []entities.ProductForm{}

	var err error
	var rows *sql.Rows

	tx := extractTx(ctx)
	if tx != nil {
		rows, err = tx.QueryContext(ctx, qFindAllProductForm)
	} else {
		rows, err = r.db.QueryContext(ctx, qFindAllProductForm)
	}
	defer rows.Close()

	for rows.Next() {
		pf := entities.ProductForm{}
		rows.Scan(&pf.Id, &pf.Name)
		pfs = append(pfs, pf)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return pfs, nil
}

func (r *ProductFieldRepositoryPostgres) FindAllClassification(ctx context.Context) ([]entities.ProductClassification, error) {
	pcs := []entities.ProductClassification{}

	var err error
	var rows *sql.Rows

	tx := extractTx(ctx)
	if tx != nil {
		rows, err = tx.QueryContext(ctx, qFindAllClassification)
	} else {
		rows, err = r.db.QueryContext(ctx, qFindAllClassification)
	}
	defer rows.Close()

	for rows.Next() {
		pc := entities.ProductClassification{}
		rows.Scan(&pc.Id, &pc.Name)
		pcs = append(pcs, pc)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return pcs, nil
}

func (r *ProductFieldRepositoryPostgres) FindAllManufacture(ctx context.Context) ([]entities.Manufacture, error) {
	ms := []entities.Manufacture{}

	var err error
	var rows *sql.Rows

	tx := extractTx(ctx)
	if tx != nil {
		rows, err = tx.QueryContext(ctx, qFindAllManufacture)
	} else {
		rows, err = r.db.QueryContext(ctx, qFindAllManufacture)
	}
	defer rows.Close()

	for rows.Next() {
		m := entities.Manufacture{}
		rows.Scan(&m.Id, &m.Name)
		ms = append(ms, m)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return ms, nil
}

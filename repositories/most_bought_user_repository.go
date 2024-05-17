package repositories

import (
	"context"
	"database/sql"

	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
)

type MostBoughtUserRepoOpts struct {
	Db *sql.DB
}
type MostBoughtUserRepository interface {
	FindMostBought(ctx context.Context, id int64, pagination entities.PaginationParams, req entities.NearestPharmacyParams) ([]dtos.ProductResponse, error)
}

type MostBoughtUserRepositoryPostgres struct {
	db *sql.DB
}

func NewMostBoughtUserRepositoryPostgres(mbOpts *MostBoughtUserRepoOpts) MostBoughtUserRepository {
	return &MostBoughtUserRepositoryPostgres{
		db: mbOpts.Db,
	}
}

func (r *MostBoughtUserRepositoryPostgres) FindMostBought(ctx context.Context, id int64, pagination entities.PaginationParams, req entities.NearestPharmacyParams) ([]dtos.ProductResponse, error) {
	products := []dtos.ProductResponse{}

	var err error
	var rows *sql.Rows
	offset := pagination.Limit * (pagination.Page - 1)

	tx := extractTx(ctx)
	if tx != nil {
		rows, err = tx.QueryContext(ctx, qFindMostBought, id, offset, pagination.Limit)
	} else {
		rows, err = r.db.QueryContext(ctx, qFindMostBought, id, offset, pagination.Limit)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		product := dtos.ProductResponse{}
		err := rows.Scan(&product.PharmacyProductId, &product.Price, &product.SellingUnit, &product.SlugId, &product.ProductId, &product.Name, &product.ProductPicture, &product.Day, &product.QuantitySold, &product.Total)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return products, nil
}

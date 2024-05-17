package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/jackc/pgx/v5/pgconn"
)

type ProductCategoryRepoOpts struct {
	Db *sql.DB
}

type ProductCategoryRepository interface {
	CreateOne(ctx context.Context, pc entities.ProductCategory) error
	UpdateOne(ctx context.Context, pc entities.ProductCategory) error
	DeleteOneByProductId(ctx context.Context, productId int64) error
	FindAllByProductId(ctx context.Context, productItd int64) ([]entities.ProductCategory, error)
}

type ProductCategoryRepositoryPostgres struct {
	db *sql.DB
}

func NewProductCategoryRepositoryPostgres(pcOpts *ProductCategoryRepoOpts) ProductCategoryRepository {
	return &ProductCategoryRepositoryPostgres{
		db: pcOpts.Db,
	}
}

func (r *ProductCategoryRepositoryPostgres) CreateOne(ctx context.Context, pc entities.ProductCategory) error {
	newPc := entities.ProductCategory{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qCreateOneProductCategory, pc.ProductId, pc.CategoryId).Scan(&newPc.Id)
	} else {
		err = r.db.QueryRowContext(ctx, qCreateOneProductCategory, pc.ProductId, pc.CategoryId).Scan(&newPc.Id)
	}

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == constants.ViolatesUniqueConstraintPgErrCode {
			return custom_errors.BadRequest(err, constants.ProductNotUniqueErrMg)
		}
		return err
	}

	return err
}

func (r *ProductCategoryRepositoryPostgres) UpdateOne(ctx context.Context, pc entities.ProductCategory) error {

	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qUpdateOneProductCategoryByProductId)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qUpdateOneProductCategoryByProductId)
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, pc.ProductId, pc.CategoryId)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return custom_errors.NotFound(sql.ErrNoRows)
	}

	return nil
}

func (r *ProductCategoryRepositoryPostgres) DeleteOneByProductId(ctx context.Context, productId int64) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qDeleteProductCategoryByProductId)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qDeleteProductCategoryByProductId)
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, productId)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return custom_errors.NotFound(sql.ErrNoRows)
	}

	return nil
}

func (r *ProductCategoryRepositoryPostgres) FindAllByProductId(ctx context.Context, productId int64) ([]entities.ProductCategory, error) {
	pcs := []entities.ProductCategory{}

	var err error
	var rows *sql.Rows

	tx := extractTx(ctx)
	if tx != nil {
		rows, err = tx.QueryContext(ctx, qFindAllProductCategoryByProductId, productId)
	} else {
		rows, err = r.db.QueryContext(ctx, qFindAllProductCategoryByProductId, productId)
	}
	defer rows.Close()

	for rows.Next() {
		pc := entities.ProductCategory{}
		rows.Scan(&pc.Id, &pc.ProductId, pc.CategoryId)
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

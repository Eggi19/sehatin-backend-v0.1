package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/jackc/pgx/v5/pgconn"
)

type ProductRepoOpts struct {
	Db *sql.DB
}

type ProductRepository interface {
	CreateOne(ctx context.Context, product entities.Product) (*int64, error)
	UpdateOne(ctx context.Context, product entities.Product) error
	DeleteOne(ctx context.Context, productId int64) error
	FindOneById(ctx context.Context, productId int64) (*entities.Product, error)
	FindAll(ctx context.Context, params entities.ProductCategoryParams) ([]entities.Product, int, error)
}

type ProductRepositoryPostgres struct {
	db *sql.DB
}

func NewProductRepositoryPostgres(pOpts *ProductRepoOpts) ProductRepository {
	return &ProductRepositoryPostgres{
		db: pOpts.Db,
	}
}

func (r *ProductRepositoryPostgres) CreateOne(ctx context.Context, product entities.Product) (*int64, error) {
	newProduct := entities.Product{}

	values := []interface{}{}
	values = append(values, product.Name)
	values = append(values, product.GenericName)
	values = append(values, product.Content)
	values = append(values, product.Description)
	values = append(values, product.UnitInPack)
	values = append(values, product.SellingUnit)
	values = append(values, product.Weight)
	values = append(values, product.Height)
	values = append(values, product.Length)
	values = append(values, product.Width)
	values = append(values, product.ProductPicture)
	values = append(values, product.SlugId)
	values = append(values, product.ProductForm.Id)
	values = append(values, product.ProductClassification.Id)
	values = append(values, product.Manufacture.Id)

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qCreateOneProduct, values...).Scan(&newProduct.Id)
	} else {
		err = r.db.QueryRowContext(ctx, qCreateOneProduct, values...).Scan(&newProduct.Id)
	}

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == constants.ViolatesUniqueConstraintPgErrCode {
			return nil, custom_errors.BadRequest(err, constants.ProductNotUniqueErrMg)
		}
		return nil, err
	}

	return &newProduct.Id, nil
}

func (r *ProductRepositoryPostgres) UpdateOne(ctx context.Context, product entities.Product) error {
	values := []interface{}{}
	values = append(values, product.Id)
	values = append(values, product.Name)
	values = append(values, product.GenericName)
	values = append(values, product.Content)
	values = append(values, product.Description)
	values = append(values, product.UnitInPack)
	values = append(values, product.SellingUnit)
	values = append(values, product.Weight)
	values = append(values, product.Height)
	values = append(values, product.Length)
	values = append(values, product.Width)
	values = append(values, product.SlugId)
	values = append(values, product.ProductForm.Id)
	values = append(values, product.ProductClassification.Id)
	values = append(values, product.Manufacture.Id)

	var err error
	var stmt *sql.Stmt

	numberOfArgs := 16

	var sb strings.Builder
	sb.WriteString(qUpdateOneProductColl)
	if product.ProductPicture != "" {
		sb.WriteString(`product_picture = `)
		sb.WriteString(fmt.Sprintf(`$%d ,`, numberOfArgs))
		values = append(values, product.ProductPicture)
		numberOfArgs++
	}
	sb.WriteString(qUpdateOneProductCommand)

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, sb.String())
	} else {
		stmt, err = r.db.PrepareContext(ctx, sb.String())
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, values...)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == constants.ViolatesUniqueConstraintPgErrCode {
			return custom_errors.BadRequest(err, constants.ProductNotUniqueErrMg)
		}
		return err
	}
	if rows == 0 {
		return custom_errors.NotFound(sql.ErrNoRows)
	}

	return nil
}

func (r *ProductRepositoryPostgres) DeleteOne(ctx context.Context, productId int64) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qDeleteOneProduct)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qDeleteOneProduct)
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

func (r *ProductRepositoryPostgres) FindOneById(ctx context.Context, productId int64) (*entities.Product, error) {
	pc := entities.Product{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindOneProductById, productId).Scan(&pc.Id, &pc.Name, &pc.GenericName, &pc.Content, &pc.Description, &pc.UnitInPack,
			&pc.SellingUnit, &pc.Weight, &pc.Height, &pc.Length, &pc.Width, &pc.ProductPicture, &pc.SlugId, &pc.ProductForm.Id, &pc.ProductForm.Name, &pc.ProductClassification.Id, &pc.ProductClassification.Name, &pc.Manufacture.Id, &pc.Manufacture.Name)
	} else {
		err = r.db.QueryRowContext(ctx, qFindOneProductById, productId).Scan(&pc.Id, &pc.Name, &pc.GenericName, &pc.Content, &pc.Description, &pc.UnitInPack,
			&pc.SellingUnit, &pc.Weight, &pc.Height, &pc.Length, &pc.Width, &pc.ProductPicture, &pc.SlugId, &pc.ProductForm.Id, &pc.ProductForm.Name, &pc.ProductClassification.Id, &pc.ProductClassification.Name, &pc.Manufacture.Id, &pc.Manufacture.Name)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &pc, nil
}

func (r *ProductRepositoryPostgres) FindAll(ctx context.Context, params entities.ProductCategoryParams) ([]entities.Product, int, error) {
	products := []entities.Product{}

	var totalRows int

	var sb strings.Builder
	sb.WriteString(qCountTotalRows)
	sb.WriteString(qProductColl)
	sb.WriteString(qProductCommands)

	var sbTotalRows strings.Builder
	sbTotalRows.WriteString(qCountTotalRows)
	sbTotalRows.WriteString(qProductCommands)

	values := []interface{}{}
	valuesCountTotal := []interface{}{}

	numberOfArgs := 1

	if params.Keyword != "" {
		sb.WriteString(`AND p.name ILIKE '%`)
		sb.WriteString(params.Keyword)
		sb.WriteString(`%' `)

		sbTotalRows.WriteString(`AND p.name ILIKE '%`)
		sbTotalRows.WriteString(params.Keyword)
		sbTotalRows.WriteString(`%' `)
	}

	var sortBy string
	switch params.SortBy {
	case "p.name":
		sortBy = `p.name `
	default:
		sortBy = `p.id `
	}
	sb.WriteString(fmt.Sprintf(`ORDER BY %s `, sortBy))
	sbTotalRows.WriteString(fmt.Sprintf(`ORDER BY %s `, sortBy))

	if params.Sort == "" {
		params.Sort = `ASC `
	}
	sb.WriteString(fmt.Sprintf(`%s `, params.Sort))
	sbTotalRows.WriteString(fmt.Sprintf(`%s `, params.Sort))

	if params.Limit != 0 {
		sb.WriteString(`LIMIT `)
		sb.WriteString(fmt.Sprintf(`$%d `, numberOfArgs))
		values = append(values, params.Limit)
		numberOfArgs++
	}

	if params.Page != 0 {
		sb.WriteString(`OFFSET `)
		sb.WriteString(fmt.Sprintf(`$%d`, numberOfArgs))
		values = append(values, params.Limit*(params.Page-1))
		numberOfArgs++
	}

	rows, err := r.db.QueryContext(ctx, sb.String(), values...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		pc := entities.Product{}
		err = rows.Scan(&totalRows,
			&pc.Id, &pc.Name, &pc.GenericName, &pc.Content, &pc.Description, &pc.UnitInPack, &pc.SellingUnit, &pc.Weight, &pc.Height, &pc.Length,
			&pc.Width, &pc.ProductPicture, &pc.SlugId, &pc.ProductForm.Id, &pc.ProductForm.Name, &pc.ProductClassification.Id, &pc.ProductClassification.Name, &pc.Manufacture.Id, &pc.Manufacture.Name)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, pc)
	}

	err = rows.Err()
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	if totalRows == 0 {
		err := r.db.QueryRowContext(ctx, sbTotalRows.String(), valuesCountTotal...).Scan(&totalRows)
		if err != nil && err != sql.ErrNoRows {
			return nil, 0, err
		}
	}

	return products, totalRows, nil
}

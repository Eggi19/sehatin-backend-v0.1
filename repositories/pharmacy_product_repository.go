package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/jackc/pgx/v5/pgconn"
)

type PharmacyProductRepoOpt struct {
	Db *sql.DB
}

type PharmacyProductRepository interface {
	GetPharmacyProduct(ctx context.Context, id int64, pagination entities.PaginationParams, req entities.NearestPharmacyParams) ([]dtos.ProductResponse, error)
	GetProductDetail(ctx context.Context, pharmacyProductId int64) (*dtos.ProductDetail, error)
	GetOnePharmacyProduct(ctx context.Context, id int64) (*entities.PharmacyProduct, error)
	CreateOnePharmacyProduct(ctx context.Context, pp entities.PharmacyProduct) (*int64, error)
	UpdateOnePharmacyProduct(ctx context.Context, pp entities.PharmacyProduct) error
	DeleteOnePharmacyProduct(ctx context.Context, pharmacyProductId int64) error
	FindOneByPharmacyAndProductId(ctx context.Context, pharmacyId int64, productId int64) (*entities.PharmacyProduct, error)
	DeletedAllPharmacyProduct(ctx context.Context, pharmacyId int64) error
	FindPharmacyProductsByPharmacyId(ctx context.Context, pharmacyId, pharmacyManagerId int64, params entities.PharmacyProductParams) ([]entities.PharmacyProduct, int, error)
	FindAllNearest(ctx context.Context, params entities.NearestPharmacyProductsParams) ([]dtos.ProductResponse, int, error)
	UpdateTotalStock(ctx context.Context, pharmacyId, productId int64, totalStock int) error
	FindOnePharmacyProduct(ctx context.Context, id int64) (*entities.PharmacyProduct, error)
	IncreaseStock(ctx context.Context, quantity int, pharmacyProductId int64) error
	DecreaseStock(ctx context.Context, quantity int, pharmacyProductId int64) error
	LockRow(ctx context.Context, pharmacyProductId int64) error
	FindPharmacyProductByPharmacyId(ctx context.Context, pharmacyId int64) ([]entities.PharmacyProduct, error)
}

type PharmacyProductRepositoryPostgres struct {
	db *sql.DB
}

func (r *PharmacyProductRepositoryPostgres) FindPharmacyProductByPharmacyId(ctx context.Context, pharmacyId int64) ([]entities.PharmacyProduct, error) {
	products := []entities.PharmacyProduct{}

	var err error
	var rows *sql.Rows

	tx := extractTx(ctx)
	if tx != nil {
		rows, err = tx.QueryContext(ctx, qFindPharmacyProductByPharmacyId, pharmacyId)
	} else {
		rows, err = r.db.QueryContext(ctx, qFindPharmacyProductByPharmacyId, pharmacyId)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		pp := entities.PharmacyProduct{}
		err := rows.Scan(&pp.Product.ProductPicture, &pp.Product.Name, &pp.Price, &pp.Product.SellingUnit, &pp.Product.SlugId, &pp.Id, &pp.Product.Id)
		if err != nil {
			return nil, err
		}

		products = append(products, pp)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return products, nil
}

func NewPharmacyProductRepositoryPostgres(pprOpt *PharmacyProductRepoOpt) PharmacyProductRepository {
	return &PharmacyProductRepositoryPostgres{
		db: pprOpt.Db,
	}
}

func (r *PharmacyProductRepositoryPostgres) GetPharmacyProduct(ctx context.Context, id int64, pagination entities.PaginationParams, req entities.NearestPharmacyParams) ([]dtos.ProductResponse, error) {
	products := []dtos.ProductResponse{}

	var err error
	var rows *sql.Rows
	offset := pagination.Limit * (pagination.Page - 1)

	tx := extractTx(ctx)
	if tx != nil {
		if req.CategoryId != 0 {
			rows, err = tx.QueryContext(ctx, qFindPharmacyProductByCategory, id, offset, pagination.Limit, req.CategoryId)
		} else {
			rows, err = tx.QueryContext(ctx, qFindPharmacyProduct, id, offset, pagination.Limit)
		}
	} else {
		if req.CategoryId != 0 {
			rows, err = tx.QueryContext(ctx, qFindPharmacyProductByCategory, id, offset, pagination.Limit, req.CategoryId)
		} else {
			rows, err = r.db.QueryContext(ctx, qFindPharmacyProduct, id, offset, pagination.Limit)
		}
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		product := dtos.ProductResponse{}
		err := rows.Scan(&product.ProductPicture, &product.Name, &product.Price, &product.SellingUnit, &product.SlugId, &product.PharmacyProductId, &product.ProductId, &product.Total)
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

func (r *PharmacyProductRepositoryPostgres) GetProductDetail(ctx context.Context, pharmacyProductId int64) (*dtos.ProductDetail, error) {
	productDetail := dtos.ProductDetail{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindPharmacyProductDetail, pharmacyProductId).Scan(&productDetail.Id, &productDetail.Name, &productDetail.GenericName, &productDetail.Content, &productDetail.Description, &productDetail.UnitInPack, &productDetail.SellingUnit, &productDetail.Weight, &productDetail.Height, &productDetail.Length, &productDetail.Width, &productDetail.ProductPicture, &productDetail.SlugId, &productDetail.ProductForm, &productDetail.ProductClassification, &productDetail.Manufacture, &productDetail.Price, &productDetail.TotalStock)
	} else {
		err = r.db.QueryRowContext(ctx, qFindPharmacyProductDetail, pharmacyProductId).Scan(&productDetail.Id, &productDetail.Name, &productDetail.GenericName, &productDetail.Content, &productDetail.Description, &productDetail.UnitInPack, &productDetail.SellingUnit, &productDetail.Weight, &productDetail.Height, &productDetail.Length, &productDetail.Width, &productDetail.ProductPicture, &productDetail.SlugId, &productDetail.ProductForm, &productDetail.ProductClassification, &productDetail.Manufacture, &productDetail.Price, &productDetail.TotalStock)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &productDetail, nil
}

func (r *PharmacyProductRepositoryPostgres) CreateOnePharmacyProduct(ctx context.Context, pp entities.PharmacyProduct) (*int64, error) {
	newPp := entities.PharmacyProduct{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qCreatePharmacyProduct, pp.Price, pp.TotalStock, pp.IsAvailable, pp.Product.Id, pp.Pharmacy.Id).Scan(&newPp.Id)
	} else {
		err = r.db.QueryRowContext(ctx, qCreatePharmacyProduct, pp.Price, pp.TotalStock, pp.IsAvailable, pp.Product.Id, pp.Pharmacy.Id).Scan(&newPp.Id)
	}

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == constants.ViolatesUniqueConstraintPgErrCode {
			return nil, custom_errors.BadRequest(err, constants.ProductNotUniqueErrMg)
		}
		return nil, err
	}

	return &newPp.Id, nil
}

func (r *PharmacyProductRepositoryPostgres) UpdateOnePharmacyProduct(ctx context.Context, pp entities.PharmacyProduct) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qUpdatePharmacyProduct)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qUpdatePharmacyProduct)
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, pp.Id, pp.Price, pp.TotalStock, pp.IsAvailable, pp.Product.Id, pp.Pharmacy.Id)
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

func (r *PharmacyProductRepositoryPostgres) DeleteOnePharmacyProduct(ctx context.Context, pharmacyProductId int64) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qDeletePharmacyProduct)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qDeletePharmacyProduct)
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, pharmacyProductId)
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

func (r *PharmacyProductRepositoryPostgres) GetOnePharmacyProduct(ctx context.Context, id int64) (*entities.PharmacyProduct, error) {
	pp := entities.PharmacyProduct{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindOnePharmacyProduct, id).Scan(
			&pp.Id, &pp.Pharmacy.Id, &pp.Product.Id, &pp.TotalStock,
		)
	} else {
		err = r.db.QueryRowContext(ctx, qFindOnePharmacyProduct, id).Scan(
			&pp.Id, &pp.Pharmacy.Id, &pp.Product.Id, &pp.TotalStock,
		)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &pp, nil
}

func (r *PharmacyProductRepositoryPostgres) FindOneByPharmacyAndProductId(ctx context.Context, pharmacyId int64, productId int64) (*entities.PharmacyProduct, error) {
	pp := entities.PharmacyProduct{
		Product: entities.Product{},
		Pharmacy: entities.Pharmacy{
			PharmacyManager: entities.PharmacyManager{},
		},
	}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindOneByPharmacyAndProductId, pharmacyId, productId).Scan(
			&pp.Id, &pp.TotalStock, &pp.IsAvailable, &pp.Pharmacy.Id, &pp.Pharmacy.Name, &pp.Product.Id, &pp.Product.Name, &pp.Pharmacy.PharmacyManager.Id, &pp.DeletedAt,
		)
	} else {
		err = r.db.QueryRowContext(ctx, qFindOneByPharmacyAndProductId, pharmacyId, productId).Scan(
			&pp.Id, &pp.TotalStock, &pp.IsAvailable, &pp.Pharmacy.Id, &pp.Pharmacy.Name, &pp.Product.Id, &pp.Product.Name, &pp.Pharmacy.PharmacyManager.Id, &pp.DeletedAt,
		)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &pp, nil
}

func (r *PharmacyProductRepositoryPostgres) DeletedAllPharmacyProduct(ctx context.Context, pharmacyId int64) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qDeletedAllPharmacyProduct)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qDeletedAllPharmacyProduct)
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, pharmacyId)
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

func (r *PharmacyProductRepositoryPostgres) FindPharmacyProductsByPharmacyId(ctx context.Context, pharmacyId, pharmacyManagerId int64, params entities.PharmacyProductParams) ([]entities.PharmacyProduct, int, error) {
	pharmacyProducts := []entities.PharmacyProduct{}

	var totalRows int

	var sb strings.Builder
	sb.WriteString(qCountTotalRows)
	sb.WriteString(qPharmacyProductByPharmacyIdColl)
	sb.WriteString(qPharmacyProductByPharmacyIdCommand)

	var sbTotalRows strings.Builder
	sbTotalRows.WriteString(qCountTotalRows)
	sbTotalRows.WriteString(qPharmacyProductByPharmacyIdCommand)

	values := []interface{}{}
	valuesCountTotal := []interface{}{}

	numberOfArgs := 3
	values = append(values, pharmacyId)
	valuesCountTotal = append(valuesCountTotal, pharmacyId)
	values = append(values, pharmacyManagerId)
	valuesCountTotal = append(valuesCountTotal, pharmacyManagerId)

	if params.Keyword != "" {
		sb.WriteString(`AND pd.name ILIKE '%`)
		sb.WriteString(params.Keyword)
		sb.WriteString(`%' `)

		sbTotalRows.WriteString(`AND pd.name ILIKE '%`)
		sbTotalRows.WriteString(params.Keyword)
		sbTotalRows.WriteString(`%' `)
	}

	var sortBy string
	switch params.SortBy {
	case "product-name":
		sortBy = `pd.name `
	case "price":
		sortBy = `pp.price `
	default:
		sortBy = `pp.id `
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
		pp := entities.PharmacyProduct{}
		err := rows.Scan(&totalRows,
			&pp.Id, &pp.TotalStock, &pp.IsAvailable, &pp.Price,
			&pp.Product.Id, &pp.Product.Name, &pp.Product.Content, &pp.Product.Description, &pp.Product.UnitInPack, &pp.Product.SellingUnit,
			&pp.Product.Weight, &pp.Product.Height, &pp.Product.Length, &pp.Product.Width, &pp.Product.ProductPicture, &pp.Product.SlugId, &pp.Product.ProductClassification.Name,
			&pp.Product.ProductForm.Name, &pp.Product.Manufacture.Name,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, 0, custom_errors.NotFound(err)
			}
			return nil, 0, err
		}
		pharmacyProducts = append(pharmacyProducts, pp)
	}

	if totalRows == 0 {
		err := r.db.QueryRowContext(ctx, sbTotalRows.String(), valuesCountTotal...).Scan(&totalRows)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, 0, custom_errors.NotFound(err)
			}
			return nil, 0, err
		}
	}

	return pharmacyProducts, totalRows, nil
}

func (r *PharmacyProductRepositoryPostgres) FindAllNearest(ctx context.Context, params entities.NearestPharmacyProductsParams) ([]dtos.ProductResponse, int, error) {
	pharmacyProducts := []dtos.ProductResponse{}

	var totalRows int

	var sb strings.Builder
	sb.WriteString(qCountTotalRows)
	sb.WriteString(qNearestPharmacyProductColl)
	sb.WriteString(qNearestPharmacyProductCommands)

	var sbTotalRows strings.Builder
	sbTotalRows.WriteString(qCountTotalRows)
	sbTotalRows.WriteString(qNearestPharmacyProductCommands)

	values := []interface{}{}
	values = append(values, params.Longitude)
	values = append(values, params.Latitude)
	values = append(values, params.Radius)
	valuesCountTotal := []interface{}{}
	valuesCountTotal = append(valuesCountTotal, params.Longitude)
	valuesCountTotal = append(valuesCountTotal, params.Latitude)
	valuesCountTotal = append(valuesCountTotal, params.Radius)

	numberOfArgs := 4

	if params.CategoryId != 0 {
		sb.WriteString(`AND pc.category_id = $4`)
		sbTotalRows.WriteString(`AND pc.category_id = $4`)

		values = append(values, params.CategoryId)
		valuesCountTotal = append(valuesCountTotal, params.CategoryId)
		numberOfArgs++
	}

	sb.WriteString(qNearestPharmacyProductCommandsSecond)
	sbTotalRows.WriteString(qNearestPharmacyProductCommandsSecond)

	if params.Keyword != "" {
		sb.WriteString(`WHERE name ILIKE '%`)
		sb.WriteString(params.Keyword)
		sb.WriteString(`%' `)

		sbTotalRows.WriteString(`WHERE name ILIKE '%`)
		sbTotalRows.WriteString(params.Keyword)
		sbTotalRows.WriteString(`%' `)
	}

	var sortBy string
	switch params.SortBy {
	case "name":
		sortBy = `name `
	case "price":
		sortBy = `price `
	default:
		sortBy = `product_id `
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
		pp := dtos.ProductResponse{}

		err = rows.Scan(&totalRows,
			&pp.PharmacyProductId, &pp.ProductId, &pp.Name, &pp.Price, &pp.ProductPicture, &pp.SellingUnit, &pp.SlugId)

		if err != nil {
			return nil, 0, err
		}

		pharmacyProducts = append(pharmacyProducts, pp)
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

	return pharmacyProducts, totalRows, nil
}

func (r *PharmacyProductRepositoryPostgres) UpdateTotalStock(ctx context.Context, pharmacyId, productId int64, totalStock int) error {

	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qUpdatePharmacyProductStock)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qUpdatePharmacyProductStock)
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, pharmacyId, productId, totalStock)
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

func (r *PharmacyProductRepositoryPostgres) FindOnePharmacyProduct(ctx context.Context, id int64) (*entities.PharmacyProduct, error) {
	pp := entities.PharmacyProduct{
		Product: entities.Product{
			ProductForm:           entities.ProductForm{},
			ProductClassification: entities.ProductClassification{},
			Manufacture:           entities.Manufacture{},
			Categories:            []entities.Category{},
			CreatedAt:             time.Time{},
			UpdatedAt:             time.Time{},
			DeletedAt:             sql.NullTime{},
		},
		Pharmacy: entities.Pharmacy{
			PharmacyManager: entities.PharmacyManager{},
		},
	}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindOnePharmacyProductById, id).Scan(
			&pp.Id, &pp.TotalStock, &pp.IsAvailable, &pp.Price,
			&pp.Product.Id, &pp.Product.Name, &pp.Product.Content, &pp.Product.Description, &pp.Product.UnitInPack, &pp.Product.SellingUnit,
			&pp.Product.Weight, &pp.Product.Height, &pp.Product.Length, &pp.Product.Width, &pp.Product.ProductPicture, &pp.Product.SlugId, &pp.Product.ProductClassification.Name,
			&pp.Product.ProductForm.Name, &pp.Product.Manufacture.Name, &pp.DeletedAt,
		)
	} else {
		err = r.db.QueryRowContext(ctx, qFindOnePharmacyProductById, id).Scan(
			&pp.Id, &pp.TotalStock, &pp.IsAvailable, &pp.Price,
			&pp.Product.Id, &pp.Product.Name, &pp.Product.Content, &pp.Product.Description, &pp.Product.UnitInPack, &pp.Product.SellingUnit,
			&pp.Product.Weight, &pp.Product.Height, &pp.Product.Length, &pp.Product.Width, &pp.Product.ProductPicture, &pp.Product.SlugId, &pp.Product.ProductClassification.Name,
			&pp.Product.ProductForm.Name, &pp.Product.Manufacture.Name, &pp.DeletedAt,
		)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &pp, err
}

func (r *PharmacyProductRepositoryPostgres) IncreaseStock(ctx context.Context, quantity int, pharmacyProductId int64) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qIncreaseStock)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qIncreaseStock)
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, quantity, pharmacyProductId)
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

func (r *PharmacyProductRepositoryPostgres) DecreaseStock(ctx context.Context, quantity int, pharmacyProductId int64) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qDecreaseStock)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qDecreaseStock)
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, quantity, pharmacyProductId)
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

func (r *PharmacyProductRepositoryPostgres) LockRow(ctx context.Context, pharmacyProductId int64) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qLockPharmacyProductRow)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qLockPharmacyProductRow)
	}

	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, pharmacyProductId)
	if err != nil {
		return err
	}

	return nil
}

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

type StockHistoryRepoOpts struct {
	Db *sql.DB
}

type StockHistoryRepository interface {
	CreateOne(ctx context.Context, stockHistory entities.StockHistory) error
	UpdateOne(ctx context.Context, stockHistory entities.StockHistory) error
	FindOneByPharmacyIdAndPharmacyProductId(ctx context.Context, pharmacyId int64, pharmacyProductId int64) (*entities.StockHistory, error)
	DeleteAllByPharmacyId(ctx context.Context, pharmacyId int64) error
	FindStockHistoriesByPharmacyId(ctx context.Context, pharmacyId, pharmacyManagerId int64, params entities.StockHistoryParams) ([]entities.StockHistory, int, error)
}

type StockHistoryRepositoryPostgres struct {
	db *sql.DB
}

func NewStockHistoryRepositoryPostgres(shOpts *StockHistoryRepoOpts) StockHistoryRepository {
	return &StockHistoryRepositoryPostgres{
		db: shOpts.Db,
	}
}

func (r *StockHistoryRepositoryPostgres) CreateOne(ctx context.Context, stockHistory entities.StockHistory) error {
	newStockHistory := entities.StockHistory{}

	values := []interface{}{}
	values = append(values, stockHistory.PharmacyProduct.Id)
	values = append(values, stockHistory.Pharmacy.Id)
	values = append(values, stockHistory.Quantity)
	values = append(values, stockHistory.Description)

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qCreateOneStockHistory, values...).Scan(&newStockHistory.Id)
	} else {
		err = r.db.QueryRowContext(ctx, qCreateOneStockHistory, values...).Scan(&newStockHistory.Id)
	}

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == constants.ViolatesUniqueConstraintPgErrCode {
			return custom_errors.BadRequest(err, constants.PharmacyNotUniqueErrMg)
		}
		return err
	}

	return nil
}

func (r *StockHistoryRepositoryPostgres) UpdateOne(ctx context.Context, stockHistory entities.StockHistory) error {
	values := []interface{}{}
	values = append(values, stockHistory.PharmacyProduct.Id)
	values = append(values, stockHistory.Pharmacy.Id)
	values = append(values, stockHistory.PharmacyProduct.Id)
	values = append(values, stockHistory.Pharmacy.Id)
	values = append(values, stockHistory.Quantity)
	values = append(values, stockHistory.Description)

	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qUpdateStockHistory)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qUpdateStockHistory)
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
		return err
	}
	if rows == 0 {
		return custom_errors.NotFound(sql.ErrNoRows)
	}

	return nil
}

func (r *StockHistoryRepositoryPostgres) FindOneByPharmacyIdAndPharmacyProductId(ctx context.Context, pharmacyId int64, pharmacyProductId int64) (*entities.StockHistory, error) {
	sh := entities.StockHistory{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindStockHistoryByPharmacyIdAndPharmacyProductId, pharmacyId, pharmacyProductId).
			Scan(
				&sh.Id, &sh.Quantity, &sh.Description, &sh.PharmacyProduct.Id, &sh.PharmacyProduct.Price, &sh.PharmacyProduct.TotalStock,
				&sh.PharmacyProduct.Product.Id, &sh.PharmacyProduct.Product.Name,
				&sh.Pharmacy.Id, &sh.Pharmacy.Name, &sh.Pharmacy.PharmacyManager.Id,
			)
	} else {
		err = r.db.QueryRowContext(ctx, qFindStockHistoryByPharmacyIdAndPharmacyProductId, pharmacyId, pharmacyProductId).
			Scan(&sh.Id, &sh.Quantity, &sh.Description, &sh.PharmacyProduct.Id, &sh.PharmacyProduct.Price, &sh.PharmacyProduct.TotalStock,
				&sh.Pharmacy.Id, &sh.Pharmacy.Name,
			)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &sh, nil
}

func (r *StockHistoryRepositoryPostgres) DeleteAllByPharmacyId(ctx context.Context, pharmacyId int64) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qDeleteAllStockHistoryByPharmacyId)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qDeleteAllStockHistoryByPharmacyId)
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

func (r *StockHistoryRepositoryPostgres) FindStockHistoriesByPharmacyId(ctx context.Context, pharmacyId, pharmacyManagerId int64, params entities.StockHistoryParams) ([]entities.StockHistory, int, error) {
	shs := []entities.StockHistory{}

	var totalRows int

	var sb strings.Builder
	sb.WriteString(qCountTotalRows)
	sb.WriteString(qFindStockHistoriesByPharmacyIdColl)
	sb.WriteString(qFindStockHistoriesByPharmacyIdCommand)

	var sbTotalRows strings.Builder
	sbTotalRows.WriteString(qCountTotalRows)
	sbTotalRows.WriteString(qFindStockHistoriesByPharmacyIdCommand)

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
		sh := entities.StockHistory{}
		err := rows.Scan(&totalRows,
			&sh.Id, &sh.Quantity, &sh.Description, &sh.PharmacyProduct.Id, &sh.PharmacyProduct.Price,
			&sh.PharmacyProduct.Product.Id, &sh.PharmacyProduct.Product.Name,
			&sh.Pharmacy.Id, &sh.Pharmacy.Name, &sh.Pharmacy.PharmacyManager.Id,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, 0, custom_errors.NotFound(err)
			}
			return nil, 0, err
		}
		shs = append(shs, sh)
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

	return shs, totalRows, nil
}

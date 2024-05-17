package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
)

type StockHistoryReportOpts struct {
	Db *sql.DB
}

type StockHistoryReportRepository interface {
	FindAll(ctx context.Context, params entities.StockHistoryReportParams) ([]entities.StockHistoryReport, int, error)
}

type StockHistoryReportRepositoryPostgres struct {
	db *sql.DB
}

func NewStockHistoryReportRepositoryPostgres(shrOpts *StockHistoryReportOpts) StockHistoryReportRepository {
	return &StockHistoryReportRepositoryPostgres{
		db: shrOpts.Db,
	}
}

func (r *StockHistoryReportRepositoryPostgres) FindAll(ctx context.Context, params entities.StockHistoryReportParams) ([]entities.StockHistoryReport, int, error) {
	shs := []entities.StockHistoryReport{}

	var totalRows int

	var sb strings.Builder
	sb.WriteString(qCountTotalRows)
	sb.WriteString(qStockHistoryReportsColl)
	sb.WriteString(qStockHistoryReportsCommand)

	var sbTotalRows strings.Builder
	sbTotalRows.WriteString(qCountTotalRows)
	sbTotalRows.WriteString(qStockHistoryReportsCommand)

	values := []interface{}{}
	valuesCountTotal := []interface{}{}

	numberOfArgs := 1

	if params.Keyword != "" {
		sb.WriteString(`AND ph.name ILIKE '%`)
		sb.WriteString(params.Keyword)
		sb.WriteString(`%' `)

		sbTotalRows.WriteString(`AND ph.name ILIKE '%`)
		sbTotalRows.WriteString(params.Keyword)
		sbTotalRows.WriteString(`%' `)
	}

	if params.PharmacyId != 0 {
		sb.WriteString(`AND ph.id = `)
		sb.WriteString(fmt.Sprintf(`$%d `, numberOfArgs))
		values = append(values, params.PharmacyId)

		sbTotalRows.WriteString(`AND ph.id = `)
		sbTotalRows.WriteString(fmt.Sprintf(`$%d `, numberOfArgs))
		valuesCountTotal = append(valuesCountTotal, params.PharmacyId)

		numberOfArgs++
	}

	sb.WriteString(qStockHistoryReportsGroup)
	sbTotalRows.WriteString(qStockHistoryReportsGroup)

	var sortBy string
	switch params.SortBy {
	case "pharmacy":
		sortBy = `ph.name `
	default:
		sortBy = `month `
	}
	sb.WriteString(fmt.Sprintf(`ORDER BY %s `, sortBy))
	sbTotalRows.WriteString(fmt.Sprintf(`ORDER BY %s `, sortBy))

	if params.Sort == "" {
		params.Sort = `ASC `
	}
	sb.WriteString(fmt.Sprintf(`%s `, params.Sort))
	sbTotalRows.WriteString(fmt.Sprintf(`%s `, params.Sort))

	sb.WriteString(`LIMIT `)
	sb.WriteString(fmt.Sprintf(`$%d `, numberOfArgs))
	values = append(values, params.Limit)
	numberOfArgs++

	sb.WriteString(`OFFSET `)
	sb.WriteString(fmt.Sprintf(`$%d`, numberOfArgs))
	values = append(values, params.Limit*(params.Page-1))
	numberOfArgs++

	rows, err := r.db.QueryContext(ctx, sb.String(), values...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		sh := entities.StockHistoryReport{PharmacyProduct: entities.PharmacyProduct{Pharmacy: entities.Pharmacy{}, Product: entities.Product{}}}
		err := rows.Scan(&totalRows,
			&sh.TotalAddition, &sh.TotalDeduction, &sh.PharmacyProduct.TotalStock, &sh.PharmacyProduct.Product.Id, &sh.PharmacyProduct.Product.Name, &sh.PharmacyProduct.Pharmacy.Id, &sh.PharmacyProduct.Pharmacy.Name, &sh.Month, &sh.Year,
		)
		if err != nil {
			return nil, 0, err
		}
		shs = append(shs, sh)
	}

	if len(shs) == 0 {
		return nil, 0, custom_errors.BadRequest(err, constants.ResponseMsgErrorNotFound)
	}

	err = rows.Err()
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	if totalRows == 0 {
		err := r.db.QueryRowContext(ctx, sbTotalRows.String(), valuesCountTotal...).Scan(&totalRows)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, 0, nil
			}
			return nil, 0, err
		}
	}

	return shs, totalRows, nil
}

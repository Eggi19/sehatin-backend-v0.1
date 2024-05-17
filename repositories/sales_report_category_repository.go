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

type SalesReportCatgoryRepoOpts struct {
	Db *sql.DB
}
type SalesReportCategoryRepository interface {
	FindAll(ctx context.Context, params entities.SalesReportCategoryParams) ([]entities.SalesReportCategory, int, error)
}

type SalesReportCategoryRepositoryPostgres struct {
	db *sql.DB
}

func NewSalesReportCategoryRepositoryPostgres(srOpts *SalesReportCatgoryRepoOpts) SalesReportCategoryRepository {
	return &SalesReportCategoryRepositoryPostgres{
		db: srOpts.Db,
	}
}

func (r *SalesReportCategoryRepositoryPostgres) FindAll(ctx context.Context, params entities.SalesReportCategoryParams) ([]entities.SalesReportCategory, int, error) {
	src := []entities.SalesReportCategory{}

	var totalRows int

	var sb strings.Builder
	sb.WriteString(qCountTotalRows)
	sb.WriteString(qSalesReportCategoryColl)
	sb.WriteString(qSalesReportCategoryCommand)

	var sbTotalRows strings.Builder
	sbTotalRows.WriteString(qCountTotalRows)
	sbTotalRows.WriteString(qSalesReportCategoryCommand)

	values := []interface{}{}
	valuesCountTotal := []interface{}{}

	numberOfArgs := 1

	if params.CategoryId != 0 {
		sb.WriteString(`AND c.id = `)
		sb.WriteString(fmt.Sprintf(`$%d `, numberOfArgs))
		values = append(values, params.CategoryId)

		sbTotalRows.WriteString(`AND c.id = `)
		sbTotalRows.WriteString(fmt.Sprintf(`$%d `, numberOfArgs))
		valuesCountTotal = append(valuesCountTotal, params.CategoryId)

		numberOfArgs++
	}

	if params.Keyword != "" {
		sb.WriteString(`AND c.name ILIKE '%`)
		sb.WriteString(params.Keyword)
		sb.WriteString(`%' `)
		sbTotalRows.WriteString(`AND c.name ILIKE '%`)
		sbTotalRows.WriteString(params.Keyword)
		sbTotalRows.WriteString(`%' `)
	}

	sb.WriteString(qSalesReportCategoryGroup)
	sbTotalRows.WriteString(qSalesReportCategoryGroup)

	var sortBy string
	switch params.SortBy {
	case "totalSold":
		sortBy = `oi.quantity `
	case "total_sale":
		sortBy = `total_sales_amount `
	case "categoryName":
		sortBy = `c.Name `
	case "year":
		sortBy = `year `
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
		sr := entities.SalesReportCategory{Category: entities.Category{}}
		err := rows.Scan(&totalRows,
			&sr.Category.Id, &sr.Category.Name, &sr.Month, &sr.Year, &sr.TotalSold,
		)
		if err != nil {
			return nil, 0, err
		}
		src = append(src, sr)
	}

	if len(src) == 0 {
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

	return src, totalRows, nil

}

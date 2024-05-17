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

type SalesReportRepoOpts struct {
	Db *sql.DB
}
type SalesReportRepository interface {
	FindAll(ctx context.Context, params entities.SalesReportParams) ([]entities.SalesReport, int, error)
}

type SalesReportRepositoryPostgres struct {
	db *sql.DB
}

func NewSalesReportRepositoryPostgres(srOpts *SalesReportRepoOpts) SalesReportRepository {
	return &SalesReportRepositoryPostgres{
		db: srOpts.Db,
	}
}

func (r *SalesReportRepositoryPostgres) FindAll(ctx context.Context, params entities.SalesReportParams) ([]entities.SalesReport, int, error) {
	srs := []entities.SalesReport{}

	var totalRows int

	var sb strings.Builder
	sb.WriteString(qCountTotalRows)
	sb.WriteString(qSalesReportColl)
	sb.WriteString(qSalesReportCommand)

	var sbTotalRows strings.Builder
	sbTotalRows.WriteString(qCountTotalRows)
	sbTotalRows.WriteString(qSalesReportCommand)

	values := []interface{}{}
	valuesCountTotal := []interface{}{}

	numberOfArgs := 1

	if params.PharmacyId != 0 {
		sb.WriteString(`AND ph.id = `)
		sb.WriteString(fmt.Sprintf(`$%d `, numberOfArgs))
		values = append(values, params.PharmacyId)

		sbTotalRows.WriteString(`AND ph.id = `)
		sbTotalRows.WriteString(fmt.Sprintf(`$%d `, numberOfArgs))
		valuesCountTotal = append(valuesCountTotal, params.PharmacyId)

		numberOfArgs++
	}

	if params.ProductId != 0 {
		sb.WriteString(`AND pd.id = `)
		sb.WriteString(fmt.Sprintf(`$%d `, numberOfArgs))
		values = append(values, params.ProductId)

		sbTotalRows.WriteString(`AND pd.id = `)
		sbTotalRows.WriteString(fmt.Sprintf(`$%d `, numberOfArgs))
		valuesCountTotal = append(valuesCountTotal, params.ProductId)

		numberOfArgs++
	}

	if params.Keyword != "" {
		sb.WriteString(`AND ph.name ILIKE '%`)
		sb.WriteString(params.Keyword)
		sb.WriteString(`%' `)

		sbTotalRows.WriteString(`AND ph.name ILIKE '%`)
		sbTotalRows.WriteString(params.Keyword)
		sbTotalRows.WriteString(`%' `)
	}

	sb.WriteString(qSalesReportGroup)
	sbTotalRows.WriteString(qSalesReportGroup)

	var sortBy string
	switch params.SortBy {
	case "total-sale":
		sortBy = `total_sales_amount `
	case "total-quantity":
		sortBy = `total_quantity_sold`
	case "product":
		sortBy = `pd.name `
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
		sr := entities.SalesReport{PharmacyProduct: entities.PharmacyProduct{Pharmacy: entities.Pharmacy{}, Product: entities.Product{}}}
		err := rows.Scan(&totalRows,
			&sr.PharmacyProduct.Pharmacy.Id, &sr.PharmacyProduct.Pharmacy.Name,
			&sr.PharmacyProduct.Product.Id, &sr.PharmacyProduct.Product.Name, &sr.PharmacyProduct.Product.Content, &sr.PharmacyProduct.Product.Description, &sr.PharmacyProduct.Product.UnitInPack, &sr.PharmacyProduct.Product.SellingUnit,
			&sr.PharmacyProduct.Product.Weight, &sr.PharmacyProduct.Product.Height, &sr.PharmacyProduct.Product.Length, &sr.PharmacyProduct.Product.Width, &sr.PharmacyProduct.Product.ProductPicture, &sr.PharmacyProduct.Product.SlugId, &sr.PharmacyProduct.Product.ProductClassification.Name,
			&sr.PharmacyProduct.Product.ProductForm.Name, &sr.PharmacyProduct.Product.Manufacture.Name,
			&sr.TotalSalesAmount, &sr.TotalQuantitySold, &sr.Month, &sr.Year,
		)
		if err != nil {
			return nil, 0, err
		}
		srs = append(srs, sr)
	}

	if len(srs) == 0 {
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

	return srs, totalRows, nil
}
